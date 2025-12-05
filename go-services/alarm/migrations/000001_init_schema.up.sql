CREATE TABLE alarm_definitions (
    id SERIAL PRIMARY KEY,
    tag VARCHAR(255) NOT NULL,
    threshold DOUBLE PRECISION NOT NULL,
    alarm_type VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_alarm_definitions_tag_type ON alarm_definitions(tag, alarm_type);

CREATE TABLE active_alarms (
    id SERIAL PRIMARY KEY,
    definition_id INTEGER NOT NULL REFERENCES alarm_definitions(id),
    state VARCHAR(50) NOT NULL,
    activation_time TIMESTAMP WITH TIME ZONE NOT NULL,
    ack_time TIMESTAMP WITH TIME ZONE,
    shelved_until TIMESTAMP WITH TIME ZONE,
    value DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_active_alarms_state ON active_alarms(state);
