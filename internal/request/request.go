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
	State       parsedState
}
type parsedState string

const (
	StateInit  parsedState = "INIT"
	StateDone  parsedState = "DONE"
	StateError parsedState = "ERROR"
)

var SEPERATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		State: StateInit}
}
func (r *Request) done() bool {
	return r.State == StateDone
}
func (r *Request) error() bool {
	return r.State == StateError
}
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateError:
			return 0, errors.ErrParse
		case StateInit:
			r1, n, err := parseRequestLine(data[read:])

			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *r1
			read += n
			r.State = StateDone
		case StateDone:
			break outer
		}

	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() && !request.error() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {

			return nil, err
		}
		bufLen += n
		read, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[read:bufLen])
		bufLen -= read
	}
	// data, err := io.ReadAll(reader)
	// if err != nil {
	// 	return nil, err
	// }
	// req, msg, err := parseRequestLine(data)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("Left Over Message ", string(msg))
	return request, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {

		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPERATOR)
	parts := bytes.Fields(startLine)
	if len(parts) != 3 {
		return nil, read, errors.ErrPartsMissingStartLine
	}
	method := parts[0]
	path := parts[1]
	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		fmt.Println("Byts", string(b))
		return nil, read, errors.ErrHttpPartsMissing
	}
	return &RequestLine{
		RequestTarget: string(path),
		Method:        string(method),
		HttpVersion:   string(httpParts[1]),
	}, read, nil
}
