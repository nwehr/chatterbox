package repos

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nwehr/chatterbox"
)

func seedDatabase(db *sql.DB) error {
	contents, err := os.ReadFile("../../seed.sql")
	if err != nil {
		return fmt.Errorf("could not read seed file: %w", err)
	}

	if _, err = db.Exec(string(contents)); err != nil {
		return fmt.Errorf("could not execute seed file: %w", err)
	}

	return nil
}

// func TestListMessages(t *testing.T) {
// 	db, err := sql.Open("postgres", "host=localhost port=5432 user=chatterbox password=chatterbox dbname=chatterbox sslmode=disable")
// 	if err != nil {
// 		t.Errorf("could not connect to database: %s", err.Error())
// 	}

// 	defer db.Close()

// 	if err = seedDatabase(db); err != nil {
// 		t.Errorf("could not seed database: %s", err.Error())
// 	}

// 	r := NewMessageRepo(db)

// 	msgs, err := r.ListMessages(chatterbox.Identity("@nate.errorcode.io"), "a83087325cec029b2c39464aa0d2ea94f9282531cb53ad41143ea8bcea78205a")
// 	if err != nil {
// 		t.Errorf("coult not list messages: %s", err.Error())
// 	}

// 	if len(msgs) != 3 {
// 		t.Errorf("expected 3 messages; got %d", len(msgs))
// 	}
// }

func TestSaveMessage(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=chatterbox password=chatterbox dbname=chatterbox sslmode=disable")
	if err != nil {
		t.Errorf("could not connect to database: %s", err.Error())
	}

	defer db.Close()

	if err = seedDatabase(db); err != nil {
		t.Errorf("could not seed database: %s", err.Error())
	}

	r := NewMessageRepo(db)

	msg := chatterbox.Msg(chatterbox.Identity("@nate.errorcode.io"), chatterbox.Identities{"@nate.errorcode.io"}, "Just another note to myself")

	if err = r.SaveMessage(msg); err != nil {
		t.Errorf("could not save message: %s", err.Error())
	}

	msgs, err := r.ListMessages(chatterbox.Identity("@nate.errorcode.io"), "76359067c82a3f8e1dadc0e93570f2113c4f6b2b66994ce5e115d9be6d983de6")
	if err != nil {
		t.Errorf("could not list messages: %s", err.Error())
	}

	if len(msgs) != 1 {
		t.Errorf("expected 1 messages; got %d", len(msgs))
	}
}
