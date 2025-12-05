package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/core"
	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/repository"
	pb "github.com/ahmetsah/industrial-historian/go-services/pkg/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

type AuditConsumer struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	repo   repository.Repository
	hasher core.Hasher
}

func NewAuditConsumer(natsUrl string, repo repository.Repository, hasher core.Hasher) (*AuditConsumer, error) {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &AuditConsumer{
		nc:     nc,
		js:     js,
		repo:   repo,
		hasher: hasher,
	}, nil
}

func (c *AuditConsumer) Start(ctx context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:     "AUDIT_EVENTS",
		Subjects: []string{"sys.audit.>", "sys.auth.login", "sys.alarm.>"},
		Storage:  jetstream.FileStorage,
	}

	_, err := c.js.CreateOrUpdateStream(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	cons, err := c.js.CreateOrUpdateConsumer(ctx, "AUDIT_EVENTS", jetstream.ConsumerConfig{
		Durable:   "AuditServiceConsumer",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		if err := c.processMessage(ctx, msg); err != nil {
			log.Printf("Failed to process message: %v", err)
			msg.Nak()
		} else {
			msg.Ack()
		}
	})
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	return nil
}

func (c *AuditConsumer) processMessage(ctx context.Context, msg jetstream.Msg) error {
	var actor, action string
	var detailsBytes []byte
	var timestamp time.Time

	if len(msg.Subject()) >= 9 && msg.Subject()[:9] == "sys.alarm" {
		// Handle Protobuf AlarmEvent
		var event pb.AlarmEvent
		if err := proto.Unmarshal(msg.Data(), &event); err != nil {
			return fmt.Errorf("invalid AlarmEvent proto: %w", err)
		}

		actor = "system" // Alarms are usually system generated or triggered by sensor
		action = fmt.Sprintf("alarm_%s", event.State)
		timestamp = time.UnixMilli(event.TimestampMs)

		details := map[string]interface{}{
			"alarm_id":      event.AlarmId,
			"definition_id": event.DefinitionId,
			"value":         event.Value,
			"message":       event.Message,
			"state":         event.State,
		}
		detailsBytes, _ = json.Marshal(details)

	} else {
		// Handle JSON (default)
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Data(), &payload); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		actor, _ = payload["actor"].(string)
		action, _ = payload["action"].(string)

		if action == "" {
			action = msg.Subject()
		}
		if actor == "" {
			actor = "system"
		}

		// Remove actor and action from payload to avoid duplication in details
		delete(payload, "actor")
		delete(payload, "action")

		detailsBytes, _ = json.Marshal(payload)
		timestamp = time.Now()
	}

	logEntry := &core.LogEntry{
		Timestamp: timestamp,
		Actor:     actor,
		Action:    action,
		Details:   json.RawMessage(detailsBytes),
	}

	return c.repo.AppendLog(ctx, logEntry, c.hasher)
}

func (c *AuditConsumer) Close() {
	c.nc.Close()
}
