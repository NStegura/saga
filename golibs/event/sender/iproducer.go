package sender

type Producer interface {
	PushMsg(msg []byte, topic string) error
}
