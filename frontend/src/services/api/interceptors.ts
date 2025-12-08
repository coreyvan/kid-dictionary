import type { Interceptor } from '@connectrpc/connect';
import { useAuthStore } from '@/stores/auth';

export const authInterceptor: Interceptor = (next) => async (req) => {
  // Get auth store - this works because Pinia is initialized before API calls
  const authStore = useAuthStore();

  // Add auth header if authenticated
  if (authStore.accessToken) {
    req.header.set('Authorization', `Bearer ${authStore.accessToken}`);
  }

  try {
    return await next(req);
  } catch (error) {
    // Re-throw the error - token refresh will be handled by the auth service
    throw error;
  }
};