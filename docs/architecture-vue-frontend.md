# Vue Frontend Architecture Decision

**Date:** 2025-12-11
**Status:** Approved
**Decision Maker:** PM (Ahmet)

---

## Context

The current Historian platform has a React-based frontend (`/viz`) that was developed during Epic 4. Based on recent implementation experience and production stability issues (white screen errors, error boundary requirements), we need to rebuild the web interface with a modern, production-ready approach.

### Current System Overview

**Backend Infrastructure:**
- **Engine Service (Rust):** Port 8081 - HTTP Export API, Metadata API, gRPC Query API
- **Auth Service (Go):** Port 8080 - JWT authentication, FDA re-auth
- **Alarm Service (Go):** Port 8083 - ISA 18.2 alarm management
- **Audit Service (Go):** Port 8082 - FDA 21 CFR Part 11 audit trail
- **NATS JetStream:** Port 4222 - Message bus for real-time events
- **MinIO:** Port 9000 - S3-compatible object storage for tiered data
- **PostgreSQL:** Port 5432 - Relational data (users, alarms, audit)

**Current Frontend (React - Deprecated):**
- React 19.2 + Vite 7
- Zustand for state management
- TailwindCSS 4 for styling
- uPlot for time-series charts
- NATS WebSocket for real-time data

**Issues with Current Implementation:**
- Error boundary needed for production stability
- White screen errors in production
- TypeScript strict mode issues
- Complex state management patterns
- Limited mobile responsiveness

---

## Decision

**We will rebuild the web interface using Vue 3 with TypeScript.**

### Technology Stack

**Core Framework:**
- Vue 3.4+ (Composition API with `<script setup>`)
- TypeScript 5.9+ (Strict mode)
- Vite 7+ (Build tool)

**State Management:**
- Pinia (Official Vue state management)

**UI Framework:**
- TailwindCSS 4.x
- Headless UI (Vue)
- Heroicons
- VueUse (Composition utilities)

**Data Visualization:**
- uPlot (High-performance time-series charts)
- D3.js (Optional for custom visualizations)

**Real-time Communication:**
- nats.ws (NATS WebSocket client)
- Axios (HTTP client)

**Testing:**
- Vitest (Unit tests)
- @vue/test-utils (Component tests)
- Playwright (E2E tests)

**Deployment:**
- Docker (Multi-stage build)
- Nginx (Static serving + API proxy)

---

## Rationale

### Why Vue 3 over React?

1. **Better Developer Experience:**
   - Composition API is more intuitive than React hooks
   - `<script setup>` reduces boilerplate
   - Official router and state management (no decision fatigue)
   - Better TypeScript integration out of the box

2. **Performance:**
   - Smaller bundle size (~40KB vs ~100KB for React)
   - Faster initial load time
   - Better reactivity system (Proxy-based)

3. **Ecosystem:**
   - Official solutions for common needs (Vue Router, Pinia)
   - Excellent documentation
   - Growing enterprise adoption

4. **Lessons from React Implementation:**
   - Zustand worked well → Pinia is similar but official
   - uPlot performed excellently → Keep it
   - TailwindCSS was good → Keep it
   - Error boundaries needed → Vue has better error handling

### Why Pinia over Vuex?

- Official recommendation from Vue team
- TypeScript-first design
- Simpler API (no mutations, just actions)
- Better DevTools integration
- Modular stores by default

### Why Keep uPlot?

- Proven performance (60 FPS with 10k+ points)
- Canvas-based rendering (better than SVG for time-series)
- Small bundle size (~50KB)
- Framework-agnostic (works with Vue)

---

## Architecture Decisions

### 1. Project Structure

