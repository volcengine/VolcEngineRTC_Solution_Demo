package svc_service

import (
	"context"
	"errors"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

const (
	FinishTypeNormal        = 1
	FinishTypeTimeout       = 2
	FinishTypeReviewIllegal = 3
)

var roomServiceClient *RoomService

type SeatInfo struct {
	Status    int   `json:"status"`
	GuestInfo *User `json:"guest_info"`
}

type RoomService struct {
	roomFactory *RoomFactory
	userFactory *UserFactory
	seatFactory *SeatFactory
}

func GetRoomService() *RoomService {
	if roomServiceClient == nil {
		roomServiceClient = &RoomService{
			roomFactory: GetRoomFactory(),
			userFactory: GetUserFactory(),
			seatFactory: GetSeatFactory(),
		}
	}
	return roomServiceClient
}

func (rs *RoomService) CreateRoom(ctx context.Context, roomName, roomBackgroundImageName, hostUserID, hostUserName, hostDeviceID string) (*Room, *User, error) {
	room, err := rs.roomFactory.NewRoom(ctx, roomName, roomBackgroundImageName, hostUserID, hostUserName)
	if err != nil || room == nil {
		logs.CtxError(ctx, "create room failed,error:%s", err)
		return nil, nil, custom_error.InternalError(errors.New("create room failed"))
	}
	logs.CtxInfo(ctx, "room:%#v", room.GetDbRoom())
	err = rs.roomFactory.Save(ctx, room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return nil, nil, err
	}

	host := rs.userFactory.NewUser(ctx, room.GetRoomID(), hostUserID, hostUserName, hostDeviceID, svc_db.UserRoleHost)
	err = rs.userFactory.Save(ctx, host)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return nil, nil, err
	}

	for seatID := 1; seatID <= 8; seatID++ {
		seat := rs.seatFactory.NewSeat(ctx, room.GetRoomID(), seatID)
		err = rs.seatFactory.Save(ctx, seat)
		if err != nil {
			logs.CtxError(ctx, "save seat failed,error:%s", err)
			return nil, nil, err
		}
	}

	return room, host, nil
}

func (rs *RoomService) StartLive(ctx context.Context, roomID string) error {
	room, err := rs.roomFactory.GetRoomByRoomID(ctx, roomID)
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist")
		return custom_error.ErrRoomNotExist
	}

	room.Start()
	err = rs.roomFactory.Save(ctx, room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return custom_error.InternalError(err)
	}

	host, err := rs.userFactory.GetUserByRoomIDUserID(ctx, room.GetRoomID(), room.GetHostUserID())
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if host == nil {
		logs.CtxError(ctx, "user is not exist")
		return custom_error.ErrUserNotExist
	}

	host.StartLive()
	err = rs.userFactory.Save(ctx, host)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return err
	}

	//todo start review

	return nil
}

func (rs *RoomService) FinishLive(ctx context.Context, roomID string, finishType int) error {
	room, err := rs.roomFactory.GetRoomByRoomID(ctx, roomID)
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist")
		return custom_error.ErrRoomNotExist
	}

	//inform
	informer := GetInformService()
	data := &InformFinishLive{
		RoomID: room.GetRoomID(),
		Type:   finishType,
	}
	informer.BroadcastRoom(ctx, roomID, OnFinishLive, data)

	//update
	room.Finish()
	err = rs.roomFactory.Save(ctx, room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return custom_error.InternalError(err)
	}

	err = rs.userFactory.UpdateUsersByRoomID(ctx, roomID, map[string]interface{}{
		"net_status":      svc_db.UserNetStatusOffline,
		"interact_status": svc_db.UserInteractStatusNormal,
		"seat_id":         0,
		"leave_time":      time.Now().UnixNano() / 1e6,
		"update_time":     time.Now(),
	})

	if err != nil {
		logs.CtxError(ctx, "update user failed,error:%s", err)
	}

	//todo stop review

	return nil
}

