import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { AuthSession, Customer } from '../types/account'
import { useApi } from '../utils/api'

// The signed-in customer, kept across visits. Signing in is optional: guests
// can still check out.
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<Customer | null>(null)
  const isSignedIn = computed(() => !!token.value)

  function setSession(session: AuthSession) {
    token.value = session.accessToken
    user.value = session.user
    useApi().setAuthToken(session.accessToken)
  }

  function setUser(customer: Customer) {
    user.value = customer
  }

  function signOut() {
    token.value = null
    user.value = null
    useApi().setAuthToken(null)
  }

  return { token, user, isSignedIn, setSession, setUser, signOut }
}, {
  persist: true,
})
