package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Catalog map[string]string

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

func (c Catalog) Text(key string, kv ...any) string {
	s, ok := c[key]
	if !ok {
		return "⟦" + key + "⟧"
	}
	if len(kv) == 0 {
		return s
	}
	repls := make([]string, 0, len(kv))
	for i := 0; i+1 < len(kv); i += 2 {
		k := fmt.Sprint(kv[i])
		v := fmt.Sprint(kv[i+1])
		repls = append(repls, "{"+k+"}", v)
	}
	return strings.NewReplacer(repls...).Replace(s)
}
