package cs_service

import (
	"context"
	"errors"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_facade"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type RoomService struct {
	room         *cs_entity.CsRoom
	roomRepo     cs_facade.RoomRepositoryInterface
	roomUserRepo cs_facade.RoomUserRepositoryInterface
	interactRepo cs_facade.InteractRepositoryInterface
}

func NewRoomService(ctx context.Context, room *cs_entity.CsRoom) *RoomService {
	roomService := &RoomService{
		room:         room,
		roomRepo:     cs_facade.GetRoomRepository(),
		roomUserRepo: cs_facade.GetRoomUserRepository(),
		interactRepo: cs_facade.GetInteractRepository(),
	}
	return roomService
}

func NewRoomServiceByRoomID(ctx context.Context, roomID string) (*RoomService, error) {
	roomService := &RoomService{
		roomRepo:     cs_facade.GetRoomRepository(),
		roomUserRepo: cs_facade.GetRoomUserRepository(),
		interactRepo: cs_facade.GetInteractRepository(),
	}
	room, err := roomService.roomRepo.GetRoomByRoomID(ctx, roomID)
	if err != nil || room == nil {
		if err == nil {
			err = errors.New("room is not exist")
		}
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	roomUserCount, _ := roomService.roomUserRepo.GetUserCountByRoomID(ctx, roomID)
	room.UserCount = roomUserCount
	roomService.room = room
	return roomService, nil
}

func (rs *RoomService) Start(ctx context.Context) error {
	rs.room.Start()
	err := rs.roomRepo.Save(ctx, rs.room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (rs *RoomService) Finish(ctx context.Context) error {
	rs.roomUserRepo.UpdateUsersByRoomID(ctx, rs.room.GetRoomID(), map[string]interface{}{
		"net_status": cs_entity.UserNetStatusOffline,
	})

	rs.room.Finish()
	err := rs.roomRepo.Save(ctx, rs.room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (rs *RoomService) Join(ctx context.Context, userID, userName string) error {
	user := &cs_entity.CsRoomUser{
		AppID:      rs.room.GetAppID(),
		RoomID:     rs.room.GetRoomID(),
		UserID:     userID,
		UserName:   userName,
		UserRole:   cs_entity.UserRoleAudience,
		Mic:        0,
		Camera:     0,
		NetStatus:  cs_entity.UserNetStatusOnline,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err := rs.roomUserRepo.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (rs *RoomService) Leave(ctx context.Context, userID string) error {
	user, err := rs.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, rs.room.GetRoomID(), userID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	if user == nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return custom_error.InternalError(errors.New("user is not exist"))
	}
	user.Leave()
	err = rs.roomUserRepo.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	return nil
}

func (rs *RoomService) ListAudiences(ctx context.Context) ([]*cs_entity.CsRoomUser, error) {
	users, err := rs.roomUserRepo.GetActiveUsersByRoomIDRole(ctx, rs.room.RoomID, cs_entity.UserRoleAudience)
	if err != nil {
		logs.CtxError(ctx, "get users failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return users, nil
}

func (rs *RoomService) ListUsers(ctx context.Context) ([]*cs_entity.CsRoomUser, error) {
	users, err := rs.roomUserRepo.GetActiveUsersByRoomID(ctx, rs.room.RoomID)
	if err != nil {
		logs.CtxError(ctx, "get users failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return users, nil
}

func (rs *RoomService) ListAudiencesByInteractStatus(ctx context.Context, status []int) ([]*cs_entity.CsRoomUser, error) {
	users, err := rs.roomUserRepo.GetActiveAudiencesByRoomIDStatus(ctx, rs.room.GetRoomID(), status)
	if err != nil {
		logs.CtxError(ctx, "get users failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return users, nil
}

func (rs *RoomService) ListUsersByInteractStatus(ctx context.Context, status []int) ([]*cs_entity.CsRoomUser, error) {
	users, err := rs.roomUserRepo.GetActiveUsersByRoomIDStatus(ctx, rs.room.GetRoomID(), status)
	if err != nil {
		logs.CtxError(ctx, "get users failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return users, nil
}

func (rs *RoomService) CreateInteract(ctx context.Context, interactID string) (*cs_entity.CsInteract, error) {
	interact := &cs_entity.CsInteract{
		InteractID:  interactID,
		OwnerUserID: rs.room.GetOwnerUserID(),
		OwnerRoomID: rs.room.GetRoomID(),
		RtcAppID:    rs.room.GetAppID(),
		RtcRoomID:   rs.room.GetRoomID(),
	}
	err := rs.interactRepo.Save(ctx, interact)
	if err != nil {
		logs.CtxError(ctx, "save interact failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return interact, nil
}

func (rs *RoomService) ChangeHost(ctx context.Context) (*cs_entity.CsRoomUser, error) {
	host, err := rs.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, rs.room.GetRoomID(), rs.room.GetOwnerUserID())
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	if host == nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, custom_error.InternalError(errors.New("user is not exist"))
	}
	host.Leave()
	err = rs.roomUserRepo.Save(ctx, host)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}

	nextHost, err := rs.roomUserRepo.GetUserByStatusOrderUtime(ctx, rs.room.GetRoomID(), cs_models.UserInteractStatusOnMicrophone)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	if nextHost == nil {
		err = rs.Finish(ctx)
		if err != nil {
			logs.CtxError(ctx, "finish room failed,error:%s", err)
			return nil, custom_error.InternalError(err)
		}
		return nil, nil
	}

	nextHost.UserRole = cs_entity.UserRoleHost
	nextHost.InteractStatus = cs_models.UserInteractStatusOnMicrophone
	err = rs.roomUserRepo.Save(ctx, nextHost)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	rs.room.SetOwnerUserID(nextHost.GetUserID())
	rs.room.SetOwnerUserName(nextHost.GetUserName())
	err = rs.roomRepo.Save(ctx, rs.room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return nextHost, nil
}

func (rs *RoomService) GetRoom() *cs_entity.CsRoom {
	return rs.room
}
