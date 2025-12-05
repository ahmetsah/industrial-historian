package transport

import (
	"log"

	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/core"
	pb "github.com/ahmetsah/industrial-historian/go-services/pkg/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type NatsTransport struct {
	conn    *nats.Conn
	service *core.AlarmService
}

func NewNatsTransport(url string) (*NatsTransport, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsTransport{conn: nc}, nil
}

func (t *NatsTransport) SetService(s *core.AlarmService) {
	t.service = s
}

func (t *NatsTransport) Start() error {
	_, err := t.conn.Subscribe("enterprise.>", func(msg *nats.Msg) {
		var sensorData pb.SensorData
		if err := proto.Unmarshal(msg.Data, &sensorData); err != nil {
			log.Printf("Failed to unmarshal sensor data: %v", err)
			return
		}
		if t.service != nil {
			if err := t.service.ProcessValue(sensorData.SensorId, sensorData.Value); err != nil {
				log.Printf("Failed to process value for %s: %v", sensorData.SensorId, err)
			}
		}
	})
	return err
}

func (t *NatsTransport) Close() {
	t.conn.Close()
}

// Implement EventPublisher
func (t *NatsTransport) PublishAlarmEvent(event *pb.AlarmEvent) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	return t.conn.Publish("sys.alarm.events", data)
}
