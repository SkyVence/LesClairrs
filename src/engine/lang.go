package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Catalog map[string]string

var ph = regexp.MustCompile(`\{[A-Za-z0-9_.-]+\}`)

// NestedData represents the nested JSON structure from language files
type NestedData map[string]interface{}

func Load(lang string) (Catalog, error) {
	b, err := os.ReadFile("assets/interface/" + lang + ".json")
	if err != nil {
		log.Printf("Failed to read language file %s: %v", lang, err)
		return nil, err
	}

	// Parse nested JSON structure
	var nested NestedData
	if err := json.Unmarshal(b, &nested); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Flatten nested structure into dot-notation keys
	catalog := make(Catalog)
	flattenMap(nested, "", catalog)

	return catalog, nil
}

// flattenMap recursively flattens nested maps into dot-notation keys
func flattenMap(data map[string]interface{}, prefix string, result Catalog) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			// Direct string value
			result[fullKey] = v
		case map[string]interface{}:
			// Nested object - recurse
			flattenMap(v, fullKey, result)
		default:
			// Convert other types to string
			result[fullKey] = fmt.Sprint(v)
		}
	}
}

// Text replaces placeholders like {player}, {hp}, {max} in encounter order.
func (c Catalog) Text(key string, args ...any) string {
	s, ok := c[key]
	if !ok {
		return "⟦" + key + "⟧"
	}
	if len(args) == 0 {
		return s
	}
	idx := 0
	return ph.ReplaceAllStringFunc(s, func(match string) string {
		if idx < len(args) {
			v := fmt.Sprint(args[idx])
			idx++
			return v
		}
		// not enough args: leave the placeholder as-is
		return match
	})
}
