import { messageClient } from './client'
import { getErrorMessage } from './errors'
import type { Message } from '@/types'

export interface SendMessageResult {
  userMessage: Message
  assistantMessage: Message
}

export async function sendMessage(
  conversationId: string,
  content: string
): Promise<SendMessageResult> {
  try {
    const response = await messageClient.sendMessage({
      conversationId,
      content,
    })

    if (!response.userMessage || !response.assistantMessage) {
      throw new Error('Invalid response from server')
    }

    return {
      userMessage: response.userMessage as Message,
      assistantMessage: response.assistantMessage as Message,
    }
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}