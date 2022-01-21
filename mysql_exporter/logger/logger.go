package logger

import (
	"github.com/siruspen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func SetAppLog(filepath string, maxsize, maxage int, localtime, compress bool) *logrus.Logger {
	logfile := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    maxsize,
		MaxAge:     maxage,
		MaxBackups: 6,
		LocalTime:  localtime,
		Compress:   compress,
	}

	defer logfile.Close()

	Logger := logrus.New()

	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetFormatter(&logrus.TextFormatter{})
	Logger.SetOutput(logfile)
	return Logger

}
