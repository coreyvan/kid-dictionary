<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import OfflineIndicator from '@/components/common/OfflineIndicator.vue'

const authStore = useAuthStore()

// PWA install prompt
const deferredPrompt = ref<Event | null>(null)
const showInstallPrompt = ref(false)

onMounted(() => {
  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    deferredPrompt.value = e
    showInstallPrompt.value = true
  })

  window.addEventListener('appinstalled', () => {
    showInstallPrompt.value = false
    deferredPrompt.value = null
  })
})

async function installApp() {
  if (!deferredPrompt.value) return

  const promptEvent = deferredPrompt.value as BeforeInstallPromptEvent
  promptEvent.prompt()

  const { outcome } = await promptEvent.userChoice
  if (outcome === 'accepted') {
    showInstallPrompt.value = false
  }
  deferredPrompt.value = null
}

function dismissInstallPrompt() {
  showInstallPrompt.value = false
}

// Type for the install prompt event
interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <!-- Offline indicator -->
    <OfflineIndicator />
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 px-4 py-3">
      <div class="max-w-4xl mx-auto flex items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold text-gray-900">
          Kid Dictionary
        </RouterLink>

        <nav class="flex items-center gap-4">
          <template v-if="authStore.isAuthenticated">
            <RouterLink
              to="/conversations"
              class="text-gray-600 hover:text-gray-900"
            >
              History
            </RouterLink>
            <RouterLink
              to="/settings"
              class="text-gray-600 hover:text-gray-900"
            >
              Settings
            </RouterLink>
            <button
              type="button"
              class="min-h-[44px] px-4 text-gray-600 hover:text-gray-900"
              @click="authStore.clearAuth()"
            >
              Logout
            </button>
          </template>
          <template v-else>
            <RouterLink
              to="/login"
              class="min-h-[44px] px-4 flex items-center text-gray-600 hover:text-gray-900"
            >
              Login
            </RouterLink>
            <RouterLink
              to="/register"
              class="min-h-[44px] px-4 flex items-center bg-bracket-preteen text-white rounded-lg hover:opacity-90"
            >
              Register
            </RouterLink>
          </template>
        </nav>
      </div>
    </header>

    <!-- Main content -->
    <main class="flex-1">
      <RouterView />
    </main>

    <!-- PWA Install prompt -->
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="opacity-0 translate-y-full"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-full"
    >
      <div
        v-if="showInstallPrompt"
        class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-lg p-4 z-50"
      >
        <div class="max-w-4xl mx-auto flex items-center justify-between gap-4">
          <div class="flex-1">
            <p class="font-medium text-gray-900">Install Kid Dictionary</p>
            <p class="text-sm text-gray-500">Add to your home screen for quick access</p>
          </div>
          <div class="flex gap-2">
            <button
              type="button"
              class="min-h-[44px] px-4 py-2 text-gray-600 hover:text-gray-900"
              @click="dismissInstallPrompt"
            >
              Not now
            </button>
            <button
              type="button"
              class="min-h-[44px] px-4 py-2 bg-bracket-preteen text-white rounded-lg hover:opacity-90"
              @click="installApp"
            >
              Install
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>