package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rakot9/ya-practicum-shortener-url/internal/ctxkey"
	"go.uber.org/zap"
)

var secretKey = []byte("super-secret-key-at-least-32-bytes-long!")

const cookieName = "user_id"

// type contextKey string

// const userIDKey contextKey = "userID"

var ErrInvalidCookie = errors.New("невалидная подпись или формат cookie")
var ErrEmptyUserID = errors.New("cookie не содержит ID пользователя")

// generateUID генерирует криптографически стойкий уникальный ID (16 байт в hex)
func generateUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// signValue создает подписанную строку в формате Base64(HMAC + Value)
func signValue(value string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(value))
	signature := mac.Sum(nil) // 32 байта

	// Объединяем подпись и само значение
	raw := append(signature, []byte(value)...)
	return base64.URLEncoding.EncodeToString(raw)
}

// verifyAndExtract проверяет подпись и возвращает исходное значение
func verifyAndExtract(signedValue string, secret []byte) (string, error) {
	raw, err := base64.URLEncoding.DecodeString(signedValue)
	if err != nil {
		return "", ErrInvalidCookie
	}

	// Длина SHA-256 подписи всегда 32 байта
	if len(raw) <= sha256.Size {
		return "", ErrInvalidCookie
	}

	signature := raw[:sha256.Size]
	value := string(raw[sha256.Size:])

	// Вычисляем ожидаемую подпись
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	// Используем hmac.Equal для защиты от Timing Attacks
	if !hmac.Equal(signature, expectedSignature) {
		return "", ErrInvalidCookie
	}

	// Подпись верна, но проверяем, не пустой ли сам ID
	if strings.TrimSpace(value) == "" {
		return "", ErrEmptyUserID
	}

	return value, nil
}

// func GetUserID(ctx context.Context) (string, bool) {
// 	val, ok := ctx.Value(userIDKey).(string)
// 	return val, ok
// }

// EnsureUserCookieMiddleware проверяет куку, и если её нет или она подделана — выдает новую
func EnsureUserCookieMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userUUID string
			var needNewCookie bool

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				needNewCookie = true
			} else {
				// Кука есть, проверяем подпись
				extractedID, err := verifyAndExtract(cookie.Value, secretKey)
				if err != nil {
					// Если кука есть, подпись верна, но ID пустой — отдаем 401
					if errors.Is(err, ErrEmptyUserID) {
						http.Error(w, "Unauthorized: cookie does not contain user ID", http.StatusUnauthorized)
						return
					}
					// В остальных случаях (неверная подпись, битый base64) перевыпускаем куку
					needNewCookie = true
				} else {
					userUUID = extractedID
				}

				logger.Info("ZAP: Found user in cookie",
					zap.String("uuid", userUUID),
				)
			}

			// Если куки не было или она не прошла проверку
			if needNewCookie {
				/*
					new_user_uuid := uuid.New().String()
					new_user_email := new_user_uuid + "@test.ru"
					new_user_password, err := generateUID()

					if err != nil {
						http.Error(w, "Ошибка генерации ID", http.StatusInternalServerError)
						return
					}

					logger.Info("ZAP: Created new user",
						zap.String("uuid", new_user_uuid),
						zap.String("email", new_user_email),
					)

					userID, err := h.UrlService.CreateUser(r.Context(), new_user_uuid, new_user_email, new_user_password)

					if err != nil || userID == nil {
						logger.Error("Ошибка вставки пользователя в БД")
						http.Error(w, "internal server error", http.StatusInternalServerError)
						return
					}

					userUUID = new_user_uuid
				*/

				new_user_uuid := uuid.New().String()
				userUUID = new_user_uuid

				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    signValue(new_user_uuid, secretKey),
					Path:     "/",
					MaxAge:   86400 * 30, // 30 дней
					HttpOnly: true,
					Secure:   false, // В продакшене с HTTPS обязательно ставьте true!
					SameSite: http.SameSiteLaxMode,
				})
			}

			// Кладем валидный ID в контекст, чтобы хендлеры могли им пользоваться
			// ctx := context.WithValue(r.Context(), userIDKey, userUUID)
			ctx := ctxkey.SetUserID(r.Context(), userUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// import (
// 	// "crypto/hmac"
// 	// "crypto/rand"
// 	// "crypto/sha256"
// 	// "encoding/base64"
// 	// "encoding/hex"
// 	// "fmt"
// 	// "strconv"

// Секретный ключ для подписи (в реальном проекте берите из переменных окружения)
// var secretKey = []byte("super-secret-key-at-least-32-bytes-long!")

// const cookieName = "user_session_id"

// // verifyAndExtract проверяет подпись и возвращает исходное значение
// func verifyAndExtract(signedValue string, secret []byte) (string, error) {
// 	parts := strings.Split(signedValue, ".")
// 	if len(parts) != 2 {
// 		return "", errors.New("неверный формат куки")
// 	}

// 	valueBytes, err := base64.URLEncoding.DecodeString(parts[0])
// 	if err != nil {
// 		return "", errors.New("ошибка декодирования значения")
// 	}

// 	sigBytes, err := base64.URLEncoding.DecodeString(parts[1])
// 	if err != nil {
// 		return "", errors.New("ошибка декодирования подписи")
// 	}

// 	// Вычисляем ожидаемую подпись
// 	mac := hmac.New(sha256.New, secret)
// 	mac.Write(valueBytes)
// 	expectedSig := mac.Sum(nil)

// 	// Используем hmac.Equal для защиты от атак по времени (timing attacks)
// 	if !hmac.Equal(sigBytes, expectedSig) {
// 		return "", errors.New("подпись недействительна")
// 	}

// 	return string(valueBytes), nil
// }

// // signValue формирует строку формата "base64(value).base64(signature)"
// func signValue(value string, secret []byte) string {
// 	mac := hmac.New(sha256.New, secret)
// 	mac.Write([]byte(value))
// 	signature := mac.Sum(nil)

// 	encodedValue := base64.URLEncoding.EncodeToString([]byte(value))
// 	encodedSig := base64.URLEncoding.EncodeToString(signature)

// 	return encodedValue + "." + encodedSig
// }

// Handler
// Cookie
// cookie, err := req.Cookie(cookieName)
// if err == nil {
// 	// Кука найдена, проверяем подпись
// 	userID, err := verifyAndExtract(cookie.Value, secretKey)
// 	if err == nil {
// 		fmt.Fprintf(res, "С возвращением! Ваш проверенный ID: %s\n", userID)
// 		return
// 	}
// 	// Если подпись невалидна, логируем инцидент и перевыпускаем куку
// 	fmt.Printf("Попытка подделки куки: %v\n", err)
// }

// user, err := h.urlService.GetUser()
// if err != nil {
// 	http.Error(res, err.Error(), http.StatusUnauthorized)
// 	return
// }

// signedCookieValue := signValue(strconv.Itoa(user), secretKey)

// http.SetCookie(res, &http.Cookie{
// 	Name:     cookieName,
// 	Value:    signedCookieValue,
// 	Path:     "/",
// 	MaxAge:   86400 * 30, // 30 дней
// 	HttpOnly: true,       // Защита от доступа через JS (XSS)
// 	Secure:   true,       // Передача только по HTTPS (для localhost можно временно отключить)
// 	SameSite: http.SameSiteLaxMode,
// })
