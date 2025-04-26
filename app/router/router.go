package router

import (
	"bufio"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/handlers"
	"github.com/codecrafters-io/http-server-starter-go/app/types"
)

var routes = map[string]string{
	"/":           "Welcome to the index page!",
	"/index.html": "Welcome to the index page!",
	"/hello":      "Hello, visitor!",
}

type Router struct {
	handlers *handlers.Handlers
}

func NewRouter(fileDir string) *Router {
	return &Router{
		handlers: handlers.NewHandlers(fileDir),
	}
}

func (r *Router) Route(rw *bufio.ReadWriter, req *types.Request) error {
	// 1) static routes
	if msg, ok := routes[req.Path]; ok {
		return r.handlers.HandleStaticRoute(rw.Writer, req, msg)
	}

	// 2) dynamic routes
	switch {
	case strings.HasPrefix(req.Path, "/echo/"):
		return r.handlers.HandleEcho(rw.Writer, req)

	case strings.HasPrefix(req.Path, "/user-agent"):
		return r.handlers.HandleUserAgent(rw.Writer, req)

	case req.Method == "GET" && strings.HasPrefix(req.Path, "/files/"):
		return r.handlers.HandleGetFile(rw.Writer, req)

	case req.Method == "POST" && strings.HasPrefix(req.Path, "/files/"):
		return r.handlers.HandlePostFile(rw.Writer, req)

	default:
		return r.handlers.HandleNotFound(rw.Writer, req)
	}
}
