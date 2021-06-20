package apiserver

import (
	"context"

	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"apigateway/api/apigateway"
	authenticationserviceapi "apigateway/api/authenticationservice"
	contentserviceapi "apigateway/api/contentservice"
	playlistserviceapi "apigateway/api/playlistservice"
	userserviceapi "apigateway/api/userservice"
	"apigateway/pkg/apigateway/infrastructure/auth"
)

const (
	authorizationHeaderName = "authorization"
)

func NewAPIGatewayServer(
	contentServiceClient contentserviceapi.ContentServiceClient,
	userServiceClient userserviceapi.UserServiceClient,
	playlistServiceClient playlistserviceapi.PlayListServiceClient,
	authenticationServiceClient authenticationserviceapi.AuthenticationServiceClient,
	authenticationService auth.AuthenticationService,
	userDescriptorSerializer commonauth.UserDescriptorSerializer,
) apigateway.APIGatewayServer {
	return &apiGatewayServer{
		contentServiceClient:        contentServiceClient,
		userServiceClient:           userServiceClient,
		playlistServiceClient:       playlistServiceClient,
		authenticationServiceClient: authenticationServiceClient,
		authenticationService:       authenticationService,
		userDescriptorSerializer:    userDescriptorSerializer,
	}
}

type apiGatewayServer struct {
	contentServiceClient        contentserviceapi.ContentServiceClient
	userServiceClient           userserviceapi.UserServiceClient
	playlistServiceClient       playlistserviceapi.PlayListServiceClient
	authenticationServiceClient authenticationserviceapi.AuthenticationServiceClient
	authenticationService       auth.AuthenticationService
	userDescriptorSerializer    commonauth.UserDescriptorSerializer
}

func (server *apiGatewayServer) AuthenticateUser(ctx context.Context, req *apigateway.AuthenticateUserRequest) (*apigateway.AuthenticateUserResponse, error) {
	resp, err := server.authenticationServiceClient.AuthenticateUser(ctx, &authenticationserviceapi.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	err = grpc.SendHeader(ctx, metadata.Pairs(authorizationHeaderName, resp.UserID))

	return &apigateway.AuthenticateUserResponse{UserID: resp.UserID}, err
}

func (server *apiGatewayServer) authenticateUser(ctx context.Context) (commonauth.UserDescriptor, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return commonauth.UserDescriptor{}, status.Errorf(codes.PermissionDenied, "missing context")
	}

	authorizationHeaders := md.Get(authorizationHeaderName)
	if len(authorizationHeaders) == 0 {
		return commonauth.UserDescriptor{}, status.Errorf(codes.PermissionDenied, "missing authentication header")
	}
	header := authorizationHeaders[0]

	token, err := server.authenticationService.ReceiveUserID(header)
	if err != nil {
		return commonauth.UserDescriptor{}, err
	}

	return token, nil
}
