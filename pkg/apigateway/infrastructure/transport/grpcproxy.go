package transport

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
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
	opts := []grpc.DialOption{grpc.WithInsecure()}
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
