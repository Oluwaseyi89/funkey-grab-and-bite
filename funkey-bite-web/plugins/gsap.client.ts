import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { TextPlugin } from 'gsap/TextPlugin'
// import type { GSAP } from 'gsap'
// import { gsap } from 'gsap'
import { defineNuxtPlugin } from 'nuxt/app'

// Define proper types for your animations
interface AnimationHelpers {
  fadeIn: (element: HTMLElement, delay?: number) => void
  staggerChildren: (parent: HTMLElement, childSelector: string) => void
  floatAnimation: (element: HTMLElement) => void
  scrollReveal: (element: HTMLElement) => void
}

export default defineNuxtPlugin((nuxtApp) => {
  if (import.meta.client) {
    gsap.registerPlugin(ScrollTrigger, TextPlugin)
    
    // Global animation helpers
    const animations: AnimationHelpers = {
      fadeIn(element: HTMLElement, delay: number = 0) {
        gsap.from(element, {
          opacity: 0,
          y: 20,
          duration: 0.8,
          delay,
          ease: 'power2.out',
        })
      },
      
      staggerChildren(parent: HTMLElement, childSelector: string) {
        gsap.from(parent.querySelectorAll(childSelector), {
          opacity: 0,
          y: 30,
          duration: 0.6,
          stagger: 0.1,
          ease: 'power2.out',
        })
      },
      
      floatAnimation(element: HTMLElement) {
        gsap.to(element, {
          y: -10,
          duration: 2,
          repeat: -1,
          yoyo: true,
          ease: 'power1.inOut',
        })
      },
      
      scrollReveal(element: HTMLElement) {
        gsap.from(element, {
          scrollTrigger: {
            trigger: element,
            start: 'top 80%',
          },
          opacity: 0,
          y: 50,
          duration: 1,
          ease: 'power2.out',
        })
      },
    }

    // Inject with proper typing
    nuxtApp.provide('gsap', gsap)
    nuxtApp.provide('animations', animations)
  }
})

























// import { gsap } from 'gsap'
// import { ScrollTrigger } from 'gsap/ScrollTrigger'
// import { TextPlugin } from 'gsap/TextPlugin'
// import { defineNuxtPlugin } from 'nuxt/app'

// export default defineNuxtPlugin(() => {
//   if (import.meta.client) {
//     gsap.registerPlugin(ScrollTrigger, TextPlugin)
    
//     // Global animation helpers
//     const animations = {
//       fadeIn(element: HTMLElement, delay: number = 0) {
//         gsap.from(element, {
//           opacity: 0,
//           y: 20,
//           duration: 0.8,
//           delay,
//           ease: 'power2.out',
//         })
//       },
      
//       staggerChildren(parent: HTMLElement, childSelector: string) {
//         gsap.from(parent.querySelectorAll(childSelector), {
//           opacity: 0,
//           y: 30,
//           duration: 0.6,
//           stagger: 0.1,
//           ease: 'power2.out',
//         })
//       },
      
//       floatAnimation(element: HTMLElement) {
//         gsap.to(element, {
//           y: -10,
//           duration: 2,
//           repeat: -1,
//           yoyo: true,
//           ease: 'power1.inOut',
//         })
//       },
      
//       scrollReveal(element: HTMLElement) {
//         gsap.from(element, {
//           scrollTrigger: {
//             trigger: element,
//             start: 'top 80%',
//           },
//           opacity: 0,
//           y: 50,
//           duration: 1,
//           ease: 'power2.out',
//         })
//       },
//     }

//     return {
//       provide: {
//         gsap,
//         animations,
//       },
//     }
//   }
// })