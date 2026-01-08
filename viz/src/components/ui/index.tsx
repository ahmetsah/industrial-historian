import React from 'react';
import {
    Activity,
    AlertTriangle,
    Bell,
    CheckCircle,
    ChevronRight,
    Clock,
    Cpu,
    Database,
    HardDrive,
    Loader2,
    type LucideIcon,
} from 'lucide-react';
import './ui.css';

// Button Component
interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'success';
    size?: 'sm' | 'md' | 'lg';
    loading?: boolean;
    icon?: LucideIcon;
    iconPosition?: 'left' | 'right';
}

export const Button: React.FC<ButtonProps> = ({
    children,
    variant = 'primary',
    size = 'md',
    loading = false,
    icon: Icon,
    iconPosition = 'left',
    className = '',
    disabled,
    ...props
}) => {
    return (
        <button
            className={`btn btn-${variant} btn-${size} ${className}`}
            disabled={disabled || loading}
            {...props}
        >
            {loading ? (
                <Loader2 className="btn-icon animate-spin" />
            ) : (
                Icon && iconPosition === 'left' && <Icon className="btn-icon" />
            )}
            {children}
            {!loading && Icon && iconPosition === 'right' && <Icon className="btn-icon" />}
        </button>
    );
};

// Card Component
interface CardProps {
    children: React.ReactNode;
    className?: string;
    padding?: 'none' | 'sm' | 'md' | 'lg';
    hover?: boolean;
    onClick?: () => void;
}

export const Card: React.FC<CardProps> = ({
    children,
    className = '',
    padding = 'md',
    hover = false,
    onClick,
}) => {
    return (
        <div
            className={`card-component padding-${padding} ${hover ? 'card-hover' : ''} ${className}`}
            onClick={onClick}
            role={onClick ? 'button' : undefined}
            tabIndex={onClick ? 0 : undefined}
        >
            {children}
        </div>
    );
};

// Stat Card Component
interface StatCardProps {
    title: string;
    value: string | number;
    unit?: string;
    icon: LucideIcon;
    trend?: { value: number; direction: 'up' | 'down' };
    status?: 'normal' | 'warning' | 'critical' | 'info';
    animate?: boolean;
}

export const StatCard: React.FC<StatCardProps> = ({
    title,
    value,
    unit,
    icon: Icon,
    trend,
    status = 'normal',
    animate = false,
}) => {
    return (
        <div className={`stat-card status-${status} ${animate ? 'animate-slide-up' : ''}`}>
            <div className="stat-card-header">
                <span className="stat-card-title">{title}</span>
                <div className={`stat-card-icon-wrapper status-${status}`}>
                    <Icon className="stat-card-icon" />
                </div>
            </div>
            <div className="stat-card-body">
                <span className="stat-card-value">{value}</span>
                {unit && <span className="stat-card-unit">{unit}</span>}
            </div>
            {trend && (
                <div className={`stat-card-trend trend-${trend.direction}`}>
                    <span>{trend.direction === 'up' ? '↑' : '↓'}</span>
                    <span>{Math.abs(trend.value)}%</span>
                </div>
            )}
        </div>
    );
};

