package core

import (
	"fmt"
	"time"
)

type AlarmState string

const (
	StateNormal      AlarmState = "Normal"
	StateUnackActive AlarmState = "UnackActive"
	StateAckActive   AlarmState = "AckActive"
	StateUnackRTN    AlarmState = "UnackRTN"
	StateShelved     AlarmState = "Shelved"
	StateSuppressed  AlarmState = "Suppressed"
)

type AlarmEvent string

const (
	EventTrigger  AlarmEvent = "Trigger"
	EventClear    AlarmEvent = "Clear"
	EventAck      AlarmEvent = "Ack"
	EventShelve   AlarmEvent = "Shelve"
	EventUnshelve AlarmEvent = "Unshelve"
)

type AlarmFSM struct {
	State        AlarmState
	ShelvedUntil time.Time
}

func NewAlarmFSM(initialState AlarmState) *AlarmFSM {
	return &AlarmFSM{
		State: initialState,
	}
}

func (fsm *AlarmFSM) Transition(event AlarmEvent) (AlarmState, error) {
	switch fsm.State {
	case StateNormal:
		switch event {
		case EventTrigger:
			fsm.State = StateUnackActive
		case EventShelve:
			fsm.State = StateShelved
		default:
			return fsm.State, fmt.Errorf("invalid transition from %s with event %s", fsm.State, event)
		}

	case StateUnackActive:
		switch event {
		case EventAck:
			fsm.State = StateAckActive
		case EventClear:
			fsm.State = StateUnackRTN
		case EventShelve:
			fsm.State = StateShelved
		default:
			return fsm.State, fmt.Errorf("invalid transition from %s with event %s", fsm.State, event)
		}

	case StateAckActive:
		switch event {
		case EventClear:
			fsm.State = StateNormal
		case EventShelve:
			fsm.State = StateShelved
		default:
			return fsm.State, fmt.Errorf("invalid transition from %s with event %s", fsm.State, event)
		}

	case StateUnackRTN:
		switch event {
		case EventAck:
			fsm.State = StateNormal
		case EventTrigger:
			fsm.State = StateUnackActive
		case EventShelve:
			fsm.State = StateShelved
		default:
			return fsm.State, fmt.Errorf("invalid transition from %s with event %s", fsm.State, event)
		}

	case StateShelved:
		switch event {
		case EventUnshelve:
			fsm.State = StateNormal
			fsm.ShelvedUntil = time.Time{}
		default:
			return fsm.State, fmt.Errorf("invalid transition from %s with event %s", fsm.State, event)
		}

	default:
		// Fallback for unknown states, e.g. Suppressed
		if event == EventShelve {
			fsm.State = StateShelved
			return fsm.State, nil
		}
		return fsm.State, fmt.Errorf("invalid transition from %s with event %s", fsm.State, event)
	}

	return fsm.State, nil
}

func (fsm *AlarmFSM) Shelve(duration time.Duration) error {
	if fsm.State == StateSuppressed {
		return fmt.Errorf("cannot shelve suppressed alarm")
	}
	fsm.State = StateShelved
	fsm.ShelvedUntil = time.Now().Add(duration)
	return nil
}

func (fsm *AlarmFSM) IsShelved() bool {
	if fsm.State != StateShelved {
		return false
	}
	if !fsm.ShelvedUntil.IsZero() && time.Now().After(fsm.ShelvedUntil) {
		return false
	}
	return true
}
