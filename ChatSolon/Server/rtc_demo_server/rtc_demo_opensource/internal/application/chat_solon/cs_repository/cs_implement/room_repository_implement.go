package cs_implement

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const RoomTable = "cs_room"

type RoomRepositoryImpl struct {
}

func (rri *RoomRepositoryImpl) Save(ctx context.Context, room *cs_entity.CsRoom) error {
	room.UpdateTime = time.Now()
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}},
			UpdateAll: true,
		}).Create(&room).Error
	return err
}

func (rri *RoomRepositoryImpl) GetRoomByRoomID(ctx context.Context, roomID string) (*cs_entity.CsRoom, error) {
	var rs *cs_entity.CsRoom
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).Where("room_id= ?   and status = ?", roomID, cs_entity.RoomStatusStart).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (rri *RoomRepositoryImpl) GetRooms(ctx context.Context) ([]*cs_entity.CsRoom, error) {
	var rs []*cs_entity.CsRoom
	err := db.Client.WithContext(ctx).Debug().Table(RoomTable).Where(" status = ?", cs_entity.RoomStatusStart).Order("create_time desc").Find(&rs).Error
	return rs, err
}
