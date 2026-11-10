<template>
  <section class="container-narrow py-6 md:py-12" aria-labelledby="catering-cta-title">
    <div class="relative overflow-hidden rounded-[2rem] bg-ink-900 px-6 py-8 text-white dark:bg-ink-800 md:grid md:grid-cols-[1.3fr_1fr] md:items-center md:gap-10 md:px-12 md:py-14">
      <div aria-hidden="true" class="pointer-events-none absolute -right-20 -top-24 h-72 w-72 rounded-full bg-brand-600/50 blur-3xl" />
      <div aria-hidden="true" class="pointer-events-none absolute -bottom-24 left-10 h-64 w-64 rounded-full bg-leaf-500/30 blur-3xl" />

      <div class="relative">
        <p class="text-xs font-semibold uppercase tracking-[0.14em] text-accent-300">Catering</p>
        <h2 id="catering-cta-title" class="mt-2 font-display text-3xl font-bold leading-tight md:text-4xl">
          Feeding a crowd? <span class="text-accent-300">We’ve got it.</span>
        </h2>
        <p class="mt-3 max-w-md text-white/80">
          Birthdays, weddings, naming ceremonies, church and office events. Pick a package or tell us what you want, and we’ll send you a quote.
        </p>
        <div class="mt-6 flex flex-col gap-3 sm:flex-row">
          <NuxtLink to="/catering" class="btn btn-lg bg-white text-ink-900 hover:bg-ink-100">
            Get a quote <ArrowRight class="h-5 w-5" aria-hidden="true" />
          </NuxtLink>
          <NuxtLink to="/catering/request" class="btn btn-lg border border-white/30 text-white hover:bg-white/10">
            Check my booking
          </NuxtLink>
        </div>
      </div>

      <ol class="relative mt-8 space-y-3 md:mt-0">
        <li v-for="(step, index) in steps" :key="step" class="flex items-center gap-3 rounded-2xl bg-white/10 px-4 py-3 backdrop-blur">
          <span class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-accent-400 font-display text-sm font-bold text-ink-900">{{ index + 1 }}</span>
          <span class="text-sm font-medium">{{ step }}</span>
        </li>
      </ol>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next'

const steps = [
  'Tell us the date, guests and food',
  'Get your quote by email',
  'Accept and pay the deposit online',
  'We cater, you enjoy',
]
</script>
