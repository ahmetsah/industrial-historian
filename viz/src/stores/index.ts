import { create } from 'zustand';

// Types
export interface SensorMetadata {
    id: string;
    desc: string;
    factory: string;
    line: string;
    machine: string;
    type: string;
    unit: string;
}

export interface SensorDataPoint {
    timestamp: number;
    value: number;
}

export interface ActiveAlarm {
    id: number;
    definition_id: number;
    state: 'UnackActive' | 'AckActive' | 'UnackRTN' | 'Shelved' | 'Normal';
    activation_time: string;
    ack_time?: string;
    shelved_until?: string;
    value: number;
    tag?: string;
    priority?: string;
    type?: string;
}

export interface AlarmDefinition {
    id: number;
    tag: string;
    threshold: number;
    type: 'High' | 'Low';
    priority: 'Critical' | 'Warning';
}

export interface AuditLogEntry {
    id: string;
    timestamp: string;
    actor: string;
    action: string;
    details: Record<string, unknown>;
    prev_hash: string;
    curr_hash: string;
}

export interface SystemStats {
    activeAlarms: number;
    criticalAlarms: number;
    totalSensors: number;
    dataPointsToday: number;
    uptime: string;
    cpuUsage: number;
    memoryUsage: number;
}

// Sensor Store
interface SensorState {
    sensors: SensorMetadata[];
    selectedSensors: string[];
    sensorData: Record<string, SensorDataPoint[]>;
    loading: boolean;
    error: string | null;
    setSensors: (sensors: SensorMetadata[]) => void;
    selectSensor: (sensorId: string) => void;
    deselectSensor: (sensorId: string) => void;
    clearSelection: () => void;
    setSensorData: (sensorId: string, data: SensorDataPoint[]) => void;
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
}

export const useSensorStore = create<SensorState>((set) => ({
    sensors: [],
    selectedSensors: [],
    sensorData: {},
    loading: false,
    error: null,
    setSensors: (sensors) => set({ sensors }),
    selectSensor: (sensorId) =>
        set((state) => ({
            selectedSensors: state.selectedSensors.includes(sensorId)
                ? state.selectedSensors
                : [...state.selectedSensors, sensorId]
        })),
    deselectSensor: (sensorId) =>
        set((state) => ({
            selectedSensors: state.selectedSensors.filter(id => id !== sensorId)
        })),
    clearSelection: () => set({ selectedSensors: [] }),
    setSensorData: (sensorId, data) =>
        set((state) => ({
            sensorData: { ...state.sensorData, [sensorId]: data }
        })),
    setLoading: (loading) => set({ loading }),
    setError: (error) => set({ error }),
}));

// Alarm Store
interface AlarmState {
    activeAlarms: ActiveAlarm[];
    alarmDefinitions: AlarmDefinition[];
    loading: boolean;
    error: string | null;
    setActiveAlarms: (alarms: ActiveAlarm[]) => void;
    setAlarmDefinitions: (definitions: AlarmDefinition[]) => void;
    acknowledgeAlarm: (alarmId: number) => void;
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
}

export const useAlarmStore = create<AlarmState>((set) => ({
    activeAlarms: [],
    alarmDefinitions: [],
    loading: false,
    error: null,
    setActiveAlarms: (alarms) => set({ activeAlarms: alarms }),
    setAlarmDefinitions: (definitions) => set({ alarmDefinitions: definitions }),
    acknowledgeAlarm: (alarmId) =>
        set((state) => ({
            activeAlarms: state.activeAlarms.map(alarm =>
                alarm.id === alarmId
                    ? { ...alarm, state: 'AckActive' as const, ack_time: new Date().toISOString() }
                    : alarm
            )
        })),
    setLoading: (loading) => set({ loading }),
    setError: (error) => set({ error }),
}));

// Audit Store
interface AuditState {
    logs: AuditLogEntry[];
    loading: boolean;
    error: string | null;
    setLogs: (logs: AuditLogEntry[]) => void;
    addLog: (log: AuditLogEntry) => void;
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
}

export const useAuditStore = create<AuditState>((set) => ({
    logs: [],
    loading: false,
    error: null,
    setLogs: (logs) => set({ logs }),
    addLog: (log) => set((state) => ({ logs: [log, ...state.logs] })),
    setLoading: (loading) => set({ loading }),
    setError: (error) => set({ error }),
}));

// System Stats Store
interface SystemState {
    stats: SystemStats;
    connected: boolean;
    lastUpdate: Date | null;
    setStats: (stats: Partial<SystemStats>) => void;
    setConnected: (connected: boolean) => void;
    setLastUpdate: (date: Date) => void;
}

export const useSystemStore = create<SystemState>((set) => ({
    stats: {
        activeAlarms: 0,
        criticalAlarms: 0,
        totalSensors: 0,
        dataPointsToday: 0,
        uptime: '0d 0h 0m',
        cpuUsage: 0,
        memoryUsage: 0,
    },
    connected: false,
    lastUpdate: null,
    setStats: (newStats) =>
        set((state) => ({
            stats: { ...state.stats, ...newStats }
        })),
    setConnected: (connected) => set({ connected }),
    setLastUpdate: (date) => set({ lastUpdate: date }),
}));

// UI Store
interface UIState {
    sidebarOpen: boolean;
    darkMode: boolean;
    activeTab: string;
    modalOpen: boolean;
    modalContent: React.ReactNode | null;
    toggleSidebar: () => void;
    setSidebarOpen: (open: boolean) => void;
    toggleDarkMode: () => void;
    setActiveTab: (tab: string) => void;
    openModal: (content: React.ReactNode) => void;
    closeModal: () => void;
}

export const useUIStore = create<UIState>((set) => ({
    sidebarOpen: true,
    darkMode: true,
    activeTab: 'dashboard',
    modalOpen: false,
    modalContent: null,
    toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
    setSidebarOpen: (open) => set({ sidebarOpen: open }),
    toggleDarkMode: () => set((state) => ({ darkMode: !state.darkMode })),
    setActiveTab: (tab) => set({ activeTab: tab }),
    openModal: (content) => set({ modalOpen: true, modalContent: content }),
    closeModal: () => set({ modalOpen: false, modalContent: null }),
}));
