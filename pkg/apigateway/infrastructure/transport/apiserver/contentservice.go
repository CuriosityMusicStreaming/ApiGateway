package apiserver

import (
	"apigateway/api/apigateway"
	contentserviceapi "apigateway/api/contentservice"
	"context"
	"github.com/pkg/errors"
	context2 "golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrUnknownContentType             = errors.New("unknown content type")
	ErrUnknownContentAvailabilityType = errors.New("unknown content availability type")
)

func (server *apiGatewayServer) AddContent(ctx context.Context, req *apigateway.AddContentRequest) (*apigateway.AddContentResponse, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	contentType, ok := apiServiceToContentServiceContentTypeMap[req.Type]
	if !ok {
		return nil, ErrUnknownContentType
	}

	availabilityType, ok := apiServiceToContentServiceAvailabilityTypeMap[req.AvailabilityType]
	if !ok {
		return nil, ErrUnknownContentAvailabilityType
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	resp, err := server.contentServiceClient.AddContent(ctx, &contentserviceapi.AddContentRequest{
		Name:             req.Name,
		Type:             contentType,
		AvailabilityType: availabilityType,
		UserToken:        serializedToken,
	})
	if err != nil {
		return nil, err
	}

	return &apigateway.AddContentResponse{ContentID: resp.ContentID}, nil
}

func (server *apiGatewayServer) DeleteContent(ctx context.Context, req *apigateway.DeleteContentRequest) (*emptypb.Empty, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	_, err = server.contentServiceClient.DeleteContent(ctx, &contentserviceapi.DeleteContentRequest{
		ContentID: req.ContentID,
		UserToken: serializedToken,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (server *apiGatewayServer) SetContentAvailabilityType(ctx context.Context, req *apigateway.SetContentAvailabilityTypeRequest) (*emptypb.Empty, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	availabilityType, ok := apiServiceToContentServiceAvailabilityTypeMap[req.NewContentAvailabilityType]
	if !ok {
		return nil, ErrUnknownContentAvailabilityType
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	_, err = server.contentServiceClient.SetContentAvailabilityType(ctx, &contentserviceapi.SetContentAvailabilityTypeRequest{
		ContentID:                  req.ContentID,
		NewContentAvailabilityType: availabilityType,
		UserToken:                  serializedToken,
	})
	return &emptypb.Empty{}, err
}

func (server *apiGatewayServer) GetAuthorContent(ctx context2.Context, _ *apigateway.GetAuthorContentRequest) (*apigateway.GetAuthorContentResponse, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)

	resp, err := server.contentServiceClient.GetAuthorContent(ctx, &contentserviceapi.GetAuthorContentRequest{UserToken: serializedToken})
	if err != nil {
		return nil, err
	}

	res := make([]*apigateway.Content, 0, len(resp.Contents))
	for _, content := range resp.Contents {
		res = append(res, &apigateway.Content{
			ContentID:        content.ContentID,
			Name:             content.Name,
			AuthorID:         content.AuthorID,
			Type:             contentServiceContentTypeToApiServiceMap[content.Type],
			AvailabilityType: contentServiceAvailabilityTypeToApiServiceMap[content.AvailabilityType],
		})
	}

	return &apigateway.GetAuthorContentResponse{Contents: res}, nil
}

var apiServiceToContentServiceContentTypeMap = map[apigateway.ContentType]contentserviceapi.ContentType{
	apigateway.ContentType_Song:    contentserviceapi.ContentType_Song,
	apigateway.ContentType_Podcast: contentserviceapi.ContentType_Podcast,
}

var apiServiceToContentServiceAvailabilityTypeMap = map[apigateway.ContentAvailabilityType]contentserviceapi.ContentAvailabilityType{
	apigateway.ContentAvailabilityType_Public:  contentserviceapi.ContentAvailabilityType_Public,
	apigateway.ContentAvailabilityType_Private: contentserviceapi.ContentAvailabilityType_Private,
}

var contentServiceContentTypeToApiServiceMap = map[contentserviceapi.ContentType]apigateway.ContentType{
	contentserviceapi.ContentType_Song:    apigateway.ContentType_Song,
	contentserviceapi.ContentType_Podcast: apigateway.ContentType_Podcast,
}

var contentServiceAvailabilityTypeToApiServiceMap = map[contentserviceapi.ContentAvailabilityType]apigateway.ContentAvailabilityType{
	contentserviceapi.ContentAvailabilityType_Public:  apigateway.ContentAvailabilityType_Public,
	contentserviceapi.ContentAvailabilityType_Private: apigateway.ContentAvailabilityType_Private,
}
