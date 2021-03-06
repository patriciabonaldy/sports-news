package logger

import (
	"github.com/pkg/errors"
	"log"
	"os"
)

// Logger is the standard logger interface.
type Logger interface {
	Error(args ...interface{})
	ErrorTrace(err error)
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

type lg struct {
	logger *log.Logger
}

// New initializes a new logger.
func New() Logger {
	const flag = 5
	return &lg{logger: log.New(os.Stdout, "", flag)}
}

func (l *lg) Error(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *lg) ErrorTrace(err error) {
	l.Error(errors.Cause(err))
}

func (l *lg) Errorf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *lg) Info(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *lg) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
