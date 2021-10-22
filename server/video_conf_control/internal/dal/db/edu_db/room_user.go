package edu_db

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/pkg/public"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	UserStatusOnline       = 0
	UserStatusOffline      = 1
	UserStatusReconnecting = 2
)

const (
	UserRoleTeacher = 0
	UserRoleStudent = 1
)

const RoomUserTable = "edu_user_room_info"

func GetActiveUserByConnID(ctx context.Context, connID string) (*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	var rs *db.EduUserRoomInfo
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("user_status = ? and conn_id = ?", UserStatusOnline, connID).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func GetReconnectingUserByDeviceID(ctx context.Context, deviceID string) (*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	var rs *db.EduUserRoomInfo
	err := db.Client.WithContext(ctx).Table(RoomUserTable).
		Where("device_id = ? and user_status = ?", deviceID, UserStatusReconnecting).
		First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func CreateOrUpdateUser(ctx context.Context, user *db.EduUserRoomInfo) error {
	defer public.CheckPanic()
	err := db.Client.WithContext(ctx).Table(RoomUserTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}, {Name: "user_id"}},
			UpdateAll: true,
		}).Create(&user).Error
	return err
}

func GetInteractStudentsCountByParentRoomID(ctx context.Context, parentRoomID string) (int, error) {
	defer public.CheckPanic()
	var count int64
	err := db.Client.WithContext(ctx).Table(RoomUserTable).Where("parent_room_id = ? and is_interact = 1 and user_status in ?", parentRoomID, []int{UserStatusOnline, UserStatusReconnecting}).Count(&count).Error
	return int(count), err
}

func GetInteractStudentsByParentRoomID(ctx context.Context, parentRoomID string) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).Where("parent_room_id = ? and is_interact = 1 and user_status in ?", parentRoomID, []int{UserStatusOnline, UserStatusReconnecting}).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, nil
		}
		return nil, err
	}
	return rs, nil
}

func GetHandsUpStudentsByParentRoomID(ctx context.Context, parentRoomID string) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).Where("parent_room_id = ? and is_hands_up = 1", parentRoomID).Limit(200).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, nil
		}
		return nil, err
	}
	return rs, nil
}

func UpdateUsersByParentRoomID(ctx context.Context, parentRoomID string, ups map[string]interface{}) error {
	defer public.CheckPanic()
	return db.Client.WithContext(ctx).Table(RoomUserTable).Where("parent_room_id = ?", parentRoomID).Updates(ups).Error
}

func GetUserByParentRoomIDUserID(ctx context.Context, parentRoomID, userID string) (*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	var rs *db.EduUserRoomInfo
	err := db.Client.WithContext(ctx).Table(RoomUserTable).Where("parent_room_id = ? and user_id = ?", parentRoomID, userID).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func GetActiveUsersByUserID(ctx context.Context, userID string) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).Where("user_id = ? and user_status in ?", userID, []int{UserStatusOnline, UserStatusReconnecting}).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func GetActiveStudentCountByParentRoomID(ctx context.Context, parentRoomID string) (int, error) {
	defer public.CheckPanic()
	var count int64
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("parent_room_id = ? and user_role = ? and user_status in ?", parentRoomID, UserRoleStudent, []int{UserStatusOnline, UserStatusReconnecting}).Count(&count).Error
	return int(count), err
}

//pageNumber==0 and pageSize =0 means get all
func GetActiveStudentsByParentRoomID(ctx context.Context, parentRoomID string, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	dbt := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).
		Where("parent_room_id = ? and user_role = ? and user_status in ?", parentRoomID, UserRoleStudent, []int{UserStatusOnline, UserStatusReconnecting})
	if pageNumber != 0 && pageSize != 0 {
		dbt = dbt.Offset((pageNumber - 1) * pageSize).Limit(pageSize)
	}
	err := dbt.Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

//pageNumber==0 and pageSize =0 means get all
func GetActiveStudentsByRoomID(ctx context.Context, roomID string, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	dbt := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).
		Where("room_id = ? and user_role = ? and user_status in ?", roomID, UserRoleStudent, []int{UserStatusOnline, UserStatusReconnecting})
	if pageNumber != 0 && pageSize != 0 {
		dbt = dbt.Offset((pageNumber - 1) * pageSize).Limit(pageSize)
	}
	err := dbt.Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func GetActiveStudentsByRoomIDRange(ctx context.Context, minRoomID, maxRoomID string) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).
		Where("room_id  between ? and ? and  user_role = ? and user_status in ?", minRoomID, maxRoomID, UserRoleStudent, []int{UserStatusOnline, UserStatusReconnecting}).
		Order("create_time").
		Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func GetActiveStudentsByRoomIDSet(ctx context.Context, roomIDSet []string) ([]*db.EduUserRoomInfo, error) {
	defer public.CheckPanic()
	rs := make([]*db.EduUserRoomInfo, 0)
	err := db.Client.WithContext(ctx).Table(RoomUserTable).
		Where("room_id in ? and user_role = ? and user_status in ?", roomIDSet, UserRoleStudent, []int{UserStatusOnline, UserStatusReconnecting}).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
