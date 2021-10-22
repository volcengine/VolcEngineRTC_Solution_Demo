package edu_db

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/pkg/public"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const RoomTable = "edu_room_info"
const RoomChildTable = "edu_room_child_info"

const (
	RoomTypeSingleRoomClass = 0
	RoomTypeGroupRoomClass  = 1
)

const (
	ClassPending = 0
	ClassRunning = 1
	ClassFinish  = 2
	ClassDelete  = 3
)

func GetActiveRoomByRoomID(ctx context.Context, roomID string) (*db.EduRoomInfo, error) {
	defer public.CheckPanic()
	var rs *db.EduRoomInfo
	err := db.Client.WithContext(ctx).Table(RoomTable).Where("room_id = ? and status in ?", roomID, []int{ClassPending, ClassRunning}).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func UpdateRoomWithMap(ctx context.Context, roomID string, ups map[string]interface{}) error {
	defer public.CheckPanic()
	return db.Client.WithContext(ctx).Table(RoomTable).Where("room_id = ?", roomID).Updates(ups).Error
}

func GetCreatedRooms(ctx context.Context, userID string) ([]*db.EduRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomTable).Where("create_user_id = ? and status in ?", userID, []int{ClassPending, ClassRunning}).
		Order("create_time desc").Find(&rs).Error
	return rs, err
}

func GetActiveRooms(ctx context.Context) ([]*db.EduRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomTable).Where("status in ?", []int{ClassPending, ClassRunning}).
		Order("create_time desc").Find(&rs).Error
	return rs, err
}

func GetHistoryRoomsByUserID(ctx context.Context, userID string) ([]*db.EduRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduRoomInfo, 0)
	parentRoomIds := make([]string, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).Select("parent_room_id").Where("user_id=?", userID).Group("parent_room_id").Find(&parentRoomIds).Error
	if err != nil {
		return rs, err
	}
	err = db.Client.WithContext(ctx).Table(RoomTable).Where("status = ? and room_id in ?", ClassFinish, parentRoomIds).Find(&rs).Error
	return rs, err
}

func CreateOrUpdateRoom(ctx context.Context, room *db.EduRoomInfo) error {
	defer public.CheckPanic()
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}},
			UpdateAll: true,
		}).Create(&room).Error
	return err
}

func GetActiveRoomsByCreateUserID(ctx context.Context, createUserID string) ([]*db.EduRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomTable).Where("create_user_id = ? and status in ?", createUserID, []int{ClassPending, ClassRunning}).
		Order("create_time desc").Find(&rs).Error
	return rs, err
}

func GetActiveGroupRoomIDSet(ctx context.Context, parentRoomID string) ([]string, error) {
	defer public.CheckPanic()
	type gRoom struct {
		GroupRoomID string `gorm:"column:room_id"`
		Count       int64  `gorm:"column:c"`
	}
	rs := make([]*gRoom, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).
		Select("room_id,count(id) as c").
		Where("parent_room_id = ? and user_role = ? and user_status in ?", parentRoomID, UserRoleStudent, []int{UserStatusOnline, UserStatusReconnecting}).
		Group("room_id").
		Having("c>0").
		Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := make([]string, 0)
	for _, r := range rs {
		res = append(res, r.GroupRoomID)
	}
	return res, nil
}
