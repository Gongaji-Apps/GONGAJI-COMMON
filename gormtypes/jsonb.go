// Package gormtypes menyediakan tipe kolom GORM/database-sql kustom untuk Postgres
// (mis. jsonb, text[]) tanpa dependensi eksternal — cukup implementasi
// driver.Valuer + sql.Scanner sehingga bisa dipakai langsung sebagai field model.
package gormtypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONB memetakan kolom Postgres jsonb tanpa dependensi tambahan.
//   - Value/Scan             : tulis/baca jsonb (nil → NULL).
//   - Marshal/Unmarshal JSON : pass-through agar tampil sebagai JSON bersarang
//     (bukan string ter-escape) saat di-serialize ke respons API.
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = append((*j)[:0], v...)
	case string:
		*j = append((*j)[:0], v...)
	default:
		return fmt.Errorf("JSONB: tipe tak didukung %T", value)
	}
	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSONB) UnmarshalJSON(b []byte) error {
	if j == nil {
		return errors.New("JSONB: UnmarshalJSON pada penerima nil")
	}
	*j = append((*j)[:0], b...)
	return nil
}
