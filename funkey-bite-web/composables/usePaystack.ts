import { useRuntimeConfig } from 'nuxt/app'

// Paystack's popup: the customer enters their card in Paystack's own window,
// so card details never touch this site. The script loads only when needed.
const SCRIPT_URL = 'https://js.paystack.co/v2/inline.js'

type PopupResult = 'paid' | 'cancelled'

interface PaystackPopInstance {
  resumeTransaction: (
    accessCode: string,
    callbacks: { onSuccess: () => void; onCancel: () => void; onError?: (error: { message?: string }) => void },
  ) => void
}

declare global {
  interface Window {
    PaystackPop?: new () => PaystackPopInstance
  }
}

let loading: Promise<void> | null = null

const loadScript = () => {
  if (window.PaystackPop) return Promise.resolve()
  loading ??= new Promise<void>((resolve, reject) => {
    const script = document.createElement('script')
    script.src = SCRIPT_URL
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => {
      loading = null
      reject(new Error("We couldn't load the card payment window. Check your connection and try again."))
    }
    document.head.appendChild(script)
  })
  return loading
}

export const usePaystack = () => {
  const config = useRuntimeConfig()
  const cardPaymentsEnabled = !!config.public.paystackPublicKey

  // Opens the popup for a payment the server started. Resolves 'paid' when
  // Paystack reports success (the server still confirms it) or 'cancelled'.
  const payWithCard = async (accessCode: string): Promise<PopupResult> => {
    await loadScript()
    if (!window.PaystackPop) throw new Error('The card payment window is unavailable right now.')
    const popup = new window.PaystackPop()
    return new Promise<PopupResult>((resolve, reject) => {
      popup.resumeTransaction(accessCode, {
        onSuccess: () => resolve('paid'),
        onCancel: () => resolve('cancelled'),
        onError: (error) => reject(new Error(error?.message || 'The card payment failed. Please try again.')),
      })
    })
  }

  return { cardPaymentsEnabled, payWithCard }
}
