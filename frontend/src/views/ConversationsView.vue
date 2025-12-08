<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConversationsStore } from '@/stores/conversations'
import ConversationList from '@/components/conversations/ConversationList.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'

const router = useRouter()
const conversationsStore = useConversationsStore()

onMounted(() => {
  conversationsStore.refresh()
})

function handleSelect(id: string) {
  router.push(`/chat/${id}`)
}

async function handleDelete(id: string) {
  try {
    await conversationsStore.deleteConversation(id)
  } catch {
    // Error is already set in the store
  }
}

function handleLoadMore() {
  conversationsStore.loadMore()
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Conversation History</h1>
      <RouterLink
        to="/chat"
        class="min-h-[44px] px-4 py-2 bg-bracket-preteen text-white rounded-lg font-medium hover:opacity-90 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Chat
      </RouterLink>
    </div>

    <ErrorMessage
      v-if="conversationsStore.error"
      :message="conversationsStore.error"
      dismissible
      class="mb-4"
      @dismiss="conversationsStore.clearError()"
    />

    <ConversationList
      :conversations="conversationsStore.conversations"
      :is-loading="conversationsStore.isLoading"
      :has-more="conversationsStore.hasMore"
      :is-empty="conversationsStore.isEmpty"
      @select="handleSelect"
      @delete="handleDelete"
      @load-more="handleLoadMore"
    />
  </div>
</template>
