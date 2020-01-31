package servers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebScoketsRequestHandler -
type WebScoketsRequestHandler func(inChannel chan []byte, outChannel chan []byte)

// WebScoketsHandler -
type WebScoketsHandler struct {
	Route   string
	Handler WebScoketsRequestHandler
}

// WebScoketsServer -
type WebScoketsServer struct {
	handlers       []WebScoketsHandler
	routeToHandler map[string]WebScoketsHandler
	HTTPServer     IServer
	peers          []*websocket.Conn
	peersSync      sync.RWMutex
	enableCORS     bool
}

// NewWebScoketsServer -
func NewWebScoketsServer(handlers []WebScoketsHandler) IServer {

	var thisRef = &WebScoketsServer{
		handlers:       handlers,
		routeToHandler: map[string]WebScoketsHandler{},
		HTTPServer:     nil,
		peers:          []*websocket.Conn{},
		peersSync:      sync.RWMutex{},
	}

	var lowLevelRequestHelper = func(rw http.ResponseWriter, r *http.Request) {
		r.Header["Origin"] = nil

		var handler WebScoketsHandler = thisRef.routeToHandler[r.URL.Path]

		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return thisRef.enableCORS },
		}
		ws, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Print("upgrade: ", err)
			return
		}

		thisRef.setupCommunication(ws, &handler)
	}

	var HTTPHandlers = []HTTPHandler{}

	for _, handler := range thisRef.handlers {
		thisRef.routeToHandler[handler.Route] = handler

		HTTPHandlers = append(HTTPHandlers, HTTPHandler{
			Route:   handler.Route,
			Handler: lowLevelRequestHelper,
			Verb:    "GET",
		})
	}

	thisRef.HTTPServer = NewHTTPServer(HTTPHandlers)

	return thisRef
}

// Run - Implement `IServer`
func (thisRef *WebScoketsServer) Run(ipPort string, enableCORS bool) error {
	thisRef.enableCORS = enableCORS
	return thisRef.HTTPServer.Run(ipPort, enableCORS)
}

// PrepareRoutes - Implement `IServer`
func (thisRef *WebScoketsServer) PrepareRoutes(router *mux.Router) {
	thisRef.HTTPServer.PrepareRoutes(router)
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *WebScoketsServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	thisRef.HTTPServer.RunOnExistingListenerAndRouter(listener, router, enableCORS)
}

func (thisRef *WebScoketsServer) setupCommunication(ws *websocket.Conn, handler *WebScoketsHandler) {
	thisRef.addPeer(ws)

	var inChannel = make(chan []byte)
	var outChannel = make(chan []byte)

	var once sync.Once
	closeInChannel := func() {
		close(inChannel)
	}

	var wg sync.WaitGroup

	// outChannel -> PEER
	wg.Add(1)
	go func() {
		fmt.Println("SEND-TO-PEER - START")

		for {
			data, readOk := <-outChannel
			if !readOk {
				break
			}

			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				break
			}
		}

		fmt.Println("SEND-TO-PEER - END")
		once.Do(closeInChannel)
		wg.Done()
	}()

	// PEER -> inChannel
	wg.Add(1)
	go func() {
		fmt.Println("READ-FROM-PEER - START")

		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				break
			}

			var haveToStop = false
			select {
			case inChannel <- []byte(data):
			default:
				haveToStop = true
				break
			}

			if haveToStop {
				break
			}
		}

		fmt.Println("READ-FROM-PEER - END")
		once.Do(closeInChannel)
		wg.Done()
	}()

	go handler.Handler(inChannel, outChannel)

	wg.Wait()
	fmt.Println("setupCommunication - DONE")
	thisRef.removePeer(ws)
}

// SendToAllPeers -
func (thisRef *WebScoketsServer) SendToAllPeers(data []byte) {
	thisRef.peersSync.RLock()
	defer thisRef.peersSync.RUnlock()

	for _, conn := range thisRef.peers {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (thisRef *WebScoketsServer) addPeer(peer *websocket.Conn) {
	thisRef.peersSync.Lock()
	defer thisRef.peersSync.Unlock()

	thisRef.peers = append(thisRef.peers, peer)
}

func (thisRef *WebScoketsServer) removePeer(peer *websocket.Conn) {
	thisRef.peersSync.Lock()
	defer thisRef.peersSync.Unlock()

	index := -1
	for i, conn := range thisRef.peers {
		if conn == peer {
			index = i
			break
		}
	}
	if index != -1 {
		thisRef.peers = append(thisRef.peers[:index], thisRef.peers[index+1:]...)
	}

	peer.Close()
}
