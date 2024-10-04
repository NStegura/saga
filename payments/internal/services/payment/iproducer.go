package outbox

type Producer interface {
	PushPaymentEvent(int64, bool) error
}
