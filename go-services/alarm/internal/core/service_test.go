package core

import (
	"testing"
	"time"

	pb "github.com/ahmetsah/industrial-historian/go-services/pkg/proto"
)

// MockRepo implements AlarmRepository for testing
type MockRepo struct {
	definitions  map[int]*AlarmDefinition
	activeAlarms map[int]*ActiveAlarm
	nextDefID    int
	nextAlarmID  int
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		definitions:  make(map[int]*AlarmDefinition),
		activeAlarms: make(map[int]*ActiveAlarm),
		nextDefID:    1,
		nextAlarmID:  1,
	}
}

func (m *MockRepo) CreateDefinition(def *AlarmDefinition) error {
	def.ID = m.nextDefID
	m.nextDefID++
	m.definitions[def.ID] = def
	return nil
}

func (m *MockRepo) GetDefinition(id int) (*AlarmDefinition, error) {
	return m.definitions[id], nil
}

func (m *MockRepo) ListDefinitions() ([]*AlarmDefinition, error) {
	var defs []*AlarmDefinition
	for _, d := range m.definitions {
		defs = append(defs, d)
	}
	return defs, nil
}

func (m *MockRepo) GetDefinitionsByTag(tag string) ([]*AlarmDefinition, error) {
	var defs []*AlarmDefinition
	for _, d := range m.definitions {
		if d.Tag == tag {
			defs = append(defs, d)
		}
	}
	return defs, nil
}

func (m *MockRepo) CreateActiveAlarm(alarm *ActiveAlarm) error {
	alarm.ID = m.nextAlarmID
	m.nextAlarmID++
	m.activeAlarms[alarm.ID] = alarm
	return nil
}

func (m *MockRepo) UpdateActiveAlarmState(id int, state string) error {
	if a, ok := m.activeAlarms[id]; ok {
		a.State = state
	}
	return nil
}

func (m *MockRepo) AckActiveAlarm(id int, ackTime time.Time) error {
	if a, ok := m.activeAlarms[id]; ok {
		a.State = "AckActive"
		a.AckTime = &ackTime
	}
	return nil
}

func (m *MockRepo) ShelveActiveAlarm(id int, shelvedUntil time.Time) error {
	if a, ok := m.activeAlarms[id]; ok {
		a.State = "Shelved"
		a.ShelvedUntil = &shelvedUntil
	}
	return nil
}

func (m *MockRepo) GetActiveAlarms() ([]*ActiveAlarm, error) {
	var alarms []*ActiveAlarm
	for _, a := range m.activeAlarms {
		if a.State != "Normal" {
			alarms = append(alarms, a)
		}
	}
	return alarms, nil
}

// MockPublisher implements EventPublisher for testing
type MockPublisher struct {
	events []*pb.AlarmEvent
}

func (m *MockPublisher) PublishAlarmEvent(event *pb.AlarmEvent) error {
	m.events = append(m.events, event)
	return nil
}

func TestAlarmService_ProcessValue(t *testing.T) {
	repo := NewMockRepo()
	publisher := &MockPublisher{}
	svc := NewAlarmService(repo, publisher)

	// Create Definition
	def := &AlarmDefinition{
		Tag:       "sensor1",
		Threshold: 100,
		Type:      "High",
		Priority:  "Critical",
	}
	repo.CreateDefinition(def)
	svc.LoadDefinitions()

	// 1. Normal -> Trigger -> UnackActive
	svc.ProcessValue("sensor1", 101)

	alarms := svc.GetActiveAlarms()
	if len(alarms) != 1 {
		t.Fatalf("Expected 1 active alarm, got %d", len(alarms))
	}
	if alarms[0].State != "UnackActive" {
		t.Errorf("Expected state UnackActive, got %s", alarms[0].State)
	}
	if len(publisher.events) != 1 {
		t.Errorf("Expected 1 event published, got %d", len(publisher.events))
	}

	// 2. UnackActive -> Ack -> AckActive
	alarmID := alarms[0].ID
	err := svc.Acknowledge(alarmID)
	if err != nil {
		t.Fatalf("Failed to acknowledge: %v", err)
	}

	alarms = svc.GetActiveAlarms()
	if alarms[0].State != "AckActive" {
		t.Errorf("Expected state AckActive, got %s", alarms[0].State)
	}
	if len(publisher.events) != 2 {
		t.Errorf("Expected 2 events published, got %d", len(publisher.events))
	}

	// 3. AckActive -> Clear -> Normal
	svc.ProcessValue("sensor1", 99)

	alarms = svc.GetActiveAlarms()
	if len(alarms) != 0 {
		t.Errorf("Expected 0 active alarms (Normal), got %d", len(alarms))
	}
	if len(publisher.events) != 3 {
		t.Errorf("Expected 3 events published, got %d", len(publisher.events))
	}
}

func TestAlarmService_Shelve(t *testing.T) {
	repo := NewMockRepo()
	publisher := &MockPublisher{}
	svc := NewAlarmService(repo, publisher)

	def := &AlarmDefinition{
		Tag:       "sensor1",
		Threshold: 100,
		Type:      "High",
	}
	repo.CreateDefinition(def)
	svc.LoadDefinitions()

	// Trigger
	svc.ProcessValue("sensor1", 101)
	alarms := svc.GetActiveAlarms()
	alarmID := alarms[0].ID

	// Shelve
	err := svc.Shelve(alarmID, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to shelve: %v", err)
	}

	alarms = svc.GetActiveAlarms()
	if alarms[0].State != "Shelved" {
		t.Errorf("Expected state Shelved, got %s", alarms[0].State)
	}
	if alarms[0].ShelvedUntil == nil {
		t.Error("Expected ShelvedUntil to be set")
	}
}
