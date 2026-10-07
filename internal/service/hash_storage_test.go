package service_test

import (
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"testing"
)

func TestURLStorage_Find(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "Проверка нахождения ключа",
			key:     "4b90906a4f8dbe74fca39107f330b069",
			want:    "http://n1qttzvbn3.yandex/arqay",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s handler.Storage = service.URLStorage{}
			key, _ := service.Generate(tt.want, false, 1)

			_, err := s.Save(tt.want, key)
			if err != nil {
				t.Errorf("error prepare url %s", tt.want)
			}

			got, gotErr := s.Find(tt.key)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Find() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Find() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
