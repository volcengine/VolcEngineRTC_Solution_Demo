package edu_models

type InformEvent string

const (
	OnTeacherJoinClass      InformEvent = "onTeacherJoinClass"
	OnTeacherLeaveClass     InformEvent = "onTeacherLeaveClass"
	OnBeginClass            InformEvent = "onBeginClass"
	OnEndClass              InformEvent = "onEndClass"
	OnOpenGroupSpeech       InformEvent = "onOpenGroupSpeech"
	OnCloseGroupSpeech      InformEvent = "onCloseGroupSpeech"
	OnOpenVideoInteract     InformEvent = "onOpenVideoInteract"
	OnCloseVideoInteract    InformEvent = "onCloseVideoInteract"
	OnStartInteract         InformEvent = "onStuMicOn"
	OnFinishInteract        InformEvent = "onStuMicOff"
	OnTeacherMicOn          InformEvent = "onTeacherMicOn"
	OnTeacherMicOff         InformEvent = "onTeacherMicOff"
	OnTeacherCameraOn       InformEvent = "onTeacherCameraOn"
	OnTeacherCameraOff      InformEvent = "onTeacherCameraOff"
	OnStudentJoinGroupRoom  InformEvent = "onStudentJoinGroupRoom"
	OnStudentLeaveGroupRoom InformEvent = "onStudentLeaveGroupRoom"
	OnLogInElsewhere        InformEvent = "onLogInElsewhere"
)

type NoticeRoom struct {
	RoomID   string `json:"room_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}
