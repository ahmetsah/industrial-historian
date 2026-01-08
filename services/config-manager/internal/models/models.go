package models

import (
	"time"

	"github.com/google/uuid"
)

// Protocol types
type Protocol string

const (
	ProtocolModbus Protocol = "modbus"
	ProtocolOPC    Protocol = "opc"
	ProtocolS7     Protocol = "s7"
	ProtocolMQTT   Protocol = "mqtt"
)

// Device status
type DeviceStatus string

const (
	StatusActive   DeviceStatus = "active"
	StatusInactive DeviceStatus = "inactive"
	StatusError    DeviceStatus = "error"
	StatusPending  DeviceStatus = "pending"
)

// Base Device model
type Device struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Description string       `json:"description"`
	Protocol    Protocol     `json:"protocol" gorm:"type:protocol_type;not null"`
	Status      DeviceStatus `json:"status" gorm:"type:device_status;default:'inactive'"`
	Enabled     bool         `json:"enabled" gorm:"default:true"`
	Metadata    string       `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CreatedBy   string       `json:"created_by"`
	UpdatedBy   string       `json:"updated_by"`
}

// ============================================================================
// MODBUS MODELS
// ============================================================================

type ModbusDevice struct {
	ID             uuid.UUID        `json:"id" gorm:"type:uuid;primary_key"`
	IP             string           `json:"ip" gorm:"not null"`
	Port           int              `json:"port" gorm:"default:502"`
	UnitID         int              `json:"unit_id" gorm:"not null"`
	PollIntervalMs int              `json:"poll_interval_ms" gorm:"default:1000"`
	TimeoutMs      int              `json:"timeout_ms" gorm:"default:5000"`
	RetryCount     int              `json:"retry_count" gorm:"default:3"`
	ConnectionType string           `json:"connection_type" gorm:"default:'TCP'"`
	Registers      []ModbusRegister `json:"registers" gorm:"foreignKey:DeviceID"`
	Device         Device           `json:"device" gorm:"foreignKey:ID"`
	RuntimeStatus  string           `json:"runtime_status" gorm:"-"`
}

type ModbusRegister struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DeviceID     uuid.UUID `json:"device_id" gorm:"type:uuid;not null"`
	Address      int       `json:"address" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null"`
	DataType     string    `json:"data_type" gorm:"not null"`
	RegisterType string    `json:"register_type" gorm:"default:'holding'"`
	ScaleFactor  float64   `json:"scale_factor" gorm:"default:1.0"`
	Offset       float64   `json:"offset" gorm:"default:0.0"`
	Unit         string    `json:"unit"`
	Description  string    `json:"description"`
	Enabled      bool      `json:"enabled" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
}

// ============================================================================
// OPC UA MODELS
// ============================================================================

type OPCDevice struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	EndpointURL       string    `json:"endpoint_url" gorm:"uniqueIndex;not null"`
	SecurityMode      string    `json:"security_mode" gorm:"default:'None'"`
	SecurityPolicy    string    `json:"security_policy" gorm:"default:'None'"`
	Username          string    `json:"username"`
	PasswordEncrypted string    `json:"password_encrypted"`
	CertificatePath   string    `json:"certificate_path"`
	PollIntervalMs    int       `json:"poll_interval_ms" gorm:"default:1000"`
	TimeoutMs         int       `json:"timeout_ms" gorm:"default:5000"`
	Nodes             []OPCNode `json:"nodes" gorm:"foreignKey:DeviceID"`
	Device            Device    `json:"device" gorm:"foreignKey:ID"`
}

type OPCNode struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DeviceID       uuid.UUID `json:"device_id" gorm:"type:uuid;not null"`
	NodeID         string    `json:"node_id" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	DataType       string    `json:"data_type"`
	NamespaceIndex int       `json:"namespace_index" gorm:"default:0"`
	Unit           string    `json:"unit"`
	Description    string    `json:"description"`
	Enabled        bool      `json:"enabled" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
}

// ============================================================================
// S7 MODELS
// ============================================================================

type S7Device struct {
	ID             uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	IP             string        `json:"ip" gorm:"uniqueIndex;not null"`
	Rack           int           `json:"rack" gorm:"default:0"`
	Slot           int           `json:"slot" gorm:"default:1"`
	PollIntervalMs int           `json:"poll_interval_ms" gorm:"default:1000"`
	TimeoutMs      int           `json:"timeout_ms" gorm:"default:5000"`
	PLCType        string        `json:"plc_type" gorm:"default:'S7-1200'"`
	DataBlocks     []S7DataBlock `json:"data_blocks" gorm:"foreignKey:DeviceID"`
	Device         Device        `json:"device" gorm:"foreignKey:ID"`
}

type S7DataBlock struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DeviceID     uuid.UUID `json:"device_id" gorm:"type:uuid;not null"`
	DBNumber     int       `json:"db_number" gorm:"not null"`
	StartAddress int       `json:"start_address" gorm:"not null"`
	Length       int       `json:"length" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null"`
	DataType     string    `json:"data_type" gorm:"not null"`
	Unit         string    `json:"unit"`
	Description  string    `json:"description"`
	Enabled      bool      `json:"enabled" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
}

// ============================================================================
// CONFIG GENERATION MODELS
// ============================================================================

type ConfigGeneration struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DeviceID      uuid.UUID  `json:"device_id" gorm:"type:uuid;not null"`
	ConfigHash    string     `json:"config_hash" gorm:"not null"`
	ConfigContent string     `json:"config_content" gorm:"type:text;not null"`
	FilePath      string     `json:"file_path" gorm:"not null"`
	GeneratedAt   time.Time  `json:"generated_at"`
	DeployedAt    *time.Time `json:"deployed_at"`
	Status        string     `json:"status" gorm:"default:'pending'"`
	ErrorMessage  string     `json:"error_message"`
	GeneratedBy   string     `json:"generated_by"`
}

// ============================================================================
// REQUEST/RESPONSE DTOs
// ============================================================================

// CreateModbusDeviceRequest represents the request to create a Modbus device
type CreateModbusDeviceRequest struct {
	Name           string                        `json:"name" binding:"required"`
	Description    string                        `json:"description"`
	IP             string                        `json:"ip" binding:"required,ip"`
	Port           int                           `json:"port" binding:"required,min=1,max=65535"`
	UnitID         int                           `json:"unit_id" binding:"required,min=0,max=255"`
	PollIntervalMs int                           `json:"poll_interval_ms" binding:"required,min=100"`
	Registers      []CreateModbusRegisterRequest `json:"registers" binding:"required,dive"`
}

type CreateModbusRegisterRequest struct {
	Address     int     `json:"address" binding:"min=0,max=65535"`
	Name        string  `json:"name" binding:"required"`
	DataType    string  `json:"data_type" binding:"required,oneof=Int16 UInt16 Int32 UInt32 Float32 Float64"`
	ScaleFactor float64 `json:"scale_factor"`
	Offset      float64 `json:"offset"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
}

