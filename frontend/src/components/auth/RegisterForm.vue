<script setup lang="ts">
import { ref, computed } from 'vue'
import type { RegisterFormData, FormErrors } from '@/types'
import { AgeBracket } from '@/types'
import AgeBracketSelector from '@/components/chat/AgeBracketSelector.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'

const props = defineProps<{
  isLoading?: boolean
  error?: string | null
}>()

const emit = defineEmits<{
  submit: [data: RegisterFormData]
}>()

const formData = ref<RegisterFormData>({
  email: '',
  password: '',
  confirmPassword: '',
  defaultAgeBracket: AgeBracket.GROWING_MINDS,
})

const formErrors = ref<FormErrors>({})
const touched = ref<Record<string, boolean>>({})

function validateEmail(email: string): string | undefined {
  if (!email) return 'Email is required'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return 'Invalid email format'
  return undefined
}

function validatePassword(password: string): string | undefined {
  if (!password) return 'Password is required'
  if (password.length < 8) return 'Password must be at least 8 characters'
  return undefined
}

function validateConfirmPassword(confirmPassword: string): string | undefined {
  if (!confirmPassword) return 'Please confirm your password'
  if (confirmPassword !== formData.value.password) return 'Passwords do not match'
  return undefined
}

function validate(): boolean {
  formErrors.value = {
    email: validateEmail(formData.value.email),
    password: validatePassword(formData.value.password),
    confirmPassword: validateConfirmPassword(formData.value.confirmPassword),
  }
  return (
    !formErrors.value.email &&
    !formErrors.value.password &&
    !formErrors.value.confirmPassword
  )
}

function handleBlur(field: keyof typeof touched.value) {
  touched.value[field] = true
  if (field === 'email') {
    formErrors.value.email = validateEmail(formData.value.email)
  } else if (field === 'password') {
    formErrors.value.password = validatePassword(formData.value.password)
  } else if (field === 'confirmPassword') {
    formErrors.value.confirmPassword = validateConfirmPassword(formData.value.confirmPassword)
  }
}

function handleSubmit() {
  if (validate() && !props.isLoading) {
    emit('submit', formData.value)
  }
}

const isValid = computed(() => {
  return (
    formData.value.email &&
    formData.value.password &&
    formData.value.confirmPassword &&
    !formErrors.value.email &&
    !formErrors.value.password &&
    !formErrors.value.confirmPassword
  )
})
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <ErrorMessage
      v-if="error"
      :message="error"
      class="mb-4"
    />

    <div>
      <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
        Email
      </label>
      <input
        id="email"
        v-model="formData.email"
        type="email"
        autocomplete="email"
        :disabled="isLoading"
        class="w-full min-h-[44px] px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-bracket-preteen focus:border-transparent disabled:bg-gray-100"
        :class="touched.email && formErrors.email ? 'border-red-500' : 'border-gray-300'"
        @blur="handleBlur('email')"
      />
      <p v-if="touched.email && formErrors.email" class="mt-1 text-sm text-red-600">
        {{ formErrors.email }}
      </p>
    </div>

    <div>
      <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
        Password
      </label>
      <input
        id="password"
        v-model="formData.password"
        type="password"
        autocomplete="new-password"
        :disabled="isLoading"
        class="w-full min-h-[44px] px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-bracket-preteen focus:border-transparent disabled:bg-gray-100"
        :class="touched.password && formErrors.password ? 'border-red-500' : 'border-gray-300'"
        @blur="handleBlur('password')"
      />
      <p v-if="touched.password && formErrors.password" class="mt-1 text-sm text-red-600">
        {{ formErrors.password }}
      </p>
    </div>

    <div>
      <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">
        Confirm Password
      </label>
      <input
        id="confirmPassword"
        v-model="formData.confirmPassword"
        type="password"
        autocomplete="new-password"
        :disabled="isLoading"
        class="w-full min-h-[44px] px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-bracket-preteen focus:border-transparent disabled:bg-gray-100"
        :class="touched.confirmPassword && formErrors.confirmPassword ? 'border-red-500' : 'border-gray-300'"
        @blur="handleBlur('confirmPassword')"
      />
      <p v-if="touched.confirmPassword && formErrors.confirmPassword" class="mt-1 text-sm text-red-600">
        {{ formErrors.confirmPassword }}
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">
        Default Age Bracket
      </label>
      <p class="text-sm text-gray-500 mb-3">
        This will be used as the default for new conversations.
      </p>
      <AgeBracketSelector
        v-model="formData.defaultAgeBracket"
        :disabled="isLoading"
      />
    </div>

    <button
      type="submit"
      :disabled="isLoading || !isValid"
      class="w-full min-h-[44px] px-4 py-2 bg-bracket-preteen text-white rounded-lg font-medium hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-bracket-preteen focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
    >
      <LoadingSpinner v-if="isLoading" size="sm" />
      {{ isLoading ? 'Creating account...' : 'Create Account' }}
    </button>
  </form>
</template>