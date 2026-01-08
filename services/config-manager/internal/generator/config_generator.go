package generator

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	"github.com/historian/config-manager/internal/models"
)

type ConfigGenerator struct {
	configDir string
}

func NewConfigGenerator(configDir string) *ConfigGenerator {
	return &ConfigGenerator{
		configDir: configDir,
	}
}

// ============================================================================
// MODBUS CONFIG GENERATION
// ============================================================================

const modbusConfigTemplate = `# Modbus Device: {{.Device.Name}}
# Generated at: {{.GeneratedAt}}
# Device ID: {{.Device.ID}}

[[modbus_devices]]
ip = "{{.IP}}"
port = {{.Port}}
unit_id = {{.UnitID}}
poll_interval_ms = {{.PollIntervalMs}}
timeout_ms = {{.TimeoutMs}}
retry_count = {{.RetryCount}}

{{range .Registers}}
[[modbus_devices.registers]]
address = {{.Address}}
name = "{{.Name}}"
data_type = "{{.DataType}}"
register_type = "{{.RegisterType}}"
scale_factor = {{.ScaleFactor}}
offset = {{.Offset}}
{{if .Unit}}unit = "{{.Unit}}"{{end}}
{{if .Description}}description = "{{.Description}}"{{end}}
{{end}}

[nats]
url = "nats://nats:4222"
subject = "data.modbus.{{.Device.Name}}"

[buffer]
memory_capacity = 10000
disk_path = "/data/buffer/{{.Device.Name}}.wal"
`

