import type { Ref } from 'vue'
import { ref } from 'vue'
import type { MenuCategory, MenuItem } from '~/types/menu'
import type { Order, OrderCreatePayload, OrderLookupResult } from '~/types/order'
import type {
  CateringBooking,
  CateringBookingAuth,
  CateringPackage,
  CateringPaymentCharge,
  CateringRequest,
  CateringRequestPayload,
} from '~/types/catering'
import type { TeamMember } from '~/types/team'
import type { BusinessSettings } from '~/types/settings'
import { mockCategories, mockMenuItems, mockOrders, mockPromotions } from './mockData'
import { useNuxtApp, useRuntimeConfig } from 'nuxt/app'

type ToastNotify = (endpoint: string, reason: FallbackReason, requestSucceeded: boolean) => void
type FallbackReason = 'backend-unavailable' | 'empty-backend-data'

type BackendMenuCategory = {
  id: number | string
  name: string
  description?: string
  displayOrder?: number
  display_order?: number
  isActive?: boolean
  is_active?: boolean
}

type BackendMenuItem = {
  id: number | string
  categoryId?: number | string
  category_id?: number | string
  name: string
  description?: string
  price: number
  imageUrl?: string
  image_url?: string
  isAvailable?: boolean
  is_available?: boolean
  stockQuantity?: number
  stock_quantity?: number
  isPreOrder?: boolean
  is_pre_order?: boolean
  preparationTime?: number
  preparation_time?: number
  tags?: string[]
  nutritionalInfo?: MenuItem['nutritionalInfo']
  nutritional_info?: MenuItem['nutritionalInfo']
}

type BackendOrderItem = {
  menuItemId?: string | number
  menu_item_id?: string | number
  name: string
  quantity: number
  unitPrice?: number
  unit_price?: number
  specialInstructions?: string
  special_instructions?: string
}

type BackendOrder = {
  id: string | number
  orderNumber?: string
  order_number?: string
  customerName?: string
  customer_name?: string
  customerPhone?: string
  customer_phone?: string
  customerEmail?: string
  customer_email?: string
  orderType?: Order['orderType']
  order_type?: Order['orderType']
  status: Order['status']
  totalAmount?: number
  total_amount?: number
  notes?: string
  pickupTime?: string
  pickup_time?: string
  createdAt?: string
  created_at?: string
  estimatedReadyTime?: string
  estimated_ready_time?: string
  items?: BackendOrderItem[]
  paymentMethod?: Order['paymentMethod']
  payment_method?: Order['paymentMethod']
  paymentStatus?: Order['paymentStatus']
  payment_status?: Order['paymentStatus']
  paymentReference?: string
  payment_reference?: string
  paymentAccountNumber?: string
  payment_account_number?: string
  paymentAccountName?: string
  payment_account_name?: string
  paymentBankName?: string
  payment_bank_name?: string
  paymentExpiresAt?: string
  payment_expires_at?: string
  paymentPaidAt?: string
  payment_paid_at?: string
  paymentCollectedVia?: Order['paymentCollectedVia']
  deliveryAddress?: string
  invoiceNumber?: string
  subtotalAmount?: number
  deliveryFee?: number
  taxAmount?: number
  discountAmount?: number
}

// A rejected API request, with what's needed to decide whether to retry it.
export class ApiRequestError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly code?: string,
    readonly retryAfterSeconds?: number,
  ) {
    super(message)
    this.name = 'ApiRequestError'
  }
}

// Checkout is retried at most this many times in all, and only when it is
// safe: the request never arrived, the server failed, or the first attempt is
// still running. The Idempotency-Key makes a retry return the same order.
const CHECKOUT_ATTEMPTS = 3

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

export class ApiService {
  private isBackendAvailable = true
  private checkBackendStatus: Ref<boolean> = ref(false)
  private baseURL: string
  private notifyMockUsage?: ToastNotify
  private notifiedFallbacks = new Set<string>()
  private headers = {
    'Content-Type': 'application/json',
  }

  constructor(baseURL: string = '', notifyMockUsage?: ToastNotify) {
    // Align default port with nuxt.config.ts (8080), though config should always provide value
    this.baseURL = baseURL || 'http://localhost:8080/api/v1'
    this.notifyMockUsage = notifyMockUsage
    this.checkBackend()
  }

