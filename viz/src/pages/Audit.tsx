import React, { useEffect, useState, useCallback } from 'react';
import {
    Check,
    ChevronLeft,
    ChevronRight,
    FileText,
    Filter,
    Link2,
    RefreshCw,
    Search,
    Shield,
    User,
    AlertCircle,
} from 'lucide-react';
import { Button, Card, EmptyState } from '../components/ui';
import { useAuditStore } from '../stores';
import { auditAPI } from '../api';
import type { AuditLogEntry } from '../stores';
import './Audit.css';

const actionIcons: Record<string, React.ReactNode> = {
    'ALARM_ACK': <Check className="action-icon ack" />,
    'ALARM_TRIGGER': <Shield className="action-icon trigger" />,
    'SETPOINT_CHANGE': <FileText className="action-icon change" />,
    'LOGIN': <User className="action-icon login" />,
    'LOGOUT': <User className="action-icon logout" />,
};

const AuditPage: React.FC = () => {
    const { logs, setLogs } = useAuditStore();
    const [isLoading, setIsLoading] = useState(true);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [searchTerm, setSearchTerm] = useState('');
    const [filterAction, setFilterAction] = useState<string>('all');
    const [currentPage, setCurrentPage] = useState(1);
    const [selectedLog, setSelectedLog] = useState<AuditLogEntry | null>(null);
    const [integrityStatus, setIntegrityStatus] = useState<'verified' | 'pending' | 'failed'>('pending');
    const [apiError, setApiError] = useState<string | null>(null);
    const logsPerPage = 20;

    // Load logs from API
    const loadLogs = useCallback(async () => {
        try {
            setApiError(null);
            const fetchedLogs = await auditAPI.getLogs(100, 0);
            setLogs(fetchedLogs);
            setIntegrityStatus('verified');
        } catch (err) {
            console.error('Failed to load audit logs:', err);
            setApiError(err instanceof Error ? err.message : 'Failed to load audit logs');
            setIntegrityStatus('failed');
        }
    }, [setLogs]);

    // Initial load
    useEffect(() => {
        const init = async () => {
            setIsLoading(true);
            await loadLogs();
            setIsLoading(false);
        };
        init();
    }, [loadLogs]);

    // Unique actions for filter
    const uniqueActions = Array.from(new Set(logs.map(log => log.action)));

    // Filter logs
    const filteredLogs = logs.filter(log => {
        if (filterAction !== 'all' && log.action !== filterAction) return false;

        if (searchTerm) {
            const term = searchTerm.toLowerCase();
            return (
                log.actor.toLowerCase().includes(term) ||
                log.action.toLowerCase().includes(term) ||
                JSON.stringify(log.details).toLowerCase().includes(term)
            );
        }

        return true;
    });

    // Pagination
    const totalPages = Math.ceil(filteredLogs.length / logsPerPage);
    const paginatedLogs = filteredLogs.slice(
        (currentPage - 1) * logsPerPage,
        currentPage * logsPerPage
    );

    const handleRefresh = async () => {
        setIsRefreshing(true);
        await loadLogs();
        setIsRefreshing(false);
    };

    const handleVerifyIntegrity = async () => {
        setIntegrityStatus('pending');
        try {
            const result = await auditAPI.verifyIntegrity();
            setIntegrityStatus(result.valid ? 'verified' : 'failed');
        } catch (err) {
            console.error('Integrity verification failed:', err);
            setIntegrityStatus('failed');
        }
    };

    const formatTimestamp = (dateStr: string) => {
        const date = new Date(dateStr);
        return {
            date: date.toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' }),
            time: date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }),
        };
    };

    if (isLoading) {
        return (
            <div className="audit-page">
                <div className="audit-loading">
                    <FileText className="loading-icon" />
                    <p>Loading audit logs...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="audit-page">
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
            <div className="audit-header">
                <div className="audit-title-section">
                    <h2>Audit Trail</h2>
                    <div className={`integrity-badge ${integrityStatus}`}>
                        {integrityStatus === 'verified' && (
                            <>
                                <Check className="integrity-icon" />
                                <span>Chain Verified</span>
                            </>
                        )}
                        {integrityStatus === 'pending' && (
                            <>
                                <RefreshCw className="integrity-icon spin" />
                                <span>Verifying...</span>
                            </>
                        )}
                        {integrityStatus === 'failed' && (
                            <>
                                <Shield className="integrity-icon" />
                                <span>Integrity Issue</span>
                            </>
                        )}
                    </div>
                </div>

                <div className="audit-actions">
                    <Button
                        variant="ghost"
                        size="sm"
                        icon={RefreshCw}
                        onClick={handleRefresh}
                        loading={isRefreshing}
                    >
                        Refresh
                    </Button>
                    <Button
                        variant="secondary"
                        size="sm"
                        icon={Link2}
                        onClick={handleVerifyIntegrity}
                        disabled={integrityStatus === 'pending'}
                    >
                        Verify Integrity
                    </Button>
                </div>
            </div>

            {/* FDA Compliance Notice */}
            <div className="fda-notice">
                <div className="fda-badge">FDA 21 CFR Part 11</div>
                <span>All records are immutable with cryptographic chained-hash verification</span>
            </div>

            {/* Filters */}
            <div className="audit-filters">
                <div className="search-box">
                    <Search className="search-icon" />
                    <input
                        type="text"
                        placeholder="Search by actor, action, or details..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>

                <div className="action-filter">
                    <Filter className="filter-icon" />
                    <select
                        value={filterAction}
                        onChange={(e) => setFilterAction(e.target.value)}
                    >
                        <option value="all">All Actions</option>
                        {uniqueActions.map(action => (
                            <option key={action} value={action}>{action}</option>
                        ))}
                    </select>
                </div>

                <div className="log-count">
                    {filteredLogs.length} {filteredLogs.length === 1 ? 'record' : 'records'}
                </div>
            </div>

            {/* Logs Table */}
            <Card className="audit-table-card">
                {paginatedLogs.length === 0 ? (
                    <EmptyState
                        icon={FileText}
                        title="No audit logs found"
                        description="Try adjusting your search or filters"
                    />
                ) : (
                    <div className="audit-table">
                        <div className="audit-table-header">
                            <div className="th-time">Timestamp</div>
                            <div className="th-action">Action</div>
                            <div className="th-actor">Actor</div>
                            <div className="th-details">Details</div>
                            <div className="th-hash">Chain Hash</div>
                        </div>

                        <div className="audit-table-body">
                            {paginatedLogs.map((log, index) => {
                                const { date, time } = formatTimestamp(log.timestamp);

                                return (
                                    <div
                                        key={log.id}
                                        className={`audit-table-row ${selectedLog?.id === log.id ? 'selected' : ''}`}
                                        onClick={() => setSelectedLog(selectedLog?.id === log.id ? null : log)}
                                        style={{ animationDelay: `${index * 30}ms` }}
                                    >
                                        <div className="td-time">
                                            <span className="log-date">{date}</span>
                                            <span className="log-time">{time}</span>
                                        </div>

                                        <div className="td-action">
                                            {actionIcons[log.action] || <FileText className="action-icon" />}
                                            <span className="action-name">{log.action}</span>
                                        </div>

                                        <div className="td-actor">
                                            <User className="actor-icon" />
                                            <span className="actor-name">{log.actor}</span>
                                        </div>

                                        <div className="td-details">
                                            <span className="details-preview">
                                                {JSON.stringify(log.details).slice(0, 50)}...
                                            </span>
                                        </div>

                                        <div className="td-hash">
                                            <code className="hash-preview" title={log.curr_hash}>
                                                {log.curr_hash.slice(0, 12)}...
                                            </code>
                                            <Link2 className="hash-link-icon" />
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                )}

                {/* Pagination */}
                {totalPages > 1 && (
                    <div className="audit-pagination">
                        <Button
                            variant="ghost"
                            size="sm"
                            icon={ChevronLeft}
                            onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                            disabled={currentPage === 1}
                        />
                        <span className="page-info">
                            Page {currentPage} of {totalPages}
                        </span>
                        <Button
                            variant="ghost"
                            size="sm"
                            icon={ChevronRight}
                            iconPosition="right"
                            onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                            disabled={currentPage === totalPages}
                        />
                    </div>
                )}
            </Card>

            {/* Log Details Panel */}
            {selectedLog && (
                <Card className="log-details-panel">
                    <div className="details-header">
                        <h3>Log Details</h3>
                        <button className="close-details" onClick={() => setSelectedLog(null)}>×</button>
                    </div>

                    <div className="details-content">
                        <div className="detail-row">
                            <span className="detail-label">ID</span>
                            <code className="detail-value">{selectedLog.id}</code>
                        </div>
                        <div className="detail-row">
                            <span className="detail-label">Timestamp</span>
                            <span className="detail-value">{new Date(selectedLog.timestamp).toISOString()}</span>
                        </div>
                        <div className="detail-row">
                            <span className="detail-label">Actor</span>
                            <span className="detail-value">{selectedLog.actor}</span>
                        </div>
                        <div className="detail-row">
                            <span className="detail-label">Action</span>
                            <span className="detail-value">{selectedLog.action}</span>
                        </div>
                        <div className="detail-row">
                            <span className="detail-label">Previous Hash</span>
                            <code className="detail-value hash">{selectedLog.prev_hash}</code>
                        </div>
                        <div className="detail-row">
                            <span className="detail-label">Current Hash</span>
                            <code className="detail-value hash">{selectedLog.curr_hash}</code>
                        </div>
                        <div className="detail-row full-width">
                            <span className="detail-label">Details</span>
                            <pre className="detail-json">{JSON.stringify(selectedLog.details, null, 2)}</pre>
                        </div>
                    </div>
                </Card>
            )}
        </div>
    );
};

export default AuditPage;
