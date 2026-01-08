import React, { useEffect, useState, useCallback } from 'react';
import {
    AlertTriangle,
    Bell,
    BellOff,
    Check,
    Clock,
    Filter,
    RefreshCw,
    Search,
    AlertCircle,
} from 'lucide-react';
import { Button, Card, StatusBadge, EmptyState } from '../components/ui';
import { useAlarmStore } from '../stores';
import { alarmAPI } from '../api';
import type { ActiveAlarm } from '../stores';
import './Alarms.css';

type AlarmFilter = 'all' | 'unacked' | 'acked' | 'shelved' | 'critical' | 'warning';

const AlarmsPage: React.FC = () => {
    const { activeAlarms, setActiveAlarms, acknowledgeAlarm } = useAlarmStore();
    const [filter, setFilter] = useState<AlarmFilter>('all');
    const [searchTerm, setSearchTerm] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [apiError, setApiError] = useState<string | null>(null);

    // Load alarms from API
    const loadAlarms = useCallback(async () => {
        try {
            setApiError(null);
            const alarms = await alarmAPI.getActiveAlarms();
            setActiveAlarms(alarms);
        } catch (err) {
            console.error('Failed to load alarms:', err);
            setApiError(err instanceof Error ? err.message : 'Failed to load alarms');
        }
    }, [setActiveAlarms]);

    // Initial load
    useEffect(() => {
        const init = async () => {
            setIsLoading(true);
            await loadAlarms();
            setIsLoading(false);
        };
        init();
    }, [loadAlarms]);

    // Filter alarms
    const filteredAlarms = activeAlarms.filter(alarm => {
        // Apply status filter
        if (filter === 'unacked' && !alarm.state.includes('Unack')) return false;
        if (filter === 'acked' && !alarm.state.includes('Ack')) return false;
        if (filter === 'shelved' && alarm.state !== 'Shelved') return false;
        if (filter === 'critical' && alarm.priority !== 'Critical') return false;
        if (filter === 'warning' && alarm.priority !== 'Warning') return false;

        // Apply search filter
        if (searchTerm) {
            const term = searchTerm.toLowerCase();
            const tag = (alarm.tag || '').toLowerCase();
            return tag.includes(term);
        }

        return true;
    });

    // Alarm counts
    const counts = {
        all: activeAlarms.length,
        unacked: activeAlarms.filter(a => a.state.includes('Unack')).length,
        acked: activeAlarms.filter(a => a.state.includes('Ack')).length,
        shelved: activeAlarms.filter(a => a.state === 'Shelved').length,
        critical: activeAlarms.filter(a => a.priority === 'Critical').length,
        warning: activeAlarms.filter(a => a.priority === 'Warning').length,
    };

    const handleRefresh = async () => {
        setIsRefreshing(true);
        await loadAlarms();
        setIsRefreshing(false);
    };

    const handleAcknowledge = async (alarmId: number) => {
        try {
            await alarmAPI.acknowledgeAlarm(alarmId);
            acknowledgeAlarm(alarmId);
        } catch (err) {
            console.error('Failed to acknowledge alarm:', err);
            alert('Failed to acknowledge alarm');
        }
    };

    const handleShelve = async (alarmId: number) => {
        try {
            await alarmAPI.shelveAlarm(alarmId, 60 * 60); // Shelve for 60 minutes (3600 seconds)
            await loadAlarms(); // Reload to get updated state
        } catch (err) {
            console.error('Failed to shelve alarm:', err);
            alert('Failed to shelve alarm');
        }
    };

    const handleAcknowledgeAll = async () => {
        const unackedAlarms = activeAlarms.filter(a => a.state === 'UnackActive');
        for (const alarm of unackedAlarms) {
            try {
                await alarmAPI.acknowledgeAlarm(alarm.id);
                acknowledgeAlarm(alarm.id);
            } catch (err) {
                console.error(`Failed to acknowledge alarm ${alarm.id}:`, err);
            }
        }
    };

    const timeAgo = (dateStr: string) => {
        const diff = Date.now() - new Date(dateStr).getTime();
        const seconds = Math.floor(diff / 1000);
        if (seconds < 60) return `${seconds}s ago`;
        const minutes = Math.floor(seconds / 60);
        if (minutes < 60) return `${minutes}m ago`;
        const hours = Math.floor(minutes / 60);
        if (hours < 24) return `${hours}h ago`;
        return `${Math.floor(hours / 24)}d ago`;
    };

    const getStatusType = (alarm: ActiveAlarm): 'normal' | 'warning' | 'critical' | 'shelved' | 'acknowledged' => {
        if (alarm.state === 'Shelved') return 'shelved';
        if (alarm.state.includes('Ack')) return 'acknowledged';
        if (alarm.priority === 'Critical') return 'critical';
        return 'warning';
    };

    if (isLoading) {
        return (
            <div className="alarms-page">
                <div className="alarms-loading">
                    <AlertTriangle className="loading-icon" />
                    <p>Loading alarms...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="alarms-page">
            {/* API Error */}
            {apiError && (
                <div className="api-error-banner">
                    <AlertCircle className="error-icon" />
                    <span>{apiError}</span>
                    <Button variant="ghost" size="sm" icon={RefreshCw} onClick={handleRefresh}>
                        Retry
                    </Button>
                </div>
            )}

            {/* Header */}
            <div className="alarms-header">
                <div className="alarms-title-section">
                    <h2>Alarm Management</h2>
                    <div className="alarm-summary">
                        {counts.critical > 0 && (
                            <span className="summary-critical">
                                <AlertTriangle /> {counts.critical} Critical
                            </span>
                        )}
                        {counts.unacked > 0 && (
                            <span className="summary-unacked">
                                <Bell /> {counts.unacked} Unacknowledged
                            </span>
                        )}
                    </div>
                </div>

                <div className="alarms-actions">
                    <Button
                        variant="ghost"
                        size="sm"
                        icon={RefreshCw}
                        onClick={handleRefresh}
                        loading={isRefreshing}
                    >
                        Refresh
                    </Button>
                    {counts.unacked > 0 && (
                        <Button
                            variant="primary"
                            size="sm"
                            icon={Check}
                            onClick={handleAcknowledgeAll}
                        >
                            Acknowledge All
                        </Button>
                    )}
                </div>
            </div>

            {/* Filters */}
            <div className="alarms-filters">
                <div className="filter-tabs">
                    <button
                        className={`filter-tab ${filter === 'all' ? 'active' : ''}`}
                        onClick={() => setFilter('all')}
                    >
                        <Filter />
                        All <span className="filter-count">{counts.all}</span>
                    </button>
                    <button
                        className={`filter-tab ${filter === 'unacked' ? 'active' : ''}`}
                        onClick={() => setFilter('unacked')}
                    >
                        <Bell />
                        Unacked <span className="filter-count">{counts.unacked}</span>
                    </button>
                    <button
                        className={`filter-tab ${filter === 'acked' ? 'active' : ''}`}
                        onClick={() => setFilter('acked')}
                    >
                        <Check />
                        Acked <span className="filter-count">{counts.acked}</span>
                    </button>
                    <button
                        className={`filter-tab ${filter === 'shelved' ? 'active' : ''}`}
                        onClick={() => setFilter('shelved')}
                    >
                        <BellOff />
                        Shelved <span className="filter-count">{counts.shelved}</span>
                    </button>
                    <button
                        className={`filter-tab priority-critical ${filter === 'critical' ? 'active' : ''}`}
                        onClick={() => setFilter('critical')}
                    >
                        Critical <span className="filter-count">{counts.critical}</span>
                    </button>
                    <button
                        className={`filter-tab priority-warning ${filter === 'warning' ? 'active' : ''}`}
                        onClick={() => setFilter('warning')}
                    >
                        Warning <span className="filter-count">{counts.warning}</span>
                    </button>
                </div>

                <div className="search-box">
                    <Search className="search-icon" />
                    <input
                        type="text"
                        placeholder="Search by tag..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
            </div>

            {/* Alarm List */}
            <Card className="alarms-list-card">
                {filteredAlarms.length === 0 ? (
                    <EmptyState
                        icon={Bell}
                        title="No alarms found"
                        description={filter !== 'all' ? "Try adjusting your filters" : "All systems are operating normally"}
                    />
                ) : (
                    <div className="alarms-table">
                        <div className="alarms-table-header">
                            <div className="th-status">Status</div>
                            <div className="th-priority">Priority</div>
                            <div className="th-tag">Tag</div>
                            <div className="th-value">Value</div>
                            <div className="th-time">Time</div>
                            <div className="th-actions">Actions</div>
                        </div>

                        <div className="alarms-table-body">
                            {filteredAlarms.map((alarm, index) => (
                                <div
                                    key={alarm.id}
                                    className={`alarm-table-row priority-${alarm.priority?.toLowerCase() || 'warning'} ${alarm.state === 'UnackActive' && alarm.priority === 'Critical' ? 'pulse' : ''
                                        }`}
                                    style={{ animationDelay: `${index * 50}ms` }}
                                >
                                    <div className="td-status">
                                        <StatusBadge
                                            status={getStatusType(alarm)}
                                            pulse={alarm.state === 'UnackActive' && alarm.priority === 'Critical'}
                                        />
                                    </div>

                                    <div className="td-priority">
                                        <span className={`priority-badge priority-${alarm.priority?.toLowerCase()}`}>
                                            {alarm.priority}
                                        </span>
                                    </div>

                                    <div className="td-tag">
                                        <span className="tag-name">{alarm.tag || `Alarm-${alarm.definition_id}`}</span>
                                        <span className="tag-type">{alarm.type} Limit</span>
                                    </div>

                                    <div className="td-value">
                                        <span className="value-number">{alarm.value.toFixed(2)}</span>
                                    </div>

                                    <div className="td-time">
                                        <Clock className="time-icon" />
                                        <div className="time-info">
                                            <span className="time-ago">{timeAgo(alarm.activation_time)}</span>
                                            <span className="time-date">
                                                {new Date(alarm.activation_time).toLocaleString([], {
                                                    month: 'short',
                                                    day: 'numeric',
                                                    hour: '2-digit',
                                                    minute: '2-digit'
                                                })}
                                            </span>
                                        </div>
                                    </div>

                                    <div className="td-actions">
                                        {alarm.state === 'UnackActive' && (
                                            <Button
                                                variant="success"
                                                size="sm"
                                                onClick={() => handleAcknowledge(alarm.id)}
                                            >
                                                Acknowledge
                                            </Button>
                                        )}
                                        {alarm.state !== 'Shelved' && (
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                onClick={() => handleShelve(alarm.id)}
                                            >
                                                Shelve
                                            </Button>
                                        )}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </Card>

            {/* ISA 18.2 Info */}
            <div className="isa-info">
                <div className="isa-badge">ISA 18.2</div>
                <span>Alarm state machine compliant (Normal → UnackActive → AckActive → Normal)</span>
            </div>
        </div>
    );
};

export default AlarmsPage;
