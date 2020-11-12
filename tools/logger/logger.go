package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logrusLogger *logrus.Logger
var logFileDesc *os.File

var logLevels = map[int]logrus.Level{
	10: logrus.DebugLevel,
	20: logrus.InfoLevel,
	30: logrus.WarnLevel,
	40: logrus.ErrorLevel,
	50: logrus.FatalLevel,
}

func InitLogger(filename string, logLevel int) {
	if logrusLogger != nil {
		return
	}

	// Open log file with creating dirs if necessary
	err := os.MkdirAll(filepath.Dir(filename), 0777)
	if err != nil {
		log.Fatal(err)
	}
	logFileDesc, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}

	logFile := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,  // megabytes after which new file is created
		MaxBackups: 3,  // number of backups
		MaxAge:     28, //days
	}

	logrusLogger = logrus.New()
	logrusLogger.SetLevel(logLevels[logLevel])
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	mw := io.MultiWriter(os.Stdout, logFile)
	logrusLogger.SetOutput(mw)
	logrus.RegisterExitHandler(closeLogFile)
}

func DisableLogger() {
	if logrusLogger == nil {
		logrusLogger = logrus.New()
	}
	logrusLogger.SetOutput(ioutil.Discard)
}

func Info(args ...interface{}) {
	logrusLogger.Info(args...)
}

func Error(args ...interface{}) {
	logrusLogger.Error(args...)
}

func Debug(args ...interface{}) {
	logrusLogger.Debug(args...)
}

func Warn(args ...interface{}) {
	logrusLogger.Warn(args...)
}

func Print(args ...interface{}) {
	logrusLogger.Print(args...)
}

func Println(args ...interface{}) {
	logrusLogger.Println(args...)
}

func closeLogFile() {
	if logFileDesc != nil {
		logFileDesc.Close()
	}
}
