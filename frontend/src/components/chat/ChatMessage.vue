<script setup lang="ts">
import { computed } from 'vue'
import { MessageRole, ContentTier } from '@/types'
import ContentTierIndicator from './ContentTierIndicator.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

const props = defineProps<{
  role: MessageRole
  content: string
  contentTier?: ContentTier
  isLoading?: boolean
  hasError?: boolean
  errorMessage?: string
}>()

const emit = defineEmits<{
  retry: []
}>()

const isUser = computed(() => props.role === MessageRole.USER)
const isAssistant = computed(() => props.role === MessageRole.ASSISTANT)

const showTierIndicator = computed(() => {
  return isAssistant.value && props.contentTier && props.contentTier !== ContentTier.NORMAL
})
</script>

<template>
  <div
    class="flex gap-3"
    :class="isUser ? 'flex-row-reverse' : ''"
  >
    <!-- Avatar -->
    <div
      class="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-medium"
      :class="isUser ? 'bg-bracket-preteen' : 'bg-gray-600'"
    >
      {{ isUser ? 'U' : 'AI' }}
    </div>

    <!-- Message bubble -->
    <div
      class="max-w-[80%] rounded-lg px-4 py-3"
      :class="[
        isUser
          ? 'bg-bracket-preteen text-white'
          : 'bg-gray-100 text-gray-900',
        hasError ? 'border-2 border-red-300' : '',
      ]"
    >
      <!-- Loading state -->
      <div v-if="isLoading" class="flex items-center gap-2">
        <LoadingSpinner size="sm" />
        <span class="text-gray-500">Thinking...</span>
      </div>

      <!-- Error state -->
      <div v-else-if="hasError" class="space-y-2">
        <p class="text-red-600 text-sm">{{ errorMessage || 'Failed to send message' }}</p>
        <button
          type="button"
          class="min-h-[44px] px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200"
          @click="emit('retry')"
        >
          Retry
        </button>
      </div>

      <!-- Content -->
      <template v-else>
        <p class="whitespace-pre-wrap">{{ content }}</p>
        <ContentTierIndicator
          v-if="showTierIndicator"
          :tier="contentTier!"
          class="mt-2"
        />
      </template>
    </div>
  </div>
</template>