package handler

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type StorageProtocol interface {
	Save(URL string) (string, error)
	Find(key string) (string, error)
}

type FlagEnv struct {
	FlagRunAddr           string
	FlagRunShorternerAddr string
}

type Handler struct {
	storage StorageProtocol
	flagEnv *FlagEnv
}

func NewHandler(storage StorageProtocol, flagEnv *FlagEnv) *Handler {
	return &Handler{
		storage: storage,
		flagEnv: flagEnv,
	}
}

func (h *Handler) MainPage(res http.ResponseWriter, req *http.Request) {
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

	URLtoShorten := bodyText

	key, err := h.storage.Save(URLtoShorten)
	if err != nil {
		http.Error(res, "server error", http.StatusInternalServerError)
		return
	}

	url, err := url.JoinPath(h.flagEnv.FlagRunShorternerAddr, key)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(url))
}

func (h *Handler) PageByID(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	url, err := h.storage.Find(id)
	if err != nil {
		slog.Error("error find key.", slog.Any("error", err))

		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
