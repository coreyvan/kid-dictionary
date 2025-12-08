# Research: Kid Dictionary Frontend

**Feature**: 006-frontend-app
**Date**: 2025-12-07

## Technology Decisions

### 1. Frontend Framework: Vue 3 with Composition API

**Decision**: Vue 3.4+ with Composition API and `<script setup>` syntax

**Rationale**:
- Already specified in `docs/frontend.md` as the intended stack
- Composition API provides better TypeScript integration and code organization
- `<script setup>` reduces boilerplate and improves developer experience
- Strong ecosystem with Pinia, Vue Router, and Vite integration

**Alternatives Considered**:
- React: More ecosystem options but Vue already specified; switching adds no value
- Svelte: Smaller bundle but less mature ecosystem for PWA features
- Vue 2 Options API: Less composable, worse TypeScript support

### 2. Build Tool: Vite 5.x

**Decision**: Vite 5.x with vite-plugin-pwa

**Rationale**:
- Native ESM development server for fast hot module replacement
- Optimized production builds with Rollup
- First-class Vue support
- vite-plugin-pwa provides service worker generation and PWA manifest handling

**Alternatives Considered**:
- Webpack: Slower development experience, more configuration overhead
- Parcel: Less control over PWA configuration
- Turbopack: Still experimental for Vue

### 3. State Management: Pinia

**Decision**: Pinia 2.x

**Rationale**:
- Official Vue 3 state management solution (replaces Vuex)
- Simpler API without mutations
- Full TypeScript support with type inference
- DevTools integration for debugging
- Modular stores align with feature-based organization

**Alternatives Considered**:
- Vuex 4: More verbose, mutations add complexity without benefit
- Composables only: Insufficient for cross-component state sharing
- Zustand/Jotai: React-focused, not Vue ecosystem

### 4. API Client: Connect-Web

**Decision**: @connectrpc/connect-web with generated TypeScript clients

**Rationale**:
- Backend uses Connect-Go; Connect-Web provides protocol-compatible client
- Type-safe API calls generated from proto definitions
- Automatic serialization/deserialization
- Works with existing buf toolchain

**Alternatives Considered**:
- Fetch API + manual types: Error-prone, no type safety from protos
- Axios: Still requires manual type definitions
- gRPC-Web: Connect-Web is simpler and works without proxy

### 5. Styling: Tailwind CSS 3.x

**Decision**: Tailwind CSS 3.x with custom design tokens

**Rationale**:
- Utility-first approach speeds mobile-responsive development
- JIT compiler keeps bundle size small
- Easy to implement age bracket color theming
- No component library lock-in

**Alternatives Considered**:
- CSS Modules: More verbose for responsive designs
- Styled Components: Requires additional runtime
- Bootstrap/Material: Heavier, harder to customize for unique age bracket theming

### 6. Component Library: Headless UI

**Decision**: @headlessui/vue for accessible primitives only (dialogs, dropdowns)

**Rationale**:
- Unstyled, accessible components
- Integrates with Tailwind CSS
- Only use for complex accessible patterns (modals, menus)
- Build simple components (buttons, inputs) custom for control

**Alternatives Considered**:
- Vuetify/Quasar: Too opinionated, hard to match age bracket theming
- PrimeVue: Adds significant bundle size
- Fully custom: Reinventing accessibility patterns is error-prone

### 7. Testing: Vitest + Playwright

**Decision**: Vitest for unit/component tests, Playwright for E2E

**Rationale**:
- Vitest: Fast, Vite-native, Jest-compatible API
- Vue Test Utils integration for component testing
- Playwright: Cross-browser E2E, reliable, good mobile emulation
- Both support TypeScript natively

**Alternatives Considered**:
- Jest: Slower, requires additional Vite configuration
- Cypress: E2E only; component testing less mature than Vitest
- WebDriverIO: More complex setup

### 8. PWA Strategy: vite-plugin-pwa with Workbox

**Decision**: vite-plugin-pwa with generateSW strategy

**Rationale**:
- Automatic service worker generation
- Workbox provides caching strategies out of the box
- Network-first for API calls, cache-first for static assets
- Precaching for offline shell

**Alternatives Considered**:
- Manual service worker: More error-prone, maintenance overhead
- @vite-pwa/assets-generator: Will use for icon generation
- Offline-first: Too complex for initial release; network-first with cache fallback sufficient

### 9. Authentication Token Storage

**Decision**: localStorage for tokens, sessionStorage for anonymous conversations

**Rationale**:
- localStorage persists across browser sessions (user expectation for "stay logged in")
- sessionStorage for anonymous data aligns with clarification (cleared on browser close)
- Tokens have expiry; refresh flow handles security
- httpOnly cookies would require backend proxy changes

**Alternatives Considered**:
- Cookies: Require CORS/proxy configuration changes
- IndexedDB: Overkill for token storage
- Memory only: Loses state on refresh (poor UX)

### 10. Anonymous User Identification

**Decision**: Backend handles rate limiting via IP; frontend sends no identifier

**Rationale**:
- Simplest approach per Principle III (Simplicity)
- Backend already receives client IP from request headers
- No fingerprinting or tracking required on frontend
- Rate limit errors displayed as friendly registration prompts

**Alternatives Considered**:
- Client fingerprinting: Privacy concerns, complexity
- Anonymous session tokens: Requires backend changes
- Local storage counter: Easily bypassed, not reliable

## Integration Patterns

### Backend API Communication

1. **Connect-Web Client Setup**:
   - Generate TypeScript clients from proto files using `buf generate`
   - Configure base URL from environment variable
   - Add interceptor for auth token injection
   - Handle Connect error codes → user-friendly messages

2. **Authentication Flow**:
   - Login/Register → receive access_token + refresh_token
   - Store tokens in localStorage
   - Inject access_token in Authorization header
   - Auto-refresh when token near expiry (check on each request)
   - On refresh failure → redirect to login

3. **Anonymous Mode**:
   - No auth header sent
   - Conversations stored in sessionStorage only
   - On registration → clear sessionStorage (no migration)

### Offline Support

1. **Service Worker Caching**:
   - Precache: App shell (HTML, JS, CSS, fonts)
   - Runtime cache: API responses for conversation list (authenticated only)
   - Network-first strategy for API calls
   - Stale-while-revalidate for static assets

2. **Offline Detection**:
   - `navigator.onLine` + online/offline events
   - Display offline indicator when offline
   - Queue failed messages for retry (stretch goal)

## Open Questions Resolved

| Question | Resolution |
|----------|------------|
| Which Vue version? | Vue 3.4+ with Composition API |
| How to consume Connect-Go API? | Connect-Web with buf-generated clients |
| Component library? | Headless UI for accessibility primitives only |
| How to identify anonymous users? | Backend uses IP; frontend sends nothing |
| PWA caching strategy? | Network-first for API, cache-first for static |

## Dependencies Summary

```json
{
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "@connectrpc/connect": "^1.0.0",
    "@connectrpc/connect-web": "^1.0.0",
    "@headlessui/vue": "^1.7.0"
  },
  "devDependencies": {
    "vite": "^5.0.0",
    "vite-plugin-pwa": "^0.17.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.3.0",
    "vitest": "^1.0.0",
    "@vue/test-utils": "^2.4.0",
    "playwright": "^1.40.0",
    "@bufbuild/buf": "^1.28.0",
    "@bufbuild/protoc-gen-es": "^1.6.0",
    "@connectrpc/protoc-gen-connect-es": "^1.2.0"
  }
}
```