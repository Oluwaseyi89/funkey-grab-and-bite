import type { Plugin } from '#app'

declare module '#app' {
  interface NuxtApp {
    $toast: {
      success: (message: string, options?: any) => void
      error: (message: string, options?: any) => void
      warning: (message: string, options?: any) => void
      info: (message: string, options?: any) => void
    }
  }
}

export {} // Important: This makes it a module