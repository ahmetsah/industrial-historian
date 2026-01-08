package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/historian/config-manager/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// ============================================================================
// GENERIC DEVICE OPERATIONS
// ============================================================================

func (r *DeviceRepository) ListDevices(protocol *models.Protocol, page, pageSize int) ([]models.Device, int64, error) {
	var devices []models.Device
	var total int64

	query := r.db.Model(&models.Device{})

	if protocol != nil {
		query = query.Where("protocol = ?", *protocol)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *DeviceRepository) GetDeviceByID(id uuid.UUID) (*models.Device, error) {
	var device models.Device
	if err := r.db.First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) DeleteDevice(id uuid.UUID) error {
	return r.db.Delete(&models.Device{}, "id = ?", id).Error
}

func (r *DeviceRepository) UpdateDeviceStatus(id uuid.UUID, status models.DeviceStatus) error {
	return r.db.Model(&models.Device{}).Where("id = ?", id).Update("status", status).Error
}

// ============================================================================
// MODBUS OPERATIONS
// ============================================================================

func (r *DeviceRepository) CreateModbusDevice(req *models.CreateModbusDeviceRequest, createdBy string) (*models.ModbusDevice, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create base device
	device := models.Device{
		Name:        req.Name,
		Description: req.Description,
		Protocol:    models.ProtocolModbus,
		Status:      models.StatusInactive,
		Enabled:     true,
		CreatedBy:   createdBy,
	}

	if err := tx.Create(&device).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Create Modbus-specific config
	modbusDevice := models.ModbusDevice{
		ID:             device.ID,
		IP:             req.IP,
		Port:           req.Port,
		UnitID:         req.UnitID,
		PollIntervalMs: req.PollIntervalMs,
		ConnectionType: "TCP",
	}

	if err := tx.Create(&modbusDevice).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create modbus device: %w", err)
	}

	// Create registers
	for _, regReq := range req.Registers {
		register := models.ModbusRegister{
			DeviceID:    device.ID,
			Address:     regReq.Address,
			Name:        regReq.Name,
			DataType:    regReq.DataType,
			ScaleFactor: regReq.ScaleFactor,
			Offset:      regReq.Offset,
			Unit:        regReq.Unit,
			Description: regReq.Description,
			Enabled:     true,
		}

		if err := tx.Create(&register).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create register: %w", err)
		}

		modbusDevice.Registers = append(modbusDevice.Registers, register)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	modbusDevice.Device = device
	return &modbusDevice, nil
}

func (r *DeviceRepository) GetModbusDevice(id uuid.UUID) (*models.ModbusDevice, error) {
	var modbusDevice models.ModbusDevice

	if err := r.db.Preload("Device").Preload("Registers").First(&modbusDevice, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &modbusDevice, nil
}

func (r *DeviceRepository) ListModbusDevices() ([]models.ModbusDevice, error) {
	var devices []models.ModbusDevice

	if err := r.db.Preload("Device").Preload("Registers").Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *DeviceRepository) UpdateModbusDevice(id uuid.UUID, req *models.CreateModbusDeviceRequest) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update base device
	if err := tx.Model(&models.Device{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update Modbus config
	if err := tx.Model(&models.ModbusDevice{}).Where("id = ?", id).Updates(map[string]interface{}{
		"ip":               req.IP,
		"port":             req.Port,
		"unit_id":          req.UnitID,
		"poll_interval_ms": req.PollIntervalMs,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete old registers and create new ones
	if err := tx.Where("device_id = ?", id).Delete(&models.ModbusRegister{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, regReq := range req.Registers {
		register := models.ModbusRegister{
			DeviceID:    id,
			Address:     regReq.Address,
			Name:        regReq.Name,
			DataType:    regReq.DataType,
			ScaleFactor: regReq.ScaleFactor,
			Offset:      regReq.Offset,
			Unit:        regReq.Unit,
			Description: regReq.Description,
			Enabled:     true,
		}

		if err := tx.Create(&register).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// ============================================================================
// OPC UA OPERATIONS
// ============================================================================

func (r *DeviceRepository) CreateOPCDevice(req *models.CreateOPCDeviceRequest, createdBy string) (*models.OPCDevice, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create base device
	device := models.Device{
		Name:        req.Name,
		Description: req.Description,
		Protocol:    models.ProtocolOPC,
		Status:      models.StatusInactive,
		Enabled:     true,
		CreatedBy:   createdBy,
	}

	if err := tx.Create(&device).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Create OPC-specific config
	opcDevice := models.OPCDevice{
		ID:             device.ID,
		EndpointURL:    req.EndpointURL,
		SecurityMode:   req.SecurityMode,
		Username:       req.Username,
		PollIntervalMs: req.PollIntervalMs,
	}

	// TODO: Encrypt password before storing
	if req.Password != "" {
		opcDevice.PasswordEncrypted = req.Password // Should be encrypted!
	}

	if err := tx.Create(&opcDevice).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create opc device: %w", err)
	}

	// Create nodes
	for _, nodeReq := range req.Nodes {
		node := models.OPCNode{
			DeviceID:    device.ID,
			NodeID:      nodeReq.NodeID,
			Name:        nodeReq.Name,
			DataType:    nodeReq.DataType,
			Unit:        nodeReq.Unit,
			Description: nodeReq.Description,
			Enabled:     true,
		}

		if err := tx.Create(&node).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create node: %w", err)
		}

		opcDevice.Nodes = append(opcDevice.Nodes, node)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	opcDevice.Device = device
	return &opcDevice, nil
}

func (r *DeviceRepository) GetOPCDevice(id uuid.UUID) (*models.OPCDevice, error) {
	var opcDevice models.OPCDevice

	if err := r.db.Preload("Device").Preload("Nodes").First(&opcDevice, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &opcDevice, nil
}

func (r *DeviceRepository) ListOPCDevices() ([]models.OPCDevice, error) {
	var devices []models.OPCDevice

	if err := r.db.Preload("Device").Preload("Nodes").Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// ============================================================================
// CONFIG GENERATION OPERATIONS
// ============================================================================

func (r *DeviceRepository) SaveConfigGeneration(config *models.ConfigGeneration) error {
	return r.db.Create(config).Error
}

func (r *DeviceRepository) GetLatestConfig(deviceID uuid.UUID) (*models.ConfigGeneration, error) {
	var config models.ConfigGeneration

	if err := r.db.Where("device_id = ?", deviceID).Order("generated_at DESC").First(&config).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

func (r *DeviceRepository) UpdateConfigStatus(configID uuid.UUID, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	return r.db.Model(&models.ConfigGeneration{}).Where("id = ?", configID).Updates(updates).Error
}
