package core

import (
	"testing"
	"time"
)

func TestFSM_Transitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  AlarmState
		event         AlarmEvent
		expectedState AlarmState
	}{
		{"Normal -> Trigger -> UnackActive", StateNormal, EventTrigger, StateUnackActive},
		{"UnackActive -> Ack -> AckActive", StateUnackActive, EventAck, StateAckActive},
		{"AckActive -> Clear -> Normal", StateAckActive, EventClear, StateNormal},
		{"UnackActive -> Clear -> UnackRTN", StateUnackActive, EventClear, StateUnackRTN},
		{"UnackRTN -> Ack -> Normal", StateUnackRTN, EventAck, StateNormal},
		{"UnackRTN -> Trigger -> UnackActive", StateUnackRTN, EventTrigger, StateUnackActive},
		// Shelving logic might be separate or part of FSM
		{"Normal -> Shelve -> Shelved", StateNormal, EventShelve, StateShelved},
		{"Shelved -> Unshelve -> Normal", StateShelved, EventUnshelve, StateNormal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := NewAlarmFSM(tt.initialState)
			newState, err := fsm.Transition(tt.event)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if newState != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, newState)
			}
		})
	}
}

func TestFSM_InvalidTransitions(t *testing.T) {
	fsm := NewAlarmFSM(StateNormal)
	_, err := fsm.Transition(EventAck) // Cannot ack Normal
	if err == nil {
		t.Error("Expected error for invalid transition Normal -> Ack")
	}
}

func TestFSM_Shelving(t *testing.T) {
	fsm := NewAlarmFSM(StateNormal)
	// Shelve for 1 hour
	err := fsm.Shelve(1 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to shelve: %v", err)
	}
	if fsm.State != StateShelved {
		t.Errorf("Expected Shelved, got %v", fsm.State)
	}

	// Check if shelved
	if !fsm.IsShelved() {
		t.Error("Expected IsShelved to be true")
	}
}
