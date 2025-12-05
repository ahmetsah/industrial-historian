import React, { useMemo } from 'react';
import { Responsive, WidthProvider, type Layout } from 'react-grid-layout';
import { useDashboardStore } from '../stores/useDashboardStore';
import { DashboardWidget } from './DashboardWidget';
import { WidgetPlaceholder } from './WidgetPlaceholder';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';

const ResponsiveGridLayout = WidthProvider(Responsive);

export const DashboardLayout: React.FC = () => {
    const { widgets, isEditMode, updateLayout } = useDashboardStore();

    const layouts = useMemo(() => {
        return {
            lg: widgets.map((w) => ({ i: w.id, ...w.layout })),
        };
    }, [widgets]);

    const handleLayoutChange = (currentLayout: Layout[]) => {
        const updates = currentLayout.map(item => ({
            i: item.i,
            x: item.x,
            y: item.y,
            w: item.w,
            h: item.h
        }));

        updateLayout(updates);
    };

    return (
        <div className="w-full min-h-screen p-6">
            <ResponsiveGridLayout
                className="layout"
                layouts={layouts}
                breakpoints={{ lg: 1200, md: 996, sm: 768, xs: 480, xxs: 0 }}
                cols={{ lg: 12, md: 10, sm: 6, xs: 4, xxs: 2 }}
                rowHeight={30}
                isDraggable={isEditMode}
                isResizable={isEditMode}
                draggableHandle=".draggable-handle"
                onLayoutChange={handleLayoutChange}
                margin={[16, 16]}
            >
                {widgets.map((widget) => (
                    <div key={widget.id}>
                        <DashboardWidget widget={widget} className="h-full">
                            <WidgetPlaceholder type={widget.type} />
                        </DashboardWidget>
                    </div>
                ))}
            </ResponsiveGridLayout>
        </div>
    );
};
