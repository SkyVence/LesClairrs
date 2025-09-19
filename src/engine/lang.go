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

// Load reads and parses language file from assets/interface/{lang}.json
// Returns a flattened Catalog with dot-notation keys for text retrieval
func Load(lang string) (Catalog, error) {
	b, err := os.ReadFile("assets/interface/" + lang + ".json")
	if err != nil {
		log.Printf("Failed to read language file %s: %v", lang, err)
		return nil, err
	}

	var nested NestedData
	if err := json.Unmarshal(b, &nested); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	catalog := make(Catalog)
	flattenMap(nested, "", catalog)

	return catalog, nil
}

// flattenMap recursively converts nested maps to dot-notation keys in result catalog
func flattenMap(data map[string]interface{}, prefix string, result Catalog) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			flattenMap(v, fullKey, result)
		default:
			result[fullKey] = fmt.Sprint(v)
		}
	}
}

// Text replaces placeholders like {player}, {hp}, {max} with provided args in order
// Returns the formatted string or a placeholder notation if key not found
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
		return match
	})
}
