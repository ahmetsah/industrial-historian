# Modern Vue Web Interface - Development Plan

**Status:** Planning
**Epic:** Epic 7: Modern Web Interface with Vue
**Created:** 2025-12-11
**PM:** Ahmet

---

## 1. Executive Summary

This document outlines the complete plan for developing a modern, high-performance web interface for the Historian platform using **Vue 3** with **TypeScript**. The new interface will replace the existing React-based `/viz` directory and provide a premium, state-of-the-art user experience.

### Key Objectives
- ✅ Modern, premium UI/UX design
- ✅ High-performance real-time data visualization
- ✅ Responsive and mobile-friendly
- ✅ Integration with existing backend services (Engine, Auth, Alarm, Audit)
- ✅ FDA 21 CFR Part 11 compliance support
- ✅ ISA 18.2 Alarm Management interface

---

## 2. Current System Analysis

### 2.1 Existing Infrastructure

**Backend Services (Rust):**
- **Engine Service** (Port 8081)
  - HTTP Export API: `GET /api/v1/export`
  - Metadata API: `GET /api/v1/metadata`
  - gRPC Query API (Port 50051)
  - RocksDB + MinIO tiered storage
  - NATS JetStream integration

**Backend Services (Go):**
- **Auth Service** (Port 8080)
  - JWT-based authentication
  - FDA re-authentication support
  - PostgreSQL user database
  
- **Alarm Service** (Port 8083)
  - ISA 18.2 compliant alarm management
  - NATS event publishing
  
- **Audit Service** (Port 8082)
  - FDA 21 CFR Part 11 audit trail
  - PostgreSQL audit log storage

**Infrastructure:**
- NATS JetStream (Port 4222, 8222)
- MinIO S3 (Port 9000, 9001)
- PostgreSQL (Port 5432)
- Docker Compose orchestration

### 2.2 Recent Changes & Learnings

From the previous React implementation (now deprecated):
- ✅ Metadata API successfully provides sensor hierarchy
- ✅ Real-time streaming via NATS WebSocket works
- ✅ uPlot provides excellent performance for time-series charts
- ✅ Zustand state management was lightweight and effective
- ✅ TailwindCSS v4 with PostCSS integration
- ⚠️ Error boundary needed for production stability
- ⚠️ TypeScript strict mode revealed type safety issues

### 2.3 Data Flow Architecture

```
┌─────────────┐     Modbus/OPC UA      ┌──────────────┐
│   Sensors   │ ──────────────────────▶│  Ingestor    │
└─────────────┘                        └──────┬───────┘
                                              │
                                              │ NATS
                                              ▼
┌─────────────┐                        ┌──────────────┐
│   Vue UI    │◀───── HTTP/WS ─────────│    Engine    │
│             │                        │  (RocksDB +  │
│  - Charts   │                        │   MinIO)     │
│  - Alarms   │                        └──────────────┘
│  - Audit    │                              │
└─────────────┘                              │ NATS Events
                                              ▼
                                       ┌──────────────┐
                                       │ Alarm/Audit  │
                                       │   Services   │
                                       └──────────────┘
```

---

## 3. Technology Stack

### 3.1 Core Framework
- **Vue 3.4+** - Composition API with `<script setup>`
- **TypeScript 5.9+** - Strict mode enabled
- **Vite 7+** - Ultra-fast build tool and dev server

### 3.2 State Management
- **Pinia** - Official Vue state management (replaces Vuex)
  - Lightweight, TypeScript-first
  - DevTools integration
  - Modular stores for different features

### 3.3 UI Framework & Styling
- **TailwindCSS 4.x** - Utility-first CSS
- **Headless UI (Vue)** - Unstyled, accessible components
- **Heroicons** - Beautiful hand-crafted SVG icons
- **VueUse** - Collection of essential Vue composition utilities

### 3.4 Data Visualization
- **uPlot** - High-performance time-series charts (Canvas-based)
- **D3.js** (optional) - For complex custom visualizations
- **Chart.js** (optional) - For simple charts (pie, bar, etc.)

### 3.5 Real-time Communication
- **nats.ws** - NATS WebSocket client
- **Axios** - HTTP client for REST APIs
- **@grpc/grpc-js** + **@grpc/proto-loader** - gRPC client (if needed)

