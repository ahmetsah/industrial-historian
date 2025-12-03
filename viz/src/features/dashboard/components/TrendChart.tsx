import React, { useMemo } from 'react';
import { UPlotWrapper } from './uPlotWrapper';
import { ChartControls } from './ChartControls';
import { useTrendData } from '../hooks/useTrendData';
import { createTrendChartOptions } from '../config/uPlotConfig';

interface TrendChartProps {
    title?: string;
    sensorId?: string;
}

export const TrendChart: React.FC<TrendChartProps> = ({ title, sensorId = "test-sensor" }) => {
    const { dataRef, registerChart, isPaused, togglePause, reset } = useTrendData();

    const options = useMemo(() => createTrendChartOptions(title), [title]);

    const handleExport = () => {
        const end = Date.now();
        const start = end - 3600 * 1000; // Last 1 hour
        const url = `http://localhost:8080/api/v1/export?sensor_id=${sensorId}&start_ts=${start}&end_ts=${end}`;
        window.open(url, '_blank');
    };

    return (
        <div className="flex flex-col h-full w-full overflow-hidden relative group">
            <ChartControls
                isPaused={isPaused}
                onTogglePause={togglePause}
                onReset={reset}
                onExport={handleExport}
            />
            <div className="flex-1 min-h-0">
                <UPlotWrapper
                    options={options}
                    data={dataRef.current}
                    onCreate={registerChart}
                />
            </div>
        </div>
    );
};
