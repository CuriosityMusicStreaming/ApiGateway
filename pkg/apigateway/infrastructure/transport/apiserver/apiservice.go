package apiserver

import (
	"apigateway/api/apigateway"
	contentserviceapi "apigateway/api/contentservice"
	userserviceapi "apigateway/api/userservice"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeaderName = "authorization"
)

func NewApiGatewayServer(
	contentServiceClient contentserviceapi.ContentServiceClient,
	userServiceClient userserviceapi.UserServiceClient,
) apigateway.ApiGatewayServer {
	return &apiGatewayServer{
		contentServiceClient: contentServiceClient,
		userServiceClient:    userServiceClient,
	}
}

type apiGatewayServer struct {
	contentServiceClient contentserviceapi.ContentServiceClient
	userServiceClient    userserviceapi.UserServiceClient
}

func (server *apiGatewayServer) AuthenticateUser(ctx context.Context, req *apigateway.AuthenticateUserRequest) (*apigateway.AuthenticateUserResponse, error) {
	resp, err := server.userServiceClient.AuthenticateUser(ctx, &userserviceapi.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	err = grpc.SendHeader(ctx, metadata.Pairs(authorizationHeaderName, resp.UserID))

	return &apigateway.AuthenticateUserResponse{UserID: resp.UserID}, err
}

func (server *apiGatewayServer) authenticateUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.PermissionDenied, "missing authentication token")
	}

	fmt.Println(md)

	authorizationHeaders := md.Get(authorizationHeaderName)
	if len(authorizationHeaders) == 0 {
		return "", status.Errorf(codes.PermissionDenied, "missing authentication token")
	}
	token := authorizationHeaders[0]

	return token, nil
}
