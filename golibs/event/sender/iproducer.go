package sender

type Producer interface {
	PushMsg(msg []byte) error
}
