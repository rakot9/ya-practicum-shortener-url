package handler

import (
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"io"
	"net/http"
)

const URL = "http://localhost:8080/"

func MainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if "text/plain" != req.Header.Get("Content-Type") {
			http.Error(res, "Неверный http метод запроса", http.StatusMethodNotAllowed)
		}

		req.Body = http.MaxBytesReader(res, req.Body, 1048576)

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Превышен размер тела сообщения: "+err.Error(), http.StatusBadRequest)
			return
		}

		bodyText := string(bodyBytes)

		shortUrl := bodyText

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)

		suffix, err := service.Save(shortUrl)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.Write([]byte(URL + suffix))
	} else {
		http.Error(res, "Неверный http метод запроса", http.StatusMethodNotAllowed)
	}
}

func PageById(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {

		url, err := service.Find(req.PathValue("id"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "Неверный http метод запроса", http.StatusMethodNotAllowed)
	}
}
