<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  disabled?: boolean
  placeholder?: string
}>()

const emit = defineEmits<{
  send: [content: string]
}>()

const content = ref('')

function handleSubmit() {
  const trimmed = content.value.trim()
  if (trimmed && !props.disabled) {
    emit('send', trimmed)
    content.value = ''
  }
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSubmit()
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="flex gap-2">
    <textarea
      v-model="content"
      :disabled="disabled"
      :placeholder="placeholder || 'Ask a question...'"
      rows="1"
      class="flex-1 min-h-[44px] max-h-32 px-4 py-3 border border-gray-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-bracket-preteen focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
      @keydown="handleKeydown"
    />
    <button
      type="submit"
      :disabled="disabled || !content.trim()"
      class="min-h-[44px] min-w-[44px] px-4 bg-bracket-preteen text-white rounded-lg font-medium hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-bracket-preteen focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-opacity"
    >
      <svg class="w-5 h-5 mx-auto" fill="currentColor" viewBox="0 0 20 20">
        <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" />
      </svg>
      <span class="sr-only">Send</span>
    </button>
  </form>
</template>