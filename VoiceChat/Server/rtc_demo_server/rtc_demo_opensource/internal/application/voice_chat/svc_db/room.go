package svc_db

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/util"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const RoomTable = "svc_room"

const (
	RoomStatusPrepare = 1
	RoomStatusStart   = 2
	RoomStatusFinish  = 3
)

type DbRoomRepo struct{}

func (repo *DbRoomRepo) Save(ctx context.Context, liveRoom *svc_entity.SvcRoom) error {
	defer util.CheckPanic()
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}},
			UpdateAll: true,
		}).Create(&liveRoom).Error
	return err
}

func (repo *DbRoomRepo) GetRoomByRoomID(ctx context.Context, roomID string) (*svc_entity.SvcRoom, error) {
	defer util.CheckPanic()
	var rs *svc_entity.SvcRoom
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).Where("room_id = ? and status <> ?", roomID, RoomStatusFinish).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (repo *DbRoomRepo) GetActiveRooms(ctx context.Context) ([]*svc_entity.SvcRoom, error) {
	defer util.CheckPanic()
	var rs []*svc_entity.SvcRoom
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).Where("status = ?", RoomStatusStart).Order("create_time desc").Limit(200).Find(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