### 3.6 Testing & Quality
- **Vitest** - Vite-native unit testing
- **@vue/test-utils** - Vue component testing
- **Playwright** - E2E testing
- **ESLint** + **Prettier** - Code quality and formatting

### 3.7 Build & Deployment
- **Docker** - Multi-stage build with Nginx
- **Nginx** - Static file serving + API proxy
- **Docker Compose** - Local development orchestration

---

## 4. Architecture Design

### 4.1 Project Structure

```
web-ui/                              # New Vue application
├── .vscode/                         # VSCode settings
├── public/                          # Static assets
│   ├── favicon.ico
│   └── logo.svg
├── src/
│   ├── main.ts                      # Application entry point
│   ├── App.vue                      # Root component
│   ├── router/                      # Vue Router configuration
│   │   └── index.ts
│   ├── stores/                      # Pinia stores
│   │   ├── auth.ts                  # Authentication state
│   │   ├── sensors.ts               # Sensor metadata & filtering
│   │   ├── realtime.ts              # Real-time data stream
│   │   ├── alarms.ts                # Alarm state
│   │   └── audit.ts                 # Audit trail state
│   ├── composables/                 # Reusable composition functions
│   │   ├── useNatsStream.ts         # NATS WebSocket connection
│   │   ├── useHistoryQuery.ts       # Historical data fetching
│   │   ├── useAuth.ts               # Authentication logic
│   │   └── useTheme.ts              # Dark/Light mode
│   ├── services/                    # API clients
│   │   ├── api.ts                   # Axios instance configuration
│   │   ├── engine.ts                # Engine API client
│   │   ├── auth.ts                  # Auth API client
│   │   ├── alarm.ts                 # Alarm API client
│   │   └── audit.ts                 # Audit API client
│   ├── components/                  # Reusable components
│   │   ├── ui/                      # Base UI components
│   │   │   ├── Button.vue
│   │   │   ├── Input.vue
│   │   │   ├── Modal.vue
│   │   │   ├── Card.vue
│   │   │   └── Badge.vue
│   │   ├── charts/                  # Chart components
│   │   │   ├── TrendChart.vue       # Main time-series chart
│   │   │   ├── Sparkline.vue        # Mini inline chart
│   │   │   └── UPlotWrapper.vue     # uPlot base wrapper
│   │   ├── layout/                  # Layout components
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppSidebar.vue
│   │   │   └── AppFooter.vue
│   │   └── common/                  # Common components
│   │       ├── ErrorBoundary.vue
│   │       ├── LoadingSpinner.vue
│   │       └── EmptyState.vue
│   ├── views/                       # Page components
│   │   ├── DashboardView.vue        # Main dashboard
│   │   ├── SensorsView.vue          # Sensor list & detail
│   │   ├── AlarmsView.vue           # Alarm management
│   │   ├── AuditView.vue            # Audit trail viewer
│   │   ├── LoginView.vue            # Login page
│   │   └── NotFoundView.vue         # 404 page
│   ├── types/                       # TypeScript type definitions
│   │   ├── sensor.ts
│   │   ├── alarm.ts
│   │   ├── audit.ts
│   │   └── api.ts
│   ├── utils/                       # Utility functions
│   │   ├── format.ts                # Date/number formatting
│   │   ├── validation.ts            # Input validation
│   │   └── constants.ts             # App constants
│   └── assets/                      # Images, fonts, etc.
│       └── styles/
│           └── main.css             # Global styles + Tailwind
├── tests/
│   ├── unit/                        # Unit tests
│   └── e2e/                         # E2E tests
├── nginx.conf                       # Nginx configuration
├── Dockerfile                       # Multi-stage Docker build
├── .env.example                     # Environment variables template
├── .eslintrc.cjs                    # ESLint configuration
├── .prettierrc                      # Prettier configuration
├── tailwind.config.js               # Tailwind configuration
├── tsconfig.json                    # TypeScript configuration
├── vite.config.ts                   # Vite configuration
├── vitest.config.ts                 # Vitest configuration
└── package.json                     # Dependencies
```

### 4.2 Routing Structure

