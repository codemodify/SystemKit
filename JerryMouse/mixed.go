package servers

import (
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"

	helpersReflect "github.com/codemodify/SystemKit/Helpers"
)

// MixedServer -
type MixedServer struct {
	servers []IServer
}

// NewMixedServer -
func NewMixedServer(servers []IServer) IServer {
	return &MixedServer{
		servers: servers,
	}
}

// Run - Implement `IServer`
func (thisRef *MixedServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	var router = mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *MixedServer) PrepareRoutes(router *mux.Router) {
	for _, server := range thisRef.servers {
		server.PrepareRoutes(router)
	}
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *MixedServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	if enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		err := http.Serve(listener, corsSetterHandler)
		if err != nil {
			logging.Instance().LogFatalWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": err.Error(),
			})

			os.Exit(-1)
		}
	} else {
		err := http.Serve(listener, router)
		if err != nil {
			logging.Instance().LogFatalWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": err.Error(),
			})

			os.Exit(-1)
		}
	}
}
