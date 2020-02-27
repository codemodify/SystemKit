package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/codemodify/systemkit-cryptography/gocrypto/ssh"
	"github.com/gorilla/mux"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"

	helpersReflect "github.com/codemodify/SystemKit/Helpers"
)

// SSHTunnelServer -
type SSHTunnelServer struct {
	sshServerConfig *ssh.ServerConfig
	server          IServer
	router          *mux.Router
}

// NewSSHTunnelServer -
func NewSSHTunnelServer(sshServerConfig *ssh.ServerConfig, server IServer) IServer {
	return &SSHTunnelServer{
		sshServerConfig: sshServerConfig,
		server:          server,
		router:          mux.NewRouter(),
	}
}

// Run - Implement `IServer`
func (thisRef *SSHTunnelServer) Run(ipPort string, enableCORS bool) error {

	//
	// BASED-ON: https://godoc.org/github.com/codemodify/systemkit-cryptography/gocrypto/ssh#example-NewServerConn
	//

	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	thisRef.PrepareRoutes(thisRef.router)
	thisRef.RunOnExistingListenerAndRouter(listener, thisRef.router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *SSHTunnelServer) PrepareRoutes(router *mux.Router) {
	thisRef.server.PrepareRoutes(router)
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *SSHTunnelServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("JM-SSH: failed to accept incoming connection: %s", err),
			})

			continue
		}

		go thisRef.runSSH(connection)
	}
}

type customResponseWriter struct {
	http.ResponseWriter
	sshChannel ssh.Channel
}

func (thisRef *customResponseWriter) Write(data []byte) (int, error) {
	logging.Instance().LogTraceWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("JM-SSH: sending back %d bytes", len(data)),
	})

	return thisRef.sshChannel.Write(data)
}

func (thisRef *SSHTunnelServer) runSSH(connection net.Conn) {
	// Before use, a handshake must be performed on the incoming connection
	sshServerConnection, chans, reqs, err := ssh.NewServerConn(connection, thisRef.sshServerConfig)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("JM-SSH: failed to handshake: %s", err),
		})

		return
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("JM-SSH: Connection %s", sshServerConnection.RemoteAddr().String()),
	})

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, _, err := newChannel.Accept()
		if err != nil {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("JM-SSH: could not accept channel: %v", err),
			})
			break
		}

		go func(ch ssh.Channel) {
			logging.Instance().LogTraceWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("JM-SSH: newChannel.Accept()"),
			})

			defer ch.Close()

			for {
				data := make([]byte, 1000000)
				len, err := ch.Read(data)
				if err != nil {
					if strings.Compare(err.Error(), "EOF") == 0 {
						logging.Instance().LogInfoWithFields(loggingC.Fields{
							"method":  helpersReflect.GetThisFuncName(),
							"message": fmt.Sprintf("JM-SSH: TRANSFER-FINISHED: %v", err),
						})
						break
					} else {
						logging.Instance().LogErrorWithFields(loggingC.Fields{
							"method":  helpersReflect.GetThisFuncName(),
							"message": fmt.Sprintf("JM-SSH: DATA-ERROR: %v", err),
						})
						break
					}
				}

				data = data[0:len]
				logging.Instance().LogDebugWithFields(loggingC.Fields{
					"method":  helpersReflect.GetThisFuncName(),
					"message": fmt.Sprintf("JM-SSH: DATA-TO-PASS-ON: %s", string(data)),
				})

				apiEndpoing := APIEndpoint{}
				err = json.Unmarshal(data, &apiEndpoing)
				if err != nil {
					logging.Instance().LogErrorWithFields(loggingC.Fields{
						"method":  helpersReflect.GetThisFuncName(),
						"message": fmt.Sprintf("JM-SSH: Missing ROUTE: %s", err.Error()),
					})
				}

				// Make `http.Request`
				request, err := http.NewRequest("POST", apiEndpoing.Value, bytes.NewBuffer(data))
				if err != nil {
					logging.Instance().LogErrorWithFields(loggingC.Fields{
						"method":  helpersReflect.GetThisFuncName(),
						"message": fmt.Sprintf("JM-SSH: SSH-DATA-ERROR: %s", err.Error()),
					})
					break
				}

				route := thisRef.router.Get(apiEndpoing.Value)
				if route == nil {
					logging.Instance().LogErrorWithFields(loggingC.Fields{
						"method":  helpersReflect.GetThisFuncName(),
						"message": fmt.Sprintf("JM-SSH: Missing ROUTE: %s", apiEndpoing.Value),
					})
					break
				}

				logging.Instance().LogTraceWithFields(loggingC.Fields{
					"method":  helpersReflect.GetThisFuncName(),
					"message": fmt.Sprintf("JM-SSH: ServeHTTP()"),
				})
				route.GetHandler().ServeHTTP(&customResponseWriter{sshChannel: ch}, request)

				break
			}
		}(channel)
	}
}
