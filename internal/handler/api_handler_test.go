package handler

import (
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_ApiShortURL(t *testing.T) {
	flags := config.Flags{
		FlagRunAddr:           "localhost:8080",
		FlagRunShorternerAddr: "http://localhost:8080",
		FlagFileStoragePath:   "shortener.txt",
		FlagInMemoryStorage:   true,
	}

	repo := initRepository(flags)
	h := setupTestHandler(flags, repo)

	tests := []struct {
		name         string
		contentType  string
		body         string
		wantStatus   int
		wantBodyPart string
	}{
		{
			name:         "Успешное сокращение URL",
			contentType:  "application/json",
			body:         "{\"url\":\"https://yandex.ru\"}",
			wantStatus:   http.StatusCreated,
			wantBodyPart: "http://localhost:8080/",
		},
		{
			name:         "Неверный Content-Type",
			contentType:  "plain/text",
			body:         "{\"url\":\"https://yandex.ru\"}",
			wantStatus:   http.StatusBadRequest,
			wantBodyPart: "Неверный http метод запроса",
		},
		{
			name:         "Пустой тип контента",
			contentType:  "",
			body:         "{\"url\":\"https://yandex.ru\"}",
			wantStatus:   http.StatusBadRequest,
			wantBodyPart: "Неверный http метод запроса",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			if tt.contentType != "" {
				request.Header.Set("Content-Type", tt.contentType)
			}

			recorder := httptest.NewRecorder()
			h.APIShortURL(recorder, request)

			res := recorder.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			assert.Contains(t, string(resBody), tt.wantBodyPart)

			if tt.wantStatus == http.StatusCreated {
				assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			}
		})
	}
}
