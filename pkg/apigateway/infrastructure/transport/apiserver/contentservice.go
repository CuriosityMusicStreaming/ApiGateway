package apiserver

import (
	"apigateway/api/apigateway"
	contentserviceapi "apigateway/api/contentservice"
	"context"
	"github.com/pkg/errors"
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
	_, err = server.contentServiceClient.AddContent(ctx, &contentserviceapi.AddContentRequest{
		Name:             req.Name,
		Type:             contentType,
		AvailabilityType: availabilityType,
		UserToken:        serializedToken,
	})
	if err != nil {
		return nil, err
	}

	return &apigateway.AddContentResponse{}, nil
}

var apiServiceToContentServiceContentTypeMap = map[apigateway.ContentType]contentserviceapi.ContentType{
	apigateway.ContentType_Song:    contentserviceapi.ContentType_Song,
	apigateway.ContentType_Podcast: contentserviceapi.ContentType_Podcast,
}

var apiServiceToContentServiceAvailabilityTypeMap = map[apigateway.ContentAvailabilityType]contentserviceapi.ContentAvailabilityType{
	apigateway.ContentAvailabilityType_Public:  contentserviceapi.ContentAvailabilityType_Public,
	apigateway.ContentAvailabilityType_Private: contentserviceapi.ContentAvailabilityType_Private,
}
