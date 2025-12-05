import React from 'react';
import { Plus, Edit2, Check } from 'lucide-react';
import { DashboardLayout } from '../features/dashboard/components/DashboardLayout';
import { useDashboardStore } from '../features/dashboard/stores/useDashboardStore';

export const DashboardPage: React.FC = () => {
    const { isEditMode, toggleEditMode, addWidget } = useDashboardStore();

    return (
        <div className="flex flex-col h-screen bg-slate-950 text-slate-50 font-sans">
            <header className="flex items-center justify-between px-6 py-4 bg-slate-900/50 backdrop-blur-md border-b border-white/5 sticky top-0 z-50">
                <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-industrial-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-industrial-500/20">
                        <div className="w-3 h-3 bg-white rounded-full"></div>
                    </div>
                    <h1 className="text-xl font-bold tracking-tight text-white">
                        Semper<span className="text-industrial-400">Historian</span>
                    </h1>
                </div>
                <div className="flex items-center gap-3">
                    <button
                        onClick={() => addWidget('chart', 'New Metric')}
                        className="flex items-center gap-2 px-4 py-2 bg-industrial-600 hover:bg-industrial-500 text-white rounded-lg transition-all shadow-lg shadow-industrial-900/40 border border-white/10 hover:border-white/20 text-sm font-medium active:scale-95"
                    >
                        <Plus size={16} />
                        Add Widget
                    </button>
                    <button
                        onClick={toggleEditMode}
                        className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-all text-sm font-medium border active:scale-95 ${isEditMode
                                ? 'bg-green-500/10 text-green-400 border-green-500/50 hover:bg-green-500/20'
                                : 'bg-slate-800 text-slate-300 border-white/5 hover:bg-slate-700 hover:text-white'
                            }`}
                    >
                        {isEditMode ? <Check size={16} /> : <Edit2 size={16} />}
                        {isEditMode ? 'Done' : 'Edit Layout'}
                    </button>
                </div>
            </header>
            <main className="flex-1 overflow-auto bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-slate-900 via-slate-950 to-black">
                <DashboardLayout />
            </main>
        </div>
    );
};
