package config

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Flags
	}{
		{
			name: "Проверка дефолтных значений",
			args: []string{},
			expected: Flags{
				FlagRunAddr:           "localhost:8080",
				FlagRunShorternerAddr: "http://localhost:8080",
				FlagFileStoragePath:   "",
			},
		},
		{
			name: "Проверка флага -a",
			args: []string{"-a=localhost:8081"},
			expected: Flags{
				FlagRunAddr:           "localhost:8081",
				FlagRunShorternerAddr: "http://localhost:8080",
				FlagFileStoragePath:   "",
			},
		},
		{
			name: "Проверка флага -b",
			args: []string{"-b=http://localhost:8081"},
			expected: Flags{
				FlagRunAddr:           "localhost:8080",
				FlagRunShorternerAddr: "http://localhost:8081",
				FlagFileStoragePath:   "",
			},
		},
		{
			name: "Проверка флага -a -b",
			args: []string{"-a=localhost:8081", "-b=http://localhost:8081"},
			expected: Flags{
				FlagRunAddr:           "localhost:8081",
				FlagRunShorternerAddr: "http://localhost:8081",
				FlagFileStoragePath:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := ParseFlags(tt.args)

			if flags.FlagRunAddr != tt.expected.FlagRunAddr ||
				flags.FlagRunShorternerAddr != tt.expected.FlagRunShorternerAddr ||
				flags.FlagFileStoragePath != tt.expected.FlagFileStoragePath {
				t.Errorf("got %+v, want %+v", flags, tt.expected)
			}
		})
	}
}
