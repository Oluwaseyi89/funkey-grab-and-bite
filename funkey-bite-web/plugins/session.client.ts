import { defineNuxtPlugin } from 'nuxt/app'
import { useAuthStore } from '../stores/auth'
import { useApi } from '../utils/api'

// Runs after the saved stores are restored (plugins load in name order), so a
// returning customer's requests carry their token straight away.
export default defineNuxtPlugin(() => {
  const auth = useAuthStore()
  useApi().setAuthToken(auth.token)
})
