package task

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/record"
)

func Start() {
	ctx := context.Background()
	c := cron.New()
	addEduTask(ctx, c)
	c.Start()
}

func addEduTask(ctx context.Context, c *cron.Cron) {
	c.AddFunc("@every 1m", func() {
		activeRooms, err := edu_db.GetActiveRooms(ctx)
		if err != nil {
			logs.Errorf("cron get rooms failed,error:%s", err)
			return
		}
		for _, r := range activeRooms {
			logs.Infof("cron get active room::%v", r.RoomID)
			room, err := edu_models.GetRoom(ctx, r.RoomID)
			if err != nil {
				logs.Errorf("get room failed,error:%s", err)
				continue
			}

			if room.GetBeginClassRealTime() == 0 {
				//创建一定时间没开课,删除
				if time.Now().Sub(room.GetCreateTime()) >= time.Minute*time.Duration(config.Config.EduCreatedExpireTime) {
					logs.Infof("cron delete room:%s", room.GetRoomID())
					room.InformRoom(ctx, edu_models.OnEndClass, &edu_models.NoticeRoom{RoomID: room.GetRoomID()})
					edu_models.StopRecord(ctx, config.Config.AppID, room.GetRoomID(), room.GetTeacherUserID())
					err = room.DeleteRoom(ctx)
					if err != nil {
						logs.Errorf("cron delete room failed,error:%s", err)
						continue
					}
					err = room.Save(ctx)
					if err != nil {
						logs.Errorf("cron delete room failed,error:%s", err)
						continue
					}
				}
			} else {
				//开课30分钟后，且老师离线的，调用end_class
				startTime := time.Unix(0, room.GetBeginClassRealTime())
				teacher := room.GetTeacherInfo(ctx)
				if time.Now().Sub(startTime) > time.Minute*time.Duration(config.Config.EduCreatedExpireTime) && teacher.UserStatus == edu_db.UserStatusOffline {
					logs.Infof("auto end class,room:%s", room.GetRoomID())
					room.InformRoom(ctx, edu_models.OnEndClass, &edu_models.NoticeRoom{RoomID: room.GetRoomID()})
					edu_models.StopRecord(ctx, config.Config.AppID, room.GetRoomID(), room.GetTeacherUserID())
					err = room.EndClass(ctx)
					if err != nil {
						logs.Errorf("cron end class failed,error:%s", err)
						continue
					}
					err = room.Save(ctx)
					if err != nil {
						logs.Errorf("cron save failed,error:%s", err)
						continue
					}
				}
			}
		}
	})

	c.AddFunc("@every 1m", func() {
		records, err := edu_db.QueryAllTimeOutRecord(ctx)
		if err != nil {
			logs.Errorf("cron get records failed,error:%s", err)
			return
		}
		for _, r := range records {
			if err := record.StopRecord(ctx, config.Config.AppID, r.RoomID, r.TaskID); err != nil {
				logs.Errorf("stop record failed,err:%s", err)
			}
		}
	})
}
