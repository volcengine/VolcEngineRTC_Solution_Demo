package rtc_openapi

import (
	"net/http"
	"net/url"
	"time"

	"github.com/volcengine/volc-sdk-golang/base"
)

const (
	DefaultRegion          = base.RegionCnNorth1
	ServiceVersion20201201 = "2020-12-01"
	ServiceName            = "rtc"
	ServiceHost            = "open.volcengineapi.com"
)

const (
	startRecord     = "StartRecord"
	updateRecord    = "UpdateRecord"
	sendUnicast     = "SendUnicast"     //房间外点对点消息
	sendRoomUnicast = "SendRoomUnicast" //房间内点对点消息
	sendBroadcast   = "SendBroadcast"   //房间内广播消息
)

var (
	serviceInfo = &base.ServiceInfo{
		Timeout: 5 * time.Second,
		Host:    ServiceHost,
		Header: http.Header{
			"Accept": []string{"application/json"},
		},
		Credentials: base.Credentials{Region: DefaultRegion, Service: ServiceName},
	}

	defaultApiInfoList = map[string]*base.ApiInfo{
		sendUnicast: {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{sendUnicast},
				"Version": []string{ServiceVersion20201201},
			},
		},
		sendRoomUnicast: {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{sendRoomUnicast},
				"Version": []string{ServiceVersion20201201},
			},
		},
		sendBroadcast: {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{sendBroadcast},
				"Version": []string{ServiceVersion20201201},
			},
		},
	}
)

type RTC struct {
	*base.Client
}

func NewInstance() *RTC {
	instance := &RTC{}
	instance.Client = base.NewClient(serviceInfo, defaultApiInfoList)
	return instance
}
