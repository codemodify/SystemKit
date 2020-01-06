package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	jm "github.com/codemodify/SystemKit/JerryMouse"

	"github.com/gorilla/mux"
	"go.isomorphicgo.org/go/isokit"
)

var expressServer *ExpressServer

// NewExpressServer -
func NewExpressServer(
	rootFolder string,
	templatesFolder string,
	staticPathRewrites map[string]string,
	staticFilesRewrites map[string]string,
	servers []jm.IServer,
) *ExpressServer {
	expressServer = &ExpressServer{
		rootFolder:          rootFolder,
		templatesFolder:     templatesFolder,
		staticPathRewrites:  staticPathRewrites,
		staticFilesRewrites: staticFilesRewrites,
		templates:           template.Must(template.ParseGlob(templatesFolder + "/*.html")),
		templateSet:         isokit.NewTemplateSet(),
		servers:             servers,
	}

	return expressServer
}

// GetExpressServer -
func GetExpressServer() *ExpressServer {
	return expressServer
}

// Run -
func (thisRef *ExpressServer) Run(ipPort string) {
	isokit.TemplateFilesPath = thisRef.templatesFolder
	isokit.TemplateFileExtension = ".html"

	thisRef.templateSet.GatherTemplates()

	router := mux.NewRouter()

	server := jm.NewMixedServer(thisRef.servers)
	server.PrepareRoutes(router)

	// Static paths
	for k, v := range thisRef.staticPathRewrites {
		var from = k
		var to = thisRef.rootFolder + v

		fmt.Println(fmt.Sprintf("Static paths: FROM [%s] TO [%s]", from, to))

		router.PathPrefix(from).Handler(http.StripPrefix(from,
			http.FileServer(http.Dir(to)),
		))
	}

	// Static files
	for k, v := range thisRef.staticFilesRewrites {
		var from = k
		var to = thisRef.rootFolder + v

		fmt.Println(fmt.Sprintf("Static files: FROM [%s] TO [%s]", from, to))

		router.Handle(from, func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, to)
			})
		}())
	}

	// template-bundle
	router.Handle(
		"/template-bundle",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var templateContentItemsBuffer bytes.Buffer
			enc := gob.NewEncoder(&templateContentItemsBuffer)
			m := thisRef.templateSet.Bundle().Items()
			err := enc.Encode(&m)
			if err != nil {
				log.Print("TemplateBundleHandler encoding err: ", err)
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(templateContentItemsBuffer.Bytes())
		}),
	)

	// Listen and RUN
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		fmt.Printf("Can't RUN: %s", err.Error())

		return
	}

	server.RunOnExistingListenerAndRouter(listener, router, true)
}

// RenderTemplate -
func (thisRef *ExpressServer) RenderTemplate(rw http.ResponseWriter, templateFile string, templateData interface{}) {
	err := thisRef.templates.ExecuteTemplate(rw, templateFile, templateData)
	if err != nil {
		fmt.Printf("RenderTemplate: %s", err.Error())
	}
}
