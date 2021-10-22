package edu_service

import (
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
)

var eventHandler = map[string]service_utils.EventHandlerFunc{
	"eduReconnect":                   reconnect,
	"eduCreateClass":                 eduCreateClass,
	"eduGetCreatedClass":             getCreatedClass,
	"eduTeacherJoinClass":            teacherJoinClass,
	"eduTeacherGetStudentsInfo":      teacherGetStudentsInfo,
	"eduTeacherGetGroupStudentsInfo": teacherGetGroupStudents,
	"eduBeginClass":                  beginClass,
	"eduEndClass":                    endClass,
	"eduDeleteRecord":                deleteRecord,
	"eduGetGroupClassInfo":           getGroupClassInfo,
	"eduOpenGroupSpeech":             openGroupSpeech,
	"eduCloseGroupSpeech":            closeGroupSpeech,
	"eduOpenVideoInteract":           openVideoInteract,
	"eduCloseVideoInteract":          closeVideoInteract,
	"eduGetHandsUpList":              getHandsUpList,
	"eduGetStuMicOnList":             teacherGetStuMicOnList,
	"eduGetActiveClass":              getActiveClass,
	"eduApproveMic":                  onApproveInteract,
	"eduForceMicOff":                 onFinishInteract,
	"eduGetHistoryRecordList":        getHistoryRecordList,
	"eduGetHistoryRoomList":          getHistoryRoomList,
	"eduCancelHandsUp":               cancelHandsUp,
	"eduHandsUp":                     handsUp,
	"eduJoinClass":                   joinClass,
	"eduLeaveClass":                  leaveClass,
	"eduTurnOffCamera":               turnOffCamera,
	"eduTurnOnCamera":                turnOnCamera,
	"eduTurnOffMic":                  turnOffMic,
	"eduTurnOnMic":                   turnOnMic,
}

func GetHandlerByEventName(event string) (service_utils.EventHandlerFunc, bool) {
	h, ok := eventHandler[event]
	return h, ok
}
