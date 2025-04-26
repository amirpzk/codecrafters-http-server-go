package handlers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"github.com/codecrafters-io/http-server-starter-go/app/types"
)

type Handlers struct {
	fileDir string // moved configuration inside the struct
}

func NewHandlers(dir string) *Handlers {
	return &Handlers{fileDir: dir} // initialize field rather than global variable
}

// Updated signature: receive full request.
func (h *Handlers) HandleGetFile(w *bufio.Writer, req *types.Request) error {
	name := req.Path[len("/files/"):]
	full := filepath.Join(h.fileDir, name) // updated usage: h.fileDir
	data, err := os.ReadFile(full)
	if err != nil {
		return http.WritePlainTextResponse(w, req, 404, "Not Found", "File not found")
	}
	return http.WriteResponse(w, req, 200, "OK", "application/octet-stream", data)
}

func (h *Handlers) HandlePostFile(w *bufio.Writer, req *types.Request) error {
	name := req.Path[len("/files/"):]
	full := filepath.Join(h.fileDir, name)
	if err := os.WriteFile(full, req.Body, 0644); err != nil {
		return http.WritePlainTextResponse(w, req, 500, "Internal Server Error", "Could not write file")
	}
	return http.WritePlainTextResponse(w, req, 201, "Created", "File saved")
}

func (h *Handlers) HandleNotFound(w *bufio.Writer, req *types.Request) error {
	return http.WritePlainTextResponse(w, req, 404, "Not Found", "404 page not found")
}

// Now using req to pass along the connection header.
func (h *Handlers) HandleUserAgent(w *bufio.Writer, req *types.Request) error {
	return http.WritePlainTextResponse(w, req, 200, "OK", req.Headers["User-Agent"])
}

func (h *Handlers) HandleEcho(w *bufio.Writer, req *types.Request) error {
	// prepare body
	body := []byte(req.Path[len("/echo/"):])

	// If Accept-Encoding contains gzip compress the response.
	if supportsGzip(req) {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, err := gz.Write(body)
		if err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}
		body = buf.Bytes()
		// Mark response as gzipped.
		req.Headers["Response-Content-Encoding"] = "gzip"
	}
	return http.WriteResponse(w, req, 200, "OK", "text/plain", body)
}

// Updated: pass req along.
func (h *Handlers) HandleStaticRoute(w *bufio.Writer, req *types.Request, msg string) error {
	return http.WritePlainTextResponse(w, req, 200, "OK", msg)
}

// supportsGzip reports whether the client sent “gzip” in Accept-Encoding.
func supportsGzip(req *types.Request) bool {
	if ae, ok := req.Headers["Accept-Encoding"]; ok {
		for _, enc := range strings.Split(ae, ",") {
			if strings.TrimSpace(enc) == "gzip" {
				return true
			}
		}
	}
	return false
}
