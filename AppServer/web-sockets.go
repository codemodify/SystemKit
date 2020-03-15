package servers

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/codemodify/SystemKit/Helpers/channels"
)

// WebScoketsRequestHandler -
type WebScoketsRequestHandler func(inChannel chan []byte, outChannel chan []byte)

// WebSocketsHandler -
type WebSocketsHandler struct {
	Route   string
	Handler WebScoketsRequestHandler
}

// WebScoketsServer -
type WebScoketsServer struct {
	handlers       []WebSocketsHandler
	routeToHandler map[string]WebSocketsHandler
	HTTPServer     IServer
	peers          []*websocket.Conn
	peersSync      sync.RWMutex
	enableCORS     bool
}

// NewWebSocketsServer -
func NewWebSocketsServer(handlers []WebSocketsHandler) IServer {

	var thisRef = &WebScoketsServer{
		handlers:       handlers,
		routeToHandler: map[string]WebSocketsHandler{},
		HTTPServer:     nil,
		peers:          []*websocket.Conn{},
		peersSync:      sync.RWMutex{},
	}

	var lowLevelRequestHelper = func(rw http.ResponseWriter, r *http.Request) {
		r.Header["Origin"] = nil

		var handler WebSocketsHandler = thisRef.routeToHandler[r.URL.Path]

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

func (thisRef *WebScoketsServer) setupCommunication(ws *websocket.Conn, handler *WebSocketsHandler) {
	// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: setupCommunication - START")

	thisRef.addPeer(ws)

	var inChannel = make(chan []byte)  // data from WS
	var outChannel = make(chan []byte) // data to WS

	go handler.Handler(inChannel, outChannel)

	var once sync.Once
	closeInChannel := func() {
		close(inChannel)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: SEND-TO-PEER - START")

		for {
			data, readOk := <-outChannel
			if !readOk { // if CHANNEL closed - means communication ended by the handler
				// DEBUG: fmt.Println(fmt.Sprint("AppServer-WebScokets-DEBUG: SEND-TO-PEER - communication ended by the handler"))
				break
			}

			// DEBUG: fmt.Println(fmt.Sprint("AppServer-WebScokets-DEBUG: SEND-TO-PEER - DATA: ", string(data)))

			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil { // if can't send - means communication ended by the peer
				// DEBUG: fmt.Println(fmt.Sprint("AppServer-WebScokets-DEBUG: SEND-TO-PEER - send - communication ended by the peer"))
				break
			}
		}

		once.Do(closeInChannel)

		// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: SEND-TO-PEER - END")
		wg.Done()

	}()

	wg.Add(1)
	go func() {
		// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: READ-FROM-PEER - START")

		for {
			_, data, err := ws.ReadMessage()
			if err != nil { // if can't read - means communication ended by the peer
				// DEBUG: fmt.Println(fmt.Sprint("AppServer-WebScokets-DEBUG: SEND-FROM-PEER - read - communication ended by the peer"))
				once.Do(closeInChannel)
				break
			}

			// DEBUG: fmt.Println(fmt.Sprint("AppServer-WebScokets-DEBUG: READ-FROM-PEER - DATA: ", string(data)))

			if channels.IsClosed(inChannel) {
				break
			}
			inChannel <- []byte(data)
		}

		// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: READ-FROM-PEER - END")
		wg.Done()
	}()

	wg.Wait()
	thisRef.removePeer(ws)

	// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: setupCommunication - DONE")
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

	// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: addPeer")

	thisRef.peers = append(thisRef.peers, peer)
}

func (thisRef *WebScoketsServer) removePeer(peer *websocket.Conn) {
	thisRef.peersSync.Lock()
	defer thisRef.peersSync.Unlock()

	// DEBUG: fmt.Println("AppServer-WebScokets-DEBUG: removePeer")

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
