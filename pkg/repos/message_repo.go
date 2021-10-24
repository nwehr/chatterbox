package repos

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/nwehr/chatterbox"
)

type messageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) chatterbox.MessageRepo {
	return messageRepo{
		db: db,
	}
}

func (r messageRepo) ListMessages(recipient chatterbox.Identity, conversationID string) ([]chatterbox.Message, error) {
	query := `select 
		m.uuid
		, m.conversation_id
		, m.recipients
		, m.from
		, m.encoding
		, m.length
		, m.data
	from message_recipients r

	left join messages m
		on m.uuid = r.message_uuid

	where 
		r.recipient = $1
		and r.conversation_id = $2
		
	order by m.created_at desc;`

	rows, err := r.db.Query(query, recipient.String(), conversationID)
	if err != nil {
		return nil, fmt.Errorf("could not query messages: %w", err)
	}

	defer rows.Close()

	msgs := []chatterbox.Message{}

	for rows.Next() {
		uuid := ""
		conversationID := ""
		recipients := pq.StringArray{}
		from := ""
		encoding := ""
		length := 0
		data := ""

		err = rows.Scan(&uuid, &conversationID, &recipients, &from, &encoding, &length, &data)
		if err != nil {
			return msgs, fmt.Errorf("could not scan into msg: %w", err)
		}

		msg := chatterbox.Message{
			Type: "MSG",
			Args: chatterbox.Args{
				"Uuid":           []string{uuid},
				"ConversationId": []string{conversationID},
				"Recipients":     recipients,
				"From":           []string{from},
				"Encoding":       []string{encoding},
				"Length":         []string{fmt.Sprintf("%d", length)},
			},
			Data: []byte(data),
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (r messageRepo) SaveMessage(msg chatterbox.Message) error {
	insert := `insert into messages ("uuid", "conversation_id", "recipients", "from", "encoding", "length", "data") values ($1, $2, $3, $4, $5, $6, $7);`

	convId := chatterbox.NewIdentities(msg.Args["Recipients"]).ConversationID()

	_, err := r.db.Exec(insert,
		msg.Args.First("Uuid"),
		convId,
		pq.StringArray(msg.Args["Recipients"]),
		msg.Args.First("From"),
		msg.Args.First("Encoding"),
		msg.Args.First("Length"),
		msg.Data,
	)
	if err != nil {
		return fmt.Errorf("could not insert into messages: %w", err)
	}

	for _, recipient := range msg.Args["Recipients"] {
		if chatterbox.Identity(recipient).Host() == "errorcode.io" {
			insert = `insert into message_recipients (message_uuid, conversation_id, recipient) values ($1, $2, $3);`

			if _, err = r.db.Exec(insert, msg.Args.First("Uuid"), convId, recipient); err != nil {
				return fmt.Errorf("could not insert into message_recipients: %w", err)
			}
		}
	}

	return nil
}
