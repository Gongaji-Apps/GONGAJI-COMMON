package converter_test

import (
	"os"
	"testing"

	"github.com/Gongaji-Apps/GONGAJI-FRAMEWORK/converter"
)

func TestStringToInt(t *testing.T) {
	os.Setenv("ADDITIONAL_ERR_500", "extra error")

	tests := []struct {
		name    string
		input   string
		want    *int
		wantErr bool
	}{
		{
			name:  "valid int",
			input: "123",
			want:  ptrInt(123),
		},
		{
			name:    "invalid int",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.StringToInt(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
				return
			}

			if err == nil && *result != *tt.want {
				t.Errorf("expected %d, got %d", *tt.want, *result)
			}
		})
	}
}

func ptrInt(v int) *int {
	return &v
}

func TestStringToDecimal(t *testing.T) {
	if d, err := converter.StringToDecimal("12.50"); err != nil || d == nil || d.String() != "12.5" {
		t.Fatalf("valid: got %v err %v", d, err)
	}
	if _, err := converter.StringToDecimal("abc"); err == nil {
		t.Fatalf("invalid harus error")
	}
}
