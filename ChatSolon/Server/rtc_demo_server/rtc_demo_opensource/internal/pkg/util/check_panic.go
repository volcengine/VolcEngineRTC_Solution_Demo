package util

import (
	"runtime"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

// const values
const (
	MaxStack = 5
)

// CheckPanic checks panics
func CheckPanic() {
	if err := recover(); err != nil {
		for i := 0; i <= MaxStack; i++ {
			if pc, file, line, ok := runtime.Caller(i); ok {
				f := runtime.FuncForPC(pc)
				logs.Error("panic error happened, stack:%d, med:%s, file:%s, line:%d, error:%s", i, f.Name(), file, line, err)
			}

		}
		logs.Error("panic error happened, error:%s", err)
	}
}
