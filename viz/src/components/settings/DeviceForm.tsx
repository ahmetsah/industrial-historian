import React, { useState, useEffect } from 'react';
import { Plus, X, Server, Save, HelpCircle } from 'lucide-react';
import { Button, Card } from '../ui';
import { configAPI, type ModbusDevice, type ModbusDeviceBase, type Register } from '../../api/configAPI';

interface DeviceFormProps {
    device: ModbusDevice | null;
    onSubmit: () => void;
    onCancel: () => void;
}

// Protocol options
const protocols = [
    { value: 'modbus', label: 'Modbus TCP', enabled: true },
    { value: 's7net', label: 'S7Net (Siemens)', enabled: false },
    { value: 'opcua', label: 'OPC UA Client', enabled: false },
] as const;

// All supported Modbus data types with descriptions
const dataTypes = [
    { value: 'Bool', label: 'Bool', registers: 1, description: '0 veya 1' },
    { value: 'Int16', label: 'Int16', registers: 1, description: '-32768 to 32767' },
    { value: 'UInt16', label: 'UInt16', registers: 1, description: '0 to 65535' },
    { value: 'Int32', label: 'Int32', registers: 2, description: '-2B to 2B' },
    { value: 'UInt32', label: 'UInt32', registers: 2, description: '0 to 4B' },
    { value: 'Float32', label: 'Float32', registers: 2, description: 'IEEE 754 float' },
    { value: 'Int64', label: 'Int64', registers: 4, description: '64-bit signed' },
    { value: 'UInt64', label: 'UInt64', registers: 4, description: '64-bit unsigned' },
    { value: 'Float64', label: 'Float64', registers: 4, description: 'IEEE 754 double' },
] as const;

// Register type options
const registerTypes = [
    { value: 'holding', label: 'Holding Register (4x)', prefix: 4 },
    { value: 'input', label: 'Input Register (3x)', prefix: 3 },
    { value: 'coil', label: 'Coil (0x)', prefix: 0 },
    { value: 'discrete', label: 'Discrete Input (1x)', prefix: 1 },
] as const;

