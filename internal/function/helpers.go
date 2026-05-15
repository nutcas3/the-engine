package function

// GetString extracts a string value from a nested map
func GetString(obj map[string]any, path ...string) string {
	current := obj
	for i, key := range path {
		if i == len(path)-1 {
			if val, ok := current[key].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

// GetFloat64 extracts a float64 value from a nested map
func GetFloat64(obj map[string]any, path ...string) float64 {
	current := obj
	for i, key := range path {
		if i == len(path)-1 {
			switch val := current[key].(type) {
			case float64:
				return val
			case float32:
				return float64(val)
			case int:
				return float64(val)
			case int64:
				return float64(val)
			default:
				return 0
			}
		}
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return 0
		}
	}
	return 0
}

// GetBool extracts a bool value from a nested map
func GetBool(obj map[string]any, path ...string) bool {
	current := obj
	for i, key := range path {
		if i == len(path)-1 {
			if val, ok := current[key].(bool); ok {
				return val
			}
			return false
		}
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return false
		}
	}
	return false
}
