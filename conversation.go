package chatterbox

type Conversation struct {
	ID         string
	Recipients Identities
	Unread     uint
}

type ConversationRepo interface {
	ListConversations(recipient Identity) ([]Conversation, error)
}
