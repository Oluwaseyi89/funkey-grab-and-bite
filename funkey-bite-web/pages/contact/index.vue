<template>
  <div>
    <PageHeader eyebrow="Contact" title="We’re here to help" subtitle="Questions about an order, catering or anything else? Call, email or come and see us.">
      <StoreStatus />
    </PageHeader>

    <div class="container-narrow grid gap-6 py-6 md:grid-cols-[1.2fr_1fr] md:py-10">
      <div class="space-y-4">
        <!-- Ways to reach us -->
        <ul class="grid gap-3 sm:grid-cols-3 md:grid-cols-1 lg:grid-cols-3">
          <li v-if="phoneHref">
            <a :href="phoneHref" class="card flex h-full items-center gap-4 p-4 transition hover:shadow-lift sm:flex-col sm:items-start md:flex-row lg:flex-col">
              <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-brand-50 text-brand-700 dark:bg-brand-900/30 dark:text-brand-300"><Phone class="h-6 w-6" aria-hidden="true" /></span>
              <span class="min-w-0">
                <span class="block font-semibold text-ink-900 dark:text-white">Call us</span>
                <span class="block truncate text-sm text-ink-500 dark:text-ink-400">{{ settings?.phoneNumber }}</span>
              </span>
            </a>
          </li>
          <li v-if="settings?.email">
            <a :href="`mailto:${settings.email}`" class="card flex h-full items-center gap-4 p-4 transition hover:shadow-lift sm:flex-col sm:items-start md:flex-row lg:flex-col">
              <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-accent-50 text-accent-700 dark:bg-accent-700/20 dark:text-accent-300"><Mail class="h-6 w-6" aria-hidden="true" /></span>
              <span class="min-w-0">
                <span class="block font-semibold text-ink-900 dark:text-white">Email us</span>
                <span class="block truncate text-sm text-ink-500 dark:text-ink-400">{{ settings.email }}</span>
              </span>
            </a>
          </li>
          <li v-if="mapsUrl">
            <a :href="mapsUrl" target="_blank" rel="noopener noreferrer" class="card flex h-full items-center gap-4 p-4 transition hover:shadow-lift sm:flex-col sm:items-start md:flex-row lg:flex-col">
              <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-leaf-50 text-leaf-700 dark:bg-leaf-800/30 dark:text-leaf-300"><MapPin class="h-6 w-6" aria-hidden="true" /></span>
              <span class="min-w-0">
                <span class="block font-semibold text-ink-900 dark:text-white">Visit us</span>
                <span class="block text-sm text-ink-500 dark:text-ink-400">{{ fullAddress }}</span>
                <span class="sr-only">(opens Google Maps)</span>
              </span>
            </a>
          </li>
        </ul>
        <p v-if="status === 'unavailable'" class="card p-5 text-sm text-ink-600 dark:text-ink-300" role="alert">
          Our contact details didn’t load. Please refresh the page to try again.
        </p>

        <!-- Common questions -->
        <section class="card p-5" aria-labelledby="faq-title">
          <h2 id="faq-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Common questions</h2>
          <div class="mt-2 divide-y divide-ink-100 dark:divide-ink-800">
            <details v-for="faq in faqs" :key="faq.question" class="group py-1">
              <summary class="flex min-h-[48px] cursor-pointer list-none items-center justify-between gap-4 font-medium text-ink-900 dark:text-white">
                {{ faq.question }}
                <ChevronDown class="h-5 w-5 shrink-0 text-ink-400 transition group-open:rotate-180" aria-hidden="true" />
              </summary>
              <p class="pb-3 text-sm text-ink-600 dark:text-ink-300">
                {{ faq.answer }}
                <NuxtLink v-if="faq.link" :to="faq.link.to" class="link">{{ faq.link.label }}</NuxtLink>
              </p>
            </details>
          </div>
        </section>
      </div>

      <!-- Opening hours -->
      <section class="card h-fit p-5" aria-labelledby="hours-title">
        <h2 id="hours-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Opening hours</h2>
        <dl v-if="status === 'ready'" class="mt-3 divide-y divide-ink-100 dark:divide-ink-800">
          <div v-for="group in hoursGroups" :key="group.label" class="flex justify-between gap-4 py-2.5 text-sm" :class="group.days.includes(today) ? 'font-semibold text-ink-900 dark:text-white' : 'text-ink-600 dark:text-ink-300'">
            <dt>{{ group.label }}<span v-if="group.days.includes(today)" class="ml-2 badge bg-brand-50 text-brand-700 dark:bg-brand-900/30 dark:text-brand-300">Today</span></dt>
            <dd class="text-right">{{ group.hours }}</dd>
          </div>
        </dl>
        <div v-else-if="status === 'loading'" class="mt-3 space-y-2">
          <div v-for="i in 4" :key="i" class="skeleton h-5" />
        </div>
        <p v-else class="mt-3 text-sm text-ink-600 dark:text-ink-300">Our hours didn’t load. Please call us to check.</p>
        <p v-if="timezoneLabel" class="mt-3 text-xs text-ink-500 dark:text-ink-400">Times are in {{ timezoneLabel }}.</p>
      </section>
    </div>

    <!-- The team, once the admin has added people. -->
    <section v-if="team.length" class="container-narrow pb-12" aria-labelledby="team-title">
      <h2 id="team-title" class="section-title">Meet the team</h2>
      <ul class="mt-5 grid grid-cols-2 gap-4 md:grid-cols-4">
        <li v-for="member in team" :key="member.id" class="card p-5 text-center">
          <img v-if="member.photoUrl" :src="member.photoUrl" :alt="member.name" class="mx-auto h-20 w-20 rounded-full object-cover" loading="lazy" />
          <span v-else class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-br from-brand-500 to-accent-500 font-display text-2xl font-bold text-white" aria-hidden="true">{{ initials(member.name) }}</span>
          <h3 class="mt-3 font-display font-semibold text-ink-900 dark:text-white">{{ member.name }}</h3>
          <p class="text-sm text-brand-600 dark:text-brand-400">{{ member.role }}</p>
          <p v-if="member.bio" class="mt-1 text-xs text-ink-500 dark:text-ink-400">{{ member.bio }}</p>
        </li>
      </ul>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ChevronDown, Mail, MapPin, Phone } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import StoreStatus from '../../components/layout/StoreStatus.vue'
