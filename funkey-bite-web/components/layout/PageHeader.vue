<template>
  <header class="relative overflow-hidden border-b border-ink-200/60 bg-gradient-to-b from-brand-50/70 to-transparent dark:border-ink-800 dark:from-brand-900/10">
    <div aria-hidden="true" class="pointer-events-none absolute -right-16 -top-24 h-64 w-64 rounded-full bg-accent-200/40 blur-3xl dark:bg-accent-500/10" />
    <div aria-hidden="true" class="pointer-events-none absolute -left-20 top-10 h-48 w-48 rounded-full bg-leaf-200/40 blur-3xl dark:bg-leaf-500/10" />
    <div class="container-narrow relative py-7 md:py-12" :class="alignment === 'center' ? 'text-center' : ''">
      <p v-if="eyebrow" class="eyebrow mb-2">{{ eyebrow }}</p>
      <h1 class="font-display text-[1.75rem] font-bold leading-tight text-ink-900 dark:text-white md:text-5xl" :class="titleClass">
        <slot name="title">
          <template v-if="highlightText">{{ titleBefore }} <span class="text-gradient">{{ highlightText }}</span> {{ titleAfter }}</template>
          <template v-else>{{ title }}</template>
        </slot>
      </h1>
      <p
        v-if="subtitle || $slots.subtitle"
        class="mt-2 max-w-2xl text-base text-ink-600 dark:text-ink-300 md:mt-3 md:text-lg"
        :class="[alignment === 'center' ? 'mx-auto' : '', subtitleClass]"
      >
        <slot name="subtitle">{{ subtitle }}</slot>
      </p>
      <div v-if="$slots.default" class="mt-5">
        <slot />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
// The title band at the top of a page. Compact on phones so content starts
// near the top of the screen.
withDefaults(defineProps<{
  title?: string
  titleBefore?: string
  highlightText?: string
  titleAfter?: string
  subtitle?: string
  eyebrow?: string
  alignment?: 'center' | 'left'
  titleClass?: string
  subtitleClass?: string
  // Older options, accepted so existing pages keep working.
  variant?: string
  narrow?: boolean
  headerClass?: string
  solidColor?: string
  solidColorDark?: string
}>(), {
  title: '',
  titleBefore: '',
  highlightText: '',
  titleAfter: '',
  subtitle: '',
  eyebrow: '',
  alignment: 'left',
  titleClass: '',
  subtitleClass: '',
})
</script>
