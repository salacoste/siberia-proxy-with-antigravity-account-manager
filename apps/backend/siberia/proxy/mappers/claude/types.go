package claude

// Claude Messages API Types

type MessageRequest struct {
	Model         string           `json:"model"`
	Messages      []Message        `json:"messages"`
	System        any              `json:"system,omitempty"` // string or []ContentBlock
	MaxTokens     int              `json:"max_tokens"`
	Metadata      *MessageMetadata `json:"metadata,omitempty"`
	StopSequences []string         `json:"stop_sequences,omitempty"`
	Stream        bool             `json:"stream,omitempty"`
	Temperature   float64          `json:"temperature,omitempty"`
	TopP          float64          `json:"top_p,omitempty"`
	TopK          int              `json:"top_k,omitempty"`
	Tools         []Tool           `json:"tools,omitempty"`
	ToolChoice    any              `json:"tool_choice,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string or []ContentBlock
}

type ContentBlock struct {
	Type         string        `json:"type"` // "text", "image", "tool_use", "tool_result"
	Text         string        `json:"text,omitempty"`
	Source       *ImageSource  `json:"source,omitempty"`
	ID           string        `json:"id,omitempty"`            // For tool_use/tool_result
	Name         string        `json:"name,omitempty"`          // For tool_use
	Input        any           `json:"input,omitempty"`         // For tool_use (JSON)
	ToolUseID    string        `json:"tool_use_id,omitempty"`   // For tool_result
	Content      any           `json:"content,omitempty"`       // For tool_result
	CacheControl *CacheControl `json:"cache_control,omitempty"` // To be stripped

	// Thinking fields (Claude 3.7+ custom / Vertex)
	Thinking  string `json:"thinking,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type ImageSource struct {
	Type      string `json:"type"` // "base64"
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type CacheControl struct {
	Type string `json:"type"` // "ephemeral"
}

type MessageMetadata struct {
	UserID string `json:"user_id,omitempty"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	InputSchema any    `json:"input_schema"` // JSON Schema
}

// Response Structures

type MessageResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"` // "message"
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason,omitempty"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        *Usage         `json:"usage,omitempty"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
