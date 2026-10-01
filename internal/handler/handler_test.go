package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
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
		name   string
		method string
		want   want
	}{
		{
			name:   "[Negative] Проверка на http метод GET -> POST",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    `Неверный http метод запроса`,
				contentType: "text/plain",
			},
		},
		{
			name:   "[Negative] Проверка на Content-Type",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    `Неверный Content-Type`,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			response := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, "/", nil)

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

	flagEnv := &FlagEnv{
		FlagRunAddr:           "localhost:8080",
		FlagRunShorternerAddr: "http://localhost:8080",
	}

	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "[Negative] Проверка на http метод GET -> POST",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    `Неверный http метод запроса`,
				contentType: "text/plain",
			},
		},
		{
			name:   "[Negative] Проверка на Content-Type",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    `Неверный Content-Type`,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			response := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, "/", nil)

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
