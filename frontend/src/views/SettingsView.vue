<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { AgeBracket } from '@/types'
import AgeBracketSelector from '@/components/chat/AgeBracketSelector.vue'

const authStore = useAuthStore()

const defaultAgeBracket = ref<AgeBracket>(AgeBracket.GROWING_MINDS)
const saved = ref(false)

function handleBracketChange(_old: AgeBracket, newValue: AgeBracket) {
  defaultAgeBracket.value = newValue
  // In a real app, this would be saved to the user's profile via API
  localStorage.setItem('defaultAgeBracket', String(newValue))
  saved.value = true
  setTimeout(() => {
    saved.value = false
  }, 2000)
}

// Load saved preference on mount
const savedBracket = localStorage.getItem('defaultAgeBracket')
if (savedBracket) {
  defaultAgeBracket.value = Number(savedBracket) as AgeBracket
}
</script>

<template>
  <div class="max-w-md mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Settings</h1>

    <!-- Account Info -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">Account</h2>
      <div class="bg-white border border-gray-200 rounded-lg p-4">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-bracket-preteen rounded-full flex items-center justify-center">
            <span class="text-white text-lg font-bold">
              {{ authStore.user?.email?.charAt(0).toUpperCase() || 'U' }}
            </span>
          </div>
          <div>
            <p class="font-medium text-gray-900">{{ authStore.user?.email || 'User' }}</p>
            <p class="text-sm text-gray-500">Registered user</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Default Age Bracket -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-gray-800 mb-2">Default Age Bracket</h2>
      <p class="text-sm text-gray-500 mb-4">
        Choose the default age bracket for new conversations.
      </p>
      <AgeBracketSelector
        v-model="defaultAgeBracket"
        @change="handleBracketChange"
      />
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <p v-if="saved" class="mt-2 text-sm text-green-600">
          Preference saved!
        </p>
      </Transition>
    </section>

    <!-- Logout -->
    <section>
      <button
        type="button"
        class="w-full min-h-[44px] px-4 py-2 border border-red-300 text-red-600 rounded-lg hover:bg-red-50 font-medium"
        @click="authStore.clearAuth()"
      >
        Sign Out
      </button>
    </section>
  </div>
</template>