package core

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/ahmetsah/industrial-historian/go-services/pkg/proto"
)

type EventPublisher interface {
	PublishAlarmEvent(event *pb.AlarmEvent) error
}

type AlarmService struct {
	repo         AlarmRepository
	publisher    EventPublisher
	definitions  map[string][]*AlarmDefinition
	activeAlarms map[int]*ActiveAlarm
	mu           sync.RWMutex
}

func NewAlarmService(repo AlarmRepository, publisher EventPublisher) *AlarmService {
	return &AlarmService{
		repo:         repo,
		publisher:    publisher,
		definitions:  make(map[string][]*AlarmDefinition),
		activeAlarms: make(map[int]*ActiveAlarm),
	}
}

func (s *AlarmService) StartBackgroundTasks() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.checkShelvedAlarms()
		}
	}()
}

func (s *AlarmService) checkShelvedAlarms() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for _, active := range s.activeAlarms {
		if active.State == string(StateShelved) && active.ShelvedUntil != nil && now.After(*active.ShelvedUntil) {
			// Unshelve
			log.Printf("Unshelving alarm %d", active.ID)

			// Transition to Normal first, then let next evaluation decide if it should be Active
			// In a real system, we might want to check the current value immediately,
			// but resetting to Normal and letting the next value update trigger logic is safer/simpler.
			// Or better: Transition to UnackActive if we don't know the value, or Normal?
			// ISA 18.2 says unshelving returns to the state it would be in.
			// Since we don't have the *current* value right here without looking it up or waiting,
			// we'll transition to Normal and let the next ProcessValue update it.
			// However, if the condition is still present, it should go to UnackActive.

			// For now, we just update state to Normal so it's "active" again in the system's eyes
			// and will be re-evaluated on next sensor update.

			newState := StateNormal
			active.State = string(newState)
			active.ShelvedUntil = nil
			active.UpdatedAt = now

			if err := s.repo.UpdateActiveAlarmState(active.ID, string(newState)); err != nil {
				log.Printf("Failed to update unshelved alarm: %v", err)
				continue
			}

			// If it went to normal, remove from active map?
			// Wait, if we just set it to Normal, it's effectively "cleared" until next trigger.
			// But we need to keep it in map if we want to track it?
			// Actually, our logic says "if newState == StateNormal { delete }".
			// So let's do that.

			// But wait, if we delete it, and the value is still High, it will re-trigger as a NEW alarm (UnackActive)
			// on the next reading. This is acceptable behavior for unshelving.

			delete(s.activeAlarms, active.DefinitionID)

			// Publish Event
			if s.publisher != nil {
				eventPayload := &pb.AlarmEvent{
					AlarmId:      int32(active.ID),
					DefinitionId: int32(active.DefinitionID),
					State:        string(newState),
					Value:        active.Value,
					TimestampMs:  now.UnixMilli(),
					Message:      "Alarm unshelved (expired)",
				}
				s.publisher.PublishAlarmEvent(eventPayload)
			}
		}
	}
}

func (s *AlarmService) LoadDefinitions() error {
	defs, err := s.repo.ListDefinitions()
	if err != nil {
		return err
	}

	// Also load active alarms
	active, err := s.repo.GetActiveAlarms()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.definitions = make(map[string][]*AlarmDefinition)
	for _, def := range defs {
		s.definitions[def.Tag] = append(s.definitions[def.Tag], def)
	}

	s.activeAlarms = make(map[int]*ActiveAlarm)
	for _, a := range active {
		s.activeAlarms[a.DefinitionID] = a
	}

	return nil
}

