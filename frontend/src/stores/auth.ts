import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User } from '@/types';

const ACCESS_TOKEN_KEY = 'auth:accessToken';
const REFRESH_TOKEN_KEY = 'auth:refreshToken';
const USER_KEY = 'auth:user';

// Token expiry buffer (refresh 5 minutes before expiry)
const TOKEN_REFRESH_BUFFER_MS = 5 * 60 * 1000;

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const accessToken = ref<string | null>(null);
  const refreshToken = ref<string | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const tokenExpiresAt = ref<number | null>(null);

  const isAuthenticated = computed(() => !!accessToken.value && !!user.value);

  function loadFromStorage() {
    const storedAccessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);

    if (storedAccessToken) {
      accessToken.value = storedAccessToken;
    }
    if (storedRefreshToken) {
      refreshToken.value = storedRefreshToken;
    }
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser);
      } catch {
        localStorage.removeItem(USER_KEY);
      }
    }
  }

  function saveToStorage() {
    if (accessToken.value) {
      localStorage.setItem(ACCESS_TOKEN_KEY, accessToken.value);
    } else {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
    }
    if (refreshToken.value) {
      localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken.value);
    } else {
      localStorage.removeItem(REFRESH_TOKEN_KEY);
    }
    if (user.value) {
      localStorage.setItem(USER_KEY, JSON.stringify(user.value));
    } else {
      localStorage.removeItem(USER_KEY);
    }
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access;
    refreshToken.value = refresh;
    saveToStorage();
  }

  function setUser(userData: User) {
    user.value = userData;
    saveToStorage();
  }

  function clearAuth() {
    user.value = null;
    accessToken.value = null;
    refreshToken.value = null;
    tokenExpiresAt.value = null;
    error.value = null;
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  function shouldRefreshToken(): boolean {
    if (!accessToken.value || !refreshToken.value) {
      return false;
    }
    if (!tokenExpiresAt.value) {
      // If we don't know when it expires, don't proactively refresh
      return false;
    }
    return Date.now() >= tokenExpiresAt.value - TOKEN_REFRESH_BUFFER_MS;
  }

  function setError(err: string | null) {
    error.value = err;
  }

  function setLoading(loading: boolean) {
    isLoading.value = loading;
  }

  // Initialize from storage on store creation
  loadFromStorage();

  return {
    // State
    user,
    accessToken,
    refreshToken,
    isLoading,
    error,

    // Computed
    isAuthenticated,

    // Actions
    setTokens,
    setUser,
    clearAuth,
    shouldRefreshToken,
    setError,
    setLoading,
    loadFromStorage,
  };
});