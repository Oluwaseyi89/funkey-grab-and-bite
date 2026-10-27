<template>
  <form class="space-y-8" novalidate @submit.prevent="submit">
    <!-- 1. Package or custom -->
    <fieldset class="space-y-3">
      <legend class="text-xl font-bold text-gray-900 dark:text-white mb-2">1. What would you like?</legend>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <label
          v-for="pkg in packages"
          :key="pkg.id"
          class="flex items-start gap-3 p-4 border-2 rounded-xl cursor-pointer transition-colors"
          :class="selection === pkg.id ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20' : 'border-gray-200 dark:border-slate-700 hover:border-brand-300'"
        >
          <input v-model="selection" type="radio" name="catering-choice" :value="pkg.id" class="mt-1" />
          <span>
            <span class="block font-bold text-gray-900 dark:text-white">{{ pkg.name }}</span>
            <span class="block text-sm text-gray-600 dark:text-gray-400">{{ naira(pkg.pricePerPerson) }} per guest · {{ guestRange(pkg) }}</span>
          </span>
        </label>
        <label
          class="flex items-start gap-3 p-4 border-2 rounded-xl cursor-pointer transition-colors"
          :class="selection === 'custom' ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20' : 'border-gray-200 dark:border-slate-700 hover:border-brand-300'"
        >
          <input v-model="selection" type="radio" name="catering-choice" value="custom" class="mt-1" />
          <span>
            <span class="block font-bold text-gray-900 dark:text-white">Custom catering</span>
            <span class="block text-sm text-gray-600 dark:text-gray-400">Tell us what you want and we'll quote for it</span>
          </span>
        </label>
      </div>
    </fieldset>

    <!-- 2. Event -->
    <fieldset class="space-y-4">
      <legend class="text-xl font-bold text-gray-900 dark:text-white mb-2">2. Your event</legend>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="cr-event-type" class="form-label">Type of event *</label>
          <select id="cr-event-type" v-model="form.eventType" class="form-input">
            <option v-for="type in EVENT_TYPES" :key="type.id" :value="type.id">{{ type.name }}</option>
          </select>
        </div>
        <div v-if="form.eventType === 'other'">
          <label for="cr-event-other" class="form-label">What kind of event? *</label>
          <input id="cr-event-other" v-model="form.eventTypeOther" maxlength="100" class="form-input" placeholder="e.g. House warming" />
        </div>
        <div :class="form.eventType === 'other' ? 'md:col-span-2' : ''">
          <label for="cr-event-name" class="form-label">Occasion name</label>
          <input id="cr-event-name" v-model="form.eventName" maxlength="200" class="form-input" placeholder="e.g. Mum's 60th birthday" />
        </div>
        <div>
          <label for="cr-date" class="form-label">Event date *</label>
          <input id="cr-date" v-model="form.eventDate" type="date" :min="earliestDate" class="form-input" />
          <p class="form-hint">At least {{ noticeDays }} days' notice{{ selectedPackage ? ` for ${selectedPackage.name}` : '' }}.</p>
        </div>
        <div>
          <label for="cr-time" class="form-label">Food served from *</label>
          <input id="cr-time" v-model="form.eventTime" type="time" class="form-input" />
        </div>
        <div>
          <label for="cr-duration" class="form-label">How long? *</label>
          <select id="cr-duration" v-model.number="form.durationHours" class="form-input">
            <option v-for="hours in DURATIONS" :key="hours" :value="hours">{{ hours }} hours</option>
          </select>
        </div>
        <div>
          <label for="cr-guests" class="form-label">Number of guests *</label>
          <input id="cr-guests" v-model.number="form.guestCount" type="number" :min="minGuests" :max="maxGuests" class="form-input" />
          <p class="form-hint">{{ minGuests }}–{{ maxGuests }} guests{{ selectedPackage ? ` for ${selectedPackage.name}` : '' }}.</p>
        </div>
      </div>
    </fieldset>

    <!-- 3. Venue -->
    <fieldset class="space-y-4">
      <legend class="text-xl font-bold text-gray-900 dark:text-white mb-2">3. Where is it?</legend>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="md:col-span-2">
          <label for="cr-venue-name" class="form-label">Venue name</label>
          <input id="cr-venue-name" v-model="form.venueName" maxlength="200" class="form-input" placeholder="e.g. Civic Centre, or 'Home'" />
        </div>
        <div class="md:col-span-2">
          <label for="cr-venue-address" class="form-label">Street address *</label>
          <input id="cr-venue-address" v-model="form.venueAddress" maxlength="500" autocomplete="street-address" class="form-input" placeholder="House number and street" />
        </div>
        <div>
          <label for="cr-venue-area" class="form-label">Area / city *</label>
          <input id="cr-venue-area" v-model="form.venueArea" maxlength="100" class="form-input" placeholder="e.g. Lekki" />
        </div>
        <div>
          <label for="cr-venue-state" class="form-label">State *</label>
          <input id="cr-venue-state" v-model="form.venueState" maxlength="100" class="form-input" placeholder="e.g. Lagos" />
        </div>
        <div class="md:col-span-2">
          <label for="cr-venue-notes" class="form-label">Access notes</label>
          <input id="cr-venue-notes" v-model="form.venueNotes" maxlength="500" class="form-input" placeholder="Landmark, gate, parking or setup time" />
        </div>
      </div>
    </fieldset>

    <!-- 4. Food & service -->
    <fieldset class="space-y-4">
      <legend class="text-xl font-bold text-gray-900 dark:text-white mb-2">4. Food & service</legend>

      <template v-if="selectedPackage">
        <p class="text-sm text-gray-700 dark:text-gray-300">
          Served as <strong>{{ serviceStyleName(selectedPackage.serviceStyle) }}</strong>.
          <span v-if="selectedPackage.inclusions.length">Includes {{ selectedPackage.inclusions.join(', ') }}.</span>
        </p>
        <div v-if="selectedPackage.dishChoices > 0">
          <p class="form-label">
            Choose {{ selectedPackage.dishChoices }} {{ selectedPackage.dishChoices === 1 ? 'dish' : 'dishes' }} *
            <span class="font-normal text-gray-500">({{ form.selectedItemIds.length }} chosen)</span>
          </p>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            <label
              v-for="item in selectedPackage.menuItems"
              :key="item.menuItemId"
              class="flex items-center gap-2 p-3 border rounded-lg"
              :class="item.isAvailable ? 'border-gray-200 dark:border-slate-700 cursor-pointer' : 'border-gray-100 dark:border-slate-800 opacity-50'"
            >
              <input
                type="checkbox"
                :value="item.menuItemId"
                v-model="form.selectedItemIds"
                :disabled="!item.isAvailable || (!form.selectedItemIds.includes(item.menuItemId) && form.selectedItemIds.length >= selectedPackage.dishChoices)"
              />
              <span class="text-gray-900 dark:text-white">{{ item.name }}</span>
              <span v-if="!item.isAvailable" class="ml-auto text-xs text-red-600">Unavailable</span>
            </label>
          </div>
        </div>
        <p v-else-if="selectedPackage.menuItems.length" class="text-sm text-gray-700 dark:text-gray-300">
          Menu: {{ selectedPackage.menuItems.map((item) => item.name).join(', ') }}
        </p>
      </template>

      <template v-else>
        <div>
          <p class="form-label">How should we serve? *</p>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <label
              v-for="style in SERVICE_STYLES"
              :key="style.id"
              class="p-3 border-2 rounded-xl cursor-pointer"
              :class="form.serviceStyle === style.id ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20' : 'border-gray-200 dark:border-slate-700'"
            >
              <input v-model="form.serviceStyle" type="radio" name="service-style" :value="style.id" class="sr-only" />
              <span class="block font-medium text-gray-900 dark:text-white">{{ style.name }}</span>
              <span class="block text-xs text-gray-600 dark:text-gray-400">{{ style.description }}</span>
            </label>
          </div>
        </div>
        <div>
          <label for="cr-menu-notes" class="form-label">What food would you like? *</label>
          <textarea id="cr-menu-notes" v-model="form.menuNotes" rows="3" maxlength="1000" class="form-input" placeholder="e.g. Jollof rice, fried rice, peppered chicken, small chops and drinks"></textarea>
        </div>
        <div v-if="menuItems.length">
          <p class="form-label">Or tick dishes from our menu <span class="font-normal text-gray-500">(optional, up to 20)</span></p>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 max-h-56 overflow-y-auto pr-1">
            <label v-for="item in menuItems" :key="item.id" class="flex items-center gap-2 p-2 border border-gray-200 dark:border-slate-700 rounded-lg cursor-pointer">
              <input type="checkbox" :value="Number(item.id)" v-model="form.selectedItemIds" :disabled="!form.selectedItemIds.includes(Number(item.id)) && form.selectedItemIds.length >= 20" />
              <span class="text-sm text-gray-900 dark:text-white">{{ item.name }}</span>
            </label>
          </div>
        </div>
      </template>

      <div>
        <p class="form-label">Dietary needs</p>
        <div class="flex flex-wrap gap-2">
          <label v-for="option in DIETARY_OPTIONS" :key="option.id" class="flex items-center gap-2 px-3 py-1.5 border border-gray-200 dark:border-slate-700 rounded-full cursor-pointer text-sm">
            <input type="checkbox" :value="option.id" v-model="form.dietaryRequirements" />
            {{ option.name }}
          </label>
        </div>
        <input v-model="form.dietaryNotes" maxlength="500" class="form-input mt-2" placeholder="Allergies or other needs (optional)" aria-label="Dietary notes" />
      </div>
    </fieldset>

    <!-- 5. Contact -->
    <fieldset class="space-y-4">
      <legend class="text-xl font-bold text-gray-900 dark:text-white mb-2">5. Your details</legend>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="cr-name" class="form-label">Full name *</label>
          <input id="cr-name" v-model="form.contactName" autocomplete="name" maxlength="200" class="form-input" />
        </div>
        <div>
          <label for="cr-phone" class="form-label">Phone number *</label>
          <input id="cr-phone" v-model="form.contactPhone" type="tel" autocomplete="tel" class="form-input" placeholder="0801 234 5678" />
        </div>
        <div>
          <label for="cr-email" class="form-label">Email address *</label>
          <input id="cr-email" v-model="form.contactEmail" type="email" autocomplete="email" maxlength="200" class="form-input" />
        </div>
        <div>
          <label for="cr-org" class="form-label">Organisation</label>
          <input id="cr-org" v-model="form.organization" autocomplete="organization" maxlength="200" class="form-input" placeholder="For company or church events" />
        </div>
        <div>
          <label for="cr-budget" class="form-label">Budget (₦)</label>
          <input id="cr-budget" v-model.number="form.budget" type="number" min="0" step="1000" class="form-input" placeholder="Optional" />
        </div>
        <div class="md:col-span-2">
          <label for="cr-special" class="form-label">Anything else?</label>
          <textarea id="cr-special" v-model="form.specialRequests" rows="2" maxlength="1000" class="form-input" placeholder="Theme, timing, equipment you need…"></textarea>
        </div>
      </div>
    </fieldset>

    <div class="p-4 rounded-xl bg-gray-50 dark:bg-slate-800 text-sm">
      <template v-if="selectedPackage">
        <p class="text-gray-600 dark:text-gray-400">Estimate for {{ form.guestCount || 0 }} guests</p>
        <p class="text-2xl font-bold text-brand-500">{{ naira(estimate) }}</p>
        <p class="text-gray-500 dark:text-gray-400">We'll confirm the final price in your quote.</p>
      </template>
      <p v-else class="text-gray-700 dark:text-gray-300">We'll review your request and send you a quote.</p>
    </div>

    <p v-if="error" class="text-red-600 text-sm" role="alert">{{ error }}</p>

    <button type="submit" :disabled="isSubmitting" class="w-full btn-primary py-4 text-lg font-bold disabled:opacity-50 dark:text-white">
      <Loader2 v-if="isSubmitting" class="w-5 h-5 animate-spin inline mr-2" />
      {{ isSubmitting ? 'Sending request…' : 'Request a quote' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import type { CateringPackage, CateringRequest, CateringRequestPayload, ServiceStyle } from '../../types/catering'
import type { MenuItem } from '../../types/menu'
import { CUSTOM_NOTICE_DAYS, DIETARY_OPTIONS, EVENT_TYPES, SERVICE_STYLES } from '../../types/catering'
import { useApi } from '../../utils/api'
import { useBusinessSettings } from '../../composables/useBusinessSettings'

const DURATIONS = [2, 3, 4, 5, 6, 8, 10, 12]

const props = defineProps<{
  packages: CateringPackage[]
  // A package id, or 'custom'.
  initialSelection?: number | 'custom'
}>()

const emit = defineEmits<{
  submitted: [request: CateringRequest]
}>()

const api = useApi()
const { settings } = useBusinessSettings()

const selection = ref<number | 'custom'>(props.initialSelection ?? props.packages[0]?.id ?? 'custom')
watch(() => props.initialSelection, (value) => {
  if (value !== undefined) selection.value = value
})

const emptyForm = () => ({
  eventType: 'party',
  eventTypeOther: '',
  eventName: '',
  eventDate: '',
  eventTime: '',
  durationHours: 4,
  guestCount: 50,
  venueName: '',
  venueAddress: '',
  venueArea: '',
  venueState: 'Lagos',
  venueNotes: '',
  serviceStyle: '' as ServiceStyle | '',
  selectedItemIds: [] as number[],
  menuNotes: '',
  dietaryRequirements: [] as string[],
  dietaryNotes: '',
  contactName: '',
  contactPhone: '',
  contactEmail: '',
  organization: '',
  budget: undefined as number | undefined,
  specialRequests: '',
})
const form = ref(emptyForm())
const error = ref('')
const isSubmitting = ref(false)
const menuItems = ref<MenuItem[]>([])

const selectedPackage = computed(() =>
  typeof selection.value === 'number' ? props.packages.find((pkg) => pkg.id === selection.value) : undefined,
)

// Dish choices belong to one package; changing the choice starts them over.
watch(selection, () => {
  form.value.selectedItemIds = []
  if (selectedPackage.value && form.value.guestCount < selectedPackage.value.minGuests) {
    form.value.guestCount = selectedPackage.value.minGuests
  }
})

const minGuests = computed(() => selectedPackage.value?.minGuests ?? 1)
const maxGuests = computed(() => selectedPackage.value?.maxGuests ?? 1000)
const noticeDays = computed(() => selectedPackage.value?.minNoticeDays ?? CUSTOM_NOTICE_DAYS)
const estimate = computed(() => (selectedPackage.value?.pricePerPerson ?? 0) * (form.value.guestCount || 0))

// "Today" in the business timezone, so the earliest date matches the API.
const earliestDate = computed(() => {
  const timeZone = settings.value?.timezone || 'Africa/Lagos'
  const today = new Date(new Intl.DateTimeFormat('en-CA', { timeZone }).format(new Date()) + 'T00:00:00')
  today.setDate(today.getDate() + noticeDays.value)
  return [today.getFullYear(), String(today.getMonth() + 1).padStart(2, '0'), String(today.getDate()).padStart(2, '0')].join('-')
})

const naira = (amount: number) => `₦${amount.toLocaleString('en-NG')}`
const guestRange = (pkg: CateringPackage) => pkg.maxGuests ? `${pkg.minGuests}–${pkg.maxGuests} guests` : `${pkg.minGuests}+ guests`
const serviceStyleName = (style: ServiceStyle) => SERVICE_STYLES.find((s) => s.id === style)?.name ?? style
const cleanPhone = (value: string) => value.trim().replace(/(?!^\+)[^\d]/g, '')
const optional = (value: string) => value.trim() || undefined

// Mirrors the API's rules so most mistakes are caught before sending.
function validate(): string {
  const f = form.value
  if (f.eventType === 'other' && f.eventTypeOther.trim().length < 2) return 'Tell us what kind of event it is.'
  if (!f.eventDate) return 'Choose the event date.'
  if (f.eventDate < earliestDate.value) return `We need at least ${noticeDays.value} days' notice; please choose ${earliestDate.value} or later.`
  if (!/^([01]\d|2[0-3]):[0-5]\d$/.test(f.eventTime)) return 'Choose the time food should be served.'
  if (!Number.isInteger(f.guestCount) || f.guestCount < minGuests.value || f.guestCount > maxGuests.value) {
    return `Enter between ${minGuests.value} and ${maxGuests.value} guests.`
  }
  if (f.venueAddress.trim().length < 5) return 'Enter the venue street address.'
  if (f.venueArea.trim().length < 2) return 'Enter the area or city of the venue.'
  if (f.venueState.trim().length < 2) return 'Enter the state of the venue.'
  if (selectedPackage.value) {
    const choices = selectedPackage.value.dishChoices
    if (choices > 0 && f.selectedItemIds.length !== choices) return `Choose ${choices} ${choices === 1 ? 'dish' : 'dishes'} for the ${selectedPackage.value.name} package.`
  } else {
    if (!f.serviceStyle) return 'Choose how you would like the food served.'
    if (f.selectedItemIds.length === 0 && f.menuNotes.trim().length < 10) return 'Tell us what food you would like.'
  }
  if (f.contactName.trim().length < 2) return 'Enter your full name.'
  if (!/^\+?\d{7,15}$/.test(cleanPhone(f.contactPhone))) return 'Enter a valid phone number, e.g. 0801 234 5678.'
  if (!/^\S+@\S+\.\S+$/.test(f.contactEmail.trim())) return 'Enter a valid email address.'
  if (f.budget != null && (Number.isNaN(f.budget) || f.budget < 0)) return 'Budget cannot be negative.'
  return ''
}

async function submit() {
  error.value = validate()
  if (error.value) return

  const f = form.value
  const payload: CateringRequestPayload = {
    requestType: selectedPackage.value ? 'package' : 'custom',
    packageId: selectedPackage.value?.id,
    contactName: f.contactName.trim(),
    contactPhone: f.contactPhone.trim(),
    contactEmail: f.contactEmail.trim(),
    organization: optional(f.organization),
    eventType: f.eventType,
    eventTypeOther: f.eventType === 'other' ? optional(f.eventTypeOther) : undefined,
    eventName: optional(f.eventName),
    eventDate: f.eventDate,
    eventTime: f.eventTime,
    durationHours: f.durationHours,
    guestCount: f.guestCount,
    venueName: optional(f.venueName),
    venueAddress: f.venueAddress.trim(),
    venueArea: f.venueArea.trim(),
    venueState: f.venueState.trim(),
    venueNotes: optional(f.venueNotes),
    serviceStyle: selectedPackage.value ? undefined : (f.serviceStyle || undefined),
    selectedItemIds: f.selectedItemIds,
    menuNotes: optional(f.menuNotes),
    dietaryRequirements: f.dietaryRequirements,
    dietaryNotes: optional(f.dietaryNotes),
    budget: f.budget || undefined,
    specialRequests: optional(f.specialRequests),
  }

  isSubmitting.value = true
  try {
    const created = await api.submitCateringRequest(payload)
    form.value = emptyForm()
    emit('submitted', created)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Your request could not be sent. Please try again.'
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  try {
    const items = await api.getMenuItems()
    menuItems.value = items.filter((item) => item.isAvailable && /^\d+$/.test(item.id))
  } catch {
    // Dish preferences are optional; the notes field is enough.
  }
})
</script>

<style scoped>
.form-label {
  @apply block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2;
}
.form-input {
  @apply w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400;
}
.form-hint {
  @apply mt-1 text-xs text-gray-500 dark:text-gray-400;
}
</style>
