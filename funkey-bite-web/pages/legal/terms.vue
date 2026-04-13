<template>
  <div>
    <PageHeader
      title-before="Terms of"
      highlight-text="Service"
      subtitle="The terms that govern use of the Funkey Grab & Bite website and ordering experience."
      variant="simple"
      header-class="pb-12"
    />

    <section class="section-padding mx-5 my-5 md:mx-8 md:my-8 px-3 py-3 md:px-8 md:py-8">
      <div class="container-narrow max-w-4xl space-y-8 text-gray-700 dark:text-gray-300">
        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-gray-200 dark:border-slate-700 p-6 md:p-8 space-y-6">
          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Orders and Availability</h2>
            <p>
              All orders are subject to availability, confirmation, and operational capacity. Menu items, prices, offers, and delivery coverage may change without prior notice.
            </p>
          </div>

          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Payments and Promotions</h2>
            <p>
              Promotions are subject to their stated conditions and validity periods. Misuse of promotional offers or fraudulent activity may result in cancellation of orders or restricted access to services.
            </p>
          </div>

          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Cancellations and Changes</h2>
            <p>
              Order changes or cancellations may be limited once preparation has started. Catering bookings may require additional notice depending on event size, staffing, and procurement commitments.
            </p>
          </div>

          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Use of the Website</h2>
            <p>
              You agree to use the website only for lawful purposes and not to interfere with platform security, availability, or other users’ experience.
            </p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import PageHeader from '../../components/layout/PageHeader.vue'
</script>