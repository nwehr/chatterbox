package chatterbox

import (
	"bufio"
	"fmt"
	"io"
)

type Args map[string]string

func (h Args) Set(key, value string) {
	h[key] = value
}

type Request struct {
	Type string
	Args Args
	Data []byte
}

func (r Request) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", r.mainLine())

	if len(r.Data) > 0 {
		w.Write(r.Data)
		w.Write([]byte{'\n'})
	}

	return nil
}

func (r Request) mainLine() string {
	line := r.Type

	for key, value := range r.Args {
		line += fmt.Sprintf(" %s=%s", key, value)
	}

	return line
}

func (req *Request) Read(r io.Reader) error {
	reader := bufio.NewReader(r)

	mainLine, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	err = req.parseMainLine(mainLine)
	if err != nil {
		return err
	}

	if len(req.Args["Length"]) > 0 {
		req.Data, err = reader.ReadBytes('\n')
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Request) parseMainLine(buf []byte) error {
	r.Args = Args{}

	key := ""
	value := ""

	readKey := true

	for _, ch := range buf {
		if ch == '\n' {
			return nil
		}

		if ch == ' ' {
			if len(r.Type) == 0 {
				r.Type = key

				key = ""
				value = ""

				continue
			}

			r.Args.Set(key, value)

			key = ""
			value = ""
			readKey = true

		}

		if ch == '=' {
			readKey = false
			continue
		}

		if readKey {
			key += string(ch)
		} else {
			value += string(ch)
		}

	}

	return nil
}

func LoginRequest(ident Identity, password string) Request {
	return Request{
		Type: "LOGIN",
		Args: Args{
			"Identity": string(ident),
			"Password": password,
		},
	}
}
