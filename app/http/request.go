package http

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/types"
)

func ParseRequest(r *bufio.Reader) (*types.Request, error) {
	// Request line
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil, fmt.Errorf("malformed request line: %q", line)
	}

	req := &types.Request{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
		Headers: make(map[string]string),
	}

	// Headers
	for {
		hdr, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		hdr = strings.TrimSpace(hdr)
		if hdr == "" {
			break
		}
		if kv := strings.SplitN(hdr, ":", 2); len(kv) == 2 {
			req.Headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	// Body (if any)
	if cl, ok := req.Headers["Content-Length"]; ok {
		n, err := strconv.Atoi(cl)
		if err != nil {
			return nil, fmt.Errorf("invalid Content-Length: %v", err)
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		req.Body = buf
	}

	return req, nil
}
