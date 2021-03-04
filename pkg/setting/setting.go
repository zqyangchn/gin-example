package setting

import (
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type setting struct {
	vp  *viper.Viper
	set map[string]interface{}
}

func New() (*setting, error) {
	vp := viper.New()

	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("configs/")
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	return &setting{vp: vp, set: make(map[string]interface{}, 10)}, nil
}

type Server struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type App struct {
	DefaultPageSize      int
	MaxPageSize          int
	UploadSavePath       string
	UploadServerUrl      string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

var AppSetting = &App{}

type JWT struct {
	Secret string
	Issuer string
	Expire time.Duration
}

var JWTSetting = &JWT{}

type Logger struct {
	Level    string
	Stdout   bool
	FilePath string
}

var LoggerSetting = &Logger{}

type Database struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int

	GormForceGormZapLog  bool
	GormLogSlowThreshold time.Duration
}

var DatabaseSetting = &Database{}

type SessionRedis struct {
	Type      string
	Address   string
	Addresses []string
	Password  string
}

var SessionRedisSetting = &SessionRedis{}

func (s *setting) createSection() {
	s.set = map[string]interface{}{
		"Server":       ServerSetting,
		"App":          AppSetting,
		"JWT":          JWTSetting,
		"Log":          LoggerSetting,
		"Database":     DatabaseSetting,
		"SessionRedis": SessionRedisSetting,
	}
}

func (s *setting) readSection(key string, rawVal interface{}) error {
	if err := s.vp.UnmarshalKey(key, rawVal); err != nil {
		return err
	}
	return nil
}

func Setup() error {
	s, err := New()
	if err != nil {
		return err
	}
	s.createSection()

	for key, setting := range s.set {
		if err := s.readSection(key, setting); err != nil {
			return err
		}
		switch key {
		case "Server":
			s := reflect.ValueOf(setting).Elem().Addr().Interface().(*Server)
			s.ReadTimeout *= time.Second
			s.WriteTimeout *= time.Second
		case "JWT":
			j := reflect.ValueOf(setting).Elem().Addr().Interface().(*JWT)
			j.Expire *= time.Second
		}
	}

	return nil
}
