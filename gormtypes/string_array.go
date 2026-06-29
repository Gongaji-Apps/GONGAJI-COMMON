package gormtypes

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray memetakan kolom Postgres text[]/varchar[] tanpa dependensi tambahan
// (tanpa lib/pq). Cukup untuk nilai sederhana — token/enum tanpa koma atau kutip
// di dalam elemen. Untuk elemen kompleks, pakai jsonb (JSONB) atau lib/pq.
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *StringArray) Scan(value any) error {
	if value == nil {
		*a = nil
		return nil
	}

	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("StringArray: tipe tak didukung %T", value)
	}

	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		*a = StringArray{}
		return nil
	}

	parts := strings.Split(s, ",")
	out := make(StringArray, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, `"`)
		out = append(out, p)
	}
	*a = out
	return nil
}
