package login_facade

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_repository/login_implement"
)

var userRepository UserRepositoryInterface

func GetUserRepository() UserRepositoryInterface {
	if userRepository == nil {
		userRepository = &login_implement.UserRepositoryImpl{}
	}
	return userRepository
}

type UserRepositoryInterface interface {
	Save(ctx context.Context, user *login_entity.UserProfile) error
	GetUser(ctx context.Context, userID string) (*login_entity.UserProfile, error)
}
