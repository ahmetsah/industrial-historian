-- Industrial Historian - Config Management Database Schema
-- Version: 2.0
-- Description: Multi-protocol device configuration management

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Protocol types
CREATE TYPE protocol_type AS ENUM ('modbus', 'opc', 's7', 'mqtt');
CREATE TYPE device_status AS ENUM ('active', 'inactive', 'error', 'pending');
CREATE TYPE config_status AS ENUM ('pending', 'generated', 'deployed', 'failed');

-- ============================================================================
-- MAIN TABLES
-- ============================================================================

-- Main devices table
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    protocol protocol_type NOT NULL,
    status device_status DEFAULT 'inactive',
    enabled BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255)
);

-- ============================================================================
-- MODBUS PROTOCOL
-- ============================================================================

CREATE TABLE modbus_devices (
    id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    ip VARCHAR(45) NOT NULL,
    port INTEGER DEFAULT 502 CHECK (port > 0 AND port < 65536),
    unit_id INTEGER NOT NULL CHECK (unit_id >= 0 AND unit_id < 256),
    poll_interval_ms INTEGER DEFAULT 1000 CHECK (poll_interval_ms > 0),
    timeout_ms INTEGER DEFAULT 5000 CHECK (timeout_ms > 0),
    retry_count INTEGER DEFAULT 3 CHECK (retry_count >= 0),
    connection_type VARCHAR(20) DEFAULT 'TCP' CHECK (connection_type IN ('TCP', 'RTU')),
    UNIQUE(ip, port, unit_id)
);

CREATE TABLE modbus_registers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID REFERENCES modbus_devices(id) ON DELETE CASCADE,
    address INTEGER NOT NULL CHECK (address >= 0 AND address < 65536),
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL CHECK (data_type IN ('Int16', 'UInt16', 'Int32', 'UInt32', 'Float32', 'Float64')),
    register_type VARCHAR(20) DEFAULT 'holding' CHECK (register_type IN ('holding', 'input', 'coil', 'discrete')),
    scale_factor FLOAT DEFAULT 1.0,
    "offset" FLOAT DEFAULT 0.0,
    unit VARCHAR(50),
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id, address, name)
);

-- ============================================================================
-- OPC UA PROTOCOL
-- ============================================================================

CREATE TABLE opc_devices (
    id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    endpoint_url VARCHAR(512) NOT NULL UNIQUE,
    security_mode VARCHAR(50) DEFAULT 'None' CHECK (security_mode IN ('None', 'Sign', 'SignAndEncrypt')),
    security_policy VARCHAR(50) DEFAULT 'None',
    username VARCHAR(255),
    password_encrypted TEXT,
    certificate_path VARCHAR(512),
    poll_interval_ms INTEGER DEFAULT 1000 CHECK (poll_interval_ms > 0),
    timeout_ms INTEGER DEFAULT 5000 CHECK (timeout_ms > 0)
);

CREATE TABLE opc_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID REFERENCES opc_devices(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50),
    namespace_index INTEGER DEFAULT 0,
    unit VARCHAR(50),
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id, node_id)
);

-- ============================================================================
-- SIEMENS S7 PROTOCOL
-- ============================================================================

CREATE TABLE s7_devices (
    id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    ip VARCHAR(45) NOT NULL UNIQUE,
    rack INTEGER DEFAULT 0 CHECK (rack >= 0),
    slot INTEGER DEFAULT 1 CHECK (slot >= 0),
    poll_interval_ms INTEGER DEFAULT 1000 CHECK (poll_interval_ms > 0),
    timeout_ms INTEGER DEFAULT 5000 CHECK (timeout_ms > 0),
    plc_type VARCHAR(50) DEFAULT 'S7-1200' CHECK (plc_type IN ('S7-200', 'S7-300', 'S7-400', 'S7-1200', 'S7-1500'))
);

CREATE TABLE s7_data_blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID REFERENCES s7_devices(id) ON DELETE CASCADE,
    db_number INTEGER NOT NULL CHECK (db_number >= 0),
    start_address INTEGER NOT NULL CHECK (start_address >= 0),
    length INTEGER NOT NULL CHECK (length > 0),
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL CHECK (data_type IN ('BOOL', 'BYTE', 'INT', 'DINT', 'REAL', 'STRING')),
    unit VARCHAR(50),
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id, db_number, start_address)
);

-- ============================================================================
-- CONFIG GENERATION & DEPLOYMENT
-- ============================================================================

CREATE TABLE config_generations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    config_hash VARCHAR(64) NOT NULL,
    config_content TEXT NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    generated_at TIMESTAMP DEFAULT NOW(),
    deployed_at TIMESTAMP,
    status config_status DEFAULT 'pending',
    error_message TEXT,
    generated_by VARCHAR(255)
);

CREATE TABLE deployment_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    config_generation_id UUID REFERENCES config_generations(id),
    deployed_at TIMESTAMP DEFAULT NOW(),
    deployed_by VARCHAR(255),
    instance_id VARCHAR(255),
    success BOOLEAN DEFAULT false,
    error_message TEXT
);

-- ============================================================================
-- TEMPLATES
-- ============================================================================

