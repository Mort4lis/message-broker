package core

type eventType int

const (
	sendMessageEventType eventType = iota
	subscribeEventType
	unsubscribeEventType
)

type event struct {
	Type     eventType
	Consumer *Consumer
}
