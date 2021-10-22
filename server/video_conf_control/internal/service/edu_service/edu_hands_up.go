package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type handsUpReq struct {
	UserID     string `json:"user_id"`
	RoomID     string `json:"room_id"`
	LoginToken string `json:"login_token"`
}

func handsUp(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("handsUp:%+v", param)
	var p handsUpReq
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
	//检查是否开启集体发言
	if !room.IsEnableInteract() {
		logs.Infof("room  not allow group speech")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrorForbidHandsUp)
		return
	}

	user, err := edu_models.GetUser(ctx, param.ConnId)
	if err != nil {
		logs.Errorf("get user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	user.HandsUp()
	err = user.Save(ctx)
	if err != nil {
		logs.Errorf("hands up failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	service_utils.Push2Client(ctx, param, err, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
}
