package core

import "time"

type AlarmDefinition struct {
	ID        int       `json:"id"`
	Tag       string    `json:"tag"`
	Threshold float64   `json:"threshold"`
	Type      string    `json:"type"`     // High, Low
	Priority  string    `json:"priority"` // Critical, Warning
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ActiveAlarm struct {
	ID             int        `json:"id"`
	DefinitionID   int        `json:"definition_id"`
	State          string     `json:"state"`
	ActivationTime time.Time  `json:"activation_time"`
	AckTime        *time.Time `json:"ack_time,omitempty"`
	ShelvedUntil   *time.Time `json:"shelved_until,omitempty"`
	Value          float64    `json:"value"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type AlarmRepository interface {
	CreateDefinition(def *AlarmDefinition) error
	GetDefinition(id int) (*AlarmDefinition, error)
	ListDefinitions() ([]*AlarmDefinition, error)
	GetDefinitionsByTag(tag string) ([]*AlarmDefinition, error)

	CreateActiveAlarm(alarm *ActiveAlarm) error
	UpdateActiveAlarmState(id int, state string) error
	AckActiveAlarm(id int, ackTime time.Time) error
	ShelveActiveAlarm(id int, shelvedUntil time.Time) error
	GetActiveAlarms() ([]*ActiveAlarm, error)
}
