package servers

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"

	helpersReflect "github.com/codemodify/systemkit-helpers"
)

// HTTPRequestHandler -
type HTTPRequestHandler func(rw http.ResponseWriter, r *http.Request)

// HTTPHandler -
type HTTPHandler struct {
	Route   string
	Verb    string
	Handler HTTPRequestHandler
}

// HTTPServer -
type HTTPServer struct {
	handlers []HTTPHandler
}

// NewHTTPServer -
func NewHTTPServer(handlers []HTTPHandler) IServer {
	return &HTTPServer{
		handlers: handlers,
	}
}

// Run - Implement `IServer`
func (thisRef *HTTPServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *HTTPServer) PrepareRoutes(router *mux.Router) {
	for _, handler := range thisRef.handlers {
		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s - for %s", handler.Route, handler.Verb),
		})

		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb, "OPTIONS").Name(handler.Route)
	}
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *HTTPServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
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
