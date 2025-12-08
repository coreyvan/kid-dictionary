import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { AuthService } from '@/gen/kiddictionary/v1/auth_pb';
import { ConversationService } from '@/gen/kiddictionary/v1/conversation_pb';
import { MessageService } from '@/gen/kiddictionary/v1/message_pb';

const baseUrl = import.meta.env.VITE_API_BASE_URL || '';

// Create transport without interceptor initially to avoid circular dependency
// Auth interceptor will be added when making authenticated requests
const transport = createConnectTransport({
  baseUrl,
});

export const authClient = createClient(AuthService, transport);
export const conversationClient = createClient(ConversationService, transport);
export const messageClient = createClient(MessageService, transport);

export { transport };