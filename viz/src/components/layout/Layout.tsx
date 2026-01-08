import React from 'react';
import {
    LayoutDashboard,
    Activity,
    AlertTriangle,
    FileText,
    Settings,
    ChevronLeft,
    ChevronRight,
    Zap,
    Menu,
} from 'lucide-react';
import './layout.css';

interface LayoutProps {
    children: React.ReactNode;
    activeTab: string;
    onTabChange: (tab: string) => void;
    sidebarOpen: boolean;
    onToggleSidebar: () => void;
}

const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { id: 'trends', label: 'Trends', icon: Activity },
    { id: 'alarms', label: 'Alarms', icon: AlertTriangle },
    { id: 'audit', label: 'Audit Trail', icon: FileText },
    { id: 'settings', label: 'Settings', icon: Settings },
];

const Layout: React.FC<LayoutProps> = ({
    children,
    activeTab,
    onTabChange,
    sidebarOpen,
    onToggleSidebar,
}) => {
    return (
        <div className="layout">
            {/* Sidebar */}
            <aside className={`sidebar ${sidebarOpen ? 'open' : 'collapsed'}`}>
                <div className="sidebar-header">
                    <div className="logo">
                        <div className="logo-icon">
                            <Zap className="logo-icon-svg" />
                        </div>
                        {sidebarOpen && (
                            <div className="logo-text">
                                <span className="logo-title">Historian</span>
                                <span className="logo-subtitle">Industrial Platform</span>
                            </div>
                        )}
                    </div>
                    <button className="sidebar-toggle" onClick={onToggleSidebar}>
                        {sidebarOpen ? <ChevronLeft /> : <ChevronRight />}
                    </button>
                </div>

                <nav className="sidebar-nav">
                    {navItems.map((item) => (
                        <button
                            key={item.id}
                            className={`nav-item ${activeTab === item.id ? 'active' : ''}`}
                            onClick={() => onTabChange(item.id)}
                            title={!sidebarOpen ? item.label : undefined}
                        >
                            <item.icon className="nav-icon" />
                            {sidebarOpen && <span className="nav-label">{item.label}</span>}
                        </button>
                    ))}
                </nav>

                <div className="sidebar-footer">
                    {sidebarOpen && (
                        <div className="system-status">
                            <div className="status-dot connected" />
                            <div className="status-text">
                                <span className="status-label">System Status</span>
                                <span className="status-value">All Systems Operational</span>
                            </div>
                        </div>
                    )}
                </div>
            </aside>

            {/* Main Content */}
            <div className="main-wrapper">
                {/* Top Bar */}
                <header className="topbar">
                    <button className="mobile-menu-btn" onClick={onToggleSidebar}>
                        <Menu />
                    </button>

                    <div className="topbar-title">
                        <h1>{navItems.find(n => n.id === activeTab)?.label || 'Dashboard'}</h1>
                    </div>

                    <div className="topbar-actions">
                        <div className="live-badge">
                            <span className="live-dot" />
                            <span>LIVE</span>
                        </div>

                        <div className="time-display">
                            <Clock />
                        </div>
                    </div>
                </header>

                {/* Page Content */}
                <main className="main-content">
                    {children}
                </main>
            </div>
        </div>
    );
};

// Clock component for topbar
const Clock: React.FC = () => {
    const [time, setTime] = React.useState(new Date());

    React.useEffect(() => {
        const interval = setInterval(() => setTime(new Date()), 1000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="clock">
            <span className="clock-time">
                {time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
            </span>
            <span className="clock-date">
                {time.toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' })}
            </span>
        </div>
    );
};

export default Layout;
