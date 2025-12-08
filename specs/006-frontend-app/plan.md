# Implementation Plan: Kid Dictionary Frontend Application

**Branch**: `006-frontend-app` | **Date**: 2025-12-07 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-frontend-app/spec.md`

## Summary

Build a mobile-first Progressive Web App frontend for Kid Dictionary that enables parents to get age-appropriate explanations for their children's questions. The app supports anonymous usage with session-based persistence and optional registration for conversation history sync. Uses Vue 3 with Composition API, Vite for build, Tailwind CSS for styling, and Pinia for state management.

## Technical Context

**Language/Version**: TypeScript 5.x with Vue 3.4+
**Primary Dependencies**: Vue 3 (Composition API), Vite 5.x, Tailwind CSS 3.x, Pinia 2.x, Vue Router 4.x, vite-plugin-pwa
**Storage**: Browser sessionStorage (anonymous), localStorage (auth tokens), Backend API (authenticated data)
**Testing**: Vitest for unit tests, Playwright for E2E tests
**Target Platform**: Modern browsers (Chrome, Safari, Firefox, Edge), iOS Safari, Android Chrome
**Project Type**: Web application (frontend only - consumes existing Go backend API)
**Performance Goals**: <3s initial load on 4G, <100ms UI interactions, <15s AI response display
**Constraints**: Mobile-first design, 44px minimum touch targets, offline conversation list caching, PWA installable
**Scale/Scope**: Single-page application with 5 main views, ~15 components

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | PASS | Backend API already defined in proto files; frontend consumes via Connect-Web client |
| II. Test-Alongside Development | PASS | Plan includes Vitest unit tests and Playwright E2E tests alongside implementation |
| III. Simplicity & YAGNI | PASS | Using established Vue ecosystem (no custom abstractions); minimal dependencies |
| IV. Observability | PASS | Will integrate structured logging for API calls, error tracking |
| V. Content Safety | PASS | Frontend displays content tier indicators; backend handles classification |

**Additional Constraints Check**:
- Technology Stack: Frontend extends existing Go backend (compatible)
- Security: JWT tokens stored in localStorage, auto-refresh before expiry
- Documentation: README will be updated with frontend setup instructions

## Project Structure

### Documentation (this feature)

```text
specs/006-frontend-app/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (API client specs)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── components/
│   │   ├── chat/
│   │   │   ├── ChatMessage.vue
│   │   │   ├── ChatInput.vue
│   │   │   ├── AgeBracketSelector.vue
│   │   │   └── ContentTierIndicator.vue
│   │   ├── auth/
│   │   │   ├── LoginForm.vue
│   │   │   └── RegisterForm.vue
│   │   ├── conversations/
│   │   │   ├── ConversationList.vue
│   │   │   └── ConversationCard.vue
│   │   └── common/
│   │       ├── LoadingSpinner.vue
│   │       ├── ErrorMessage.vue
│   │       ├── OfflineIndicator.vue
│   │       └── EmptyState.vue
│   ├── views/
│   │   ├── HomeView.vue
│   │   ├── LoginView.vue
│   │   ├── RegisterView.vue
│   │   ├── ConversationsView.vue
│   │   ├── ChatView.vue
│   │   └── SettingsView.vue
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── conversations.ts
│   │   └── messages.ts
│   ├── composables/
│   │   ├── useApi.ts
│   │   ├── useAuth.ts
│   │   └── useOffline.ts
│   ├── services/
│   │   └── api/
│   │       ├── client.ts
│   │       ├── auth.ts
│   │       ├── conversations.ts
│   │       └── messages.ts
│   ├── router/
│   │   └── index.ts
│   ├── types/
│   │   └── index.ts
│   ├── App.vue
│   └── main.ts
├── public/
│   └── manifest.json
├── tests/
│   ├── unit/
│   └── e2e/
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── package.json
```

**Structure Decision**: Standalone frontend directory at repository root. Vue 3 with TypeScript, following Vue community conventions. API services layer wraps Connect-Web client for backend communication.

## Complexity Tracking

> No constitution violations requiring justification.

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| State Management | Pinia (not Vuex) | Simpler API, better TypeScript support, Vue 3 recommended |
| API Client | Connect-Web | Matches backend Connect-Go protocol, type-safe from protos |
| Styling | Tailwind CSS | Utility-first, mobile-responsive, no custom CSS framework |
| Testing | Vitest + Playwright | Fast unit tests, reliable E2E, Vue ecosystem standard |