package logger

import (
	"fmt"
	"os"
	"path/filepath"
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
	runtimeLog, err := os.OpenFile(
		filepath.Join("data", "logs", fmt.Sprintf("%s.log", conf.RunSince.Format("2006010215040507"))),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664)
	if err != nil {
		panic(fmt.Errorf("failed to open logfile err: %+w", err))
	}

	out := zerolog.MultiLevelWriter(os.Stdout, runtimeLog)
	return zerolog.New(out).With().Timestamp().Str("svc", conf.ServiceName).Logger().Hook()
}