```
web-ui/                              # New Vue application
├── src/
│   ├── main.ts                      # Entry point
│   ├── App.vue                      # Root component
│   ├── router/                      # Vue Router
│   ├── stores/                      # Pinia stores
│   │   ├── auth.ts
│   │   ├── sensors.ts
│   │   ├── realtime.ts
│   │   ├── alarms.ts
│   │   └── audit.ts
│   ├── composables/                 # Reusable composition functions
│   │   ├── useNatsStream.ts
│   │   ├── useHistoryQuery.ts
│   │   └── useAuth.ts
│   ├── services/                    # API clients
│   │   ├── engine.ts
│   │   ├── auth.ts
│   │   ├── alarm.ts
│   │   └── audit.ts
│   ├── components/                  # Reusable components
│   │   ├── ui/                      # Base UI components
│   │   ├── charts/                  # Chart components
│   │   ├── layout/                  # Layout components
│   │   └── common/                  # Common components
│   ├── views/                       # Page components
│   │   ├── DashboardView.vue
│   │   ├── SensorsView.vue
│   │   ├── AlarmsView.vue
│   │   ├── AuditView.vue
│   │   └── LoginView.vue
│   ├── types/                       # TypeScript types
│   └── utils/                       # Utility functions
└── tests/
    ├── unit/
    └── e2e/
```

### 2. State Management Architecture

**Pinia Stores:**

1. **Auth Store:** User authentication, JWT tokens, permissions
2. **Sensors Store:** Sensor metadata, filtering, search
3. **Realtime Store:** NATS connection, real-time data buffer
4. **Alarms Store:** Active alarms, alarm history
5. **Audit Store:** Audit trail entries, filtering

**Pattern:**
```typescript
// stores/sensors.ts
import { defineStore } from 'pinia'

export const useSensorsStore = defineStore('sensors', () => {
  // State
  const sensors = ref<Sensor[]>([])
  const selectedSensor = ref<Sensor | null>(null)
  const filters = ref<SensorFilters>({})

  // Getters
  const filteredSensors = computed(() => {
    return sensors.value.filter(/* filter logic */)
  })

  // Actions
  async function fetchSensors() {
    const data = await engineApi.getMetadata()
    sensors.value = data.sensors
  }

  return {
    sensors,
    selectedSensor,
    filters,
    filteredSensors,
    fetchSensors
  }
})
```

### 3. API Integration

**Base URLs (via Nginx proxy):**
- Engine API: `/api/*` → `http://engine:8081/api/v1/*`
- Auth API: `/auth/*` → `http://auth:8080/api/v1/*`
- Alarm API: `/alarm/*` → `http://alarm:8083/api/v1/*`
- Audit API: `/audit/*` → `http://audit:8082/api/v1/*`
- NATS WebSocket: `/nats/*` → `ws://nats:8222/*`

**API Client Pattern:**
```typescript
// services/engine.ts
import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const engineApi = {
  async getMetadata() {
    const { data } = await client.get('/metadata')
    return data
  },
  
  async exportData(params: ExportParams) {
    const { data } = await client.get('/export', { params })
    return data
  }
}
```

### 4. Real-time Data Flow

```
┌─────────────┐
│   Sensors   │
└──────┬──────┘
       │ Modbus/OPC UA
       ▼
┌──────────────┐
│  Ingestor    │
└──────┬───────┘
       │ NATS Publish
       ▼
┌──────────────┐      WebSocket      ┌─────────────┐
│     NATS     │ ◀──────────────────▶│   Vue UI    │
│  JetStream   │                     │             │
└──────┬───────┘                     │ - Pinia     │
       │                             │ - uPlot     │
       ▼                             │ - Realtime  │
┌──────────────┐      HTTP API       └─────────────┘
│    Engine    │ ◀──────────────────▶
│  (RocksDB)   │
└──────────────┘
```

**Composable Pattern:**
```typescript
// composables/useNatsStream.ts
import { connect } from 'nats.ws'
import { useRealtimeStore } from '@/stores/realtime'

export function useNatsStream(subject: string) {
  const store = useRealtimeStore()
  const connection = ref<NatsConnection | null>(null)

  async function subscribe() {
    connection.value = await connect({ servers: '/nats' })
    const sub = connection.value.subscribe(subject)
    
    for await (const msg of sub) {
      const data = decodeSensorData(msg.data)
      store.addDataPoint(data)
    }
  }

  onUnmounted(() => {
    connection.value?.close()
  })

  return { subscribe }
}
```

