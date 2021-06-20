package transport

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/activity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
