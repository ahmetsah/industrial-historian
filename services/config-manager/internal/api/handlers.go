package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/historian/config-manager/internal/generator"
	"github.com/historian/config-manager/internal/models"
	"github.com/historian/config-manager/internal/repository"
)

type Handler struct {
	repo      *repository.DeviceRepository
	generator *generator.ConfigGenerator
}

func NewHandler(repo *repository.DeviceRepository, gen *generator.ConfigGenerator) *Handler {
	return &Handler{
		repo:      repo,
		generator: gen,
	}
}

// ============================================================================
// GENERIC DEVICE HANDLERS
// ============================================================================

func (h *Handler) ListDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var protocol *models.Protocol
	if p := c.Query("protocol"); p != "" {
		proto := models.Protocol(p)
		protocol = &proto
	}

	devices, total, err := h.repo.ListDevices(protocol, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.DeviceListResponse{
		Devices:    devices,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *Handler) GetDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *Handler) DeleteDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	// Get device info before deleting (to get config file path)
	device, err := h.repo.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	// Get latest config to find file path (Legacy support)
	config, err := h.repo.GetLatestConfig(id)
	if err == nil && config != nil && config.FilePath != "" {
		// Try to delete physical config file (ignore error if not exists)
		_ = h.generator.DeleteConfigFile(config.FilePath)
	}

	// Delete Docker Container
	containerName := "ops-ingestor-modbus-" + strings.ToLower(strings.ReplaceAll(device.Name, " ", ""))
	// docker rm -f containerName
	exec.Command("docker", "rm", "-f", containerName).Run() // Ignore error if container doesn't exist

	// Delete from database (cascade will handle related records)
	if err := h.repo.DeleteDevice(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message":        "Device and container deleted successfully",
		"device_name":    device.Name,
		"container_name": containerName,
	}
	if config != nil {
		response["config_file"] = config.FilePath
	}

	c.JSON(http.StatusOK, response)
}

// ============================================================================
// MODBUS HANDLERS
// ============================================================================

func (h *Handler) ListModbusDevices(c *gin.Context) {
	devices, err := h.repo.ListModbusDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	containerStatuses := h.getContainerStatuses()
	var responseList []map[string]interface{}

	for _, dev := range devices {
		// Convert struct to map manually to avoid model pollution and ensure fields exist
		var devMap map[string]interface{}
		devBytes, _ := json.Marshal(dev)
		_ = json.Unmarshal(devBytes, &devMap)

		containerName := sanitizeContainerName(dev.Device.Name)

		deployStatus := "not_deployed"
		connStatus := "idle" // Default when not deployed

		// Check Container Status
		if state, exists := containerStatuses[containerName]; exists && state == "running" {
			deployStatus = "deployed"

			// Check Connectivity (TCP Ping) only if deployed
			target := fmt.Sprintf("%s:%d", dev.IP, dev.Port)
			conn, err := net.DialTimeout("tcp", target, 500*time.Millisecond)
			if err == nil {
				connStatus = "connected"
				conn.Close()
			} else {
				connStatus = "disconnected"
			}
		}

		devMap["deployment_status"] = deployStatus
		devMap["connection_status"] = connStatus

		responseList = append(responseList, devMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": responseList,
		"total":   len(devices),
	})
}

func (h *Handler) CreateModbusDevice(c *gin.Context) {
	var req models.CreateModbusDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (for now, use default)
	createdBy := c.GetString("user")
	if createdBy == "" {
		createdBy = "admin"
	}

	device, err := h.repo.CreateModbusDevice(&req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate config
	config, err := h.generator.GenerateModbusConfig(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Device created but config generation failed",
			"details": err.Error(),
			"device":  device,
		})
		return
	}

	// Save config generation record
	if err := h.repo.SaveConfigGeneration(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Config generated but failed to save record",
			"details": err.Error(),
			"device":  device,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"device": device,
		"config": gin.H{
			"id":        config.ID,
			"file_path": config.FilePath,
			"hash":      config.ConfigHash,
		},
		"message": "Device created and config generated successfully",
	})
}

func (h *Handler) GetModbusDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetModbusDevice(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *Handler) UpdateModbusDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var req models.CreateModbusDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateModbusDevice(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Regenerate config
	device, err := h.repo.GetModbusDevice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated device"})
		return
	}

	config, err := h.generator.GenerateModbusConfig(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to regenerate config"})
		return
	}

	if err := h.repo.SaveConfigGeneration(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device":  device,
		"config":  config,
		"message": "Device updated and config regenerated",
	})
}

// ============================================================================
// OPC UA HANDLERS
// ============================================================================

func (h *Handler) ListOPCDevices(c *gin.Context) {
	devices, err := h.repo.ListOPCDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
		"total":   len(devices),
	})
}

func (h *Handler) CreateOPCDevice(c *gin.Context) {
	var req models.CreateOPCDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := c.GetString("user")
	if createdBy == "" {
		createdBy = "admin"
	}

	device, err := h.repo.CreateOPCDevice(&req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate config
	config, err := h.generator.GenerateOPCConfig(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Device created but config generation failed",
			"details": err.Error(),
			"device":  device,
		})
		return
	}

	if err := h.repo.SaveConfigGeneration(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Config generated but failed to save record",
			"details": err.Error(),
			"device":  device,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"device": device,
		"config": gin.H{
			"id":        config.ID,
			"file_path": config.FilePath,
			"hash":      config.ConfigHash,
		},
		"message": "OPC UA device created and config generated successfully",
	})
}

