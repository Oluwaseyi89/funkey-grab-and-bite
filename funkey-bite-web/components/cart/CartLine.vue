<template>
  <li class="flex gap-3 py-4">
    <MenuImage :src="item.menuItem.imageUrl" :alt="item.menuItem.name" class="h-16 w-16 shrink-0 rounded-xl" />
    <div class="min-w-0 flex-1">
      <div class="flex items-start justify-between gap-3">
        <p class="font-semibold leading-snug text-ink-900 dark:text-white">{{ item.menuItem.name }}</p>
        <p class="shrink-0 font-semibold tabular-nums text-ink-900 dark:text-white">{{ naira(item.menuItem.price * item.quantity) }}</p>
      </div>
      <p class="text-xs text-ink-500 dark:text-ink-400">{{ naira(item.menuItem.price) }} each</p>

      <div v-if="editingNote" class="mt-2">
        <label :for="noteId" class="sr-only">Note for {{ item.menuItem.name }}</label>
        <input
          :id="noteId"
          v-model="note"
          maxlength="200"
          class="form-input min-h-[40px] py-2 text-sm"
          placeholder="e.g. no onions, extra pepper"
          @keydown.enter.prevent="saveNote"
          @blur="saveNote"
        />
      </div>
      <button
        v-else
        type="button"
        class="mt-1 inline-flex min-h-[32px] items-center gap-1 text-left text-xs font-medium text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-white"
        @click="startNote"
      >
        <PencilLine class="h-3.5 w-3.5" aria-hidden="true" />
        <span v-if="item.specialInstructions" class="line-clamp-1 italic">“{{ item.specialInstructions }}”</span>
        <span v-else>Add a note</span>
      </button>

      <div class="mt-2 flex items-center justify-between gap-3">
        <QuantityStepper
          :name="item.menuItem.name"
          :quantity="item.quantity"
          :at-max="!cart.canAddMore(item.menuItem)"
          @increase="cart.increase(item.menuItem)"
          @decrease="cart.decrease(item.menuItem.id)"
        />
        <button
          v-if="item.quantity > 1"
          type="button"
          class="min-h-[44px] px-2 text-sm font-medium text-ink-500 hover:text-brand-600 dark:text-ink-400"
          @click="cart.removeItem(item.menuItem.id)"
        >
          Remove<span class="sr-only"> {{ item.menuItem.name }}</span>
        </button>
      </div>
    </div>
  </li>
</template>

<script setup lang="ts">
import { nextTick, ref, useId } from 'vue'
import { PencilLine } from 'lucide-vue-next'
import MenuImage from '../menu/MenuImage.vue'
import QuantityStepper from '../ui/QuantityStepper.vue'
import { useCartStore } from '../../stores/cart'
import { naira } from '../../utils/format'
import type { CartItem } from '../../types/menu'

// One line in the cart: change the quantity (the minus becomes a bin at one),
// remove it, or add a note for the kitchen.
const props = defineProps<{ item: CartItem }>()
const cart = useCartStore()

const noteId = `note-${useId()}`
const editingNote = ref(false)
const note = ref('')

const startNote = async () => {
  note.value = props.item.specialInstructions ?? ''
  editingNote.value = true
  await nextTick()
  document.getElementById(noteId)?.focus()
}

const saveNote = () => {
  if (!editingNote.value) return
  cart.setInstructions(props.item.menuItem.id, note.value)
  editingNote.value = false
}
</script>
