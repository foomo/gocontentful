package test

import (
	"fmt"
	"os"

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
	testFile, err := GetTestFile("./test-space-export.json")
	if err != nil {
		return nil, fmt.Errorf("getTestClient could not read space export file: %v", err)
	}
	return testapi.NewOfflineContentfulClient(testFile, GetContenfulLogger(testLogger), LogDebug, true, true)
}

func GetTestFile(filename string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("getTestFile could not read space export file: %v", err)
	}
	return fileBytes, nil
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
