package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	ServeRESTAddress                 string `envconfig:"serve_rest_address" default:":8001"`
	ServeGRPCAddress                 string `envconfig:"serve_grpc_address" default:":8002"`
	ContentServiceGRPCAddress        string `envconfig:"content_service_grpc_address"`
	UserServiceGRPCAddress           string `envconfig:"user_service_grpc_address"`
	AuthenticationServiceGRPCAddress string `envconfig:"authentication_service_grpc_address"`
	PlaylistServiceGRPCAddress       string `envconfig:"playlist_service_grpc_address"`
}
