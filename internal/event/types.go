package event

import "time"

type Category string

const (
	SessionStart   Category = "session.start"
	TaskAcknowledge Category = "task.acknowledge"
	TaskComplete   Category = "task.complete"
	TaskError      Category = "task.error"
	InputRequired  Category = "input.required"
	ResourceLimit  Category = "resource.limit"
	UserSpam       Category = "user.spam"
	None           Category = ""
)

type HookEvent struct {
	EventName        string `json:"hook_event_name"`
	SessionID        string `json:"session_id"`
	NotificationType string `json:"notification_type,omitempty"`
	ToolName         string `json:"tool_name,omitempty"`
	Error            string `json:"error,omitempty"`
	Source           string `json:"source,omitempty"`
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

type FixedClock struct {
	T time.Time
}

func (c FixedClock) Now() time.Time { return c.T }

func (c *FixedClock) Advance(d time.Duration) {
	c.T = c.T.Add(d)
}
