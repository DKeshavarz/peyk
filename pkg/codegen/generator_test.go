package codegen

import (
	"regexp"
	"testing"
)

func TestGenerate_WithRegex(t *testing.T) {
	tests := []struct {
		name      string
		groups    []int
		separator string
		charset   string
		pattern   string // regex pattern to match
	}{
		{
			name:      "default xxx-xxxx-xxx lowercase",
			groups:    []int{3, 4, 3},
			separator: "-",
			charset:   "abcdefghijklmnopqrstuvwxyz",
			pattern:   `^[a-z]{3}-[a-z]{4}-[a-z]{3}$`,
		},
		{
			name:      "xxx-xxxx-xxx uppercase",
			groups:    []int{3, 4, 3},
			separator: "-",
			charset:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			pattern:   `^[A-Z]{3}-[A-Z]{4}-[A-Z]{3}$`,
		},
		{
			name:      "xxxx-xxxx numeric",
			groups:    []int{4, 4},
			separator: "-",
			charset:   "0123456789",
			pattern:   `^\d{4}-\d{4}$`,
		},
		{
			name:      "xxx.xxx.xxx with dots",
			groups:    []int{3, 3, 3},
			separator: ".",
			charset:   "abcdefghijklmnopqrstuvwxyz",
			pattern:   `^[a-z]{3}\.[a-z]{3}\.[a-z]{3}$`,
		},
		{
			name:      "xxx_xxx_xxx with underscore",
			groups:    []int{3, 3, 3},
			separator: "_",
			charset:   "abcdefghijklmnopqrstuvwxyz",
			pattern:   `^[a-z]{3}_[a-z]{3}_[a-z]{3}$`,
		},
		{
			name:      "alphanumeric mix",
			groups:    []int{4, 4},
			separator: "-",
			charset:   "abcdefghijklmnopqrstuvwxyz0123456789",
			pattern:   `^[a-z0-9]{4}-[a-z0-9]{4}$`,
		},
		{
			name:      "single group no separator",
			groups:    []int{8},
			separator: "",
			charset:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			pattern:   `^[A-Z]{8}$`,
		},
		{
			name:      "hex format",
			groups:    []int{4, 4},
			separator: "-",
			charset:   "0123456789ABCDEF",
			pattern:   `^[0-9A-F]{4}-[0-9A-F]{4}$`,
		},
		{
			name:      "custom charset with symbols",
			groups:    []int{5},
			separator: "",
			charset:   "ABCDEFGHJKLMNPQRSTUVWXYZ23456789", // Base32 without ambiguous chars
			pattern:   `^[ABCDEFGHJKLMNPQRSTUVWXYZ23456789]{5}$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			g := &codeGenerator{
				groups:    tt.groups,
				separator: tt.separator,
				charset:   tt.charset,
			}

			
			for range 100 {
				code, err := g.Generate()
				if err != nil {
					t.Fatalf("Generate() error = %v", err)
				}

				matched, err := regexp.MatchString(tt.pattern, code)
				if err != nil {
					t.Fatalf("Regex error: %v", err)
				}

				if !matched {
					t.Errorf("Code %q does not match pattern %q", code, tt.pattern)
				}

				
				for _, char := range code {
					charStr := string(char)
					if charStr != tt.separator && !contains(tt.charset, charStr) {
						t.Errorf("Code contains character %q not in charset %q", charStr, tt.charset)
					}
				}
			}
		})
	}
}

// Helper function
func contains(charset, char string) bool {
	for _, c := range charset {
		if string(c) == char {
			return true
		}
	}
	return false
}
