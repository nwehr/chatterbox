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

type Message struct {
	Type string
	Args Args
	Data []byte
}

func (msg *Message) Read(r io.Reader) error {
	reader := bufio.NewReader(r)

	mainLine, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	err = msg.parseMainLine(mainLine)
	if err != nil {
		return err
	}

	if len(msg.Args["Length"]) > 0 {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		msg.Data = data[:len(data)-1]
	}

	return nil
}

func (msg Message) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", msg.mainLine())

	if len(msg.Data) > 0 {
		w.Write(msg.Data)
		w.Write([]byte{'\n'})
	}

	return nil
}

func (msg Message) mainLine() string {
	line := msg.Type

	for key, values := range msg.Args {
		line += fmt.Sprintf(" %s=%s", key, strings.Join(values, ";"))
	}

	return line
}

func (msg *Message) parseMainLine(buf []byte) error {
	msg.Args = Args{}

	key := ""
	value := ""

	readKey := true

	for _, ch := range buf {
		if ch == '\n' {
			break
		}

		if ch == ' ' {
			if msg.Type == "" {
				msg.Type = key
			} else {
				msg.Args.Add(key, value)
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
			msg.Args.Add(key, value)
			value = ""
			continue
		}

		if readKey {
			key += string(ch)
		} else {
			value += string(ch)
		}

	}

	msg.Args.Add(key, value)
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

func Login(ident Identity, password string) Message {
	return Message{
		Type: "LOGIN",
		Args: Args{
			"Identity": []string{string(ident)},
			"Password": []string{password},
		},
	}
}

func Send(from Identity, to []Identity, msg string) Message {
	strTo := []string{}

	for _, ident := range to {
		strTo = append(strTo, string(ident))
	}

	return Message{
		Type: "SEND",
		Args: Args{
			"To":     strTo,
			"From":   []string{string(from)},
			"Length": []string{fmt.Sprintf("%d", len(msg))},
		},
		Data: []byte(msg),
	}
}

func Ok() Message {
	return Message{
		Type: "OK",
	}
}
