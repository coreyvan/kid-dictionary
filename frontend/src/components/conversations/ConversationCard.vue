<script setup lang="ts">
import type { Conversation } from '@/types'
import { AgeBracketDisplay } from '@/types'

const props = defineProps<{
  conversation: Conversation
}>()

defineEmits<{
  click: []
  delete: []
}>()

function formatDate(timestamp: unknown): string {
  if (!timestamp) return ''

  // Handle protobuf Timestamp type
  const date = timestamp instanceof Date
    ? timestamp
    : new Date((timestamp as { seconds?: bigint }).seconds
        ? Number((timestamp as { seconds: bigint }).seconds) * 1000
        : Date.now())

  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  } else if (days === 1) {
    return 'Yesterday'
  } else if (days < 7) {
    return date.toLocaleDateString([], { weekday: 'long' })
  } else {
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
  }
}

const bracketInfo = AgeBracketDisplay[props.conversation.ageBracket as keyof typeof AgeBracketDisplay]
</script>

<template>
  <div
    class="bg-white border border-gray-200 rounded-lg p-4 hover:border-gray-300 hover:shadow-sm transition-all cursor-pointer"
    @click="$emit('click')"
  >
    <div class="flex items-start justify-between gap-3">
      <div class="flex-1 min-w-0">
        <h3 class="font-medium text-gray-900 truncate">
          {{ conversation.title || 'Untitled conversation' }}
        </h3>
        <div class="flex items-center gap-2 mt-1">
          <span
            v-if="bracketInfo"
            class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
            :class="`bg-${bracketInfo.color}/10 text-${bracketInfo.color}`"
            :style="{ backgroundColor: `color-mix(in srgb, var(--color-${bracketInfo.color}) 10%, white)`, color: `var(--color-${bracketInfo.color})` }"
          >
            {{ bracketInfo.label }}
          </span>
          <span class="text-xs text-gray-500">
            {{ formatDate(conversation.updatedAt || conversation.createdAt) }}
          </span>
        </div>
      </div>
      <button
        type="button"
        class="min-h-[44px] min-w-[44px] p-2 text-gray-400 hover:text-red-500 rounded-lg hover:bg-red-50 transition-colors"
        @click.stop="$emit('delete')"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          />
        </svg>
        <span class="sr-only">Delete</span>
      </button>
    </div>
  </div>
</template>
