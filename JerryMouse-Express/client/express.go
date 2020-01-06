package client

import (
	"go.isomorphicgo.org/go/isokit"
	"honnef.co/go/js/dom"
)

// RealtimeClient -
type RealtimeClient struct {
	Address string
	Peers   []string
}

var expressClient *ExpressClient

// NewExpressClient -
func NewExpressClient(
	pageLoadedHandler ExpressHandler,
	appContainerID string,
	handlers map[string]isokit.Handler,
) *ExpressClient {

	expressClient = &ExpressClient{}

	var pageLoaded = func() {
		templateSetChannel := make(chan *isokit.TemplateSet)
		go isokit.FetchTemplateBundle(templateSetChannel)

		// App Context
		expressClient.templateSet = <-templateSetChannel
		expressClient.Window = dom.GetWindow()
		expressClient.Document = dom.GetWindow().Document()
		expressClient.AppContainer = expressClient.Document.GetElementByID(appContainerID)

		routes := isokit.NewRouter()
		for route, handler := range handlers {
			routes.Handle(route, handler)
		}
		routes.Listen()

		pageLoadedHandler()
	}

	var domDocument = dom.GetWindow().Document().(dom.HTMLDocument)
	switch readyState := domDocument.ReadyState(); readyState {
	case "loading":
		domDocument.AddEventListener("DOMContentLoaded", false, func(dom.Event) {
			go pageLoaded()
		})
	case "interactive", "complete":
		go pageLoaded()
	default:
		println("Unexpected document.ReadyState value!")
	}

	return expressClient
}

// GetExpressClient -
func GetExpressClient() *ExpressClient {
	return expressClient
}

// func NewRealtimeClient() *RealtimeClient {
// 	return &RealtimeClient{
// 		Address: "",
// 		Peers:   []string{},
// 	}
// }

// func (rtc *RealtimeClient) ConnectToPeer(url url.URL, WebScoketsHandler *WebScoketsHandler) {
// 	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}

// 	helpers_RealtimeCommunicationHandler(c, WebScoketsHandler)
// }
