package core

type eventType int

const (
	sendMessageEventType eventType = iota
	subscribeEventType
	unsubscribeEventType
)

type event struct {
	Consumer *Consumer
	Type     eventType
}
