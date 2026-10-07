package service

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		regenKey        bool
		countRandLetter int
		wantHashLen     int
	}{
		{
			name:            "Standard MD5 generation without random suffix",
			url:             "https://yandex.ru",
			regenKey:        false,
			countRandLetter: 0,
			wantHashLen:     32,
		},
		{
			name:            "MD5 generation with 3 random letters",
			url:             "https://yandex.ru",
			regenKey:        true,
			countRandLetter: 3,
			wantHashLen:     32 + 4,
		},
		{
			name:            "Empty URL generation",
			url:             "",
			regenKey:        false,
			countRandLetter: 0,
			wantHashLen:     32,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Generate(tc.url, tc.regenKey, tc.countRandLetter)
			if err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}

			// Проверяем длину получившейся строки
			if len(got) != tc.wantHashLen {
				t.Errorf("Generate() length = %d, want %d", len(got), tc.wantHashLen)
			}
		})
	}
}
