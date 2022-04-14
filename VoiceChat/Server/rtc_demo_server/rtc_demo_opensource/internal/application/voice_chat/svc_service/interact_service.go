package svc_service

import (
	"context"
	"errors"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

const (
	InteractReplyTypeAccept  = 1
	InteractReplyTypeReject  = 2
	InteractReplyTypeTimeout = 3
)

const (
	InteractManageTypeLockSeat   = 1
	InteractManageTypeUnlockSeat = 2
	InteractManageTypeMute       = 3
	InteractManageTypeUnmute     = 4
	InteractManageTypeKick       = 5
)

const (
	InteractFinishTypeHost = 1
	InteractFinishTypeSelf = 2
)

var interactServiceClient *InteractService

type InteractService struct {
	userFactory *UserFactory
	seatFactory *SeatFactory
}

func GetInteractService() *InteractService {
	if interactServiceClient == nil {
		interactServiceClient = &InteractService{
			userFactory: GetUserFactory(),
			seatFactory: GetSeatFactory(),
		}
	}
	return interactServiceClient
}

func (is *InteractService) Invite(ctx context.Context, roomID, hostUserID, audienceUserID string, seatID int) error {
	seats, err := is.seatFactory.GetSeatsByRoomID(ctx, roomID)
	if err != nil {
		logs.CtxError(ctx, "get seats failed,error:%s", err)
		return err
	}
	seatAvailable := false
	for _, s := range seats {
		if s.IsEnableAlloc() {
			seatAvailable = true
			break
		}
	}
	if !seatAvailable {
		return custom_error.ErrUserOnMicExceedLimit
	}

	audience, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, audienceUserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if audience == nil {
		logs.CtxError(ctx, "user is not exist", err)
		return custom_error.ErrUserNotExist
	}
	if !audience.IsEnableInvite() {
		logs.CtxError(ctx, "user is interacting,roomID:%s,userID:%s", roomID, audienceUserID)
		return custom_error.InternalError(errors.New("user is interacting"))
	}

	audience.SetInteract(svc_db.UserInteractStatusInviting, seatID)
	err = is.userFactory.Save(ctx, audience)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return err
	}

	//inform
	informer := GetInformService()

	host, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, hostUserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if host == nil {
		logs.CtxError(ctx, "user is not exist", err)
		return custom_error.ErrUserNotExist
	}
	data := &InformInviteInteract{
		HostInfo: host,
		SeatID:   seatID,
	}
	informer.UnicastUser(ctx, audience.GetRoomID(), audience.GetUserID(), OnInviteInteract, data)
	return nil
}

func (is *InteractService) Apply(ctx context.Context, roomID, hostUserID, audienceUserID string, seatID int) error {
	audience, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, audienceUserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if audience == nil {
		logs.CtxError(ctx, "user is not exist", err)
		return custom_error.ErrUserNotExist
	}
	if !audience.IsEnableApply() {
		logs.CtxError(ctx, "user is interacting,roomID:%s,userID:%s", roomID, audienceUserID)
		return custom_error.InternalError(errors.New("user is interacting"))
	}

	audience.SetInteract(svc_db.UserInteractStatusApplying, seatID)
	err = is.userFactory.Save(ctx, audience)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return err
	}

	return nil
}

