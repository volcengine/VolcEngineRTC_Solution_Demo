package edu_service

import (
	"context"
	"encoding/json"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type getCreatedClassReq struct {
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

type getCreatedClassResp struct {
	RoomList []*db.EduRoomInfo `json:"room_list"`
	Token    string            `json:"token"`
}

func getCreatedClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("getCreatedClass:%+v", param)
	var p getCreatedClassReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.UserID == "" {
		logs.Warnf("input user_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	rooms, err := edu_models.GetCreatedRooms(ctx, p.UserID)
	if err != nil {
		logs.Errorf("get created rooms failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	token := ""
	if len(rooms) > 0 {
		token, err = genToken(rooms[0].RoomID, p.UserID)
		if err != nil {
			logs.Errorf("gen token failed,error:%s", err)
			service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
			return
		}
	}

	res := &getCreatedClassResp{
		RoomList: rooms,
		Token:    token,
	}

	service_utils.Push2Client(ctx, param, err, res)
}
