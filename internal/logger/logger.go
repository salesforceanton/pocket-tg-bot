package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func LogIssueWithPoint(point string, err error) {
	logrus.WithFields(logrus.Fields{
		"point":   point,
		"problem": fmt.Sprintf("Error appeared in [%s]: %s", point, err.Error()),
	}).Error(err)
}

func LogInfoWithPoint(point, info string) {
	logrus.Info(logrus.Fields{
		"point": point,
		"info":  info,
	})
}
