package server

import (
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/zekroTJA/echo/pkg/verbosity"
)

//go:embed templates
var pages embed.FS

var tpl = template.Must(template.New("").ParseFS(pages, "templates/*.html"))

type Server struct {
	addr      string
	verbosity verbosity.Verbosity

	router *gin.Engine
}

func New(addr string, verb verbosity.Verbosity, debug bool) *Server {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		addr:      addr,
		verbosity: verb,
		router:    gin.Default(),
	}

	s.registerHandlers()

	return s
}

func (s *Server) registerHandlers() {
	s.router.Any("/*path", s.echoHandler)
}

func (s *Server) echoHandler(c *gin.Context) {
	req := c.Request

	query := req.URL.Query()

	verb, err := verbosity.FromString(query.Get("verbosity"))
	if err != nil {
		verb = s.verbosity
	}

	var echo echoObject

	if verb >= verbosity.Minimal {
		echo.Method = req.Method
		echo.Path = req.URL.Path
	}

	if verb >= verbosity.Normal {
		echo.Host = req.Host
		echo.Query = query
		echo.Header = req.Header
		echo.RemoteAddress = req.RemoteAddr
	}

	if verb >= verbosity.Detailed {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			respondError(c, http.StatusInternalServerError, err)
			return
		}
		echo.BodyString = string(body)
	}

	respondRes(c, &echo)
}

func respondError(ctx *gin.Context, status int, err error) {
	ctx.String(status, err.Error())
}

func respondRes(ctx *gin.Context, res *echoObject) {

	accept := ctx.GetHeader("Accept")

	var data []byte
	var contentType string

	if strings.Contains(accept, "text/html") {
		buf := bytes.NewBuffer(data)
		err := tpl.ExecuteTemplate(buf, "main.html", res)
		if err != nil {
			respondError(ctx, http.StatusInternalServerError, err)
			return
		}
		data = buf.Bytes()
		contentType = "text/html; charset=utf-8"
	} else {
		data, _ = json.MarshalIndent(res, "", "  ")
		contentType = "application/json; charset=utf-8"
	}

	ctx.Data(http.StatusOK, contentType, data)
}

func (s *Server) Run() error {
	return s.router.Run(s.addr)
}
