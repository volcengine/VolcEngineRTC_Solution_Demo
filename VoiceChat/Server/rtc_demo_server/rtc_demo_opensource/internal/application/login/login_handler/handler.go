package login_handler

var handler *EventHandler

type EventHandler struct {
}

func NewEventHandler() *EventHandler {
	if handler == nil {
		handler = &EventHandler{}
	}
	return handler
}
