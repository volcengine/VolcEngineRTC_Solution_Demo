package edu_models

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/redis/edu_redis"
)

var UserDsClient UserDs
var RoomDsClient RoomDs

func init() {
	UserDsClient = &MixUserDs{}
	RoomDsClient = &MixRoomDs{}
}

// user datasource model
type UserDs interface {
	//save
	Save(ctx context.Context, user *db.EduUserRoomInfo) error //create or update

	//update users
	UpdateUsersWithMapByParentRoomID(ctx context.Context, parentRoomID string, ups map[string]interface{}) error

	//get user
	GetUserByParentRoomIDUserID(ctx context.Context, parentRoomID, userID string) (*db.EduUserRoomInfo, error)
	GetActiveUserByConnID(ctx context.Context, connID string) (*db.EduUserRoomInfo, error)
	GetReconnectUserByDeviceID(ctx context.Context, deviceID string) (*db.EduUserRoomInfo, error)

	//get users
	GetActiveUsersByUserID(ctx context.Context, userID string) ([]*db.EduUserRoomInfo, error)

	GetActiveStudentsCountByParentRoomID(ctx context.Context, parentRoomID string) (int, error)
	GetActiveStudentsByParentRoomID(ctx context.Context, parentRoomID string, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error)
	GetActiveStudentsByRoomID(ctx context.Context, roomID string, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error)
	GetActiveStudentsByRoomIDRange(ctx context.Context, parentRoomID string, mixIdx, maxIdx int) (map[int][]*db.EduUserRoomInfo, error)
	GetActiveStudentsByRoomIDSet(ctx context.Context, roomIDSet []string) (map[int][]*db.EduUserRoomInfo, error)

	//举手、互动默认active
	GetInteractStudentsCountByParentRoomID(ctx context.Context, parentRoomID string) (int, error)
	GetInteractStudentsByParentRoomID(ctx context.Context, parentRoomID string) ([]*db.EduUserRoomInfo, error)
	GetHandsUpStudentsByParentRoomID(ctx context.Context, parentRoomID string) ([]*db.EduUserRoomInfo, error)
}

//mix implete,use redis and mysql
type MixUserDs struct{}

func (ds *MixUserDs) Save(ctx context.Context, user *db.EduUserRoomInfo) error {
	return edu_db.CreateOrUpdateUser(ctx, user)
}

func (ds *MixUserDs) UpdateUsersWithMapByParentRoomID(ctx context.Context, parentRoomID string, ups map[string]interface{}) error {
	return edu_db.UpdateUsersByParentRoomID(ctx, parentRoomID, ups)
}

func (ds *MixUserDs) GetUserByParentRoomIDUserID(ctx context.Context, parentRoomID, userID string) (*db.EduUserRoomInfo, error) {
	return edu_db.GetUserByParentRoomIDUserID(ctx, parentRoomID, userID)
}

func (ds *MixUserDs) GetActiveUserByConnID(ctx context.Context, connID string) (*db.EduUserRoomInfo, error) {
	return edu_db.GetActiveUserByConnID(ctx, connID)
}

func (ds *MixUserDs) GetReconnectUserByDeviceID(ctx context.Context, deviceID string) (*db.EduUserRoomInfo, error) {
	return edu_db.GetReconnectingUserByDeviceID(ctx, deviceID)
}

func (ds *MixUserDs) GetActiveUsersByUserID(ctx context.Context, userID string) ([]*db.EduUserRoomInfo, error) {
	return edu_db.GetActiveUsersByUserID(ctx, userID)
}

func (ds *MixUserDs) GetActiveStudentsCountByParentRoomID(ctx context.Context, parentRoomID string) (int, error) {
	return edu_db.GetActiveStudentCountByParentRoomID(ctx, parentRoomID)
}

func (ds *MixUserDs) GetActiveStudentsByParentRoomID(ctx context.Context, parentRoomID string, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error) {
	return edu_db.GetActiveStudentsByParentRoomID(ctx, parentRoomID, pageNumber, pageSize)
}

func (ds *MixUserDs) GetActiveStudentsByRoomID(ctx context.Context, roomID string, pageNumber, pageSize int) ([]*db.EduUserRoomInfo, error) {
	return edu_db.GetActiveStudentsByRoomID(ctx, roomID, pageNumber, pageSize)
}

func (ds *MixUserDs) GetActiveStudentsByRoomIDRange(ctx context.Context, parentRoomID string, minIdx, maxIdx int) (map[int][]*db.EduUserRoomInfo, error) {
	res := make(map[int][]*db.EduUserRoomInfo)

	minRoomID := edu_redis.GetGroupRoomIDByIdx(ctx, parentRoomID, minIdx)
	maxRoomID := edu_redis.GetGroupRoomIDByIdx(ctx, parentRoomID, maxIdx)

	students, err := edu_db.GetActiveStudentsByRoomIDRange(ctx, minRoomID, maxRoomID)
	if err != nil {
		return nil, err
	}
	for _, stu := range students {
		idx := edu_redis.GetIdxByGroupRoomID(ctx, stu.RoomID)
		res[idx] = append(res[idx], stu)
	}

	maxNotEmptyPos := minIdx
	for i := maxIdx; i >= minIdx; i-- {
		if _, ok := res[i]; ok {
			maxNotEmptyPos = i
			break
		}
	}

	for i := minIdx; i < maxNotEmptyPos; i++ {
		if _, ok := res[i]; !ok {
			res[i] = make([]*db.EduUserRoomInfo, 0)
		}
	}

	return res, nil
}

