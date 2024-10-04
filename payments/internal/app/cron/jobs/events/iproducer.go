package events

type Producer interface {
	PushMsg(msg []byte) error
}
