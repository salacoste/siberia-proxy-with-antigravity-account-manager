package mappers

// Gemini Internal API Structs
// These match Google's v1internal structure used by Cloud Code

type GeminiRequest struct {
	Contents          []GeminiContent   `json:"contents"`
	SystemInstruction *GeminiContent    `json:"system_instruction,omitempty"`
	Tools             []GeminiTool      `json:"tools,omitempty"`
	ToolConfig        *GeminiToolConfig `json:"tool_config,omitempty"`
	GenerationConfig  *GeminiGenConfig  `json:"generation_config,omitempty"`
	SafetySettings    []GeminiSafetySet `json:"safety_settings,omitempty"`
}

type GeminiContent struct {
	Role  string       `json:"role,omitempty"` // "user", "model"
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text             string          `json:"text,omitempty"`
	InlineData       *GeminiBlob     `json:"inline_data,omitempty"`
	FunctionCall     *GeminiFuncCall `json:"function_call,omitempty"`
	FunctionResponse *GeminiFuncResp `json:"function_response,omitempty"`
}

type GeminiBlob struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // Base64
}

type GeminiFuncCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type GeminiFuncResp struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

type GeminiTool struct {
	FunctionDeclarations []GeminiFuncDecl `json:"function_declarations"`
}

type GeminiFuncDecl struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters"` // OpenAPI Schema
}

type GeminiToolConfig struct {
	FunctionCallingConfig *GeminiFuncCallingConfig `json:"function_calling_config,omitempty"`
}

type GeminiFuncCallingConfig struct {
	Mode string `json:"mode"` // "AUTO", "ANY", "NONE"
}

type GeminiGenConfig struct {
	Temperature     float64  `json:"temperature,omitempty"`
	TopP            float64  `json:"top_p,omitempty"`
	TopK            int      `json:"top_k,omitempty"`
	CandidateCount  int      `json:"candidate_count,omitempty"`
	StopSequences   []string `json:"stop_sequences,omitempty"`
	MaxOutputTokens int      `json:"max_output_tokens,omitempty"`
}

type GeminiSafetySet struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// Response Structures

type GeminiResponse struct {
	Candidates    []GeminiCandidate    `json:"candidates"`
	UsageMetadata *GeminiUsageMetadata `json:"usage_metadata,omitempty"`
}

type GeminiCandidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finish_reason"`
	Index        int           `json:"index"`
}

type GeminiUsageMetadata struct {
	PromptTokenCount     int `json:"prompt_token_count"`
	CandidatesTokenCount int `json:"candidates_token_count"`
	TotalTokenCount      int `json:"total_token_count"`
}
