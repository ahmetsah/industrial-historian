import { useEffect, useRef, useState, useCallback } from 'react';
import uPlot from 'uplot';

interface UseTrendDataConfig {
    maxPoints?: number;
    updateRateMs?: number;
    initialPoints?: number;
}

export const useTrendData = ({
    maxPoints = 100000,
    updateRateMs = 16,
    initialPoints = 1000,
}: UseTrendDataConfig = {}) => {
    const chartRef = useRef<uPlot | null>(null);
    const dataRef = useRef<uPlot.AlignedData>([[], []]); // [Time, Value]
    const [isPaused, setIsPaused] = useState(false);
    const requestRef = useRef<number | undefined>(undefined);
    const lastUpdateRef = useRef<number>(0);

    // Internal state for the mock generator
    const counterRef = useRef<number>(0);

    // Initialize Data
    useEffect(() => {
        let startTime = Date.now() / 1000;
        const times: number[] = [];
        const values: number[] = [];

        for (let i = 0; i < initialPoints; i++) {
            times.push(startTime - (initialPoints - i) * 0.1);
            values.push(Math.sin((counterRef.current + i) * 0.1) * 10 + 50);
        }
        counterRef.current += initialPoints;
        dataRef.current = [times, values];
    }, [initialPoints]);

    // Animation Loop
    useEffect(() => {
        const update = (timestamp: number) => {
            if (isPaused) {
                requestRef.current = requestAnimationFrame(update);
                return;
            }

            if (timestamp - lastUpdateRef.current >= updateRateMs) {
                const now = Date.now() / 1000;
                const newVal = Math.sin(counterRef.current * 0.1) * 10 + 50;
                counterRef.current++;

                const [ts, vals] = dataRef.current as [number[], number[]];

                ts.push(now);
                vals.push(newVal);

                if (ts.length > maxPoints) {
                    ts.shift();
                    vals.shift();
                }

                if (chartRef.current) {
                    const chart = chartRef.current;
                    chart.setData(dataRef.current);

                    const windowSize = 60;
                    const end = ts[ts.length - 1];
                    const start = end - windowSize;
                    chart.setScale('x', { min: start, max: end });
                }

                lastUpdateRef.current = timestamp;
            }

            requestRef.current = requestAnimationFrame(update);
        };

        requestRef.current = requestAnimationFrame(update);

        return () => {
            if (requestRef.current) cancelAnimationFrame(requestRef.current);
        };
    }, [isPaused, maxPoints, updateRateMs]);

    const registerChart = useCallback((chart: uPlot) => {
        chartRef.current = chart;
    }, []);

    const togglePause = useCallback(() => {
        setIsPaused((prev) => !prev);
    }, []);

    const reset = useCallback(() => {
        // Reset logic if needed
        // For streaming, maybe clear data?
        // For now, just a placeholder as in original code
    }, []);

    return {
        dataRef,
        registerChart,
        isPaused,
        togglePause,
        reset
    };
};
