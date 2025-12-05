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
                        'bg-slate-900 border border-white/5 rounded-xl shadow-xl shadow-black/20 flex flex-col overflow-hidden backdrop-blur-sm transition-shadow hover:shadow-2xl hover:border-white/10 group',
                        className
                    )}
                    {...props}
                >
                    <div className="flex-none flex items-center justify-between px-4 py-3 border-b border-white/5 bg-slate-900/50 cursor-move draggable-handle group-hover:bg-slate-800/50 transition-colors">
                        <div className="flex items-center gap-2">
                            <div className="w-1.5 h-1.5 rounded-full bg-industrial-400 shadow-[0_0_8px_rgba(56,189,248,0.5)]"></div>
                            <h3 className="text-sm font-medium text-slate-200 tracking-wide truncate select-none">
                                {widget.title}
                            </h3>
                        </div>
                        <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setIsConfigOpen(true);
                                }}
                                onMouseDown={(e) => e.stopPropagation()}
                                className="p-1.5 text-slate-400 hover:text-industrial-400 hover:bg-white/5 rounded-lg transition-all cursor-pointer"
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
                                    className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-all cursor-pointer"
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
