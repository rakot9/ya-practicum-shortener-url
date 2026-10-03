package config

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "Проверка флага -a",
			args: []string{},
			want: []string{"localhost:8080", "http://localhost:8080"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseFlags(tt.args)

			if FlagRunAddr != tt.want[0] && FlagRunShorternerAddr != tt.want[1] {
				t.Error("ParseFlags with no arg")
			}
		})
	}
}
