package svc_db

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/util"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const UserTable = "svc_user"

const (
	UserRoleHost     = 1
	UserRoleAudience = 2
)
const (
	UserNetStatusOnline       = 1
	UserNetStatusOffline      = 2
	UserNetStatusReconnecting = 3
)

const (
	UserInteractStatusNormal      = 1
	UserInteractStatusInteracting = 2
	UserInteractStatusApplying    = 3
	UserInteractStatusInviting    = 4
)

type DbUserRepo struct{}

func (repo *DbUserRepo) Save(ctx context.Context, user *svc_entity.SvcUser) error {
	defer util.CheckPanic()
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}, {Name: "user_id"}},
			UpdateAll: true,
		}).Create(&user).Error
	return err
}

func (repo *DbUserRepo) UpdateUsersWithMapByRoomID(ctx context.Context, roomID string, ups map[string]interface{}) error {
	defer util.CheckPanic()
	return db.Client.WithContext(ctx).Debug().Table(UserTable).Where("room_id = ?", roomID).Updates(ups).Error
}

func (repo *DbUserRepo) GetActiveUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*svc_entity.SvcUser, error) {
	defer util.CheckPanic()
	var rs *svc_entity.SvcUser
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).Where("room_id = ? and user_id = ? and net_status in ?", roomID, userID, []int{UserNetStatusOnline, UserNetStatusReconnecting}).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (repo *DbUserRepo) GetUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*svc_entity.SvcUser, error) {
	defer util.CheckPanic()
	var rs *svc_entity.SvcUser
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).Where("room_id = ? and user_id = ?", roomID, userID).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (repo *DbUserRepo) GetActiveUserByUserID(ctx context.Context, userID string) (*svc_entity.SvcUser, error) {
	defer util.CheckPanic()
	var rs *svc_entity.SvcUser
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).Where("user_id = ? and net_status in ?", userID, []int{UserNetStatusOnline, UserNetStatusReconnecting}).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (repo *DbUserRepo) GetActiveUsersByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcUser, error) {
	defer util.CheckPanic()
	var rs []*svc_entity.SvcUser
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).Where("room_id = ? and net_status = ?", roomID, UserNetStatusOnline).Order("create_time desc").Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (repo *DbUserRepo) GetAudiencesWithoutApplyByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcUser, error) {
	defer util.CheckPanic()
	var rs []*svc_entity.SvcUser
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).Where("room_id = ? and net_status = ? and interact_status <> ? and user_role = ?", roomID, UserNetStatusOnline, UserInteractStatusApplying, UserRoleAudience).Order("create_time desc").Limit(200).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (repo *DbUserRepo) GetApplyAudiencesByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcUser, error) {
	defer util.CheckPanic()
	var rs []*svc_entity.SvcUser
	err := db.Client.WithContext(ctx).Debug().Table(UserTable).Where("room_id= ? and net_status = ? and interact_status = ? and user_role = ?", roomID, UserNetStatusOnline, UserInteractStatusApplying, UserRoleAudience).Order("create_time desc").Limit(200).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
