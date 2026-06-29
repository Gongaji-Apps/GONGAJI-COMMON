package gormtypes

import "testing"

func TestJSONB_ValueScanRoundTrip(t *testing.T) {
	src := JSONB(`{"a":1,"b":[2,3]}`)
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	var dst JSONB
	if err := dst.Scan(v); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if string(dst) != string(src) {
		t.Fatalf("round-trip: got %q want %q", dst, src)
	}
}

func TestJSONB_NilIsNull(t *testing.T) {
	var j JSONB
	v, err := j.Value()
	if err != nil || v != nil {
		t.Fatalf("nil JSONB harus NULL, got %v err %v", v, err)
	}
	if err := j.Scan(nil); err != nil || j != nil {
		t.Fatalf("Scan(nil) harus nil, got %v err %v", j, err)
	}
}

func TestStringArray_ValueScanRoundTrip(t *testing.T) {
	src := StringArray{"online", "offline"}
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if v.(string) != `{"online","offline"}` {
		t.Fatalf("Value format: %v", v)
	}
	var dst StringArray
	if err := dst.Scan(v); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(dst) != 2 || dst[0] != "online" || dst[1] != "offline" {
		t.Fatalf("round-trip: %v", dst)
	}
}

func TestStringArray_Empty(t *testing.T) {
	var a StringArray
	if err := a.Scan("{}"); err != nil || len(a) != 0 {
		t.Fatalf("scan empty: %v err %v", a, err)
	}
}
