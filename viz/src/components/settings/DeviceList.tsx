import React from 'react';
import { Edit2, Rocket, Trash2, Database, Square } from 'lucide-react';
import type { ModbusDevice } from '../../api/configAPI';

interface DeviceListProps {
    devices: ModbusDevice[];
    onEdit: (deviceId: string) => void;
    onDelete: (deviceId: string, deviceName: string) => void;
    onDeploy: (deviceId: string) => void;
    onStop: (deviceId: string, deviceName: string) => void;
}

const DeviceList: React.FC<DeviceListProps> = ({
    devices,
    onEdit,
    onDelete,
    onDeploy,
    onStop,
}) => {
    const getDeployStatus = (device: ModbusDevice) => {
        const status = device.deployment_status || 'not_deployed';
        return status === 'deployed' ? 'deployed' : 'not-deployed';
    };

    const getDeployText = (device: ModbusDevice) => {
        const status = device.deployment_status || 'not_deployed';
        return status === 'deployed' ? 'Aktif' : 'Pasif';
    };

    const getConnectionStatus = (device: ModbusDevice) => {
        const status = device.connection_status || 'idle';
        if (status === 'connected') return 'connected';
        if (status === 'disconnected') return 'disconnected';
        return 'idle';
    };

    const getConnectionText = (device: ModbusDevice) => {
        const status = device.connection_status || 'idle';
        if (status === 'connected') return 'Bağlı';
        if (status === 'disconnected') return 'Bağlı Değil';
        return 'Beklemede';
    };

    const isDeployed = (device: ModbusDevice) => {
        return device.deployment_status === 'deployed';
    };

    return (
        <table className="device-table">
            <thead>
                <tr>
                    <th>Cihaz</th>
                    <th>IP Adresi</th>
                    <th>Poll Interval</th>
                    <th>Registers</th>
                    <th>Durum</th>
                    <th>İşlemler</th>
                </tr>
            </thead>
            <tbody>
                {devices.map(device => (
                    <tr key={device.device.id}>
                        <td>
                            <div className="device-name">{device.device.name}</div>
                            {device.device.description && (
                                <div className="device-desc">{device.device.description}</div>
                            )}
                        </td>
                        <td>
                            <span className="device-ip">
                                {device.ip}:{device.port}
                            </span>
                            <div className="device-desc">Unit ID: {device.unit_id}</div>
                        </td>
                        <td>{device.poll_interval_ms}ms</td>
                        <td>
                            <span className="device-registers">
                                <Database size={14} />
                                {device.registers?.length || 0}
                            </span>
                        </td>
                        <td>
                            <div className="status-badges">
                                <span className={`device-status-badge ${getDeployStatus(device)}`}>
                                    {getDeployText(device)}
                                </span>
                                <span className={`device-status-badge ${getConnectionStatus(device)}`}>
                                    {getConnectionText(device)}
                                </span>
                            </div>
                        </td>
                        <td>
                            <div className="device-actions">
                                <button
                                    className="action-btn edit"
                                    onClick={() => onEdit(device.device.id)}
                                    title="Düzenle"
                                >
                                    <Edit2 />
                                </button>
                                {isDeployed(device) ? (
                                    <button
                                        className="action-btn stop"
                                        onClick={() => onStop(device.device.id, device.device.name)}
                                        title="Durdur"
                                    >
                                        <Square />
                                    </button>
                                ) : (
                                    <button
                                        className="action-btn deploy"
                                        onClick={() => onDeploy(device.device.id)}
                                        title="Başlat"
                                    >
                                        <Rocket />
                                    </button>
                                )}
                                <button
                                    className="action-btn delete"
                                    onClick={() => onDelete(device.device.id, device.device.name)}
                                    title="Sil"
                                >
                                    <Trash2 />
                                </button>
                            </div>
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
};

export default DeviceList;
