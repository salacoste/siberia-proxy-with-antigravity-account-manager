package oauth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
)

// Server handles the local loopback for OAuth flow.
// It listens on a random free port, returns the port,
// and blocks until a callback is received or context is cancelled.
type Server struct {
	mu         sync.Mutex
	server     *http.Server
	resultChan chan string
	errChan    chan error
}

func NewServer() *Server {
	return &Server{
		resultChan: make(chan string, 1),
		errChan:    make(chan error, 1),
	}
}

// Start listens on 127.0.0.1:0 (ephemeral port) and returns the port number.
func (s *Server) Start(expectedState string) (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("failed to listen on loopback: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		state := query.Get("state")

		// Basic Validation
		if state != expectedState {
			http.Error(w, "State mismatch (possible CSRF)", http.StatusBadRequest)
			s.errChan <- fmt.Errorf("state mismatch")
			return
		}
		if code == "" {
			http.Error(w, "Missing Auth Code", http.StatusBadRequest)
			s.errChan <- fmt.Errorf("missing code")
			return
		}

		// Capture Logic
		s.resultChan <- code

		// Serve Success Page
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(successHTML))

		// Auto-Shutdown in background
		go func() {
			s.server.Shutdown(context.Background())
		}()
	})

	s.server = &http.Server{
		Handler: mux,
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[OAuth] Error: %v\n", err)
		}
	}()

	return port, nil
}

// WaitForCode blocks until a code is received or ctx is cancelled.
func (s *Server) WaitForCode(ctx context.Context) (string, error) {
	select {
	case code := <-s.resultChan:
		return code, nil
	case err := <-s.errChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

const successHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Siberia Authentication</title>
    <style>
        body { font-family: sans-serif; text-align: center; padding-top: 50px; background: #111; color: #fff; }
        .card { background: #222; padding: 40px; border-radius: 12px; display: inline-block; }
        h1 { color: #4ade80; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Authentication Successful</h1>
        <p>You can close this window and return to Siberia.</p>
    </div>
    <script>window.setTimeout(function(){window.close();}, 2000);</script>
</body>
</html>
`
