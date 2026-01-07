package mappers

// ImageRequest represents a generic request for image generation
type ImageRequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`               // Number of images to generate
	Size           string `json:"size,omitempty"`            // 256x256, 512x512, 1024x1024
	ResponseFormat string `json:"response_format,omitempty"` // url or b64_json
	Style          string `json:"style,omitempty"`           // vivid, natural (OpenAI specific but useful to map)
	User           string `json:"user,omitempty"`
}

// ImageResponse represents the generic response
type ImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

type ImageData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}
