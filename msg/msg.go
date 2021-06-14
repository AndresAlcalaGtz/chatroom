package msg

const (
	Login  = "Joined the chat!"
	Logout = "Left the chat!"
)

type Message struct {
	Author string
	Text   string
	Bytes  []byte
	File   bool
}

func (msg *Message) ToString() string {
	return msg.Author + ": " + msg.Text
}
