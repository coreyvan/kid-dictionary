import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Conversation } from '@/types'
import {
  listConversations as apiListConversations,
  getConversation as apiGetConversation,
  deleteConversation as apiDeleteConversation,
} from '@/services/api/conversations'

export const useConversationsStore = defineStore('conversations', () => {
  const conversations = ref<Conversation[]>([])
  const currentConversation = ref<Conversation | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const nextPageToken = ref<string | null>(null)
  const hasMore = computed(() => !!nextPageToken.value)
  const isEmpty = computed(() => conversations.value.length === 0 && !isLoading.value)

  async function loadConversations(refresh: boolean = false): Promise<void> {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const result = await apiListConversations(
        20,
        refresh ? undefined : nextPageToken.value || undefined
      )

      if (refresh) {
        conversations.value = result.conversations
      } else {
        conversations.value = [...conversations.value, ...result.conversations]
      }
      nextPageToken.value = result.nextPageToken
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load conversations'
    } finally {
      isLoading.value = false
    }
  }

  async function loadMore(): Promise<void> {
    if (!hasMore.value || isLoading.value) return
    await loadConversations(false)
  }

  async function refresh(): Promise<void> {
    nextPageToken.value = null
    await loadConversations(true)
  }

  async function loadConversation(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const result = await apiGetConversation(id)
      currentConversation.value = result.conversation

      // Return messages to be used by the messages store
      return result.messages as unknown as void
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load conversation'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function deleteConversation(id: string): Promise<void> {
    try {
      await apiDeleteConversation(id)
      conversations.value = conversations.value.filter((c) => c.id !== id)

      if (currentConversation.value?.id === id) {
        currentConversation.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete conversation'
      throw err
    }
  }

  function clearError() {
    error.value = null
  }

  function reset() {
    conversations.value = []
    currentConversation.value = null
    nextPageToken.value = null
    error.value = null
    isLoading.value = false
  }

  return {
    // State
    conversations,
    currentConversation,
    isLoading,
    error,
    nextPageToken,

    // Computed
    hasMore,
    isEmpty,

    // Actions
    loadConversations,
    loadMore,
    refresh,
    loadConversation,
    deleteConversation,
    clearError,
    reset,
  }
})
