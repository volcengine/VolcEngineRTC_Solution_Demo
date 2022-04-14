package cs_handler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"
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
			hall := cs_service.GetHall()
			rooms, err := hall.ListRooms(ctx)
			if err != nil {
				logs.CtxError(ctx, "cron get rooms failed,error:%s", err)
				return
			}
			for _, r := range rooms {
				if time.Now().Sub(r.GetCreateTime()) >= time.Duration(config.Configs().CsExperienceTime)*time.Minute {
					roomService := cs_service.NewRoomService(ctx, r)
					err = roomService.Finish(ctx)
					if err != nil {
						logs.CtxError(ctx, "cron room finish failed,error:%s", err)
						continue
					}
					informer := cs_service.GetInformService()
					informer.BroadcastRoom(ctx, r.GetRoomID(), cs_models.OnCsMeetingEnd, map[string]string{
						"room_id": r.GetRoomID(),
					})
				}
			}

		})
	}
	return handler
}
