package cs_models

type InformEvent string

const (
	OnCsJoinMeeting   InformEvent = "onCsJoinMeeting"
	OnCsLeaveMeeting  InformEvent = "onCsLeaveMeeting"
	OnCsRaiseHandsMic InformEvent = "onCsRaiseHandsMic"
	OnCsInviteMic     InformEvent = "onCsInviteMic"
	OnCsMicOn         InformEvent = "onCsMicOn"
	OnCsMicOff        InformEvent = "onCsMicOff"
	OnCsMuteMic       InformEvent = "onCsMuteMic"
	OnCsUnmuteMic     InformEvent = "onCsUnmuteMic"
	OnCsMeetingEnd    InformEvent = "onCsMeetingEnd"
	OnCsHostChange    InformEvent = "onCsHostChange"
)
