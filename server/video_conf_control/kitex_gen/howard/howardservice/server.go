// Code generated by Kitex v0.0.3. DO NOT EDIT.
package howardservice

import (
	"github.com/cloudwego/kitex/server"
	"github.com/volcengine/VolcEngineRTC/server/video_conf_control/kitex_gen/howard"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler howard.HowardService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
