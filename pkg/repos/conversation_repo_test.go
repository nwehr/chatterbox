package repos

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	"github.com/nwehr/chatterbox"
)

func TestListConversations(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=chatterbox password=chatterbox dbname=chatterbox sslmode=disable")
	if err != nil {
		t.Error("could not connect to database")
	}

	defer db.Close()

	if err = seedDatabase(db); err != nil {
		t.Errorf("could not seed database: %s", err.Error())
	}

	{
		r := NewMessageRepo(db)

		msg := chatterbox.Msg(chatterbox.Identity("@nate.errorcode.io"), chatterbox.Identities{"@nate.errorcode.io"}, "Just another note to myself")

		if err = r.SaveMessage(msg); err != nil {
			t.Errorf("could not save message: %s", err.Error())
		}
	}

	r := NewConversationRepo(db)

	convs, err := r.ListConversations(chatterbox.Identity("@nate.errorcode.io"))
	if err != nil {
		t.Errorf("could not list conversations: %s", err.Error())
	}

	if len(convs) != 1 {
		t.Errorf("expected 1 conversations; got %d", len(convs))
	}

	for _, conv := range convs {
		if conv.ID == "76359067c82a3f8e1dadc0e93570f2113c4f6b2b66994ce5e115d9be6d983de6" {
			if conv.Unread != 1 {
				t.Errorf("expected unread to be 1; got %d", conv.Unread)
			}
		}
	}
}
