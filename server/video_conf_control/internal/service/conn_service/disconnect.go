package conn_service

import (
	"context"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/edu_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/vc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/vc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

func disconnect(ctx context.Context, param *vc_control.TEventParam) {
	_, err := vc_db.GetUserByConnID(ctx, param.ConnId)
	if err == nil {
		vc_service.Disconnect(ctx, param)

		return
	}

	user, err := edu_db.GetActiveUserByConnID(ctx, param.ConnId)
	if err == nil && user != nil {
		edu_service.Disconnect(ctx, param)

		return
	}
}
