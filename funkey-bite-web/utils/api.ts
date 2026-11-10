import type { MenuCategory, MenuItem, Promotion, PromotionCheck } from '~/types/menu'
import type { Order, OrderCreatePayload, OrderLookupResult } from '~/types/order'
import type { AuthSession, Customer, PaymentRefresh, RegistrationDetails, RegistrationResult } from '~/types/account'
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
import { useRuntimeConfig } from 'nuxt/app'

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
  paymentAccessCode?: string
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
  private baseURL: string
  private authToken: string | null = null

  // Requests carry the signed-in customer's token, so orders placed while
  // signed in are linked to their account.
  private get headers(): Record<string, string> {
    return this.authToken
      ? { 'Content-Type': 'application/json', Authorization: `Bearer ${this.authToken}` }
      : { 'Content-Type': 'application/json' }
  }

  setAuthToken(token: string | null) {
    this.authToken = token
  }

  constructor(baseURL: string = '') {
    this.baseURL = baseURL || 'http://localhost:8080/api/v1'
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
      // Card orders have no account; an empty value means none.
      paymentAccountNumber: (order.paymentAccountNumber ?? order.payment_account_number) || undefined,
      paymentAccountName: order.paymentAccountName ?? order.payment_account_name,
      paymentBankName: order.paymentBankName ?? order.payment_bank_name,
      paymentExpiresAt: order.paymentExpiresAt ?? order.payment_expires_at,
      paymentAccessCode: order.paymentAccessCode,
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

  private async fetchJSON<T>(endpoint: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseURL}${endpoint}`, init)
    if (!res.ok) {
      const { message, code } = await this.errorDetails(res, endpoint)
      const retryAfter = Number(res.headers.get('Retry-After'))
      throw new ApiRequestError(message, res.status, code, Number.isFinite(retryAfter) && retryAfter > 0 ? retryAfter : undefined)
    }

    return this.unwrapPayload<T>(await res.json())
  }

  // Lists the menu and the like. A network failure becomes a plain message;
  // the site never shows made-up items in place of the real menu.
  private async fetchList<T>(endpoint: string): Promise<T> {
    try {
      return await this.fetchJSON<T>(endpoint)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        throw new Error("We couldn't reach our kitchen's server. Check your connection and try again.")
      }
      if (err instanceof ApiRequestError && (err.status === 429 || err.status >= 500)) {
        throw new Error("We're a little busy right now. Please try again in a moment.")
      }
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
    const payload = await this.fetchList<BackendMenuCategory[]>('/menu/categories')
    return (payload || []).map((item) => this.normalizeMenuCategory(item))
  }

  async getMenuItems(categoryId?: string): Promise<MenuItem[]> {
    const endpoint = categoryId ? `/menu/category?category_id=${encodeURIComponent(categoryId)}` : '/menu'
    const payload = await this.fetchList<BackendMenuItem[]>(endpoint)
    return (payload || []).map((item) => this.normalizeMenuItem(item))
  }

  async getMenuItem(id: string): Promise<MenuItem | null> {
    const payload = await this.fetchList<BackendMenuItem | null>(`/menu/${encodeURIComponent(id)}`)
    return payload ? this.normalizeMenuItem(payload) : null
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

  // Payment details for an unpaid online transfer order, when getting them at
  // checkout failed or the payment window closed. Details still valid are
  // returned as they are, so pressing twice is safe.
  async getPaymentDetails(orderNumber: string, phone: string): Promise<Order> {
    try {
      const payload = await this.fetchJSON<BackendOrder>('/orders/payment-details', {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify({ orderNumber, phone }),
      })
      return this.normalizeOrder(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        throw new Error("We couldn't reach our server. Check your connection and try again.")
      }
      throw err
    }
  }

  // ---- Customer accounts ----

  private post<T>(endpoint: string, body: unknown): Promise<T> {
    return this.fetchJSON<T>(endpoint, { method: 'POST', headers: this.headers, body: JSON.stringify(body) })
  }

  async register(details: RegistrationDetails): Promise<RegistrationResult> {
    const payload = await this.post<AuthSession | { verificationRequired: true; sentTo: string }>('/auth/register', details)
    if ('verificationRequired' in payload) return { kind: 'verify', sentTo: payload.sentTo }
    return { kind: 'signedIn', session: payload }
  }

  // Completes a registration for a number with guest orders, linking them.
  confirmRegistration(details: RegistrationDetails, code: string): Promise<AuthSession> {
    return this.post<AuthSession>('/auth/register/confirm', { ...details, code })
  }

  login(phone: string, password: string): Promise<AuthSession> {
    return this.post<AuthSession>('/auth/login', { phone, password })
  }

  async forgotPassword(email: string): Promise<void> {
    await this.post('/auth/password/forgot', { email })
  }

  async resetPassword(email: string, code: string, password: string): Promise<void> {
    await this.post('/auth/password/reset', { email, code, password })
  }

  async getProfile(): Promise<Customer> {
    const payload = await this.fetchJSON<{ user: Customer }>('/auth/profile', { headers: this.headers })
    return payload.user
  }

  async getMyOrders(page = 1, limit = 10): Promise<{ orders: Order[]; totalPages: number }> {
    const payload = await this.fetchJSON<{ orders: BackendOrder[]; pagination: { totalPages: number } }>(
      `/auth/orders?page=${page}&limit=${limit}`,
      { headers: this.headers },
    )
    return { orders: (payload.orders || []).map((o) => this.normalizeOrder(o)), totalPages: payload.pagination?.totalPages ?? 0 }
  }

  // One of the signed-in customer's orders with its invoice (null for orders
  // paid online). Another customer's order comes back as not found.
  async getMyOrder(orderId: string): Promise<OrderLookupResult> {
    const payload = await this.fetchJSON<{ order: BackendOrder; invoice: OrderLookupResult['invoice'] }>(
      `/auth/orders/${encodeURIComponent(orderId)}`,
      { headers: this.headers },
    )
    return { order: this.normalizeOrder(payload.order), invoice: payload.invoice ?? null }
  }

  // ---- Payment status ----

  // Asks our server to check an unpaid online order with Paystack, by order
  // number and phone; returns the order as it now stands.
  async refreshPaymentStatus(orderNumber: string, phone: string): Promise<{ order: Order } & PaymentRefresh> {
    const payload = await this.post<{ order: BackendOrder } & PaymentRefresh>('/orders/payment-status', { orderNumber, phone })
    return { ...payload, order: this.normalizeOrder(payload.order) }
  }

  // The same for a signed-in customer's own order.
  async refreshMyPaymentStatus(orderId: string): Promise<{ order: Order } & PaymentRefresh> {
    const payload = await this.fetchJSON<{ order: BackendOrder } & PaymentRefresh>(`/auth/orders/${orderId}/payment-status`, {
      method: 'POST',
      headers: this.headers,
    })
    return { ...payload, order: this.normalizeOrder(payload.order) }
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
  async payCateringOnline(
    auth: CateringBookingAuth,
    portion: 'deposit' | 'balance',
    method: 'transfer' | 'card' = 'transfer',
  ): Promise<CateringPaymentCharge> {
    return this.fetchJSON<CateringPaymentCharge>('/catering/booking/pay', {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify({ ...auth, portion, method }),
    })
  }

  // Has the server confirm an online payment with Paystack (e.g. right after
  // the card popup reports success) and returns the booking as it now stands.
  async verifyCateringPayment(auth: CateringBookingAuth, reference: string): Promise<CateringBooking> {
    return this.fetchJSON<CateringBooking>('/catering/booking/pay/verify', {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify({ ...auth, paymentReference: reference }),
    })
  }

  // Promotions are extras: if they can't load, the site simply hides them.
  async getActivePromotions(): Promise<Promotion[]> {
    try {
      const payload = await this.fetchJSON<Promotion[]>('/promotions/active')
      return (payload || []).filter((promo) => promo.isActive && new Date(promo.validUntil).getTime() > Date.now())
    } catch (err) {
      console.warn('[API] Could not load promotions:', err instanceof Error ? err.message : err)
      return []
    }
  }

  // Asks the server whether a code would apply to this order total. The
  // discount itself is worked out again when the order is placed.
  async checkPromotion(code: string, amount: number): Promise<PromotionCheck> {
    const params = new URLSearchParams({ code: code.trim(), amount: String(Math.max(0, amount)) })
    try {
      return await this.fetchJSON<PromotionCheck>(`/promotions/validate?${params}`)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        throw new Error("We couldn't check that code. Check your connection and try again.")
      }
      throw err
    }
  }

  // Contact details and opening hours never fall back to demo data: a made-up
  // address or hours would mislead customers, so callers get null instead.
  async getBusinessSettings(): Promise<BusinessSettings | null> {
    try {
      return await this.fetchJSON<BusinessSettings>('/settings')
    } catch (err) {
      console.warn('[API] Could not load business settings:', err instanceof Error ? err.message : err)
      return null
    }
  }
}

let apiInstance: ApiService | null = null

export function useApi() {
  if (!apiInstance) {
    const config = useRuntimeConfig()
    apiInstance = new ApiService(config.public.apiBaseUrl as string)
  }
  return apiInstance
}
