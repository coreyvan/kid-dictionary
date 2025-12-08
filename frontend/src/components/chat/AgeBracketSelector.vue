<script setup lang="ts">
import { AgeBracket, AgeBracketDisplay } from '@/types'

const props = defineProps<{
  modelValue: AgeBracket
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: AgeBracket]
  'change': [oldValue: AgeBracket, newValue: AgeBracket]
}>()

const brackets = [
  { value: AgeBracket.LITTLE_ONES, ...AgeBracketDisplay[1] },
  { value: AgeBracket.GROWING_MINDS, ...AgeBracketDisplay[2] },
  { value: AgeBracket.PRE_TEENS, ...AgeBracketDisplay[3] },
]

function selectBracket(newValue: AgeBracket) {
  if (newValue !== props.modelValue) {
    const oldValue = props.modelValue
    emit('update:modelValue', newValue)
    emit('change', oldValue, newValue)
  }
}
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="bracket in brackets"
      :key="bracket.value"
      type="button"
      :disabled="disabled"
      class="min-h-[44px] px-4 py-2 rounded-lg font-medium transition-all duration-300 ease-in-out transform focus:outline-none focus:ring-2 focus:ring-offset-2"
      :class="[
        modelValue === bracket.value
          ? `bg-${bracket.color} text-white focus:ring-${bracket.color} scale-105 shadow-md`
          : `bg-gray-100 text-gray-700 hover:bg-gray-200 focus:ring-gray-400 hover:scale-102`,
        disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
      ]"
      @click="selectBracket(bracket.value)"
    >
      <span class="block text-sm font-semibold">{{ bracket.label }}</span>
      <span class="block text-xs opacity-75">{{ bracket.description }}</span>
    </button>
  </div>
</template>