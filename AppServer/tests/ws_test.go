package tests

import (
	"fmt"
	"testing"

	appS "github.com/codemodify/SystemKit/AppServer"
)

func Test_01(t *testing.T) {
	a := NewWebSocketTransport(64000, messageEventHandler)
	a.Serve()
}

func messageEventHandler(incomingData []byte) []byte {
	//
	// build the right object
	//

	fmt.Println(fmt.Sprintf("messageEventHandler(): %s", string(incomingData)))

	return incomingData
}

// OnMessageEventHandler -
type OnMessageEventHandler func([]byte) []byte

// WebSocketTransport -
type WebSocketTransport struct {
	port                int
	server              appS.IServer
	messageEventHandler OnMessageEventHandler
}

// NewWebSocketTransport -
func NewWebSocketTransport(port int, messageEventHandler OnMessageEventHandler) *WebSocketTransport {
	return &WebSocketTransport{
		port:                port,
		server:              nil,
		messageEventHandler: messageEventHandler,
	}
}

func (thisRef WebSocketTransport) Serve() error {
	thisRef.server = appS.NewWebScoketsServer([]appS.WebScoketsHandler{
		appS.WebScoketsHandler{
			Route:   "/",
			Handler: thisRef.rawRequestHandler,
		},
	})
	return thisRef.server.Run(fmt.Sprintf(":%d", thisRef.port), true)
}

func (thisRef WebSocketTransport) rawRequestHandler(inChannel chan []byte, outChannel chan []byte) { //, done chan bool) {
	// DOX:
	//    To indicate DONE close the `outChannel`
	//    If error when reading on `inChannel` means connection was closed, do not send data

	for {
		data, readOk := <-inChannel
		if !readOk {
			break
		}

		println("rawRequestHandler-DEBUG: RECV: " + string(data))

		if thisRef.messageEventHandler != nil {
			dataToSend := thisRef.messageEventHandler(data)
			println("rawRequestHandler-DEBUG: SEND: " + string(dataToSend))

			outChannel <- []byte(dataToSend)
		}
	}

	println("rawRequestHandler-DEBUG: DONE")
	close(outChannel)
}
