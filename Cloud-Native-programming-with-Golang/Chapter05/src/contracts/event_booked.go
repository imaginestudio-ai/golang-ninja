package contracts

// EventCourseedEvent is emitted whenever an event is booked
type EventCourseedEvent struct {
	EventID string `json:"eventId"`
	UserID  string `json:"userId"`
}

// EventName returns the event's name
func (c *EventCourseedEvent) EventName() string {
	return "eventCourseed"
}
