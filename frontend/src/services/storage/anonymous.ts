import type { AnonymousConversation } from '@/types'
import { AgeBracket } from '@/types'

const CONVERSATIONS_KEY = 'anonymous:conversations'
const CURRENT_ID_KEY = 'anonymous:currentId'

// Serialized format for storage (dates as ISO strings)
interface StoredConversation {
  id: string
  title: string
  ageBracket: AgeBracket
  messages: StoredMessage[]
  createdAt: string
}

interface StoredMessage {
  id: string
  conversationId: string
  role: number
  content: string
  createdAt?: string
  contentTier: number
}

export function getAnonymousConversations(): AnonymousConversation[] {
  const stored = sessionStorage.getItem(CONVERSATIONS_KEY)
  if (!stored) return []

  try {
    const parsed: StoredConversation[] = JSON.parse(stored)
    return parsed.map((conv) => ({
      id: conv.id,
      title: conv.title,
      ageBracket: conv.ageBracket,
      createdAt: new Date(conv.createdAt),
      messages: conv.messages.map((msg) => ({
        id: msg.id,
        conversationId: msg.conversationId,
        role: msg.role,
        content: msg.content,
        contentTier: msg.contentTier,
        // Note: createdAt from proto is Timestamp, we store as ISO string
      })),
    })) as AnonymousConversation[]
  } catch {
    return []
  }
}

export function saveAnonymousConversations(conversations: AnonymousConversation[]): void {
  const toStore: StoredConversation[] = conversations.map((conv) => ({
    id: conv.id,
    title: conv.title,
    ageBracket: conv.ageBracket,
    createdAt: conv.createdAt.toISOString(),
    messages: conv.messages.map((msg) => ({
      id: msg.id,
      conversationId: msg.conversationId,
      role: msg.role,
      content: msg.content,
      contentTier: msg.contentTier,
    })),
  }))
  sessionStorage.setItem(CONVERSATIONS_KEY, JSON.stringify(toStore))
}

export function getAnonymousConversation(id: string): AnonymousConversation | null {
  const conversations = getAnonymousConversations()
  const found = conversations.find((c) => c.id === id)
  return found ?? null
}

export function createAnonymousConversation(
  id: string,
  title: string,
  ageBracket: AgeBracket
): AnonymousConversation {
  const conversation: AnonymousConversation = {
    id,
    title,
    ageBracket,
    messages: [],
    createdAt: new Date(),
  }

  const conversations = getAnonymousConversations()
  conversations.unshift(conversation)
  saveAnonymousConversations(conversations)

  return conversation
}

export function updateAnonymousConversation(
  id: string,
  updates: Partial<Pick<AnonymousConversation, 'title' | 'ageBracket' | 'messages'>>
): AnonymousConversation | null {
  const conversations = getAnonymousConversations()
  const existing = conversations.find((c) => c.id === id)

  if (!existing) return null

  const updated: AnonymousConversation = {
    id: existing.id,
    title: updates.title ?? existing.title,
    ageBracket: updates.ageBracket ?? existing.ageBracket,
    messages: updates.messages ?? existing.messages,
    createdAt: existing.createdAt,
  }

  const index = conversations.findIndex((c) => c.id === id)
  conversations[index] = updated
  saveAnonymousConversations(conversations)

  return updated
}

export function addMessageToAnonymousConversation(
  conversationId: string,
  message: AnonymousConversation['messages'][0]
): void {
  const conversations = getAnonymousConversations()
  const conv = conversations.find((c) => c.id === conversationId)

  if (conv) {
    conv.messages.push(message)
    saveAnonymousConversations(conversations)
  }
}

export function getCurrentAnonymousConversationId(): string | null {
  return sessionStorage.getItem(CURRENT_ID_KEY)
}

export function setCurrentAnonymousConversationId(id: string | null): void {
  if (id) {
    sessionStorage.setItem(CURRENT_ID_KEY, id)
  } else {
    sessionStorage.removeItem(CURRENT_ID_KEY)
  }
}

export function clearAnonymousData(): void {
  sessionStorage.removeItem(CONVERSATIONS_KEY)
  sessionStorage.removeItem(CURRENT_ID_KEY)
}