// Status Badge Component
interface StatusBadgeProps {
    status: 'normal' | 'warning' | 'critical' | 'shelved' | 'acknowledged';
    text?: string;
    pulse?: boolean;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({
    status,
    text,
    pulse = false,
}) => {
    const statusText = text || {
        normal: 'Normal',
        warning: 'Warning',
        critical: 'Critical',
        shelved: 'Shelved',
        acknowledged: 'Ack',
    }[status];

    return (
        <span className={`status-badge badge-${status} ${pulse ? 'pulse' : ''}`}>
            {status === 'critical' && <AlertTriangle className="badge-icon" />}
            {status === 'warning' && <Bell className="badge-icon" />}
            {status === 'normal' && <CheckCircle className="badge-icon" />}
            {status === 'shelved' && <Clock className="badge-icon" />}
            {statusText}
        </span>
    );
};

// Alarm Row Component
interface AlarmRowProps {
    id: number;
    tag: string;
    value: number;
    state: string;
    priority: string;
    activationTime: string;
    onAcknowledge?: () => void;
    onShelve?: () => void;
}

export const AlarmRow: React.FC<AlarmRowProps> = ({
    tag,
    value,
    state,
    priority,
    activationTime,
    onAcknowledge,
    onShelve,
}) => {
    const getStatusFromState = (state: string): 'normal' | 'warning' | 'critical' | 'shelved' | 'acknowledged' => {
        if (state === 'Shelved') return 'shelved';
        if (state.includes('Ack')) return 'acknowledged';
        if (priority === 'Critical') return 'critical';
        return 'warning';
    };

    const timeAgo = (dateStr: string) => {
        const diff = Date.now() - new Date(dateStr).getTime();
        const minutes = Math.floor(diff / 60000);
        if (minutes < 60) return `${minutes}m ago`;
        const hours = Math.floor(minutes / 60);
        if (hours < 24) return `${hours}h ago`;
        return `${Math.floor(hours / 24)}d ago`;
    };

    return (
        <div className={`alarm-row priority-${priority.toLowerCase()}`}>
            <div className="alarm-row-main">
                <StatusBadge
                    status={getStatusFromState(state)}
                    pulse={state === 'UnackActive' && priority === 'Critical'}
                />
                <div className="alarm-row-details">
                    <span className="alarm-tag">{tag}</span>
                    <span className="alarm-meta">
                        <span className="alarm-value">{value.toFixed(2)}</span>
                        <span className="alarm-time">{timeAgo(activationTime)}</span>
                    </span>
                </div>
            </div>
            <div className="alarm-row-actions">
                {state === 'UnackActive' && (
                    <Button variant="ghost" size="sm" onClick={onAcknowledge}>
                        Acknowledge
                    </Button>
                )}
                {state !== 'Shelved' && (
                    <Button variant="ghost" size="sm" onClick={onShelve}>
                        Shelve
                    </Button>
                )}
                <Button variant="ghost" size="sm" icon={ChevronRight} />
            </div>
        </div>
    );
};

// Progress Ring Component
interface ProgressRingProps {
    value: number;
    max?: number;
    size?: number;
    strokeWidth?: number;
    color?: 'primary' | 'success' | 'warning' | 'danger';
    showValue?: boolean;
    unit?: string;
}

export const ProgressRing: React.FC<ProgressRingProps> = ({
    value,
    max = 100,
    size = 100,
    strokeWidth = 8,
    color = 'primary',
    showValue = true,
    unit,
}) => {
    const radius = (size - strokeWidth) / 2;
    const circumference = radius * 2 * Math.PI;
    const percentage = Math.min(value / max, 1);
    const offset = circumference - percentage * circumference;

    return (
        <div className="progress-ring-wrapper" style={{ width: size, height: size }}>
            <svg className="progress-ring" width={size} height={size}>
                <circle
                    className="progress-ring-bg"
                    strokeWidth={strokeWidth}
                    stroke="currentColor"
                    fill="transparent"
                    r={radius}
                    cx={size / 2}
                    cy={size / 2}
                />
                <circle
                    className={`progress-ring-progress color-${color}`}
                    strokeWidth={strokeWidth}
                    strokeDasharray={circumference}
                    strokeDashoffset={offset}
                    strokeLinecap="round"
                    stroke="currentColor"
                    fill="transparent"
                    r={radius}
                    cx={size / 2}
                    cy={size / 2}
                />
            </svg>
            {showValue && (
                <div className="progress-ring-value">
                    <span className="progress-ring-number">{Math.round(value)}</span>
                    {unit && <span className="progress-ring-unit">{unit}</span>}
                </div>
            )}
        </div>
    );
};

// System Health Widget
interface SystemHealthProps {
    cpuUsage: number;
    memoryUsage: number;
    diskUsage?: number;
    connected: boolean;
}

export const SystemHealth: React.FC<SystemHealthProps> = ({
    cpuUsage,
    memoryUsage,
    diskUsage = 45,
    connected,
}) => {
    const getColor = (value: number): 'success' | 'warning' | 'danger' => {
        if (value < 60) return 'success';
        if (value < 80) return 'warning';
        return 'danger';
    };

    return (
        <Card className="system-health-card">
            <div className="system-health-header">
                <h3>System Health</h3>
                <div className={`connection-indicator ${connected ? 'connected' : 'disconnected'}`}>
                    <Activity className="connection-icon" />
                    <span>{connected ? 'Connected' : 'Disconnected'}</span>
                </div>
            </div>
            <div className="system-health-metrics">
                <div className="health-metric">
                    <ProgressRing value={cpuUsage} color={getColor(cpuUsage)} size={80} />
                    <div className="health-metric-label">
                        <Cpu className="metric-icon" />
                        <span>CPU</span>
                    </div>
                </div>
                <div className="health-metric">
                    <ProgressRing value={memoryUsage} color={getColor(memoryUsage)} size={80} />
                    <div className="health-metric-label">
                        <Database className="metric-icon" />
                        <span>Memory</span>
                    </div>
                </div>
                <div className="health-metric">
                    <ProgressRing value={diskUsage} color={getColor(diskUsage)} size={80} />
                    <div className="health-metric-label">
                        <HardDrive className="metric-icon" />
                        <span>Disk</span>
                    </div>
                </div>
            </div>
        </Card>
    );
};

// Empty State Component
interface EmptyStateProps {
    icon?: LucideIcon;
    title: string;
    description?: string;
    action?: {
        label: string;
        onClick: () => void;
    };
}

export const EmptyState: React.FC<EmptyStateProps> = ({
    icon: Icon = Database,
    title,
    description,
    action,
}) => {
    return (
        <div className="empty-state">
            <Icon className="empty-state-icon" />
            <h4 className="empty-state-title">{title}</h4>
            {description && <p className="empty-state-description">{description}</p>}
            {action && (
                <Button variant="primary" onClick={action.onClick}>
                    {action.label}
                </Button>
            )}
        </div>
    );
};

// Loading Skeleton
interface SkeletonProps {
    width?: string | number;
    height?: string | number;
    className?: string;
}

export const Skeleton: React.FC<SkeletonProps> = ({
    width = '100%',
    height = '1rem',
    className = '',
}) => {
    return (
        <div
            className={`skeleton ${className}`}
            style={{ width, height }}
        />
    );
};
