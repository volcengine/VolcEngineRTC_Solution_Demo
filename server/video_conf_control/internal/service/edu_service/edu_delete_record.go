package edu_service

import (
	"context"
	"encoding/json"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/video"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type deleteRecordReq struct {
	Vid        string `json:"vid"`
	LoginToken string `json:"login_token"`
}

func deleteRecord(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("deleteRecord:%+v", param)
	var p deleteRecordReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.Vid == "" {
		logs.Warnf("input vid error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}
	r, err := edu_db.QueryRecordByVid(ctx, p.Vid)
	if err != nil || r == nil {
		logs.Errorf("record not found, vid:%s, err:%s", p.Vid, err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	err = vCloudDeleteVideoRecord(ctx, p.Vid)
	if err != nil {
		logs.Errorf("delete record from vcloud failed,vid:%s,err:%s", p.Vid, err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	err = edu_db.DeleteRecord(ctx, p.Vid)
	if err != nil {
		logs.Errorf("delete record failed,vid:%s,err:%s", p.Vid, err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	service_utils.Push2Client(ctx, param, err, p.Vid)

}

func vCloudDeleteVideoRecord(ctx context.Context, vid string) error {
	_, err := video.DeleteRecord(ctx, vid)
	if err != nil {
		return err
	}
	return nil
}
