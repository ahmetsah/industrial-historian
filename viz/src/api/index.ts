import axios from 'axios';
import type { SensorMetadata, SensorDataPoint, ActiveAlarm, AlarmDefinition, AuditLogEntry } from '../stores';

// Backend ports (from Docker):
// - Engine (Rust): 8081
// - Alarm (Go): 8083
// - Audit (Go): 8082
// - Auth (Go): 8080

const isDev = import.meta.env.DEV;

// In development, we use full URLs because different services are on different ports
// In production, a reverse proxy would route everything through a single domain
const ENGINE_API = isDev ? 'http://localhost:8081' : (import.meta.env.VITE_ENGINE_API || '');
const ALARM_API = isDev ? 'http://localhost:8083' : (import.meta.env.VITE_ALARM_API || '');
const AUDIT_API = isDev ? 'http://localhost:8082' : (import.meta.env.VITE_AUDIT_API || '');
const AUTH_API = isDev ? 'http://localhost:8080' : (import.meta.env.VITE_AUTH_API || '');

// Create axios instances
const engineClient = axios.create({
    baseURL: ENGINE_API,
    timeout: 10000,
});

const alarmClient = axios.create({
    baseURL: ALARM_API,
    timeout: 10000,
});

const auditClient = axios.create({
    baseURL: AUDIT_API,
    timeout: 10000,
});

const authClient = axios.create({
    baseURL: AUTH_API,
    timeout: 10000,
});

// Add auth interceptor
const addAuthInterceptor = (client: typeof axios) => {
    client.interceptors.request.use((config) => {
        const token = localStorage.getItem('auth_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    });
};

[engineClient, alarmClient, auditClient, authClient].forEach(client => {
    addAuthInterceptor(client as typeof axios);
});

// Engine API (Rust service on port 8081)
export const engineAPI = {
    // Get all sensor metadata
    getMetadata: async (): Promise<{ sensors: SensorMetadata[] }> => {
        const response = await engineClient.get('/api/v1/metadata');
        return response.data;
    },

    // Query sensor data
    querySensorData: async (
        sensorId: string,
        startTs: number,
        endTs: number,
        maxPoints: number = 1000
    ): Promise<SensorDataPoint[]> => {
        const response = await engineClient.get('/api/v1/query', {
            params: { sensor_id: sensorId, start_ts: startTs, end_ts: endTs, max_points: maxPoints }
        });
        return response.data.points || [];
    },

    // Export data as CSV
    exportCSV: async (sensorId: string, startTs: number, endTs: number): Promise<Blob> => {
        const response = await engineClient.get('/api/v1/export', {
            params: { sensor_id: sensorId, start_ts: startTs, end_ts: endTs, format: 'csv' },
            responseType: 'blob',
        });
        return response.data;
    },
};

// Alarm API (Go service on port 8083)
export const alarmAPI = {
    // Get active alarms
    getActiveAlarms: async (): Promise<ActiveAlarm[]> => {
        const response = await alarmClient.get('/api/v1/alarms/active');
        return response.data || [];
    },

    // Get alarm definitions
    getDefinitions: async (): Promise<AlarmDefinition[]> => {
        const response = await alarmClient.get('/api/v1/alarms/definitions');
        return response.data || [];
    },

    // Acknowledge an alarm
    acknowledgeAlarm: async (alarmId: number): Promise<void> => {
        await alarmClient.post(`/api/v1/alarms/${alarmId}/ack`);
    },

    // Shelve an alarm
    shelveAlarm: async (alarmId: number, durationSeconds: number): Promise<void> => {
        await alarmClient.post(`/api/v1/alarms/${alarmId}/shelve`, { duration_seconds: durationSeconds });
    },

    // Create alarm definition
    createDefinition: async (definition: Omit<AlarmDefinition, 'id'>): Promise<AlarmDefinition> => {
        const response = await alarmClient.post('/api/v1/alarms/definitions', definition);
        return response.data;
    },
};

// Audit API (Go service on port 8082)
// Note: Currently only verify endpoint exists in backend
export const auditAPI = {
    // Get audit logs - Note: This endpoint may not exist yet in backend
    getLogs: async (limit: number = 100, offset: number = 0): Promise<AuditLogEntry[]> => {
        try {
            const response = await auditClient.get('/api/v1/audit/logs', { params: { limit, offset } });
            return response.data || [];
        } catch (err) {
            console.warn('Audit logs endpoint not available:', err);
            return [];
        }
    },

    // Log an action
    logAction: async (action: string, details: Record<string, unknown>): Promise<void> => {
        await auditClient.post('/api/v1/audit/logs', { action, details });
    },

    // Verify log integrity
    verifyIntegrity: async (): Promise<{ valid: boolean; broken_id?: string }> => {
        const response = await auditClient.get('/api/v1/audit/verify');
        return response.data;
    },
};

// Auth API (Go service on port 8080)
export const authAPI = {
    // Login
    login: async (username: string, password: string): Promise<{ token: string; user: { id: string; name: string; role: string } }> => {
        const response = await authClient.post('/auth/login', { username, password });
        return response.data;
    },

    // Re-authenticate (FDA 21 CFR Part 11)
    reAuthenticate: async (password: string): Promise<{ signing_token: string }> => {
        const response = await authClient.post('/auth/re-auth', { password });
        return response.data;
    },

    // Logout
    logout: async (): Promise<void> => {
        await authClient.post('/auth/logout');
        localStorage.removeItem('auth_token');
    },

    // Validate token
    validateToken: async (): Promise<boolean> => {
        try {
            await authClient.get('/auth/validate');
            return true;
        } catch {
            return false;
        }
    },
};
