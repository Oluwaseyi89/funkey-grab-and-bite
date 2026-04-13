<template>
  <header 
    class="section-padding bg-transparent dark:from-slate-900 dark:to-slate-800 px-8 md:px-12 py-8 md:py-12"
    :class="[
      headerClass,
      variant === 'solid' ? `bg-${solidColor || 'amber-50'} dark:bg-${solidColorDark || 'amber-900/20'}` : '',
      variant === 'image' && backgroundImage ? `bg-cover bg-center` : ''
    ]"
    :style="variant === 'image' && backgroundImage ? `background-image: url('${backgroundImage}')` : ''"
  >
    <div 
      :class="[
        'text-center',
        narrow ? 'container-narrow' : 'container',
        alignment === 'left' ? 'text-left' : '',
        alignment === 'right' ? 'text-right' : ''
      ]"
    >

      <h1 
        class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4"
        :class="titleClass"
      >
        <slot name="title">

          <template v-if="highlightText">
            {{ titleBefore }}
            <span class="text-gradient">{{ highlightText }}</span>
            {{ titleAfter }}
          </template>
          <template v-else>
            {{ title }}
          </template>
        </slot>
      </h1>

      <p 
        v-if="showSubtitle"
        class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto mb-8"
        :class="[subtitleClass, alignment !== 'center' ? 'mx-0' : '']"
      >
        <slot name="subtitle">
          {{ subtitle }}
        </slot>
      </p>

      <div class="mt-6">
        <slot />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

interface Props {
  title?: string
  titleBefore?: string
  highlightText?: string
  titleAfter?: string
  subtitle?: string
  
  variant?: 'gradient' | 'solid' | 'image' | 'simple'
  solidColor?: string
  solidColorDark?: string
  backgroundImage?: string
  alignment?: 'center' | 'left' | 'right'
  narrow?: boolean
  headerClass?: string
  titleClass?: string
  subtitleClass?: string
  
  showSubtitle?: boolean
  padding?: 'small' | 'medium' | 'large' | 'none'
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  titleBefore: '',
  highlightText: '',
  titleAfter: '',
  subtitle: '',
  variant: 'gradient',
  solidColor: 'amber-50',
  solidColorDark: 'amber-900/20',
  alignment: 'center',
  narrow: true,
  showSubtitle: true,
  padding: 'medium'
})

const showSubtitle = computed(() => {
  return props.showSubtitle && (props.subtitle || useSlots().subtitle)
})
</script>

<style scoped>

.text-gradient {
  background: linear-gradient(135deg, #f59e0b 0%, #ef4444 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
</style>