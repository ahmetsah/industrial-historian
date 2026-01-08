import React, { useEffect, useState, useCallback } from 'react';
import {
    Activity,
    AlertTriangle,
    BarChart3,
    Database,
    Gauge,
    TrendingUp,
    Zap,
    RefreshCw,
} from 'lucide-react';
import { StatCard, Card, AlarmRow, SystemHealth, Button } from '../components/ui';
import TrendChart from '../components/charts/TrendChart';
import { useSensorStore, useAlarmStore, useSystemStore } from '../stores';
import { engineAPI, alarmAPI } from '../api';
import './Dashboard.css';

// Color palette for charts
const chartColors = [
    '#00d4ed', // Cyan
    '#10b981', // Green
    '#f59e0b', // Amber
    '#8b5cf6', // Purple
    '#ec4899', // Pink
    '#3b82f6', // Blue
];

const Dashboard: React.FC = () => {
    const { sensors, setSensors, sensorData, setSensorData, setError: setSensorError } = useSensorStore();
    const { activeAlarms, setActiveAlarms, setError: setAlarmError } = useAlarmStore();
    const { stats, setStats, connected, setConnected } = useSystemStore();
    const [isLoading, setIsLoading] = useState(true);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [apiError, setApiError] = useState<string | null>(null);

    // Load data from real APIs
    const loadData = useCallback(async () => {
        setApiError(null);

        try {
            // Fetch sensor metadata from Engine API
            const metadataResponse = await engineAPI.getMetadata();
            const fetchedSensors = metadataResponse.sensors || [];
            setSensors(fetchedSensors);

            // Fetch sensor data for each sensor (last 24 hours)
            const now = Date.now();
            const oneDayAgo = now - 24 * 60 * 60 * 1000;

            for (const sensor of fetchedSensors.slice(0, 6)) {
                try {
                    const points = await engineAPI.querySensorData(
                        sensor.id,
                        oneDayAgo,
                        now,
                        1000
                    );
                    setSensorData(sensor.id, points);
                } catch (err) {
                    console.warn(`Failed to fetch data for sensor ${sensor.id}:`, err);
                }
            }

            // Fetch active alarms from Alarm API
            const alarms = await alarmAPI.getActiveAlarms();
            setActiveAlarms(alarms);

            // Calculate stats
            setStats({
                activeAlarms: alarms.length,
                criticalAlarms: alarms.filter(a => a.priority === 'Critical').length,
                totalSensors: fetchedSensors.length,
                dataPointsToday: 0, // Would need a separate API call
                uptime: '0d 0h 0m', // Would need a separate API call
                cpuUsage: 0,
                memoryUsage: 0,
            });

            setConnected(true);
        } catch (err) {
            console.error('Failed to load dashboard data:', err);
            setApiError(err instanceof Error ? err.message : 'Failed to connect to backend services');
            setSensorError('Failed to fetch sensors');
            setAlarmError('Failed to fetch alarms');
            setConnected(false);
        }
    }, [setSensors, setSensorData, setActiveAlarms, setStats, setConnected, setSensorError, setAlarmError]);

    // Initial load
    useEffect(() => {
        const init = async () => {
            setIsLoading(true);
            await loadData();
            setIsLoading(false);
        };
        init();
    }, [loadData]);

    // Refresh handler
    const handleRefresh = async () => {
        setIsRefreshing(true);
        await loadData();
        setIsRefreshing(false);
    };

    // Get selected sensors for main chart (first 3)
    const chartSensors = sensors.slice(0, 3).map((sensor, index) => ({
        sensorId: sensor.id,
        label: sensor.desc,
        color: chartColors[index % chartColors.length],
        unit: sensor.unit,
        points: sensorData[sensor.id] || [],
    }));

    // Get current values for gauge display
    const getCurrentValue = (sensorId: string): number => {
        const data = sensorData[sensorId];
        if (!data || data.length === 0) return 0;
        return data[data.length - 1].value;
    };

    if (isLoading) {
        return (
            <div className="dashboard-loading">
                <div className="loading-spinner">
                    <Zap className="spinner-icon" />
                </div>
                <p>Loading dashboard...</p>
            </div>
        );
    }

    return (
        <div className="dashboard">
            {/* API Error Banner */}
            {apiError && (
                <div className="api-error-banner">
                    <AlertTriangle className="error-icon" />
                    <span>{apiError}</span>
                    <Button variant="ghost" size="sm" icon={RefreshCw} onClick={handleRefresh}>
                        Retry
                    </Button>
                </div>
            )}

            {/* Stats Overview */}
            <section className="dashboard-stats">
                <StatCard
                    title="Active Alarms"
                    value={stats.activeAlarms}
                    icon={AlertTriangle}
                    status={stats.criticalAlarms > 0 ? 'critical' : stats.activeAlarms > 0 ? 'warning' : 'normal'}
                    animate
                />
                <StatCard
                    title="Total Sensors"
                    value={stats.totalSensors}
                    icon={Gauge}
                    status="info"
                    animate
                />
                <StatCard
                    title="Data Points Today"
                    value={formatNumber(stats.dataPointsToday)}
                    icon={Database}
                    status="normal"
                    trend={{ value: 12, direction: 'up' }}
                    animate
                />
                <StatCard
                    title="System Uptime"
                    value={stats.uptime}
                    icon={Activity}
                    status="normal"
                    animate
                />
            </section>

            {/* Main Chart */}
            <section className="dashboard-chart-section">
                <div className="chart-header-actions">
                    <Button
                        variant="ghost"
                        size="sm"
                        icon={RefreshCw}
                        onClick={handleRefresh}
                        loading={isRefreshing}
                    >
                        Refresh
                    </Button>
                </div>
                <TrendChart
                    title="Real-Time Process Trends"
                    data={chartSensors}
                    height={350}
                    showLegend
                    enableZoom
                />
            </section>

            {/* Two Column Layout */}
            <section className="dashboard-grid">
                {/* Sensor Values */}
                <Card className="sensor-values-card">
                    <div className="card-header">
                        <h3><TrendingUp className="header-icon" /> Live Sensor Values</h3>
                    </div>
                    {sensors.length === 0 ? (
                        <div className="no-data">
                            <p>No sensors available. Please check backend connection.</p>
                        </div>
                    ) : (
                        <div className="sensor-values-grid">
                            {sensors.slice(0, 6).map((sensor, index) => {
                                const value = getCurrentValue(sensor.id);
                                const prevValue = sensorData[sensor.id]?.[sensorData[sensor.id].length - 2]?.value || value;
                                const trend = value > prevValue ? 'up' : value < prevValue ? 'down' : 'stable';

                                return (
                                    <div
                                        key={sensor.id}
                                        className="sensor-value-item"
                                        style={{ animationDelay: `${index * 50}ms` }}
                                    >
                                        <div className="sensor-value-header">
                                            <span className="sensor-type">{sensor.type}</span>
                                            <span className={`sensor-trend trend-${trend}`}>
                                                {trend === 'up' ? '↑' : trend === 'down' ? '↓' : '→'}
                                            </span>
                                        </div>
                                        <div className="sensor-value-main">
                                            <span className="sensor-value">{value.toFixed(1)}</span>
                                            <span className="sensor-unit">{sensor.unit}</span>
                                        </div>
                                        <div className="sensor-value-footer">
                                            <span className="sensor-name">{sensor.machine}</span>
                                            <span className="sensor-line">{sensor.line}</span>
                                        </div>
                                        {/* Mini sparkline visualization */}
                                        <div className="sensor-sparkline">
                                            <SparkLine
                                                data={sensorData[sensor.id]?.slice(-20) || []}
                                                color={chartColors[index % chartColors.length]}
                                            />
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </Card>

                {/* Active Alarms */}
                <Card className="alarms-card">
                    <div className="card-header">
                        <h3><AlertTriangle className="header-icon" /> Active Alarms</h3>
                        <span className="alarm-count">
                            {stats.criticalAlarms} Critical / {stats.activeAlarms} Total
                        </span>
                    </div>
                    <div className="alarms-list">
                        {activeAlarms.length === 0 ? (
                            <div className="no-alarms">
                                <Zap className="no-alarms-icon" />
                                <p>No active alarms</p>
                            </div>
                        ) : (
                            activeAlarms.map((alarm) => (
                                <AlarmRow
                                    key={alarm.id}
                                    id={alarm.id}
                                    tag={alarm.tag || `Tag-${alarm.definition_id}`}
                                    value={alarm.value}
                                    state={alarm.state}
                                    priority={alarm.priority || 'Warning'}
                                    activationTime={alarm.activation_time}
                                />
                            ))
                        )}
                    </div>
                </Card>
            </section>

            {/* Bottom Section */}
            <section className="dashboard-bottom">
                <SystemHealth
                    cpuUsage={stats.cpuUsage}
                    memoryUsage={stats.memoryUsage}
                    connected={connected}
                />

                <Card className="quick-stats-card">
                    <div className="card-header">
                        <h3><BarChart3 className="header-icon" /> Performance Metrics</h3>
                    </div>
                    <div className="performance-metrics">
                        <div className="metric-item">
                            <span className="metric-label">Write Speed</span>
                            <span className="metric-value">-- <span className="metric-suffix">events/s</span></span>
                        </div>
                        <div className="metric-item">
                            <span className="metric-label">Query Latency</span>
                            <span className="metric-value">-- <span className="metric-suffix">ms</span></span>
                        </div>
                        <div className="metric-item">
                            <span className="metric-label">Storage Used</span>
                            <span className="metric-value">-- <span className="metric-suffix">TB</span></span>
                        </div>
                        <div className="metric-item">
                            <span className="metric-label">Compression Ratio</span>
                            <span className="metric-value">--</span>
                        </div>
                    </div>
                </Card>
            </section>
        </div>
    );
};

// Helper function to format large numbers
function formatNumber(num: number): string {
    if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
}

// Simple sparkline component
interface SparkLineProps {
    data: { timestamp: number; value: number }[];
    color: string;
}

const SparkLine: React.FC<SparkLineProps> = ({ data, color }) => {
    if (data.length < 2) return null;

    const values = data.map(d => d.value);
    const min = Math.min(...values);
    const max = Math.max(...values);
    const range = max - min || 1;

    const width = 100;
    const height = 24;
    const points = values.map((v, i) => {
        const x = (i / (values.length - 1)) * width;
        const y = height - ((v - min) / range) * height;
        return `${x},${y}`;
    }).join(' ');

    return (
        <svg width={width} height={height} className="sparkline-svg">
            <polyline
                points={points}
                fill="none"
                stroke={color}
                strokeWidth="1.5"
                strokeLinecap="round"
                strokeLinejoin="round"
            />
        </svg>
    );
};

export default Dashboard;
