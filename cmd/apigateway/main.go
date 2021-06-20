package main

import (
	"context"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	jsonlog "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/logger"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/server"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"apigateway/api/apigateway"
	"apigateway/api/authenticationservice"
	contentserviceapi "apigateway/api/contentservice"
	playlistserviceapi "apigateway/api/playlistservice"
	userserviceapi "apigateway/api/userservice"
	"apigateway/pkg/apigateway/infrastructure/auth"
	"apigateway/pkg/apigateway/infrastructure/transport"
	"apigateway/pkg/apigateway/infrastructure/transport/apiserver"
)

var appID = "UNKNOWN"

func main() {
	logger, err := initLogger()
	if err != nil {
		stdlog.Fatal("failed to initialize logger")
	}

	config, err := parseEnv()
	if err != nil {
		logger.FatalError(err)
	}

	err = runService(config, logger)
	if err == server.ErrStopped {
		logger.Info("service is successfully stopped")
	} else if err != nil {
		logger.FatalError(err)
	}
}

func runService(config *config, logger log.MainLogger) error {
	stopChan := make(chan struct{})
	listenForKillSignal(stopChan)

	serverHub := server.NewHub(stopChan)

	apiServer, err := initAPIServer(config)
	if err != nil {
		return err
	}

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(transport.NewLoggerServerInterceptor(logger)))
	apigateway.RegisterAPIGatewayServer(baseServer, apiServer)

	serverHub.AddServer(server.NewGrpcServer(
		baseServer,
		server.GrpcServerConfig{ServeAddress: config.ServeGRPCAddress},
		logger,
	))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var httpServer *http.Server

	serverHub.AddServer(&server.FuncServer{
		ServeImpl: func() error {
			grpcGatewayMux := runtime.NewServeMux()
			opts := []grpc.DialOption{grpc.WithInsecure()}
			err := apigateway.RegisterAPIGatewayHandlerFromEndpoint(ctx, grpcGatewayMux, config.ServeGRPCAddress, opts)
			if err != nil {
				return err
			}

			router := mux.NewRouter()
			router.PathPrefix("/api/").Handler(grpcGatewayMux)

			router.HandleFunc("/resilience/ready", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
			}).Methods(http.MethodGet)

			httpServer = &http.Server{
				Handler:      router,
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

func initLogger() (log.MainLogger, error) {
	return jsonlog.NewLogger(&jsonlog.Config{AppName: appID}), nil
}

func initAPIServer(config *config) (apigateway.APIGatewayServer, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	contentServiceClient, err := initContentServiceClient(opts, config)
	if err != nil {
		return nil, err
	}

	userServiceClient, err := initUserServiceClient(opts, config)
	if err != nil {
		return nil, err
	}

	playlistServiceClient, err := initPlaylistServiceClient(opts, config)
	if err != nil {
		return nil, err
	}

	authenticationServiceClient, err := initAuthenticationServiceClient(opts, config)
	if err != nil {
		return nil, err
	}

	return apiserver.NewAPIGatewayServer(
		contentServiceClient,
		userServiceClient,
		playlistServiceClient,
		authenticationServiceClient,
		auth.NewAuthenticationService(auth.TypeBearer),
		commonauth.NewUserDescriptorSerializer(),
	), nil
}

func initContentServiceClient(commonOpts []grpc.DialOption, config *config) (contentserviceapi.ContentServiceClient, error) {
	conn, err := grpc.Dial(config.ContentServiceGRPCAddress, commonOpts...)
	if err != nil {
		return nil, err
	}

	return contentserviceapi.NewContentServiceClient(conn), nil
}

func initUserServiceClient(commonOpts []grpc.DialOption, config *config) (userserviceapi.UserServiceClient, error) {
	conn, err := grpc.Dial(config.UserServiceGRPCAddress, commonOpts...)
	if err != nil {
		return nil, err
	}

	return userserviceapi.NewUserServiceClient(conn), nil
}

func initPlaylistServiceClient(commonOpts []grpc.DialOption, config *config) (playlistserviceapi.PlayListServiceClient, error) {
	conn, err := grpc.Dial(config.PlaylistServiceGRPCAddress, commonOpts...)
	if err != nil {
		return nil, err
	}

	return playlistserviceapi.NewPlayListServiceClient(conn), nil
}

func initAuthenticationServiceClient(commonOpts []grpc.DialOption, config *config) (authenticationservice.AuthenticationServiceClient, error) {
	conn, err := grpc.Dial(config.AuthenticationServiceGRPCAddress, commonOpts...)
	if err != nil {
		return nil, err
	}

	return authenticationservice.NewAuthenticationServiceClient(conn), nil
}
