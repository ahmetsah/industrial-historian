import React, { useState, useEffect, useCallback } from 'react';
import {
    Server,
    Plus,
    RefreshCw,
    Settings as SettingsIcon,
    AlertCircle,
    Database,
    Activity,
    Cpu,
} from 'lucide-react';
import { Button, Card, StatCard, EmptyState } from '../components/ui';
import { configAPI, type ModbusDevice, type DeviceStats } from '../api/configAPI';
import DeviceList from '../components/settings/DeviceList';
import DeviceForm from '../components/settings/DeviceForm';
import './Settings.css';

type SettingsTab = 'datasources' | 'system';

const Settings: React.FC = () => {
    const [activeTab, setActiveTab] = useState<SettingsTab>('datasources');
    const [devices, setDevices] = useState<ModbusDevice[]>([]);
    const [stats, setStats] = useState<DeviceStats>({ totalDevices: 0, activeDevices: 0, deployedDevices: 0, totalRegisters: 0 });
    const [isLoading, setIsLoading] = useState(true);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [apiError, setApiError] = useState<string | null>(null);
    const [showForm, setShowForm] = useState(false);
    const [editingDevice, setEditingDevice] = useState<ModbusDevice | null>(null);

    // Load devices
    const loadDevices = useCallback(async () => {
        try {
            setApiError(null);
            const fetchedDevices = await configAPI.getDevices();
            setDevices(fetchedDevices);
            setStats(configAPI.calculateStats(fetchedDevices));
        } catch (err) {
            console.error('Failed to load devices:', err);
            setApiError(err instanceof Error ? err.message : 'Cihazlar yüklenemedi');
        }
    }, []);

    // Initial load
    useEffect(() => {
        const init = async () => {
            setIsLoading(true);
            await loadDevices();
            setIsLoading(false);
        };
        init();
    }, [loadDevices]);

    // Auto refresh every 10 seconds
    useEffect(() => {
        const interval = setInterval(loadDevices, 10000);
        return () => clearInterval(interval);
    }, [loadDevices]);

    // Handle refresh
    const handleRefresh = async () => {
        setIsRefreshing(true);
        await loadDevices();
        setIsRefreshing(false);
    };

    // Handle add new device
    const handleAddDevice = () => {
        setEditingDevice(null);
        setShowForm(true);
    };

    // Handle edit device
    const handleEditDevice = async (deviceId: string) => {
        try {
            const device = await configAPI.getDevice(deviceId);
            setEditingDevice(device);
            setShowForm(true);
        } catch (err) {
            console.error('Failed to load device:', err);
            alert('Cihaz bilgileri alınamadı');
        }
    };

    // Handle delete device
    const handleDeleteDevice = async (deviceId: string, deviceName: string) => {
        if (!confirm(`"${deviceName}" cihazını silmek istediğinize emin misiniz?\n\nBu işlem:\n- Cihazı veritabanından silecek\n- Config dosyasını silecek\n- Ingestor'u durduracak\n\nBu işlem geri alınamaz!`)) {
            return;
        }

        try {
            await configAPI.deleteDevice(deviceId);
            await loadDevices();
        } catch (err) {
            console.error('Failed to delete device:', err);
            alert('Cihaz silinemedi');
        }
    };

    // Handle deploy device
    const handleDeployDevice = async (deviceId: string) => {
        try {
            const result = await configAPI.deployDevice(deviceId);
            alert(`Başarılı: ${result.message}`);
            await loadDevices();
        } catch (err) {
            console.error('Failed to deploy device:', err);
            alert('Deploy işlemi başarısız');
        }
    };

    // Handle stop device
    const handleStopDevice = async (deviceId: string, deviceName: string) => {
        if (!confirm(`"${deviceName}" cihazını durdurmak istiyor musunuz?`)) {
            return;
        }

        try {
            const result = await configAPI.stopDevice(deviceId);
            alert(`Başarılı: ${result.message}`);
            await loadDevices();
        } catch (err) {
            console.error('Failed to stop device:', err);
            alert('Durdurma işlemi başarısız');
        }
    };

    // Handle form submit (create or update)
    const handleFormSubmit = async () => {
        setShowForm(false);
        setEditingDevice(null);
        await loadDevices();
    };

    // Handle form cancel
    const handleFormCancel = () => {
        setShowForm(false);
        setEditingDevice(null);
    };

    if (isLoading) {
        return (
            <div className="settings-page">
                <div className="settings-loading">
                    <SettingsIcon className="loading-icon" />
                    <p>Ayarlar yükleniyor...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="settings-page">
            {/* Header */}
            <div className="settings-header">
                <div className="settings-title-section">
                    <h2>Ayarlar</h2>
                </div>
                <div className="settings-actions">
                    <Button
                        variant="ghost"
                        size="sm"
                        icon={RefreshCw}
                        onClick={handleRefresh}
                        loading={isRefreshing}
                    >
                        Yenile
                    </Button>
                </div>
            </div>

            {/* Tabs */}
            <div className="settings-tabs">
                <button
                    className={`settings-tab ${activeTab === 'datasources' ? 'active' : ''}`}
                    onClick={() => setActiveTab('datasources')}
                >
                    <Server className="tab-icon" />
                    Veri Kaynakları
                </button>
                <button
                    className={`settings-tab ${activeTab === 'system' ? 'active' : ''}`}
                    onClick={() => setActiveTab('system')}
                >
                    <Cpu className="tab-icon" />
                    Sistem
                </button>
            </div>

            {/* API Error */}
            {apiError && (
                <div className="api-error-banner">
                    <AlertCircle className="error-icon" />
                    <span>{apiError}</span>
                    <Button variant="ghost" size="sm" icon={RefreshCw} onClick={handleRefresh}>
                        Tekrar Dene
                    </Button>
                </div>
            )}

            {/* Content */}
            {activeTab === 'datasources' && (
                <div className="datasources-content">
                    {/* Stats */}
                    <div className="settings-stats">
                        <StatCard
                            title="Toplam Cihaz"
                            value={stats.totalDevices}
                            icon={Server}
                            status="info"
                        />
                        <StatCard
                            title="Aktif Cihaz"
                            value={stats.activeDevices}
                            icon={Activity}
                            status={stats.activeDevices > 0 ? 'normal' : 'warning'}
                        />
                        <StatCard
                            title="Toplam Register"
                            value={stats.totalRegisters}
                            icon={Database}
                            status="normal"
                        />
                    </div>

                    {/* Add Device Button */}
                    <div className="add-device-section">
                        <Button
                            variant="primary"
                            icon={Plus}
                            onClick={handleAddDevice}
                        >
                            Yeni Cihaz Ekle
                        </Button>
                    </div>

                    {/* Device List or Form */}
                    {showForm ? (
                        <DeviceForm
                            device={editingDevice}
                            onSubmit={handleFormSubmit}
                            onCancel={handleFormCancel}
                        />
                    ) : (
                        <Card className="device-list-card">
                            {devices.length === 0 ? (
                                <EmptyState
                                    icon={Server}
                                    title="Henüz cihaz yok"
                                    description="İlk Modbus cihazınızı ekleyerek başlayın"
                                    action={{
                                        label: 'Cihaz Ekle',
                                        onClick: handleAddDevice,
                                    }}
                                />
                            ) : (
                                <DeviceList
                                    devices={devices}
                                    onEdit={handleEditDevice}
                                    onDelete={handleDeleteDevice}
                                    onDeploy={handleDeployDevice}
                                    onStop={handleStopDevice}
                                />
                            )}
                        </Card>
                    )}
                </div>
            )}

            {activeTab === 'system' && (
                <div className="system-content">
                    <Card>
                        <div className="system-placeholder">
                            <Cpu className="placeholder-icon" />
                            <h3>Sistem Ayarları</h3>
                            <p>Bu bölüm yakında eklenecek.</p>
                        </div>
                    </Card>
                </div>
            )}
        </div>
    );
};

export default Settings;
