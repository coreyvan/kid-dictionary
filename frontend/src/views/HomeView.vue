<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-12">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-gray-900 mb-4">
        Kid Dictionary
      </h1>
      <p class="text-xl text-gray-600 mb-8">
        Age-appropriate explanations for curious kids
      </p>

      <RouterLink
        to="/chat"
        class="inline-flex items-center justify-center min-h-[44px] px-8 py-3 bg-bracket-preteen text-white text-lg font-semibold rounded-lg hover:opacity-90 transition-opacity"
      >
        Start Asking Questions
      </RouterLink>

      <div class="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="p-6 bg-bracket-little/10 rounded-lg">
          <h3 class="font-semibold text-gray-900 mb-2">Little Ones</h3>
          <p class="text-gray-600 text-sm">Ages 0-5: Simple words and short answers</p>
        </div>
        <div class="p-6 bg-bracket-growing/10 rounded-lg">
          <h3 class="font-semibold text-gray-900 mb-2">Growing Minds</h3>
          <p class="text-gray-600 text-sm">Ages 5-10: Fun facts and easy examples</p>
        </div>
        <div class="p-6 bg-bracket-preteen/10 rounded-lg">
          <h3 class="font-semibold text-gray-900 mb-2">Pre-Teens</h3>
          <p class="text-gray-600 text-sm">Ages 10+: Deeper explanations and context</p>
        </div>
      </div>
    </div>
  </div>
</template>