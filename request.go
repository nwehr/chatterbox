package chatterbox

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Args map[string][]string

func (h Args) Set(key, value string) {
	h[key] = []string{value}
}

func (h Args) Add(key, value string) {
	h[key] = append(h[key], value)
}

type Request struct {
	Type string
	Args Args
	Data []byte
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
		data, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		req.Data = data[:len(data)-1]
	}

	return nil
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

	for key, values := range r.Args {
		line += fmt.Sprintf(" %s=%s", key, strings.Join(values, ";"))
	}

	return line
}

func (r *Request) parseMainLine(buf []byte) error {
	r.Args = Args{}

	key := ""
	value := ""

	readKey := true

	for _, ch := range buf {
		if ch == '\n' {
			break
		}

		if ch == ' ' {
			if r.Type == "" {
				r.Type = key
			} else {
				r.Args.Add(key, value)
			}

			key = ""
			value = ""
			readKey = true

			continue
		}

		if ch == '=' {
			readKey = false
			continue
		}

		if ch == ';' {
			r.Args.Add(key, value)
			value = ""
			continue
		}

		if readKey {
			key += string(ch)
		} else {
			value += string(ch)
		}

	}

	r.Args.Add(key, value)
	return nil
}

func parseMainLine(buf []byte) (string, Args, error) {
	rType := ""
	args := Args{}

	key := ""
	value := ""

	readKey := true

	for _, ch := range buf {
		if ch == '\n' {
			break
		}

		if ch == ' ' {
			if rType == "" {
				rType = key
			} else {
				args.Add(key, value)
			}

			key = ""
			value = ""
			readKey = true

			continue
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

	args.Add(key, value)

	fmt.Println("args", args)

	return rType, args, nil
}

func LoginRequest(ident Identity, password string) Request {
	return Request{
		Type: "LOGIN",
		Args: Args{
			"Identity": []string{string(ident)},
			"Password": []string{password},
		},
	}
}

func SendRequest(from Identity, to []Identity, msg string) Request {
	strTo := []string{}

	for _, ident := range to {
		strTo = append(strTo, string(ident))
	}

	return Request{
		Type: "SEND",
		Args: Args{
			"To":     strTo,
			"From":   []string{string(from)},
			"Length": []string{fmt.Sprintf("%d", len(msg))},
		},
		Data: []byte(msg),
	}
}
