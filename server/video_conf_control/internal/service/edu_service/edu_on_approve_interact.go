package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type micReq struct {
	UserID     string `json:"user_id"`
	RoomID     string `json:"room_id"`
	LoginToken string `json:"login_token"`
}

func onApproveInteract(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("onApproveMic:%+v", param)
	var p micReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.UserID == "" || p.RoomID == "" {
		logs.Warnf("input user_id or room_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	user, err := edu_models.GetUserByRoomIDUserID(ctx, p.RoomID, p.UserID)
	if err != nil || user == nil {
		logs.Errorf("get user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	//是否重复给麦
	if user.IsInteract() {
		logs.Errorf("user already mic on")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrorDuplicateApproveMic)
		return
	}

	//校验是否已有6个麦
	if room.GetInteractUserCount(ctx) >= 6 {
		logs.Errorf("approve mic limit 6")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrorReachMicOnLimit)
		return
	}

	//校验学生是否已经取消举手
	if !user.IsHandsUp() {
		logs.Errorf("user already cancel hands up")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrorStudentNotHandsUp)
		return
	}

	user.StartInteract()
	err = user.Save(ctx)
	if err != nil {
		logs.Errorf("save user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	edu_models.StartRecord(ctx, config.Config.AppID, room.GetRoomID(), user.GetUserID(), room.GetRoomName())
	room.InformRoom(ctx, edu_models.OnStartInteract, edu_models.NoticeRoom{RoomID: room.GetRoomID(), UserID: user.GetUserID(), UserName: user.GetUserName()})
	service_utils.Push2Client(ctx, param, err, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
}
