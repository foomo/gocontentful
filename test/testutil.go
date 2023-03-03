package test

import (
	"github.com/foomo/gocontentful/test/testapi"
	"github.com/sirupsen/logrus"
)

const (
	LogDebug = 0
	LogInfo  = 1
	LogWarn  = 2
	LogError = 3
)

var testLogger = logrus.StandardLogger()

func getTestClient() (*testapi.ContentfulClient, error) {
	return testapi.NewOfflineContentfulClient("./test-space-export.json", GetContenfulLogger(testLogger), LogDebug, true, true)
}

func GetContenfulLogger(log *logrus.Logger) func(fields map[string]interface{}, level int, args ...interface{}) {
	return func(fields map[string]interface{}, level int, args ...interface{}) {
		if args == nil {
			return
		}
		switch level {
		case LogInfo:
			log.WithFields(fields).Info(args[0])
		case LogWarn:
			log.WithFields(fields).Warn(args[0])
		case LogError:
			log.WithFields(fields).Error(args[0])
		default:
			return
		}
	}
}
