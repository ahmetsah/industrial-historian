import React, { useState } from 'react';
import { X } from 'lucide-react';
import { type DashboardWidget } from '../types';
import { useDashboardStore } from '../stores/useDashboardStore';

interface WidgetConfigModalProps {
    widget: DashboardWidget;
    isOpen: boolean;
    onClose: () => void;
}

export const WidgetConfigModal: React.FC<WidgetConfigModalProps> = ({ widget, isOpen, onClose }) => {
    const { updateWidget } = useDashboardStore();
    const [sensorId, setSensorId] = useState<string>((widget.config?.sensorId as string) || '');
    const [title, setTitle] = useState<string>(widget.title);

    if (!isOpen) return null;

    const handleSave = () => {
        updateWidget(widget.id, {
            title,
            config: {
                ...widget.config,
                sensorId,
            },
        });

        onClose();
    };

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={onClose}>
            <div
                className="bg-white dark:bg-slate-800 rounded-lg shadow-xl w-full max-w-md p-6"
                onClick={(e) => e.stopPropagation()}
            >
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-lg font-semibold text-slate-800 dark:text-slate-100">
                        Configure Widget
                    </h2>
                    <button
                        onClick={onClose}
                        className="p-1 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 rounded transition-colors"
                    >
                        <X size={20} />
                    </button>
                </div>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                            Widget Title
                        </label>
                        <input
                            type="text"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            className="w-full px-3 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            placeholder="Enter widget title"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                            NATS Subject / Sensor ID
                        </label>
                        <input
                            type="text"
                            value={sensorId}
                            onChange={(e) => setSensorId(e.target.value)}
                            className="w-full px-3 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            placeholder="e.g., acme.site.line1.sensor1"
                        />
                        <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">
                            Enter the NATS subject or sensor ID to monitor
                        </p>
                    </div>
                </div>

                <div className="flex justify-end gap-2 mt-6">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handleSave}
                        className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                    >
                        Save
                    </button>
                </div>
            </div>
        </div>
    );
};
