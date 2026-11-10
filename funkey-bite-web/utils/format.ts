// One way to show money and dates across the site.

const wholeNaira = new Intl.NumberFormat('en-NG', { maximumFractionDigits: 0 })
const exactNaira = new Intl.NumberFormat('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 })

// ₦4,200 for whole amounts, ₦4,200.50 when there are kobo.
export const naira = (amount: number | null | undefined) => {
  const value = Number(amount ?? 0)
  return `₦${Number.isInteger(Math.round(value * 100) / 100) ? wholeNaira.format(value) : exactNaira.format(value)}`
}

// 10 Nov 2026, 1:14 PM
export const formatDateTime = (value?: string | null) =>
  value ? new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' }) : ''

// 10 Nov 2026
export const formatDay = (value?: string | null) =>
  value ? new Date(value).toLocaleDateString('en-NG', { dateStyle: 'medium' }) : ''

// 1:14 PM
export const formatTime = (value?: string | null) =>
  value ? new Date(value).toLocaleTimeString('en-NG', { hour: 'numeric', minute: '2-digit' }) : ''

export const plural = (count: number, one: string, many = `${one}s`) => `${count} ${count === 1 ? one : many}`
