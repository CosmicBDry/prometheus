package Logger

import (
	"github.com/siruspen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func SetLog() *logrus.Logger {
	logfile := &lumberjack.Logger{
		Filename:  "./logs/configAgent.log",
		MaxSize:   100,
		MaxAge:    30,
		Compress:  true,
		LocalTime: true,
	}
	defer logfile.Close()
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(logfile)
	return logger
}
