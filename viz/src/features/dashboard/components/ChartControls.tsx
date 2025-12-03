import React from 'react';
import { Pause, Play, RotateCcw, Download } from 'lucide-react';

interface ChartControlsProps {
    isPaused: boolean;
    onTogglePause: () => void;
    onReset: () => void;
    onExport: () => void;
}

export const ChartControls: React.FC<ChartControlsProps> = ({
    isPaused,
    onTogglePause,
    onReset,
    onExport
}) => {
    return (
        <div className="absolute top-2 right-2 z-10 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
                onClick={onTogglePause}
                className="p-1.5 bg-white/80 dark:bg-slate-800/80 backdrop-blur rounded shadow hover:bg-white dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200"
                title={isPaused ? "Resume" : "Pause"}
            >
                {isPaused ? <Play size={16} /> : <Pause size={16} />}
            </button>
            <button
                onClick={onReset}
                className="p-1.5 bg-white/80 dark:bg-slate-800/80 backdrop-blur rounded shadow hover:bg-white dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200"
                title="Reset View"
            >
                <RotateCcw size={16} />
            </button>
            <button
                onClick={onExport}
                className="p-1.5 bg-white/80 dark:bg-slate-800/80 backdrop-blur rounded shadow hover:bg-white dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200"
                title="Export CSV"
            >
                <Download size={16} />
            </button>
        </div>
    );
};
