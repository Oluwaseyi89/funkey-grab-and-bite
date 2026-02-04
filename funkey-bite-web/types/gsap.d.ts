import type { GSAP } from 'gsap'

declare module '#app' {
  interface NuxtApp {
    $gsap: typeof GSAP
    $animations: {
      fadeIn: (element: HTMLElement, delay?: number) => void
      staggerChildren: (parent: HTMLElement, childSelector: string) => void
      floatAnimation: (element: HTMLElement) => void
      scrollReveal: (element: HTMLElement) => void
    }
  }
}

export {}