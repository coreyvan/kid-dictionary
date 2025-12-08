import { messageClient } from './client'
import { getErrorMessage } from './errors'
import type { Message, Conversation } from '@/types'
import { AgeBracket as ProtoAgeBracket } from '@/gen/kiddictionary/v1/common_pb'
import type { AgeBracket } from '@/types'

export interface SendMessageResult {
  userMessage: Message
  assistantMessage: Message
  conversation?: Conversation // Present when conversation was auto-created
}

// Map frontend AgeBracket to proto AgeBracket
function toProtoAgeBracket(ageBracket: AgeBracket): ProtoAgeBracket {
  // The values should match, but we'll be explicit
  return ageBracket as ProtoAgeBracket
}

export async function sendMessage(
  content: string,
  options: {
    conversationId?: string
    ageBracket?: AgeBracket
  } = {}
): Promise<SendMessageResult> {
  try {
    const response = await messageClient.sendMessage({
      conversationId: options.conversationId ?? '',
      content,
      ageBracket: options.ageBracket ? toProtoAgeBracket(options.ageBracket) : ProtoAgeBracket.UNSPECIFIED,
    })

    if (!response.userMessage || !response.assistantMessage) {
      throw new Error('Invalid response from server')
    }

    const result: SendMessageResult = {
      userMessage: response.userMessage as Message,
      assistantMessage: response.assistantMessage as Message,
    }

    // Include conversation if it was auto-created
    if (response.conversation) {
      result.conversation = response.conversation as unknown as Conversation
    }

    return result
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}