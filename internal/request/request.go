package request

import (
	"bytes"
	"fmt"
	"io"

	"github.com/KaranMali2001/Http-Server/internal/errors"
)

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}
type Request struct {
	RequestLine RequestLine
}

var SEPERATOR = []byte("\r\n")

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	req, msg, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}
	fmt.Println("Left Over Message ", string(msg))
	return &Request{
		RequestLine: *req,
	}, nil
}
func parseRequestLine(b []byte) (*RequestLine, []byte, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, errors.ErrStartLine
	}
	startLine := b[:idx]
	restOfLine := b[idx+len(SEPERATOR):]
	parts := bytes.Fields(startLine)
	if len(parts) != 3 {
		return nil, restOfLine, errors.ErrPartsMissingStartLine
	}
	method := parts[0]
	path := parts[1]
	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, restOfLine, errors.ErrHttpPartsMissing
	}
	return &RequestLine{
		RequestTarget: string(path),
		Method:        string(method),
		HttpVersion:   string(httpParts[1]),
	}, restOfLine, nil
}
