package session

import (
	"testing"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFingerprint(t *testing.T) {
	t.Run("Deterministic", func(t *testing.T) {
		h1 := GenerateFingerprint("gemini-pro", "Hello World")
		h2 := GenerateFingerprint("gemini-pro", "Hello World")
		assert.Equal(t, h1, h2)
	})

	t.Run("Case Insensitivity Model", func(t *testing.T) {
		h1 := GenerateFingerprint("GEMINI-PRO", "Hello World")
		h2 := GenerateFingerprint("gemini-pro", "Hello World")
		assert.Equal(t, h1, h2)
	})

	t.Run("Different Content", func(t *testing.T) {
		h1 := GenerateFingerprint("gemini-pro", "Hello World")
		h2 := GenerateFingerprint("gemini-pro", "Goodbye World")
		assert.NotEqual(t, h1, h2)
	})

	t.Run("Empty Content", func(t *testing.T) {
		h := GenerateFingerprint("model", "")
		assert.Equal(t, "", h)
	})
}

func TestExtractFirstUserMessage(t *testing.T) {
	t.Run("Extracts Content", func(t *testing.T) {
		req := &mappers.GeminiRequest{
			Contents: []mappers.GeminiContent{
				{
					Role: "user",
					Parts: []mappers.GeminiPart{
						{Text: "Hello "},
						{Text: "World"},
					},
				},
			},
		}
		content := ExtractFirstUserMessage(req)
		assert.Equal(t, "Hello World", content)
	})

	t.Run("Ignores System/Model", func(t *testing.T) {
		req := &mappers.GeminiRequest{
			Contents: []mappers.GeminiContent{
				{Role: "model", Parts: []mappers.GeminiPart{{Text: "Hi"}}},
				{Role: "user", Parts: []mappers.GeminiPart{{Text: "My Prompt"}}},
			},
		}
		content := ExtractFirstUserMessage(req)
		assert.Equal(t, "My Prompt", content)
	})

	t.Run("Empty", func(t *testing.T) {
		req := &mappers.GeminiRequest{}
		content := ExtractFirstUserMessage(req)
		assert.Equal(t, "", content)
	})
}
