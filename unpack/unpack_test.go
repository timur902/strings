package unpack

import "testing"

func TestUnpack(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"basic", "a4bc2d5e", "aaaabccddddde", false},
		{"no numbers", "abcd", "abcd", false},
		{"empty", "", "", false},
		{"starts with digit", "3abc", "", true},
		{"two digits", "aaa10b", "", true},
		{"zero repeat", "a0b", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %s got %s", tt.want, got)
			}
		})
	}
}