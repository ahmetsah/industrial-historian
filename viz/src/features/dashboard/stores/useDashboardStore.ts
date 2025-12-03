import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { type DashboardWidget } from '../types';

interface DashboardState {
    widgets: DashboardWidget[];
    isEditMode: boolean;
}

interface DashboardActions {
    addWidget: (type?: string, title?: string, config?: Record<string, unknown>) => void;
    removeWidget: (id: string) => void;
    updateWidget: (id: string, updates: Partial<Omit<DashboardWidget, 'id'>>) => void;
    updateLayout: (layouts: { i: string; x: number; y: number; w: number; h: number }[]) => void;
    toggleEditMode: () => void;
}

export const useDashboardStore = create<DashboardState & DashboardActions>()(
    persist(
        (set) => ({
            widgets: [],
            isEditMode: false,
            addWidget: (type = 'chart', title = 'New Widget', config = {}) =>
                set((state) => {
                    // Generate UUID with fallback for older browsers
                    const id = typeof crypto !== 'undefined' && crypto.randomUUID
                        ? crypto.randomUUID()
                        : `widget-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

                    // Find a suitable position (basic implementation: add to bottom)
                    const y = state.widgets.reduce((acc, w) => Math.max(acc, w.layout.y + w.layout.h), 0);

                    return {
                        widgets: [
                            ...state.widgets,
                            {
                                id,
                                type,
                                title,
                                layout: { x: 0, y, w: 6, h: 4 }, // Default size
                                config,
                            },
                        ],
                    };
                }),
            removeWidget: (id) =>
                set((state) => ({
                    widgets: state.widgets.filter((w) => w.id !== id),
                })),
            updateWidget: (id, updates) =>
                set((state) => ({
                    widgets: state.widgets.map((w) =>
                        w.id === id ? { ...w, ...updates } : w
                    ),
                })),
            updateLayout: (layouts) =>
                set((state) => ({
                    widgets: state.widgets.map((w) => {
                        const layoutUpdate = layouts.find((l) => l.i === w.id);
                        if (layoutUpdate) {
                            return {
                                ...w,
                                layout: {
                                    x: layoutUpdate.x,
                                    y: layoutUpdate.y,
                                    w: layoutUpdate.w,
                                    h: layoutUpdate.h,
                                },
                            };
                        }
                        return w;
                    }),
                })),
            toggleEditMode: () => set((state) => ({ isEditMode: !state.isEditMode })),
        }),
        {
            name: 'dashboard-storage',
        }
    )
);
