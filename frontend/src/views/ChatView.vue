<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessagesStore } from '@/stores/messages'
import { MessageRole, AgeBracket } from '@/types'
import AgeBracketSelector from '@/components/chat/AgeBracketSelector.vue'
import ChatMessage from '@/components/chat/ChatMessage.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'

const route = useRoute()
const messagesStore = useMessagesStore()
const messagesContainer = ref<HTMLElement | null>(null)

// Load conversation if ID provided in route
onMounted(() => {
  const conversationId = route.params.id as string | undefined
  if (conversationId) {
    messagesStore.loadConversation(conversationId)
  } else {
    messagesStore.initialize()
  }
})

// Auto-scroll to bottom when new messages arrive
watch(
  () => messagesStore.messages.length,
  () => {
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })
  }
)

async function handleSend(content: string) {
  await messagesStore.sendMessage(content)
}

function handleRetry() {
  messagesStore.retryPendingMessage()
}

function startNewChat() {
  messagesStore.startNewConversation()
}

function handleAgeBracketChange(_oldValue: AgeBracket, newValue: AgeBracket) {
  messagesStore.setAgeBracket(newValue)
}
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-64px)]">
    <!-- Header with age bracket selector -->
    <div class="border-b border-gray-200 bg-white px-4 py-3">
      <div class="max-w-4xl mx-auto">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-lg font-semibold text-gray-900">
            {{ messagesStore.hasMessages ? 'Chat' : 'Start a Conversation' }}
          </h1>
          <button
            v-if="messagesStore.hasMessages"
            type="button"
            class="min-h-[44px] px-4 text-sm text-gray-600 hover:text-gray-900"
            @click="startNewChat"
          >
            New Chat
          </button>
        </div>

        <AgeBracketSelector
          v-model="messagesStore.currentAgeBracket"
          :disabled="messagesStore.isLoading"
          @change="handleAgeBracketChange"
        />
      </div>
    </div>

    <!-- Messages area -->
    <div
      ref="messagesContainer"
      class="flex-1 overflow-y-auto px-4 py-6"
    >
      <div class="max-w-4xl mx-auto space-y-4">
        <!-- Empty state -->
        <div
          v-if="!messagesStore.hasMessages && !messagesStore.pendingMessage"
          class="text-center py-12"
        >
          <p class="text-gray-500 mb-2">
            Ask any question and get an age-appropriate answer!
          </p>
          <p class="text-gray-400 text-sm">
            Select an age bracket above to customize the response.
          </p>
        </div>

        <!-- Messages -->
        <ChatMessage
          v-for="message in messagesStore.messages"
          :key="message.id"
          :role="message.role"
          :content="message.content"
          :content-tier="message.contentTier"
        />

        <!-- Pending message (optimistic) -->
        <ChatMessage
          v-if="messagesStore.pendingMessage"
          :role="MessageRole.USER"
          :content="messagesStore.pendingMessage.content"
          :has-error="messagesStore.pendingMessage.status === 'error'"
          :error-message="messagesStore.pendingMessage.error"
          @retry="handleRetry"
        />

        <!-- Loading indicator for assistant response -->
        <ChatMessage
          v-if="messagesStore.isLoading && messagesStore.pendingMessage?.status === 'sending'"
          :role="MessageRole.ASSISTANT"
          content=""
          :is-loading="true"
        />
      </div>
    </div>

    <!-- Error banner -->
    <div v-if="messagesStore.error && !messagesStore.pendingMessage" class="px-4 py-2 bg-red-50">
      <div class="max-w-4xl mx-auto">
        <ErrorMessage
          :message="messagesStore.error"
          dismissible
          @dismiss="messagesStore.clearError()"
        />
      </div>
    </div>

    <!-- Input area (fixed at bottom) -->
    <div class="border-t border-gray-200 bg-white px-4 py-3">
      <div class="max-w-4xl mx-auto">
        <ChatInput
          :disabled="messagesStore.isLoading"
          :placeholder="messagesStore.hasMessages ? 'Ask a follow-up question...' : 'Ask a question...'"
          @send="handleSend"
        />
      </div>
    </div>

    <!-- Toast notification -->
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-2"
    >
      <div
        v-if="messagesStore.toastMessage"
        class="fixed bottom-24 left-1/2 -translate-x-1/2 px-4 py-2 rounded-lg shadow-lg text-white text-sm font-medium z-50"
        :class="messagesStore.toastType === 'success' ? 'bg-green-600' : 'bg-red-600'"
        @click="messagesStore.clearToast()"
      >
        {{ messagesStore.toastMessage }}
      </div>
    </Transition>
  </div>
</template>