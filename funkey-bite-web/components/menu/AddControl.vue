<template>
  <div :class="block ? 'w-full' : ''">
    <button
      v-if="quantity === 0"
      type="button"
      class="inline-flex min-h-[40px] items-center justify-center gap-1.5 rounded-full px-4 text-sm font-semibold shadow-soft transition active:scale-95 disabled:shadow-none"
      :class="[
        block ? 'w-full min-h-[44px]' : '',
        stock.orderable
          ? 'bg-white text-brand-700 ring-1 ring-ink-200 hover:bg-brand-600 hover:text-white hover:ring-brand-600 dark:bg-ink-900 dark:text-brand-300 dark:ring-ink-700'
          : 'cursor-not-allowed bg-ink-100 text-ink-500 dark:bg-ink-800 dark:text-ink-400',
      ]"
      :disabled="!stock.orderable"
      :aria-label="stock.orderable ? `Add ${item.name} to cart` : `${item.name}: ${stock.label}`"
      @click="cart.increase(item)"
    >
      <Plus v-if="stock.orderable" class="h-4 w-4" aria-hidden="true" />
      {{ stock.orderable ? 'Add' : 'Sold out' }}
    </button>
    <div v-else :class="block ? 'flex justify-between' : ''">
      <QuantityStepper
        :name="item.name"
        :quantity="quantity"
        :at-max="!cart.canAddMore(item)"
        size="sm"
        :class="block ? 'w-full justify-between' : 'shadow-soft'"
        @increase="cart.increase(item)"
        @decrease="cart.decrease(item.id)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Plus } from 'lucide-vue-next'
import QuantityStepper from '../ui/QuantityStepper.vue'
import { useCartStore } from '../../stores/cart'
import type { StockState } from '../../utils/menuStock'
import type { MenuItem } from '../../types/menu'

// "Add" on a dish, turning into - n + once it's in the cart, so customers can
// add and remove without leaving the menu.
const props = withDefaults(defineProps<{ item: MenuItem; stock: StockState; block?: boolean }>(), { block: false })
const cart = useCartStore()

// The saved cart is only known in the browser.
const mounted = ref(false)
onMounted(() => { mounted.value = true })
const quantity = computed(() => (mounted.value ? cart.itemCount(props.item.id) : 0))
</script>
