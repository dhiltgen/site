package static

import (
	"net/http"
	//"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (srv *Server) GenerateMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	for route := range srv.Templates {
		if strings.HasSuffix(route, IndexFilename) {
			route = strings.TrimSuffix(route, IndexFilename)
			//route = path.Join(srv.SitePrefix, route)
			log.Debugf("Adding route %s", route)
			mux.Handle(route, srv)
		}
	}

	return mux, nil
}
