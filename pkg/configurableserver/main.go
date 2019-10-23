package configurableserver

import (
	"context"
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
	log := myLogger.GetLogger()

	// Create the server so we can pass the logger and configuration to it
	srv := NewConfigurableServer(myConfiguration, myLogger)
	// Block while server starts
	wg := &sync.WaitGroup{}
	wg.Add(1)
	log.Debugf("StartHTTPServer()...")
	http := srv.StartHTTPServer()
	wg.Wait() // This is blocking

	// We will never get past here!!!

	// Close the server gracefully ("shutdown")
	// timeout could be given with a proper context
	// (in real world you shouldn't use TODO()).
	if err := http.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	log.Printf("done. exiting. goodbye. and goodnight.")
}
