package http

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/types" // added import
)

// Updated WriteResponse: include Content-Encoding header if set.
func WriteResponse(w *bufio.Writer, req *types.Request, status int, statusText, contentType string, body []byte) error {
	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", status, statusText)
	if req != nil {
		if conn, ok := req.Headers["Connection"]; ok && strings.ToLower(conn) == "close" {
			fmt.Fprintf(w, "Connection: close\r\n")
		}
		// Write Content-Encoding if the handler marked it.
		if encoding, ok := req.Headers["Response-Content-Encoding"]; ok {
			fmt.Fprintf(w, "Content-Encoding: %s\r\n", encoding)
		}
	}
	fmt.Fprintf(w, "Content-Type: %s\r\n", contentType)
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(w, "\r\n")
	if _, err := w.Write(body); err != nil {
		return err
	}
	return w.Flush()
}

func WritePlainTextResponse(w *bufio.Writer, req *types.Request, status int, statusText, body string) error {
	return WriteResponse(w, req, status, statusText, "text/plain", []byte(body))
}

func RespondBadRequest(w *bufio.Writer) {
	_ = WritePlainTextResponse(w, nil, 400, "Bad Request", "Malformed request")
}
