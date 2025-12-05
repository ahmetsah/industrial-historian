import React, { useEffect, useRef, useLayoutEffect } from 'react';
import uPlot from 'uplot';
import 'uplot/dist/uPlot.min.css';

interface UPlotWrapperProps {
    options: uPlot.Options;
    data: uPlot.AlignedData;
    onCreate?: (chart: uPlot) => void;
    onDelete?: (chart: uPlot) => void;
    className?: string;
}

export const UPlotWrapper: React.FC<UPlotWrapperProps> = ({
    options,
    data,
    onCreate,
    onDelete,
    className,
}) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const chartRef = useRef<uPlot | null>(null);

    // Initialize chart
    useLayoutEffect(() => {
        if (!containerRef.current) return;

        // Destroy existing chart if any (shouldn't happen with correct deps, but safety first)
        if (chartRef.current) {
            if (onDelete) onDelete(chartRef.current);
            chartRef.current.destroy();
            chartRef.current = null;
        }

        const chart = new uPlot(options, data, containerRef.current);
        chartRef.current = chart;

        if (onCreate) {
            onCreate(chart);
        }

        return () => {
            if (chartRef.current) {
                if (onDelete) onDelete(chartRef.current);
                chartRef.current.destroy();
                chartRef.current = null;
            }
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [options, onCreate, onDelete]); // Re-create chart only if options change. Data updates should be handled manually via chart instance.

    // Handle resizing
    useEffect(() => {
        if (!containerRef.current || !chartRef.current) return;

        const resizeObserver = new ResizeObserver((entries) => {
            if (!chartRef.current) return;
            for (const entry of entries) {
                const { width, height } = entry.contentRect;
                chartRef.current.setSize({ width, height });
            }
        });

        resizeObserver.observe(containerRef.current);

        return () => {
            resizeObserver.disconnect();
        };
    }, []);

    return <div ref={containerRef} className={className} style={{ width: '100%', height: '100%' }} />;
};
