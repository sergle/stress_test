package app

import (
	"time"
)

const EventSenderTaskQueue = "EV_SENDER_TASK_QUEUE"

const EventSenderWorkflowType = "SingleSendEvent"

const EventSenderSignalName = "SIGNAL_1"

const EventStatusQuery = "QUERY_NUM_PROCESSED"

type EventDetails struct {
	TypeName  string    `json:"event_type"`
	EventID   int64     `json:"event_id"`
	UniqueID  string    `json:"unique_id"`
	Created   time.Time `json:"-"`
	Effective time.Time `json:"-"`
        CreatedNano int64   `json:"time_nano"`
	//Variables map[string]string `json:"variables"`
	Variables map[string]interface{} `json:"variables"`
}
