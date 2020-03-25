package servers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"

	helpersReflect "github.com/codemodify/systemkit-helpers"
)

// JSONResponse -
type JSONResponse struct {
	HasError bool        `json:"hasError"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// JSONRequestHandler -
type JSONRequestHandler func(data []byte) JSONResponse

// JSONHandler -
type JSONHandler struct {
	Route    string
	Template interface{}
	Handler  JSONRequestHandler
}

// JSONServer -
type JSONServer struct {
	handlers       []JSONHandler
	routeToHandler map[string]JSONHandler
	HTTPServer     IServer
}

// NewJSONServer -
func NewJSONServer(handlers []JSONHandler) *JSONServer {
	var thisRef = &JSONServer{
		handlers:       handlers,
		routeToHandler: map[string]JSONHandler{},
		HTTPServer:     nil,
	}

	var lowLevelRequestHelper = func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var jsonHandler = thisRef.routeToHandler[r.URL.Path]

		// Pass Object
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("Error reading body: %v", err),
			})

			http.Error(rw, "can't read body", http.StatusBadRequest)
			return
		}

		var jsonResponse = jsonHandler.Handler(body)
		json.NewEncoder(rw).Encode(jsonResponse)
	}

	var HTTPHandlers = []HTTPHandler{}

	for _, handler := range thisRef.handlers {
		thisRef.routeToHandler[handler.Route] = handler

		HTTPHandlers = append(HTTPHandlers, HTTPHandler{
			Route:   handler.Route,
			Handler: lowLevelRequestHelper,
			Verb:    "POST",
		})
	}

	thisRef.HTTPServer = NewHTTPServer(HTTPHandlers)

	return thisRef
}

// Run - Implement `IServer`
func (thisRef *JSONServer) Run(ipPort string, enableCORS bool) error {
	return thisRef.HTTPServer.Run(ipPort, enableCORS)
}

// PrepareRoutes - Implement `IServer`
func (thisRef *JSONServer) PrepareRoutes(router *mux.Router) {
	thisRef.HTTPServer.PrepareRoutes(router)
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *JSONServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	thisRef.HTTPServer.RunOnExistingListenerAndRouter(listener, router, enableCORS)
}

// func (jsonData *JsonData) ToObject(objectInstance interface{}) {
// 	// Do JSON to Object Mapping
// 	objectValue := reflect.ValueOf(objectInstance).Elem()
// 	for i := 0; i < objectValue.NumField(); i++ {
// 		field := objectValue.Field(i)
// 		fieldName := objectValue.Type().Field(i).Name

// 		if valueToCopy, ok := (*jsonData)[fieldName]; ok {
// 			if !field.CanInterface() {
// 				continue
// 			}
// 			switch field.Interface().(type) {
// 			case string:
// 				valueToCopyAsString := reflect.ValueOf(valueToCopy).String()
// 				objectValue.Field(i).SetString(valueToCopyAsString)
// 				break
// 			case int:
// 				valueToCopyAsInt := int64(reflect.ValueOf(valueToCopy).Float())
// 				objectValue.Field(i).SetInt(valueToCopyAsInt)
// 				break
// 			case float64:
// 				valueToCopyAsFloat := reflect.ValueOf(valueToCopy).Float()
// 				objectValue.Field(i).SetFloat(valueToCopyAsFloat)
// 				break
// 			default:
// 			}
// 		}
// 	}
// }

// Get JSON fields
//var jsonData JsonData
//_ = json.NewDecoder(r.Body).Decode(&jsonData)

// TRACE
// if false {
// 	reqAsJSON, _ := json.Marshal(req)
// 	fmt.Println(fmt.Sprintf("%s -> %s", Utils.CallStack(), string(reqAsJSON)))
// }

//jsonData.ToObject(jsonHandler.Template)

// Pass Object
//var response JsonResponse = jsonHandler.Handler(jsonData)
