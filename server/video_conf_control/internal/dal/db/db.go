package db

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	INACTIVE = iota
	ACTIVE
	RECONNECTING
)

type UserStatusType int

const (
	UserStatusAudience UserStatusType = iota
	UserStatusRaiseHands
	UserStatusOnMicrophone
)

var Client *gorm.DB

func Open(dsn string) {
	var err error
	Client, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
}

type MeetingConnection struct {
	ID        int64     `gorm:"column:id" json:"id"`
	AppID     string    `gorm:"column:app_id" json:"app_id"`
	ConnID    string    `gorm:"column:conn_id" json:"conn_id"`
	Addr      string    `gorm:"column:addr" json:"addr"`
	Addr6     string    `gorm:"column:addr6" json:"addr6"`
	State     int       `gorm:"column:state" json:"state"`
	DeviceID  string    `gorm:"column:device_id" json:"device_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type MeetingUser struct {
	ID         int64     `gorm:"column:id" json:"id"`
	AppID      string    `gorm:"column:app_id" json:"app_id"`
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	RoomID     string    `gorm:"column:room_id" json:"room_id"`
	ConnID     string    `gorm:"column:conn_id" json:"conn_id"`
	State      int       `gorm:"column:state" json:"state"`
	IsHost     bool      `gorm:"column:is_host" json:"is_host"`
	IsMicOn    bool      `gorm:"column:is_mic_on" json:"is_mic_on"`
	IsCameraOn bool      `gorm:"column:is_camera_on" json:"is_camera_on"`
	IsSharing  bool      `gorm:"column:is_sharing" json:"is_sharing"`
	DeviceID   string    `gorm:"column:device_id" json:"device_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type MeetingRoom struct {
	ID        int64     `gorm:"column:id" json:"id"`
	AppID     string    `gorm:"column:app_id" json:"app_id"`
	RoomID    string    `gorm:"column:room_id" json:"room_id"`
	State     bool      `gorm:"column:state" json:"state"`
	Record    bool      `gorm:"column:record" json:"record"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type MeetingVideoRecord struct {
	ID        int64     `gorm:"column:id" json:"id"`
	AppID     string    `gorm:"column:app_id" json:"app_id"`
	RoomID    string    `gorm:"column:room_id" json:"room_id"`
	VID       string    `gorm:"column:vid" json:"vid"`
	State     int       `gorm:"column:state" json:"state"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CsMeetingUser struct {
	ID           int64     `gorm:"column:id" json:"id"`
	AppID        string    `gorm:"column:app_id" json:"app_id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	UserName     string    `gorm:"column:user_name" json:"user_name"`
	UserStatus   int       `gorm:"column:user_status" json:"user_status"`
	RoomID       string    `gorm:"column:room_id" json:"room_id"`
	ConnID       string    `gorm:"column:conn_id" json:"conn_id"`
	State        int       `gorm:"column:state" json:"state"`
	IsHost       bool      `gorm:"column:is_host" json:"is_host"`
	IsMicOn      bool      `gorm:"column:is_mic_on" json:"is_mic_on"`
	DeviceID     string    `gorm:"column:device_id" json:"device_id"`
	RaiseHandsAt time.Time `gorm:"column:raise_hands_at" json:"raise_hands_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CsMeetingRoom struct {
	ID        int64     `gorm:"column:id" json:"id"`
	AppID     string    `gorm:"column:app_id" json:"app_id"`
	RoomID    string    `gorm:"column:room_id" json:"room_id"`
	RoomName  string    `gorm:"column:room_name" json:"room_name"`
	State     bool      `gorm:"column:state" json:"state"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type UserProfile struct {
	ID        int64     `gorm:"column:id" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	UserName  string    `gorm:"column:user_name" json:"user_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type EduTime time.Time

const TimeFormat = "2006-01-02 15:04:05 +0800 CST"

//for json
func (t EduTime) MarshalJSON() ([]byte, error) {
	stapm := fmt.Sprintf("%d", time.Time(t).UnixNano())
	return []byte(stapm), nil
}

func (t *EduTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = EduTime(time.Time{})
		return
	}
	// 指定解析的格式
	ns, _ := strconv.ParseInt(string(data), 10, 64)
	now := time.Unix(0, ns)
	*t = EduTime(now)
	return
}

//for mysql
func (t EduTime) Value() (driver.Value, error) {
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *EduTime) Scan(v interface{}) error {
	tTime, _ := time.ParseInLocation("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String(), time.Local)
	*t = EduTime(tTime)
	return nil
}

type EduRoomInfo struct {
	ID                 int64   `gorm:"column:id" json:"id"`
	AppID              string  `gorm:"column:app_id" json:"app_id"`
	RoomID             string  `gorm:"column:room_id" json:"room_id"`
	RoomName           string  `gorm:"column:room_name" json:"room_name"`
	RoomType           int     `gorm:"column:room_type" json:"room_type"`
	Status             int     `gorm:"column:status" json:"status"`
	CreateUserID       string  `gorm:"column:create_user_id" json:"create_user_id"`
	CreatedTime        EduTime `gorm:"column:create_time" gorm:"autoCreateTime"  json:"create_time"`
	UpdatedTime        EduTime `gorm:"column:update_time" gorm:"autoUpdateTime"  json:"update_time"`
	BeginClassTime     int64   `gorm:"column:begin_class_time" json:"begin_class_time"`
	EndClassTime       int64   `gorm:"column:end_class_time" json:"end_class_time"`
	AudioMuteAll       bool    `gorm:"column:audio_mute_all" json:"audio_mute_all"`
	VideoMuteAll       bool    `gorm:"column:video_mute_all" json:"video_mute_all"`
	EnableGroupSpeech  bool    `gorm:"column:enable_group_speech" json:"enable_group_speech"`
	EnableInteractive  bool    `gorm:"column:enable_interactive" json:"enable_interactive"`
	IsRecording        bool    `gorm:"column:is_recording" json:"is_recording"`
	TeacherName        string  `gorm:"column:teacher_name" json:"teacher_name"`
	BeginClassTimeReal int64   `gorm:"column:begin_class_time_real" json:"begin_class_time_real"`
	EndClassTimeReal   int64   `gorm:"column:end_class_time_real" json:"end_class_time_real"`
	Token              string  `gorm:"column:token" json:"token"`
	GroupNum           int     `gorm:"column:group_num" json "group_num"`
	GroupLimit         int     `gorm:"column:group_limit" json:"group_limit"`
	//RoomChildInfo      []*EduRoomChildInfo `gorm:"-" sql:"-" json:"room_child_info,omitempty"`
}

type EduRoomChildInfo struct {
	ID           int64  `gorm:"column:id" json:"id"`
	ParentRoomID string `gorm:"column:parent_room_id" json:"parent_room_id"`
	AppID        string `gorm:"column:app_id" json:"app_id"`
	RoomID       string `gorm:"column:room_id" json:"room_id"`
	RoomName     string `gorm:"column:room_name" json:"room_name"`
	RoomIdx      int    `gorm:"column:room_idx" json:"room_idx"`
}

type EduUserRoomInfo struct {
	ID                 int64   `gorm:"column:id" json:"id"`
	AppID              string  `gorm:"column:app_id" json:"app_id"`
	RoomID             string  `gorm:"column:room_id" json:"room_id"`
	DeviceID           string  `gorm:"device_id" json:"device_id"`
	UserID             string  `gorm:"column:user_id" json:"user_id"`
	UserName           string  `gorm:"column:user_name" json:"user_name"`
	UserRole           int     `gorm:"column:user_role" json:"user_role"`
	UserStatus         int     `gorm:"column:user_status" json:"user_status"`
	CreatedTime        EduTime `gorm:"column:create_time" gorm:"autoCreateTime"  json:"create_time"`
	UpdatedTime        EduTime `gorm:"column:update_time" gorm:"autoUpdateTime"  json:"update_time"`
	JoinTime           int64   `gorm:"column:join_time" json:"join_time"`
	LeaveTime          int64   `gorm:"column:leave_time" json:"leave_time"`
	IsMicOn            bool    `gorm:"column:is_mic_on" json:"is_mic_on"`
	IsCameraOn         bool    `gorm:"column:is_camera_on" json:"is_camera_on"`
	IsHandsUp          bool    `gorm:"column:is_hands_up" json:"is_hands_up"`
	IsInteract         bool    `gorm:"column:is_interact" json:"is_interact"`
	GroupSpeechJoinRtc bool    `gorm:"column:group_speech_join_rtc" json:"group_speech_join_rtc"`
	RtcToken           string  `gorm:"column:rtc_token" json:"rtc_token"`
	ConnID             string  `gorm:"column:conn_id" json:"conn_id"`
	ParentRoomID       string  `gorm:"column:parent_room_id" json:"parent_room_id"`
}

type EduRecordInfo struct {
	ID              int64   `gorm:"column:id" json:"id"`
	AppID           string  `gorm:"column:app_id" json:"app_id"`
	RoomID          string  `gorm:"column:room_id" json:"room_id"`
	ParentRoomID    string  `gorm:"column:parent_room_id" json:"parent_room_id"`
	UserID          string  `gorm:"column:user_id" json:"user_id"`
	RoomName        string  `gorm:"column:room_name" json:"room_name"`
	RecordStatus    int     `gorm:"column:record_status" json:"record_status"`
	CreatedTime     EduTime `gorm:"column:create_time" gorm:"autoCreateTime" json:"create_time"`
	UpdatedTime     EduTime `gorm:"column:update_time" gorm:"autoUpdateTime" json:"update_time"`
	RecordBeginTime int64   `gorm:"column:record_begin_time" json:"record_begin_time"`
	RecordEndTime   int64   `gorm:"column:record_end_time" json:"record_end_time"`
	TaskID          string  `gorm:"column:task_id" json:"task_id"`
	Vid             string  `gorm:"column:vid" json:"vid"`
	VideoURL        string  `gorm:"-" sql:"-" json:"video_url"`
}
