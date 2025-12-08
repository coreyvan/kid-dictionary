import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Message, PendingMessage } from '@/types'
import { AgeBracket, AgeBracketDisplay } from '@/types'
import { sendMessage as apiSendMessage } from '@/services/api/messages'
import { updateConversation as apiUpdateConversation } from '@/services/api/conversations'
import {
  createAnonymousConversation,
  getAnonymousConversation,
  addMessageToAnonymousConversation,
  updateAnonymousConversation,
  getCurrentAnonymousConversationId,
  setCurrentAnonymousConversationId,
} from '@/services/storage/anonymous'
import { useAuthStore } from './auth'

export const useMessagesStore = defineStore('messages', () => {
  const messages = ref<Message[]>([])
  const pendingMessage = ref<PendingMessage | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const currentConversationId = ref<string | null>(null)
  const currentAgeBracket = ref<AgeBracket>(AgeBracket.GROWING_MINDS)

  const hasMessages = computed(() => messages.value.length > 0)

  // Toast notification state
  const toastMessage = ref<string | null>(null)
  const toastType = ref<'success' | 'error'>('success')

  function showToast(message: string, type: 'success' | 'error' = 'success') {
    toastMessage.value = message
    toastType.value = type
    setTimeout(() => {
      toastMessage.value = null
    }, 3000)
  }

  function clearToast() {
    toastMessage.value = null
  }

  async function setAgeBracket(bracket: AgeBracket): Promise<void> {
    const previousBracket = currentAgeBracket.value
    currentAgeBracket.value = bracket

    // Update conversation if one exists
    if (currentConversationId.value) {
      const authStore = useAuthStore()
      if (authStore.isAuthenticated) {
        // Update via API for authenticated users
        try {
          await apiUpdateConversation(currentConversationId.value, {
            ageBracket: bracket,
          })
          const bracketInfo = AgeBracketDisplay[bracket as keyof typeof AgeBracketDisplay]
          showToast(`Age bracket changed to ${bracketInfo?.label || 'new bracket'}`)
        } catch (err) {
          // Revert on error
          currentAgeBracket.value = previousBracket
          showToast('Failed to update age bracket', 'error')
        }
      } else {
        // Update local storage for anonymous users
        updateAnonymousConversation(currentConversationId.value, {
          ageBracket: bracket,
        })
        const bracketInfo = AgeBracketDisplay[bracket as keyof typeof AgeBracketDisplay]
        showToast(`Age bracket changed to ${bracketInfo?.label || 'new bracket'}`)
      }
    }
  }

  function loadConversation(conversationId: string) {
    const authStore = useAuthStore()

    if (!authStore.isAuthenticated) {
      // Load from anonymous storage
      const conv = getAnonymousConversation(conversationId)
      if (conv) {
        messages.value = conv.messages as Message[]
        currentConversationId.value = conversationId
        currentAgeBracket.value = conv.ageBracket
        setCurrentAnonymousConversationId(conversationId)
      }
    }
    // For authenticated users, messages will be loaded from API in ChatView
  }

  function startNewConversation() {
    messages.value = []
    currentConversationId.value = null
    error.value = null
    pendingMessage.value = null
    setCurrentAnonymousConversationId(null)
  }

  async function sendMessage(content: string): Promise<void> {
    const authStore = useAuthStore()
    error.value = null

    // Create optimistic user message
    const tempId = crypto.randomUUID()
    pendingMessage.value = {
      id: tempId,
      content,
      status: 'sending',
    }

    try {
      isLoading.value = true

      // If no conversation yet, create one (anonymous or authenticated)
      if (!currentConversationId.value) {
        if (!authStore.isAuthenticated) {
          // Create anonymous conversation
          const title = content.slice(0, 50) + (content.length > 50 ? '...' : '')
          const conv = createAnonymousConversation(title, currentAgeBracket.value)
          currentConversationId.value = conv.id
        }
        // For authenticated users, the conversation will be created by the backend
      }

      // Send to API
      const result = await apiSendMessage(
        currentConversationId.value || '',
        content
      )

      // Clear pending and add real messages
      pendingMessage.value = null

      // Add user message
      messages.value.push(result.userMessage)

      // Add assistant message
      messages.value.push(result.assistantMessage)

      // Save to anonymous storage if not authenticated
      if (!authStore.isAuthenticated && currentConversationId.value) {
        addMessageToAnonymousConversation(currentConversationId.value, result.userMessage)
        addMessageToAnonymousConversation(currentConversationId.value, result.assistantMessage)
      }
    } catch (err) {
      // Mark pending message as error
      if (pendingMessage.value) {
        pendingMessage.value.status = 'error'
        pendingMessage.value.error = err instanceof Error ? err.message : 'Failed to send message'
      }
      error.value = err instanceof Error ? err.message : 'Failed to send message'
    } finally {
      isLoading.value = false
    }
  }

  async function retryPendingMessage(): Promise<void> {
    if (!pendingMessage.value || pendingMessage.value.status !== 'error') {
      return
    }

    const content = pendingMessage.value.content
    pendingMessage.value = null
    await sendMessage(content)
  }

  function clearError() {
    error.value = null
  }

  // Initialize from storage on store creation
  function initialize() {
    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) {
      const currentId = getCurrentAnonymousConversationId()
      if (currentId) {
        loadConversation(currentId)
      }
    }
  }

  return {
    // State
    messages,
    pendingMessage,
    isLoading,
    error,
    currentConversationId,
    currentAgeBracket,
    toastMessage,
    toastType,

    // Computed
    hasMessages,

    // Actions
    setAgeBracket,
    loadConversation,
    startNewConversation,
    sendMessage,
    retryPendingMessage,
    clearError,
    clearToast,
    initialize,
  }
})