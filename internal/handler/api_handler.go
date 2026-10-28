package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rakot9/ya-practicum-shortener-url/internal/model"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handler) APIShortURL(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		http.Error(res, "Неверный http метод запроса", http.StatusBadRequest)
		return
	}

	var shortURL model.ShortURL
	var buf bytes.Buffer

	req.Body = http.MaxBytesReader(res, req.Body, 1048576)
	defer req.Body.Close()

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			http.Error(res, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &shortURL); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.urlService.ShortenURL(shortURL.URL)

	if err != nil {
		if errors.Is(err, service.ErrCountAttemptExceed) {
			slog.Error("count attempt (10) exceed for generate key. URL", slog.Any("error", err))
			http.Error(res, "server error", http.StatusInternalServerError)
			return
		}
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	url, err := url.JoinPath(h.flagEnv.FlagRunShorternerAddr, key)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var shortURLResponse model.ShortURLResponse
	shortURLResponse.URL = url

	resp, err := json.Marshal(shortURLResponse)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	res.Write(resp)
}
