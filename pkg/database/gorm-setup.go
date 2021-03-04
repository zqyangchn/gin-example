package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gin-example/pkg/gorm-zap"
	"gin-example/pkg/logging"
	"gin-example/pkg/setting"
)

var gormDB *gorm.DB

func Setup() error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		setting.DatabaseSetting.UserName,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.DBName,
		setting.DatabaseSetting.Charset,
		setting.DatabaseSetting.ParseTime,
	)
	mysqlConfig := mysql.Config{
		DSN: dsn,
	}

	gormConfig := gorm.Config{
		//SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.DatabaseSetting.TablePrefix,
			SingularTable: true,
		},
		Logger: gormLogger(),
	}
	gormDB, err = gorm.Open(mysql.New(mysqlConfig), &gormConfig)
	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(setting.DatabaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(setting.DatabaseSetting.MaxOpenConns)

	return nil
}

func GetGormDB() *gorm.DB {
	g := gormDB
	return g
}

func gormLogLevel() gormzap.LogLevel {
	switch setting.LoggerSetting.Level {
	case zap.DebugLevel.String():
		return gormzap.Info
	case zap.ErrorLevel.String():
		return gormzap.Error
	case zap.WarnLevel.String():
		return gormzap.Warn
	case zap.InfoLevel.String():
		return gormzap.Warn
	default:
		return gormzap.Silent
	}
}

func gormLogger() logger.Interface {
	if setting.DatabaseSetting.GormForceGormZapLog {
		return gormzap.New(logging.GormLogger, gormzap.Config{
			SlowThreshold: setting.DatabaseSetting.GormLogSlowThreshold * time.Millisecond,
			LogLevel:      gormLogLevel(),
		})
	}

	switch setting.ServerSetting.RunMode {
	case "debug":
		return logger.Default.LogMode(logger.Info)
	default:
		return gormzap.New(logging.GormLogger, gormzap.Config{
			SlowThreshold: setting.DatabaseSetting.GormLogSlowThreshold * time.Millisecond,
			LogLevel:      gormLogLevel(),
		})
	}
}
