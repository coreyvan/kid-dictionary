<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg';
}>();
</script>

<template>
  <div
    class="inline-block animate-spin rounded-full border-2 border-current border-t-transparent"
    :class="{
      'h-4 w-4': size === 'sm',
      'h-6 w-6': size === 'md' || !size,
      'h-8 w-8': size === 'lg',
    }"
    role="status"
    aria-label="Loading"
  >
    <span class="sr-only">Loading...</span>
  </div>
</template>