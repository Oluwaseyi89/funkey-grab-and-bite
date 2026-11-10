<template>
  <BottomSheet :open="!!item" :title="item?.name ?? 'Item'" hide-title @close="$emit('close')">
    <template v-if="item">
      <MenuImage :src="item.imageUrl" :alt="item.name" eager class="-mx-5 -mt-1 aspect-[16/10] md:mx-0 md:rounded-2xl" />

      <div class="mt-4 flex items-start justify-between gap-4">
        <h2 class="font-display text-xl font-bold leading-tight text-ink-900 dark:text-white" aria-hidden="true">{{ item.name }}</h2>
        <p class="shrink-0 text-xl font-bold tabular-nums text-brand-700 dark:text-brand-300">{{ naira(item.price) }}</p>
      </div>
      <p v-if="item.description" class="mt-2 text-ink-600 dark:text-ink-300">{{ item.description }}</p>

      <ul class="mt-3 flex flex-wrap gap-2 text-xs">
        <li v-if="item.preparationTime" class="badge bg-ink-100 text-ink-700 dark:bg-ink-800 dark:text-ink-200">
          <Clock class="h-3.5 w-3.5" aria-hidden="true" /> Ready in about {{ item.preparationTime }} min
        </li>
        <li v-if="stock.tone === 'low'" class="badge bg-accent-100 text-accent-700 dark:bg-accent-700/30 dark:text-accent-200">{{ stock.label }}</li>
        <li v-for="tag in item.tags" :key="tag" class="badge bg-ink-100 capitalize text-ink-700 dark:bg-ink-800 dark:text-ink-200">{{ tag }}</li>
      </ul>

      <p v-if="item.isPreOrder" class="mt-4 rounded-xl bg-accent-50 px-4 py-3 text-sm text-accent-700 dark:bg-accent-700/20 dark:text-accent-200">
        This one is made to order, so please allow at least {{ item.preparationTime || 30 }} minutes.
      </p>

      <dl v-if="item.nutritionalInfo" class="mt-4 grid grid-cols-4 gap-2 text-center">
        <div v-for="fact in nutrition" :key="fact.label" class="surface-muted px-1 py-2">
          <dt class="text-[11px] text-ink-500 dark:text-ink-400">{{ fact.label }}</dt>
          <dd class="font-semibold text-ink-900 dark:text-white">{{ fact.value }}</dd>
        </div>
      </dl>

      <div class="mt-5">
        <label for="item-note" class="form-label">Note for the kitchen <span class="font-normal text-ink-500">(optional)</span></label>
        <textarea id="item-note" v-model="note" rows="2" maxlength="200" class="form-input resize-none" placeholder="e.g. no onions, extra pepper" />
      </div>
    </template>

    <template v-if="item" #footer>
      <div v-if="stock.orderable" class="flex items-center gap-3">
        <QuantityStepper :name="item.name" :quantity="quantity" :min="1" :at-max="atMax" @increase="quantity++" @decrease="quantity--" />
        <button type="button" class="btn-primary btn-lg flex-1" @click="add">
          Add {{ quantity }} · {{ naira(item.price * quantity) }}
        </button>
      </div>
      <p v-else class="rounded-xl bg-ink-100 px-4 py-3 text-center font-medium text-ink-700 dark:bg-ink-800 dark:text-ink-200">
        {{ stock.label }}
      </p>
    </template>
  </BottomSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Clock } from 'lucide-vue-next'
import BottomSheet from '../ui/BottomSheet.vue'
import QuantityStepper from '../ui/QuantityStepper.vue'
import MenuImage from './MenuImage.vue'
import { useCartStore } from '../../stores/cart'
import { menuStock } from '../../utils/menuStock'
import { naira } from '../../utils/format'
import type { MenuItem } from '../../types/menu'

const props = defineProps<{ item: MenuItem | null }>()
const emit = defineEmits<{ close: [] }>()
const cart = useCartStore()

const quantity = ref(1)
const note = ref('')

watch(() => props.item, (item) => {
  quantity.value = 1
  note.value = item ? cart.items.find((line) => line.menuItem.id === item.id)?.specialInstructions ?? '' : ''
})

const stock = computed(() => (props.item ? menuStock(props.item) : { orderable: false, maxQuantity: 0, label: '', tone: null }))
// Never let the picker go past what's left once the cart is counted.
const atMax = computed(() => {
  if (!props.item || stock.value.maxQuantity == null) return false
  return cart.itemCount(props.item.id) + quantity.value >= stock.value.maxQuantity
})

const nutrition = computed(() => {
  const info = props.item?.nutritionalInfo
  if (!info) return []
  return [
    { label: 'Calories', value: info.calories },
    { label: 'Protein', value: `${info.protein}g` },
    { label: 'Carbs', value: `${info.carbs}g` },
    { label: 'Fat', value: `${info.fat}g` },
  ]
})

const add = () => {
  if (!props.item) return
  const added = cart.addItem(props.item, quantity.value, note.value.trim())
  if (added > 0) emit('close')
}
</script>
