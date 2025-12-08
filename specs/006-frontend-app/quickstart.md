# Quickstart: Kid Dictionary Frontend

**Feature**: 006-frontend-app
**Date**: 2025-12-07

## Prerequisites

- Node.js 20+ (LTS)
- pnpm 8+ (or npm/yarn)
- Backend API running locally on port 8080

## Setup

### 1. Create Frontend Directory

```bash
# From repository root
mkdir frontend
cd frontend
```

### 2. Initialize Project

```bash
pnpm create vite@latest . --template vue-ts
```

### 3. Install Dependencies

```bash
pnpm add vue-router@4 pinia @connectrpc/connect @connectrpc/connect-web @headlessui/vue

pnpm add -D tailwindcss postcss autoprefixer vite-plugin-pwa vitest @vue/test-utils playwright @bufbuild/buf @bufbuild/protoc-gen-es @connectrpc/protoc-gen-connect-es
```

### 4. Configure Tailwind CSS

```bash
npx tailwindcss init -p
```

Update `tailwind.config.js`:
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
        'bracket-little': '#FFB347',    // Warm orange
        'bracket-growing': '#77DD77',   // Bright green
        'bracket-preteen': '#6CA0DC',   // Confident blue
      },
      minHeight: {
        'touch': '44px',
      },
      minWidth: {
        'touch': '44px',
      },
    },
  },
  plugins: [],
}
```

### 5. Generate API Client

Add to `frontend/buf.gen.yaml`:
```yaml
version: v1
plugins:
  - plugin: es
    out: src/gen
    opt: target=ts
  - plugin: connect-es
    out: src/gen
    opt: target=ts
```

Generate clients:
```bash
cd ..  # back to repo root
buf generate --template frontend/buf.gen.yaml proto
```

### 6. Configure Vite

Update `vite.config.ts`:
```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [
    vue(),
    VitePWA({
      registerType: 'autoUpdate',
      manifest: {
        name: 'Kid Dictionary',
        short_name: 'KidDict',
        description: 'Age-appropriate explanations for curious kids',
        theme_color: '#6CA0DC',
        icons: [
          { src: '/icon-192.png', sizes: '192x192', type: 'image/png' },
          { src: '/icon-512.png', sizes: '512x512', type: 'image/png' },
        ],
      },
      workbox: {
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/api\.kiddictionary\.example\.com\/.*/i,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              expiration: { maxEntries: 50, maxAgeSeconds: 300 },
            },
          },
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
```

### 7. Create Environment Files

`.env.development`:
```env
VITE_API_BASE_URL=http://localhost:8080
```

`.env.production`:
```env
VITE_API_BASE_URL=https://api.kiddictionary.example.com
```

## Development

### Start Backend

```bash
# In repository root
docker compose up -d
go run cmd/server/main.go
```

### Start Frontend

```bash
cd frontend
pnpm dev
```

Open http://localhost:5173

## Testing

### Unit Tests

```bash
pnpm test
```

### E2E Tests

```bash
pnpm exec playwright install
pnpm exec playwright test
```

## Build

```bash
pnpm build
```

Output in `frontend/dist/`

## Project Structure After Setup

```
frontend/
├── src/
│   ├── components/     # Vue components
│   ├── views/          # Route views
│   ├── stores/         # Pinia stores
│   ├── composables/    # Vue composables
│   ├── services/       # API client wrappers
│   ├── gen/            # Generated protobuf clients
│   ├── router/         # Vue Router config
│   ├── types/          # TypeScript types
│   ├── App.vue
│   └── main.ts
├── public/
│   ├── icon-192.png
│   ├── icon-512.png
│   └── manifest.json
├── tests/
│   ├── unit/
│   └── e2e/
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
├── buf.gen.yaml
└── package.json
```

## Common Tasks

| Task | Command |
|------|---------|
| Start dev server | `pnpm dev` |
| Run unit tests | `pnpm test` |
| Run E2E tests | `pnpm exec playwright test` |
| Build for production | `pnpm build` |
| Preview production build | `pnpm preview` |
| Regenerate API client | `buf generate --template frontend/buf.gen.yaml proto` |
| Lint | `pnpm lint` |