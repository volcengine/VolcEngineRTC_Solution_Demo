package cs_service

import (
	"context"
	"errors"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_facade"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type InteractService struct {
	interact           *cs_entity.CsInteract
	interactRepo       cs_facade.InteractRepositoryInterface
	interactDetailRepo cs_facade.InteractDetailRepositoryInterface
	roomUserRepo       cs_facade.RoomUserRepositoryInterface
}

func NewInteractServiceByRoomID(ctx context.Context, roomID string) (*InteractService, error) {
	interactService := &InteractService{
		interactRepo:       cs_facade.GetInteractRepository(),
		interactDetailRepo: cs_facade.GetInteractDetailRepository(),
		roomUserRepo:       cs_facade.GetRoomUserRepository(),
	}
	interact, err := interactService.interactRepo.GetInteractByRoomID(ctx, roomID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	if interact == nil {
		return nil, custom_error.InternalError(errors.New("room interact not exist"))
	}
	interactService.interact = interact
	return interactService, nil
}

func (is *InteractService) Invite(ctx context.Context, audienceID string) error {
	audience, err := is.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if err != nil || audience == nil {
		if err == nil {
			err = errors.New("user is not exist")
		}
		return custom_error.InternalError(err)
	}
	/*
		if audience.GetInteractStatus() != cs_models.UserInteractStatusAudience {
			logs.CtxError(ctx, "invite user status is not audience")
			return custom_error.InternalError(errors.New("user status error"))
		}

		audience.SetInteractStatus(cs_models.UserInteractStatusInviting)
		err = is.roomUserRepo.Save(ctx, audience)
		if err != nil {
			return custom_error.InternalError(err)
		}

		interactDetail := &entity.InteractDetail{
			InteractID: is.interact.GetInteractID(),
			FromRoomID: is.interact.GetOwnerRoomID(),
			FromUserID: is.interact.GetOwnerUserID(),
			ToRoomID:   is.interact.GetRtcRoomID(),
			ToUserID:   audienceID,
			Status:     cs_models.UserInteractStatusInviting,
		}
		err = is.interactDetailRepo.Save(ctx, interactDetail)
		if err != nil {
			return custom_error.InternalError(err)
		}


	*/
	informer := GetInformService()
	informer.UnicastUser(ctx, is.interact.OwnerRoomID, audienceID, cs_models.OnCsInviteMic, User2CsUserInfo(audience))
	return nil
}

func (is *InteractService) RaiseHands(ctx context.Context, audienceID string) error {
	audience, err := is.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if err != nil || audience == nil {
		if err == nil {
			err = errors.New("user is not exist")
		}
		return custom_error.InternalError(err)
	}
	/*
		if audience.GetInteractStatus() != cs_models.UserInteractStatusAudience {
			logs.CtxError(ctx, "invite user status is not audience")
			return custom_error.InternalError(errors.New("user status error"))
		}

	*/

	audience.SetInteractStatus(cs_models.UserInteractStatusRaiseHands)
	err = is.roomUserRepo.Save(ctx, audience)
	if err != nil {
		return custom_error.InternalError(err)
	}

	interactDetail := &cs_entity.CsInteractDetail{
		InteractID: is.interact.GetInteractID(),
		FromRoomID: is.interact.GetOwnerRoomID(),
		FromUserID: is.interact.GetOwnerUserID(),
		ToRoomID:   is.interact.GetRtcRoomID(),
		ToUserID:   audienceID,
		Status:     cs_models.UserInteractStatusRaiseHands,
	}
	err = is.interactDetailRepo.Save(ctx, interactDetail)
	if err != nil {
		return custom_error.InternalError(err)
	}

	informer := GetInformService()
	informer.BroadcastRoom(ctx, is.interact.OwnerRoomID, cs_models.OnCsRaiseHandsMic, User2CsUserInfo(audience))
	return nil
}

func (is *InteractService) Confirm(ctx context.Context, audienceID string) error {
	audience, err := is.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if err != nil || audience == nil {
		if err == nil {
			err = errors.New("user is not exist")
		}
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	audience.SetInteractStatus(cs_models.UserInteractStatusOnMicrophone)
	err = is.roomUserRepo.Save(ctx, audience)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}

	/*
		interactDetail, err := is.interactDetailRepo.GetInteractDetailByFromRoomIDToUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
		if err != nil || interactDetail == nil {
			if err == nil {
				err = errors.New("interact is not exist")
			}
			logs.CtxError(ctx, "get interact detail failed,error:%s", err)
			return custom_error.InternalError(err)
		}
		interactDetail.SetStatus(cs_models.UserInteractStatusOnMicrophone)
		err = is.interactDetailRepo.Save(ctx, interactDetail)
		if err != nil {
			logs.CtxError(ctx, "save interact detail failed,error:%s", err)
			return custom_error.InternalError(err)
		}

	*/

	informer := GetInformService()
	informer.BroadcastRoom(ctx, is.interact.OwnerRoomID, cs_models.OnCsMicOn, User2CsUserInfo(audience))
	return nil
}

func (is *InteractService) Agree(ctx context.Context, audienceID string) error {
	audience, err := is.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if err != nil || audience == nil {
		if err == nil {
			err = errors.New("user is not exist")
		}
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	audience.SetInteractStatus(cs_models.UserInteractStatusOnMicrophone)
	err = is.roomUserRepo.Save(ctx, audience)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}

	interactDetail, err := is.interactDetailRepo.GetInteractDetailByFromRoomIDToUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if err != nil || interactDetail == nil {
		if err == nil {
			err = errors.New("interact is not exist")
		}
		logs.CtxError(ctx, "get interact detail failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	interactDetail.SetStatus(cs_models.UserInteractStatusOnMicrophone)
	err = is.interactDetailRepo.Save(ctx, interactDetail)
	if err != nil {
		logs.CtxError(ctx, "save interact detail failed,error:%s", err)
		return custom_error.InternalError(err)
	}

	informer := GetInformService()
	informer.BroadcastRoom(ctx, is.interact.OwnerRoomID, cs_models.OnCsMicOn, User2CsUserInfo(audience))
	return nil
}

func (is *InteractService) Finish(ctx context.Context, audienceID string) error {
	audience, err := is.roomUserRepo.GetActiveUserByRoomIDUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if err != nil || audience == nil {
		if err == nil {
			err = errors.New("user is not exist")
		}
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return custom_error.InternalError(err)
	}
	audience.SetInteractStatus(cs_models.UserInteractStatusAudience)
	err = is.roomUserRepo.Save(ctx, audience)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return custom_error.InternalError(err)
	}

	interactDetail, err := is.interactDetailRepo.GetInteractDetailByFromRoomIDToUserID(ctx, is.interact.GetOwnerRoomID(), audienceID)
	if interactDetail != nil {
		interactDetail.SetStatus(cs_models.UserInteractStatusAudience)
		err = is.interactDetailRepo.Save(ctx, interactDetail)
		if err != nil {
			logs.CtxError(ctx, "save interact detail failed,error:%s", err)
			return custom_error.InternalError(err)
		}
	}

	informer := GetInformService()
	informer.BroadcastRoom(ctx, is.interact.OwnerRoomID, cs_models.OnCsMicOff, User2CsUserInfo(audience))
	return nil
}
