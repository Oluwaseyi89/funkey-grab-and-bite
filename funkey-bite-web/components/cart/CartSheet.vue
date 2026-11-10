<template>
  <BottomSheet :open="cart.isOpen" title="Your cart" desktop="right" @close="cart.closeCart()">
    <div v-if="cart.isEmpty" class="flex flex-col items-center py-10 text-center">
      <div class="mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-brand-50 dark:bg-brand-900/30">
        <ShoppingBag class="h-9 w-9 text-brand-500" aria-hidden="true" />
      </div>
      <p class="text-lg font-semibold text-ink-900 dark:text-white">Your cart is empty</p>
      <p class="mt-1 max-w-xs text-sm text-ink-500 dark:text-ink-400">Add something tasty from the menu and it’ll show up here.</p>
      <NuxtLink to="/menu" class="btn-primary mt-6" @click="cart.closeCart()">Browse the menu</NuxtLink>
    </div>

    <template v-else>
      <p v-if="closedNotice" class="mb-2 rounded-xl bg-accent-50 px-4 py-3 text-sm text-accent-700 dark:bg-accent-700/20 dark:text-accent-200">
        {{ closedNotice }}
      </p>
      <ul class="divide-y divide-ink-100 dark:divide-ink-800">
        <CartLine v-for="item in cart.items" :key="item.menuItem.id" :item="item" />
      </ul>
      <button type="button" class="mt-2 text-sm font-medium text-ink-500 hover:text-brand-600 dark:text-ink-400" @click="confirmClear">
        {{ clearing ? 'Tap again to empty your cart' : 'Empty cart' }}
      </button>
    </template>

    <template v-if="!cart.isEmpty" #footer>
      <div class="mb-3 flex items-baseline justify-between">
        <span class="text-sm text-ink-600 dark:text-ink-300">Subtotal · {{ plural(cart.totalItems, 'item') }}</span>
        <span class="text-xl font-bold tabular-nums text-ink-900 dark:text-white">{{ naira(cart.subtotal) }}</span>
      </div>
      <NuxtLink to="/order" class="btn-primary btn-lg w-full" @click="cart.closeCart()">
        Go to checkout
        <ArrowRight class="h-5 w-5" aria-hidden="true" />
      </NuxtLink>
      <p class="mt-2 text-center text-xs text-ink-500 dark:text-ink-400">Delivery fee and tax are added at checkout.</p>
    </template>
  </BottomSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowRight, ShoppingBag } from 'lucide-vue-next'
import BottomSheet from '../ui/BottomSheet.vue'
import CartLine from './CartLine.vue'
import { useCartStore } from '../../stores/cart'
import { useStoreStatus } from '../../composables/useStoreStatus'
import { naira, plural } from '../../utils/format'

const cart = useCartStore()
const route = useRoute()
const { isOpen, ready, detail } = useStoreStatus()

const closedNotice = computed(() =>
  ready.value && !isOpen.value
    ? `We’re closed right now${detail.value ? ` (${detail.value})` : ''}. Your cart is saved, so you can order as soon as we open.`
    : '',
)

// Emptying the cart takes two taps, so it can't happen by accident.
const clearing = ref(false)
let resetTimer: ReturnType<typeof setTimeout> | undefined
const confirmClear = () => {
  if (clearing.value) {
    cart.clearCart()
    clearing.value = false
    return
  }
  clearing.value = true
  clearTimeout(resetTimer)
  resetTimer = setTimeout(() => { clearing.value = false }, 3000)
}

watch(() => route.fullPath, () => cart.closeCart())
</script>
