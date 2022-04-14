package login_facade

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_repository/login_implement"
)

var loginTokenRepository LoginTokenRepositoryInterface

func GetLoginTokenRepository() LoginTokenRepositoryInterface {
	if loginTokenRepository == nil {
		loginTokenRepository = &login_implement.LoginTokenRepositoryImpl{}
	}
	return loginTokenRepository
}

type LoginTokenRepositoryInterface interface {
	Save(ctx context.Context, token *login_entity.LoginToken) error

	ExistToken(ctx context.Context, token string) (bool, error)
	GetUserID(ctx context.Context, token string) string
	GetTokenCreatedAt(ctx context.Context, token string) (int64, error)
	SetTokenExpiration(ctx context.Context, token string, duration time.Duration)
	GetTokenExpiration(ctx context.Context, token string) time.Duration
}
