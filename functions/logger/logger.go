package logger

import (
	"fmt"
	"olx-clone/conf"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

var (
	Log Logger
)

func (l *Logger) Println(message ...any) {
	l.Info().Msg(fmt.Sprint(message...))
}

func (l *Logger) Errorln(message ...any) {
	err := fmt.Sprint(message...)
	withCaller := l.With().Caller().Logger()
	withCaller.Error().Msg(err)
}

func (l *Logger) Warnln(message ...any) {
	err := fmt.Sprint(message...)
	l.Warn().Msg(err)
}

func (l *Logger) Fatal(message ...any) {
	l.Error().Msg(fmt.Sprint(message...))
	os.Exit(1)
}

func (l *Logger) Panicf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.Error().Msg(message)
	panic(message)
}

func (l *Logger) Panicln(v ...any) {
	message := fmt.Sprint(v...)
	l.Error().Msg(message)
	panic(message)
}

func init() {
	logLevel := zerolog.DebugLevel
	if conf.ENV == conf.ENV_PROD {
		logLevel = zerolog.InfoLevel
	}
	log := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)
	Log = Logger{&log}

}

func WithRequest(c *gin.Context) *Logger {
	ctx := Log.With().Str("http_method", c.Request.Method)
	ctx = ctx.Str("remote_addr", c.ClientIP())
	ctx = ctx.Str("uri", c.FullPath())
	log := ctx.Logger()
	return &Logger{&log}
}
