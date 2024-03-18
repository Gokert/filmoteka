package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type singleton struct {
	once     sync.Once
	instance *logrus.Logger
}

func GetLogger() *logrus.Logger {
	s := &singleton{}

	s.once.Do(func() {
		s.instance = logrus.New()
		s.instance.SetLevel(logrus.InfoLevel)
		s.instance.Infoln("logrus initialized")
	})

	return s.instance
}
