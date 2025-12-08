import { authClient } from './client'
import { getErrorMessage } from './errors'
import type { User } from '@/types'
import { AgeBracket } from '@/types'

export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

export interface LoginResult {
  user: User
  tokens: AuthTokens
}

export interface RegisterResult {
  user: User
  tokens: AuthTokens
}

export async function login(email: string, password: string): Promise<LoginResult> {
  try {
    const response = await authClient.login({ email, password })

    if (!response.user || !response.accessToken || !response.refreshToken) {
      throw new Error('Invalid response from server')
    }

    return {
      user: response.user as User,
      tokens: {
        accessToken: response.accessToken,
        refreshToken: response.refreshToken,
      },
    }
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}

export async function register(
  email: string,
  password: string,
  defaultAgeBracket: AgeBracket
): Promise<RegisterResult> {
  try {
    const response = await authClient.register({
      email,
      password,
      defaultAgeBracket,
    })

    if (!response.user || !response.accessToken || !response.refreshToken) {
      throw new Error('Invalid response from server')
    }

    return {
      user: response.user as User,
      tokens: {
        accessToken: response.accessToken,
        refreshToken: response.refreshToken,
      },
    }
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}

export async function refreshTokens(refreshToken: string): Promise<AuthTokens> {
  try {
    const response = await authClient.refreshToken({ refreshToken })

    if (!response.accessToken || !response.refreshToken) {
      throw new Error('Invalid response from server')
    }

    return {
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
    }
  } catch (error) {
    throw new Error(getErrorMessage(error))
  }
}
