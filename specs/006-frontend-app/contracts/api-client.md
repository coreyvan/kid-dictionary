# API Client Contract: Kid Dictionary Frontend

**Feature**: 006-frontend-app
**Date**: 2025-12-07

## Overview

The frontend consumes the Kid Dictionary backend API via Connect-Web, which provides a type-safe client generated from the protobuf definitions. This document specifies the client configuration and usage patterns.

## Backend Services

The backend exposes three Connect-Go services (defined in `proto/kiddictionary/v1/`):

### AuthService

**Endpoints**:

| Method | Request | Response | Auth Required |
|--------|---------|----------|---------------|
| `Register` | `RegisterRequest` | `RegisterResponse` | No |
| `Login` | `LoginRequest` | `LoginResponse` | No |
| `RefreshToken` | `RefreshTokenRequest` | `RefreshTokenResponse` | No (uses refresh token) |

**Usage**:
```typescript
// Register
const response = await authClient.register({
  email: 'user@example.com',
  password: 'securepassword',
  defaultAgeBracket: AgeBracket.GROWING_MINDS,
});
// Returns: { user, accessToken, refreshToken }

// Login
const response = await authClient.login({
  email: 'user@example.com',
  password: 'securepassword',
});
// Returns: { user, accessToken, refreshToken }

// Refresh Token
const response = await authClient.refreshToken({
  refreshToken: storedRefreshToken,
});
// Returns: { accessToken, refreshToken }
```

---

### ConversationService

**Endpoints**:

| Method | Request | Response | Auth Required |
|--------|---------|----------|---------------|
| `CreateConversation` | `CreateConversationRequest` | `CreateConversationResponse` | Yes |
| `GetConversation` | `GetConversationRequest` | `GetConversationResponse` | Yes |
| `ListConversations` | `ListConversationsRequest` | `ListConversationsResponse` | Yes |
| `UpdateConversation` | `UpdateConversationRequest` | `UpdateConversationResponse` | Yes |
| `DeleteConversation` | `DeleteConversationRequest` | `DeleteConversationResponse` | Yes |

**Usage**:
```typescript
// Create Conversation
const response = await conversationClient.createConversation({
  title: 'Why is the sky blue?',
  ageBracket: AgeBracket.GROWING_MINDS,
});
// Returns: { conversation }

// List Conversations (paginated)
const response = await conversationClient.listConversations({
  pageSize: 20,
  pageToken: nextPageToken, // optional
});
// Returns: { conversations, nextPageToken }

// Get Conversation with Messages
const response = await conversationClient.getConversation({
  id: conversationId,
});
// Returns: { conversation, messages }

// Update Conversation
const response = await conversationClient.updateConversation({
  id: conversationId,
  title: 'New title', // optional
  ageBracket: AgeBracket.PRE_TEENS, // optional
});
// Returns: { conversation }

// Delete Conversation
await conversationClient.deleteConversation({
  id: conversationId,
});
// Returns: {}
```

---

### MessageService

**Endpoints**:

| Method | Request | Response | Auth Required |
|--------|---------|----------|---------------|
| `SendMessage` | `SendMessageRequest` | `SendMessageResponse` | Yes* |

*Anonymous users can send messages but conversations are not persisted server-side.

**Usage**:
```typescript
// Send Message
const response = await messageClient.sendMessage({
  conversationId: conversationId,
  content: 'Why do dogs bark?',
});
// Returns: { userMessage, assistantMessage }
```

---

## Client Configuration

### Base Setup

```typescript
// services/api/client.ts
import { createConnectTransport } from '@connectrpc/connect-web';
import { createPromiseClient } from '@connectrpc/connect';
import { AuthService } from '@/gen/kiddictionary/v1/auth_connect';
import { ConversationService } from '@/gen/kiddictionary/v1/conversation_connect';
import { MessageService } from '@/gen/kiddictionary/v1/message_connect';

const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
  interceptors: [authInterceptor],
});

export const authClient = createPromiseClient(AuthService, transport);
export const conversationClient = createPromiseClient(ConversationService, transport);
export const messageClient = createPromiseClient(MessageService, transport);
```

### Auth Interceptor

```typescript
// services/api/interceptors.ts
import type { Interceptor } from '@connectrpc/connect';
import { useAuthStore } from '@/stores/auth';

export const authInterceptor: Interceptor = (next) => async (req) => {
  const authStore = useAuthStore();

  // Check if token needs refresh
  if (authStore.shouldRefreshToken()) {
    await authStore.refreshToken();
  }

  // Add auth header if authenticated
  if (authStore.accessToken) {
    req.header.set('Authorization', `Bearer ${authStore.accessToken}`);
  }

  return next(req);
};
```

---

## Error Handling

### Connect Error Codes → User Messages

```typescript
// services/api/errors.ts
import { ConnectError, Code } from '@connectrpc/connect';

export function getErrorMessage(error: unknown): string {
  if (error instanceof ConnectError) {
    switch (error.code) {
      case Code.Unauthenticated:
        return 'Please log in to continue.';
      case Code.PermissionDenied:
        return 'You do not have permission to perform this action.';
      case Code.NotFound:
        return 'The requested item was not found.';
      case Code.AlreadyExists:
        return 'This email is already registered.';
      case Code.InvalidArgument:
        return error.message || 'Invalid input provided.';
      case Code.ResourceExhausted:
        return 'Rate limit reached. Please register for unlimited access.';
      case Code.Unavailable:
        return 'Service temporarily unavailable. Please try again.';
      default:
        return 'An unexpected error occurred. Please try again.';
    }
  }
  return 'An unexpected error occurred.';
}
```

---

## Anonymous Mode API

For anonymous users, the frontend simulates API responses locally:

```typescript
// services/api/anonymous.ts
import { v4 as uuidv4 } from 'uuid';

export async function sendAnonymousMessage(
  conversationId: string,
  content: string,
  ageBracket: AgeBracket
): Promise<SendMessageResponse> {
  // Call backend without auth (rate-limited by IP)
  const response = await messageClient.sendMessage({
    conversationId: conversationId,
    content: content,
  });

  // Store in sessionStorage
  saveAnonymousMessage(conversationId, response.userMessage);
  saveAnonymousMessage(conversationId, response.assistantMessage);

  return response;
}
```

---

## Code Generation

### buf.gen.yaml Configuration

```yaml
version: v1
plugins:
  - plugin: es
    out: frontend/src/gen
    opt: target=ts
  - plugin: connect-es
    out: frontend/src/gen
    opt: target=ts
```

### Generated Files

```
frontend/src/gen/
├── kiddictionary/
│   └── v1/
│       ├── common_pb.ts          # Enums and common messages
│       ├── auth_pb.ts            # Auth request/response messages
│       ├── auth_connect.ts       # AuthService client
│       ├── conversation_pb.ts    # Conversation messages
│       ├── conversation_connect.ts
│       ├── message_pb.ts         # Message messages
│       └── message_connect.ts
```

---

## Environment Variables

```env
# .env.development
VITE_API_BASE_URL=http://localhost:8080

# .env.production
VITE_API_BASE_URL=https://api.kiddictionary.example.com
```