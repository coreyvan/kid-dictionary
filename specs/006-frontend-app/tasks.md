# Tasks: Kid Dictionary Frontend Application

**Input**: Design documents from `/specs/006-frontend-app/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Not explicitly requested in specification. Test tasks omitted.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Frontend**: `frontend/src/`, `frontend/tests/`
- All paths relative to repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and Vue 3 + Vite + Tailwind configuration

- [X] T001 Create frontend directory and initialize Vue 3 + TypeScript project with Vite in frontend/
- [X] T002 Install dependencies: vue-router, pinia, @connectrpc/connect, @connectrpc/connect-web, @headlessui/vue per frontend/package.json
- [X] T003 [P] Install dev dependencies: tailwindcss, postcss, autoprefixer, vite-plugin-pwa, vitest per frontend/package.json
- [X] T004 [P] Configure Tailwind CSS with age bracket colors in frontend/tailwind.config.js
- [X] T005 [P] Configure Vite with PWA plugin and path aliases in frontend/vite.config.ts
- [X] T006 [P] Create TypeScript config in frontend/tsconfig.json
- [X] T007 [P] Create environment files frontend/.env.development and frontend/.env.production
- [X] T008 Create buf.gen.yaml for frontend TypeScript client generation in frontend/buf.gen.yaml
- [X] T009 Generate Connect-Web clients from proto files into frontend/src/gen/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T010 Create TypeScript types for all entities in frontend/src/types/index.ts (User, Conversation, Message, AgeBracket, ContentTier enums)
- [X] T011 Create API client base setup with Connect transport in frontend/src/services/api/client.ts
- [X] T012 [P] Create error handling utilities with user-friendly messages in frontend/src/services/api/errors.ts
- [X] T013 [P] Create auth store with token management in frontend/src/stores/auth.ts
- [X] T014 Create auth interceptor for token injection in frontend/src/services/api/interceptors.ts
- [X] T015 [P] Create Vue Router configuration with route guards in frontend/src/router/index.ts
- [X] T016 [P] Create LoadingSpinner component in frontend/src/components/common/LoadingSpinner.vue
- [X] T017 [P] Create ErrorMessage component in frontend/src/components/common/ErrorMessage.vue
- [X] T018 [P] Create base CSS with Tailwind directives in frontend/src/style.css
- [X] T019 Create App.vue with router-view and global layout in frontend/src/App.vue
- [X] T020 Create main.ts entry point with Pinia and Router setup in frontend/src/main.ts

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Ask a Question About Any Topic (Priority: P1) 🎯 MVP

**Goal**: Anonymous user can open app, select age bracket, ask question, receive age-appropriate explanation

**Independent Test**: Open app → Select "Growing Minds" → Type question → Receive AI response with content tier indicator

### Implementation for User Story 1

- [X] T021 [P] [US1] Create AgeBracketSelector component with color-coded options in frontend/src/components/chat/AgeBracketSelector.vue
- [X] T022 [P] [US1] Create ChatMessage component with role styling in frontend/src/components/chat/ChatMessage.vue
- [X] T023 [P] [US1] Create ContentTierIndicator component for sensitive content labels in frontend/src/components/chat/ContentTierIndicator.vue
- [X] T024 [P] [US1] Create ChatInput component with send button (44px touch target) in frontend/src/components/chat/ChatInput.vue
- [X] T025 [US1] Create messages store for message state management in frontend/src/stores/messages.ts
- [X] T026 [US1] Create message API service wrapper in frontend/src/services/api/messages.ts
- [X] T027 [US1] Create anonymous conversation storage utilities in frontend/src/services/storage/anonymous.ts
- [X] T028 [US1] Create ChatView with message list, input, and age bracket selector in frontend/src/views/ChatView.vue
- [X] T029 [US1] Create HomeView as landing page with "Start Chat" CTA in frontend/src/views/HomeView.vue
- [X] T030 [US1] Add routes for home (/) and chat (/chat) in frontend/src/router/index.ts
- [X] T031 [US1] Implement optimistic message sending with loading state in frontend/src/stores/messages.ts
- [X] T032 [US1] Add message retry on failure in frontend/src/components/chat/ChatMessage.vue

**Checkpoint**: User Story 1 complete - Anonymous users can ask questions and receive responses

---

## Phase 4: User Story 2 - Register and Log In (Priority: P2)

**Goal**: Users can register with email/password, log in, and access authenticated features

**Independent Test**: Register with email → Logout → Login with same credentials → See conversation list

### Implementation for User Story 2

- [X] T033 [P] [US2] Create LoginForm component with validation in frontend/src/components/auth/LoginForm.vue
- [X] T034 [P] [US2] Create RegisterForm component with age bracket selection in frontend/src/components/auth/RegisterForm.vue
- [X] T035 [US2] Create auth API service for register, login, refresh in frontend/src/services/api/auth.ts
- [X] T036 [US2] Implement token storage and refresh logic in frontend/src/stores/auth.ts
- [X] T037 [US2] Create LoginView page in frontend/src/views/LoginView.vue
- [X] T038 [US2] Create RegisterView page in frontend/src/views/RegisterView.vue
- [X] T039 [US2] Add routes for login (/login) and register (/register) in frontend/src/router/index.ts
- [X] T040 [US2] Add auth navigation guard for protected routes in frontend/src/router/index.ts
- [X] T041 [US2] Add login/register/logout buttons to App.vue header in frontend/src/App.vue
- [X] T042 [US2] Handle token refresh interceptor triggering on API calls in frontend/src/services/api/interceptors.ts

**Checkpoint**: User Story 2 complete - Users can register, login, and logout

---

## Phase 5: User Story 3 - View and Manage Conversation History (Priority: P3)

**Goal**: Authenticated users can view, resume, and delete past conversations

**Independent Test**: Login → View conversation list → Open previous conversation → Delete a conversation

### Implementation for User Story 3

- [X] T043 [P] [US3] Create ConversationCard component with title and date in frontend/src/components/conversations/ConversationCard.vue
- [X] T044 [P] [US3] Create EmptyState component for no conversations in frontend/src/components/common/EmptyState.vue
- [X] T045 [US3] Create conversations store for list management in frontend/src/stores/conversations.ts
- [X] T046 [US3] Create conversation API service wrapper in frontend/src/services/api/conversations.ts
- [X] T047 [US3] Create ConversationList component with pagination in frontend/src/components/conversations/ConversationList.vue
- [X] T048 [US3] Create ConversationsView page in frontend/src/views/ConversationsView.vue
- [X] T049 [US3] Add route for conversations (/conversations) with auth guard in frontend/src/router/index.ts
- [X] T050 [US3] Update ChatView to load existing conversation messages in frontend/src/views/ChatView.vue
- [X] T051 [US3] Add delete confirmation dialog using HeadlessUI in frontend/src/components/conversations/DeleteConfirmDialog.vue
- [X] T052 [US3] Implement infinite scroll pagination in ConversationList in frontend/src/components/conversations/ConversationList.vue

**Checkpoint**: User Story 3 complete - Authenticated users can manage conversation history

---

## Phase 6: User Story 4 - Change Age Bracket During Conversation (Priority: P4)

**Goal**: Users can change age bracket mid-conversation and see visual feedback

**Independent Test**: Start conversation with "Little Ones" → Change to "Pre-Teens" → See color change → Ask question → Response matches new bracket

### Implementation for User Story 4

- [X] T053 [US4] Update AgeBracketSelector to emit change events in frontend/src/components/chat/AgeBracketSelector.vue
- [X] T054 [US4] Add updateConversation call on age bracket change in frontend/src/views/ChatView.vue
- [X] T055 [US4] Add visual transition for bracket color change in frontend/src/components/chat/AgeBracketSelector.vue
- [X] T056 [US4] Show confirmation toast/message on bracket change in frontend/src/views/ChatView.vue

**Checkpoint**: User Story 4 complete - Age bracket can be changed mid-conversation

---

## Phase 7: User Story 5 - Access App on Mobile Device (Priority: P5)

**Goal**: App is fully functional on mobile with touch-friendly UI

**Independent Test**: Open on mobile device → All buttons 44px+ → Input at bottom → Clear loading states

### Implementation for User Story 5

- [X] T057 [P] [US5] Audit and fix all touch targets to minimum 44px in all components
- [X] T058 [P] [US5] Ensure ChatInput stays fixed at bottom of viewport in frontend/src/views/ChatView.vue
- [X] T059 [P] [US5] Add responsive breakpoints to ConversationList in frontend/src/components/conversations/ConversationList.vue
- [X] T060 [US5] Create useOffline composable for network detection in frontend/src/composables/useOffline.ts
- [X] T061 [US5] Create OfflineIndicator component in frontend/src/components/common/OfflineIndicator.vue
- [X] T062 [US5] Add OfflineIndicator to App.vue layout in frontend/src/App.vue
- [X] T063 [US5] Test and optimize loading states for slow connections in all views

**Checkpoint**: User Story 5 complete - App works well on mobile devices

---

## Phase 8: User Story 6 - Install as Progressive Web App (Priority: P6)

**Goal**: App can be installed to home screen and works offline (viewing cached data)

**Independent Test**: Install PWA → Launch from home screen → View cached conversations offline

### Implementation for User Story 6

- [X] T064 [P] [US6] Create PWA icons (192x192, 512x512) in frontend/public/
- [X] T065 [P] [US6] Configure PWA manifest in vite.config.ts in frontend/vite.config.ts
- [X] T066 [US6] Configure Workbox runtime caching for conversation list API in frontend/vite.config.ts
- [X] T067 [US6] Add install prompt handling in App.vue in frontend/src/App.vue
- [X] T068 [US6] Create service worker registration in main.ts in frontend/src/main.ts
- [X] T069 [US6] Test offline functionality with cached conversation list

**Checkpoint**: User Story 6 complete - App installable and works offline

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements affecting multiple user stories

- [X] T070 [P] Add rate limit error handling with registration prompt in frontend/src/services/api/errors.ts
- [X] T071 [P] Create SettingsView for default age bracket and account info in frontend/src/views/SettingsView.vue
- [X] T072 Add route for settings (/settings) with auth guard in frontend/src/router/index.ts
- [ ] T073 [P] Update README.md with frontend setup instructions in README.md
- [X] T074 [P] Add ESLint + Prettier configuration in frontend/
- [ ] T075 Run quickstart.md validation - verify all setup steps work
- [ ] T076 Manual testing of all user stories end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-8)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → P4 → P5 → P6)
- **Polish (Phase 9)**: Depends on all desired user stories being complete

### User Story Dependencies

| User Story | Depends On | Can Run Parallel With |
|------------|------------|----------------------|
| US1 (P1) | Foundational | None (MVP) |
| US2 (P2) | Foundational | US1 (different components) |
| US3 (P3) | US2 (auth required) | - |
| US4 (P4) | US1 (needs ChatView) | US2, US3 |
| US5 (P5) | US1 (needs components) | US2, US3, US4 |
| US6 (P6) | Foundational | Any other story |

### Within Each User Story

- Components can often be built in parallel [P]
- Stores depend on types and API services
- Views depend on components and stores
- Routes depend on views

### Parallel Opportunities

**Phase 1 (Setup)**:
```
T003, T004, T005, T006, T007 can run in parallel
```

**Phase 2 (Foundational)**:
```
T012, T013, T015, T016, T017, T018 can run in parallel
```

**Phase 3 (US1 - MVP)**:
```
T021, T022, T023, T024 can run in parallel (all components)
```

**Cross-Story Parallel** (with multiple developers):
```
Developer A: US1 (T021-T032)
Developer B: US2 (T033-T042) - after Foundational
Developer C: US6 (T064-T069) - PWA setup
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (~9 tasks)
2. Complete Phase 2: Foundational (~11 tasks)
3. Complete Phase 3: User Story 1 (~12 tasks)
4. **STOP and VALIDATE**: Anonymous chat works end-to-end
5. Deploy/demo if ready (32 tasks total for MVP)

