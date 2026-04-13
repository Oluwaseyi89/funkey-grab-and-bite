<template>
  <div>
    <PageHeader
      title-before="Privacy"
      highlight-text="Policy"
      subtitle="How Funkey Grab & Bite collects, uses, and protects your information."
      variant="simple"
      header-class="pb-12"
    />

    <section class="section-padding mx-5 my-5 md:mx-8 md:my-8 px-3 py-3 md:px-8 md:py-8">
      <div class="container-narrow max-w-4xl space-y-8 text-gray-700 dark:text-gray-300">
        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-gray-200 dark:border-slate-700 p-6 md:p-8 space-y-6">
          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Information We Collect</h2>
            <p>
              We may collect details you provide when placing orders, requesting catering, or contacting us, including your name, phone number, email address, delivery information, and order preferences.
            </p>
          </div>

          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">How We Use Your Information</h2>
            <p>
              We use your information to process orders, coordinate catering requests, provide customer support, send order updates, and improve our services and customer experience.
            </p>
          </div>

          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Sharing and Security</h2>
            <p>
              We do not sell your personal information. Data may be shared only with service providers needed to operate the platform, fulfill orders, or deliver communications. We apply reasonable technical and operational safeguards to protect stored information.
            </p>
          </div>

          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Contact</h2>
            <p>
              If you have questions about this policy or your data, contact us through the contact page or our official support channels.
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