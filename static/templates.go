package static

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// HeaderFilename stores the full filename of the header template
	HeaderFilename = "header.tmpl"
	// FooterFilename stores the full filename of the footer template
	FooterFilename = "footer.tmpl"
	// IndexFilename stores the full filename of the index template
	IndexFilename = "index.tmpl"
)

type RequestSettings struct {
	Server

	NavSection string
}

func (srv *Server) LoadTemplates() error {
	root, err := filepath.Abs(srv.TemplatesDir)
	srv.Templates = map[string]*template.Template{}
	if err != nil {
		return err
	}
	err = filepath.Walk(root, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("Failed to process template root %s: %s", root, err)
			return err
		}
		// Strip off root
		routePath := strings.TrimPrefix(filename, root)
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(filename, ".tmpl") {
			log.Debugf("loading template %s", routePath)
			t, err := template.ParseFiles(filename)
			if err != nil {
				return err
			}
			srv.Templates[routePath] = t
		} else {
			log.Infof("ignoring non-template file %s in template dir", filename)
		}
		return nil
	})

	return err
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := strings.TrimPrefix(r.URL.Path, srv.SitePrefix)
	if route == "" {
		route = "/"
	}

	// TODO consider cleaning off trailing index.htm[l]

	log.Debugf("request: %s", route)

	bodyTmpl, err := srv.FindBody(route)
	if err != nil {
		log.Errorf("404: %s", err)
		// TODO make a fancy 404
		http.NotFound(w, r)
		return
	}
	params := RequestSettings{
		Server:     *srv,
		NavSection: path.Base(route),
	}

	// Find the relevant header and footers
	dir := path.Dir(route)
	headerTmpl := srv.FindHeader(dir)
	footerTmpl := srv.FindFooter(dir)

	// stitch together the response payload
	err = headerTmpl.Execute(w, params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = bodyTmpl.Execute(w, params)
	if err != nil {
		log.Errorf("%s - body %s", route, err)
		return
	}
	err = footerTmpl.Execute(w, params)
	if err != nil {
		log.Errorf("%s - footer %s", route, err)
		return
	}
}

func (srv *Server) FindInParent(dir, filename string) *template.Template {
	search := path.Join(dir, filename)
	for {
		t, ok := srv.Templates[search]
		if ok {
			return t
		}
		if search == filename {
			log.Errorf("failed to find header")
			return template.New("")
		}
		search = path.Join(path.Dir(path.Dir(search)), filename)
	}
}
func (srv *Server) FindHeader(dir string) *template.Template {
	return srv.FindInParent(dir, HeaderFilename)
}

func (srv *Server) FindFooter(dir string) *template.Template {
	return srv.FindInParent(dir, FooterFilename)
}

func (srv *Server) FindBody(route string) (*template.Template, error) {
	bodyName := path.Join(route, IndexFilename)
	t, ok := srv.Templates[bodyName]
	if !ok {
		return nil, fmt.Errorf("%s not found", bodyName)
	}
	return t, nil
}
