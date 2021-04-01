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
	ServeRESTAddress          string `envconfig:"serve_rest_address" default:":8001"`
	ContentServiceGRPCAddress string `envconfig:"content_service_grpc_address"`
}
