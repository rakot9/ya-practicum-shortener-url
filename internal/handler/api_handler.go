package handler

import (
	"crypto/hmac"
	// "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	// "encoding/hex"
	"fmt"
	"strconv"

	"bytes"
	"encoding/json"
	"errors"

	"net/http"
	"net/url"
	"strings"

	"github.com/rakot9/ya-practicum-shortener-url/internal/model"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service/txctx"
	"go.uber.org/zap"
)

// Секретный ключ для подписи (в реальном проекте берите из переменных окружения)
var secretKey = []byte("super-secret-key-at-least-32-bytes-long!")

const cookieName = "user_session_id"

// verifyAndExtract проверяет подпись и возвращает исходное значение
func verifyAndExtract(signedValue string, secret []byte) (string, error) {
	parts := strings.Split(signedValue, ".")
	if len(parts) != 2 {
		return "", errors.New("неверный формат куки")
	}

	valueBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("ошибка декодирования значения")
	}

	sigBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("ошибка декодирования подписи")
	}

	// Вычисляем ожидаемую подпись
	mac := hmac.New(sha256.New, secret)
	mac.Write(valueBytes)
	expectedSig := mac.Sum(nil)

	// Используем hmac.Equal для защиты от атак по времени (timing attacks)
	if !hmac.Equal(sigBytes, expectedSig) {
		return "", errors.New("подпись недействительна")
	}

	return string(valueBytes), nil
}

// signValue формирует строку формата "base64(value).base64(signature)"
func signValue(value string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(value))
	signature := mac.Sum(nil)

	encodedValue := base64.URLEncoding.EncodeToString([]byte(value))
	encodedSig := base64.URLEncoding.EncodeToString(signature)

	return encodedValue + "." + encodedSig
}

func (h *Handler) APIShortURL(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		http.Error(res, "Неверный http метод запроса", http.StatusBadRequest)
		return
	}

	var shortURL model.ShortURLRequest
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

	// Cookie
	cookie, err := req.Cookie(cookieName)
	if err == nil {
		// Кука найдена, проверяем подпись
		userID, err := verifyAndExtract(cookie.Value, secretKey)
		if err == nil {
			fmt.Fprintf(res, "С возвращением! Ваш проверенный ID: %s\n", userID)
			return
		}
		// Если подпись невалидна, логируем инцидент и перевыпускаем куку
		fmt.Printf("Попытка подделки куки: %v\n", err)
	}

	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &shortURL); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.urlService.ShortenURL(req.Context(), shortURL.URL, "")

	isSetOriginalURLInStorage := false
	if err != nil {
		if errors.Is(err, service.ErrOriginalURLAlreadyExist) {
			isSetOriginalURLInStorage = true
		}
		if errors.Is(err, service.ErrCountAttemptExceed) {
			h.logger.Error("count attempt (10) exceed for generate key. URL", zap.Error(err))
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

	var shortURLResponse model.ShortURLResponse
	shortURLResponse.URL = url

	resp, err := json.Marshal(shortURLResponse)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if isSetOriginalURLInStorage {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	user, err := h.urlService.GetUser()
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	signedCookieValue := signValue(strconv.Itoa(user), secretKey)

	http.SetCookie(res, &http.Cookie{
		Name:     cookieName,
		Value:    signedCookieValue,
		Path:     "/",
		MaxAge:   86400 * 30, // 30 дней
		HttpOnly: true,       // Защита от доступа через JS (XSS)
		Secure:   true,       // Передача только по HTTPS (для localhost можно временно отключить)
		SameSite: http.SameSiteLaxMode,
	})

	res.Write(resp)
}

func (h *Handler) APIShortURLBatch(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		http.Error(res, "Неверный http метод запроса", http.StatusBadRequest)
		return
	}

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
	var batchRequestItems []model.BatchRequestItem
	if err = json.Unmarshal(buf.Bytes(), &batchRequestItems); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if len(batchRequestItems) == 0 {
		http.Error(res, "Список URL пуст", http.StatusBadRequest)
		return
	}

	var BatchShortURLResponseItems []model.BatchResponseItem

	BatchShortURLResponseItems, err = h.urlService.ShortenBatch(req.Context(), batchRequestItems, h.logger, h.flagEnv.FlagRunShorternerAddr)

	var isSetOriginalURLInStorage = false

	if err != nil {
		var batchErr *service.BatchError

		if errors.As(err, &batchErr) {
			isSetOriginalURLInStorage = batchErr.IsOrigURLInStorage

			if errors.Is(batchErr, service.ErrInternalServerError) || errors.Is(batchErr, service.ErrServerError) {
				http.Error(res, batchErr.Err.Error(), http.StatusInternalServerError)
				return
			}

			if errors.Is(batchErr, service.ErrBadRequest) {
				http.Error(res, batchErr.Err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			if errors.Is(err, txctx.ErrTransaction) {
				http.Error(res, "internal server error", http.StatusInternalServerError)
				return
			} else {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	resp, err := json.Marshal(BatchShortURLResponseItems)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if isSetOriginalURLInStorage {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	res.Write(resp)
}

func (h *Handler) APIUserURLs(res http.ResponseWriter, req *http.Request) {
}
