package cs_implement

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const InteractTable = "cs_interact"

type InteractRepositoryImpl struct {
}

func (iri *InteractRepositoryImpl) Save(ctx context.Context, interact *cs_entity.CsInteract) error {
	err := db.Client.WithContext(ctx).Debug().Table(InteractTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "interact_id"}},
			UpdateAll: true,
		}).Create(&interact).Error
	return err
}

func (iri *InteractRepositoryImpl) GetInteractByRoomID(ctx context.Context, roomID string) (*cs_entity.CsInteract, error) {
	var rs *cs_entity.CsInteract
	err := db.Client.WithContext(ctx).Debug().Table(InteractTable).Where("owner_room_id= ? ", roomID).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
