package apiserver

import (
	"context"

	"github.com/pkg/errors"

	"apigateway/api/apigateway"
	userserviceapi "apigateway/api/userservice"
)

var (
	ErrUnknownRole = errors.New("unknown role")
)

func (server *apiGatewayServer) AddUser(ctx context.Context, req *apigateway.AddUserRequest) (*apigateway.AddUserResponse, error) {
	userRole, ok := apiServiceToContentServiceRoleMap[req.Role]
	if !ok {
		return nil, ErrUnknownRole
	}
	resp, err := server.userServiceClient.AddUser(ctx, &userserviceapi.AddUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     userRole,
	})

	if err != nil {
		return nil, err
	}

	return &apigateway.AddUserResponse{UserId: resp.UserId}, nil
}

var apiServiceToContentServiceRoleMap = map[apigateway.UserRole]userserviceapi.UserRole{
	apigateway.UserRole_LISTENER: userserviceapi.UserRole_LISTENER,
	apigateway.UserRole_CREATOR:  userserviceapi.UserRole_CREATOR,
}
