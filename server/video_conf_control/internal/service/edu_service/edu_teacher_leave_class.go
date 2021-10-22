package edu_service

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type teacherLeaveClassReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

func teacherLeaveRoom(ctx context.Context, param *vc_control.TEventParam) {

}
