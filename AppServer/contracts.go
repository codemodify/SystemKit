package servers

import (
	"net"

	"github.com/gorilla/mux"
)

// IServer -
type IServer interface {
	Run(ipPort string, enableCORS bool) error
	PrepareRoutes(router *mux.Router)
	RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool)
}

// APIEndpoint -
type APIEndpoint struct {
	Value string `json:"value,omitempty"`
}
