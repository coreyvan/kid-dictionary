<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { register } from '@/services/api/auth'
import type { RegisterFormData } from '@/types'
import RegisterForm from '@/components/auth/RegisterForm.vue'

const router = useRouter()
const authStore = useAuthStore()

const isLoading = ref(false)
const error = ref<string | null>(null)

async function handleSubmit(data: RegisterFormData) {
  isLoading.value = true
  error.value = null

  try {
    const result = await register(data.email, data.password, data.defaultAgeBracket)
    authStore.setTokens(result.tokens.accessToken, result.tokens.refreshToken)
    authStore.setUser(result.user)

    router.push('/')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Registration failed'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="max-w-md mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold text-gray-900 mb-2">Create Account</h1>
    <p class="text-gray-600 mb-6">
      Create an account to save your conversations and get unlimited questions.
    </p>

    <RegisterForm
      :is-loading="isLoading"
      :error="error"
      @submit="handleSubmit"
    />

    <p class="mt-6 text-center text-sm text-gray-600">
      Already have an account?
      <RouterLink to="/login" class="text-bracket-preteen hover:underline font-medium">
        Sign in
      </RouterLink>
    </p>
  </div>
</template>