package main

import (
	"fmt"
	"os"

	pb "github.com/ahmet/historian/go-services/pkg/proto"
	"google.golang.org/protobuf/proto"
)

func main() {
	sensorData := &pb.SensorData{
		SensorId:    "sensor-go-1",
		Value:       99.9,
		TimestampMs: 1678886400000,
		Quality:     1,
	}

	fmt.Printf("Sensor ID: %s\n", sensorData.SensorId)

	// Verify serialization
	_, err := proto.Marshal(sensorData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Serialization failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Go Protobuf verification passed.")
}
