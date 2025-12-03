import React, { useState } from 'react';
import { X, Settings } from 'lucide-react';
import { type DashboardWidget as WidgetType } from '../types';
import { useDashboardStore } from '../stores/useDashboardStore';
import { twMerge } from 'tailwind-merge';
import { TrendChart } from './TrendChart';
import { WidgetConfigModal } from './WidgetConfigModal';

interface DashboardWidgetProps extends React.HTMLAttributes<HTMLDivElement> {
    widget: WidgetType;
    children?: React.ReactNode;
    className?: string;
    // Props injected by react-grid-layout
    style?: React.CSSProperties;
    onMouseDown?: React.MouseEventHandler;
    onMouseUp?: React.MouseEventHandler;
    onTouchEnd?: React.TouchEventHandler;
}

// forwardRef is required by react-grid-layout
export const DashboardWidget = React.forwardRef<HTMLDivElement, DashboardWidgetProps>(
    ({ widget, children, className, style, onMouseDown, onMouseUp, onTouchEnd, ...props }, ref) => {
        const { isEditMode, removeWidget } = useDashboardStore();
        const [isConfigOpen, setIsConfigOpen] = useState(false);

        return (
            <>
                <div
                    ref={ref}
                    style={style}
                    onMouseDown={onMouseDown}
                    onMouseUp={onMouseUp}
                    onTouchEnd={onTouchEnd}
                    className={twMerge(
                        'bg-white dark:bg-slate-800 rounded-lg shadow-sm border border-slate-200 dark:border-slate-700 flex flex-col overflow-hidden',
                        className
                    )}
                    {...props}
                >
                    <div className="flex-none flex items-center justify-between p-2 border-b border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-900/50 cursor-move draggable-handle">
                        <h3 className="text-sm font-medium text-slate-700 dark:text-slate-200 truncate select-none">
                            {widget.title}
                        </h3>
                        <div className="flex items-center gap-1">
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setIsConfigOpen(true);
                                }}
                                onMouseDown={(e) => e.stopPropagation()}
                                className="p-1 text-slate-400 hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors cursor-pointer"
                                title="Configure widget"
                            >
                                <Settings size={14} />
                            </button>
                            {isEditMode && (
                                <button
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        removeWidget(widget.id);
                                    }}
                                    onMouseDown={(e) => e.stopPropagation()}
                                    className="p-1 text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors cursor-pointer"
                                    title="Remove widget"
                                >
                                    <X size={14} />
                                </button>
                            )}
                        </div>
                    </div>
                    <div className="flex-1 min-h-0 overflow-hidden relative p-4">
                        {widget.type === 'chart' ? (
                            <TrendChart title={widget.title} sensorId={widget.config?.sensorId as string} />
                        ) : (
                            children
                        )}
                    </div>
                </div>
                <WidgetConfigModal
                    widget={widget}
                    isOpen={isConfigOpen}
                    onClose={() => setIsConfigOpen(false)}
                />
            </>
        );
    }
);

DashboardWidget.displayName = 'DashboardWidget';
