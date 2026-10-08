package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMainPage(t *testing.T) {

	type want struct {
		method      string
		code        int
		contentType string
	}

	store := service.NewStorage()

	cfg := &FlagEnv{
		FlagRunAddr:           "localhost:8080",
		FlagRunShorternerAddr: "http://localhost:8080",
	}

	h := NewHandler(store, cfg)

	tests := []struct {
		name    string
		method  string
		request string
		body    string
		want    want
	}{
		{
			name:    "Верный ответ",
			method:  http.MethodPost,
			request: "/",
			body:    "http://n1qttzvbn3.yandex/arqay",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			response := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
			request.Header.Set("Content-Type", "text/plain")

			h.MainPage(response, request)

			res := response.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			responseKey, err := url.Parse(string(resBody))
			if err != nil {
				require.NoError(t, err)
			}

			storageURL, err := h.storage.Find(strings.Trim(responseKey.Path, "/"))

			if err != nil {
				require.NoError(t, err)
			}

			assert.Equal(t, test.body, storageURL)

			fmt.Printf("%d %s", res.StatusCode, string(resBody))
		})
	}
}

func TestPageById(t *testing.T) {

	type want struct {
		method      string
		code        int
		contentType string
	}

	store := service.NewStorage()

	cfg := &FlagEnv{
		FlagRunAddr:           "localhost:8080",
		FlagRunShorternerAddr: "http://localhost:8080",
	}

	h := NewHandler(store, cfg)

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Верный ответ",
			method: http.MethodGet,
			url:    "http://n1qttzvbn3.yandex/arqay",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			key, err := h.storage.Save(test.url)
			if err != nil {
				t.Errorf("error prepare url %s", test.url)
			}

			r := chi.NewRouter()
			r.Get("/{id}", h.PageByID)

			request := httptest.NewRequest(test.method, cfg.FlagRunShorternerAddr+"/"+key, nil)
			response := httptest.NewRecorder()

			r.ServeHTTP(response, request)

			res := response.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			fmt.Printf("%d %s", res.StatusCode, string(resBody))
		})
	}
}
