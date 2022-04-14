package cs_facade

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_implement"
)

var interactRepository *cs_implement.InteractRepositoryImpl

func GetInteractRepository() *cs_implement.InteractRepositoryImpl {
	if interactRepository == nil {
		interactRepository = &cs_implement.InteractRepositoryImpl{}
	}
	return interactRepository
}

type InteractRepositoryInterface interface {
	Save(ctx context.Context, interact *cs_entity.CsInteract) error

	GetInteractByRoomID(ctx context.Context, roomID string) (*cs_entity.CsInteract, error)
}