CREATE TABLE config_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    protocol protocol_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_content TEXT NOT NULL,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    UNIQUE(protocol, name)
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Devices
CREATE INDEX idx_devices_protocol ON devices(protocol);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_enabled ON devices(enabled);
CREATE INDEX idx_devices_created_at ON devices(created_at);

-- Modbus
CREATE INDEX idx_modbus_ip_port ON modbus_devices(ip, port);
CREATE INDEX idx_modbus_registers_device ON modbus_registers(device_id);
CREATE INDEX idx_modbus_registers_enabled ON modbus_registers(enabled);

-- OPC
CREATE INDEX idx_opc_endpoint ON opc_devices(endpoint_url);
CREATE INDEX idx_opc_nodes_device ON opc_nodes(device_id);
CREATE INDEX idx_opc_nodes_enabled ON opc_nodes(enabled);

-- S7
CREATE INDEX idx_s7_ip ON s7_devices(ip);
CREATE INDEX idx_s7_blocks_device ON s7_data_blocks(device_id);
CREATE INDEX idx_s7_blocks_enabled ON s7_data_blocks(enabled);

-- Config
CREATE INDEX idx_config_gen_device ON config_generations(device_id);
CREATE INDEX idx_config_gen_status ON config_generations(status);
CREATE INDEX idx_config_gen_generated_at ON config_generations(generated_at);
CREATE INDEX idx_deployment_device ON deployment_history(device_id);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_devices_updated_at BEFORE UPDATE ON devices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA
-- ============================================================================

-- Default Modbus template
INSERT INTO config_templates (protocol, name, description, template_content, is_default) VALUES
('modbus', 'Standard Modbus TCP', 'Default template for Modbus TCP devices', 
'[[modbus_devices]]
ip = "{{.IP}}"
port = {{.Port}}
unit_id = {{.UnitID}}
poll_interval_ms = {{.PollInterval}}

{{range .Registers}}
[[modbus_devices.registers]]
address = {{.Address}}
name = "{{.Name}}"
data_type = "{{.DataType}}"
{{end}}
', true);

-- Default OPC UA template
INSERT INTO config_templates (protocol, name, description, template_content, is_default) VALUES
('opc', 'Standard OPC UA', 'Default template for OPC UA devices',
'[[opc_devices]]
endpoint_url = "{{.EndpointURL}}"
security_mode = "{{.SecurityMode}}"
poll_interval_ms = {{.PollInterval}}

{{range .Nodes}}
[[opc_devices.nodes]]
node_id = "{{.NodeID}}"
name = "{{.Name}}"
{{end}}
', true);

-- Example device (for testing)
INSERT INTO devices (name, description, protocol, status, enabled) VALUES
('PLC-001', 'Main production line PLC', 'modbus', 'inactive', true);

INSERT INTO modbus_devices (id, ip, port, unit_id, poll_interval_ms) 
SELECT id, '192.168.1.10', 502, 1, 1000 FROM devices WHERE name = 'PLC-001';

INSERT INTO modbus_registers (device_id, address, name, data_type, unit, description)
SELECT id, 0, 'Temperature', 'Float32', '°C', 'Reactor temperature'
FROM devices WHERE name = 'PLC-001';

-- ============================================================================
-- VIEWS
-- ============================================================================

-- Complete device view with protocol details
CREATE OR REPLACE VIEW v_devices_complete AS
SELECT 
    d.id,
    d.name,
    d.description,
    d.protocol,
    d.status,
    d.enabled,
    d.created_at,
    d.updated_at,
    CASE 
        WHEN d.protocol = 'modbus' THEN json_build_object(
            'ip', m.ip,
            'port', m.port,
            'unit_id', m.unit_id,
            'poll_interval_ms', m.poll_interval_ms,
            'register_count', (SELECT COUNT(*) FROM modbus_registers WHERE device_id = d.id)
        )
        WHEN d.protocol = 'opc' THEN json_build_object(
            'endpoint_url', o.endpoint_url,
            'security_mode', o.security_mode,
            'poll_interval_ms', o.poll_interval_ms,
            'node_count', (SELECT COUNT(*) FROM opc_nodes WHERE device_id = d.id)
        )
        WHEN d.protocol = 's7' THEN json_build_object(
            'ip', s.ip,
            'rack', s.rack,
            'slot', s.slot,
            'poll_interval_ms', s.poll_interval_ms,
            'block_count', (SELECT COUNT(*) FROM s7_data_blocks WHERE device_id = d.id)
        )
    END as protocol_config
FROM devices d
LEFT JOIN modbus_devices m ON d.id = m.id
LEFT JOIN opc_devices o ON d.id = o.id
LEFT JOIN s7_devices s ON d.id = s.id;

-- Latest config generation per device
CREATE OR REPLACE VIEW v_latest_configs AS
SELECT DISTINCT ON (device_id)
    device_id,
    id as config_id,
    config_hash,
    file_path,
    generated_at,
    deployed_at,
    status
FROM config_generations
ORDER BY device_id, generated_at DESC;

COMMENT ON VIEW v_devices_complete IS 'Complete device information with protocol-specific details';
COMMENT ON VIEW v_latest_configs IS 'Latest configuration generation for each device';
