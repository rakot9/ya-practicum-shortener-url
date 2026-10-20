package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type WrapResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func ZapLogger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Оборачиваем стандартный ResponseWriter для перехвата статус-кода
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Передаем управление дальше по цепочке
			next.ServeHTTP(ww, r)

			// Вычисляем время выполнения запроса
			duration := time.Since(start)

			// Логируем структурированные данные с помощью Zap
			log.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
			)
		})
	}
}