```typescript
// src/router/index.ts
const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    component: LoginView,
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/sensors',
    component: SensorsView,
    meta: { requiresAuth: true }
  },
  {
    path: '/sensors/:id',
    component: SensorDetailView,
    meta: { requiresAuth: true }
  },
  {
    path: '/alarms',
    component: AlarmsView,
    meta: { requiresAuth: true, requiredRole: 'operator' }
  },
  {
    path: '/audit',
    component: AuditView,
    meta: { requiresAuth: true, requiredRole: 'admin' }
  },
  {
    path: '/:pathMatch(.*)*',
    component: NotFoundView
  }
]
```

### 4.3 State Management Architecture

**Pinia Stores:**

1. **Auth Store** (`stores/auth.ts`)
   - User authentication state
   - JWT token management
   - Role-based permissions
   - Re-authentication for FDA compliance

2. **Sensors Store** (`stores/sensors.ts`)
   - Sensor metadata cache
   - Hierarchical filtering (Factory > Line > Machine > Type)
   - Search functionality
   - Selected sensor state

3. **Realtime Store** (`stores/realtime.ts`)
   - NATS WebSocket connection state
   - Real-time data buffer (ring buffer)
   - Subscription management
   - Data merging with historical data

4. **Alarms Store** (`stores/alarms.ts`)
   - Active alarms list
   - Alarm history
   - Acknowledgment state
   - Shelving state

5. **Audit Store** (`stores/audit.ts`)
   - Audit trail entries
   - Filtering and search
   - Export functionality

---

## 5. Feature Breakdown & User Stories

### 5.1 Epic 7.1: Authentication & Authorization

**Story 7.1.1: User Login**
- As a user, I want to log in with username/password
- So that I can access the historian platform securely

**Acceptance Criteria:**
- [ ] Login form with username and password fields
- [ ] JWT token stored in localStorage/sessionStorage
- [ ] Automatic redirect to dashboard after successful login
- [ ] Error messages for invalid credentials
- [ ] "Remember me" option
- [ ] Password visibility toggle

**Story 7.1.2: FDA Re-Authentication**
- As a compliance officer, I want critical actions to require re-authentication
- So that we maintain FDA 21 CFR Part 11 compliance

**Acceptance Criteria:**
- [ ] Re-auth modal appears for critical actions (alarm ack, config changes)
- [ ] Short-lived signing token issued after re-auth
- [ ] Re-auth events logged to audit trail
- [ ] Timeout after 5 minutes of inactivity

**Story 7.1.3: Role-Based Access Control**
- As an admin, I want different user roles to have different permissions
- So that we maintain proper access control

**Acceptance Criteria:**
- [ ] Roles: Admin, Operator, Viewer
- [ ] Route guards based on user role
- [ ] UI elements hidden/disabled based on permissions
- [ ] Unauthorized access attempts logged

---

### 5.2 Epic 7.2: Dashboard & Visualization

**Story 7.2.1: Main Dashboard**
- As an operator, I want to see an overview of the entire plant
- So that I can quickly identify issues

**Acceptance Criteria:**
- [ ] KPI cards (total sensors, active alarms, system health)
- [ ] Recent alarms widget
- [ ] Top 10 sensors by deviation
- [ ] System status indicators
- [ ] Responsive grid layout
- [ ] Dark/Light mode toggle

**Story 7.2.2: Sensor List with Faceted Search**
- As a maintenance engineer, I want to filter sensors by hierarchy
- So that I can find specific sensors quickly

**Acceptance Criteria:**
- [ ] Sidebar with Factory > Line > Machine > Type filters
- [ ] Text search by Tag ID or Description
- [ ] Sensor count badges on each filter
- [ ] Virtual scrolling for 10,000+ sensors
- [ ] Sparkline preview for each sensor
- [ ] Click to open sensor detail modal

**Story 7.2.3: High-Performance Trend Chart**
- As a process engineer, I want to see interactive time-series charts
- So that I can analyze sensor behavior

**Acceptance Criteria:**
- [ ] uPlot-based chart with pan/zoom
- [ ] Time range selector (1h, 6h, 24h, 7d, 30d, custom)
- [ ] Real-time data streaming via NATS
- [ ] Historical data fetching via HTTP API
- [ ] LTTB downsampling for large datasets
- [ ] Export to CSV/JSON
- [ ] Multiple series overlay
- [ ] Crosshair with value tooltip
- [ ] 60 FPS performance

