import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
    Clock,
    Download,
    Maximize2,
    Plus,
    RefreshCw,
    X,
    AlertCircle,
    Search,
} from 'lucide-react';
import { Button, Card, Skeleton } from '../components/ui';
import TrendChart from '../components/charts/TrendChart';
import { useSensorStore } from '../stores';
import { engineAPI } from '../api';
import './Trends.css';

// Chart colors
const chartColors = [
    '#00d4ed', // Cyan
    '#10b981', // Green
    '#f59e0b', // Amber
    '#8b5cf6', // Purple
    '#ec4899', // Pink
    '#3b82f6', // Blue
    '#06b6d4', // Teal
    '#ef4444', // Red
];

// Time range presets
const timeRanges = [
    { id: '1h', label: '1 Saat', hours: 1 },
    { id: '6h', label: '6 Saat', hours: 6 },
    { id: '12h', label: '12 Saat', hours: 12 },
    { id: '24h', label: '24 Saat', hours: 24 },
    { id: '7d', label: '7 Gün', hours: 168 },
    { id: '30d', label: '30 Gün', hours: 720 },
];

const TrendsPage: React.FC = () => {
    const { sensors, setSensors, selectedSensors, selectSensor, deselectSensor, sensorData, setSensorData, clearSelection } = useSensorStore();
    const [activeTimeRange, setActiveTimeRange] = useState('24h');
    const [isLoading, setIsLoading] = useState(true);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [showDataPicker, setShowDataPicker] = useState(false);
    const [apiError, setApiError] = useState<string | null>(null);
    const [searchFilter, setSearchFilter] = useState('');

    // Load sensors from API - NO auto-select
    const loadSensors = useCallback(async () => {
        try {
            setApiError(null);
            const response = await engineAPI.getMetadata();
            const fetchedSensors = response.sensors || [];
            setSensors(fetchedSensors);
        } catch (err) {
            console.error('Failed to load data sources:', err);
            setApiError(err instanceof Error ? err.message : 'Veri kaynakları yüklenemedi');
        }
    }, [setSensors]);

    // Load sensor data for selected sensors
    const loadSensorData = useCallback(async () => {
        const range = timeRanges.find(r => r.id === activeTimeRange);
        if (!range) return;

        const now = Date.now();
        const startTs = now - range.hours * 60 * 60 * 1000;

        for (const sensorId of selectedSensors) {
            try {
                const points = await engineAPI.querySensorData(sensorId, startTs, now, 1000);
                setSensorData(sensorId, points);
            } catch (err) {
                console.warn(`Failed to fetch data for ${sensorId}:`, err);
            }
        }
    }, [activeTimeRange, selectedSensors, setSensorData]);

    // Initial load
    useEffect(() => {
        const init = async () => {
            setIsLoading(true);
            await loadSensors();
            setIsLoading(false);
        };
        init();
    }, [loadSensors]);

    // Load data when sensors or time range changes
    useEffect(() => {
        if (selectedSensors.length > 0) {
            loadSensorData();
        }
    }, [selectedSensors, activeTimeRange, loadSensorData]);

    // Prepare chart data for selected sensors
    const chartData = useMemo(() => {
        return selectedSensors.map((sensorId, index) => {
            const sensor = sensors.find(s => s.id === sensorId);
            if (!sensor) return null;

            return {
                sensorId: sensor.id,
                label: sensor.desc || sensor.id,
                color: chartColors[index % chartColors.length],
                unit: sensor.unit || '',
                points: sensorData[sensor.id] || [],
            };
        }).filter(Boolean) as {
            sensorId: string;
            label: string;
            color: string;
            unit: string;
            points: { timestamp: number; value: number }[];
        }[];
    }, [selectedSensors, sensors, sensorData]);

    // Filter sensors based on search
    const filteredSensors = useMemo(() => {
        if (!searchFilter) return sensors;
        const term = searchFilter.toLowerCase();
        return sensors.filter(s =>
            s.id.toLowerCase().includes(term) ||
            s.desc.toLowerCase().includes(term) ||
            s.type.toLowerCase().includes(term)
        );
    }, [sensors, searchFilter]);

    // Handle time range change
    const handleTimeRangeChange = async (rangeId: string) => {
        setActiveTimeRange(rangeId);
    };

    // Handle refresh
    const handleRefresh = async () => {
        setIsRefreshing(true);
        await loadSensorData();
        setIsRefreshing(false);
    };

    // Handle export
    const handleExport = async () => {
        if (selectedSensors.length === 0) return;

        const range = timeRanges.find(r => r.id === activeTimeRange);
        if (!range) return;

        const now = Date.now();
        const startTs = now - range.hours * 60 * 60 * 1000;

        try {
            const blob = await engineAPI.exportCSV(selectedSensors[0], startTs, now);
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${selectedSensors[0]}_export.csv`;
            a.click();
            URL.revokeObjectURL(url);
        } catch (err) {
            console.error('Export failed:', err);
            alert('Dışa aktarma başarısız. Lütfen tekrar deneyin.');
        }
    };

    // Toggle sensor selection
    const toggleSensor = (sensorId: string) => {
        if (selectedSensors.includes(sensorId)) {
            deselectSensor(sensorId);
        } else {
            selectSensor(sensorId);
        }
    };

    // Remove a sensor
    const handleRemoveSensor = (sensorId: string) => {
        deselectSensor(sensorId);
    };

    // Group sensors by factory/line for picker
    const groupedSensors = useMemo(() => {
        const groups: Record<string, typeof sensors> = {};
        filteredSensors.forEach(sensor => {
            const key = sensor.factory && sensor.line
                ? `${sensor.factory} / ${sensor.line}`
                : 'Diğer';
            if (!groups[key]) groups[key] = [];
            groups[key].push(sensor);
        });
        return groups;
    }, [filteredSensors]);

    if (isLoading) {
        return (
            <div className="trends-page">
                <div className="trends-header">
                    <Skeleton width={200} height={32} />
                    <Skeleton width={400} height={40} />
                </div>
                <Skeleton width="100%" height={400} className="chart-skeleton" />
            </div>
        );
    }

    return (
        <div className="trends-page">
            {/* API Error */}
            {apiError && (
                <div className="api-error-banner">
                    <AlertCircle className="error-icon" />
                    <span>{apiError}</span>
                    <Button variant="ghost" size="sm" icon={RefreshCw} onClick={() => loadSensors()}>
                        Tekrar Dene
                    </Button>
                </div>
            )}

            {/* Header */}
            <div className="trends-header">
                <div className="trends-title-section">
                    <h2>Trend Analizi</h2>
                    <span className="selected-count">
                        {selectedSensors.length} veri kaynağı seçili
                    </span>
                </div>

                <div className="trends-actions">
                    {/* Time Range Selector */}
                    <div className="time-range-selector">
                        <Clock className="time-icon" />
                        {timeRanges.map(range => (
                            <button
                                key={range.id}
                                className={`time-range-btn ${activeTimeRange === range.id ? 'active' : ''}`}
                                onClick={() => handleTimeRangeChange(range.id)}
                            >
                                {range.label}
                            </button>
                        ))}
                    </div>

                    <div className="action-buttons">
                        <Button
                            variant="ghost"
                            size="sm"
                            icon={RefreshCw}
                            onClick={handleRefresh}
                            loading={isRefreshing}
                        >
                            Yenile
                        </Button>
                        <Button
                            variant="ghost"
                            size="sm"
                            icon={Download}
                            onClick={handleExport}
                            disabled={selectedSensors.length === 0}
                        >
                            Dışa Aktar
                        </Button>
                        <Button
                            variant="ghost"
                            size="sm"
                            icon={Maximize2}
                        >
                            Tam Ekran
                        </Button>
                    </div>
                </div>
            </div>

            {/* Selected Data Tags */}
            <div className="selected-sensors">
                {selectedSensors.map((sensorId, index) => {
                    const sensor = sensors.find(s => s.id === sensorId);
                    if (!sensor) return null;

                    return (
                        <div
                            key={sensorId}
                            className="sensor-tag"
                            style={{ borderColor: chartColors[index % chartColors.length] }}
                        >
                            <span
                                className="sensor-tag-color"
                                style={{ background: chartColors[index % chartColors.length] }}
                            />
                            <span className="sensor-tag-name">{sensor.desc || sensor.id}</span>
                            <span className="sensor-tag-id">{sensor.id}</span>
                            <button
                                className="sensor-tag-remove"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleRemoveSensor(sensorId);
                                }}
                                type="button"
                            >
                                <X />
                            </button>
                        </div>
                    );
                })}
                <Button
                    variant="ghost"
                    size="sm"
                    icon={Plus}
                    onClick={() => setShowDataPicker(true)}
                >
                    Veri Ekle
                </Button>
                {selectedSensors.length > 0 && (
                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => clearSelection()}
                    >
                        Tümünü Temizle
                    </Button>
                )}
            </div>

            {/* Main Chart */}
            <Card className="main-chart-card">
                {selectedSensors.length === 0 ? (
                    <div className="no-sensors-selected">
                        <p>Görüntülenecek veri seçilmedi</p>
                        <Button
                            variant="primary"
                            icon={Plus}
                            onClick={() => setShowDataPicker(true)}
                        >
                            Veri Kaynağı Seç
                        </Button>
                    </div>
                ) : (
                    <TrendChart
                        data={chartData}
                        height={450}
                        showLegend={false}
                        enableZoom
                        enablePan
                    />
                )}

                {/* Live indicator */}
                {selectedSensors.length > 0 && (
                    <div className="chart-live-indicator">
                        <span className="live-dot" />
                        <span>CANLI</span>
                    </div>
                )}
            </Card>

            {/* Sensor Statistics */}
            {selectedSensors.length > 0 && (
                <div className="sensor-statistics">
                    {chartData.map((sensor) => {
                        const values = sensor.points.map(p => p.value);
                        const min = values.length > 0 ? Math.min(...values) : 0;
                        const max = values.length > 0 ? Math.max(...values) : 0;
                        const avg = values.length > 0 ? values.reduce((a, b) => a + b, 0) / values.length : 0;
                        const current = values.length > 0 ? values[values.length - 1] : 0;

                        return (
                            <Card key={sensor.sensorId} className="stat-item">
                                <div className="stat-header">
                                    <span
                                        className="stat-color"
                                        style={{ background: sensor.color }}
                                    />
                                    <span className="stat-name">{sensor.label}</span>
                                </div>
                                <div className="stat-grid">
                                    <div className="stat-cell">
                                        <span className="stat-label">Güncel</span>
                                        <span className="stat-value">{current.toFixed(2)}</span>
                                    </div>
                                    <div className="stat-cell">
                                        <span className="stat-label">Ortalama</span>
                                        <span className="stat-value">{avg.toFixed(2)}</span>
                                    </div>
                                    <div className="stat-cell">
                                        <span className="stat-label">Min</span>
                                        <span className="stat-value">{min.toFixed(2)}</span>
                                    </div>
                                    <div className="stat-cell">
                                        <span className="stat-label">Max</span>
                                        <span className="stat-value">{max.toFixed(2)}</span>
                                    </div>
                                </div>
                            </Card>
                        );
                    })}
                </div>
            )}

            {/* Data Picker Modal */}
            {showDataPicker && (
                <div className="sensor-picker-overlay" onClick={() => setShowDataPicker(false)}>
                    <div className="sensor-picker" onClick={e => e.stopPropagation()}>
                        <div className="sensor-picker-header">
                            <h3>Veri Kaynağı Seç</h3>
                            <button className="close-btn" onClick={() => setShowDataPicker(false)}>
                                <X />
                            </button>
                        </div>

                        {/* Search Filter */}
                        <div className="sensor-picker-search">
                            <Search className="search-icon" />
                            <input
                                type="text"
                                placeholder="ID, açıklama veya tip ile filtrele..."
                                value={searchFilter}
                                onChange={(e) => setSearchFilter(e.target.value)}
                                autoFocus
                            />
                            {searchFilter && (
                                <button className="clear-search" onClick={() => setSearchFilter('')}>
                                    <X />
                                </button>
                            )}
                        </div>

                        <div className="sensor-picker-info">
                            {filteredSensors.length} / {sensors.length} veri kaynağı
                        </div>

                        <div className="sensor-picker-content">
                            {Object.keys(groupedSensors).length === 0 ? (
                                <div className="no-sensors-message">
                                    <p>
                                        {searchFilter
                                            ? 'Filtre ile eşleşen veri kaynağı bulunamadı.'
                                            : 'Veri kaynağı bulunamadı. Backend bağlantısını kontrol edin.'}
                                    </p>
                                </div>
                            ) : (
                                Object.entries(groupedSensors).map(([group, groupSensors]) => (
                                    <div key={group} className="sensor-group">
                                        <h4 className="sensor-group-title">{group}</h4>
                                        <div className="sensor-group-items">
                                            {groupSensors.map(sensor => (
                                                <div
                                                    key={sensor.id}
                                                    className={`sensor-picker-item ${selectedSensors.includes(sensor.id) ? 'selected' : ''}`}
                                                    onClick={() => toggleSensor(sensor.id)}
                                                >
                                                    <div className="picker-item-check">
                                                        {selectedSensors.includes(sensor.id) && '✓'}
                                                    </div>
                                                    <div className="picker-item-info">
                                                        <span className="picker-item-name">{sensor.desc || sensor.id}</span>
                                                        <span className="picker-item-id">{sensor.id}</span>
                                                    </div>
                                                    <span className="picker-item-type">{sensor.type || 'Genel'}</span>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                ))
                            )}
                        </div>
                        <div className="sensor-picker-footer">
                            <span className="selected-info">{selectedSensors.length} seçili</span>
                            <Button variant="secondary" onClick={() => setShowDataPicker(false)}>
                                Kapat
                            </Button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default TrendsPage;
