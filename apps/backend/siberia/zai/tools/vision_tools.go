package tools

import (
	"fmt"

	"github.com/salacoste/siberia/siberia/zai"
)

// CallUiToArtifact generates frontend code from a UI screenshot
func CallUiToArtifact(apiKey, baseURL, imageSource string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("missing z.ai api key")
	}

	client := zai.NewVisionClient(baseURL, apiKey)

	imgPart, err := zai.ImageSourceToContent(imageSource)
	if err != nil {
		return "", fmt.Errorf("process image: %w", err)
	}

	messages := []zai.Message{
		{
			Role: "system",
			Content: []zai.ContentPart{
				{Type: "text", Text: "You are an expert Frontend Engineer. Your task is to analyze the provided UI screenshot and generate a complete, production-ready React component using Tailwind CSS that replicates the design effectively. Output ONLY the code block."},
			},
		},
		{
			Role: "user",
			Content: []zai.ContentPart{
				imgPart,
				{Type: "text", Text: "Generate the React/Tailwind code for this UI."},
			},
		},
	}

	return client.VisionChatCompletion(messages)
}

// CallExtractText performs OCR
func CallExtractText(apiKey, baseURL, imageSource string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("missing z.ai api key")
	}

	client := zai.NewVisionClient(baseURL, apiKey)

	imgPart, err := zai.ImageSourceToContent(imageSource)
	if err != nil {
		return "", fmt.Errorf("process image: %w", err)
	}

	messages := []zai.Message{
		{
			Role: "user",
			Content: []zai.ContentPart{
				imgPart,
				{Type: "text", Text: "Detailed transcription of all text in this image. Maintain layout structure where possible."},
			},
		},
	}

	return client.VisionChatCompletion(messages)
}

// CallDiagnoseError analyzes an error screenshot
func CallDiagnoseError(apiKey, baseURL, imageSource string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("missing z.ai api key")
	}

	client := zai.NewVisionClient(baseURL, apiKey)

	imgPart, err := zai.ImageSourceToContent(imageSource)
	if err != nil {
		return "", fmt.Errorf("process image: %w", err)
	}

	messages := []zai.Message{
		{
			Role: "system",
			Content: []zai.ContentPart{
				{Type: "text", Text: "You are an expert Software Architect and debugger. Analyze the error message or traceback in the screenshot. Explain the root cause and provide actionable fixes."},
			},
		},
		{
			Role: "user",
			Content: []zai.ContentPart{
				imgPart,
				{Type: "text", Text: "Diagnose this error."},
			},
		},
	}

	return client.VisionChatCompletion(messages)
}

// CallUnderstandDiagram explains a technical diagram
func CallUnderstandDiagram(apiKey, baseURL, imageSource string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("missing z.ai api key")
	}

	client := zai.NewVisionClient(baseURL, apiKey)

	imgPart, err := zai.ImageSourceToContent(imageSource)
	if err != nil {
		return "", fmt.Errorf("process image: %w", err)
	}

	messages := []zai.Message{
		{
			Role: "system",
			Content: []zai.ContentPart{
				{Type: "text", Text: "You are an expert Technical Writer. Explain the architecture or workflow shown in this diagram in clear Markdown steps."},
			},
		},
		{
			Role: "user",
			Content: []zai.ContentPart{
				imgPart,
				{Type: "text", Text: "Explain this diagram."},
			},
		},
	}

	return client.VisionChatCompletion(messages)
}
