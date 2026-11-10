<template>
  <div>
    <PageHeader
      eyebrow="Catering"
      title-before="Catering for"
      highlight-text="every occasion"
      subtitle="Parties, weddings, naming ceremonies, church and office events. Choose a package or tell us exactly what you want."
    >
      <div class="flex flex-col gap-3 sm:flex-row">
        <button type="button" class="btn-primary btn-lg" @click="scrollToForm">Request a quote</button>
        <a v-if="phoneHref" :href="phoneHref" class="btn-secondary btn-lg">
          <Phone class="h-5 w-5" aria-hidden="true" /> Call {{ settings?.phoneNumber }}
        </a>
      </div>
      <p class="mt-4 text-sm text-ink-600 dark:text-ink-300">
        Already sent a request? <NuxtLink to="/catering/request" class="link">Check your booking</NuxtLink>
      </p>
    </PageHeader>

    <section class="container-narrow py-8 md:py-12" aria-labelledby="how-title">
      <h2 id="how-title" class="section-title">How it works</h2>
      <ol class="mt-5 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        <li v-for="(step, index) in steps" :key="step.title" class="card flex gap-4 p-5 sm:flex-col">
          <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-brand-600 font-display font-bold text-white">{{ index + 1 }}</span>
          <div>
            <h3 class="font-display font-semibold text-ink-900 dark:text-white">{{ step.title }}</h3>
            <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">{{ step.description }}</p>
          </div>
        </li>
      </ol>
    </section>

    <section class="container-narrow pb-8 md:pb-12" aria-labelledby="packages-title">
      <h2 id="packages-title" class="section-title">Packages</h2>
      <div v-if="packagesStatus === 'loading'" class="mt-5 grid gap-4 md:grid-cols-2 lg:grid-cols-3" aria-busy="true">
        <div v-for="i in 3" :key="i" class="skeleton h-72 rounded-2xl" />
      </div>
      <div v-else class="mt-5 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <article v-for="pkg in packages" :key="pkg.id" class="card flex flex-col p-6">
          <h3 class="font-display text-xl font-bold text-ink-900 dark:text-white">{{ pkg.name }}</h3>
          <p class="mt-2 font-display text-3xl font-bold text-brand-700 dark:text-brand-300">{{ naira(pkg.pricePerPerson) }}<span class="text-base font-medium text-ink-500"> / guest</span></p>
          <p v-if="pkg.description" class="mt-2 text-ink-600 dark:text-ink-300">{{ pkg.description }}</p>
          <ul class="mt-4 flex-1 space-y-2 text-sm text-ink-700 dark:text-ink-200">
            <li class="flex gap-2"><Users class="mt-0.5 h-4 w-4 shrink-0 text-brand-500" aria-hidden="true" />{{ pkg.maxGuests ? `${pkg.minGuests}–${pkg.maxGuests}` : `${pkg.minGuests}+` }} guests</li>
            <li class="flex gap-2"><UtensilsCrossed class="mt-0.5 h-4 w-4 shrink-0 text-brand-500" aria-hidden="true" />{{ serviceStyleName(pkg.serviceStyle) }}</li>
            <li v-if="pkg.dishChoices > 0" class="flex gap-2"><Check class="mt-0.5 h-4 w-4 shrink-0 text-leaf-600" aria-hidden="true" />Choose {{ pkg.dishChoices }} of {{ pkg.menuItems.length }} dishes</li>
            <li v-for="inclusion in pkg.inclusions" :key="inclusion" class="flex gap-2"><Check class="mt-0.5 h-4 w-4 shrink-0 text-leaf-600" aria-hidden="true" />{{ inclusion }}</li>
            <li class="flex gap-2"><CalendarClock class="mt-0.5 h-4 w-4 shrink-0 text-brand-500" aria-hidden="true" />Book at least {{ pkg.minNoticeDays }} days ahead</li>
          </ul>
          <button type="button" class="btn-primary mt-6 w-full" @click="choose(pkg.id)">Choose {{ pkg.name }}</button>
        </article>

        <article class="flex flex-col rounded-2xl border-2 border-dashed border-brand-200 bg-brand-50/40 p-6 dark:border-brand-800 dark:bg-brand-900/10">
          <Sparkles class="h-7 w-7 text-brand-500" aria-hidden="true" />
          <h3 class="mt-3 font-display text-xl font-bold text-ink-900 dark:text-white">Something custom</h3>
          <p class="mt-2 flex-1 text-ink-600 dark:text-ink-300">Have your own menu in mind? Tell us the food, the guests and how you want it served, and we’ll quote for it.</p>
          <button type="button" class="btn-secondary mt-6 w-full" @click="choose('custom')">Make a custom request</button>
        </article>
      </div>
      <p v-if="packagesStatus === 'unavailable'" class="mt-4 text-sm text-ink-500 dark:text-ink-400">
        Our packages couldn’t be loaded right now, but you can still send a custom request.
      </p>
    </section>

    <section ref="formSection" class="container-narrow pb-8 md:pb-12" aria-labelledby="form-title">
      <div class="card mx-auto max-w-3xl p-5 md:p-8">
        <h2 id="form-title" class="font-display text-2xl font-bold text-ink-900 dark:text-white">Request a quote</h2>
        <p class="mt-1 text-ink-600 dark:text-ink-300">It takes about three minutes. We’ll email you a quote.</p>
        <div class="mt-6">
          <CateringRequestForm v-if="packagesStatus !== 'loading'" :packages="packages" :initial-selection="selection" @submitted="handleSubmitted" />
        </div>
      </div>
    </section>

    <section class="container-narrow pb-12">
      <div class="card p-5 md:p-8"><CateringFAQ /></div>
    </section>

    <BottomSheet :open="!!submitted" title="Request sent" @close="submitted = null">
      <div v-if="submitted" class="py-2 text-center">
        <span class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-leaf-100 text-leaf-700 dark:bg-leaf-800/40 dark:text-leaf-200">
          <CheckCircle class="h-8 w-8" aria-hidden="true" />
        </span>
        <p class="mt-4 text-ink-600 dark:text-ink-300">Your reference is</p>
        <p class="mt-1 font-mono text-3xl font-bold tracking-wider text-brand-700 dark:text-brand-300">{{ submitted.reference }}</p>
        <p class="mt-3 text-sm text-ink-600 dark:text-ink-300">We’ve emailed you a copy and will contact you with a quote. Keep the reference to check on your booking.</p>
      </div>
      <template #footer>
        <div class="grid gap-2">
          <NuxtLink to="/catering/request" class="btn-primary btn-lg w-full" @click="submitted = null">Check my booking</NuxtLink>
          <button type="button" class="btn-ghost w-full" @click="submitted = null">Done</button>
        </div>
      </template>
    </BottomSheet>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { CalendarClock, Check, CheckCircle, Phone, Sparkles, Users, UtensilsCrossed } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import BottomSheet from '../../components/ui/BottomSheet.vue'