**Story 7.2.4: Sensor Detail Modal**
- As an operator, I want detailed information about a sensor
- So that I can understand its current state

**Acceptance Criteria:**
- [ ] Modal with sensor metadata (ID, description, unit, location)
- [ ] Current value with quality indicator
- [ ] Trend chart (last 24h by default)
- [ ] Statistics (min, max, avg, stddev)
- [ ] Related alarms
- [ ] Export button
- [ ] Keyboard navigation (ESC to close, arrow keys to navigate)

---

### 5.3 Epic 7.3: Alarm Management

**Story 7.3.1: Alarm List View**
- As an operator, I want to see all active alarms
- So that I can respond to critical issues

**Acceptance Criteria:**
- [ ] Table with columns: Timestamp, Sensor, Message, Severity, State
- [ ] Color-coded severity (Critical=Red, High=Orange, Medium=Yellow, Low=Blue)
- [ ] Filter by severity, state (Active, Acknowledged, Shelved)
- [ ] Sort by timestamp, severity
- [ ] Real-time updates via NATS
- [ ] Alarm count badge in header

**Story 7.3.2: Alarm Acknowledgment**
- As an operator, I want to acknowledge alarms
- So that I can indicate I'm aware of the issue

**Acceptance Criteria:**
- [ ] "Acknowledge" button on each alarm
- [ ] Re-authentication required (FDA compliance)
- [ ] Acknowledgment logged to audit trail
- [ ] Alarm state changes to "Acknowledged"
- [ ] Acknowledged by user and timestamp displayed

**Story 7.3.3: Alarm Shelving**
- As a maintenance engineer, I want to shelve alarms during maintenance
- So that I don't get false alarms

**Acceptance Criteria:**
- [ ] "Shelve" button with duration selector (1h, 4h, 8h, 24h)
- [ ] Re-authentication required
- [ ] Shelving logged to audit trail
- [ ] Shelved alarms shown in separate tab
- [ ] Auto-unshelve after duration expires

---

### 5.4 Epic 7.4: Audit Trail

**Story 7.4.1: Audit Log Viewer**
- As a compliance officer, I want to view the audit trail
- So that I can verify system integrity

**Acceptance Criteria:**
- [ ] Table with columns: Timestamp, User, Action, Entity, Details
- [ ] Filter by date range, user, action type
- [ ] Search by entity ID
- [ ] Pagination (100 entries per page)
- [ ] Export to CSV
- [ ] Hash chain verification indicator

**Story 7.4.2: Audit Export**
- As a compliance officer, I want to export audit logs
- So that I can provide them to auditors

**Acceptance Criteria:**
- [ ] Export button with date range selector
- [ ] CSV format with all fields
- [ ] Filename: `audit_log_{start}_{end}.csv`
- [ ] Hash chain included in export
- [ ] Re-authentication required

---

### 5.5 Epic 7.5: System Administration

**Story 7.5.1: User Management**
- As an admin, I want to manage user accounts
- So that I can control access to the system

**Acceptance Criteria:**
- [ ] User list with username, email, role, status
- [ ] Create new user form
- [ ] Edit user (change role, disable/enable)
- [ ] Delete user (with confirmation)
- [ ] All actions logged to audit trail

**Story 7.5.2: System Settings**
- As an admin, I want to configure system settings
- So that I can customize the platform

**Acceptance Criteria:**
- [ ] Settings page with tabs (General, Security, Notifications)
- [ ] Dark/Light mode preference
- [ ] Default time range for charts
- [ ] Alarm notification settings
- [ ] Session timeout configuration
- [ ] All changes logged to audit trail

---

## 6. Design System

### 6.1 Color Palette

**Primary Colors:**
```css
/* Indigo - Professional, trustworthy */
--color-primary-50: #eef2ff;
--color-primary-100: #e0e7ff;
--color-primary-500: #6366f1;
--color-primary-600: #4f46e5;
--color-primary-700: #4338ca;
--color-primary-900: #312e81;

/* Slate - Neutral, modern */
--color-slate-50: #f8fafc;
--color-slate-100: #f1f5f9;
--color-slate-500: #64748b;
--color-slate-700: #334155;
--color-slate-900: #0f172a;
```

