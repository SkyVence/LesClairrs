package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

type Catalog map[string]string

var ph = regexp.MustCompile(`\{[A-Za-z0-9_.-]+\}`)

//go:embed assets/interface/*.json
var efs embed.FS

func Load(lang string) (Catalog, error) {
	var c Catalog
	b, err := efs.ReadFile("assets/interface/" + lang + ".json")
	if err != nil {
		log.Fatalf("read %s: %v", lang, err)
		return nil, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// T replaces placeholders like {player}, {hp}, {max} in encounter order.
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