### Incremental Delivery

| Milestone | Tasks | Cumulative | Value Delivered |
|-----------|-------|------------|-----------------|
| Setup | T001-T009 | 9 | Project structure ready |
| Foundation | T010-T020 | 20 | Core infrastructure |
| US1 (MVP) | T021-T032 | 32 | Anonymous chat works |
| US2 | T033-T042 | 42 | Auth + accounts |
| US3 | T043-T052 | 52 | Conversation history |
| US4 | T053-T056 | 56 | Age bracket switching |
| US5 | T057-T063 | 63 | Mobile optimized |
| US6 | T064-T069 | 69 | PWA installable |
| Polish | T070-T076 | 76 | Production ready |

---

## Summary

| Phase | Tasks | Parallel Tasks |
|-------|-------|----------------|
| Setup | 9 | 5 |
| Foundational | 11 | 6 |
| US1 (P1) | 12 | 4 |
| US2 (P2) | 10 | 2 |
| US3 (P3) | 10 | 2 |
| US4 (P4) | 4 | 0 |
| US5 (P5) | 7 | 3 |
| US6 (P6) | 6 | 2 |
| Polish | 7 | 4 |
| **Total** | **76** | **28** |

**MVP Scope**: 32 tasks (Setup + Foundational + US1)
**Full Scope**: 76 tasks