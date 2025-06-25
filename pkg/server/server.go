package server

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/zekroTJA/echo/pkg/verbosity"
)

//go:embed templates
var pages embed.FS

var tpl = template.Must(template.New("").ParseFS(pages, "templates/*.html"))

type Server struct {
	addr      string
	verbosity verbosity.Verbosity
	mux       *http.ServeMux
}

func New(addr string, verb verbosity.Verbosity) *Server {
	s := &Server{
		addr:      addr,
		verbosity: verb,
	}

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.echoHandler)

	return s
}

func (s *Server) echoHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	verb, err := verbosity.FromString(query.Get("verbosity"))
	if err != nil {
		verb = s.verbosity
	}

	var echo echoObject

	if verb >= verbosity.Minimal {
		echo.Method = r.Method
		echo.Path = r.URL.Path
	}

	if verb >= verbosity.Normal {
		echo.Host = r.Host
		echo.Query = query
		echo.Header = r.Header
		echo.RemoteAddress = r.RemoteAddr
	}

	if verb >= verbosity.Detailed {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		echo.BodyString = string(body)
	}

	respondRes(w, r, &echo)
}

func respondError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func respondRes(w http.ResponseWriter, r *http.Request, res *echoObject) {

	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tpl.ExecuteTemplate(w, "main.html", res)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(res)
	}
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.addr, s.mux)
}
