import uPlot from 'uplot';

export const createTrendChartOptions = (title?: string): uPlot.Options => ({
    title: title || "Real-time Trend",
    width: 800, // Initial, will be resized
    height: 400,
    series: [
        {}, // X-Axis (Time)
        {
            stroke: "red",
            width: 2,
            spanGaps: false,
        }
    ],
    scales: {
        x: {
            time: true,
            auto: false, // We control the range for streaming
        },
        y: {
            auto: true,
        }
    },
    axes: [
        {
            values: [
                // tick incr          default           year                             month    day                        hour     min                sec       mode
                [3600 * 24 * 365, "{YYYY}", null, null, null, null, null, null, 1],
                [3600 * 24 * 28, "{MMM}", "\n{YYYY}", null, null, null, null, null, 1],
                [3600 * 24, "{M}/{D}", "\n{YYYY}", null, null, null, null, null, 1],
                [3600, "{h}{aa}", "\n{M}/{D}/{YY}", null, "\n{M}/{D}", null, null, null, 1],
                [60, "{h}:{mm}{aa}", "\n{M}/{D}/{YY}", null, "\n{M}/{D}", null, null, null, 1],
                [1, "{h}:{mm}:{ss}", "\n{M}/{D}/{YY}", null, "\n{M}/{D}", null, null, null, 1],
                [0.001, ":{ss}.{fff}", "\n{M}/{D}/{YY}", null, "\n{M}/{D}", null, null, null, 1],
            ],
        },
        {
            values: (_self, ticks) => ticks.map(n => n.toFixed(1)),
        }
    ],
    cursor: {
        drag: { x: true, y: true, uni: 50 },
    }
});
