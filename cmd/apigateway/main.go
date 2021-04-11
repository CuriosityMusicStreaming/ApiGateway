package main

import (
	contentserviceapi "apigateway/api/contentservice"
	userserviceapi "apigateway/api/userservice"
	"apigateway/pkg/apigateway/infrastructure/transport"
	"context"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

var appID = "UNKNOWN"

func main() {
	logger.SetFormatter(&logger.JSONFormatter{})

	config, err := parseEnv()
	if err != nil {
		logger.Fatal(err)
	}

	err = runService(config)
	if err == server.ErrStopped {
		logger.Info("service is successfully stopped")
	} else if err != nil {
		logger.Fatal(err)
	}
}

func runService(config *config) error {
	stopChan := make(chan struct{})
	listenForKillSignal(stopChan)

	serverHub := server.NewHub(stopChan)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var httpServer *http.Server

	serverHub.AddServer(&server.FuncServer{
		ServeImpl: func() error {
			router := mux.NewRouter()

			// Implement healthcheck for Kubernetes
			router.HandleFunc("/resilience/ready", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
			}).Methods(http.MethodGet)

			grpcProxy := transport.NewGRPCProxy(router)

			err := registerContentService(ctx, grpcProxy, config)
			if err != nil {
				return err
			}

			err = registerUserService(ctx, grpcProxy, config)
			if err != nil {
				return err
			}

			httpServer = &http.Server{
				Handler:      transport.NewLoggingMiddleware(grpcProxy),
				Addr:         config.ServeRESTAddress,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			logger.Info("REST server started")
			return httpServer.ListenAndServe()
		},
		StopImpl: func() error {
			cancel()
			return httpServer.Shutdown(context.Background())
		},
	})

	return serverHub.Run()
}

func listenForKillSignal(stopChan chan<- struct{}) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		stopChan <- struct{}{}
	}()
}

func registerContentService(ctx context.Context, proxy transport.GRPCProxy, config *config) error {
	return proxy.RegisterConfiguration(
		config.ContentServiceGRPCAddress,
		"/api/cs/",
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return contentserviceapi.RegisterContentServiceHandler(ctx, mux, conn)
		},
	)
}

func registerUserService(ctx context.Context, proxy transport.GRPCProxy, config *config) error {
	return proxy.RegisterConfiguration(
		config.UserServiceGRPCAddress,
		"/api/us/",
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return userserviceapi.RegisterUserServiceHandler(ctx, mux, conn)
		},
	)
}
