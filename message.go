package chatterbox

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
	"github.com/google/uuid"
)

type Args map[string][]string

func (h Args) Set(key, value string) {
	h[key] = []string{value}
}

func (h Args) Add(key, value string) {
	h[key] = append(h[key], value)
}

func (h Args) First(key string) string {
	if len(h[key]) > 0 {
		return h[key][0]
	}

	return ""
}

type Message struct {
	Type string
	Args Args
	Data []byte
}

func (msg *Message) ReadFrom(r io.Reader) (int64, error) {
	reader := bufio.NewReader(r)

	mainLine, err := reader.ReadBytes('\n')
	if err != nil {
		return int64(len(mainLine)), err
	}

	err = msg.parseMainLine(mainLine)
	if err != nil {
		return int64(len(mainLine)), err
	}

	if len(msg.Args["Length"]) > 0 {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			return int64(len(mainLine) + len(data)), err
		}

		msg.Data = data[:len(data)-1]
	}

	return int64(len(mainLine) + len(msg.Data)), nil
}

func (msg Message) WriteTo(w io.Writer) (int64, error) {
	nMainLine, err := fmt.Fprintf(w, "%s\n", msg.mainLine())
	if err != nil {
		return int64(nMainLine), err
	}

	nDataLine := 0

	if len(msg.Data) > 0 {
		nDataLine, err = fmt.Fprintf(w, "%s\n\n", msg.Data)
	}

	return int64(nMainLine + nDataLine), err
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
			if msg.Type == "" {
				msg.Type = key
			} else {
				msg.Args.Add(key, value)
			}
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

	return nil
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

func Msg(from Identity, recipients Identities, msg string) Message {
	m := Message{
		Type: "MSG",
		Args: Args{
			"Recipients": recipients.Strings(),
			"From":       []string{string(from)},
			"Encoding":   []string{"text/plain"},
			"Uuid":       []string{uuid.NewString()},
		},
	}

	key, _ := age.ParseX25519Recipient("age1rmpjlh40vsmry47pad0h4u0lavtrm0nlypaya4adf7xy9n0rd5zqzpgkua")

	encoded := new(bytes.Buffer)

	encoder := base64.NewEncoder(base64.RawStdEncoding, encoded)
	encrypter, _ := age.Encrypt(encoder, key)
	io.WriteString(encrypter, msg)

	encrypter.Close()
	encoder.Close()

	m.Data = encoded.Bytes()
	m.Args["Length"] = []string{fmt.Sprintf("%d", len(m.Data))}

	return m
}

func LConv() Message {
	return Message{
		Type: "LCONV",
	}
}

func Conv(id string, recipients Identities, unread uint) Message {
	return Message{
		Type: "CONV",
		Args: Args{
			"Id":         []string{id},
			"Recipients": recipients.Strings(),
			"Unread":     []string{fmt.Sprintf("%d", unread)},
		},
	}
}

func Ok() Message {
	return Message{
		Type: "OK",
	}
}

type MessageRepo interface {
	ListMessages(recipient Identity, conversationID string) ([]Message, error)
	SaveMessage(Message) error
}
