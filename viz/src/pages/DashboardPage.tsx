import React from 'react';
import { Plus, Edit2, Check } from 'lucide-react';
import { DashboardLayout } from '../features/dashboard/components/DashboardLayout';
import { useDashboardStore } from '../features/dashboard/stores/useDashboardStore';

export const DashboardPage: React.FC = () => {
    const { isEditMode, toggleEditMode, addWidget } = useDashboardStore();

    return (
        <div className="flex flex-col h-screen">
            <header className="bg-white dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 p-4 flex items-center justify-between shadow-sm z-10 relative">
                <h1 className="text-xl font-semibold text-slate-800 dark:text-slate-100">
                    Industrial Historian Dashboard
                </h1>
                <div className="flex items-center gap-2">
                    <button
                        onClick={() => addWidget('chart', 'New Chart')}
                        className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium"
                    >
                        <Plus size={16} />
                        Add Widget
                    </button>
                    <button
                        onClick={toggleEditMode}
                        className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-colors text-sm font-medium border ${isEditMode
                                ? 'bg-green-50 text-green-700 border-green-200 hover:bg-green-100'
                                : 'bg-white text-slate-700 border-slate-300 hover:bg-slate-50'
                            }`}
                    >
                        {isEditMode ? <Check size={16} /> : <Edit2 size={16} />}
                        {isEditMode ? 'Done' : 'Edit Layout'}
                    </button>
                </div>
            </header>
            <main className="flex-1 overflow-auto bg-slate-50 dark:bg-slate-900">
                <DashboardLayout />
            </main>
        </div>
    );
};
