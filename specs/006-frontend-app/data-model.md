# Data Model: Kid Dictionary Frontend

**Feature**: 006-frontend-app
**Date**: 2025-12-07

## Overview

The frontend data model mirrors the backend proto definitions but adds client-side state for UI concerns (loading states, form validation, offline status). All types are derived from the protobuf definitions to ensure type safety.

## Core Entities

### User

**Source**: `kiddictionary.v1.User` (proto)

```typescript
interface User {
  id: string;
  email: string;
  defaultAgeBracket: AgeBracket;
  createdAt: Date;
  updatedAt: Date;
}
```

**Frontend State Extensions**:
```typescript
interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}
```

**Validation Rules**:
- `email`: Valid email format (RFC 5322)
- `password`: Minimum 8 characters (registration only, not stored)

**State Transitions**:
- Anonymous → Authenticated (via Register or Login)
- Authenticated → Anonymous (via Logout or token expiry)

---

### Conversation

**Source**: `kiddictionary.v1.Conversation` (proto)

```typescript
interface Conversation {
  id: string;
  userId: string;
  title: string;
  ageBracket: AgeBracket;
  createdAt: Date;
  updatedAt: Date;
}
```

**Frontend State Extensions**:
```typescript
interface ConversationsState {
  conversations: Conversation[];
  currentConversation: Conversation | null;
  isLoading: boolean;
  error: string | null;
  nextPageToken: string | null;
  hasMore: boolean;
}
```

**Anonymous Conversation** (sessionStorage):
```typescript
interface AnonymousConversation {
  id: string; // client-generated UUID
  title: string;
  ageBracket: AgeBracket;
  messages: Message[];
  createdAt: Date;
}
```

**Validation Rules**:
- `title`: Non-empty string, max 200 characters
- `ageBracket`: Must be valid enum value

**State Transitions**:
- Created → Active (messages added)
- Active → Deleted (user action, with confirmation)

---

### Message

**Source**: `kiddictionary.v1.Message` (proto)

```typescript
interface Message {
  id: string;
  conversationId: string;
  role: MessageRole;
  content: string;
  createdAt: Date;
  contentTier: ContentTier; // for assistant messages
}
```

**Frontend State Extensions**:
```typescript
interface MessagesState {
  messages: Message[];
  pendingMessage: PendingMessage | null;
  isLoading: boolean;
  error: string | null;
}

interface PendingMessage {
  id: string; // temporary client ID
  content: string;
  status: 'sending' | 'error';
  error?: string;
}
```

**Validation Rules**:
- `content`: Non-empty string, max 5000 characters

**State Transitions**:
- Pending (optimistic) → Sent (API success)
- Pending → Error (API failure, allow retry)

---

### AgeBracket (Enum)

**Source**: `kiddictionary.v1.AgeBracket` (proto)

```typescript
enum AgeBracket {
  UNSPECIFIED = 0,
  LITTLE_ONES = 1,    // 0-5 years
  GROWING_MINDS = 2,  // 5-10 years
  PRE_TEENS = 3,      // 10+ years
}
```

**Display Mapping**:
```typescript
const AgeBracketDisplay: Record<AgeBracket, { label: string; color: string; description: string }> = {
  [AgeBracket.LITTLE_ONES]: {
    label: 'Little Ones',
    color: 'bracket-little', // Tailwind class
    description: '0-5 years',
  },
  [AgeBracket.GROWING_MINDS]: {
    label: 'Growing Minds',
    color: 'bracket-growing',
    description: '5-10 years',
  },
  [AgeBracket.PRE_TEENS]: {
    label: 'Pre-Teens',
    color: 'bracket-preteen',
    description: '10+ years',
  },
};
```

---

### MessageRole (Enum)

**Source**: `kiddictionary.v1.MessageRole` (proto)

```typescript
enum MessageRole {
  UNSPECIFIED = 0,
  USER = 1,
  ASSISTANT = 2,
}
```

---

### ContentTier (Enum)

**Source**: `kiddictionary.v1.ContentTier` (proto)

```typescript
enum ContentTier {
  UNSPECIFIED = 0,
  NORMAL = 1,       // Standard topics
  SENSITIVE = 2,    // Death, divorce, etc.
  CONTEXTUAL = 3,   // Religion, politics
  REDIRECT = 4,     // Harmful content - declined
}
```

**Display Mapping**:
```typescript
const ContentTierIndicator: Record<ContentTier, { icon: string; label: string } | null> = {
  [ContentTier.NORMAL]: null, // No indicator
  [ContentTier.SENSITIVE]: { icon: 'heart', label: 'Sensitive topic' },
  [ContentTier.CONTEXTUAL]: { icon: 'info', label: 'Context-dependent' },
  [ContentTier.REDIRECT]: { icon: 'shield', label: 'Alternative suggested' },
};
```

---

## Client-Side Only State

### UI State

```typescript
interface UIState {
  isOffline: boolean;
  isMobile: boolean;
  sidebarOpen: boolean;
  theme: 'light' | 'dark' | 'system';
}
```

### Form State

```typescript
interface LoginForm {
  email: string;
  password: string;
  errors: { email?: string; password?: string };
  isSubmitting: boolean;
}

interface RegisterForm {
  email: string;
  password: string;
  confirmPassword: string;
  defaultAgeBracket: AgeBracket;
  errors: { email?: string; password?: string; confirmPassword?: string };
  isSubmitting: boolean;
}

interface ChatInputForm {
  content: string;
  isSubmitting: boolean;
}
```

---

## Storage Schema

### localStorage

```typescript
interface LocalStorageSchema {
  'auth:accessToken': string;
  'auth:refreshToken': string;
  'auth:user': string; // JSON serialized User
  'settings:theme': 'light' | 'dark' | 'system';
}
```

### sessionStorage

```typescript
interface SessionStorageSchema {
  'anonymous:conversations': string; // JSON serialized AnonymousConversation[]
  'anonymous:currentId': string; // Current conversation ID
}
```

---

## Relationships

```
User (1) ──────── (N) Conversation
                        │
Conversation (1) ─── (N) Message

AnonymousConversation (1) ─── (N) Message (embedded)
```

---

## Data Flow

1. **Anonymous User**:
   - All data in sessionStorage
   - No backend persistence
   - Cleared on browser close

2. **Authenticated User**:
   - Auth tokens in localStorage
   - Conversations/messages fetched from backend
   - Conversation list cached for offline (service worker)

3. **Token Refresh**:
   - Access token checked before each API call
   - If near expiry, refresh token used to obtain new tokens
   - On refresh failure, user redirected to login