const DeviceForm: React.FC<DeviceFormProps> = ({ device, onSubmit, onCancel }) => {
    const isEditing = !!device;

    // Form state
    const [protocol, setProtocol] = useState<string>('modbus');
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [ip, setIp] = useState('');
    const [port, setPort] = useState(502);
    const [unitId, setUnitId] = useState(1);
    const [pollInterval, setPollInterval] = useState(1000);
    const [registers, setRegisters] = useState<(Register & { registerType?: string })[]>([
        { address: 0, name: '', data_type: 'Float32', registerType: 'holding' },
    ]);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Populate form when editing
    useEffect(() => {
        if (device) {
            setName(device.device.name || '');
            setDescription(device.device.description || '');
            setIp(device.ip || '');
            setPort(device.port || 502);
            setUnitId(device.unit_id || 1);
            setPollInterval(device.poll_interval_ms || 1000);
            setRegisters(device.registers?.length > 0
                ? device.registers.map(r => ({ ...r, registerType: 'holding' }))
                : [{ address: 0, name: '', data_type: 'Float32', registerType: 'holding' }]);
        }
    }, [device]);

    // Add register
    const addRegister = () => {
        const lastReg = registers[registers.length - 1];
        const lastType = dataTypes.find(t => t.value === lastReg?.data_type);
        const nextAddress = (lastReg?.address || 0) + (lastType?.registers || 1);

        setRegisters([...registers, {
            address: nextAddress,
            name: '',
            data_type: 'Float32',
            registerType: 'holding'
        }]);
    };

    // Remove register
    const removeRegister = (index: number) => {
        if (registers.length > 1) {
            setRegisters(registers.filter((_, i) => i !== index));
        }
    };

    // Update register field
    const updateRegister = (index: number, field: keyof Register | 'registerType', value: string | number) => {
        const updated = [...registers];
        if (field === 'address') {
            updated[index] = { ...updated[index], address: Number(value) };
        } else if (field === 'data_type') {
            updated[index] = { ...updated[index], data_type: value as Register['data_type'] };
        } else if (field === 'name') {
            updated[index] = { ...updated[index], name: String(value) };
        } else if (field === 'registerType') {
            updated[index] = { ...updated[index], registerType: String(value) };
        }
        setRegisters(updated);
    };

    // Auto-generate tag name based on device and address
    const generateTagName = (index: number) => {
        const reg = registers[index];
        const tagBase = name.replace(/\s+/g, '') || 'Device';
        const tagName = `${tagBase}.Tag${reg.address}`;
        updateRegister(index, 'name', tagName);
    };

    // Handle submit
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        // Validate
        if (!name.trim()) {
            setError('Cihaz adı zorunludur');
            setIsSubmitting(false);
            return;
        }

        if (!ip.trim()) {
            setError('IP adresi zorunludur');
            setIsSubmitting(false);
            return;
        }

        // Filter out empty registers and remove registerType
        const validRegisters = registers
            .filter(r => r.name.trim())
            .map(({ registerType, ...rest }) => rest);

        const deviceData: ModbusDeviceBase = {
            name: name.trim(),
            description: description.trim() || undefined,
            ip: ip.trim(),
            port,
            unit_id: unitId,
            poll_interval_ms: pollInterval,
            registers: validRegisters,
        };

        try {
            if (isEditing && device) {
                await configAPI.updateDevice(device.device.id, deviceData);
            } else {
                await configAPI.createDevice(deviceData);
            }
            onSubmit();
        } catch (err) {
            console.error('Failed to save device:', err);
            setError(err instanceof Error ? err.message : 'Cihaz kaydedilemedi');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Card className="device-form-card">
            <div className="device-form-header">
                <h3>
                    <Server size={20} />
                    {isEditing ? `Düzenle: ${device?.device.name}` : 'Yeni Cihaz Ekle'}
                </h3>
                <Button variant="ghost" size="sm" onClick={onCancel}>
                    İptal
                </Button>
            </div>

            {error && (
                <div className="api-error-banner" style={{ marginBottom: '1rem' }}>
                    <span>{error}</span>
                </div>
            )}

            <form onSubmit={handleSubmit}>
                {/* Protocol Selection */}
                <div className="protocol-selection">
                    <label>Haberleşme Protokolü</label>
                    <div className="protocol-options">
                        {protocols.map(p => (
                            <button
                                key={p.value}
                                type="button"
                                className={`protocol-option ${protocol === p.value ? 'active' : ''} ${!p.enabled ? 'disabled' : ''}`}
                                onClick={() => p.enabled && setProtocol(p.value)}
                                disabled={!p.enabled}
                            >
                                {p.label}
                                {!p.enabled && <span className="coming-soon">Yakında</span>}
                            </button>
                        ))}
                    </div>
                </div>

                <div className="form-grid">
                    <div className="form-group">
                        <label>Cihaz Adı *</label>
                        <input
                            type="text"
                            value={name}
                            onChange={e => setName(e.target.value)}
                            placeholder="PLC-001"
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Açıklama</label>
                        <input
                            type="text"
                            value={description}
                            onChange={e => setDescription(e.target.value)}
                            placeholder="Ana üretim hattı PLC"
                        />
                    </div>

                    <div className="form-group">
                        <label>IP Adresi *</label>
                        <input
                            type="text"
                            value={ip}
                            onChange={e => setIp(e.target.value)}
                            placeholder="192.168.1.10"
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Port</label>
                        <input
                            type="number"
                            value={port}
                            onChange={e => setPort(Number(e.target.value))}
                            min={1}
                            max={65535}
                        />
                    </div>

                    <div className="form-group">
                        <label>Unit ID (Slave)</label>
                        <input
                            type="number"
                            value={unitId}
                            onChange={e => setUnitId(Number(e.target.value))}
                            min={0}
                            max={255}
                        />
                    </div>

                    <div className="form-group">
                        <label>Poll Interval (ms)</label>
                        <input
                            type="number"
                            value={pollInterval}
                            onChange={e => setPollInterval(Number(e.target.value))}
                            min={100}
                            step={100}
                        />
                    </div>
                </div>

                {/* Registers Section */}
                <div className="registers-section">
                    <div className="registers-header">
                        <h4>Registers</h4>
                        <div className="registers-header-actions">
                            <span className="register-count">{registers.length} register</span>
                            <Button variant="ghost" size="sm" icon={Plus} type="button" onClick={addRegister}>
                                Register Ekle
                            </Button>
                        </div>
                    </div>

                    {/* Register Header */}
                    <div className="register-header-row">
                        <span className="reg-col-type">Tip</span>
                        <span className="reg-col-addr">Adres</span>
                        <span className="reg-col-name">Tag Adı</span>
                        <span className="reg-col-dtype">Veri Tipi</span>
                        <span className="reg-col-action"></span>
                    </div>

                    <div className="register-list">
                        {registers.map((register, index) => (
                            <div key={index} className="register-item">
                                <select
                                    className="reg-col-type"
                                    value={register.registerType || 'holding'}
                                    onChange={e => updateRegister(index, 'registerType', e.target.value)}
                                    title="Register Tipi"
                                >
                                    {registerTypes.map(rt => (
                                        <option key={rt.value} value={rt.value}>{rt.label}</option>
                                    ))}
                                </select>

                                <input
                                    className="reg-col-addr"
                                    type="number"
                                    value={register.address}
                                    onChange={e => updateRegister(index, 'address', e.target.value)}
                                    placeholder="0"
                                    min={0}
                                    title="Register Adresi"
                                />

                                <div className="reg-col-name-wrapper">
                                    <input
                                        className="reg-col-name"
                                        type="text"
                                        value={register.name}
                                        onChange={e => updateRegister(index, 'name', e.target.value)}
                                        placeholder="Factory1.Line1.PLC001.Temperature.T001"
                                        title="Tag Adı (Sensör ID)"
                                    />
                                    <button
                                        type="button"
                                        className="auto-name-btn"
                                        onClick={() => generateTagName(index)}
                                        title="Otomatik İsim Oluştur"
                                    >
                                        <HelpCircle size={14} />
                                    </button>
                                </div>

                                <select
                                    className="reg-col-dtype"
                                    value={register.data_type}
                                    onChange={e => updateRegister(index, 'data_type', e.target.value)}
                                    title="Veri Tipi"
                                >
                                    {dataTypes.map(type => (
                                        <option key={type.value} value={type.value}>
                                            {type.label} ({type.registers}R)
                                        </option>
                                    ))}
                                </select>

                                <button
                                    type="button"
                                    className="register-remove"
                                    onClick={() => removeRegister(index)}
                                    disabled={registers.length <= 1}
                                    title="Sil"
                                >
                                    <X size={14} />
                                </button>
                            </div>
                        ))}
                    </div>

                    {/* Quick Add Section */}
                    <div className="quick-add-section">
                        <span className="quick-add-label">Hızlı Ekle:</span>
                        <button type="button" className="quick-add-btn" onClick={() => {
                            const last = registers[registers.length - 1];
                            const nextAddr = (last?.address || 0) + 2;
                            setRegisters([...registers,
                            { address: nextAddr, name: `${name || 'Device'}.Tag${nextAddr}`, data_type: 'Float32', registerType: 'holding' }
                            ]);
                        }}>+1 Float32</button>
                        <button type="button" className="quick-add-btn" onClick={() => {
                            const last = registers[registers.length - 1];
                            let addr = (last?.address || 0) + 2;
                            const newRegs = [];
                            for (let i = 0; i < 5; i++) {
                                newRegs.push({ address: addr, name: `${name || 'Device'}.Tag${addr}`, data_type: 'Float32' as const, registerType: 'holding' });
                                addr += 2;
                            }
                            setRegisters([...registers, ...newRegs]);
                        }}>+5 Float32</button>
                        <button type="button" className="quick-add-btn" onClick={() => {
                            const last = registers[registers.length - 1];
                            let addr = (last?.address || 0) + 1;
                            const newRegs = [];
                            for (let i = 0; i < 10; i++) {
                                newRegs.push({ address: addr, name: `${name || 'Device'}.Tag${addr}`, data_type: 'Int16' as const, registerType: 'holding' });
                                addr += 1;
                            }
                            setRegisters([...registers, ...newRegs]);
                        }}>+10 Int16</button>
                    </div>
                </div>

                {/* Form Actions */}
                <div className="form-actions">
                    <Button
                        type="submit"
                        variant="primary"
                        icon={Save}
                        loading={isSubmitting}
                    >
                        {isEditing ? 'Güncelle' : 'Oluştur'}
                    </Button>
                    <Button type="button" variant="secondary" onClick={onCancel}>
                        İptal
                    </Button>
                </div>
            </form>
        </Card>
    );
};

export default DeviceForm;
