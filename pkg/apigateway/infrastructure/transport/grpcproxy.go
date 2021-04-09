package transport

import (
	"context"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/activity"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
)

type GRPCProxy interface {
	RegisterConfiguration(targetAddress string, handlePrefix string, register func(mux *runtime.ServeMux, conn *grpc.ClientConn) error) error
	io.Closer
	http.Handler
}

func NewGRPCProxy(router *mux.Router) GRPCProxy {
	return &grpcProxy{
		router: router,
	}
}

type grpcProxy struct {
	router  *mux.Router
	closers []io.Closer
}

func (c *grpcProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.router.ServeHTTP(w, req)
}

func (c *grpcProxy) RegisterConfiguration(targetAddress string, handlePrefix string, register func(mux *runtime.ServeMux, conn *grpc.ClientConn) error) error {
	grpcGatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(NewUnaryInterceptor()),
	}

	conn, err := grpc.Dial(targetAddress, opts...)
	if err != nil {
		return err
	}

	err = register(grpcGatewayMux, conn)
	if err != nil {
		return err
	}
	c.closers = append(c.closers, conn)

	c.router.PathPrefix(handlePrefix).Handler(grpcGatewayMux)

	return nil
}

func (c *grpcProxy) Close() error {
	var err error
	for _, closer := range c.closers {
		err2 := closer.Close()
		if err2 != nil {
			err = err2
		}
	}
	return err
}

func NewUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		activityID := activity.NewActivityID()
		md := metadata.New(map[string]string{"activityID": activityID.String()})

		oldMd, _ := metadata.FromOutgoingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(oldMd, md))

		err := invoker(ctx, method, req, reply, cc, opts...)

		entry := logger.WithFields(logger.Fields{
			"activityID":   activityID.String(),
			"proxy_target": cc.Target(),
			"method":       method,
		})

		if err != nil {
			entry.Error("call failed")
		} else {
			entry.Info("call finished")
		}

		return err
	}
}
