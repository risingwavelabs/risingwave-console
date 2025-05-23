package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogAgent struct {
	name   string
	fileds []zap.Field
}

var log *zap.Logger

func init() {
	logger, err := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	log = logger
}

func NewLogAgent(name string) *LogAgent {
	return &LogAgent{name: name, fileds: []zap.Field{zap.String("module", name)}}
}

func (a *LogAgent) AppendFiled(field zap.Field) *LogAgent {
	a.fileds = append(a.fileds, field)
	return a
}

// provide basic observability
func (a *LogAgent) Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, append(a.fileds, fields...)...)
}

// expected situation but worth a look
func (a *LogAgent) Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, append(a.fileds, fields...)...)
}

// unexpected error causing broken connection
func (a *LogAgent) Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, append(a.fileds, fields...)...)
}

// fatal error causing application shutdown
func (a *LogAgent) Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, append(a.fileds, fields...)...)
}

// provide basic observability
func (a *LogAgent) Infof(msg string, args ...any) {
	log.Info(fmt.Sprintf(msg, args...), a.fileds...)
}

// expected situation but worth a look
func (a *LogAgent) Warnf(msg string, args ...any) {
	log.Warn(fmt.Sprintf(msg, args...), a.fileds...)
}

// unexpected error causing broken connection
func (a *LogAgent) Errorf(msg string, args ...any) {
	log.Error(fmt.Sprintf(msg, args...), a.fileds...)
}
