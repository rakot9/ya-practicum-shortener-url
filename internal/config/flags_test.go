package config

import (
	// "github.com/stretchr/testify/assert"
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
			args: []string{"-a", "localhost:8081"},
			want: []string{"localhost:8081"},
		},
		{
			name: "Проверка флага -b",
			args: []string{"-b", "http://localhost:8080"},
			want: []string{"http://localhost:8080"},
		},
		{
			name: "Проверка отсутствия флагов",
			args: []string{},
			want: []string{"localhost:8080", "http://localhost:8080"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := ParseFlags(tt.args)

			if len(tt.args) == 0 {
				if flags.FlagRunAddr != tt.want[0] || flags.FlagRunShorternerAddr != tt.want[1] {
					t.Error("ParseFlags error default values")
				}
			}

			if len(tt.args) == 2 {
				t.Logf("### Flags %+v", flags)
				if tt.args[0] == "-a" {
					if flags.FlagRunAddr != tt.want[0] {
						t.Error("ParseFlags error default values for 1 arguments with -a argument")
					}
				}
				if tt.args[0] == "-b" {
					if flags.FlagRunShorternerAddr != tt.want[0] {
						t.Error("ParseFlags error default values for 1 arguments with -b argument")
					}
				}
			}

			// t.Logf("### Flags %+v", flags)
			// if flags.FlagRunAddr != tt.want[0] || flags.FlagRunShorternerAddr != tt.want[1] {
			// 	t.Error("ParseFlags with no arg")
			// }
		})
	}
}
