import React from 'react';

export const WidgetPlaceholder: React.FC<{ type: string }> = ({ type }) => {
    return (
        <div className="h-full w-full flex flex-col items-center justify-center text-slate-400 p-4 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-lg">
            <span className="text-lg font-medium mb-2">Placeholder</span>
            <span className="text-xs font-mono bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded">
                Type: {type}
            </span>
        </div>
    );
};
