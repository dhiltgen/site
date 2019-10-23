package static

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	Port         int
	SitePrefix   string
	BlobPrefix   string
	TemplatesDir string
	BlobDir      string

	// Templates is keyed by path, and contains the template
	Templates map[string]*template.Template
}

func (srv *Server) Serve() error {
	log.Infof("Starting server on %d with prefix %s", srv.Port, srv.SitePrefix)
	err := srv.LoadTemplates()
	if err != nil {
		return err
	}
	mux, err := srv.GenerateMux()
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", srv.Port), mux)
}
