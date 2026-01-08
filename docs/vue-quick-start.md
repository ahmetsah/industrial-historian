# Quick Start Guide: Modern Vue Web Interface

**Last Updated:** 2025-12-11
**Status:** Ready for Development

---

## 📋 Overview

This guide provides quick commands and references for developing the new Vue-based web interface for the Historian platform.

---

## 🚀 Quick Start Commands

### 1. Initialize New Vue Project

```bash
# Navigate to project root
cd /home/ahmet/historian

# Create new Vue project with Vite
npm create vue@latest web-ui

# When prompted, select:
# ✅ TypeScript
# ✅ Vue Router
# ✅ Pinia
# ✅ Vitest
# ✅ ESLint
# ✅ Prettier

# Navigate to project
cd web-ui

# Install dependencies
npm install

# Install additional dependencies
npm install -D tailwindcss@latest postcss autoprefixer @tailwindcss/postcss
npm install @headlessui/vue @heroicons/vue @vueuse/core
npm install uplot axios nats.ws
npm install -D @types/node vitest @vue/test-utils jsdom
npm install -D playwright @playwright/test

# Initialize TailwindCSS
npx tailwindcss init -p
```

### 2. Development Server

```bash
# Start dev server (http://localhost:5173)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm run test:unit
npm run test:e2e

# Lint and format
npm run lint
npm run format
```

### 3. Docker Build

```bash
# Build Docker image
docker build -t historian-web-ui:latest .

# Run container
docker run -p 3000:80 historian-web-ui:latest

# Build and run with docker-compose
cd /home/ahmet/historian/ops
docker-compose up --build web-ui
```

---

## 📁 Project Structure Reference

```
web-ui/
├── src/
│   ├── main.ts                    # Entry point
│   ├── App.vue                    # Root component
│   ├── router/index.ts            # Routes
│   ├── stores/                    # Pinia stores
│   │   ├── auth.ts
│   │   ├── sensors.ts
│   │   ├── realtime.ts
│   │   ├── alarms.ts
│   │   └── audit.ts
│   ├── composables/               # Reusable logic
│   │   ├── useNatsStream.ts
│   │   ├── useHistoryQuery.ts
│   │   └── useAuth.ts
│   ├── services/                  # API clients
│   │   ├── api.ts
│   │   ├── engine.ts
│   │   ├── auth.ts
│   │   ├── alarm.ts
│   │   └── audit.ts
│   ├── components/
│   │   ├── ui/                    # Base components
│   │   ├── charts/                # Chart components
│   │   ├── layout/                # Layout
│   │   └── common/                # Common
│   ├── views/                     # Pages
│   │   ├── DashboardView.vue
│   │   ├── SensorsView.vue
│   │   ├── AlarmsView.vue
│   │   ├── AuditView.vue
│   │   └── LoginView.vue
│   ├── types/                     # TypeScript types
│   └── utils/                     # Utilities
└── tests/
    ├── unit/
    └── e2e/
```

---

## 🔧 Configuration Files

### vite.config.ts

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/auth/, '/api/v1')
      },
      '/alarm': {
        target: 'http://localhost:8083',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/alarm/, '/api/v1')
      },
      '/audit': {
        target: 'http://localhost:8082',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/audit/, '/api/v1')
      },
      '/nats': {
        target: 'ws://localhost:8222',
        ws: true
      }
    }
  }
})
```

### tailwind.config.js

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eef2ff',
          100: '#e0e7ff',
          500: '#6366f1',
          600: '#4f46e5',
          700: '#4338ca',
          900: '#312e81',
        },
        alarm: {
          critical: '#dc2626',
          high: '#ea580c',
          medium: '#eab308',
          low: '#3b82f6',
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      }
    },
  },
  plugins: [],
}
```

### tsconfig.json

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

### Dockerfile

```dockerfile
# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build application
RUN npm run build

# Stage 2: Production
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
```

### nginx.conf

```nginx
server {
  listen 80;
  server_name _;
  root /usr/share/nginx/html;
  index index.html;

  # Gzip compression
  gzip on;
  gzip_vary on;
  gzip_min_length 1024;
  gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

  # SPA routing
  location / {
    try_files $uri $uri/ /index.html;
  }

  # API proxies
  location /api/ {
    proxy_pass http://engine:8081/api/v1/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
  }

  location /auth/ {
    proxy_pass http://auth:8080/api/v1/;
  }

  location /alarm/ {
    proxy_pass http://alarm:8083/api/v1/;
  }

  location /audit/ {
    proxy_pass http://audit:8082/api/v1/;
  }

  # NATS WebSocket
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

---

## 🎨 Code Templates

### Pinia Store Template

```typescript
// stores/example.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useExampleStore = defineStore('example', () => {
  // State
  const items = ref<Item[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const itemCount = computed(() => items.value.length)

  // Actions
  async function fetchItems() {
    loading.value = true
    error.value = null
    try {
      const data = await api.getItems()
      items.value = data
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  return {
    items,
    loading,
    error,
    itemCount,
    fetchItems
  }
})
```

### Composable Template

```typescript
// composables/useExample.ts
import { ref, onMounted, onUnmounted } from 'vue'

export function useExample() {
  const data = ref<Data | null>(null)
  const loading = ref(false)

  async function load() {
    loading.value = true
    try {
      data.value = await fetchData()
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    load()
  })

  return {
    data,
    loading,
    load
  }
}
```

### Component Template

```vue
<!-- components/Example.vue -->
<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  title: string
  count?: number
}

