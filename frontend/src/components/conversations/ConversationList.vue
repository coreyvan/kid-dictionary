<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { Conversation } from '@/types'
import ConversationCard from './ConversationCard.vue'
import DeleteConfirmDialog from './DeleteConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

const props = defineProps<{
  conversations: Conversation[]
  isLoading: boolean
  hasMore: boolean
  isEmpty: boolean
}>()

const emit = defineEmits<{
  select: [id: string]
  delete: [id: string]
  loadMore: []
}>()

const deleteDialogOpen = ref(false)
const conversationToDelete = ref<string | null>(null)
const isDeleting = ref(false)
const listContainer = ref<HTMLElement | null>(null)

function openDeleteDialog(id: string) {
  conversationToDelete.value = id
  deleteDialogOpen.value = true
}

function closeDeleteDialog() {
  deleteDialogOpen.value = false
  conversationToDelete.value = null
}

async function confirmDelete() {
  if (!conversationToDelete.value) return

  isDeleting.value = true
  emit('delete', conversationToDelete.value)

  // Wait a bit for the delete to complete
  setTimeout(() => {
    isDeleting.value = false
    closeDeleteDialog()
  }, 500)
}

// Infinite scroll
function handleScroll() {
  if (!listContainer.value || props.isLoading || !props.hasMore) return

  const { scrollTop, scrollHeight, clientHeight } = listContainer.value
  if (scrollTop + clientHeight >= scrollHeight - 100) {
    emit('loadMore')
  }
}

onMounted(() => {
  listContainer.value?.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  listContainer.value?.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <div
    ref="listContainer"
    class="h-full overflow-y-auto"
  >
    <EmptyState
      v-if="isEmpty"
      title="No conversations yet"
      description="Start a new chat to see your conversation history here."
      icon="chat"
    >
      <RouterLink
        to="/chat"
        class="inline-flex items-center min-h-[44px] px-4 py-2 bg-bracket-preteen text-white rounded-lg font-medium hover:opacity-90"
      >
        Start a Conversation
      </RouterLink>
    </EmptyState>

    <div v-else class="space-y-2 sm:space-y-3">
      <ConversationCard
        v-for="conversation in conversations"
        :key="conversation.id"
        :conversation="conversation"
        @click="emit('select', conversation.id)"
        @delete="openDeleteDialog(conversation.id)"
      />

      <!-- Loading indicator for infinite scroll -->
      <div
        v-if="isLoading && conversations.length > 0"
        class="py-4 text-center"
      >
        <LoadingSpinner />
      </div>

      <!-- End of list indicator -->
      <p
        v-if="!hasMore && conversations.length > 0"
        class="py-4 text-center text-sm text-gray-500"
      >
        No more conversations
      </p>
    </div>

    <!-- Initial loading state -->
    <div
      v-if="isLoading && conversations.length === 0"
      class="py-12 text-center"
    >
      <LoadingSpinner size="lg" />
      <p class="mt-4 text-gray-500">Loading conversations...</p>
    </div>

    <DeleteConfirmDialog
      :is-open="deleteDialogOpen"
      :is-deleting="isDeleting"
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>
