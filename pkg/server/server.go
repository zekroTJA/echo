package server

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/k0kubun/pp/v3"
	"github.com/zekroTJA/echo/pkg/util"
	"github.com/zekroTJA/echo/pkg/verbosity"
)

//go:embed templates
var pages embed.FS

var tpl = template.Must(template.New("").ParseFS(pages, "templates/*.html"))

type Server struct {
	addr      string
	verbosity verbosity.Verbosity
	bodyLimit int
	mux       *http.ServeMux
	pp        *pp.PrettyPrinter
}

func New(addr string, verb verbosity.Verbosity, bodyLimit int) *Server {
	s := &Server{
		addr:      addr,
		verbosity: verb,
		bodyLimit: bodyLimit,
	}

	s.pp = pp.New()
	s.pp.SetColoringEnabled(false)

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.echoHandler)

	return s
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) echoHandler(w http.ResponseWriter, r *http.Request) {

	slog.Info("received request", "method", r.Method, "path", r.URL)

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
		raw, parsed, err := s.parseBody(r)
		if err != nil {
			respondError(w, err)
			return
		}

		echo.BodyString = string(raw)
		echo.BodyParsed = parsed
	}

	respondRes(w, r, &echo)
}

func respondError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	if errors.Is(err, util.ErrLimitReached) {
		status = http.StatusBadRequest
	}

	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func respondRes(w http.ResponseWriter, r *http.Request, res *echoObject) {

	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tpl.ExecuteTemplate(w, "main.html", res)
		if err != nil {
			respondError(w, err)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(res)
	}
}

func (s *Server) parseBody(r *http.Request) (raw []byte, parsed any, err error) {
	rdr := &util.LimitReader{Reader: r.Body, Limit: s.bodyLimit}

	raw, err = io.ReadAll(rdr)
	if err != nil {
		return nil, nil, err
	}

	contentType, contentTypeRest := parseContentType(r.Header.Get("Content-Type"))

	switch contentType {
	case "application/json":
		err = json.Unmarshal(raw, &parsed)
	case "multipart/form-data":
		boundary := strings.TrimPrefix(contentTypeRest, "boundary=")
		parsed, err = parseMultipartFormdata(raw, boundary)
	case "application/x-www-form-urlencoded":
		parsed, err = url.ParseQuery(string(raw))
	}

	if err != nil {
		parsed = fmt.Sprintf("<error: %s>", err.Error())
	}

	return raw, parsed, nil
}

func parseContentType(ct string) (typ string, rest string) {
	if ct == "" {
		return "", ""
	}

	ct = strings.ToLower(ct)

	i := strings.Index(ct, ";")
	if i < 0 {
		return strings.TrimSpace(ct), ""
	}

	return strings.TrimSpace(ct[:i]), strings.TrimSpace(ct[i+1:])
}

func parseMultipartFormdata(raw []byte, boundary string) (map[string]string, error) {
	rdr := multipart.NewReader(bytes.NewBuffer(raw), boundary)

	parsed := make(map[string]string)
	for {
		part, err := rdr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := part.FormName()
		data, err := io.ReadAll(part)
		if err != nil {
			return nil, err
		}

		if utf8.Valid(data) {
			parsed[name] = string(data)
		} else {
			parsed[name] = "data:" + base64.StdEncoding.EncodeToString(data)
		}
	}

	return parsed, nil
}