### 5. Component Architecture

**Base UI Components (Headless UI + TailwindCSS):**
- Button, Input, Select, Checkbox, Radio
- Modal, Dropdown, Tooltip, Popover
- Card, Badge, Alert, Toast
- Table, Pagination

**Chart Components (uPlot):**
- TrendChart: Main time-series chart with zoom/pan
- Sparkline: Inline mini chart for sensor list
- UPlotWrapper: Base wrapper for uPlot lifecycle

**Layout Components:**
- AppHeader: Top navigation with user menu
- AppSidebar: Collapsible sidebar with navigation
- AppFooter: Footer with system info

**Feature Components:**
- SensorList: Virtual scrolling list with filters
- SensorDetail: Modal with trend chart and metadata
- AlarmList: Real-time alarm table
- AuditLog: Paginated audit trail viewer

### 6. Routing

```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    redirect: '/dashboard',
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/sensors',
    component: () => import('@/views/SensorsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/alarms',
    component: () => import('@/views/AlarmsView.vue'),
    meta: { requiresAuth: true, requiredRole: 'operator' }
  },
  {
    path: '/audit',
    component: () => import('@/views/AuditView.vue'),
    meta: { requiresAuth: true, requiredRole: 'admin' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Route guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiredRole && !authStore.hasRole(to.meta.requiredRole)) {
    next('/unauthorized')
  } else {
    next()
  }
})

export default router
```

---

## Integration with Existing Backend

### Engine Service Integration

**Metadata API:**
```typescript
// Fetch sensor metadata on app load
const sensorsStore = useSensorsStore()
await sensorsStore.fetchSensors()

// Response format (from existing API):
{
  sensors: [
    {
      id: "factory1.line1.machine1.temp.sensor1",
      desc: "Temperature Sensor 1",
      factory: "factory1",
      line: "line1",
      machine: "machine1",
      type: "temp",
      unit: "°C"
    }
  ]
}
```

**Export API:**
```typescript
// Export sensor data to CSV
async function exportToCsv(sensorId: string, start: number, end: number) {
  const response = await fetch(
    `/api/export?tag_id=${sensorId}&start=${start}&end=${end}&format=csv`
  )
  const blob = await response.blob()
  downloadFile(blob, `${sensorId}_${start}_${end}.csv`)
}
```

### Auth Service Integration

**Login Flow:**
```typescript
// stores/auth.ts
async function login(username: string, password: string) {
  const { data } = await authApi.login({ username, password })
  token.value = data.token
  user.value = data.user
  localStorage.setItem('token', data.token)
  router.push('/dashboard')
}
```

**Re-Authentication (FDA Compliance):**
```typescript
// composables/useReAuth.ts
export function useReAuth() {
  const authStore = useAuthStore()
  
  async function requireReAuth(action: () => Promise<void>) {
    const modal = await showReAuthModal()
    const password = await modal.getPassword()
    
    const { signing_token } = await authApi.reAuth(password)
    
    // Execute action with signing token
    await action()
    
    // Log to audit trail
    auditStore.logAction('re-auth', { action: action.name })
  }
  
  return { requireReAuth }
}
```

### Alarm Service Integration

**Real-time Alarms:**
```typescript
// composables/useAlarms.ts
export function useAlarms() {
  const alarmsStore = useAlarmsStore()
  const { subscribe } = useNatsStream('sys.alarm.>')
  
  onMounted(async () => {
    // Fetch initial alarms
    await alarmsStore.fetchAlarms()
    
    // Subscribe to real-time updates
    subscribe((msg) => {
      const alarm = decodeAlarmEvent(msg.data)
      alarmsStore.addAlarm(alarm)
    })
  })
}
```

