package cs_facade

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_implement"
)

var interactDetailRepository *cs_implement.InteractDetailRepositoryImpl

func GetInteractDetailRepository() *cs_implement.InteractDetailRepositoryImpl {
	if interactDetailRepository == nil {
		interactDetailRepository = &cs_implement.InteractDetailRepositoryImpl{}
	}
	return interactDetailRepository
}

type InteractDetailRepositoryInterface interface {
	Save(ctx context.Context, room *cs_entity.CsInteractDetail) error

	GetInteractDetailByFromRoomIDToUserID(ctx context.Context, fromRoomID, toUserID string) (*cs_entity.CsInteractDetail, error)
}