**Semantic Colors:**
```css
/* Success */
--color-success: #10b981;

/* Warning */
--color-warning: #f59e0b;

/* Error */
--color-error: #ef4444;

/* Info */
--color-info: #3b82f6;
```

**Alarm Severity Colors:**
```css
--alarm-critical: #dc2626;  /* Red */
--alarm-high: #ea580c;      /* Orange */
--alarm-medium: #eab308;    /* Yellow */
--alarm-low: #3b82f6;       /* Blue */
```

### 6.2 Typography

**Font Family:**
- Primary: `Inter` (Google Fonts)
- Monospace: `JetBrains Mono` (for sensor IDs, timestamps)

**Font Sizes:**
```css
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
--text-4xl: 2.25rem;   /* 36px */
```

### 6.3 Spacing

```css
--spacing-1: 0.25rem;  /* 4px */
--spacing-2: 0.5rem;   /* 8px */
--spacing-3: 0.75rem;  /* 12px */
--spacing-4: 1rem;     /* 16px */
--spacing-6: 1.5rem;   /* 24px */
--spacing-8: 2rem;     /* 32px */
--spacing-12: 3rem;    /* 48px */
```

### 6.4 Shadows

```css
--shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
--shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
--shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
--shadow-xl: 0 20px 25px -5px rgb(0 0 0 / 0.1);
```

### 6.5 Animations

```css
/* Smooth transitions */
--transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1);
--transition-base: 300ms cubic-bezier(0.4, 0, 0.2, 1);
--transition-slow: 500ms cubic-bezier(0.4, 0, 0.2, 1);

/* Micro-interactions */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
```

---

## 7. Performance Requirements

### 7.1 Loading Performance
- **Time to Interactive (TTI):** < 2 seconds
- **First Contentful Paint (FCP):** < 1 second
- **Largest Contentful Paint (LCP):** < 2.5 seconds

### 7.2 Runtime Performance
- **Chart Rendering:** 60 FPS for real-time updates
- **List Rendering:** Virtual scrolling for 10,000+ items
- **Memory Usage:** < 200MB for typical session
- **Bundle Size:** < 500KB (gzipped)

### 7.3 Data Handling
- **Real-time Updates:** < 100ms latency from NATS
- **Historical Query:** < 500ms for 1 year of data
- **Downsampling:** LTTB to max 5,000 points per series
- **Export:** Streaming for > 1M rows

---

## 8. Implementation Phases

### Phase 1: Foundation (Week 1)
**Goal:** Set up project infrastructure

**Tasks:**
1. Initialize Vue 3 + Vite + TypeScript project
2. Configure TailwindCSS, ESLint, Prettier
3. Set up Pinia stores
4. Configure Vue Router
5. Create base UI components (Button, Input, Card, etc.)
6. Set up Docker build and Nginx configuration
7. Create development environment documentation

**Deliverables:**
- Working dev server with hot reload
- Dockerized build
- Base component library
- Project documentation

---

### Phase 2: Authentication & Layout (Week 2)
**Goal:** Implement authentication and app shell

**Tasks:**
1. Create Auth API client
2. Implement Auth store (Pinia)
3. Build Login view
4. Create route guards
5. Build app layout (Header, Sidebar, Footer)
6. Implement dark/light mode
7. Add error boundary
8. Create loading states

**Deliverables:**
- Working login flow
- Protected routes
- Responsive app layout
- Theme switcher

---

### Phase 3: Dashboard & Sensors (Week 3-4)
**Goal:** Core visualization features

**Tasks:**
1. Create Engine API client
2. Implement Sensors store
3. Build Dashboard view with KPI cards
4. Create Sensor list with faceted search
5. Implement Sparkline component
6. Build TrendChart with uPlot
7. Create Sensor detail modal
8. Implement NATS real-time streaming
9. Add data export functionality

**Deliverables:**
- Working dashboard
- Sensor list with search/filter
- Interactive trend charts
- Real-time data updates

---

### Phase 4: Alarms (Week 5)
**Goal:** Alarm management interface

**Tasks:**
1. Create Alarm API client
2. Implement Alarms store
3. Build Alarms view
4. Implement alarm acknowledgment
5. Add alarm shelving
6. Create alarm notification system
7. Add real-time alarm updates via NATS

