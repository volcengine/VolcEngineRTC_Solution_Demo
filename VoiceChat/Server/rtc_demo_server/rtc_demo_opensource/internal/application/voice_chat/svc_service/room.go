package svc_service

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/token"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type Room struct {
	*svc_entity.SvcRoom
	AudienceCount int `json:"audience_count"`
	isDirty       bool
}

func (r *Room) IsDirty() bool {
	return r.isDirty
}

func (r *Room) SetIsDirty(isDirty bool) {
	r.isDirty = isDirty
}

func (r *Room) SetUpdateTime(time time.Time) {
	r.UpdateTime = time
	r.isDirty = true
}

func (r *Room) GetDbRoom() *svc_entity.SvcRoom {
	return r.SvcRoom
}

func (r *Room) GetCreateTime() time.Time {
	return r.CreateTime
}

func (r *Room) Start() {
	r.Status = svc_db.RoomStatusStart
	r.isDirty = true
}

func (r *Room) Finish() {
	r.Status = svc_db.RoomStatusFinish
	r.FinishTime = time.Now().UnixNano() / 1e6
	r.isDirty = true
}

func (r *Room) OpenApply() {
	r.EnableAudienceInteractApply = 1
	r.isDirty = true
}

func (r *Room) CloseApply() {
	r.EnableAudienceInteractApply = 0
	r.isDirty = true
}

func (r *Room) IsNeedApply() bool {
	return r.EnableAudienceInteractApply == 1
}

func (r *Room) GetRoomID() string {
	return r.RoomID
}

func (r *Room) GetAppID() string {
	return r.AppID
}

func (r *Room) GetHostUserID() string {
	return r.HostUserID
}

func (r *Room) GenerateToken(ctx context.Context, userID string) string {
	param := &token.GenerateParam{
		AppID:        config.Configs().SvcAppID,
		AppKey:       config.Configs().SvcAppKey,
		RoomID:       r.RoomID,
		UserID:       userID,
		ExpireAt:     7 * 24 * 3600,
		CanPublish:   true,
		CanSubscribe: true,
	}
	rtcToken, err := token.GenerateToken(param)
	if err != nil {
		logs.CtxError(ctx, "generate token failed,error:%s", err)
	}
	return rtcToken
}
