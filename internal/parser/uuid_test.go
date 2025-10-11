package parser

import (
	"testing"

	"github.com/google/uuid"
)

func TestParse_ValidFormats(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantV  string
		wantOK bool
	}{
		{
			name:   "canonical",
			input:  "550e8400-e29b-41d4-a716-446655440000",
			wantV:  "4",
			wantOK: true,
		},
		{
			name:   "without dashes",
			input:  "550e8400e29b41d4a716446655440000",
			wantV:  "4",
			wantOK: true,
		},
		{
			name:   "URN",
			input:  "urn:uuid:550e8400-e29b-41d4-a716-446655440000",
			wantV:  "4",
			wantOK: true,
		},
		{
			name:   "braces",
			input:  "{550e8400-e29b-41d4-a716-446655440000}",
			wantV:  "4",
			wantOK: true,
		},
		{
			name:   "base58",
			input:  "BWBeN28Vb7cMEx7Ym8AUzs",
			wantV:  "4",
			wantOK: true,
		},
		{
			name:   "base32",
			input:  "KUHIIAHCTNA5JJYWIRTFKRAAAA",
			wantV:  "4",
			wantOK: true,
		},
		{
			name:   "uuid v1",
			input:  "c232ab00-9414-11ec-b3c8-9f6bdeced846",
			wantV:  "1",
			wantOK: true,
		},
		{
			name:   "uuid v7",
			input:  "018e8f2d-5e77-7000-8000-000000000000",
			wantV:  "7",
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Parse(tt.input)
			if tt.wantOK {
				if info.Error != "" {
					t.Errorf("Parse() error = %v, want success", info.Error)
					return
				}
				if info.Version != tt.wantV {
					t.Errorf("Parse() version = %v, want %v", info.Version, tt.wantV)
				}
				if info.Canonical == "" {
					t.Errorf("Parse() canonical is empty")
				}
			} else {
				if info.Error == "" {
					t.Errorf("Parse() expected error, got success")
				}
			}
		})
	}
}

func TestParse_InvalidFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"invalid hex", "550e8400-e29b-41d4-a716-44665544000g"},
		{"too short", "550e8400-e29b"},
		{"invalid base58", "invalid-base58-string"},
		{"random string", "not-a-uuid-at-all"},
		{"wrong length", "550e8400e29b41d4a71644665544"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Parse(tt.input)
			if info.Error == "" {
				t.Errorf("Parse(%q) expected error, got success", tt.input)
			}
		})
	}
}

func TestParse_OutputFormats(t *testing.T) {
	input := "550e8400-e29b-41d4-a716-446655440000"
	info := Parse(input)

	if info.Error != "" {
		t.Fatalf("Parse() error = %v", info.Error)
	}

	// Check all output formats are present
	if info.Canonical == "" {
		t.Error("Canonical is empty")
	}
	if info.Base58 == "" {
		t.Error("Base58 is empty")
	}
	if info.Base32 == "" {
		t.Error("Base32 is empty")
	}
	if info.Hex == "" {
		t.Error("Hex is empty")
	}
	if info.URN == "" {
		t.Error("URN is empty")
	}
	if info.Version == "" {
		t.Error("Version is empty")
	}
	if info.Variant == "" {
		t.Error("Variant is empty")
	}

	// Verify Base58/Base32 can be parsed back
	infoB58 := Parse(info.Base58)
	if infoB58.Error != "" {
		t.Errorf("Base58 round-trip failed: %v", infoB58.Error)
	}
	if infoB58.Canonical != info.Canonical {
		t.Errorf("Base58 round-trip mismatch: got %v, want %v", infoB58.Canonical, info.Canonical)
	}

	infoB32 := Parse(info.Base32)
	if infoB32.Error != "" {
		t.Errorf("Base32 round-trip failed: %v", infoB32.Error)
	}
	if infoB32.Canonical != info.Canonical {
		t.Errorf("Base32 round-trip mismatch: got %v, want %v", infoB32.Canonical, info.Canonical)
	}
}

func TestExtractTime_V1(t *testing.T) {
	// UUID v1 with known timestamp
	u := uuid.MustParse("c232ab00-9414-11ec-b3c8-9f6bdeced846")
	tm, err := extractTime(u, 1)
	if err != nil {
		t.Errorf("extractTime(v1) error = %v", err)
	}
	if tm.IsZero() {
		t.Error("extractTime(v1) returned zero time")
	}
}

func TestExtractTime_V6(t *testing.T) {
	// UUID v6 (reordered v1)
	u := uuid.MustParse("1ec9414c-232a-6b00-b3c8-9f6bdeced846")
	tm, err := extractTime(u, 6)
	if err != nil {
		t.Errorf("extractTime(v6) error = %v", err)
	}
	if tm.IsZero() {
		t.Error("extractTime(v6) returned zero time")
	}
}

func TestExtractTime_V7(t *testing.T) {
	// UUID v7 with Unix timestamp
	u := uuid.MustParse("018e8f2d-5e77-7000-8000-000000000000")
	tm, err := extractTime(u, 7)
	if err != nil {
		t.Errorf("extractTime(v7) error = %v", err)
	}
	if tm.IsZero() {
		t.Error("extractTime(v7) returned zero time")
	}
}

func TestExtractTime_V4(t *testing.T) {
	// UUID v4 (random) should not extract time
	u := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	_, err := extractTime(u, 4)
	if err == nil {
		t.Error("extractTime(v4) expected error, got nil")
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name string
		uuid uuid.UUID
		want bool
	}{
		{
			name: "valid v4",
			uuid: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			want: true,
		},
		{
			name: "valid v1",
			uuid: uuid.MustParse("c232ab00-9414-11ec-b3c8-9f6bdeced846"),
			want: true,
		},
		{
			name: "invalid version 0",
			uuid: uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: false,
		},
		{
			name: "invalid version 15",
			uuid: uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUUID(tt.uuid); got != tt.want {
				t.Errorf("isValidUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse_TimeExtraction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		hasTime bool
	}{
		{
			name:    "v1 has time",
			input:   "c232ab00-9414-11ec-b3c8-9f6bdeced846",
			hasTime: true,
		},
		{
			name:    "v4 no time",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			hasTime: false,
		},
		{
			name:    "v6 has time",
			input:   "1ec9414c-232a-6b00-b3c8-9f6bdeced846",
			hasTime: true,
		},
		{
			name:    "v7 has time",
			input:   "018e8f2d-5e77-7000-8000-000000000000",
			hasTime: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Parse(tt.input)
			if info.Error != "" {
				t.Fatalf("Parse() error = %v", info.Error)
			}
			hasTime := info.Time != ""
			if hasTime != tt.hasTime {
				t.Errorf("Parse() time extraction = %v, want %v", hasTime, tt.hasTime)
			}
		})
	}
}
