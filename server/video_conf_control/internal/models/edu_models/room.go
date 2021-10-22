package edu_models

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/conn_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/response"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/rpc/frontier"

	logs "github.com/sirupsen/logrus"
)

const (
	retryDelay    = 8 * time.Millisecond
	maxRetryDelay = 128 * time.Millisecond
	maxRetryNum   = 10
)

type Room struct {
	r       *db.EduRoomInfo
	isDirty bool
}

func GetCreatedRooms(ctx context.Context, userID string) ([]*db.EduRoomInfo, error) {
	return RoomDsClient.GetActiveRoomsByCreateUserID(ctx, userID)
}

func GetParentRoomID(ctx context.Context, groupRoomID string) string {
	return RoomDsClient.GetParentRoomID(ctx, groupRoomID)
}

func NewRoom(dbRoom *db.EduRoomInfo) *Room {
	room := &Room{
		r:       dbRoom,
		isDirty: true,
	}

	return room
}

func GetRoom(ctx context.Context, roomID string) (*Room, error) {
	dbRoom, err := RoomDsClient.GetActiveRoomByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if dbRoom == nil {
		return nil, errors.New("room not exist")
	}
	room := &Room{
		r:       dbRoom,
		isDirty: false,
	}
	return room, nil
}

func (r *Room) Save(ctx context.Context) error {
	if r.isDirty {
		r.r.UpdatedTime = db.EduTime(time.Now())
		err := RoomDsClient.Save(ctx, r.r)
		if err != nil {
			return err
		}
	}
	r.isDirty = false
	return nil
}

func (r *Room) BeginClass(ctx context.Context) {
	r.r.Status = edu_db.ClassRunning
	r.r.BeginClassTimeReal = time.Now().UnixNano()
	r.r.IsRecording = true

	r.isDirty = true
}

func (r *Room) EndClass(ctx context.Context) error {
	r.r.Status = edu_db.ClassFinish
	r.r.EndClassTimeReal = time.Now().UnixNano()
	r.isDirty = true
	err := UserDsClient.UpdateUsersWithMapByParentRoomID(ctx, r.r.RoomID, map[string]interface{}{
		"user_status": edu_db.UserStatusOffline,
		"leave_time":  time.Now().UnixNano(),
		"is_hands_up": 0,
		"is_interact": 0,
	})
	if err != nil {
		logs.Errorf("update users failed,error:%s", err)
		return err
	}
	return nil
}

func (r *Room) ListStudents(ctx context.Context, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error) {
	empty := make([]*db.EduUserRoomInfo, 0)
	users, err := UserDsClient.GetActiveStudentsByParentRoomID(ctx, r.r.RoomID, pageNumber, pageSize)
	if err != nil {
		logs.Errorf("get active students failed,error:%s", err)
		return empty, err
	}
	return users, nil
}

func (r *Room) ListStudentsByGroup(ctx context.Context, pageNumber, pageSize int) (map[int][]*db.EduUserRoomInfo, error) {
	minIdx := (pageNumber - 1) * pageSize
	maxIdx := pageNumber*pageSize - 1
	groupUserMap, err := UserDsClient.GetActiveStudentsByRoomIDRange(ctx, r.r.RoomID, minIdx, maxIdx)
	if err != nil {
		logs.Errorf("get students by range failed,error:%s", err)
		return nil, err
	}
	return groupUserMap, nil
}

func (r *Room) GetRoomStudentCount(ctx context.Context) (int, error) {
	return UserDsClient.GetActiveStudentsCountByParentRoomID(ctx, r.r.RoomID)
}

func (r *Room) ListGroupRoomStudents(ctx context.Context, groupRoomID string) ([]*db.EduUserRoomInfo, error) {
	empty := make([]*db.EduUserRoomInfo, 0)
	users, err := UserDsClient.GetActiveStudentsByRoomID(ctx, groupRoomID, 0, 0)
	if err != nil {
		logs.Errorf("get group room users id failed,error:%s", err)
		return empty, nil
	}
	return users, nil
}

