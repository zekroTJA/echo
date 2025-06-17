package server

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/zekroTJA/echo/pkg/verbosity"
	"gopkg.in/yaml.v2"
)

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
			respondError(c, 500, err)
			return
		}
		echo.BodyString = string(body)
	}

	respondRes(c, echo)
}

func respondError(ctx *gin.Context, status int, err error) {
	ctx.String(status, err.Error())
}

func respondRes(ctx *gin.Context, res interface{}) {
	typ := ctx.Query("type")
	var data []byte
	var contentType string

	switch typ {
	case "yml", "yaml":
		data, _ = yaml.Marshal(res)
		contentType = "text/yaml"
	default:
		data, _ = json.MarshalIndent(res, "", "  ")
		contentType = "application/json"
	}

	ctx.Data(200, contentType, data)
}

func (s *Server) Run() error {
	return s.router.Run(s.addr)
}
