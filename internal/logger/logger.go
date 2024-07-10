package logger

import (
	"os"
	"runtime"
	"strings"

	"github.com/nmluci/da-farm-be/internal/config"
	"github.com/rs/zerolog"
)

func callerNameHook() zerolog.HookFunc {
	return func(e *zerolog.Event, level zerolog.Level, message string) {
		pc, _, _, ok := runtime.Caller(4)
		if !ok {
			return
		}

		funcname := runtime.FuncForPC(pc).Name()
		fn := funcname[strings.LastIndex(funcname, "/")+1:]
		e.Str("caller", fn)
	}
}

func New(conf *config.Config) zerolog.Logger {
	out := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	return zerolog.New(out).With().Timestamp().Str("svc", conf.ServiceName).Logger().Hook()
}
