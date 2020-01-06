package client

import (
	"go.isomorphicgo.org/go/isokit"
	"honnef.co/go/js/dom"
)

// ExpressHandler -
type ExpressHandler func()

// ExpressClient -
type ExpressClient struct {
	templateSet  *isokit.TemplateSet
	handlers     map[string]isokit.Handler
	Window       dom.Window
	Document     dom.Document
	AppContainer dom.Element
}
