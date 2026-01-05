package types

// ProxyRequestEvent represents a captured HTTP request/response exchange
type ProxyRequestEvent struct {
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Status       int               `json:"status"`
	Duration     int64             `json:"duration_ms"` // milliseconds
	Time         string            `json:"time"`
	Size         int64             `json:"size"`
	ReqHeaders   map[string]string `json:"req_headers"`
	RespHeaders  map[string]string `json:"resp_headers"`
	ReqBody      string            `json:"req_body"`
	RespBody     string            `json:"resp_body"`
	ConnectionID string            `json:"connection_id"`
}