**Deliverables:**
- Working alarm management
- ISA 18.2 compliant interface
- Real-time alarm notifications

---

### Phase 5: Audit Trail (Week 6)
**Goal:** Compliance and audit features

**Tasks:**
1. Create Audit API client
2. Implement Audit store
3. Build Audit log viewer
4. Add filtering and search
5. Implement audit export
6. Add hash chain verification UI
7. Create re-authentication modal

**Deliverables:**
- Working audit trail viewer
- FDA 21 CFR Part 11 compliant interface
- Audit export functionality

---

### Phase 6: Testing & Polish (Week 7-8)
**Goal:** Quality assurance and optimization

**Tasks:**
1. Write unit tests (Vitest)
2. Write E2E tests (Playwright)
3. Performance optimization
4. Accessibility audit (WCAG 2.1 AA)
5. Browser compatibility testing
6. Mobile responsiveness testing
7. Security audit
8. Documentation completion

**Deliverables:**
- 80%+ test coverage
- Performance benchmarks met
- Accessibility compliant
- Production-ready application

---

## 9. API Integration Specifications

### 9.1 Engine Service APIs

**Base URL:** `http://engine:8081/api/v1`

**Endpoints:**

1. **Get Metadata**
   ```
   GET /metadata
   Response: {
     sensors: [
       {
         id: string,
         desc: string,
         factory: string,
         line: string,
         machine: string,
         type: string,
         unit: string
       }
     ]
   }
   ```

2. **Export Data**
   ```
   GET /export?tag_id={id}&start={ts}&end={ts}&format={csv|json}
   Response: Stream (CSV or JSON)
   ```

3. **Query (gRPC - Optional)**
   ```
   gRPC: HistorianQuery/Query
   Request: { sensor_id, start, end, max_points }
   Response: { points: [{ts, val, quality}] }
   ```

### 9.2 Auth Service APIs

**Base URL:** `http://auth:8080/api/v1`

**Endpoints:**

1. **Login**
   ```
   POST /login
   Body: { username, password }
   Response: { token, user: { id, username, role } }
   ```

2. **Re-Authenticate**
   ```
   POST /re-auth
   Headers: { Authorization: Bearer {token} }
   Body: { password }
   Response: { signing_token, expires_at }
   ```

3. **Logout**
   ```
   POST /logout
   Headers: { Authorization: Bearer {token} }
   Response: { success: true }
   ```

### 9.3 Alarm Service APIs

**Base URL:** `http://alarm:8083/api/v1`

**Endpoints:**

1. **Get Alarms**
   ```
   GET /alarms?state={active|acknowledged|shelved}&severity={critical|high|medium|low}
   Response: {
     alarms: [
       {
         id, sensor_id, message, severity, state,
         triggered_at, acknowledged_at, acknowledged_by,
         shelved_until
       }
     ]
   }
   ```

2. **Acknowledge Alarm**
   ```
   POST /alarms/{id}/acknowledge
   Headers: { Authorization: Bearer {signing_token} }
   Response: { success: true }
   ```

3. **Shelve Alarm**
   ```
   POST /alarms/{id}/shelve
   Headers: { Authorization: Bearer {signing_token} }
   Body: { duration_ms }
   Response: { success: true, shelved_until }
   ```

### 9.4 Audit Service APIs

**Base URL:** `http://audit:8082/api/v1`

**Endpoints:**

1. **Get Audit Logs**
   ```
   GET /audit?start={ts}&end={ts}&user={username}&action={type}
   Response: {
     entries: [
       {
         id, timestamp, user, action, entity_type,
         entity_id, details, hash, prev_hash
       }
     ],
     total, page, per_page
   }
   ```

2. **Export Audit Logs**
   ```
   GET /audit/export?start={ts}&end={ts}
   Headers: { Authorization: Bearer {signing_token} }
   Response: Stream (CSV)
   ```

### 9.5 NATS WebSocket Subscriptions

**URL:** `ws://nats:8222`

**Subjects:**

1. **Real-time Sensor Data**
   ```
   Subject: data.{factory}.{line}.{machine}.{type}.{sensor_id}
   Payload: Protobuf SensorData {
     sensor_id, timestamp_ms, value, quality
   }
   ```

