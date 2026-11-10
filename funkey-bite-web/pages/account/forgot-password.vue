<template>
  <AccountCard title-before="Reset your" highlight="Password" subtitle="We'll email you a code">
    <form v-if="!codeSent" class="space-y-4" @submit.prevent="requestCode">
      <div>
        <label for="reset-email" class="form-label">Email on your account</label>
        <input id="reset-email" v-model="email" type="email" required autocomplete="email" class="form-input" />
      </div>
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <button type="submit" class="btn-primary btn-lg w-full" :disabled="busy">
        <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />Send code
      </button>
    </form>

    <form v-else-if="!done" class="space-y-4" @submit.prevent="reset">
      <p class="text-sm text-ink-700 dark:text-ink-200">If an account uses {{ email }}, we've sent it a 6-digit code.</p>
      <div>
        <label for="reset-code" class="form-label">Code</label>
        <input id="reset-code" v-model="code" inputmode="numeric" autocomplete="one-time-code" maxlength="6" required class="form-input tracking-widest text-center text-lg" />
      </div>
      <div>
        <label for="reset-password" class="form-label">New password</label>
        <input id="reset-password" v-model="password" type="password" required minlength="8" autocomplete="new-password" class="form-input" />
        <p class="text-xs text-ink-500 dark:text-ink-400 mt-1">At least 8 characters, with upper and lower case letters and a number.</p>
      </div>
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <button type="submit" class="btn-primary btn-lg w-full" :disabled="busy">
        <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />Change password
      </button>
    </form>

    <div v-else class="space-y-4">
      <p class="text-ink-700 dark:text-ink-200" role="status">Your password has been changed.</p>
      <NuxtLink to="/account/login" class="btn-primary inline-block dark:text-white">Sign in</NuxtLink>
    </div>
  </AccountCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import AccountCard from '../../components/account/AccountCard.vue'
import { useApi } from '../../utils/api'

const email = ref('')
const code = ref('')
const password = ref('')
const codeSent = ref(false)
const done = ref(false)
const error = ref('')
const busy = ref(false)

const message = (err: unknown) => (err instanceof Error ? err.message : 'Something went wrong. Please try again.')

const requestCode = async () => {
  error.value = ''
  busy.value = true
  try {
    await useApi().forgotPassword(email.value.trim())
    codeSent.value = true
  } catch (err) {
    error.value = message(err)
  } finally {
    busy.value = false
  }
}

const reset = async () => {
  error.value = ''
  busy.value = true
  try {
    await useApi().resetPassword(email.value.trim(), code.value.trim(), password.value)
    done.value = true
  } catch (err) {
    error.value = message(err)
  } finally {
    busy.value = false
  }
}
</script>