func (rs *RoomService) JoinRoom(ctx context.Context, roomID, userID, userName, userDeviceID string) error {
	room, err := rs.roomFactory.GetRoomByRoomID(ctx, roomID)
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return err
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist")
		return custom_error.ErrRoomNotExist
	}

	user := rs.userFactory.NewUser(ctx, roomID, userID, userName, userDeviceID, svc_db.UserRoleAudience)
	user.JoinRoom(room.GetRoomID())
	err = rs.userFactory.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s")
		return custom_error.InternalError(err)
	}

	err = rs.roomFactory.IncrRoomAudienceCount(ctx, room.GetRoomID(), 1)
	if err != nil {
		logs.CtxError(ctx, "incr room audience count failed,error:%s", err)
	}

	//inform
	userCountMap, err := rs.roomFactory.GetRoomsAudienceCount(ctx, []string{room.GetRoomID()})
	if err != nil {
		logs.CtxError(ctx, "get user count failed,error:%s")
	}
	userCount := userCountMap[room.GetRoomID()]
	data := &InformJoinRoom{
		UserInfo:      user,
		AudienceCount: userCount,
	}
	informer := GetInformService()
	informer.BroadcastRoom(ctx, roomID, OnAudienceJoinRoom, data)
	return nil
}

func (rs *RoomService) LeaveRoom(ctx context.Context, roomID, userID string) error {
	room, err := rs.roomFactory.GetRoomByRoomID(ctx, roomID)
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return err
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist")
		return custom_error.ErrRoomNotExist
	}

	user, err := rs.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, userID)
	if err != nil || user == nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil
	}

	if user.IsInteract() {
		interactService := GetInteractService()
		err = interactService.FinishInteract(ctx, user.GetRoomID(), user.GetSeatID(), InteractFinishTypeSelf)
		if err != nil {
			logs.CtxError(ctx, "finish interact failed,error:%s", err)
		}
	}

	user.LeaveRoom()
	err = rs.userFactory.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s")
		return custom_error.InternalError(err)
	}

	err = rs.roomFactory.IncrRoomAudienceCount(ctx, room.GetRoomID(), -1)
	if err != nil {
		logs.CtxError(ctx, "incr room audience count failed,error:%s", err)
	}

	//inform
	userCountMap, err := rs.roomFactory.GetRoomsAudienceCount(ctx, []string{room.GetRoomID()})
	if err != nil {
		logs.CtxError(ctx, "get user count failed,error:%s")
	}
	userCount := userCountMap[room.GetRoomID()]
	data := &InformLeaveRoom{
		UserInfo:      user,
		AudienceCount: userCount,
	}
	informer := GetInformService()
	informer.BroadcastRoom(ctx, roomID, OnAudienceLeaveRoom, data)
	return nil
}

func (rs *RoomService) Disconnect(ctx context.Context, roomID, userID string) {
	user, err := rs.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, userID)
	if err != nil || user == nil {
		logs.CtxWarn(ctx, "get user failed,error:%s", err)
		return
	}
	logs.CtxInfo(ctx, "svc disconnect user:%#v", user.SvcUser)

	user.Disconnect()
	err = rs.userFactory.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s")
		return
	}

	go func(ctx context.Context, roomID, userID string) {
		time.Sleep(time.Duration(config.Configs().ReconnectTimeout) * time.Second)
		user, err := rs.userFactory.GetUserByRoomIDUserID(ctx, roomID, userID)
		if err != nil || user == nil {
			logs.CtxWarn(ctx, "get user failed,error:%s", err)
			return
		}
		roomService := GetRoomService()
		if user.IsReconnecting() {
			if user.IsHost() {
				err = roomService.FinishLive(ctx, user.GetRoomID(), FinishTypeNormal)
				if err != nil {
					logs.CtxError(ctx, "finish live failed,error:%s", err)
				}
			} else if user.IsAudience() {
				err = roomService.LeaveRoom(ctx, user.GetRoomID(), user.GetUserID())
				if err != nil {
					logs.CtxError(ctx, "leave room failed,error:%s", err)
				}
			}

		}
	}(ctx, user.GetRoomID(), user.GetUserID())
}