func (is *InteractService) AudienceReply(ctx context.Context, roomID, hostUserID, audienceUserID string, replyType int) error {
	audience, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, audienceUserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if audience == nil {
		logs.CtxError(ctx, "user is not exist", err)
		return custom_error.ErrUserNotExist
	}

	if replyType == InteractReplyTypeAccept {
		if !audience.IsEnableInteract() {
			logs.CtxError(ctx, "user can not interact")
			return custom_error.InternalError(errors.New("user can not interact"))
		}

		seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, audience.SeatID)
		if err != nil {
			logs.CtxError(ctx, "get seat failed,error:%s", err)
			return err
		}
		if seat == nil {
			logs.CtxError(ctx, "seat is not exist")
			return custom_error.InternalError(errors.New("seat is not exist"))
		}

		if !seat.IsEnableAlloc() {
			seat = nil
			seats, err := is.seatFactory.GetSeatsByRoomID(ctx, roomID)
			if err != nil {
				logs.CtxError(ctx, "get seats failed,error:%s", err)
				return err
			}
			for _, s := range seats {
				if s.IsEnableAlloc() {
					seat = s
					break
				}
			}
			if seat == nil {
				audience.SetInteract(svc_db.UserInteractStatusNormal, 0)
				err = is.userFactory.Save(ctx, audience)
				if err != nil {
					logs.CtxError(ctx, "save user failed,error:%s", err)
					return err
				}
				logs.CtxWarn(ctx, "no seat can be alloc")
				return custom_error.ErrUserOnMicExceedLimit
			}
		}

		audience.SetInteract(svc_db.UserInteractStatusInteracting, seat.GetSeatID())
		err = is.userFactory.Save(ctx, audience)
		if err != nil {
			logs.CtxError(ctx, "save user failed,error:%s", err)
			return err
		}
		seat.SetOwnerUserID(audience.GetUserID())
		err = is.seatFactory.Save(ctx, seat)
		if err != nil {
			logs.CtxError(ctx, "save seat failed,error:%s", err)
			return err
		}

		informer := GetInformService()
		data := &InformJoinInteract{
			UserInfo: audience,
			SeatID:   audience.GetSeatID(),
		}
		informer.BroadcastRoom(ctx, roomID, OnJoinInteract, data)

	} else {
		audience.SetInteract(svc_db.UserInteractStatusNormal, 0)
		err = is.userFactory.Save(ctx, audience)
		if err != nil {
			logs.CtxError(ctx, "save user failed,error:%s", err)
			return err
		}
	}

	informer := GetInformService()
	host, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, hostUserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if host == nil {
		logs.CtxError(ctx, "user is not exist", err)
		return custom_error.ErrUserNotExist
	}
	data := &InformInviteResult{
		UserInfo: audience,
		Reply:    replyType,
	}
	informer.UnicastUser(ctx, host.GetRoomID(), host.GetUserID(), OnInviteResult, data)

	return nil
}

func (is *InteractService) HostReply(ctx context.Context, roomID, hostUserID, audienceUserID string, replyType int) error {
	audience, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, audienceUserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return err
	}
	if audience == nil {
		logs.CtxError(ctx, "user is not exist", err)
		return custom_error.ErrUserNotExist
	}

	if replyType == InteractReplyTypeAccept {
		if !audience.IsEnableInteract() {
			logs.CtxError(ctx, "user can not interact")
			return custom_error.InternalError(errors.New("user can not interact"))
		}

		seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, audience.SeatID)
		if err != nil {
			logs.CtxError(ctx, "get seat failed,error:%s", err)
			return err
		}
		if seat == nil {
			logs.CtxError(ctx, "seat is not exist")
			return custom_error.InternalError(errors.New("seat is not exist"))
		}

		if !seat.IsEnableAlloc() {
			seat = nil
			seats, err := is.seatFactory.GetSeatsByRoomID(ctx, roomID)
			if err != nil {
				logs.CtxError(ctx, "get seats failed,error:%s", err)
				return err
			}
			for _, s := range seats {
				if s.IsEnableAlloc() {
					seat = s
					break
				}
			}
			if seat == nil {
				audience.SetInteract(svc_db.UserInteractStatusNormal, 0)
				err = is.userFactory.Save(ctx, audience)
				if err != nil {
					logs.CtxError(ctx, "save user failed,error:%s", err)
					return err
				}
				logs.CtxWarn(ctx, "no seat can be alloc")
				return custom_error.ErrUserOnMicExceedLimit
			}
		}

		audience.SetInteract(svc_db.UserInteractStatusInteracting, seat.GetSeatID())
		err = is.userFactory.Save(ctx, audience)
		if err != nil {
			logs.CtxError(ctx, "save user failed,error:%s", err)
			return err
		}
		seat.SetOwnerUserID(audience.GetUserID())
		err = is.seatFactory.Save(ctx, seat)
		if err != nil {
			logs.CtxError(ctx, "save seat failed,error:%s", err)
			return err
		}

		informer := GetInformService()
		data := &InformJoinInteract{
			UserInfo: audience,
			SeatID:   audience.GetSeatID(),
		}
		informer.BroadcastRoom(ctx, roomID, OnJoinInteract, data)
	} else {
		audience.SetInteract(svc_db.UserInteractStatusNormal, 0)
		err = is.userFactory.Save(ctx, audience)
		if err != nil {
			logs.CtxError(ctx, "save user failed,error:%s", err)
			return err
		}
	}

	return nil
}

