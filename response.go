package chatterbox

import (
	"bufio"
	"fmt"
	"io"
)

type Response struct {
	Type string
}

func (r Response) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s\n", r.Type)
	return err
}

func (r *Response) Read(rd io.Reader) error {
	line, err := bufio.NewReader(rd).ReadBytes('\n')
	if err != nil {
		return err
	}

	for _, ch := range line {
		if ch == ' ' || ch == '\n' {
			break
		}

		r.Type += string(ch)
	}

	return nil
}

func OKResponse() Response {
	return Response{
		Type: "OK",
	}
}
