package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainPage(t *testing.T) {

	type want struct {
		method      string
		code        int
		response    string
		contentType string
	}

	flagEnv := &FlagEnv{
		FlagRunAddr:           "localhost:8080",
		FlagRunShorternerAddr: "http://localhost:8080",
	}

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
				response:    "http://localhost:8080/4b90906a4f8dbe74fca39107f330b069",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			response := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
			request.Header.Set("Content-Type", "text/plain")

			handler := MainPage(flagEnv)
			handler.ServeHTTP(response, request)

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

func TestPageById(t *testing.T) {

	type want struct {
		method      string
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:    "Верный ответ",
			method:  http.MethodGet,
			request: "/4b90906a4f8dbe74fca39107f330b069",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "http://n1qttzvbn3.yandex/arqay",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			r := chi.NewRouter()
			r.Get("/{id}", PageByID)

			request := httptest.NewRequest(test.method, test.request, nil)
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
