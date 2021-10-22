package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/video"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type getHistoryRecordListReq struct {
	RoomID     string `json:"room_id"`
	LoginToken string `json:"login_token"`
}

func getHistoryRecordList(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("getHistoryRecordList:%+v", param)
	var p getHistoryRecordListReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.RoomID == "" {
		logs.Warnf("input user_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//查record
	res := make([]*db.EduRecordInfo, 0)
	records, err := edu_db.QueryRecordByRoomID(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("failed to QueryRecordByRoomID: %s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	vids := make([]string, 0)
	for _, record := range records {
		vids = append(vids, record.Vid)
	}
	durl := video.GetVideoURL(ctx, vids)
	logs.Infof("vids: %v, downloadurl: %v", vids, durl)

	for _, r := range records {
		r.VideoURL = durl[r.Vid]
		res = append(res, r)
	}

	service_utils.Push2Client(ctx, param, err, res)
}
