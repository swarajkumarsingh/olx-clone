package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"olx-clone/conf"
	"olx-clone/constants/messages"
	models "olx-clone/models/error"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var STRING string = "string"
var INT string = "int"

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

func (l *Logger) Panicln(values ...any) {
	valueType := checkValueType(values[0])
	if valueType == STRING {
		panicString(l, values[0])
	} else if isError(values[0]) {
		panicError(l, values[0])
	}
	panicCodeAndString(l, values[0], values[1])
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

func structToString(s models.CustomError) (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func convertAnyToInteger(val interface{}) (int, error) {
	if intValue, ok := val.(int); ok {
		return intValue, nil
	}
	if floatValue, ok := val.(float64); ok {
		return int(floatValue), nil
	}
	return 0, fmt.Errorf(messages.SomethingWentWrongMessage)
}

func convertAnyToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func checkValueType(value any) string {
	if fmt.Sprintf("%v", value) == value {
		return STRING
	} else {
		return INT
	}
}

func panicError(l *Logger, value interface{}) {
	err, _ := value.(error)
	panicString(l, err.Error())
}

func panicString(l *Logger, value interface{}) {
	code := http.StatusConflict

	message := convertAnyToString(value)

	model := models.CustomError{
		Status_code: code,
		Message:     message,
	}
	msg, err := structToString(model)
	if err != nil {
		l.Error().Msg(message)
		panic(message)
	}

	l.Error().Msg(msg)
	panic(msg)
}

func isError(err any) bool {
	_, isErr := err.(error)
	return isErr
}

func panicCodeAndString(l *Logger, a interface{}, b interface{}) {
	code, err := convertAnyToInteger(a)
	if err != nil {
		code = http.StatusConflict
	}

	message := convertAnyToString(b)

	model := models.CustomError{
		Status_code: code,
		Message:     message,
	}

	msg, err := structToString(model)
	if err != nil {
		l.Error().Msg(message)
		panic(message)
	}

	l.Error().Msg(msg)
	panic(msg)
}
