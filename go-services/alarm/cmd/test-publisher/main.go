package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/ahmetsah/industrial-historian/go-services/pkg/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func main() {
	sensorID := flag.String("sensor", "sensor1", "Sensor ID")
	value := flag.Float64("value", 0.0, "Sensor Value")
	natsURL := flag.String("nats", "nats://localhost:4222", "NATS URL")
	flag.Parse()

	nc, err := nats.Connect(*natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	data := &pb.SensorData{
		SensorId:    *sensorID,
		Value:       *value,
		TimestampMs: time.Now().UnixMilli(),
		Quality:     1,
	}

	bytes, err := proto.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	subject := fmt.Sprintf("enterprise.site1.area1.line1.device1.%s", *sensorID)
	if err := nc.Publish(subject, bytes); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	fmt.Printf("Published to %s: %v\n", subject, data)
}
