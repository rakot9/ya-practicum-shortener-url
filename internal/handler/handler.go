package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"io"
	"net/http"
	"strings"
)

type Storage interface {
	Save(URL string) (string, error)
	Find(hash string) (string, error)
}

type FlagEnv struct {
	FlagRunAddr           string
	FlagRunShorternerAddr string
}

func MainPage(flagEnv *FlagEnv) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var s Storage = service.URLStorage{}

		if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
			http.Error(res, "Неверный http метод запроса", http.StatusBadRequest)
			return
		}

		req.Body = http.MaxBytesReader(res, req.Body, 1048576)

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Превышен размер тела сообщения: "+err.Error(), http.StatusBadRequest)
			return
		}

		bodyText := string(bodyBytes)

		shortURL := bodyText

		suffix, err := s.Save(shortURL)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(flagEnv.FlagRunShorternerAddr + "/" + suffix))
	}
}

func PageByID(res http.ResponseWriter, req *http.Request) {

	var s Storage = service.URLStorage{}

	id := chi.URLParam(req, "id")
	url, err := s.Find(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
