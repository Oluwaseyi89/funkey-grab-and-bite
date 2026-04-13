import { createToast } from 'mosha-vue-toastify'
import 'mosha-vue-toastify/dist/style.css'
import { defineNuxtPlugin } from 'nuxt/app'

export default defineNuxtPlugin((nuxtApp) => {
  const toast = {
    success(message: string, options = {}) {
      createToast(message, {
        type: 'success',
        position: 'top-right',
        showIcon: true,
        timeout: 3000,
        ...options,
      })
    },
    
    error(message: string, options = {}) {
      createToast(message, {
        type: 'danger',
        position: 'top-right',
        showIcon: true,
        timeout: 4000,
        ...options,
      })
    },
    
    warning(message: string, options = {}) {
      createToast(message, {
        type: 'warning',
        position: 'top-right',
        showIcon: true,
        timeout: 3500,
        ...options,
      })
    },
    
    info(message: string, options = {}) {
      createToast(message, {
        type: 'info',
        position: 'top-right',
        showIcon: true,
        timeout: 3000,
        ...options,
      })
    },
  }

  nuxtApp.provide('toast', toast)
})

    
    
    

