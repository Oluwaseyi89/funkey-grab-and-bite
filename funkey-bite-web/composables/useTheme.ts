import { computed } from 'vue'
import { useColorMode } from '@vueuse/core'

// Light or dark, remembered by the browser. Starts from the phone's setting.
export function useTheme() {
  const mode = useColorMode({ emitAuto: false })
  const isDark = computed(() => mode.value === 'dark')
  const toggle = () => {
    mode.value = isDark.value ? 'light' : 'dark'
  }
  return { isDark, toggle }
}
