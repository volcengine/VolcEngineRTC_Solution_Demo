package rpc

import (
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/rpc/frontier"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/rpc/howard"
)

func Init() {
	frontier.Init()
	howard.Init()
}
