package main

import (
	"github.com/stevemcquaid/goprojectsetup/pkg/logger"
	"github.com/stevemcquaid/goprojectsetup/pkg/simpleserver"
)

func main() {
	myLogger := logger.NewLogService()
	log := myLogger.GetLogger()
	log.Debugf("Starting...")

	simpleserver.Run()
}
