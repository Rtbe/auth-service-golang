package repository

import (
	"context"

	"example.com/auth-service-go/internal/entity"
)

//Token is an interface which abstracts interaction with databases that interacts with tokens
type Token interface {
	Insert(context.Context, *entity.TokenPair) error
	DeleteUserRefreshTokens(context.Context, string) error
	DeleteRefreshToken(context.Context, string, string) error
	IsUserInDB(context.Context, string) bool
	IsRefreshTokenInDB(context.Context, string) bool
	RefreshTokenSetIsUsed(context.Context, string) error
}