func (ds *MixUserDs) GetActiveStudentsByRoomIDSet(ctx context.Context, roomIDSet []string) (map[int][]*db.EduUserRoomInfo, error) {
	res := make(map[int][]*db.EduUserRoomInfo)

	stepSize := 40
	for i := 0; i*stepSize < len(roomIDSet); i++ {
		students, err := edu_db.GetActiveStudentsByRoomIDSet(ctx, roomIDSet[i*stepSize:min(len(roomIDSet), (i+1)*stepSize)])
		if err != nil {
			return nil, err
		}
		for _, stu := range students {
			idx := edu_redis.GetIdxByGroupRoomID(ctx, stu.RoomID)
			res[idx] = append(res[idx], stu)
		}
	}
	return res, nil
}

func (ds *MixUserDs) GetInteractStudentsCountByParentRoomID(ctx context.Context, parentRoomID string) (int, error) {
	return edu_db.GetInteractStudentsCountByParentRoomID(ctx, parentRoomID)
}

func (ds *MixUserDs) GetInteractStudentsByParentRoomID(ctx context.Context, parentRoomID string) ([]*db.EduUserRoomInfo, error) {
	return edu_db.GetInteractStudentsByParentRoomID(ctx, parentRoomID)
}

func (ds *MixUserDs) GetHandsUpStudentsByParentRoomID(ctx context.Context, parentRoomID string) ([]*db.EduUserRoomInfo, error) {
	return edu_db.GetHandsUpStudentsByParentRoomID(ctx, parentRoomID)
}

//room datasource model
type RoomDs interface {
	//save
	Save(ctx context.Context, room *db.EduRoomInfo) error //create or update

	//get room
	GetActiveRoomByRoomID(ctx context.Context, roomID string) (*db.EduRoomInfo, error)

	//get rooms
	GetActiveRooms(ctx context.Context) ([]*db.EduRoomInfo, error)
	GetActiveRoomsByCreateUserID(ctx context.Context, createUserID string) ([]*db.EduRoomInfo, error)
	GetHistoryRoomsByUserID(ctx context.Context, userID string) ([]*db.EduRoomInfo, error)

	//get group room
	GetActiveGroupRoomIDSet(ctx context.Context, roomID string) ([]string, error)

	GetIdleGroupRoomID(ctx context.Context, roomID string, groupNum, groupLimit int) (string, error)
	GetIdxByGroupRoomID(ctx context.Context, groupRoomID string) int
	GetParentRoomID(ctx context.Context, groupRoomID string) string
	//GetGroupRoomIDByIdx(ctx context.Context)

}

type MixRoomDs struct{}

func (ds *MixRoomDs) Save(ctx context.Context, room *db.EduRoomInfo) error {
	return edu_db.CreateOrUpdateRoom(ctx, room)
}

func (ds *MixRoomDs) GetActiveRoomByRoomID(ctx context.Context, roomID string) (*db.EduRoomInfo, error) {
	return edu_db.GetActiveRoomByRoomID(ctx, roomID)
}

func (ds *MixRoomDs) GetActiveRooms(ctx context.Context) ([]*db.EduRoomInfo, error) {
	return edu_db.GetActiveRooms(ctx)
}

func (ds *MixRoomDs) GetActiveRoomsByCreateUserID(ctx context.Context, createUserID string) ([]*db.EduRoomInfo, error) {
	return edu_db.GetActiveRoomsByCreateUserID(ctx, createUserID)
}

func (ds *MixRoomDs) GetHistoryRoomsByUserID(ctx context.Context, userID string) ([]*db.EduRoomInfo, error) {
	return edu_db.GetHistoryRoomsByUserID(ctx, userID)
}

func (ds *MixRoomDs) GetActiveGroupRoomIDSet(ctx context.Context, roomID string) ([]string, error) {
	return edu_db.GetActiveGroupRoomIDSet(ctx, roomID)
}

func (ds *MixRoomDs) GetIdleGroupRoomID(ctx context.Context, roomID string, groupNum, groupLimit int) (string, error) {
	return edu_redis.GetIdleGroupRoomID(ctx, roomID, groupNum, groupLimit)
}

func (ds *MixRoomDs) GetIdxByGroupRoomID(ctx context.Context, groupRoomID string) int {
	return edu_redis.GetIdxByGroupRoomID(ctx, groupRoomID)
}

func (ds *MixRoomDs) GetParentRoomID(ctx context.Context, groupRoomID string) string {
	return edu_redis.GetParentRoomID(ctx, groupRoomID)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