interface Emits {
  (e: 'update', value: number): void
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  count: 0
})

const emit = defineEmits<Emits>()

const localCount = ref(props.count)

const doubled = computed(() => localCount.value * 2)

function increment() {
  localCount.value++
  emit('update', localCount.value)
}
</script>

<template>
  <div class="p-4 bg-white rounded-lg shadow">
    <h2 class="text-xl font-bold">{{ title }}</h2>
    <p>Count: {{ localCount }}</p>
    <p>Doubled: {{ doubled }}</p>
    <button 
      @click="increment"
      class="px-4 py-2 bg-primary-600 text-white rounded hover:bg-primary-700"
    >
      Increment
    </button>
  </div>
</template>
```

### API Client Template

```typescript
// services/example.ts
import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// Request interceptor
client.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const exampleApi = {
  async getItems() {
    const { data } = await client.get('/items')
    return data
  },

  async createItem(item: CreateItemDto) {
    const { data } = await client.post('/items', item)
    return data
  },

  async updateItem(id: string, item: UpdateItemDto) {
    const { data } = await client.put(`/items/${id}`, item)
    return data
  },

  async deleteItem(id: string) {
    await client.delete(`/items/${id}`)
  }
}
```

---

## 🧪 Testing Templates

### Unit Test Template

```typescript
// tests/unit/Example.test.ts
import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Example from '@/components/Example.vue'

describe('Example', () => {
  it('renders properly', () => {
    const wrapper = mount(Example, {
      props: { title: 'Test' }
    })
    expect(wrapper.text()).toContain('Test')
  })

  it('increments count on button click', async () => {
    const wrapper = mount(Example, {
      props: { title: 'Test', count: 0 }
    })
    
    await wrapper.find('button').trigger('click')
    
    expect(wrapper.emitted('update')).toBeTruthy()
    expect(wrapper.emitted('update')?.[0]).toEqual([1])
  })
})
```

### Store Test Template

```typescript
// tests/unit/stores/example.test.ts
import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useExampleStore } from '@/stores/example'

describe('Example Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with empty items', () => {
    const store = useExampleStore()
    expect(store.items).toEqual([])
    expect(store.itemCount).toBe(0)
  })

  it('fetches items', async () => {
    const store = useExampleStore()
    await store.fetchItems()
    expect(store.items.length).toBeGreaterThan(0)
  })
})
```

---

## 🔗 API Endpoints Reference

### Engine Service (Port 8081)

```
GET  /api/v1/metadata
GET  /api/v1/export?tag_id={id}&start={ts}&end={ts}&format={csv|json}
```

### Auth Service (Port 8080)

```
POST /api/v1/login
POST /api/v1/re-auth
POST /api/v1/logout
GET  /api/v1/users
POST /api/v1/users
PUT  /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

### Alarm Service (Port 8083)

```
GET  /api/v1/alarms
POST /api/v1/alarms/{id}/acknowledge
POST /api/v1/alarms/{id}/shelve
DELETE /api/v1/alarms/{id}/unshelve
```

### Audit Service (Port 8082)

```
GET  /api/v1/audit
GET  /api/v1/audit/export
```

### NATS WebSocket (Port 8222)

```
ws://nats:8222

Subjects:
- data.{factory}.{line}.{machine}.{type}.{sensor_id}
- sys.alarm.{triggered|acknowledged|shelved}
- sys.audit.{action}
```

---

## 📚 Useful Resources

### Documentation
- [Vue 3 Docs](https://vuejs.org/)
- [Pinia Docs](https://pinia.vuejs.org/)
- [Vue Router Docs](https://router.vuejs.org/)
- [Vite Docs](https://vitejs.dev/)
- [TailwindCSS Docs](https://tailwindcss.com/)
- [Headless UI Docs](https://headlessui.com/)
- [VueUse Docs](https://vueuse.org/)
- [uPlot Docs](https://github.com/leeoniya/uPlot)

### Tools
- [Vue DevTools](https://devtools.vuejs.org/)
- [Vite Plugin Vue DevTools](https://github.com/webfansplz/vite-plugin-vue-devtools)

---

## 🐛 Common Issues & Solutions

### Issue: CORS errors in development

**Solution:** Use Vite proxy configuration (already in vite.config.ts)

### Issue: WebSocket connection fails

**Solution:** Check NATS is running and proxy configuration is correct

### Issue: TypeScript errors with uPlot

**Solution:** Add type definitions:
```typescript
declare module 'uplot' {
  export default class uPlot {
    constructor(opts: any, data: any, target: HTMLElement)
    setData(data: any): void
    destroy(): void
  }
}
```

### Issue: Pinia store not reactive

**Solution:** Use `ref()` or `reactive()` for state, not plain objects

---

## 🚢 Deployment Checklist

- [ ] Build passes without errors
- [ ] All tests pass
- [ ] ESLint passes
- [ ] Bundle size < 500KB (gzipped)
- [ ] Lighthouse score > 90
- [ ] Docker image builds successfully
- [ ] Environment variables configured
- [ ] API endpoints tested
- [ ] Authentication working
- [ ] Real-time updates working
- [ ] Mobile responsive
- [ ] Browser compatibility tested
- [ ] Security headers configured

---

## 📞 Support

For questions or issues:
1. Check this guide
2. Review main planning document: `docs/sprint-artifacts/modern-vue-web-interface.md`
3. Review architecture document: `docs/architecture-vue-frontend.md`
4. Contact PM (Ahmet)

---

**Last Updated:** 2025-12-11
**Version:** 1.0
