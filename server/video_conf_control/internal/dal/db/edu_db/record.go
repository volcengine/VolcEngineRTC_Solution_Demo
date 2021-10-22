package edu_db

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/pkg/public"
)

const EduRecordInfo = "edu_record_info"

const (
	RecordStart   = 0
	RecordSuccess = 1
	RecordFail    = 2
	RecordDelete  = 3
)

type RecordInfo struct {
	VID       string `json:"Vid"`
	Duration  int64  `json:"Duration"`  // 单位ms
	Size      int64  `json:"Size"`      // 单位ms
	StartTime int64  `json:"StartTime"` // 单位ms
}

func StartRecord(ctx context.Context, appID, roomID, userID, taskID, roomName string) error {
	defer public.CheckPanic()
	record := &db.EduRecordInfo{
		AppID:           appID,
		RoomID:          roomID,
		UserID:          userID,
		RoomName:        roomName,
		RecordStatus:    RecordStart,
		CreatedTime:     db.EduTime(time.Now()),
		UpdatedTime:     db.EduTime(time.Now()),
		TaskID:          taskID,
		RecordBeginTime: time.Now().UnixNano(),
	}
	return db.Client.WithContext(ctx).Table(EduRecordInfo).Create(record).Error
}

func FinishRecord(ctx context.Context, taskID string, status int, records []*RecordInfo) error {
	defer public.CheckPanic()
	if status == RecordFail {
		db.Client.WithContext(ctx).Table(EduRecordInfo).Where("task_id = ? and record_status = ?", taskID, RecordStart).Updates(map[string]interface{}{
			"record_status": RecordFail,
		})
		return nil
	}
	recordBase := &db.EduRecordInfo{}
	err := db.Client.WithContext(ctx).Table(EduRecordInfo).Where("task_id = ? and record_status = ?", taskID, RecordStart).First(&recordBase).Error
	if err != nil {
		return err
	}
	for _, r := range records {
		record := &db.EduRecordInfo{
			AppID:           recordBase.AppID,
			RoomID:          recordBase.RoomID,
			ParentRoomID:    recordBase.ParentRoomID,
			RoomName:        recordBase.RoomName,
			UserID:          recordBase.UserID,
			RecordBeginTime: r.StartTime * 1e6,
			RecordEndTime:   (r.StartTime + r.Duration) * 1e6,
			CreatedTime:     recordBase.CreatedTime,
			TaskID:          recordBase.TaskID,
			UpdatedTime:     db.EduTime(time.Now()),
			RecordStatus:    RecordSuccess,
			Vid:             r.VID,
		}
		err := db.Client.WithContext(ctx).Table(EduRecordInfo).Create(record).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteRecord(ctx context.Context, vid string) error {
	defer public.CheckPanic()
	err := db.Client.WithContext(ctx).Table(EduRecordInfo).Where("vid=?", vid).Updates(map[string]interface{}{
		"record_status": RecordDelete,
	}).Error
	return err
}

func QueryRecordByRoomID(ctx context.Context, roomID string) ([]*db.EduRecordInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduRecordInfo, 0)
	err := db.Client.WithContext(ctx).Table(EduRecordInfo).Where("room_id = ? and record_status= ?", roomID, RecordSuccess).
		Order("create_time desc").Find(&rs).Error
	return rs, err
}

func QueryRecordByVid(ctx context.Context, vid string) (*db.EduRecordInfo, error) {
	defer public.CheckPanic()
	var rs *db.EduRecordInfo
	err := db.Client.WithContext(ctx).Table(EduRecordInfo).Where("vid = ?", vid).First(&rs).Error
	return rs, err
}

func QueryTaskId(ctx context.Context, appID, roomID, userID string) (string, error) {
	defer public.CheckPanic()
	recordInfo := db.EduRecordInfo{}
	err := db.Client.WithContext(ctx).Table(EduRecordInfo).Where("app_id=? and room_id =? and user_id =? and record_status=?", appID, roomID, userID, RecordStart).First(&recordInfo).Error
	return recordInfo.TaskID, err
}

func QueryStartRecordByTaskID(ctx context.Context, taskID string) (*db.EduRecordInfo, error) {
	defer public.CheckPanic()
	recordBase := &db.EduRecordInfo{}
	err := db.Client.Debug().WithContext(ctx).Table(EduRecordInfo).Where("task_id = ? and record_status = ?", taskID, RecordStart).First(&recordBase).Error
	if err != nil {
		return nil, err
	}
	return recordBase, nil
}

func QueryAllTimeOutRecord(ctx context.Context) ([]*db.EduRecordInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduRecordInfo, 0)
	t := time.Now().Add(-15 * time.Minute).UnixNano()
	err := db.Client.WithContext(ctx).Table(EduRecordInfo).Where("record_status=? and record_begin_time<?", RecordStart, t).Find(&rs).Error
	return rs, err
}