func (h *Handler) GetOPCDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetOPCDevice(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// ============================================================================
// CONFIG GENERATION HANDLERS
// ============================================================================

func (h *Handler) GenerateConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	var config *models.ConfigGeneration

	switch device.Protocol {
	case models.ProtocolModbus:
		modbusDevice, err := h.repo.GetModbusDevice(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Modbus device"})
			return
		}
		config, err = h.generator.GenerateModbusConfig(modbusDevice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case models.ProtocolOPC:
		opcDevice, err := h.repo.GetOPCDevice(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch OPC device"})
			return
		}
		config, err = h.generator.GenerateOPCConfig(opcDevice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported protocol"})
		return
	}

	if err := h.repo.SaveConfigGeneration(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, models.ConfigGenerationResponse{
		ConfigID:   config.ID,
		DeviceID:   config.DeviceID,
		FilePath:   config.FilePath,
		ConfigHash: config.ConfigHash,
		Status:     config.Status,
		Message:    "Config generated successfully",
	})
}

func (h *Handler) GetLatestConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	config, err := h.repo.GetLatestConfig(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No config found for this device"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// ============================================================================
// DYNAMIC CONFIG & DEPLOYMENT API
// ============================================================================

// GetDeviceConfig returns the JSON configuration for the ingestor service
func (h *Handler) GetDeviceConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if device.Protocol != models.ProtocolModbus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only Modbus supported for now"})
		return
	}

	modbusDevice, err := h.repo.GetModbusDevice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Modbus details"})
		return
	}

	// Prepare Registers
	var registers []models.RegisterConfig
	for _, reg := range modbusDevice.Registers {
		registers = append(registers, models.RegisterConfig{
			Address:      reg.Address,
			Name:         reg.Name,
			DataType:     reg.DataType,
			RegisterType: "HoldingRegister", // Default
			ScaleFactor:  reg.ScaleFactor,
			Offset:       reg.Offset,
			Unit:         reg.Unit,
			Description:  reg.Description,
		})
	}

	// Construct Config Response
	response := models.IngestorConfig{
		ModbusDevices: []models.ModbusDeviceConfig{
			{
				IP:             modbusDevice.IP,
				Port:           modbusDevice.Port,
				UnitID:         modbusDevice.UnitID,
				PollIntervalMs: modbusDevice.PollIntervalMs,
				TimeoutMs:      modbusDevice.TimeoutMs,
				RetryCount:     modbusDevice.RetryCount,
				Registers:      registers,
			},
		},
		Nats: models.NatsConfig{
			URL:     "nats://nats:4222",
			Subject: fmt.Sprintf("data.modbus.%s", modbusDevice.Device.Name),
		},
		Buffer: models.BufferConfig{
			MemoryCapacity: 10000,
			DiskPath:       fmt.Sprintf("/data/buffer/%s.wal", modbusDevice.Device.Name),
		},
	}

	c.JSON(http.StatusOK, response)
}

// DeployDevice checks/creates Docker container for the device
func (h *Handler) DeployDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	containerName := "ops-ingestor-modbus-" + strings.ToLower(strings.ReplaceAll(device.Name, " ", ""))
	configURL := fmt.Sprintf("http://config-manager:8090/api/v1/devices/%s/config", id.String())

	// Check if container exists
	checkCmd := exec.Command("docker", "inspect", containerName)
	if err := checkCmd.Run(); err == nil {
		// Container exists -> Restart (Hot Reload)
		restartCmd := exec.Command("docker", "restart", containerName)
		if out, err := restartCmd.CombinedOutput(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to restart container",
				"details": string(out),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":        "Device hot-reloaded successfully",
			"container_name": containerName,
			"action":         "restarted",
			"config_url":     configURL,
		})
		return
	}

	// Container does not exist -> Run new one
	// Note: We assume 'historian-ingestor:latest' image is available
	// docker run -d --name ... --network historian-net -e CONFIG_URL=... ...
	runCmd := exec.Command("docker", "run", "-d",
		"--name", containerName,
		"--network", "ops_historian-net",
		"--restart", "unless-stopped",
		"-e", fmt.Sprintf("CONFIG_URL=%s", configURL),
		"-e", "RUST_LOG=info",
		"--add-host", "host.docker.internal:host-gateway",
		"-v", fmt.Sprintf("ingestor_buffer_%s:/data/buffer", strings.ToLower(device.Name)),
		"historian-ingestor:latest",
	)

	if out, err := runCmd.CombinedOutput(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create container",
			"details": string(out),
			"command": runCmd.String(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Device deployed successfully",
		"container_name": containerName,
		"action":         "created",
		"config_url":     configURL,
	})
}

// StopDevice stops the Docker container for the device
func (h *Handler) StopDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.repo.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	containerName := "ops-ingestor-modbus-" + strings.ToLower(strings.ReplaceAll(device.Name, " ", ""))

	// Check if container exists
	checkCmd := exec.Command("docker", "inspect", containerName)
	if err := checkCmd.Run(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":        "Container not running",
			"container_name": containerName,
			"action":         "not_found",
		})
		return
	}

	// Stop the container
	stopCmd := exec.Command("docker", "stop", containerName)
	if out, err := stopCmd.CombinedOutput(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to stop container",
			"details": string(out),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Device stopped successfully",
		"container_name": containerName,
		"action":         "stopped",
	})
}

// ============================================================================
// HELPERS
// ============================================================================

func (h *Handler) getContainerStatuses() map[string]string {
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}|{{.State}}")
	out, err := cmd.Output()
	if err != nil {
		return make(map[string]string)
	}

	statuses := make(map[string]string)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			statuses[parts[0]] = parts[1] // e.g. "running"
		}
	}
	return statuses
}

func sanitizeContainerName(input string) string {
	return "ops-ingestor-modbus-" + strings.ToLower(strings.ReplaceAll(input, " ", ""))
}
