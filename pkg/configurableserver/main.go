package configurableserver

import (
	"sync"

	configuration "github.com/stevemcquaid/goprojectsetup/config"
	"github.com/stevemcquaid/goprojectsetup/pkg/config"
	"github.com/stevemcquaid/goprojectsetup/pkg/logger"
)

func Run() {
	// Get configuration
	myConfiguration := &configuration.Configuration{}
	err := config.SetupConfig(myConfiguration)
	if err != nil {
		panic(err)
	}

	// Configure the logger
	myLogger := logger.NewLogService(myConfiguration)
	_ = myLogger.GetLogger()

	// Create the custom server so we can pass the logger and configuration to it
	srv := NewConfigurableServer(myConfiguration, myLogger)
	_ = srv.StartHTTPServer()

	// Block while server runs
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait() // This is the blocking command
}