  private async checkBackend() {
    try {
      const health_url = `${this.baseURL}/health`
      const res = await fetch(health_url)
      this.isBackendAvailable = res.ok
      console.info(`[API Health Check] ${health_url} → ${res.status}`, this.isBackendAvailable ? '✅ Backend available' : '❌ Backend unavailable')
    } catch (err) {
      this.isBackendAvailable = false
      console.warn(`[API Health Check] Failed to reach ${this.baseURL}/health`, err instanceof Error ? err.message : 'unknown error')
    } finally {
      this.checkBackendStatus.value = true
    }
  }

  private markBackendUnavailable() {
    this.isBackendAvailable = false
    this.checkBackendStatus.value = true
  }

  private normalizeMenuCategory(item: BackendMenuCategory): MenuCategory {
    return {
      id: String(item.id),
      name: item.name,
      description: item.description || '',
      displayOrder: item.displayOrder ?? item.display_order ?? 0,
      isActive: item.isActive ?? item.is_active ?? true,
    }
  }

  private normalizeMenuItem(item: BackendMenuItem): MenuItem {
    return {
      id: String(item.id),
      categoryId: String(item.categoryId ?? item.category_id ?? ''),
      name: item.name,
      description: item.description || '',
      price: item.price,
      imageUrl: item.imageUrl ?? item.image_url ?? '',
      isAvailable: item.isAvailable ?? item.is_available ?? true,
      stockQuantity: item.stockQuantity ?? item.stock_quantity,
      isPreOrder: item.isPreOrder ?? item.is_pre_order ?? false,
      preparationTime: item.preparationTime ?? item.preparation_time ?? 0,
      tags: Array.isArray(item.tags) ? item.tags : [],
      nutritionalInfo: item.nutritionalInfo ?? item.nutritional_info,
    }
  }

  private normalizeOrder(order: BackendOrder): Order {
    return {
      id: order.id != null ? String(order.id) : (order.orderNumber ?? order.order_number ?? `local-${Date.now()}`),
      orderNumber: order.orderNumber ?? order.order_number ?? `FG-${Date.now()}`,
      customerName: order.customerName ?? order.customer_name ?? 'Guest',
      customerPhone: order.customerPhone ?? order.customer_phone ?? '',
      customerEmail: order.customerEmail ?? order.customer_email,
      orderType: order.orderType ?? order.order_type ?? 'pickup',
      status: order.status,
      totalAmount: order.totalAmount ?? order.total_amount ?? 0,
      notes: order.notes,
      pickupTime: order.pickupTime ?? order.pickup_time,
      items: (order.items || []).map((item) => ({
        menuItemId: String(item.menuItemId ?? item.menu_item_id ?? ''),
        name: item.name,
        quantity: item.quantity,
        unitPrice: item.unitPrice ?? item.unit_price ?? 0,
        specialInstructions: item.specialInstructions ?? item.special_instructions,
      })),
      createdAt: order.createdAt ?? order.created_at ?? new Date().toISOString(),
      estimatedReadyTime: order.estimatedReadyTime ?? order.estimated_ready_time,
      paymentMethod: order.paymentMethod ?? order.payment_method,
      paymentStatus: order.paymentStatus ?? order.payment_status,
      paymentReference: order.paymentReference ?? order.payment_reference,
      paymentAccountNumber: order.paymentAccountNumber ?? order.payment_account_number,
      paymentAccountName: order.paymentAccountName ?? order.payment_account_name,
      paymentBankName: order.paymentBankName ?? order.payment_bank_name,
      paymentExpiresAt: order.paymentExpiresAt ?? order.payment_expires_at,
      paymentPaidAt: order.paymentPaidAt ?? order.payment_paid_at,
      paymentCollectedVia: order.paymentCollectedVia,
      deliveryAddress: order.deliveryAddress,
      invoiceNumber: order.invoiceNumber,
      subtotalAmount: order.subtotalAmount,
      deliveryFee: order.deliveryFee,
      taxAmount: order.taxAmount,
      discountAmount: order.discountAmount,
    }
  }

  private isBackendUnavailableError(error: unknown): boolean {
    const message = error instanceof Error ? error.message.toLowerCase() : String(error || '').toLowerCase()
    return (
      message.includes('failed to fetch') ||
      message.includes('networkerror') ||
      message.includes('network error') ||
      message.includes('backend unavailable')
    )
  }

