// Code generated by Kitex v0.0.3. DO NOT EDIT.

package howardservice

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/howard"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	CleanVideo(ctx context.Context, request *howard.CleanVideoRequest, callOptions ...callopt.Option) (r *howard.CleanVideoResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kHowardServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kHowardServiceClient struct {
	*kClient
}

func (p *kHowardServiceClient) CleanVideo(ctx context.Context, request *howard.CleanVideoRequest, callOptions ...callopt.Option) (r *howard.CleanVideoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CleanVideo(ctx, request)
}
