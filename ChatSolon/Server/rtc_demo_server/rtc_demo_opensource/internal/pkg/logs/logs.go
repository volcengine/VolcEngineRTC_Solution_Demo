package logs

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs/ctxvalues"

	"github.com/sirupsen/logrus"
)

func InitLog() {
	logfile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic("logrus open file failed,error:%s" + err.Error())
	}
	logrus.SetOutput(logfile)
}

func CtxInfo(ctx context.Context, format string, args ...interface{}) {
	logID, _ := ctxvalues.LogID(ctx)

	_, file, line, _ := runtime.Caller(1)
	s := strings.Split(file, "/")
	file = s[len(s)-1]
	logrus.WithFields(logrus.Fields{
		"Location": fmt.Sprintf("%s:%v", file, line),
		"LogID":    logID,
	}).Infof(format, args...)
}

func CtxWarn(ctx context.Context, format string, args ...interface{}) {
	logID, _ := ctxvalues.LogID(ctx)
	_, file, line, _ := runtime.Caller(1)
	s := strings.Split(file, "/")
	file = s[len(s)-1]
	logrus.WithFields(logrus.Fields{
		"Location": fmt.Sprintf("%s:%v", file, line),
		"LogID":    logID,
	}).Warnf(format, args...)
}

func CtxError(ctx context.Context, format string, args ...interface{}) {
	logID, _ := ctxvalues.LogID(ctx)
	_, file, line, _ := runtime.Caller(1)
	s := strings.Split(file, "/")
	file = s[len(s)-1]
	logrus.WithFields(logrus.Fields{
		"Location": fmt.Sprintf("%s:%v", file, line),
		"LogID":    logID,
	}).Errorf(format, args...)
}
func CtxDebug(ctx context.Context, format string, args ...interface{}) {
	logID, _ := ctxvalues.LogID(ctx)
	_, file, line, _ := runtime.Caller(1)
	s := strings.Split(file, "/")
	file = s[len(s)-1]
	logrus.WithFields(logrus.Fields{
		"Location": fmt.Sprintf("%s:%v", file, line),
		"LogID":    logID,
	}).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
func Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