func (r *Room) ListHandsUpStudents(ctx context.Context) ([]*db.EduUserRoomInfo, error) {
	empty := make([]*db.EduUserRoomInfo, 0)
	users, err := UserDsClient.GetHandsUpStudentsByParentRoomID(ctx, r.r.RoomID)
	if err != nil {
		logs.Errorf("get hands up users failed,error:%s", err)
		return empty, nil
	}
	return users, nil
}

func (r *Room) ListInteractStudents(ctx context.Context) ([]*db.EduUserRoomInfo, error) {
	empty := make([]*db.EduUserRoomInfo, 0)
	users, err := UserDsClient.GetInteractStudentsByParentRoomID(ctx, r.r.RoomID)
	if err != nil {
		logs.Errorf("get interact users failed,error:%s", err)
		return empty, nil
	}
	return users, nil
}

func (r *Room) GetInteractUserCount(ctx context.Context) int {
	count, err := UserDsClient.GetInteractStudentsCountByParentRoomID(ctx, r.r.RoomID)
	if err != nil {
		return 0
	}
	return count
}

func (r *Room) ListGroupRoomID(ctx context.Context) ([]string, error) {
	return RoomDsClient.GetActiveGroupRoomIDSet(ctx, r.r.RoomID)
}

func (r *Room) GetIdleGroupRoomID(ctx context.Context) (string, error) {
	return RoomDsClient.GetIdleGroupRoomID(ctx, r.r.RoomID, r.r.GroupNum, r.r.GroupLimit)
}

func (r *Room) GetIdxByGroupRoomID(ctx context.Context, groupRoomID string) int {
	return RoomDsClient.GetIdxByGroupRoomID(ctx, groupRoomID)
}

func (r *Room) GetBeginClassRealTime() int64 {
	return r.r.BeginClassTimeReal
}

func (r *Room) GetCreateTime() time.Time {
	return time.Time(r.r.CreatedTime)
}

func (r *Room) OpenGroupSpeech() {
	r.r.EnableGroupSpeech = true
	r.isDirty = true
}

func (r *Room) CloseGroupSpeech() {
	r.r.EnableGroupSpeech = false
	r.isDirty = true
}

func (r *Room) IsEnableGroupSpeech() bool {
	return r.r.EnableGroupSpeech
}

func (r *Room) OpenInteract() {
	r.r.EnableInteractive = true
	r.isDirty = true
}

func (r *Room) CloseInteract(ctx context.Context) error {
	r.r.EnableInteractive = false
	r.isDirty = true
	err := UserDsClient.UpdateUsersWithMapByParentRoomID(ctx, r.r.RoomID, map[string]interface{}{
		"is_hands_up": 0,
		"is_interact": 0,
	})
	return err
}

func (r *Room) IsEnableInteract() bool {
	return r.r.EnableInteractive
}

func (r *Room) GetRoomInfo() *db.EduRoomInfo {
	return r.r
}

func (r *Room) GetTeacherInfo(ctx context.Context) *db.EduUserRoomInfo {
	teacherUserID := r.r.CreateUserID
	teacher, err := UserDsClient.GetUserByParentRoomIDUserID(ctx, r.r.RoomID, teacherUserID)
	if err != nil || teacher == nil {
		teacher = &db.EduUserRoomInfo{
			AppID:        config.Config.AppID,
			ParentRoomID: r.r.RoomID,
			RoomID:       r.r.RoomID,
			UserID:       r.r.CreateUserID,
			UserName:     r.r.TeacherName,
			UserStatus:   edu_db.UserStatusOffline,
		}
	}
	return teacher
}

func (r *Room) GetRoomID() string {
	return r.r.RoomID
}

func (r *Room) GetRoomName() string {
	return r.r.RoomName
}

func (r *Room) GetTeacherUserID() string {
	return r.r.CreateUserID
}

func (r *Room) GetRoomType() int {
	return r.r.RoomType
}

func (r *Room) IsGroupRoom() bool {
	return r.r.RoomType == edu_db.RoomTypeGroupRoomClass
}

