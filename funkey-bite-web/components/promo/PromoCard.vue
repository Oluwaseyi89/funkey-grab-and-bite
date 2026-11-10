<template>
  <article class="relative flex h-full flex-col overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 via-brand-500 to-accent-500 p-5 text-white shadow-glow">
    <div aria-hidden="true" class="pointer-events-none absolute -right-10 -top-10 h-40 w-40 rounded-full bg-white/15" />
    <div aria-hidden="true" class="pointer-events-none absolute -bottom-16 -left-8 h-44 w-44 rounded-full bg-black/10" />

    <div class="relative flex items-start justify-between gap-3">
      <p class="font-display text-3xl font-extrabold leading-none tracking-tight">{{ headline }}</p>
      <span class="badge bg-white/20 text-white backdrop-blur">{{ ends }}</span>
    </div>
    <h3 class="relative mt-3 font-display text-lg font-semibold leading-snug">{{ promotion.title }}</h3>
    <p v-if="promotion.description" class="relative mt-1 line-clamp-2 text-sm text-white/85">{{ promotion.description }}</p>
    <p v-if="conditions" class="relative mt-2 text-xs text-white/75">{{ conditions }}</p>

    <div class="relative mt-auto flex items-center gap-2 pt-5">
      <button
        type="button"
        class="inline-flex min-h-[44px] flex-1 items-center justify-between gap-2 rounded-xl border-2 border-dashed border-white/60 bg-white/10 px-3 font-mono text-sm font-bold tracking-wider transition hover:bg-white/20"
        :aria-label="copied ? `Code ${promotion.code} copied` : `Copy code ${promotion.code}`"
        @click="copy"
      >
        {{ promotion.code }}
        <component :is="copied ? Check : Copy" class="h-4 w-4" aria-hidden="true" />
      </button>
      <button type="button" class="inline-flex min-h-[44px] items-center rounded-xl bg-white px-4 text-sm font-bold text-brand-700 transition hover:bg-brand-50 active:scale-95" @click="use">
        {{ inUse ? 'Applied' : 'Use code' }}
      </button>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Check, Copy } from 'lucide-vue-next'
import { useCartStore } from '../../stores/cart'
import { formatDay, naira } from '../../utils/format'
import type { Promotion } from '../../types/menu'

// A deal with its code. "Use code" saves it for checkout, where the server
// checks it and works out the discount.
const props = defineProps<{ promotion: Promotion }>()
const cart = useCartStore()
const router = useRouter()
const copied = ref(false)

const headline = computed(() => {
  const p = props.promotion
  if (p.promotionType === 'percentage') return `${p.discountValue}% off`
  if (p.promotionType === 'fixed') return `${naira(p.discountValue)} off`
  return 'Buy 1, get 1'
})

const ends = computed(() => {
  const days = Math.ceil((new Date(props.promotion.validUntil).getTime() - Date.now()) / 86_400_000)
  if (days <= 1) return 'Ends today'
  if (days <= 7) return `${days} days left`
  return `Until ${formatDay(props.promotion.validUntil)}`
})

const conditions = computed(() => {
  const parts: string[] = []
  if (props.promotion.minOrderAmount) parts.push(`On orders over ${naira(props.promotion.minOrderAmount)}`)
  if (props.promotion.maxDiscount && props.promotion.promotionType === 'percentage') parts.push(`up to ${naira(props.promotion.maxDiscount)} off`)
  return parts.join(', ')
})

const inUse = computed(() => cart.promoCode.toUpperCase() === props.promotion.code.toUpperCase())

const copy = async () => {
  try {
    await navigator.clipboard.writeText(props.promotion.code)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // Clipboard blocked: the code is on screen to type.
  }
}

const use = () => {
  cart.promoCode = props.promotion.code
  router.push(cart.isEmpty ? '/menu' : '/order')
}
</script>
