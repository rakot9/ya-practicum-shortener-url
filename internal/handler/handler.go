package handler

import (
	"io"
	"net/http"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Привет!"))
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
		res.Write([]byte(shortUrl))
	} else {
		http.Error(res, "Неверный http метод запроса", http.StatusMethodNotAllowed)
	}
}

func pageById(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		shortUrl := "Orig"
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusTemporaryRedirect)
		res.Write([]byte(shortUrl))
	} else {
		http.Error(res, "Неверный http метод запроса", http.StatusMethodNotAllowed)
	}
}
