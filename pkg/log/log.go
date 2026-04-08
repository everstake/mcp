package log

import (
	"github.com/sirupsen/logrus"
)

type (
	AppLogger struct {
		logger *logrus.Logger
	}
	Value struct {
		Value interface{}
		Key   string
	}
)

var (
	Logger *AppLogger
)

// Init logger on start
func init() {
	logger := logrus.New()
	Logger = &AppLogger{logger: logger}
}

func Setup(level string) {
	Logger.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	switch level {
	case "error":
		Logger.logger.SetLevel(logrus.ErrorLevel)
	case "info":
		Logger.logger.SetLevel(logrus.InfoLevel)
	default:
		Logger.logger.SetLevel(logrus.DebugLevel)
	}
}

func (l *AppLogger) Debug(msg string, values ...Value) {
	if len(values) > 0 {
		l.logger.WithFields(convertValues(values...)).Debug(msg)
	} else {
		l.logger.Debug(msg)
	}
}

func (l *AppLogger) Info(msg string, values ...Value) {
	if len(values) > 0 {
		l.logger.WithFields(convertValues(values...)).Info(msg)
	} else {
		l.logger.Info(msg)
	}
}

func (l *AppLogger) Warn(msg string, values ...Value) {
	if len(values) > 0 {
		l.logger.WithFields(convertValues(values...)).Warn(msg)
	} else {
		l.logger.Warn(msg)
	}
}

func (l *AppLogger) Error(msg string, values ...Value) {
	if len(values) > 0 {
		l.logger.WithFields(convertValues(values...)).Error(msg)
	} else {
		l.logger.Error(msg)
	}
}

func (l *AppLogger) Fatal(msg string, values ...Value) {
	if len(values) > 0 {
		l.logger.WithFields(convertValues(values...)).Panic(msg)
	} else {
		l.logger.Fatal(msg)
	}
}

func convertValues(values ...Value) map[string]interface{} {
	fieldsMap := make(map[string]interface{})
	for _, value := range values {
		fieldsMap[value.Key] = value.Value
	}
	return fieldsMap
}

func V(key string, value interface{}) Value {
	return Value{
		Key:   key,
		Value: value,
	}
}

func E(err error) Value {
	txt := ""
	if err != nil {
		txt = err.Error()
	}
	return Value{
		Key:   "error",
		Value: txt,
	}
}
