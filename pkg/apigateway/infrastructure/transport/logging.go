package transport

import (
	logger "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		h.ServeHTTP(writer, request)

		logger.WithFields(logger.Fields{
			"duration": time.Since(now),
			"method":   request.Method,
			"url":      request.URL,
		}).Info("request finished")
	})
}
