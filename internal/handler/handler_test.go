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
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "[Negative] Проверка на http метод GET -> POST",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    `Неверный http метод запроса`,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", nil)

			MainPage(response, request)

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