func (is *InteractService) Mute(ctx context.Context, roomID string, seatID int) error {
	seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, seatID)
	if err != nil {
		logs.CtxError(ctx, "get seat failed,error:%s", err)
		return err
	}
	if seat == nil {
		logs.CtxError(ctx, "seat is not exist")
		return custom_error.InternalError(errors.New("seat is not exist"))
	}

	userID := seat.GetOwnerUserID()
	if userID == "" {
		logs.CtxError(ctx, "no user in this seat")
		return custom_error.InternalError(errors.New("no user in this seat"))
	}

	informer := GetInformService()
	data := &InformMediaOperate{
		Mic: 0,
	}
	informer.UnicastUser(ctx, roomID, userID, OnMediaOperate, data)
	return nil

}

func (is *InteractService) Unmute(ctx context.Context, roomID string, seatID int) error {
	seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, seatID)
	if err != nil {
		logs.CtxError(ctx, "get seat failed,error:%s", err)
		return err
	}
	if seat == nil {
		logs.CtxError(ctx, "seat is not exist")
		return custom_error.InternalError(errors.New("seat is not exist"))
	}

	userID := seat.GetOwnerUserID()
	if userID == "" {
		logs.CtxError(ctx, "no user in this seat")
		return custom_error.InternalError(errors.New("no user in this seat"))
	}

	informer := GetInformService()
	data := &InformMediaOperate{
		Mic: 1,
	}
	informer.UnicastUser(ctx, roomID, userID, OnMediaOperate, data)
	return nil
}

func (is *InteractService) LockSeat(ctx context.Context, roomID string, seatID int) error {
	seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, seatID)
	if err != nil {
		logs.CtxError(ctx, "get seat failed,error:%s", err)
		return err
	}
	if seat == nil {
		logs.CtxError(ctx, "seat is not exist")
		return custom_error.InternalError(errors.New("seat is not exist"))
	}

	userID := seat.GetOwnerUserID()
	if userID != "" {
		err = is.FinishInteract(ctx, roomID, seatID, InteractFinishTypeHost)
	}

	seat.Lock()
	seat.SetOwnerUserID("")
	err = is.seatFactory.Save(ctx, seat)
	if err != nil {
		logs.CtxError(ctx, "save seat failed,error:%s", err)
		return err
	}

	informer := GetInformService()
	data := &InformSeatStatusChange{
		SeatID: seatID,
		Type:   0,
	}
	informer.BroadcastRoom(ctx, roomID, OnSeatStatusChange, data)
	return nil

}

func (is *InteractService) UnlockSeat(ctx context.Context, roomID string, seatID int) error {
	seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, seatID)
	if err != nil {
		logs.CtxError(ctx, "get seat failed,error:%s", err)
		return err
	}
	if seat == nil {
		logs.CtxError(ctx, "seat is not exist")
		return custom_error.InternalError(errors.New("seat is not exist"))
	}

	seat.Unlock()
	err = is.seatFactory.Save(ctx, seat)
	if err != nil {
		logs.CtxError(ctx, "save seat failed,error:%s", err)
		return err
	}

	informer := GetInformService()
	data := &InformSeatStatusChange{
		SeatID: seatID,
		Type:   1,
	}
	informer.BroadcastRoom(ctx, roomID, OnSeatStatusChange, data)
	return nil
}

func (is *InteractService) FinishInteract(ctx context.Context, roomID string, seatID int, finishType int) error {
	seat, err := is.seatFactory.GetSeatByRoomIDSeatID(ctx, roomID, seatID)
	if err != nil {
		logs.CtxError(ctx, "get seat failed,error:%s", err)
		return err
	}
	if seat == nil {
		logs.CtxError(ctx, "seat is not exist")
		return custom_error.InternalError(errors.New("seat is not exist"))
	}

	userID := seat.GetOwnerUserID()
	if userID != "" {
		user, err := is.userFactory.GetActiveUserByRoomIDUserID(ctx, roomID, userID)
		if err != nil {
			logs.CtxError(ctx, "get user failed,error:%s", err)
			return err
		}
		if user == nil {
			logs.CtxError(ctx, "user is not exist", err)
			return custom_error.ErrUserNotExist
		}
		user.SetInteract(svc_db.UserInteractStatusNormal, 0)
		user.UnmuteMic()
		err = is.userFactory.Save(ctx, user)
		if err != nil {
			logs.CtxError(ctx, "save user failed,error:%s", err)
			return err
		}

		informer := GetInformService()
		data := &InformFinishInteract{
			UserInfo: user,
			SeatID:   seatID,
			Type:     finishType,
		}
		informer.BroadcastRoom(ctx, roomID, OnFinishInteract, data)
	}

	seat.SetOwnerUserID("")
	err = is.seatFactory.Save(ctx, seat)
	if err != nil {
		logs.CtxError(ctx, "save seat failed,error:%s", err)
		return err
	}

	return nil
}
