package handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config/db"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type FlagEnv struct {
	FlagRunAddr           string
	FlagRunShorternerAddr string
}

type Handler struct {
	urlService *service.URLService
	flagEnv    *FlagEnv
	dbStruct   *db.DbStruct
}

func NewHandler(urlService *service.URLService, flagEnv *FlagEnv, dbStruct *db.DbStruct) *Handler {
	return &Handler{
		urlService: urlService,
		flagEnv:    flagEnv,
		dbStruct:   dbStruct,
	}
}

func (h *Handler) MainPage(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
		http.Error(res, "Неверный http метод запроса", http.StatusBadRequest)
		return
	}

	req.Body = http.MaxBytesReader(res, req.Body, 1048576)
	defer req.Body.Close()

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			http.Error(res, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(res, "Превышен размер тела сообщения: "+err.Error(), http.StatusBadRequest)
		return
	}

	bodyText := string(bodyBytes)

	URLtoShorten := bodyText

	key, err := h.urlService.ShortenURL(URLtoShorten, "")

	isSetOriginalURLInStorage := false

	if err != nil {
		if errors.Is(err, service.ErrOriginalURLAlreadyExist) {
			isSetOriginalURLInStorage = true
		}

		if errors.Is(err, service.ErrCountAttemptExceed) {
			slog.Error("count attempt (10) exceed for generate key. URL", slog.Any("error", err))
			http.Error(res, "server error", http.StatusInternalServerError)
			return
		}

		if !isSetOriginalURLInStorage {
			http.Error(res, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	url, err := url.JoinPath(h.flagEnv.FlagRunShorternerAddr, key)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")

	if isSetOriginalURLInStorage {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	res.Write([]byte(url))
}

func (h *Handler) PageByID(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	url, err := h.urlService.GetURL(id)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		slog.Error("error find key.", slog.Any("error", err))

		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func (h *Handler) PingDB(res http.ResponseWriter, req *http.Request) {
	if err := h.dbStruct.Pool.Ping(h.dbStruct.Ctx); err != nil {
		slog.Error("db not accessable", slog.Any("error", err))

		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	res.WriteHeader(http.StatusOK)
}
