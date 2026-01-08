import React, { useEffect, useRef, useMemo } from 'react';
import uPlot from 'uplot';
import 'uplot/dist/uPlot.min.css';
import './trend-chart.css';

interface TrendChartProps {
    title?: string;
    data: {
        sensorId: string;
        label: string;
        color: string;
        unit: string;
        points: { timestamp: number; value: number }[];
    }[];
    height?: number;
    showLegend?: boolean;
    showGrid?: boolean;
    enableZoom?: boolean;
    enablePan?: boolean;
    timeRange?: { start: number; end: number };
    onTimeRangeChange?: (range: { start: number; end: number }) => void;
}

const TrendChart: React.FC<TrendChartProps> = ({
    title,
    data,
    height = 300,
    showLegend = true,
    showGrid = true,
    enableZoom = true,
    timeRange,
}) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const chartRef = useRef<uPlot | null>(null);

    // Transform data for uPlot format
    const uPlotData = useMemo(() => {
        if (data.length === 0 || data[0].points.length === 0) {
            return [[], []] as [number[], number[]];
        }

        // Use timestamps from first series as x-axis
        const timestamps = data[0].points.map(p => p.timestamp / 1000);
        const series = data.map(s => s.points.map(p => p.value));

        return [timestamps, ...series] as uPlot.AlignedData;
    }, [data]);

    // Chart options
    const options = useMemo((): uPlot.Options => {
        const isDarkMode = true;

        const seriesConfig: uPlot.Series[] = [
            {
                label: 'Time',
            },
            ...data.map((s) => ({
                label: s.label,
                stroke: s.color,
                width: 2,
                spanGaps: true,
                points: { show: false },
                fill: (self: uPlot) => {
                    const ctx = self.ctx;
                    const gradient = ctx.createLinearGradient(0, 0, 0, height);
                    gradient.addColorStop(0, `${s.color}33`);
                    gradient.addColorStop(1, `${s.color}00`);
                    return gradient;
                },
            })),
        ];

        return {
            width: 0, // Will be set on mount
            height,
            padding: [20, 20, 0, 0],
            cursor: {
                show: true,
                drag: enableZoom ? { x: true, y: false } : undefined,
                focus: { prox: 16 },
                points: {
                    show: true,
                    size: 8,
                    stroke: isDarkMode ? '#fff' : '#000',
                    width: 2,
                    fill: (self: uPlot, si: number) => {
                        const series = self.series[si];
                        return series?.stroke?.toString() || '#fff';
                    },
                },
            },
            legend: {
                show: showLegend,
            },
            axes: [
                {
                    // X axis
                    stroke: isDarkMode ? '#64748b' : '#94a3b8',
                    grid: {
                        show: showGrid,
                        stroke: isDarkMode ? '#334155' : '#e2e8f0',
                        width: 1,
                    },
                    ticks: {
                        show: true,
                        stroke: isDarkMode ? '#475569' : '#cbd5e1',
                    },
                    values: (_u: uPlot, ticks: number[]) =>
                        ticks.map((v) => {
                            const date = new Date(v * 1000);
                            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                        }),
                },
                {
                    // Y axis
                    stroke: isDarkMode ? '#64748b' : '#94a3b8',
                    grid: {
                        show: showGrid,
                        stroke: isDarkMode ? '#334155' : '#e2e8f0',
                        width: 1,
                    },
                    ticks: {
                        show: true,
                        stroke: isDarkMode ? '#475569' : '#cbd5e1',
                    },
                    values: (_u: uPlot, ticks: number[]) =>
                        ticks.map((v) => v.toFixed(1)),
                    size: 60,
                },
            ],
            series: seriesConfig,
            scales: {
                x: {
                    time: true,
                    range: timeRange ? [timeRange.start / 1000, timeRange.end / 1000] : undefined,
                },
                y: {
                    auto: true,
                },
            },
        };
    }, [data, height, showLegend, showGrid, enableZoom, timeRange]);

    // Create/update chart
    useEffect(() => {
        if (!containerRef.current) return;

        // Destroy previous chart
        if (chartRef.current) {
            chartRef.current.destroy();
            chartRef.current = null;
        }

        if (uPlotData[0].length === 0) return;

        const container = containerRef.current;
        const width = container.clientWidth;

        const opts = { ...options, width };
        chartRef.current = new uPlot(opts, uPlotData, container);

        // Cleanup
        return () => {
            if (chartRef.current) {
                chartRef.current.destroy();
                chartRef.current = null;
            }
        };
    }, [uPlotData, options]);

    // Handle resize
    useEffect(() => {
        if (!containerRef.current || !chartRef.current) return;

        const handleResize = () => {
            if (containerRef.current && chartRef.current) {
                chartRef.current.setSize({
                    width: containerRef.current.clientWidth,
                    height,
                });
            }
        };

        const resizeObserver = new ResizeObserver(handleResize);
        resizeObserver.observe(containerRef.current);

        return () => resizeObserver.disconnect();
    }, [height]);

    return (
        <div className="trend-chart-container">
            {title && (
                <div className="trend-chart-header">
                    <h3 className="trend-chart-title">{title}</h3>
                    <div className="trend-chart-legend">
                        {data.map((s) => (
                            <div key={s.sensorId} className="legend-item">
                                <span className="legend-color" style={{ background: s.color }} />
                                <span className="legend-label">{s.label}</span>
                                <span className="legend-unit">{s.unit}</span>
                            </div>
                        ))}
                    </div>
                </div>
            )}
            <div
                ref={containerRef}
                className="trend-chart-canvas"
                style={{ height }}
            />
            {uPlotData[0].length === 0 && (
                <div className="trend-chart-empty">
                    <p>No data available for the selected time range</p>
                </div>
            )}
        </div>
    );
};

export default TrendChart;
