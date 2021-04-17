package transport

import (
	"context"
	"fmt"
	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/activity"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
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

func NewLoggerServerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		activityID := activity.NewActivityID()
		md := metadata.New(map[string]string{"activityID": activityID.String()})

		oldMd, _ := metadata.FromOutgoingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(oldMd, md))

		start := time.Now()

		resp, err = handler(ctx, req)

		fields := log.Fields{
			"activityID": activityID.String(),
			"args":       req,
			"duration":   fmt.Sprintf("%v", time.Since(start)),
			"method":     getGRPCMethodName(info),
		}

		entry := logger.WithFields(fields)
		if err != nil {
			entry.Error(err, "call failed")
		} else {
			entry.Info("call finished")
		}

		return resp, err
	}
}

func getGRPCMethodName(info *grpc.UnaryServerInfo) string {
	method := info.FullMethod
	return method[strings.LastIndex(method, "/")+1:]
}
