package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type getHistoryRoomListReq struct {
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

func getHistoryRoomList(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("getHistoryRoomList:%+v", param)
	var p getHistoryRoomListReq
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

	//查大房间
	rooms, err := edu_models.GetHistoryRoomsByUserID(ctx, p.UserID)
	if err != nil {
		logs.Errorf("failed to QueryActiveRooms: %s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	res := make([]*db.EduRoomInfo, 0)

	for _, r := range rooms {
		records, _ := edu_db.QueryRecordByRoomID(ctx, r.RoomID)
		if len(records) != 0 {
			res = append(res, r)
		}
	}

	service_utils.Push2Client(ctx, param, err, res)
}