**Acknowledge Alarm:**
```typescript
async function acknowledgeAlarm(alarmId: string) {
  await requireReAuth(async () => {
    await alarmApi.acknowledge(alarmId)
    alarmsStore.updateAlarmState(alarmId, 'acknowledged')
  })
}
```

### Audit Service Integration

**Audit Log Viewer:**
```typescript
// views/AuditView.vue
const auditStore = useAuditStore()
const filters = ref({
  start: Date.now() - 86400000, // Last 24h
  end: Date.now(),
  user: '',
  action: ''
})

async function loadAuditLogs() {
  await auditStore.fetchLogs(filters.value)
}

async function exportAuditLogs() {
  await requireReAuth(async () => {
    const blob = await auditApi.export(filters.value)
    downloadFile(blob, `audit_${filters.value.start}_${filters.value.end}.csv`)
  })
}
```

---

## Performance Optimizations

### 1. Code Splitting
```typescript
// Lazy load routes
const routes = [
  {
    path: '/dashboard',
    component: () => import('@/views/DashboardView.vue')
  }
]
```

### 2. Virtual Scrolling
```vue
<!-- For large sensor lists -->
<template>
  <RecycleScroller
    :items="sensors"
    :item-size="48"
    key-field="id"
  >
    <template #default="{ item }">
      <SensorListItem :sensor="item" />
    </template>
  </RecycleScroller>
</template>
```

### 3. Debounced Search
```typescript
import { useDebounceFn } from '@vueuse/core'

const searchQuery = ref('')
const debouncedSearch = useDebounceFn((query: string) => {
  sensorsStore.setSearchQuery(query)
}, 300)

watch(searchQuery, debouncedSearch)
```

### 4. Memoized Computations
```typescript
// Expensive filtering computation
const filteredSensors = computed(() => {
  return sensors.value.filter(sensor => {
    // Complex filtering logic
  })
})
```

---

## Security Considerations

### 1. JWT Token Management
```typescript
// Axios interceptor for auth token
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle token expiration
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      authStore.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

### 2. Input Sanitization
```typescript
import DOMPurify from 'dompurify'

function sanitizeInput(input: string): string {
  return DOMPurify.sanitize(input)
}
```

### 3. Content Security Policy
```nginx
# nginx.conf
add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self' ws: wss:;" always;
```

---

## Migration Strategy

### Phase 1: Parallel Development
- Keep existing React app running
- Develop new Vue app in `/web-ui`
- Test with same backend services

### Phase 2: Feature Parity
- Implement all features from React app
- Add new features (alarms, audit)
- Comprehensive testing

### Phase 3: Cutover
- Update docker-compose.yml to use new Vue app
- Deprecate old React app
- Monitor for issues

### Phase 4: Cleanup
- Remove `/viz` directory
- Update documentation
- Archive old code

---

## Success Criteria

### Technical
- ✅ Build time < 30 seconds
- ✅ Bundle size < 500KB (gzipped)
- ✅ TTI < 2 seconds
- ✅ 60 FPS chart rendering
- ✅ 80%+ test coverage

### Functional
- ✅ All React features implemented
- ✅ Alarm management working
- ✅ Audit trail working
- ✅ Real-time updates working
- ✅ Mobile responsive

### Compliance
- ✅ FDA 21 CFR Part 11 compliant
- ✅ ISA 18.2 compliant
- ✅ Audit trail with hash chain
- ✅ Re-authentication working

---

## Next Steps

1. ✅ Create planning document
2. 🔲 Get stakeholder approval
3. 🔲 Initialize Vue project
4. 🔲 Set up development environment
5. 🔲 Begin Phase 1 implementation

---

## References

- [Vue 3 Documentation](https://vuejs.org/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [uPlot Documentation](https://github.com/leeoniya/uPlot)
- [TailwindCSS Documentation](https://tailwindcss.com/)
- [Headless UI Documentation](https://headlessui.com/)

---

**Document Status:** Approved
**Last Updated:** 2025-12-11
**Next Review:** After Phase 1 completion
