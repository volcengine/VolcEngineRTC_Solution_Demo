package record

import (
	"context"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/vc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/redis/vc_redis"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"

	logs "github.com/sirupsen/logrus"
)

type PostProcessingCallBack struct {
	AppID          string       `json:"AppId"`
	BussinessID    string       `json:"BusinessId"`
	RoomID         string       `json:"RoomId"`
	TaskID         string       `json:"TaskId"`
	Code           int          `json:"code"`
	ErrorMessage   string       `json:"ErrorMessage"`
	RecordFileList []recordFile `json:"RecordFileList"`
}

func HandleRecordCallback(ctx context.Context, param *vc_control.TRecordCallbackParam) (resp *vc_control.THTTPResp, err error) {
	resp = vc_control.NewTHTTPResp()

	// 只处理录制完成的回调
	var p PostProcessingCallBack
	if err2 := json.Unmarshal([]byte(param.EventData), &p); err2 != nil {
		return
	}

	logs.Infof("HandleRecordCallback, PostProcessingCallBack: %v", p)

	if len(p.RecordFileList) == 0 {
		return
	}

	//vc
	userID := vc_redis.GetUserIDByTaskID(ctx, p.TaskID)
	if userID != "" {
		// 对于只需要混流录制的场景，取返回值的第一个就可以
		r := &db.MeetingVideoRecord{
			AppID:  p.AppID,
			RoomID: p.RoomID,
			VID:    p.RecordFileList[0].VID,
			UserID: userID,
			State:  db.ACTIVE,
		}

		logs.Infof("HandleRecordCallback, MeetingVideoRecord: %v, userID: %v", *r, userID)

		if err := vc_db.CreateRecord(ctx, r); err != nil {
			logs.Warnf("create meeting record error, err: %v", err)
		}
	}

	//edu
	_, err = edu_db.QueryStartRecordByTaskID(ctx, p.TaskID)
	if err == nil {
		var status int
		if p.Code == 0 {
			status = edu_db.RecordSuccess
		} else {
			status = edu_db.RecordFail
		}
		recordInfo := make([]*edu_db.RecordInfo, 0)
		for _, r := range p.RecordFileList {
			recordInfo = append(recordInfo, &edu_db.RecordInfo{
				VID:       r.VID,
				StartTime: r.StartTime,
				Duration:  r.Duration,
				Size:      r.Size,
			})
		}
		logs.Infof("handle edu record ,status:%v,record_info:%v", status, recordInfo)
		if len(recordInfo) > 0 {
			err = edu_db.FinishRecord(ctx, p.TaskID, status, recordInfo)
		}
	}

	return
}

func HandleRecordCallbackHttp(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	// 只处理录制完成的回调
	var p PostProcessingCallBack
	if err2 := json.Unmarshal([]byte(body), &p); err2 != nil {
		return
	}

	logs.Infof("HandleRecordCallback, PostProcessingCallBack: %v", p)

	if len(p.RecordFileList) == 0 {
		return
	}

	//vc
	userID := vc_redis.GetUserIDByTaskID(ctx, p.TaskID)
	if userID != "" {
		// 对于只需要混流录制的场景，取返回值的第一个就可以
		r := &db.MeetingVideoRecord{
			AppID:  p.AppID,
			RoomID: p.RoomID,
			VID:    p.RecordFileList[0].VID,
			UserID: userID,
			State:  db.ACTIVE,
		}

		logs.Infof("HandleRecordCallback, MeetingVideoRecord: %v, userID: %v", *r, userID)

		if err := vc_db.CreateRecord(ctx, r); err != nil {
			logs.Warnf("create meeting record error, err: %v", err)
		}
	}

	//edu
	_, err := edu_db.QueryStartRecordByTaskID(ctx, p.TaskID)
	if err == nil {
		var status int
		if p.Code == 0 {
			status = edu_db.RecordSuccess
		} else {
			status = edu_db.RecordFail
		}
		recordInfo := make([]*edu_db.RecordInfo, 0)
		for _, r := range p.RecordFileList {
			recordInfo = append(recordInfo, &edu_db.RecordInfo{
				VID:       r.VID,
				StartTime: r.StartTime,
				Duration:  r.Duration,
				Size:      r.Size,
			})
		}
		logs.Infof("handle edu record ,status:%v,record_info:%v", status, recordInfo)
		if len(recordInfo) > 0 {
			err = edu_db.FinishRecord(ctx, p.TaskID, status, recordInfo)
		}
	}

	return
}
