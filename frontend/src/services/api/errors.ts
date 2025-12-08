import { ConnectError, Code } from '@connectrpc/connect';

export function getErrorMessage(error: unknown): string {
  if (error instanceof ConnectError) {
    switch (error.code) {
      case Code.Unauthenticated:
        return 'Please log in to continue.';
      case Code.PermissionDenied:
        return 'You do not have permission to perform this action.';
      case Code.NotFound:
        return 'The requested item was not found.';
      case Code.AlreadyExists:
        return 'This email is already registered.';
      case Code.InvalidArgument:
        return error.message || 'Invalid input provided.';
      case Code.ResourceExhausted:
        return 'Rate limit reached. Please register for unlimited access.';
      case Code.Unavailable:
        return 'Service temporarily unavailable. Please try again.';
      case Code.DeadlineExceeded:
        return 'Request timed out. Please try again.';
      case Code.Internal:
        return 'An unexpected error occurred. Please try again.';
      default:
        return 'An unexpected error occurred. Please try again.';
    }
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'An unexpected error occurred.';
}

export function isRateLimitError(error: unknown): boolean {
  return error instanceof ConnectError && error.code === Code.ResourceExhausted;
}

export function isAuthError(error: unknown): boolean {
  return error instanceof ConnectError && error.code === Code.Unauthenticated;
}