func (s *AlarmService) ProcessValue(sensorId string, value float64) error {
	s.mu.RLock()
	defs, ok := s.definitions[sensorId]
	s.mu.RUnlock()

	if !ok || len(defs) == 0 {
		return nil // No definitions for this tag
	}

	var errs []error
	for _, def := range defs {
		if err := s.evaluateDefinition(def, value); err != nil {
			log.Printf("Error evaluating definition %d: %v", def.ID, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during processing", len(errs))
	}
	return nil
}

func (s *AlarmService) evaluateDefinition(def *AlarmDefinition, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	active, exists := s.activeAlarms[def.ID]

	currentState := StateNormal
	if exists {
		currentState = AlarmState(active.State)
	}

	fsm := NewAlarmFSM(currentState)

	shouldFire := Evaluate(def, value)

	var event AlarmEvent
	if shouldFire {
		event = EventTrigger
	} else {
		event = EventClear
	}

	newState, err := fsm.Transition(event)
	if err != nil {
		return nil
	}

	if newState != currentState {
		// State changed!
		var alarmID int
		if !exists {
			// Create new active alarm
			newAlarm := &ActiveAlarm{
				DefinitionID:   def.ID,
				State:          string(newState),
				ActivationTime: time.Now(),
				Value:          value,
			}
			if err := s.repo.CreateActiveAlarm(newAlarm); err != nil {
				return err
			}
			s.activeAlarms[def.ID] = newAlarm
			alarmID = newAlarm.ID
		} else {
			// Update existing
			active.State = string(newState)
			active.Value = value
			active.UpdatedAt = time.Now()
			if err := s.repo.UpdateActiveAlarmState(active.ID, string(newState)); err != nil {
				return err
			}
			alarmID = active.ID

			if newState == StateNormal {
				delete(s.activeAlarms, def.ID)
			}
		}

		// Publish Event
		if s.publisher != nil {
			eventPayload := &pb.AlarmEvent{
				AlarmId:      int32(alarmID),
				DefinitionId: int32(def.ID),
				State:        string(newState),
				Value:        value,
				TimestampMs:  time.Now().UnixMilli(),
				Message:      fmt.Sprintf("Alarm %s transitioned to %s", def.Tag, newState),
			}
			if err := s.publisher.PublishAlarmEvent(eventPayload); err != nil {
				log.Printf("Failed to publish alarm event: %v", err)
			}
		}
	}

	return nil
}

func (s *AlarmService) Acknowledge(alarmID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the alarm in active alarms map
	var active *ActiveAlarm
	var defID int
	found := false
	for dID, a := range s.activeAlarms {
		if a.ID == alarmID {
			active = a
			defID = dID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("alarm not found or not active")
	}

	currentState := AlarmState(active.State)
	fsm := NewAlarmFSM(currentState)

	newState, err := fsm.Transition(EventAck)
	if err != nil {
		return fmt.Errorf("cannot acknowledge alarm in state %s: %w", currentState, err)
	}

	if newState != currentState {
		active.State = string(newState)
		now := time.Now()
		active.AckTime = &now
		active.UpdatedAt = now

		if newState == StateNormal {
			if err := s.repo.UpdateActiveAlarmState(alarmID, string(newState)); err != nil {
				return err
			}
			delete(s.activeAlarms, defID)
		} else {
			if newState == StateAckActive {
				if err := s.repo.AckActiveAlarm(alarmID, now); err != nil {
					return err
				}
			} else {
				if err := s.repo.UpdateActiveAlarmState(alarmID, string(newState)); err != nil {
					return err
				}
			}
		}

		// Publish Event
		if s.publisher != nil {
			eventPayload := &pb.AlarmEvent{
				AlarmId:      int32(alarmID),
				DefinitionId: int32(defID),
				State:        string(newState),
				Value:        active.Value,
				TimestampMs:  time.Now().UnixMilli(),
				Message:      "Alarm acknowledged",
			}
			s.publisher.PublishAlarmEvent(eventPayload)
		}
	}

	return nil
}

func (s *AlarmService) Shelve(alarmID int, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var active *ActiveAlarm
	var defID int
	found := false
	for dID, a := range s.activeAlarms {
		if a.ID == alarmID {
			active = a
			defID = dID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("alarm not found or not active")
	}

	currentState := AlarmState(active.State)
	fsm := NewAlarmFSM(currentState)

	newState, err := fsm.Transition(EventShelve)
	if err != nil {
		return fmt.Errorf("cannot shelve alarm in state %s: %w", currentState, err)
	}

	shelvedUntil := time.Now().Add(duration)

	active.State = string(newState)
	active.ShelvedUntil = &shelvedUntil
	active.UpdatedAt = time.Now()

	if err := s.repo.ShelveActiveAlarm(alarmID, shelvedUntil); err != nil {
		return err
	}

	// Publish Event
	if s.publisher != nil {
		eventPayload := &pb.AlarmEvent{
			AlarmId:      int32(alarmID),
			DefinitionId: int32(defID),
			State:        string(newState),
			Value:        active.Value,
			TimestampMs:  time.Now().UnixMilli(),
			Message:      fmt.Sprintf("Alarm shelved until %s", shelvedUntil),
		}
		s.publisher.PublishAlarmEvent(eventPayload)
	}

	return nil
}

func (s *AlarmService) CreateDefinition(def *AlarmDefinition) error {
	if err := s.repo.CreateDefinition(def); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.definitions[def.Tag] = append(s.definitions[def.Tag], def)
	return nil
}

func (s *AlarmService) GetActiveAlarms() []*ActiveAlarm {
	s.mu.RLock()
	defer s.mu.RUnlock()

	alarms := make([]*ActiveAlarm, 0, len(s.activeAlarms))
	for _, a := range s.activeAlarms {
		alarms = append(alarms, a)
	}
	return alarms
}
