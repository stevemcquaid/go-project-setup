package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stevemcquaid/goprojectsetup/config"
)

type LogService struct {
	log *logrus.Logger
}

func NewLogService(configuration *config.Configuration) *LogService {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}

	logrusLevel, err := configuration.Logger.GetLogLevel()
	if err != nil {
		log.Level = logrus.WarnLevel
	} else {
		log.Level = logrusLevel
	}

	return &LogService{log}
}

func (l *LogService) GetLogger() *logrus.Logger {
	return l.log
}

// Satisfy the middleware interface
// Middleware is a struct that has a ServeHTTP method
func (l *LogService) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	// log.Infof("Here in the logHandler")
	start := time.Now()
	next.ServeHTTP(w, r)
	l.log.Infof("%s | %s | %s | %s", time.Since(start), r.Method, r.URL.Path, r.RequestURI)
}

func (l *LogService) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Infof("Here in the logHandler")
		start := time.Now()
		next.ServeHTTP(w, r)
		l.log.Infof("%s | %s | %s | %s", time.Since(start), r.Method, r.URL.Path, r.RequestURI)
	},
	)
}

func NewTestLogService(configuration *config.Configuration) *LogService {
	ls := NewLogService(configuration)
	ls.GetLogger().Level = logrus.DebugLevel
	return ls
}
