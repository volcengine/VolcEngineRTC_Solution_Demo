package cs_service

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
)

func Room2CsRoomInfo(ctx context.Context, room *cs_entity.CsRoom) *cs_models.RoomInfo {
	roomInfo := &cs_models.RoomInfo{
		RoomID:      room.GetRoomID(),
		RoomName:    room.GetRoomName(),
		HostID:      room.GetOwnerUserID(),
		HostName:    room.GetOwnerUserName(),
		UserCounts:  room.UserCount,
		MiconCounts: 1,
		CreatedAt:   room.CreateTime.UnixNano(),
		Now:         time.Now().UnixNano(),
	}
	roomService := NewRoomService(ctx, room)
	users, _ := roomService.ListUsersByInteractStatus(ctx, []int{cs_models.UserInteractStatusOnMicrophone})
	roomInfo.MiconCounts = len(users)
	users, _ = roomService.ListUsers(ctx)
	roomInfo.UserCounts = len(users)
	return roomInfo
}
func User2CsUserInfo(user *cs_entity.CsRoomUser) *cs_models.UserInfo {
	userStatus := user.GetInteractStatus()
	if userStatus == cs_models.UserInteractStatusInviting {
		userStatus = cs_models.UserInteractStatusAudience
	}
	return &cs_models.UserInfo{
		UserID:       user.UserID,
		UserName:     user.UserName,
		IsHost:       user.IsHost(),
		UserStatus:   user.GetInteractStatus(),
		IsMicOn:      user.Mic == 1,
		CreatedAt:    user.CreateTime.UnixNano(),
		RaiseHandsAt: user.UpdateTime.UnixNano(),
	}
}
