package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"github.com/codecrafters-io/http-server-starter-go/app/router"
)

type Server struct {
	config *Config
	router *router.Router
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: router.NewRouter(config.FileDir()),
	}
}

func (s *Server) Start() {
	addr := ":" + s.config.Port()
	ln := s.createListener(addr)
	defer ln.Close()

	log.Printf("Listening on %s", addr)
	s.acceptConnections(ln)
}

// createListener initializes the listener.
func (s *Server) createListener(addr string) net.Listener {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to %s: %v", addr, err)
	}
	return ln
}

// acceptConnections handles incoming connections.
func (s *Server) acceptConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection processes a single connection.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(reader, writer)
	for {
		req, err := http.ParseRequest(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			http.RespondBadRequest(writer)
			break
		}
		if err := s.router.Route(rw, req); err != nil {
			log.Printf("Error routing %s: %v", req.Path, err)
		}
		// Exit if client asks to close the connection.
		if connHeader, ok := req.Headers["Connection"]; ok && strings.ToLower(connHeader) == "close" {
			break
		}
	}
}
