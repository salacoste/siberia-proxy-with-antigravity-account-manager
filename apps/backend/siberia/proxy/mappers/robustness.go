package mappers

// robustness.go implements logic to clean and sanitize requests for upstream compatibility (Story-71).

// DeepClean recursively traverses a JSON-compatible structure (map, slice, string, etc.)
// and performs the following cleanups:
// 1. Replaces the string literal "[undefined]" with nil (or removes it from maps/slices if possible).
// 2. Uppercases "type": "object" -> "type": "OBJECT" (required for Gemini compatibility).
func DeepClean(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range val {
			// Check for [undefined] string value
			if strVal, ok := v.(string); ok && strVal == "[undefined]" {
				continue // Skip/Remove this key
			}

			// Recurse first
			cleanedVal := DeepClean(v)

			// Logic: Type Uppercasing for Gemini Schema
			if k == "type" {
				if s, ok := cleanedVal.(string); ok && s == "object" {
					cleanedVal = "OBJECT"
				}
			}

			newMap[k] = cleanedVal
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, 0, len(val))
		for _, item := range val {
			if strVal, ok := item.(string); ok && strVal == "[undefined]" {
				continue // Skip [undefined] in slice
			}
			newSlice = append(newSlice, DeepClean(item))
		}
		return newSlice
	case string:
		if val == "[undefined]" {
			return nil
		}
		return val
	default:
		return val
	}
}

// SanitizeSchema removes unsupported fields from JSON Schemas (e.g. for tool definitions).
// Removes: "format", "strict", "additionalProperties" (if problematic).
func SanitizeSchema(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range val {
			// Blacklist fields
			if k == "format" || k == "strict" || k == "additionalProperties" {
				continue
			}
			newMap[k] = SanitizeSchema(v)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, 0, len(val))
		for _, item := range val {
			newSlice = append(newSlice, SanitizeSchema(item))
		}
		return newSlice
	default:
		return val
	}
}

// FilterWebSearch recursively removes any tool definition named "web_search" or "googleSearchRetrieval".
// This prevents conflicts when the proxy injects its own search tool.
func FilterWebSearch(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range val {
			newMap[k] = FilterWebSearch(v)
		}
		// Logic: If this map itself describes a tool with name="web_search", we can't "return nil" easily from here if the parent expects a map value.
		// So we rely on the Slice handler (below) to drop items.
		// But if the map IS the tool (e.g. at root or single property), maybe we return nil?
		// Usually tools are in a list.
		// But if we encounter { "name": "web_search" ... } here, we return it cleaned.
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, 0, len(val))
		for _, item := range val {
			// Check item content
			if m, ok := item.(map[string]interface{}); ok {
				// Gemini Check
				if name, ok := m["name"].(string); ok {
					if name == "web_search" || name == "googleSearchRetrieval" {
						continue // Drop it
					}
				}
				// OpenAI Check: tools -> item.type="function", item.function.name="..."
				if f, ok := m["function"].(map[string]interface{}); ok {
					if name, ok := f["name"].(string); ok {
						if name == "web_search" || name == "googleSearchRetrieval" {
							continue // Drop
						}
					}
				}
			}
			newSlice = append(newSlice, FilterWebSearch(item))
		}
		return newSlice
	default:
		return val
	}
}
