<template>
  <AccountCard title-before="Sign" highlight="In" subtitle="See your orders and check out faster">
    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="login-phone" class="form-label">Phone number</label>
        <input id="login-phone" v-model="phone" type="tel" required autocomplete="tel" class="form-input" placeholder="0801 234 5678" />
      </div>
      <div>
        <label for="login-password" class="form-label">Password</label>
        <input id="login-password" v-model="password" type="password" required autocomplete="current-password" class="form-input" />
      </div>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
      <button type="submit" class="btn-primary w-full dark:text-white" :disabled="busy">
        <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />Sign in
      </button>
    </form>
    <div class="flex justify-between text-sm">
      <NuxtLink to="/account/forgot-password" class="text-brand-600 hover:underline">Forgot password?</NuxtLink>
      <NuxtLink :to="{ path: '/account/register', query: next ? { next } : {} }" class="text-brand-600 hover:underline">Create an account</NuxtLink>
    </div>
    <p class="text-xs text-gray-500 dark:text-gray-400">An account is optional; you can always check out as a guest.</p>
  </AccountCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Loader2 } from 'lucide-vue-next'
import AccountCard from '../../components/account/AccountCard.vue'
import { useApi } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const next = typeof route.query.next === 'string' && route.query.next.startsWith('/') ? route.query.next : ''
const phone = ref('')
const password = ref('')
const error = ref('')
const busy = ref(false)

const submit = async () => {
  error.value = ''
  busy.value = true
  try {
    auth.setSession(await useApi().login(phone.value.trim(), password.value))
    await router.push(next || '/account')
  } catch {
    error.value = 'That phone number and password don’t match an account.'
  } finally {
    busy.value = false
  }
}
</script>