// CreateOPCDeviceRequest represents the request to create an OPC UA device
type CreateOPCDeviceRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Description    string                 `json:"description"`
	EndpointURL    string                 `json:"endpoint_url" binding:"required,url"`
	SecurityMode   string                 `json:"security_mode" binding:"oneof=None Sign SignAndEncrypt"`
	Username       string                 `json:"username"`
	Password       string                 `json:"password"`
	PollIntervalMs int                    `json:"poll_interval_ms" binding:"required,min=100"`
	Nodes          []CreateOPCNodeRequest `json:"nodes" binding:"required,dive"`
}

type CreateOPCNodeRequest struct {
	NodeID      string `json:"node_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	DataType    string `json:"data_type"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

// DeviceListResponse represents a paginated list of devices
type DeviceListResponse struct {
	Devices    []Device `json:"devices"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

type UpdateModbusDeviceRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	PollIntervalMs int    `json:"poll_interval_ms"`
	TimeoutMs      int    `json:"timeout_ms"`
	RetryCount     int    `json:"retry_count"`
}

// ============================================================================
// INGESTOR CONFIG MODELS (JSON Response for Ingestor)
// ============================================================================

type IngestorConfig struct {
	ModbusDevices []ModbusDeviceConfig `json:"modbus_devices"`
	Nats          NatsConfig           `json:"nats"`
	Buffer        BufferConfig         `json:"buffer"`
}

type ModbusDeviceConfig struct {
	IP             string           `json:"ip"`
	Port           int              `json:"port"`
	UnitID         int              `json:"unit_id"`
	PollIntervalMs int              `json:"poll_interval_ms"`
	TimeoutMs      int              `json:"timeout_ms"`
	RetryCount     int              `json:"retry_count"`
	Registers      []RegisterConfig `json:"registers"`
}

type RegisterConfig struct {
	Address      int     `json:"address"`
	Name         string  `json:"name"`
	DataType     string  `json:"data_type"`
	RegisterType string  `json:"register_type"`
	ScaleFactor  float64 `json:"scale_factor"`
	Offset       float64 `json:"offset"`
	Unit         string  `json:"unit,omitempty"`
	Description  string  `json:"description,omitempty"`
}

type NatsConfig struct {
	URL     string `json:"url"`
	Subject string `json:"subject"`
}

type BufferConfig struct {
	MemoryCapacity int    `json:"memory_capacity"`
	DiskPath       string `json:"disk_path"`
}

// ConfigGenerationResponse represents the response after generating a config
type ConfigGenerationResponse struct {
	ConfigID   uuid.UUID `json:"config_id"`
	DeviceID   uuid.UUID `json:"device_id"`
	FilePath   string    `json:"file_path"`
	ConfigHash string    `json:"config_hash"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
}
