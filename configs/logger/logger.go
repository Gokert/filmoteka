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

		//file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		//if err == nil {
		//	s.instance.Out = file
		//} else {
		//	s.instance.Info("Failed to log to file, using default stderr")
		//}

		s.instance.SetLevel(logrus.InfoLevel)
		s.instance.Infoln("logrus initialized")
	})

	return s.instance
}
