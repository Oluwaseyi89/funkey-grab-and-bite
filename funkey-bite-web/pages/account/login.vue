<template>
  <AccountCard title-before="Sign" highlight="In" subtitle="See your orders and check out faster">
    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="login-phone" class="form-label">Phone number</label>
        <input id="login-phone" v-model="phone" type="tel" required autocomplete="tel" class="form-input" placeholder="0801 234 5678" />
      </div>
      <div>
        <label for="login-password" class="form-label">Password</label>
        <div class="relative">
          <input id="login-password" v-model="password" :type="showPassword ? 'text' : 'password'" required autocomplete="current-password" class="form-input pr-12" />
          <button type="button" class="icon-btn absolute right-1 top-1/2 -translate-y-1/2" :aria-label="showPassword ? 'Hide password' : 'Show password'" :aria-pressed="showPassword" @click="showPassword = !showPassword">
            <component :is="showPassword ? EyeOff : Eye" class="h-5 w-5" aria-hidden="true" />
          </button>
        </div>
      </div>
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <button type="submit" class="btn-primary btn-lg w-full" :disabled="busy">
        <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />Sign in
      </button>
    </form>
    <div class="flex flex-wrap justify-between gap-2 text-sm">
      <NuxtLink to="/account/forgot-password" class="link">Forgot password?</NuxtLink>
      <NuxtLink :to="{ path: '/account/register', query: next ? { next } : {} }" class="link">Create an account</NuxtLink>
    </div>
    <p class="text-center text-xs text-ink-500 dark:text-ink-400">An account is optional; you can always check out as a guest.</p>
  </AccountCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import AccountCard from '../../components/account/AccountCard.vue'
import { ApiRequestError, useApi } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const next = typeof route.query.next === 'string' && route.query.next.startsWith('/') ? route.query.next : ''
const phone = ref('')
const password = ref('')
const showPassword = ref(false)
const error = ref('')
const busy = ref(false)

const submit = async () => {
  error.value = ''
  busy.value = true
  try {
    auth.setSession(await useApi().login(phone.value.trim(), password.value))
    await router.push(next || '/account')
  } catch (err) {
    // A wrong password is a 400/401; anything else (offline, server down)
    // shouldn't be blamed on what the customer typed.
    error.value = err instanceof ApiRequestError && (err.status === 400 || err.status === 401)
      ? 'That phone number and password don’t match an account.'
      : 'We couldn’t sign you in right now. Check your connection and try again.'
  } finally {
    busy.value = false
  }
}
</script>
