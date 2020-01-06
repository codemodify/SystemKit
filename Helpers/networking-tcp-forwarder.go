package Helpers

import (
	"io"
	"net"
	"time"
)

// ErrorHandler -
type ErrorHandler func(error)

// NewTCPForwarder -
func NewTCPForwarder(listenOn string, gatewayTo string, maxConnectionTimeoutIfNoActivity time.Duration, onAcceptError ErrorHandler, onTargetConnectionError ErrorHandler) error {
	listener, err := net.Listen("tcp4", listenOn)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			if onAcceptError != nil {
				onAcceptError(err)
			}
			continue
		}

		go handleConnection(connection, gatewayTo, onTargetConnectionError, maxConnectionTimeoutIfNoActivity)
	}
}

func handleConnection(connection net.Conn, gatewayTo string, onTargetConnectionError ErrorHandler, maxConnectionTimeoutIfNoActivity time.Duration) {
	passThroughConnection, err := net.Dial("tcp", gatewayTo)
	if onTargetConnectionError != nil {
		onTargetConnectionError(err)
	} else {
		defer connection.Close()
		defer passThroughConnection.Close()

		connection.SetDeadline(time.Now().Add(maxConnectionTimeoutIfNoActivity))
		passThroughConnection.SetDeadline(time.Now().Add(maxConnectionTimeoutIfNoActivity))

		go io.Copy(connection, passThroughConnection)
		io.Copy(passThroughConnection, connection)
	}
}
