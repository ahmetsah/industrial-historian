import React from 'react';
import Layout from './components/layout/Layout';
import Dashboard from './pages/Dashboard';
import TrendsPage from './pages/Trends';
import AlarmsPage from './pages/Alarms';
import AuditPage from './pages/Audit';
import Settings from './pages/Settings';
import { useUIStore } from './stores';
import './App.css';

const App: React.FC = () => {
  const { activeTab, setActiveTab, sidebarOpen, toggleSidebar } = useUIStore();

  const renderPage = () => {
    switch (activeTab) {
      case 'dashboard':
        return <Dashboard />;
      case 'trends':
        return <TrendsPage />;
      case 'alarms':
        return <AlarmsPage />;
      case 'audit':
        return <AuditPage />;
      case 'settings':
        return <Settings />;
      default:
        return <Dashboard />;
    }
  };

  return (
    <Layout
      activeTab={activeTab}
      onTabChange={setActiveTab}
      sidebarOpen={sidebarOpen}
      onToggleSidebar={toggleSidebar}
    >
      {renderPage()}
    </Layout>
  );
};

export default App;
