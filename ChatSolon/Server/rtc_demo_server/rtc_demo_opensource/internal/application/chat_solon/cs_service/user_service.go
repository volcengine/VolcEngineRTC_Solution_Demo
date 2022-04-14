package cs_service

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_facade"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type UserService struct {
	user         *cs_entity.CsRoomUser
	roomUserRepo cs_facade.RoomUserRepositoryInterface
}

func NewUserServiceByRoomIDUserID(ctx context.Context, roomID, userID string) (*UserService, error) {
	userService := &UserService{
		roomUserRepo: cs_facade.GetRoomUserRepository(),
	}
	user, err := userService.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, roomID, userID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	if user == nil {
		logs.CtxError(ctx, "user is not exist")
		return nil, custom_error.ErrUserNotExist
	}

	userService.user = user
	return userService, nil
}

func NewReconnectUserServiceByRoomIDUserID(ctx context.Context, roomID, userID string) (*UserService, error) {
	userService := &UserService{
		roomUserRepo: cs_facade.GetRoomUserRepository(),
	}
	user, err := userService.roomUserRepo.GetReconnectUserByRoomIDUserID(ctx, roomID, userID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	if user == nil {
		logs.CtxError(ctx, "user is not exist")
		return nil, custom_error.ErrUserNotExist
	}

	userService.user = user
	return userService, nil
}

func (us *UserService) Mute(ctx context.Context) error {
	us.user.Mute()
	err := us.roomUserRepo.Save(ctx, us.user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (us *UserService) Disconnect(ctx context.Context) error {
	us.user.NetStatus = cs_entity.UserNetStatusReconnecting
	err := us.roomUserRepo.Save(ctx, us.user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (us *UserService) Reconnect(ctx context.Context) error {
	us.user.NetStatus = cs_entity.UserNetStatusOnline
	err := us.roomUserRepo.Save(ctx, us.user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (us *UserService) Unmute(ctx context.Context) error {
	us.user.Unmute()
	err := us.roomUserRepo.Save(ctx, us.user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (us *UserService) GetUser() *cs_entity.CsRoomUser {
	return us.user
}