2. **Alarm Events**
   ```
   Subject: sys.alarm.{triggered|acknowledged|shelved}
   Payload: Protobuf AlarmEvent {
     alarm_id, sensor_id, severity, message, state
   }
   ```

3. **Audit Events**
   ```
   Subject: sys.audit.{action}
   Payload: Protobuf AuditEvent {
     user, action, entity_type, entity_id, timestamp
   }
   ```

---

## 10. Security Considerations

### 10.1 Authentication
- JWT tokens with 1-hour expiration
- Refresh tokens with 7-day expiration
- Secure storage (httpOnly cookies or encrypted localStorage)
- CSRF protection

### 10.2 Authorization
- Role-based access control (RBAC)
- Route guards
- API permission checks
- Audit logging of all actions

### 10.3 Data Protection
- TLS 1.2+ for all communications
- Input validation and sanitization
- XSS protection (Vue's built-in escaping)
- Content Security Policy (CSP) headers

### 10.4 FDA Compliance
- Re-authentication for critical actions
- Audit trail with hash chain
- Electronic signatures
- Tamper-evident logging

---

## 11. Deployment Strategy

### 11.1 Docker Build

**Dockerfile:**
```dockerfile
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Production
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Nginx Configuration:**
```nginx
server {
  listen 80;
  server_name _;
  root /usr/share/nginx/html;
  index index.html;

  # Gzip compression
  gzip on;
  gzip_types text/plain text/css application/json application/javascript;

  # SPA routing
  location / {
    try_files $uri $uri/ /index.html;
  }

  # API proxy
  location /api/ {
    proxy_pass http://engine:8081/api/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
  }

  # Auth API proxy
  location /auth/ {
    proxy_pass http://auth:8080/api/v1/;
  }

  # Alarm API proxy
  location /alarm/ {
    proxy_pass http://alarm:8083/api/v1/;
  }

  # Audit API proxy
  location /audit/ {
    proxy_pass http://audit:8082/api/v1/;
  }

  # NATS WebSocket proxy
  location /nats/ {
    proxy_pass http://nats:8222/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";
  }

  # Security headers
  add_header X-Frame-Options "SAMEORIGIN" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header X-XSS-Protection "1; mode=block" always;
  add_header Referrer-Policy "no-referrer-when-downgrade" always;
}
```

### 11.2 Docker Compose Integration

```yaml
web-ui:
  build:
    context: ../web-ui
    dockerfile: Dockerfile
  container_name: ops-web-ui
  ports:
    - "3000:80"
  depends_on:
    - engine
    - auth
    - alarm
    - audit
  networks:
    - historian-net
  restart: unless-stopped
```

---

## 12. Testing Strategy

### 12.1 Unit Tests (Vitest)
- Component logic testing
- Store (Pinia) testing
- Utility function testing
- API client testing (mocked)
- Target: 80%+ coverage

### 12.2 Integration Tests
- Component integration
- Store + API integration
- Router navigation
- Authentication flow

### 12.3 E2E Tests (Playwright)
- Critical user flows
- Login/logout
- Sensor search and detail view
- Alarm acknowledgment
- Audit log export

### 12.4 Performance Tests
- Lighthouse CI
- Bundle size monitoring
- Chart rendering benchmarks
- Memory leak detection

---

## 13. Documentation Requirements

### 13.1 Developer Documentation
- Setup guide
- Architecture overview
- Component API documentation
- Store documentation
- Testing guide
- Deployment guide

### 13.2 User Documentation
- User manual
- Feature guides
- FAQ
- Troubleshooting

### 13.3 API Documentation
- OpenAPI/Swagger specs
- Integration examples
- Authentication guide

---

## 14. Success Metrics

### 14.1 Technical Metrics
- ✅ Build time < 30 seconds
- ✅ Bundle size < 500KB (gzipped)
- ✅ TTI < 2 seconds
- ✅ 60 FPS chart rendering
- ✅ 80%+ test coverage
- ✅ Zero critical security vulnerabilities

### 14.2 User Experience Metrics
- ✅ < 3 clicks to any feature
- ✅ < 5 seconds to find a sensor
- ✅ < 1 second to acknowledge an alarm
- ✅ Mobile responsive (320px - 4K)
- ✅ WCAG 2.1 AA compliant

### 14.3 Business Metrics
- ✅ FDA 21 CFR Part 11 compliant
- ✅ ISA 18.2 compliant
- ✅ Zero data loss
- ✅ 99.9% uptime

---

## 15. Risk Assessment

### 15.1 Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Performance issues with 10k+ sensors | High | Medium | Virtual scrolling, lazy loading, pagination |
| Real-time data latency | High | Low | NATS WebSocket optimization, local buffering |
| Browser compatibility | Medium | Medium | Polyfills, progressive enhancement |
| Bundle size bloat | Medium | Medium | Code splitting, tree shaking, lazy loading |

### 15.2 Schedule Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Scope creep | High | High | Strict MVP definition, change control |
| Integration delays | Medium | Medium | Early API contract definition, mocking |
| Testing delays | Medium | Low | Parallel testing, automated CI/CD |

### 15.3 Compliance Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| FDA audit failure | Critical | Low | Regular compliance reviews, documentation |
| Security vulnerabilities | High | Medium | Security audits, dependency scanning |
| Data integrity issues | Critical | Low | Hash chain verification, audit logging |

---

## 16. Next Steps

### Immediate Actions (This Week)
1. ✅ Delete old `/viz` directory (React-based)
2. ✅ Create this planning document
3. 🔲 Get stakeholder approval on architecture
4. 🔲 Set up new `web-ui` project directory
5. 🔲 Initialize Vue 3 + Vite + TypeScript
6. 🔲 Configure development environment

### Phase 1 Kickoff (Next Week)
1. 🔲 Create base project structure
2. 🔲 Set up TailwindCSS and design system
3. 🔲 Build base UI component library
4. 🔲 Configure Docker build
5. 🔲 Set up CI/CD pipeline

---

## 17. Appendix

### A. Technology Comparison: Vue vs React

| Aspect | Vue 3 | React (Current) | Decision |
|--------|-------|-----------------|----------|
| Learning Curve | Easier | Moderate | ✅ Vue |
| Performance | Excellent | Excellent | ✅ Tie |
| TypeScript Support | Excellent | Excellent | ✅ Tie |
| Bundle Size | Smaller (~40KB) | Larger (~100KB) | ✅ Vue |
| Ecosystem | Growing | Mature | ⚠️ React |
| Official Router | Yes (Vue Router) | No (React Router) | ✅ Vue |
| Official State Mgmt | Yes (Pinia) | No (Redux/Zustand) | ✅ Vue |
| Developer Experience | Excellent | Good | ✅ Vue |
| Team Familiarity | New | Existing | ⚠️ React |

**Decision Rationale:**
- Vue 3 offers better DX with Composition API
- Smaller bundle size and better performance
- Official solutions for routing and state management
- Modern, clean syntax with `<script setup>`
- Better TypeScript integration out of the box

### B. Design Inspiration

**Reference Applications:**
- Grafana (Time-series visualization)
- Datadog (Monitoring dashboard)
- Sentry (Error tracking UI)
- Linear (Modern, fast UI)
- Vercel Dashboard (Clean, premium design)

**Design Principles:**
- **Clarity:** Information hierarchy, clear typography
- **Efficiency:** Keyboard shortcuts, quick actions
- **Consistency:** Design system, reusable components
- **Performance:** 60 FPS, instant feedback
- **Accessibility:** WCAG 2.1 AA, keyboard navigation

### C. Glossary

- **LTTB:** Largest-Triangle-Three-Buckets (downsampling algorithm)
- **TTI:** Time to Interactive
- **FCP:** First Contentful Paint
- **LCP:** Largest Contentful Paint
- **WCAG:** Web Content Accessibility Guidelines
- **CSP:** Content Security Policy
- **RBAC:** Role-Based Access Control
- **FDA 21 CFR Part 11:** FDA regulation for electronic records
- **ISA 18.2:** ISA standard for alarm management

---

## Document Control

**Version:** 1.0
**Last Updated:** 2025-12-11
**Author:** PM Agent (Ahmet)
**Reviewers:** Architect, UX Designer, Tech Lead
**Status:** Draft - Awaiting Approval

**Change Log:**
- 2025-12-11: Initial document creation
- TBD: Stakeholder review and approval
- TBD: Architecture review
- TBD: Final approval and phase 1 kickoff

---

**END OF DOCUMENT**
