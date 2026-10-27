<template>
  <div class="container-narrow">
    <PageHeader
      title-before="Catering for"
      highlight-text="Every"
      title-after="Occasion"
      subtitle="Parties, weddings, naming ceremonies, church and office events: choose one of our packages or tell us exactly what you want."
      variant="solid"
      solid-color="brand-50"
      solid-color-dark="brand-600"
      :narrow="true"
      alignment="center"
      header-class="dark:from-brand-500 dark:to-accent-500"
      title-class="text-gray-900 dark:text-white"
      subtitle-class="text-gray-700 dark:text-white/90"
    >
      <div class="flex flex-col sm:flex-row gap-4 justify-center mt-8">
        <button
          @click="scrollToForm"
          class="text-gray-900 dark:text-white font-bold py-3 px-8 rounded-xl transition-colors hover:scale-105 shadow-lg hover:shadow-xl"
        >
          Request a Quote
        </button>
        <a
          v-if="phoneHref"
          :href="phoneHref"
          class="bg-white/80 dark:bg-white/20 hover:bg-white dark:hover:bg-white/30 backdrop-blur-sm font-medium py-3 px-8 rounded-xl transition-colors hover:scale-105 flex items-center justify-center gap-2 border border-brand-200 dark:border-white/30 text-brand-700 dark:text-white"
        >
          <Phone class="w-5 h-5" />
          {{ settings?.phoneNumber }}
        </a>
      </div>
    </PageHeader>

    <section class="section-padding mx-4 px-4">
      <h2 class="text-3xl font-bold text-gray-900 dark:text-white text-center mb-10">How it works</h2>
      <ol class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <li v-for="(step, index) in steps" :key="step.title" class="text-center">
          <div class="w-12 h-12 rounded-full bg-brand-500 text-white font-bold flex items-center justify-center mx-auto mb-4">{{ index + 1 }}</div>
          <h3 class="font-bold text-gray-900 dark:text-white mb-2">{{ step.title }}</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ step.description }}</p>
        </li>
      </ol>
    </section>

    <section class="section-padding mx-4 px-4">
      <h2 class="text-3xl font-bold text-gray-900 dark:text-white text-center mb-10">
        Our <span class="text-gradient">Catering</span> Packages
      </h2>

      <p v-if="packagesStatus === 'loading'" class="text-center text-gray-500">Loading packages…</p>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <article
          v-for="pkg in packages"
          :key="pkg.id"
          class="flex flex-col bg-white dark:bg-slate-800 rounded-2xl p-6 border-2 border-gray-200 dark:border-slate-700"
        >
          <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ pkg.name }}</h3>
          <p class="text-2xl font-bold text-brand-500 mt-2">{{ naira(pkg.pricePerPerson) }}<span class="text-base text-gray-500"> / guest</span></p>
          <p v-if="pkg.description" class="text-gray-600 dark:text-gray-400 mt-2">{{ pkg.description }}</p>
          <ul class="mt-4 space-y-2 text-sm text-gray-700 dark:text-gray-300 flex-1">
            <li class="flex gap-2"><Users class="w-4 h-4 text-brand-500 shrink-0 mt-0.5" />{{ pkg.maxGuests ? `${pkg.minGuests}–${pkg.maxGuests}` : `${pkg.minGuests}+` }} guests</li>
            <li class="flex gap-2"><UtensilsCrossed class="w-4 h-4 text-brand-500 shrink-0 mt-0.5" />{{ serviceStyleName(pkg.serviceStyle) }}</li>
            <li v-if="pkg.dishChoices > 0" class="flex gap-2"><CheckCircle class="w-4 h-4 text-green-500 shrink-0 mt-0.5" />Choose {{ pkg.dishChoices }} of {{ pkg.menuItems.length }} dishes</li>
            <li v-for="inclusion in pkg.inclusions" :key="inclusion" class="flex gap-2"><CheckCircle class="w-4 h-4 text-green-500 shrink-0 mt-0.5" />{{ inclusion }}</li>
            <li class="flex gap-2"><CalendarClock class="w-4 h-4 text-brand-500 shrink-0 mt-0.5" />Book at least {{ pkg.minNoticeDays }} days ahead</li>
          </ul>
          <button class="btn-primary w-full mt-6 dark:text-white" @click="choose(pkg.id)">Choose {{ pkg.name }}</button>
        </article>

        <article class="flex flex-col bg-white dark:bg-slate-800 rounded-2xl p-6 border-2 border-dashed border-brand-300 dark:border-brand-700">
          <h3 class="text-xl font-bold text-gray-900 dark:text-white">Custom catering</h3>
          <p class="text-gray-600 dark:text-gray-400 mt-2 flex-1">
            Have your own menu in mind? Tell us the food, the guests and how you want it served, and we'll quote for it.
          </p>
          <button class="btn-secondary w-full mt-6" @click="choose('custom')">Make a custom request</button>
        </article>
      </div>
      <p v-if="packagesStatus === 'unavailable'" class="text-center text-sm text-gray-500 mt-4">
        Our packages couldn't be loaded right now, but you can still send a custom request.
      </p>
    </section>

    <section ref="formSection" class="section-padding mx-4 px-4 bg-white dark:bg-slate-900 rounded-2xl">
      <div class="max-w-3xl mx-auto">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Request a catering quote</h2>
        <p class="text-gray-600 dark:text-gray-400 mb-8">
          Already sent one?
          <NuxtLink to="/catering/request" class="text-brand-500 hover:text-brand-600 font-medium">Check your request</NuxtLink>
        </p>
        <CateringRequestForm
          v-if="packagesStatus !== 'loading'"
          :packages="packages"
          :initial-selection="selection"
          @submitted="handleSubmitted"
        />
      </div>
    </section>

    <section class="section-padding mx-4 px-4 bg-gray-50 dark:bg-slate-800 my-8 rounded-2xl">
      <CateringFAQ />
    </section>

    <div v-if="submitted" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" role="presentation">
      <div class="bg-white dark:bg-slate-800 rounded-2xl max-w-md w-full p-8 text-center" role="dialog" aria-modal="true" aria-labelledby="catering-sent-title">
        <div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
          <CheckCircle class="w-10 h-10 text-green-600 dark:text-green-400" />
        </div>
        <h2 id="catering-sent-title" class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Request sent!</h2>
        <p class="text-gray-600 dark:text-gray-400 mb-4">
          Your reference is
        </p>
        <p class="text-3xl font-mono font-bold text-brand-500 mb-4">{{ submitted.reference }}</p>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          We've emailed you a copy. We'll review your request and contact you with a quote. Keep your reference to check its status.
        </p>
        <div class="space-y-3">
          <NuxtLink to="/catering/request" class="block w-full btn-primary dark:text-white" @click="submitted = null">Check request status</NuxtLink>
          <button class="w-full btn-secondary" @click="submitted = null">Done</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { CheckCircle, Phone, Users, UtensilsCrossed, CalendarClock } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import CateringRequestForm from '../../components/catering/CateringRequestForm.vue'
import CateringFAQ from '../../components/catering/CateringFAQ.vue'
import type { CateringPackage, CateringRequest, ServiceStyle } from '../../types/catering'
import { SERVICE_STYLES } from '../../types/catering'
import { useApi } from '../../utils/api'
import { useBusinessSettings } from '../../composables/useBusinessSettings'

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
  { title: 'Get a quote', description: 'We review your request and contact you with a price.' },
  { title: 'Confirm', description: 'Once you accept, your booking is confirmed and we cater your event.' },
]

const naira = (amount: number) => `₦${amount.toLocaleString('en-NG')}`
const serviceStyleName = (style: ServiceStyle) => SERVICE_STYLES.find((s) => s.id === style)?.name ?? style

const scrollToForm = () => formSection.value?.scrollIntoView({ behavior: 'smooth' })

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
