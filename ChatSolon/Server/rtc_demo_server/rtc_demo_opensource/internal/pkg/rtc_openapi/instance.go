package rtc_openapi

import "github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"

var instance *RTC

func GetInstance() *RTC {
	if instance == nil {
		instance = NewInstance()
		instance.SetAccessKey(config.Configs().VolcAk)
		instance.SetSecretKey(config.Configs().VolcSk)
	}
	return instance
}
