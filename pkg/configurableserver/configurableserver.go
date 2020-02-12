package configurableserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/stevemcquaid/goprojectsetup/config"
	"github.com/stevemcquaid/goprojectsetup/pkg/logger"
	"github.com/stevemcquaid/goprojectsetup/pkg/simpleserver"
)

type ConfigurableServer struct {
	Config *config.Configuration
	Logger *logger.LogService
}

func NewConfigurableServer(config *config.Configuration, logger *logger.LogService) *ConfigurableServer {
	return &ConfigurableServer{
		Config: config,
		Logger: logger,
	}
}

type Server struct{}

func (srv *ConfigurableServer) StartHTTPServer() *http.Server {
	log := srv.Logger.GetLogger()

	// Build our routes
	r := mux.NewRouter()
	r.HandleFunc("/", srv.ServeHTTP).Methods("GET")

	// Declare & Start Server
	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", srv.Config.Server.Addr, srv.Config.Server.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithField("addr", server.Addr).Infof("server running on %s:%d", srv.Config.Server.Addr, srv.Config.Server.Port)

	go func() {
		// returns ErrServerClosed on graceful close
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.WithError(err).Fatalf("failed to start up server. ListenAndServe(): %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return server
}

// Example of middleware
// ConfigurableServer implements `http.Handler` interface because it has a method: `ServeHTTP(w http.ResponseWriter, r *http.Request)`
func (srv *ConfigurableServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	simpleserver.HelloHandler(w, r)
}
