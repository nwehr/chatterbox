package repos

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/nwehr/chatterbox"
)

type conversationRepo struct {
	db *sql.DB
}

func NewConversationRepo(db *sql.DB) chatterbox.ConversationRepo {
	return conversationRepo{
		db: db,
	}
}

func (r conversationRepo) ListConversations(recipient chatterbox.Identity) ([]chatterbox.Conversation, error) {
	query := `select 
		m.conversation_id
		, m.recipients 
		, ( select 
				count(r.*) 
			from message_recipients r 
			where 
				r.conversation_id = m.conversation_id 
				and r.recipient = $1 
				and r.read_at is null
		) as unread
	from messages m
	where 
		$1 = any(m.recipients)
	group by 
		m.conversation_id
		, m.recipients;`

	rows, err := r.db.Query(query, recipient.String())
	if err != nil {
		return nil, fmt.Errorf("could not query conversations: %w", err)
	}

	defer rows.Close()

	convs := []chatterbox.Conversation{}

	for rows.Next() {
		conv := chatterbox.Conversation{}
		recipients := pq.StringArray{}

		err = rows.Scan(&conv.ID, &recipients, &conv.Unread)
		if err != nil {
			return convs, fmt.Errorf("could not scan into conversation: %w", err)
		}

		conv.Recipients = chatterbox.NewIdentities(recipients)

		convs = append(convs, conv)
	}

	return convs, nil
}
