package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LogService struct {
	log *logrus.Logger
}

// Middleware is a struct that has a ServeHTTP method
func NewLogService() *LogService {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{}

	return &LogService{log}
}

func NewTestLogService() *LogService {
	ls := NewLogService()
	ls.GetLogger().Level = logrus.DebugLevel
	return ls
}

func (l *LogService) GetLogger() *logrus.Logger {
	return l.log
}

// Satisfy the middleware interface
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
