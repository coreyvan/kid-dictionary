import { conversationClient } from './client'
import { getErrorMessage } from './errors'
import type { Conversation, Message, AgeBracket } from '@/types'

export interface ListConversationsResult {
  conversations: Conversation[]
  nextPageToken: string | null
}

export interface GetConversationResult {
  conversation: Conversation
  messages: Message[]
}

export async function listConversations(
  pageSize: number = 20,
  pageToken?: string
): Promise<ListConversationsResult> {
  try {
    const response = await conversationClient.listConversations({
      pageSize,
      pageToken: pageToken || '',
    })

    return {
      conversations: (response.conversations || []) as Conversation[],
      nextPageToken: response.nextPageToken || null,
    }
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}

export async function getConversation(id: string): Promise<GetConversationResult> {
  try {
    const response = await conversationClient.getConversation({ id })

    if (!response.conversation) {
      throw new Error('Conversation not found')
    }

    return {
      conversation: response.conversation as Conversation,
      messages: (response.messages || []) as Message[],
    }
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}

export async function createConversation(
  title: string,
  ageBracket: AgeBracket
): Promise<Conversation> {
  try {
    const response = await conversationClient.createConversation({
      title,
      ageBracket,
    })

    if (!response.conversation) {
      throw new Error('Failed to create conversation')
    }

    return response.conversation as Conversation
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}

export async function updateConversation(
  id: string,
  updates: { title?: string; ageBracket?: AgeBracket }
): Promise<Conversation> {
  try {
    const response = await conversationClient.updateConversation({
      id,
      ...updates,
    })

    if (!response.conversation) {
      throw new Error('Failed to update conversation')
    }

    return response.conversation as Conversation
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}

export async function deleteConversation(id: string): Promise<void> {
  try {
    await conversationClient.deleteConversation({ id })
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}
