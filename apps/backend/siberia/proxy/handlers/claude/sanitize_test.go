package claude

import (
	"reflect"
	"testing"

	"github.com/salacoste/siberia/siberia/proxy/mappers/claude"
)

func TestSanitizeMessages(t *testing.T) {
	tests := []struct {
		name     string
		input    []claude.Message
		expected []claude.Message
	}{
		{
			name: "Keep normal text",
			input: []claude.Message{
				{Role: "user", Content: "Hello"},
			},
			expected: []claude.Message{
				{Role: "user", Content: "Hello"},
			},
		},
		{
			name: "Remove empty thinking block",
			input: []claude.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Hi"},
						map[string]interface{}{"type": "thinking", "thinking": "   "},
					},
				},
			},
			// Expect last thinking removed because it's empty
			expected: []claude.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Hi"},
					},
				},
			},
		},
		{
			name: "Remove short thinking block (<5 chars)",
			input: []claude.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Hi"},
						map[string]interface{}{"type": "thinking", "thinking": "Hmm"},
					},
				},
			},
			expected: []claude.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Hi"},
					},
				},
			},
		},
		{
			name: "Preserve valid thinking block in middle",
			input: []claude.Message{
				{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{"type": "thinking", "thinking": "Valid thinking process here..."},
						map[string]interface{}{"type": "text", "text": "Answer"},
					},
				},
			},
			expected: []claude.Message{
				{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{"type": "thinking", "thinking": "Valid thinking process here..."},
						map[string]interface{}{"type": "text", "text": "Answer"},
					},
				},
			},
		},
		{
			name: "Remove trailing valid thinking block (Last Message)",
			input: []claude.Message{
				{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Answer"},
						map[string]interface{}{"type": "thinking", "thinking": "Partial valid thinking..."},
					},
				},
			},
			// Should be removed because it is TRAILING block of LAST message
			expected: []claude.Message{
				{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Answer"},
					},
				},
			},
		},
		{
			name: "Keep trailing thinking block (NOT Last Message)",
			input: []claude.Message{
				{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Answer"},
						map[string]interface{}{"type": "thinking", "thinking": "Thinking..."},
					},
				},
				{
					Role:    "user",
					Content: "Next",
				},
			},
			// Thinking preserved because it's not the last message
			expected: []claude.Message{
				{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Answer"},
						map[string]interface{}{"type": "thinking", "thinking": "Thinking..."},
					},
				},
				{
					Role:    "user",
					Content: "Next",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeMessages(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("sanitizeMessages() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStripAllThinking(t *testing.T) {
	input := []claude.Message{
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{"type": "text", "text": "Hi"},
				map[string]interface{}{"type": "thinking", "thinking": "Thinking..."},
			},
		},
	}
	expected := []claude.Message{
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{"type": "text", "text": "Hi"},
			},
		},
	}

	got := stripAllThinking(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("stripAllThinking() = %v, want %v", got, expected)
	}
}
