package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/token"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type getGroupClassInfoReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

type groupTokenInfo struct {
	RoomID  string `json:"room_id"`
	RoomIdx int    `json:"room_idx"`
	Token   string `json:"token"`
}

func getGroupClassInfo(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("getGroupClassInfo:%+v", param)
	var p getGroupClassInfoReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.RoomID == "" || p.UserID == "" {
		logs.Warnf("input room_id or user_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	groupRoomIDSet, err := room.ListGroupRoomID(ctx)
	if err != nil {
		logs.Errorf("get group room id error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	res := make([]*groupTokenInfo, 0)
	for _, groupRoomID := range groupRoomIDSet {
		tokenParam := &token.GenerateParam{
			AppID:        config.Config.AppID,
			AppKey:       config.Config.AppKey,
			RoomID:       groupRoomID,
			UserID:       p.UserID,
			ExpireAt:     7 * 24 * 3600,
			CanPublish:   true,
			CanSubscribe: true,
			CanSwitch:    true,
		}
		generateToken, err := token.GenerateToken(tokenParam)
		logs.Infof("GenerateToken:%+v,%s", tokenParam, generateToken)
		if err != nil {
			logs.Errorf("failed to GenerateToken: %s", err)
			service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
			return
		}
		res = append(res, &groupTokenInfo{
			RoomID:  groupRoomID,
			RoomIdx: room.GetIdxByGroupRoomID(ctx, groupRoomID),
			Token:   generateToken,
		})
	}
	service_utils.Push2Client(ctx, param, err, res)
}