func (r *Room) DeleteRoom(ctx context.Context) error {
	err := r.EndClass(ctx)
	if err != nil {
		logs.Errorf("end class failed,error:%s", err)
		return err
	}
	r.r.Status = edu_db.ClassDelete
	r.isDirty = true
	return nil
}

func (r *Room) Lock() {

}

func (r *Room) Unlock() {

}

func (r *Room) InformRoom(ctx context.Context, event InformEvent, data interface{}) {
	users, err := r.ListStudents(ctx, 0, 0)
	if err != nil {
		logs.Errorf("get room users failed,error:%s", err)
		return
	}
	//添加老师
	users = append(users, r.GetTeacherInfo(ctx))
	logs.Infof("event:%s,get users count:%v", event, len(users))

	connIDs := make([]string, 0)
	for _, user := range users {
		if user.ConnID != "" {
			connIDs = append(connIDs, user.ConnID)
		}
	}
	conns, err := conn_db.GetConnections(ctx, connIDs)
	if err != nil {
		logs.Infof("failed to get conns, connIDs: %v, err: %v", connIDs, err)
		return
	}

	connIDMap := make(map[string][]string) // addr-addr6: []connIDs
	for _, conn := range conns {
		address := strings.Join([]string{conn.Addr, conn.Addr6}, "-")
		connIDMap[address] = append(connIDMap[address], conn.ConnID)
	}
	logs.Infof("get connID count:%v", len(connIDMap))

	for address, connIDList := range connIDMap {
		addr := strings.Split(address, "-")
		go frontier.BroadcastToClient(ctx, connIDList, addr[0], addr[1], string(event), response.NewInformToClient(string(event), data), -1)
	}
}

func (r *Room) InformGroupRoom(ctx context.Context, groupRoomID string, event InformEvent, data interface{}) {
	users, err := r.ListGroupRoomStudents(ctx, groupRoomID)
	if err != nil {
		logs.Infof("failed to get users in group room, error: %s", err)
		return
	}
	logs.Infof("event:%s,get users count:%v", event, len(users))

	connIDs := make([]string, 0)
	for _, user := range users {
		if user.ConnID != "" {
			connIDs = append(connIDs, user.ConnID)
		}
	}

	conns, err := conn_db.GetConnections(ctx, connIDs)
	if err != nil {
		logs.Infof("failed to get conns, connIDs: %v, err: %v", connIDs, err)
		return
	}

	connIDMap := make(map[string][]string) // addr-addr6: []connIDs
	for _, conn := range conns {
		address := strings.Join([]string{conn.Addr, conn.Addr6}, "-")
		connIDMap[address] = append(connIDMap[address], conn.ConnID)
	}
	logs.Infof("get connID count:%v", len(connIDMap))

	for address, connIDList := range connIDMap {
		addr := strings.Split(address, "-")
		go frontier.BroadcastToClient(ctx, connIDList, addr[0], addr[1], string(event), response.NewInformToClient(string(event), data), -1)
	}

}

func ApplyRoomIDWithRetry(ctx context.Context) (string, error) {
	roomID, err := cs_models.GenerateRoomID(ctx)
	for i := 0; roomID == 0 && i <= maxRetryNum; i++ {
		backOff := time.Duration(int64(math.Pow(2, float64(i)))) * retryDelay
		if backOff > maxRetryDelay {
			backOff = maxRetryDelay
		}
		time.Sleep(backOff)
		roomID, err = cs_models.GenerateRoomID(ctx)
	}
	if roomID == 0 {
		logs.Errorf("failed to generate roomID, err: %s", err)
		return "", custom_error.InternalError(errors.New("make room err"))
	}
	return strconv.FormatInt(roomID, 10), nil
}

func GetActiveRooms(ctx context.Context) ([]*db.EduRoomInfo, error) {
	return RoomDsClient.GetActiveRooms(ctx)
}

func GetHistoryRoomsByUserID(ctx context.Context, userID string) ([]*db.EduRoomInfo, error) {
	return RoomDsClient.GetHistoryRoomsByUserID(ctx, userID)
}
