package log

import (
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	logPath string
	Log     *logrus.Logger
)

func init() {
	dir, _ := os.Getwd()
	logPath = fmt.Sprintf("%s/logs/deeplx.log", dir)
}

func InitLog() {
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logPath,
		logrus.ErrorLevel: logPath,
	}

	log := logrus.New()
	log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	Log = log
}

func Info(args ...interface{}) {
	verbose(func() {
		Log.Info(args)
	})
}

func Infof(format string, args ...interface{}) {
	verbose(func() {
		Log.Infof(format, args)
	})
}

func Errorf(format string, args ...interface{}) {
	verbose(func() {
		Log.Errorf(format, args)
	})
}

func verbose(f func()) {
	if ok := os.Getenv("VERBOSE"); ok == "true" {
		f()
	}
}
