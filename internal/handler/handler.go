package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"io"
	"net/http"
	"strings"
)

type FlagEnv struct {
	FlagRunAddr           string
	FlagRunShorternerAddr string
}

func MainPage(flagEnv *FlagEnv) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
				http.Error(res, "Неверный http метод запроса", http.StatusBadRequest)
			}

			req.Body = http.MaxBytesReader(res, req.Body, 1048576)

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(res, "Превышен размер тела сообщения: "+err.Error(), http.StatusBadRequest)
				return
			}

			bodyText := string(bodyBytes)

			shortURL := bodyText

			res.Header().Set("Content-Type", "text/plain")
			res.WriteHeader(http.StatusCreated)

			suffix, err := service.Save(shortURL)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}

			res.Write([]byte(flagEnv.FlagRunShorternerAddr + "/" + suffix))
		} else {
			http.Error(res, "Неверный http метод запроса Main", http.StatusBadRequest)
			return
		}
	}
}

func PageByID(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {

		id := chi.URLParam(req, "id")
		url, err := service.Find(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "Неверный http метод запроса ById", http.StatusBadRequest)
	}
}
