package core

type eventType int

const (
	sendMessageEventType eventType = iota
	subscribeEventType
	unsubscribeEventType
)

type event struct {
	Type eventType
	Meta any
}

type subscriberMeta struct {
	Subscriber *Subscriber
}

func newSendMessageEvent() event {
	return event{Type: sendMessageEventType}
}

func newSubscribeEvent(sub *Subscriber) event {
	return event{
		Type: subscribeEventType,
		Meta: subscriberMeta{Subscriber: sub},
	}
}

func newUnsubscribeEvent(sub *Subscriber) event {
	return event{
		Type: unsubscribeEventType,
		Meta: subscriberMeta{Subscriber: sub},
	}
}
