// Re-export generated types from protobuf
export {
  AgeBracket,
  MessageRole,
  ContentTier,
  type User,
  type Message,
} from '@/gen/kiddictionary/v1/common_pb';

export { type Conversation } from '@/gen/kiddictionary/v1/conversation_pb';

// Frontend-specific state types

export interface AuthState {
  user: import('@/gen/kiddictionary/v1/common_pb').User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface ConversationsState {
  conversations: import('@/gen/kiddictionary/v1/conversation_pb').Conversation[];
  currentConversation: import('@/gen/kiddictionary/v1/conversation_pb').Conversation | null;
  isLoading: boolean;
  error: string | null;
  nextPageToken: string | null;
  hasMore: boolean;
}

export interface MessagesState {
  messages: import('@/gen/kiddictionary/v1/common_pb').Message[];
  pendingMessage: PendingMessage | null;
  isLoading: boolean;
  error: string | null;
}

export interface PendingMessage {
  id: string;
  content: string;
  status: 'sending' | 'error';
  error?: string;
}

// Anonymous conversation stored in sessionStorage
export interface AnonymousConversation {
  id: string;
  title: string;
  ageBracket: import('@/gen/kiddictionary/v1/common_pb').AgeBracket;
  messages: import('@/gen/kiddictionary/v1/common_pb').Message[];
  createdAt: Date;
}

// Display helpers for enums
export const AgeBracketDisplay = {
  1: { label: 'Little Ones', color: 'bracket-little', description: '0-5 years' },
  2: { label: 'Growing Minds', color: 'bracket-growing', description: '5-10 years' },
  3: { label: 'Pre-Teens', color: 'bracket-preteen', description: '10+ years' },
} as const;

export const ContentTierIndicator = {
  1: null, // NORMAL - no indicator
  2: { icon: 'heart', label: 'Sensitive topic' },
  3: { icon: 'info', label: 'Context-dependent' },
  4: { icon: 'shield', label: 'Alternative suggested' },
} as const;

// Form types
export interface LoginFormData {
  email: string;
  password: string;
}

export interface RegisterFormData {
  email: string;
  password: string;
  confirmPassword: string;
  defaultAgeBracket: import('@/gen/kiddictionary/v1/common_pb').AgeBracket;
}

export interface FormErrors {
  email?: string;
  password?: string;
  confirmPassword?: string;
}