<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="translate-y-6 opacity-0"
    leave-active-class="transition duration-200 ease-in"
    leave-to-class="translate-y-6 opacity-0"
  >
    <div v-if="visible" class="fixed inset-x-0 z-40 px-3 md:hidden" style="bottom: calc(4.5rem + env(safe-area-inset-bottom))">
      <button
        type="button"
        class="flex w-full items-center gap-3 rounded-2xl bg-brand-600 px-4 py-3 text-left text-white shadow-glow transition active:scale-[0.98]"
        @click="cart.openCart()"
      >
        <span :key="cart.bump" class="flex h-9 min-w-9 items-center justify-center rounded-xl bg-white/20 px-2 text-sm font-bold tabular-nums animate-pop">
          {{ cart.totalItems }}
        </span>
        <span class="flex-1 font-semibold">View cart</span>
        <span class="font-bold tabular-nums">{{ naira(cart.subtotal) }}</span>
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useCartStore } from '../../stores/cart'
import { naira } from '../../utils/format'

// The phone's "View cart" bar above the tab bar, shown while browsing food
// with something in the cart. Checkout has its own button, so it hides there.
const cart = useCartStore()
const route = useRoute()
const mounted = ref(false)
onMounted(() => { mounted.value = true })

const browsingFood = computed(() => route.path === '/' || route.path.startsWith('/menu') || route.path.startsWith('/promotion'))
const visible = computed(() => mounted.value && !cart.isEmpty && !cart.isOpen && browsingFood.value)
</script>
