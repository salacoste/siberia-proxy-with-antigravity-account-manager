package upstream

import (
	"context"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// Client defines the interface for the Gemini Upstream Client
type Client interface {
	// GenerateContent sends a unary request to Gemini
	GenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (*mappers.GeminiResponse, string, error)

	// StreamGenerateContent sends a streaming request
	// It returns a channel of chunks or error
	StreamGenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (<-chan *mappers.GeminiResponse, <-chan error)

	// GenerateImage sends an image generation request
	GenerateImage(ctx context.Context, req *mappers.ImageRequest) (*mappers.ImageResponse, string, error)
}
