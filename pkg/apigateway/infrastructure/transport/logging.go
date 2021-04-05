package transport

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		h.ServeHTTP(writer, request)

		logger.WithFields(logger.Fields{
			"duration": fmt.Sprintf("%v", time.Since(start)),
			"method":   request.Method,
			"url":      request.RequestURI,
		}).Info("request finished")
	})
}
