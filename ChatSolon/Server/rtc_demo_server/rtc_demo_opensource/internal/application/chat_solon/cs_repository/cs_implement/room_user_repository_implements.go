package cs_implement

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const RoomUserTable = "cs_room_user"

type RoomUserRepositoryImpl struct {
}

func (ruri *RoomUserRepositoryImpl) Save(ctx context.Context, user *cs_entity.CsRoomUser) error {
	user.UpdateTime = time.Now()
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}, {Name: "user_id"}},
			UpdateAll: true,
		}).Create(&user).Error
	return err
}

func (ruri *RoomUserRepositoryImpl) GetActiveUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*cs_entity.CsRoomUser, error) {
	var rs *cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?  and user_id= ?  and net_status=?", roomID, userID, cs_entity.UserNetStatusOnline).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (ruri *RoomUserRepositoryImpl) GetReconnectUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*cs_entity.CsRoomUser, error) {
	var rs *cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?  and user_id= ?  and net_status in ?", roomID, userID, []int{cs_entity.UserNetStatusReconnecting, cs_entity.UserNetStatusOnline}).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (ruri *RoomUserRepositoryImpl) GetActiveUsersByRoomID(ctx context.Context, roomID string) ([]*cs_entity.CsRoomUser, error) {
	var rs []*cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?   and net_status= ?", roomID, cs_entity.UserNetStatusOnline).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (ruri *RoomUserRepositoryImpl) GetActiveUsersByRoomIDStatus(ctx context.Context, roomID string, status []int) ([]*cs_entity.CsRoomUser, error) {
	var rs []*cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?  and interact_status in ?  and net_status= ?", roomID, status, cs_entity.UserNetStatusOnline).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
func (ruri *RoomUserRepositoryImpl) GetActiveAudiencesByRoomIDStatus(ctx context.Context, roomID string, status []int) ([]*cs_entity.CsRoomUser, error) {
	var rs []*cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?  and interact_status in ?  and net_status= ? and user_role = ?", roomID, status, cs_entity.UserNetStatusOnline, cs_entity.UserRoleAudience).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (ruri *RoomUserRepositoryImpl) GetActiveUsersByRoomIDRole(ctx context.Context, roomID string, role int) ([]*cs_entity.CsRoomUser, error) {
	var rs []*cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?  and user_role= ?  and net_status= ? and interact_status = ?", roomID, role, cs_entity.UserNetStatusOnline, cs_models.UserInteractStatusAudience).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (ruri *RoomUserRepositoryImpl) UpdateUsersByUserID(ctx context.Context, userID string, ups map[string]interface{}) error {
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("user_id= ? ", userID).Updates(ups).Error
	return err
}

func (ruri *RoomUserRepositoryImpl) UpdateUsersByRoomID(ctx context.Context, roomID string, ups map[string]interface{}) error {
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ? ", roomID).Updates(ups).Error
	return err
}

func (ruri *RoomUserRepositoryImpl) GetUserByStatusOrderUtime(ctx context.Context, roomID string, status int) (*cs_entity.CsRoomUser, error) {
	var rs *cs_entity.CsRoomUser
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ?  and interact_status = ?  and net_status= ?", roomID, status, cs_entity.UserNetStatusOnline).Order("update_time").First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (ruri *RoomUserRepositoryImpl) GetUserCountByRoomID(ctx context.Context, roomID string) (int, error) {
	var count int64
	err := db.Client.WithContext(ctx).Debug().Table(RoomUserTable).Where("room_id= ? and net_status= ?", roomID, cs_entity.UserNetStatusOnline).Count(&count).Error
	return int(count), err
}