  private notifyFallback(endpoint: string, reason: FallbackReason, requestSucceeded: boolean) {
    const key = `${endpoint}:${reason}`
    if (this.notifiedFallbacks.has(key)) {
      console.info(`[API] Notification already shown for ${key}, skipping duplicate`)
      return
    }
    this.notifiedFallbacks.add(key)

    if (this.notifyMockUsage) {
      console.info(`[API] Triggering toast notification: ${reason} for ${endpoint} | requestSucceeded=${requestSucceeded}`)
      this.notifyMockUsage(endpoint, reason, requestSucceeded)
      return
    }

    // Fallback to console if no toast system available
    const message = !requestSucceeded
      ? `🔴 Backend offline. Using showcase data for ${endpoint}.`
      : `🔵 No records in database for ${endpoint}. Empty result returned.`
    console.warn(message)
  }

  private async fetchJSON<T>(endpoint: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseURL}${endpoint}`, init)
    if (!res.ok) {
      const { message, code } = await this.errorDetails(res, endpoint)
      const retryAfter = Number(res.headers.get('Retry-After'))
      throw new ApiRequestError(message, res.status, code, Number.isFinite(retryAfter) && retryAfter > 0 ? retryAfter : undefined)
    }

    // Any successful API response means backend is reachable.
    this.isBackendAvailable = true
    this.checkBackendStatus.value = true

    return this.unwrapPayload<T>(await res.json())
  }

  private async fetchWithFallback<T>(
    endpoint: string,
    mockData: T,
    normalize?: (payload: T) => T,
  ): Promise<T> {
    try {
      const payload = await this.fetchJSON<T>(endpoint)
      const normalized = normalize ? normalize(payload) : payload

      // Backend is healthy - return actual data (even if empty)
      // Backend is the single source of truth; empty arrays are valid responses
      if (Array.isArray(normalized) && normalized.length === 0) {
        console.info(`[API] 🔵 ${endpoint} returned EMPTY ARRAY (backend responded successfully, database has no records)`)
        console.info(`[API] 🔵 Backend is source of truth → returning empty array to UI`)
        // Pass true because THIS request succeeded (got 200 OK)
        this.notifyFallback(endpoint, 'empty-backend-data', true)
        // Return the real empty array, not mock data - let UI handle "no items" state
        return normalized
      }

      console.info(`[API] ✅ ${endpoint} fetched ${Array.isArray(normalized) ? normalized.length : 1} real item(s) from backend`)
      return normalized
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        console.error(`[API] 🔴 NETWORK ERROR on ${endpoint}:`, err instanceof Error ? err.message : 'unknown error')
        console.error(`[API] 🔴 Backend server unreachable or offline → using mock data for demo`)
        this.markBackendUnavailable()
        // Pass false because THIS request failed (network error)
        this.notifyFallback(endpoint, 'backend-unavailable', false)
        return mockData
      }

      console.error(`[API] ⚠️ Unexpected error on ${endpoint}:`, err)
      throw err
    }
  }

  // The API explains rejections (e.g. a missing delivery address) in the
  // response body; show that instead of a bare status code.
  private async errorDetails(res: Response, endpoint: string): Promise<{ message: string; code?: string }> {
    try {
      const body = await res.json() as { error?: { message?: string; code?: string } | string; message?: string }
      if (typeof body.error === 'string') return { message: body.error }
      if (body.error?.message) return { message: body.error.message, code: body.error.code }
      if (body.message) return { message: body.message }
    } catch {
      // not JSON
    }
    return { message: `HTTP ${res.status} at ${endpoint}` }
  }

  private unwrapPayload<T>(payload: unknown): T {
    if (
      payload &&
      typeof payload === 'object' &&
      'success' in payload &&
      typeof (payload as { success?: unknown }).success === 'boolean'
    ) {
      const envelope = payload as { success: boolean; data?: T; error?: { message?: string } }
      if (!envelope.success) {
        throw new Error(envelope.error?.message || 'API request failed')
      }
      return (envelope.data ?? payload) as T
    }

    return payload as T
  }

  async getMenuCategories(): Promise<MenuCategory[]> {
    return this.fetchWithFallback('/menu/categories', mockCategories, (payload) => {
      return (payload as BackendMenuCategory[]).map((item) => this.normalizeMenuCategory(item)) as MenuCategory[]
    })
  }

  async getMenuItems(categoryId?: string): Promise<MenuItem[]> {
    if (categoryId) {
      return this.fetchWithFallback(`/menu/category?category_id=${encodeURIComponent(categoryId)}`, mockMenuItems, (payload) => {
        return (payload as BackendMenuItem[]).map((item) => this.normalizeMenuItem(item)) as MenuItem[]
      })
    }

    return this.fetchWithFallback('/menu', mockMenuItems, (payload) => {
      return (payload as BackendMenuItem[]).map((item) => this.normalizeMenuItem(item)) as MenuItem[]
    })
  }

  async getMenuItem(id: string): Promise<MenuItem | null> {
    try {
      const payload = await this.fetchJSON<BackendMenuItem>(`/menu/${id}`)
      return this.normalizeMenuItem(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback(`/menu/${id}`, 'backend-unavailable', false)
        return mockMenuItems.find(i => i.id === id) || null
      }
      throw err
    }
  }

  // Places the order. Every attempt sends the same idempotencyKey, so
  // retrying after a dropped connection or a server error returns the order
  // already placed instead of placing a second one. Never pretends to succeed.
  async createOrder(orderData: OrderCreatePayload, idempotencyKey: string): Promise<Order> {
    const body = JSON.stringify(orderData)
    for (let attempt = 1; ; attempt++) {
      try {
        const payload = await this.fetchJSON<BackendOrder>('/orders', {
          method: 'POST',
          headers: { ...this.headers, 'Idempotency-Key': idempotencyKey },
          body,
        })
        return this.normalizeOrder(payload)
      } catch (err) {
        const unreachable = this.isBackendUnavailableError(err)
        const retryable = unreachable || (err instanceof ApiRequestError &&
          (err.status >= 500 || err.code === 'REQUEST_IN_PROGRESS'))
        if (!retryable || attempt >= CHECKOUT_ATTEMPTS) {
          if (unreachable) {
            this.markBackendUnavailable()
            throw new Error("We couldn't reach our server to confirm your order. Check your connection and press Place Order again; it won't be placed twice.")
          }
          if (err instanceof ApiRequestError && err.status >= 500) {
            throw new Error("We couldn't confirm your order. Press Place Order again; it won't be placed twice.")
          }
          throw err
        }
        const waitSeconds = err instanceof ApiRequestError && err.retryAfterSeconds ? err.retryAfterSeconds : attempt
        await sleep(Math.min(waitSeconds, 5) * 1000)
      }
    }
  }

  async getOrder(phone: string, orderNumber: string): Promise<Order | null> {
    const endpoint = `/order/track/${encodeURIComponent(phone)}/${encodeURIComponent(orderNumber)}`
    try {
      const payload = await this.fetchJSON<BackendOrder>(endpoint)
      return this.normalizeOrder(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback('/order/track/:phone/:orderNumber', 'backend-unavailable', false)
        return mockOrders[0] || null
      }

      throw err
    }
  }

  // Looks up the full order and its invoice. The phone goes in the body, not
  // the URL, so it stays out of logs and browser history.
  async lookupOrder(orderNumber: string, phone: string): Promise<OrderLookupResult | null> {
    try {
      const payload = await this.fetchJSON<{ order: BackendOrder; invoice: OrderLookupResult['invoice'] }>('/orders/lookup', {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify({ orderNumber, phone }),
      })
      return { order: this.normalizeOrder(payload.order), invoice: payload.invoice ?? null }
    } catch (err) {
      if (err instanceof Error && /not found/i.test(err.message)) return null
      throw err
    }
  }

  // Catering packages are the admin's offer; with none (or the API down) the
  // page still takes custom requests, so this never falls back to demo data.
  async getCateringPackages(): Promise<CateringPackage[] | null> {
    try {
      return await this.fetchJSON<CateringPackage[]>('/catering/packages')
    } catch (err) {
      console.warn('[API] Could not load catering packages:', err instanceof Error ? err.message : err)
      return null
    }
  }

  // The team section only shows real people added by the admin, so there is
  // no demo fallback: on failure the section is simply hidden.
  async getTeamMembers(): Promise<TeamMember[]> {
    try {
      return await this.fetchJSON<TeamMember[]>('/team')
    } catch (err) {
      console.warn('[API] Could not load team members:', err instanceof Error ? err.message : err)
      return []
    }
  }

  // Never pretends to succeed: a customer must know if their request was not sent.
  async submitCateringRequest(req: CateringRequestPayload): Promise<CateringRequest> {
    try {
      return await this.fetchJSON<CateringRequest>('/catering/requests', {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(req),
      })
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        throw new Error('We could not reach our server, so your request was not sent. Please try again or call us.')
      }
      throw err
    }
  }

  // The customer's request, quote, invoice and payments. The token or phone
  // goes in the body so it stays out of URLs and logs. Null when nothing
  // matches the reference and phone.
  async getCateringBooking(auth: CateringBookingAuth): Promise<CateringBooking | null> {
    try {
      return await this.fetchJSON<CateringBooking>('/catering/booking', {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(auth),
      })
    } catch (err) {
      if (err instanceof Error && /request not found/i.test(err.message)) return null
      throw err
    }
  }

  async acceptCateringQuote(auth: CateringBookingAuth): Promise<CateringBooking> {
    return this.fetchJSON<CateringBooking>('/catering/booking/accept', {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify({ ...auth, agreeToTerms: true }),
    })
  }

  // A bank account to transfer the deposit or balance to.
  async payCateringOnline(auth: CateringBookingAuth, portion: 'deposit' | 'balance'): Promise<CateringPaymentCharge> {
    return this.fetchJSON<CateringPaymentCharge>('/catering/booking/pay', {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify({ ...auth, portion }),
    })
  }

  async getActivePromotions(): Promise<any[]> {
    return this.fetchWithFallback('/promotions/active', mockPromotions)
  }

  // Contact details and opening hours never fall back to demo data: a made-up
  // address or hours would mislead customers, so callers get null instead.
  async getBusinessSettings(): Promise<BusinessSettings | null> {
    try {
      return await this.fetchJSON<BusinessSettings>('/settings')
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
      }
      console.warn('[API] Could not load business settings:', err instanceof Error ? err.message : err)
      return null
    }
  }
}

let apiInstance: ApiService | null = null

export function useApi() {
  if (!apiInstance) {
    const config = useRuntimeConfig()
    let notifyMockUsage: ToastNotify | undefined

    if (process.client) {
      const { $toast } = useNuxtApp()
      const toast = $toast as { warning?: (message: string) => void; info?: (message: string) => void } | undefined
      notifyMockUsage = (endpoint: string, reason: FallbackReason, requestSucceeded: boolean) => {
        // requestSucceeded = true means THIS request got a valid response (backend is reachable)
        // requestSucceeded = false means THIS request failed (network error, backend offline)
        const isUnavailable = !requestSucceeded
        const dataSetEmpty = requestSucceeded && reason === 'empty-backend-data'

        console.info(`[API Fallback] ${endpoint} → reason=${reason} | requestSucceeded=${requestSucceeded} | isUnavailable=${isUnavailable} | dataSetEmpty=${dataSetEmpty}`)

        const message = isUnavailable
          ? `⚠️ Backend server is offline. Showing demo data for ${endpoint}.`
          : dataSetEmpty
            ? `ℹ️ No records found for ${endpoint}. Database is empty.`
            : `ℹ️ Could not load ${endpoint}. Please try again.`

        if (toast) {
          if (isUnavailable && typeof toast.warning === 'function') {
            toast.warning(message)
            console.info(`[API] 🔴 Backend offline → mock data for ${endpoint}`)
            return
          }
          if (dataSetEmpty && typeof toast.info === 'function') {
            toast.info(message)
            console.info(`[API] 🔵 Empty dataset (backend online) → mock data for ${endpoint}`)
            return
          }
          console.warn(`[API] Toast method not available for ${endpoint} (isUnavailable=${isUnavailable}, dataSetEmpty=${dataSetEmpty})`)
        }

        console.info(`[API] ${message}`)
      }
    }

    apiInstance = new ApiService(config.public.apiBaseUrl as string, notifyMockUsage)
  }

  return apiInstance
}

