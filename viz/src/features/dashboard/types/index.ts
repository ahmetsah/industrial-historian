export interface WidgetLayout {
    x: number;
    y: number;
    w: number;
    h: number;
}

export interface DashboardWidget {
    id: string;
    type: string;
    title: string;
    layout: WidgetLayout;
    config: Record<string, unknown>;
}
