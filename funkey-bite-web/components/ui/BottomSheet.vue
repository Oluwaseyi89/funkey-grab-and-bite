<template>
  <Teleport to="body">
    <Transition name="sheet-backdrop">
      <div v-if="open" class="fixed inset-0 z-[90] bg-ink-950/50 backdrop-blur-[2px]" aria-hidden="true" @click="close" />
    </Transition>
    <Transition :name="desktop === 'right' ? 'sheet-right' : 'sheet-center'">
      <div
        v-if="open"
        ref="panel"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="titleId"
        tabindex="-1"
        class="sheet fixed inset-x-0 bottom-0 z-[95] flex max-h-[92dvh] flex-col rounded-t-3xl bg-white shadow-lift outline-none dark:bg-ink-900"
        :class="desktop === 'right'
          ? 'md:inset-y-0 md:left-auto md:right-0 md:max-h-none md:w-[420px] md:rounded-none md:rounded-l-3xl'
          : 'md:inset-auto md:left-1/2 md:top-1/2 md:max-h-[85vh] md:w-[min(560px,92vw)] md:-translate-x-1/2 md:-translate-y-1/2 md:rounded-3xl'"
      >
        <div class="flex justify-center pt-2.5 md:hidden" aria-hidden="true">
          <span class="h-1.5 w-10 rounded-full bg-ink-200 dark:bg-ink-700" />
        </div>
        <header class="flex items-center justify-between gap-3 px-5 pb-2 pt-2 md:pt-5">
          <h2 :id="titleId" class="text-lg font-bold text-ink-900 dark:text-white" :class="{ 'sr-only': hideTitle }">{{ title }}</h2>
          <button type="button" class="icon-btn -mr-2 ml-auto" :aria-label="`Close ${title}`" @click="close">
            <X class="h-5 w-5" aria-hidden="true" />
          </button>
        </header>
        <div class="min-h-0 flex-1 overflow-y-auto overscroll-contain px-5 pb-5">
          <slot />
        </div>
        <footer v-if="$slots.footer" class="border-t border-ink-100 px-5 pb-safe pt-3 dark:border-ink-800">
          <div class="pb-3">
            <slot name="footer" />
          </div>
        </footer>
      </div>
    </Transition>
  </Teleport>
</template>

<script lang="ts">
// Open sheets across the page, newest last: only the top one answers the
// keyboard.
const stack: symbol[] = []
</script>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import { X } from 'lucide-vue-next'

// A dialog that slides up from the bottom on phones and sits in the centre
// (or on the right) on larger screens. It keeps keyboard focus inside while
// open, closes on Escape or a tap outside, and returns focus afterwards.
const props = withDefaults(defineProps<{
  open: boolean
  title: string
  hideTitle?: boolean
  desktop?: 'center' | 'right'
}>(), { hideTitle: false, desktop: 'center' })

const emit = defineEmits<{ close: [] }>()

const panel = ref<HTMLElement>()
const titleId = `sheet-${useId()}`
let returnFocusTo: HTMLElement | null = null
const self = Symbol('sheet')

const close = () => emit('close')

const focusable = () =>
  Array.from(
    panel.value?.querySelectorAll<HTMLElement>(
      'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
    ) ?? [],
  ).filter((el) => el.offsetParent !== null)

// Listens on the whole page, so Escape and Tab still work after the focused
// control disappears (e.g. a removed cart line).
const onKeydown = (event: KeyboardEvent) => {
  if (stack[stack.length - 1] !== self || !panel.value) return
  if (!panel.value.contains(document.activeElement)) {
    if (event.key === 'Tab') {
      event.preventDefault()
      ;(focusable()[0] ?? panel.value).focus()
      return
    }
  }
  if (event.key === 'Escape') {
    event.stopPropagation()
    close()
    return
  }
  if (event.key !== 'Tab') return
  const items = focusable()
  if (!items.length) return
  const first = items[0]!
  const last = items[items.length - 1]!
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

// Sheets can stack (an item opened from the cart), so the page stays locked
// until the last one closes.
const lockScroll = (locked: boolean) => {
  const body = document.body
  const count = Number(body.dataset.sheets ?? 0) + (locked ? 1 : -1)
  body.dataset.sheets = String(Math.max(count, 0))
  body.style.overflow = count > 0 ? 'hidden' : ''
}

watch(
  () => props.open,
  async (isOpen, wasOpen) => {
    if (!import.meta.client) return
    if (isOpen) {
      returnFocusTo = document.activeElement as HTMLElement | null
      stack.push(self)
      document.addEventListener('keydown', onKeydown)
      lockScroll(true)
      await nextTick()
      panel.value?.focus()
    } else if (wasOpen) {
      release()
      lockScroll(false)
      returnFocusTo?.focus?.()
    }
  },
  { immediate: true },
)

function release() {
  document.removeEventListener('keydown', onKeydown)
  const index = stack.indexOf(self)
  if (index >= 0) stack.splice(index, 1)
}

onBeforeUnmount(() => {
  if (props.open && import.meta.client) {
    release()
    lockScroll(false)
  }
})
</script>

<style scoped>
.sheet-backdrop-enter-active,
.sheet-backdrop-leave-active {
  transition: opacity 0.2s ease;
}
.sheet-backdrop-enter-from,
.sheet-backdrop-leave-to {
  opacity: 0;
}

.sheet-center-enter-active,
.sheet-center-leave-active,
.sheet-right-enter-active,
.sheet-right-leave-active {
  transition: transform 0.3s cubic-bezier(0.32, 0.72, 0, 1), opacity 0.2s ease;
}
.sheet-center-enter-from,
.sheet-center-leave-to,
.sheet-right-enter-from,
.sheet-right-leave-to {
  transform: translateY(100%);
}

@media (min-width: 768px) {
  .sheet-center-enter-from,
  .sheet-center-leave-to {
    transform: translate(-50%, -46%);
    opacity: 0;
  }
  .sheet-right-enter-from,
  .sheet-right-leave-to {
    transform: translateX(100%);
  }
}
</style>
