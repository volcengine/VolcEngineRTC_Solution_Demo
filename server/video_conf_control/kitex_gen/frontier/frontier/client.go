// Code generated by Kitex v0.0.3. DO NOT EDIT.

package frontier

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/frontier"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	PushToClient(ctx context.Context, param *frontier.TPushParam, callOptions ...callopt.Option) (r *frontier.TPushResp, err error)
	BroadcastToClient(ctx context.Context, param *frontier.TBroadCastParam, callOptions ...callopt.Option) (r *frontier.TBroadCastResp, err error)
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
	return &kFrontierClient{
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

type kFrontierClient struct {
	*kClient
}

func (p *kFrontierClient) PushToClient(ctx context.Context, param *frontier.TPushParam, callOptions ...callopt.Option) (r *frontier.TPushResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PushToClient(ctx, param)
}

func (p *kFrontierClient) BroadcastToClient(ctx context.Context, param *frontier.TBroadCastParam, callOptions ...callopt.Option) (r *frontier.TBroadCastResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.BroadcastToClient(ctx, param)
}
