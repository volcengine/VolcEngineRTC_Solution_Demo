package svc_handler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/task"
)

var handler *EventHandler

type EventHandler struct {
	c *cron.Cron
}

func NewEventHandler() *EventHandler {
	if handler == nil {
		handler = &EventHandler{
			c: task.GetCronTask(),
		}
		handler.c.AddFunc("@every 1m", func() {
			ctx := context.Background()
			roomFactory := svc_service.GetRoomFactory()
			rooms, err := roomFactory.GetActiveRoomList(ctx, false)
			if err != nil {
				logs.CtxError(ctx, "cron: get svc rooms failed,error:%s", err)
				return
			}

			roomService := svc_service.GetRoomService()
			for _, room := range rooms {
				if time.Now().Sub(room.GetCreateTime()) >= time.Duration(config.Configs().SvcExperienceTime)*time.Minute {
					err = roomService.FinishLive(ctx, room.GetRoomID(), svc_service.FinishTypeTimeout)
					if err != nil {
						logs.CtxError(ctx, "cron: finish room failed,error:%s", err)
						continue
					}
				}
			}

		})
	}
	return handler
}
