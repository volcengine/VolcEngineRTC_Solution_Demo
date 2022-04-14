package svc_service

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
)

var userFactoryClient *UserFactory

type UserFactory struct {
	userRepo UserRepo
}

func GetUserFactory() *UserFactory {
	if userFactoryClient == nil {
		userFactoryClient = &UserFactory{
			userRepo: GetUserRepo(),
		}
	}
	return userFactoryClient
}

func (uf *UserFactory) NewUser(ctx context.Context, roomID, userID, userName, deviceID string, userRole int) *User {
	dbUser := &svc_entity.SvcUser{
		AppID:          config.Configs().SvcAppID,
		RoomID:         roomID,
		UserID:         userID,
		UserName:       userName,
		UserRole:       userRole,
		NetStatus:      svc_db.UserNetStatusOffline,
		InteractStatus: svc_db.UserInteractStatusNormal,
		SeatID:         0,
		Mic:            1,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		DeviceID:       deviceID,
	}
	user := &User{
		SvcUser: dbUser,
		isDirty: true,
	}
	return user
}

func (uf *UserFactory) Save(ctx context.Context, user *User) error {
	if user.IsDirty() {
		user.SetUpdateTime(time.Now())
		err := uf.userRepo.Save(ctx, user.SvcUser)
		if err != nil {
			return custom_error.InternalError(err)
		}
		user.SetIsDirty(false)
	}
	return nil
}

func (uf *UserFactory) GetActiveUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*User, error) {
	dbUser, err := uf.userRepo.GetActiveUserByRoomIDUserID(ctx, roomID, userID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	if dbUser == nil {
		return nil, nil
	}
	user := &User{
		SvcUser: dbUser,
		isDirty: false,
	}
	return user, nil
}

func (uf *UserFactory) GetUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*User, error) {
	dbUser, err := uf.userRepo.GetUserByRoomIDUserID(ctx, roomID, userID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	if dbUser == nil {
		return nil, nil
	}
	user := &User{
		SvcUser: dbUser,
		isDirty: false,
	}
	return user, nil
}

func (uf *UserFactory) GetActiveUserByUserID(ctx context.Context, userID string) (*User, error) {
	dbUser, err := uf.userRepo.GetActiveUserByUserID(ctx, userID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	if dbUser == nil {
		return nil, nil
	}
	user := &User{
		SvcUser: dbUser,
		isDirty: false,
	}
	return user, nil
}

func (uf *UserFactory) GetAudiencesWithoutApplyByRoomID(ctx context.Context, roomID string) ([]*User, error) {
	dbAudiences, err := uf.userRepo.GetAudiencesWithoutApplyByRoomID(ctx, roomID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}

	audiences := make([]*User, len(dbAudiences))
	for i := 0; i < len(dbAudiences); i++ {
		audiences[i] = &User{
			SvcUser: dbAudiences[i],
			isDirty: false,
		}
	}
	return audiences, nil
}

func (uf *UserFactory) GetApplyAudiencesByRoomID(ctx context.Context, roomID string) ([]*User, error) {
	dbAudiences, err := uf.userRepo.GetApplyAudiencesByRoomID(ctx, roomID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}

	audiences := make([]*User, len(dbAudiences))
	for i := 0; i < len(dbAudiences); i++ {
		audiences[i] = &User{
			SvcUser: dbAudiences[i],
			isDirty: false,
		}
	}
	return audiences, nil
}

func (uf *UserFactory) GetUsersByRoomID(ctx context.Context, roomID string) ([]*User, error) {
	dbAudiences, err := uf.userRepo.GetActiveUsersByRoomID(ctx, roomID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}

	audiences := make([]*User, len(dbAudiences))
	for i := 0; i < len(dbAudiences); i++ {
		audiences[i] = &User{
			SvcUser: dbAudiences[i],
			isDirty: false,
		}
	}
	return audiences, nil
}

func (uf *UserFactory) UpdateUsersByRoomID(ctx context.Context, roomID string, ups map[string]interface{}) error {
	err := uf.userRepo.UpdateUsersWithMapByRoomID(ctx, roomID, ups)
	if err != nil {
		return custom_error.InternalError(err)
	}
	return nil
}
