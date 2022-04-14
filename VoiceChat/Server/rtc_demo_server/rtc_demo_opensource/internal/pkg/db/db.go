package db

import (
	"context"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

var Client *gorm.DB

func Open(dsn string) {
	var err error
	Client, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: NewGormCustomLogger(),
	})
	if err != nil {
		panic("db init failed,error:" + err.Error())
	}
}

type GormCustomLogger struct {
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormCustomLogger() logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &GormCustomLogger{
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *GormCustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l GormCustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logs.CtxInfo(ctx, l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

// Warn print warn messages
func (l GormCustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logs.CtxWarn(ctx, l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

// Error print error messages
func (l GormCustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logs.CtxError(ctx, l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

// Trace print sql message
func (l GormCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if rows == -1 {
		logs.CtxError(ctx, l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
	} else {
		logs.CtxInfo(ctx, l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
