package internalhttp

import (
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		logger.Info("Request handled:",
			zap.String("ip", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("protocol", r.Proto),
			zap.Int("status", lrw.statusCode),
			zap.String("latency", time.Since(start).String()),
		)
		/*
			* IP клиента;
			//* дата и время запроса;
			//* метод, path и версия HTTP;
			* код ответа;
			//* latency (время обработки запроса, посчитанное, например, с помощью middleware);
			* user agent, если есть.
		*/
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
