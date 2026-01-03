package proxy

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// WebSocketFrame represents a single WS message
type WebSocketFrame struct {
	ID        string `json:"id"`
	Time      string `json:"time"`
	Direction string `json:"direction"` // "outgoing" (Client->Server) or "incoming" (Server->Client)
	OpCode    int    `json:"opcode"`    // 1=Text, 2=Binary, 8=Close, 9=Ping, 10=Pong
	Payload   string `json:"payload"`
	Length    int64  `json:"length"`
}

// HandleWebSocketTunnel performs the Man-in-the-Middle for WebSockets
// It assumes it's called when `Upgrade: websocket` is detected.
func HandleWebSocketTunnel(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// 1. Hijack Client Connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		fmt.Printf("WS: Hijack error: %v\n", err)
		return
	}
	defer clientConn.Close()

	// 2. Dial Target Server
	var targetConn net.Conn
	host := r.Host
	if !strings.Contains(host, ":") {
		if r.TLS != nil || r.URL.Scheme == "https" || r.URL.Scheme == "wss" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	// Determine if TLS
	// NOTE: If we are in the MitM proxy, r.URL.Scheme might be http because goproxy strips SSL.
	// However, we know if the *original* connect was port 443.
	// A safe bet: if goproxy gives us this request, we are looking at the deciphered text.
	// But we need to communicate to the upstream.
	// If the upstream is secure, we need TLS.
	// Simple heuristic: if port 443, use TLS.
	if strings.HasSuffix(host, ":443") {
		targetConn, err = tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	} else {
		targetConn, err = net.Dial("tcp", host)
	}

	if err != nil {
		fmt.Printf("WS: Failed to dial target %s: %v\n", host, err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer targetConn.Close()

	// 3. Reconstruct and Send Handshake to Target
	// We need to write the original request to the target to initiate the upgrade there.
	err = r.Write(targetConn)
	if err != nil {
		fmt.Printf("WS: Failed to write handshake: %v\n", err)
		return
	}

	// 4. Pipe and Snoop
	// We need to snoop on the response from Target (Handshake response)
	// Then snoop frames.

	// Because `r.Write` sends all headers, we just expect the 101 Switching Protocols back.

	// Start Pipes
	cID := uuid.New().String()

	// Client -> Server
	go pipeAndInspect(cID, clientConn, targetConn, "outgoing", ctx)

	// Server -> Client
	// We pass clientBuf because hijacking might have buffered data?
	// Actually for the first read validation we might stick to conn.
	pipeAndInspect(cID, targetConn, clientConn, "incoming", ctx)
}

func pipeAndInspect(id string, src io.Reader, dst io.Writer, direction string, ctx context.Context) {
	// We use a TeeReader? No, we need to parse frames which is stateful.
	// Best approach: Read into buffer, Parse, Write to dst.

	buf := make([]byte, 32*1024)

	// We need a proper frame parser that doesn't consume bytes destructively from the pipe perspective?
	// Or we act as a true proxy: Read Frame -> Decode -> Re-serialize -> Write.
	// Re-serializing is safer but slower.
	// Generic TCP pipe is simplest but we can't parse easily because frames can be fragmented.

	// CopyBuffer is standard.
	// To inspect, we need a MultiWriter that writes to a "Parser" and the "Dst".

	pr, pw := io.Pipe()
	tee := io.TeeReader(src, pw)

	// Parser Goroutine
	go func() {
		defer pr.Close()
		bufReader := bufio.NewReader(pr)
		for {
			payload, opcode, _, err := readFrame(bufReader)
			if err != nil {
				return // Stream closed or error
			}

			// Emit Event
			if ctx != nil && (opcode == 1 || opcode == 2) { // Text or Binary
				msg := ""
				if opcode == 1 {
					msg = string(payload)
					if len(msg) > 500 {
						msg = msg[:500] + "..."
					}
				} else {
					msg = fmt.Sprintf("[Binary data: %d bytes]", len(payload))
				}

				runtime.EventsEmit(ctx, "proxy:ws:frame", WebSocketFrame{
					ID:        uuid.New().String(),
					Time:      time.Now().Format("15:04:05.000"),
					Direction: direction,
					OpCode:    opcode,
					Payload:   msg,
					Length:    int64(len(payload)),
				})
			}
		}
	}()

	io.CopyBuffer(dst, tee, buf)
	pw.Close() // Close the parser writer when copy finishes
}

func readFrame(r *bufio.Reader) ([]byte, int, bool, error) {
	// RFC 6455 Simplistic Parser
	b0, err := r.ReadByte()
	if err != nil {
		return nil, 0, false, err
	}

	opcode := int(b0 & 0x0F)

	b1, err := r.ReadByte()
	if err != nil {
		return nil, 0, false, err
	}

	masked := (b1 & 0x80) != 0
	payloadLen := int64(b1 & 0x7F)

	if payloadLen == 126 {
		var u16 uint16
		binary.Read(r, binary.BigEndian, &u16)
		payloadLen = int64(u16)
	} else if payloadLen == 127 {
		binary.Read(r, binary.BigEndian, &payloadLen)
	}

	// Limit payload read for safety (1MB)
	if payloadLen > 1024*1024 {
		// Skip
		io.CopyN(io.Discard, r, payloadLen)
		return nil, opcode, masked, nil
	}

	var maskKey []byte
	if masked {
		maskKey = make([]byte, 4)
		if _, err := io.ReadFull(r, maskKey); err != nil {
			return nil, 0, false, err
		}
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, 0, false, err
	}

	if masked {
		for i := 0; i < len(payload); i++ {
			payload[i] ^= maskKey[i%4]
		}
	}

	return payload, opcode, masked, nil
}
