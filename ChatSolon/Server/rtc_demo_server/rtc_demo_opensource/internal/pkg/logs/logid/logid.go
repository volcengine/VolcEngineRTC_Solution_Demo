package logid

import "github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/uuid"

func GenLogID() string {
	return uuid.GetUUID()
}
