package auth

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
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
	ReceiveUserID(header string) (string, error)
}

func NewAuthenticationService(authType Type) AuthenticationService {
	return &authenticationService{authType: authType}
}

type authenticationService struct {
	authType Type
}

func (service *authenticationService) ReceiveUserID(header string) (string, error) {
	if !strings.HasPrefix(header, string(service.authType)) {
		return "", ErrInvalidAuthorizationHeader
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", ErrInvalidAuthorizationHeader
	}

	token := parts[1]
	_, err := uuid.Parse(token)
	if err != nil {
		return "", ErrInvalidToken
	}

	return token, nil
}
