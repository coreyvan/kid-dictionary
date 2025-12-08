<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login } from '@/services/api/auth'
import type { LoginFormData } from '@/types'
import LoginForm from '@/components/auth/LoginForm.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const isLoading = ref(false)
const error = ref<string | null>(null)

async function handleSubmit(data: LoginFormData) {
  isLoading.value = true
  error.value = null

  try {
    const result = await login(data.email, data.password)
    authStore.setTokens(result.tokens.accessToken, result.tokens.refreshToken)
    authStore.setUser(result.user)

    // Redirect to intended destination or home
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Login failed'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="max-w-md mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold text-gray-900 mb-2">Sign In</h1>
    <p class="text-gray-600 mb-6">
      Sign in to access your conversation history.
    </p>

    <LoginForm
      :is-loading="isLoading"
      :error="error"
      @submit="handleSubmit"
    />

    <p class="mt-6 text-center text-sm text-gray-600">
      Don't have an account?
      <RouterLink to="/register" class="text-bracket-preteen hover:underline font-medium">
        Create one
      </RouterLink>
    </p>
  </div>
</template>