import { useBusinessSettings } from '../../composables/useBusinessSettings'
import { useApi } from '../../utils/api'
import { usePageSeo } from '../../composables/usePageSeo'
import type { TeamMember } from '../../types/team'

usePageSeo({ title: 'Contact', description: 'Call, email or visit Funkey Grab & Bite. Opening hours and answers to common questions.', path: '/contact' })

const { settings, status, phoneHref, mapsUrl, fullAddress, hoursGroups, today, timezoneLabel, load } = useBusinessSettings()
const team = ref<TeamMember[]>([])

onMounted(async () => {
  load()
  team.value = await useApi().getTeamMembers()
})

const initials = (name: string) =>
  name.split(/\s+/).filter(Boolean).slice(0, 2).map((part) => part[0]?.toUpperCase() ?? '').join('')

const faqs = computed(() => [
  {
    question: 'When can I order?',
    answer: 'Orders can be placed while we’re open; the hours are on this page. Your cart is saved if you visit while we’re closed.',
  },
  {
    question: 'How do I check on my order?',
    answer: 'Open Orders and enter your order number and phone, or sign in to see all your orders.',
    link: { to: '/order/track', label: 'Track an order' },
  },
  {
    question: 'Something’s wrong with my order. What do I do?',
    answer: 'Call us as soon as you can with your order number and we’ll sort it out.',
  },
  {
    question: 'Do you cater small events?',
    answer: 'Yes. Each package shows how many guests it serves, and a custom request can be any size.',
    link: { to: '/catering', label: 'See catering' },
  },
  {
    question: 'I have an allergy. Can you help?',
    answer: 'Add a note to the item or your order, and call us before ordering so we can tell you what’s in each dish.',
  },
])
</script>
