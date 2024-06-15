package entity

import (
	"sync"
	"time"
)

type Stats struct {
	Mu        sync.Mutex `bson:"-" json:"-"`
	Timestamp time.Time

	AllMessagesCount  int `bson:"all_messages_count"`
	AllCommandsCount  int `bson:"all_commands_count"`
	AllCallbacksCount int `bson:"all_callbacks_count"`

	LatestMessagesCount  int `bson:"latest_messages_count"`
	LatestCommandsCount  int `bson:"latest_commands_count"`
	LatestCallbacksCount int `bson:"latest_callbacks_count"`
}
