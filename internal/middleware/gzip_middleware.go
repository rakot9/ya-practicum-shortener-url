package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// 1. Обертка для чтения сжатого запроса
type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c *gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzWriter *gzip.Writer
	skipGzip bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if statusCode >= 300 && statusCode < 400 {
		w.skipGzip = true
		w.ResponseWriter.Header().Del("Content-Encoding")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.skipGzip {
		return w.ResponseWriter.Write(b)
	}

	// Инициализируем gzip.Writer только при первой записи тела ответа
	if w.gzWriter == nil {
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			return 0, err
		}
		w.gzWriter = gz
	}

	return w.gzWriter.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = &gzipReader{r: r.Body, zr: zr}
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			skipGzip:       false,
		}

		next.ServeHTTP(gzw, r)

		if gzw.gzWriter != nil {
			gzw.gzWriter.Close()
		}
	})
}
