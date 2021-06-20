package auth

import (
	"strings"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Type string

const (
	TypeBearer = "Bearer"
)

var (
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidToken               = errors.New("invalid token")
)

type AuthenticationService interface {
	ReceiveUserID(header string) (auth.UserDescriptor, error)
}

func NewAuthenticationService(authType Type) AuthenticationService {
	return &authenticationService{authType: authType}
}

type authenticationService struct {
	authType Type
}

func (service *authenticationService) ReceiveUserID(header string) (auth.UserDescriptor, error) {
	if !strings.HasPrefix(header, string(service.authType)) {
		return auth.UserDescriptor{}, ErrInvalidAuthorizationHeader
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return auth.UserDescriptor{}, ErrInvalidAuthorizationHeader
	}

	token := parts[1]
	userID, err := uuid.Parse(token)
	if err != nil {
		return auth.UserDescriptor{}, ErrInvalidToken
	}

	return auth.UserDescriptor{UserID: userID}, nil
}
