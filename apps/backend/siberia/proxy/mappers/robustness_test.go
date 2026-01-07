package mappers

import (
	"reflect"
	"testing"
)

func TestDeepClean(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "Remove [undefined] string",
			input: map[string]interface{}{"key": "[undefined]", "valid": "value"},
			want:  map[string]interface{}{"valid": "value"},
		},
		{
			name:  "Uppercase type: object",
			input: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			want:  map[string]interface{}{"type": "OBJECT", "properties": map[string]interface{}{}},
		},
		{
			name: "Nested Cleaning",
			input: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{"bad": "[undefined]", "good": "ok"},
					"[undefined]",
					"valid_string",
				},
			},
			want: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{"good": "ok"},
					"valid_string",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeepClean(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepClean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeSchema(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "Remove format and strict",
			input: map[string]interface{}{"type": "string", "format": "uri", "strict": true},
			want:  map[string]interface{}{"type": "string"},
		},
		{
			name: "Recursive Sanitization",
			input: map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": map[string]interface{}{"type": "string", "format": "date-time"},
				},
				"additionalProperties": false,
			},
			want: map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": map[string]interface{}{"type": "string"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeSchema(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SanitizeSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterWebSearch(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name: "Filter Gemini Web Search",
			input: []interface{}{
				map[string]interface{}{"name": "code_interpreter"},
				map[string]interface{}{"name": "web_search"},
			},
			want: []interface{}{
				map[string]interface{}{"name": "code_interpreter"},
			},
		},
		{
			name: "Filter OpenAI Web Search",
			input: []interface{}{
				map[string]interface{}{
					"type": "function",
					"function": map[string]interface{}{
						"name": "googleSearchRetrieval",
					},
				},
				map[string]interface{}{
					"type": "function",
					"function": map[string]interface{}{
						"name": "weather",
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"type": "function",
					"function": map[string]interface{}{
						"name": "weather",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterWebSearch(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterWebSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