import CateringRequestForm from '../../components/catering/CateringRequestForm.vue'
import CateringFAQ from '../../components/catering/CateringFAQ.vue'
import type { CateringPackage, CateringRequest, ServiceStyle } from '../../types/catering'
import { SERVICE_STYLES } from '../../types/catering'
import { useApi } from '../../utils/api'
import { naira } from '../../utils/format'
import { useBusinessSettings } from '../../composables/useBusinessSettings'
import { usePageSeo } from '../../composables/usePageSeo'

usePageSeo({
  title: 'Catering',
  description: 'Catering for parties, weddings, naming ceremonies, church and office events. Choose a package or request a custom quote.',
  path: '/catering',
})

const api = useApi()
const { settings, phoneHref, load: loadBusinessSettings } = useBusinessSettings()

const formSection = ref<HTMLElement>()
const packages = ref<CateringPackage[]>([])
const packagesStatus = ref<'loading' | 'ready' | 'unavailable'>('loading')
const selection = ref<number | 'custom' | undefined>(undefined)
const submitted = ref<CateringRequest | null>(null)

const steps = [
  { title: 'Choose', description: 'Pick a package, or describe the food you want in a custom request.' },
  { title: 'Tell us about it', description: 'Event date, venue, number of guests and any dietary needs.' },
  { title: 'Get a quote', description: 'We review your request and email you a price.' },
  { title: 'Confirm', description: 'Accept the quote and pay the deposit online to lock in your date.' },
]

const serviceStyleName = (style: ServiceStyle) => SERVICE_STYLES.find((s) => s.id === style)?.name ?? style

const scrollToForm = () => formSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })

const choose = (choice: number | 'custom') => {
  selection.value = choice
  scrollToForm()
}

const handleSubmitted = (request: CateringRequest) => {
  submitted.value = request
}

onMounted(async () => {
  loadBusinessSettings()
  const loaded = await api.getCateringPackages()
  packages.value = loaded ?? []
  packagesStatus.value = loaded ? 'ready' : 'unavailable'
})
</script>
