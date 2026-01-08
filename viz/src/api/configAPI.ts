import axios from 'axios';

// Config Manager API - uses Vite proxy in development
const isDev = import.meta.env.DEV;

// In development, use relative paths so Vite proxy handles requests
// In production, use environment variable or direct URL
const CONFIG_API = isDev ? '' : (import.meta.env.VITE_CONFIG_API || 'http://localhost:8090');

const configClient = axios.create({
    baseURL: CONFIG_API,
    timeout: 10000,
});

// Types
export interface Register {
    address: number;
    name: string;
    data_type: 'Bool' | 'Int16' | 'UInt16' | 'Int32' | 'UInt32' | 'Float32' | 'Int64' | 'UInt64' | 'Float64';
}

export interface ModbusDeviceBase {
    name: string;
    description?: string;
    ip: string;
    port: number;
    unit_id: number;
    poll_interval_ms: number;
    registers: Register[];
}

export interface ModbusDevice extends ModbusDeviceBase {
    device: {
        id: string;
        name: string;
        description: string;
    };
    deployment_status?: 'deployed' | 'not_deployed';
    connection_status?: 'connected' | 'disconnected' | 'idle';
    runtime_status?: 'running' | 'stopped';
}

export interface DeviceStats {
    totalDevices: number;
    activeDevices: number;
    deployedDevices: number;
    totalRegisters: number;
}

// Config API functions
export const configAPI = {
    // Get all Modbus devices
    getDevices: async (): Promise<ModbusDevice[]> => {
        const response = await configClient.get('/api/v1/devices/modbus');
        return response.data.devices || [];
    },

    // Get single device by ID
    getDevice: async (deviceId: string): Promise<ModbusDevice> => {
        const response = await configClient.get(`/api/v1/devices/modbus/${deviceId}`);
        return response.data;
    },

    // Create new Modbus device
    createDevice: async (device: ModbusDeviceBase): Promise<ModbusDevice> => {
        const response = await configClient.post('/api/v1/devices/modbus', device);
        return response.data;
    },

    // Update existing device
    updateDevice: async (deviceId: string, device: ModbusDeviceBase): Promise<ModbusDevice> => {
        const response = await configClient.put(`/api/v1/devices/modbus/${deviceId}`, device);
        return response.data;
    },

    // Delete device (uses generic endpoint, not /modbus)
    deleteDevice: async (deviceId: string): Promise<{ message: string }> => {
        const response = await configClient.delete(`/api/v1/devices/${deviceId}`);
        return response.data;
    },

    // Deploy device (start/restart container)
    deployDevice: async (deviceId: string): Promise<{ message: string; action: string; container_name: string }> => {
        const response = await configClient.post(`/api/v1/devices/${deviceId}/deploy`);
        return response.data;
    },

    // Stop device (stop container)
    stopDevice: async (deviceId: string): Promise<{ message: string }> => {
        const response = await configClient.post(`/api/v1/devices/${deviceId}/stop`);
        return response.data;
    },

    // Calculate stats from devices
    calculateStats: (devices: ModbusDevice[]): DeviceStats => {
        return {
            totalDevices: devices.length,
            activeDevices: devices.filter(d => d.connection_status === 'connected').length,
            deployedDevices: devices.filter(d => d.deployment_status === 'deployed').length,
            totalRegisters: devices.reduce((sum, d) => sum + (d.registers?.length || 0), 0),
        };
    },
};
