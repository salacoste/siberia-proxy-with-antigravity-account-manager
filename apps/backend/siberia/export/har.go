package export

import (
	"encoding/json"
	"time"

	"github.com/salacoste/siberia/siberia/proxy"
)

// HAR 1.2 Structures
type HAR struct {
	Log Log `json:"log"`
}

type Log struct {
	Version string  `json:"version"`
	Creator Creator `json:"creator"`
	Entries []Entry `json:"entries"`
}

type Creator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Entry struct {
	StartedDateTime string   `json:"startedDateTime"`
	Time            int64    `json:"time"`
	Request         Request  `json:"request"`
	Response        Response `json:"response"`
	Cache           Cache    `json:"cache"`
	Timings         Timings  `json:"timings"`
}

type Request struct {
	Method      string        `json:"method"`
	URL         string        `json:"url"`
	HTTPVersion string        `json:"httpVersion"`
	Headers     []NameValue   `json:"headers"`
	QueryString []NameValue   `json:"queryString"`
	Cookies     []interface{} `json:"cookies"` // Simplified
	HeaderSize  int           `json:"headerSize"`
	BodySize    int64         `json:"bodySize"`
	PostData    *PostData     `json:"postData,omitempty"`
}

type Response struct {
	Status      int           `json:"status"`
	StatusText  string        `json:"statusText"`
	HTTPVersion string        `json:"httpVersion"`
	Headers     []NameValue   `json:"headers"`
	Cookies     []interface{} `json:"cookies"`
	Content     Content       `json:"content"`
	RedirectURL string        `json:"redirectURL"`
	HeaderSize  int           `json:"headerSize"`
	BodySize    int64         `json:"bodySize"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PostData struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

type Content struct {
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
}

type Cache struct{}

type Timings struct {
	Send    int `json:"send"`
	Wait    int `json:"wait"`
	Receive int `json:"receive"`
}

func ToHAR(event proxy.ProxyRequestEvent) (string, error) {
	// Parse Time
	// event.Time is "15:04:05", which misses date.
	// ideally we want full ISO8601.
	// For now, assuming current day or just using what we have (HAR requires ISO).
	// Since event.Time is lossy (just HH:MM:SS), we'll construct a mock date or use Now if acceptable.
	// Ideally, ProxyRequestEvent should have full timestamp.
	// Let's use time.Now() formatted properly or append to today.
	// Actually, let's just make it ISO8601 of "today + time".

	now := time.Now()
	// Parse event.Time "15:04:05"
	t, err := time.Parse("15:04:05", event.Time)
	if err == nil {
		// Combine with today's date
		now = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())
	}
	startedDateTime := now.Format(time.RFC3339)

	// Headers
	reqHeaders := mapToNameValue(event.ReqHeaders)
	respHeaders := mapToNameValue(event.RespHeaders)

	// PostData
	var postData *PostData
	if event.ReqBody != "" {
		postData = &PostData{
			MimeType: getHeader(event.ReqHeaders, "Content-Type"),
			Text:     event.ReqBody,
		}
	}

	har := HAR{
		Log: Log{
			Version: "1.2",
			Creator: Creator{
				Name:    "Siberia Proxy",
				Version: "1.0",
			},
			Entries: []Entry{
				{
					StartedDateTime: startedDateTime,
					Time:            event.Duration,
					Request: Request{
						Method:      event.Method,
						URL:         event.URL,
						HTTPVersion: "HTTP/1.1",
						Headers:     reqHeaders,
						QueryString: []NameValue{}, // Not parsed separately in event
						Cookies:     []interface{}{},
						HeaderSize:  -1,
						BodySize:    -1, // Calculated
						PostData:    postData,
					},
					Response: Response{
						Status:      event.Status,
						StatusText:  "",
						HTTPVersion: "HTTP/1.1",
						Headers:     respHeaders,
						Cookies:     []interface{}{},
						Content: Content{
							Size:     event.Size,
							MimeType: getHeader(event.RespHeaders, "Content-Type"),
							Text:     event.RespBody,
						},
						RedirectURL: "",
						HeaderSize:  -1,
						BodySize:    event.Size,
					},
					Cache: Cache{},
					Timings: Timings{
						Send:    0,
						Wait:    int(event.Duration),
						Receive: 0,
					},
				},
			},
		},
	}

	bytes, err := json.Marshal(har)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func mapToNameValue(h map[string]string) []NameValue {
	res := []NameValue{}
	for k, v := range h {
		res = append(res, NameValue{Name: k, Value: v})
	}
	return res
}

func getHeader(h map[string]string, key string) string {
	// Simple lookup, ignoring case for MVP (assuming map keys are canonicalized or as provided)
	if v, ok := h[key]; ok {
		return v
	}
	return "application/octet-stream" // Default
}