func (g *ConfigGenerator) GenerateModbusConfig(device *models.ModbusDevice) (*models.ConfigGeneration, error) {
	// Parse template
	tmpl, err := template.New("modbus").Parse(modbusConfigTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	data := struct {
		*models.ModbusDevice
		GeneratedAt string
	}{
		ModbusDevice: device,
		GeneratedAt:  fmt.Sprintf("%v", device.Device.UpdatedAt),
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	configContent := buf.String()

	// Calculate hash
	hash := sha256.Sum256([]byte(configContent))
	configHash := fmt.Sprintf("%x", hash)

	// Generate file path
	fileName := fmt.Sprintf("modbus-%s.toml", device.Device.Name)
	filePath := filepath.Join(g.configDir, fileName)

	// Save to file
	if err := os.MkdirAll(g.configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(configContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	// Create config generation record
	configGen := &models.ConfigGeneration{
		ID:            uuid.New(),
		DeviceID:      device.ID,
		ConfigHash:    configHash,
		ConfigContent: configContent,
		FilePath:      filePath,
		Status:        "generated",
	}

	return configGen, nil
}

// ============================================================================
// OPC UA CONFIG GENERATION
// ============================================================================

const opcConfigTemplate = `# OPC UA Device: {{.Device.Name}}
# Generated at: {{.GeneratedAt}}
# Device ID: {{.Device.ID}}

[[opc_devices]]
endpoint_url = "{{.EndpointURL}}"
security_mode = "{{.SecurityMode}}"
security_policy = "{{.SecurityPolicy}}"
{{if .Username}}username = "{{.Username}}"{{end}}
{{if .PasswordEncrypted}}password = "${OPC_PASSWORD_{{.Device.Name}}}"{{end}}
poll_interval_ms = {{.PollIntervalMs}}
timeout_ms = {{.TimeoutMs}}

{{range .Nodes}}
[[opc_devices.nodes]]
node_id = "{{.NodeID}}"
name = "{{.Name}}"
{{if .DataType}}data_type = "{{.DataType}}"{{end}}
namespace_index = {{.NamespaceIndex}}
{{if .Unit}}unit = "{{.Unit}}"{{end}}
{{if .Description}}description = "{{.Description}}"{{end}}
{{end}}

[nats]
url = "${NATS_URL}"
subject = "data.opc.{{.Device.Name}}"

[buffer]
memory_capacity = 10000
disk_path = "/data/buffer/{{.Device.Name}}.wal"
`

func (g *ConfigGenerator) GenerateOPCConfig(device *models.OPCDevice) (*models.ConfigGeneration, error) {
	tmpl, err := template.New("opc").Parse(opcConfigTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		*models.OPCDevice
		GeneratedAt string
	}{
		OPCDevice:   device,
		GeneratedAt: fmt.Sprintf("%v", device.Device.UpdatedAt),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	configContent := buf.String()
	hash := sha256.Sum256([]byte(configContent))
	configHash := fmt.Sprintf("%x", hash)

	fileName := fmt.Sprintf("opc-%s.toml", device.Device.Name)
	filePath := filepath.Join(g.configDir, fileName)

	if err := os.MkdirAll(g.configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(configContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	configGen := &models.ConfigGeneration{
		ID:            uuid.New(),
		DeviceID:      device.ID,
		ConfigHash:    configHash,
		ConfigContent: configContent,
		FilePath:      filePath,
		Status:        "generated",
	}

	return configGen, nil
}

// ============================================================================
// S7 CONFIG GENERATION
// ============================================================================

const s7ConfigTemplate = `# S7 Device: {{.Device.Name}}
# Generated at: {{.GeneratedAt}}
# Device ID: {{.Device.ID}}

[[s7_devices]]
ip = "{{.IP}}"
rack = {{.Rack}}
slot = {{.Slot}}
plc_type = "{{.PLCType}}"
poll_interval_ms = {{.PollIntervalMs}}
timeout_ms = {{.TimeoutMs}}

{{range .DataBlocks}}
[[s7_devices.data_blocks]]
db_number = {{.DBNumber}}
start_address = {{.StartAddress}}
length = {{.Length}}
name = "{{.Name}}"
data_type = "{{.DataType}}"
{{if .Unit}}unit = "{{.Unit}}"{{end}}
{{if .Description}}description = "{{.Description}}"{{end}}
{{end}}

[nats]
url = "${NATS_URL}"
subject = "data.s7.{{.Device.Name}}"

[buffer]
memory_capacity = 10000
disk_path = "/data/buffer/{{.Device.Name}}.wal"
`

func (g *ConfigGenerator) GenerateS7Config(device *models.S7Device) (*models.ConfigGeneration, error) {
	tmpl, err := template.New("s7").Parse(s7ConfigTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		*models.S7Device
		GeneratedAt string
	}{
		S7Device:    device,
		GeneratedAt: fmt.Sprintf("%v", device.Device.UpdatedAt),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	configContent := buf.String()
	hash := sha256.Sum256([]byte(configContent))
	configHash := fmt.Sprintf("%x", hash)

	fileName := fmt.Sprintf("s7-%s.toml", device.Device.Name)
	filePath := filepath.Join(g.configDir, fileName)

	if err := os.MkdirAll(g.configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(configContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	configGen := &models.ConfigGeneration{
		ID:            uuid.New(),
		DeviceID:      device.ID,
		ConfigHash:    configHash,
		ConfigContent: configContent,
		FilePath:      filePath,
		Status:        "generated",
	}

	return configGen, nil
}

// ============================================================================
// DOCKER COMPOSE GENERATION
// ============================================================================

func (g *ConfigGenerator) GenerateDockerCompose(devices []models.Device) error {
	const dockerComposeTemplate = `version: '3.8'

services:
{{range .ModbusDevices}}
  ingestor-modbus-{{.Device.Name}}:
    image: historian-ingestor-modbus:latest
    environment:
      - NATS_URL=nats://nats:4222
      - CONFIG_FILE=/config/modbus-{{.Device.Name}}.toml
      - DEVICE_NAME={{.Device.Name}}
    volumes:
      - ./config/generated:/config:ro
      - ./data/buffer:/data/buffer
    depends_on:
      - nats
    restart: unless-stopped
{{end}}

{{range .OPCDevices}}
  ingestor-opc-{{.Device.Name}}:
    image: historian-ingestor-opc:latest
    environment:
      - NATS_URL=nats://nats:4222
      - CONFIG_FILE=/config/opc-{{.Device.Name}}.toml
      - DEVICE_NAME={{.Device.Name}}
    volumes:
      - ./config/generated:/config:ro
      - ./data/buffer:/data/buffer
    depends_on:
      - nats
    restart: unless-stopped
{{end}}

  nats:
    image: nats:latest
    ports:
      - "4222:4222"
    command: "-js -sd /data"
    volumes:
      - nats_data:/data

volumes:
  nats_data:
`

	// This would be called with actual device data
	// For now, it's a placeholder
	return nil
}

// ============================================================================
// FILE MANAGEMENT
// ============================================================================

// DeleteConfigFile removes a config file from the filesystem
func (g *ConfigGenerator) DeleteConfigFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, consider it already deleted
		return nil
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete config file %s: %w", filePath, err)
	}

	return nil
}
