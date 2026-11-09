<template>
  <AccountCard title-before="Create an" highlight="Account" subtitle="Track your orders and check out faster">
    <form v-if="!sentTo" class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="reg-name" class="form-label">Full name</label>
        <input id="reg-name" v-model="details.fullName" required minlength="2" autocomplete="name" class="form-input" />
      </div>
      <div>
        <label for="reg-phone" class="form-label">Phone number</label>
        <input id="reg-phone" v-model="details.phone" type="tel" required autocomplete="tel" class="form-input" placeholder="0801 234 5678" />
      </div>
      <div>
        <label for="reg-email" class="form-label">Email</label>
        <input id="reg-email" v-model="details.email" type="email" required autocomplete="email" class="form-input" />
      </div>
      <div>
        <label for="reg-password" class="form-label">Password</label>
        <input id="reg-password" v-model="details.password" type="password" required minlength="8" autocomplete="new-password" class="form-input" />
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">At least 8 characters, with upper and lower case letters and a number.</p>
      </div>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
      <button type="submit" class="btn-primary w-full dark:text-white" :disabled="busy">
        <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />Create account
      </button>
    </form>

    <form v-else class="space-y-4" @submit.prevent="confirm">
      <p class="text-sm text-gray-700 dark:text-gray-300">
        This number has orders from before. To link them to your account, enter the 6-digit code we sent to <strong>{{ sentTo }}</strong>.
      </p>
      <div>
        <label for="reg-code" class="form-label">Code</label>
        <input id="reg-code" v-model="code" inputmode="numeric" autocomplete="one-time-code" maxlength="6" required class="form-input tracking-widest text-center text-lg" />
      </div>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
      <button type="submit" class="btn-primary w-full dark:text-white" :disabled="busy">
        <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />Confirm and create account
      </button>
      <button type="button" class="text-sm text-brand-600 hover:underline" :disabled="busy" @click="submit">Send a new code</button>
    </form>

    <p class="text-sm">Already have an account? <NuxtLink to="/account/login" class="text-brand-600 hover:underline">Sign in</NuxtLink></p>
  </AccountCard>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Loader2 } from 'lucide-vue-next'
import AccountCard from '../../components/account/AccountCard.vue'
import { useApi } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'
import type { RegistrationDetails } from '../../types/account'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const next = typeof route.query.next === 'string' && route.query.next.startsWith('/') ? route.query.next : ''
const details = reactive<RegistrationDetails>({ fullName: '', phone: '', email: '', password: '' })
const sentTo = ref('')
const code = ref('')
const error = ref('')
const busy = ref(false)

const message = (err: unknown) => (err instanceof Error ? err.message : 'Something went wrong. Please try again.')

const submit = async () => {
  error.value = ''
  busy.value = true
  try {
    const result = await useApi().register({ ...details, phone: details.phone.trim(), email: details.email.trim() })
    if (result.kind === 'verify') {
      sentTo.value = result.sentTo
      code.value = ''
    } else {
      auth.setSession(result.session)
      await router.push(next || '/account')
    }
  } catch (err) {
    error.value = message(err)
  } finally {
    busy.value = false
  }
}

const confirm = async () => {
  error.value = ''
  busy.value = true
  try {
    auth.setSession(await useApi().confirmRegistration({ ...details, phone: details.phone.trim(), email: details.email.trim() }, code.value.trim()))
    await router.push(next || '/account')
  } catch (err) {
    error.value = message(err)
  } finally {
    busy.value = false
  }
}
</script>
