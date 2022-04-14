package cs_implement

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const InteractDetailTable = "cs_interact_detail"

type InteractDetailRepositoryImpl struct {
}

func (idri *InteractDetailRepositoryImpl) Save(ctx context.Context, interactDetail *cs_entity.CsInteractDetail) error {
	interactDetail.UpdateTime = time.Now()
	err := db.Client.WithContext(ctx).Debug().Table(InteractDetailTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "interact_id"}, {Name: "from_room_id"}, {Name: "from_user_id"}, {Name: "to_room_id"}, {Name: "to_user_id"}},
			UpdateAll: true,
		}).Create(&interactDetail).Error
	return err
}

func (idri *InteractDetailRepositoryImpl) GetInteractDetailByFromRoomIDToUserID(ctx context.Context, fromRoomID, toUserID string) (*cs_entity.CsInteractDetail, error) {
	var rs *cs_entity.CsInteractDetail
	err := db.Client.WithContext(ctx).Debug().Table(InteractDetailTable).Where("from_room_id= ? and to_user_id = ? ", fromRoomID, toUserID).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
