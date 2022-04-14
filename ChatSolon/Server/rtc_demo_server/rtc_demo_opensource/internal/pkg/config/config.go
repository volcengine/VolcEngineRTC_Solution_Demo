package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/configor"
)

const (
	DefaultConfDir  = "conf"
	DefaultConfFile = "config.yaml"
)

type Configuration struct {
	//general
	MysqlDSN      string `yaml:"mysql_dsn"`
	RedisAddr     string `yaml:"redis_addr"`
	RedisPassword string `yaml:"redis_password"`
	ServerUrl     string `yaml:"server_url"`
	Addr          string `yaml:"addr"`
	Port          string `yaml:"port"`

	VolcAk        string `yaml:"volc_ak"`
	VolcSk        string `yaml:"volc_sk"`
	VolcEngineUrl string `yaml:"volc_engine_url"`

	CheckLoginToken bool `yaml:"check_login_token"`

	//login
	AuditorPhoneCode map[string]string `yaml:"auditor_phone_code"`

	//chat solon
	CsAppID          string `yaml:"cs_app_id"`
	CsAppKey         string `yaml:"cs_app_key"`
	CsExperienceTime int    `yaml:"cs_experience_time"`

	//unknown
	RoomUserLimit    int `yaml:"room_user_limit"`
	ReconnectTimeout int `yaml:"reconnect_timeout"`
}

var configs *Configuration

func InitConfig(file string) {
	configs = &Configuration{}
	if err := configor.Load(configs, file); err != nil {
		fmt.Fprintf(os.Stderr, "ParseConfig failed: err=%v\n", err)
		os.Exit(1)
	}
}

func Configs() *Configuration {
	if configs == nil {
		panic("config not init")
	}
	return configs
}
