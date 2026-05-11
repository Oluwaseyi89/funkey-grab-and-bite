import type { Ref } from 'vue'
import { ref } from 'vue'
import type { MenuCategory, MenuItem } from '~/types/menu'
import type { Order, CateringRequest } from '~/types/order'
import { mockCategories, mockMenuItems, mockOrders, mockPromotions } from './mockData'
import { useNuxtApp, useRuntimeConfig } from 'nuxt/app'

type ToastNotify = (message: string) => void
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
}

type BackendCateringRequest = {
  id: string | number
  eventName?: string
  event_name?: string
  contactName: string
  contact_name?: string
  contactPhone: string
  contact_phone?: string
  contactEmail?: string
  contact_email?: string
  eventDate: string
  event_date?: string
  eventTime?: string
  event_time?: string
  guestCount: number
  guest_count?: number
  eventType: string
  event_type?: string
  budget?: number
  specialRequests?: string
  special_requests?: string
  status: CateringRequest['status']
  createdAt: string
  created_at?: string
}

type CateringRequestPayload = {
  eventName?: string
  contactName: string
  contactPhone: string
  contactEmail?: string
  eventDate: string
  eventTime?: string
  guestCount: number
  eventType: string
  package?: string
  budget?: number
  specialRequests?: string
}

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
    this.baseURL = baseURL || 'http://localhost:3000/api'
    this.notifyMockUsage = notifyMockUsage
    this.checkBackend()
  }

  private async checkBackend() {
    try {
      const res = await fetch(`${this.baseURL}/health`)
      this.isBackendAvailable = res.ok
    } catch {
      this.isBackendAvailable = false
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
      isPreOrder: item.isPreOrder ?? item.is_pre_order ?? false,
      preparationTime: item.preparationTime ?? item.preparation_time ?? 0,
      tags: Array.isArray(item.tags) ? item.tags : [],
      nutritionalInfo: item.nutritionalInfo ?? item.nutritional_info,
    }
  }

  private normalizeOrder(order: BackendOrder): Order {
    return {
      id: String(order.id),
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
    }
  }

  private normalizeCateringRequest(request: BackendCateringRequest): CateringRequest {
    return {
      id: String(request.id),
      eventName: request.eventName ?? request.event_name,
      contactName: request.contactName ?? request.contact_name ?? '',
      contactPhone: request.contactPhone ?? request.contact_phone ?? '',
      contactEmail: request.contactEmail ?? request.contact_email,
      eventDate: request.eventDate ?? request.event_date ?? '',
      eventTime: request.eventTime ?? request.event_time,
      guestCount: request.guestCount ?? request.guest_count ?? 0,
      eventType: request.eventType ?? request.event_type ?? '',
      budget: request.budget,
      specialRequests: request.specialRequests ?? request.special_requests,
      status: request.status,
      createdAt: request.createdAt ?? request.created_at ?? new Date().toISOString(),
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

  private notifyFallback(endpoint: string, reason: FallbackReason) {
    const key = `${endpoint}:${reason}`
    if (this.notifiedFallbacks.has(key)) {
      return
    }
    this.notifiedFallbacks.add(key)

    const message = reason === 'backend-unavailable'
      ? `Backend unavailable. Using showcase data for ${endpoint}.`
      : `Backend returned no records for ${endpoint}. Using showcase data.`

    if (this.notifyMockUsage) {
      this.notifyMockUsage(message)
      return
    }

    console.warn(message)
  }

  private async fetchJSON<T>(endpoint: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseURL}${endpoint}`, init)
    if (!res.ok) {
      throw new Error(`HTTP ${res.status} at ${endpoint}`)
    }

    return this.unwrapPayload<T>(await res.json())
  }

  private async fetchWithFallback<T>(
    endpoint: string,
    mockData: T,
    normalize?: (payload: T) => T,
  ): Promise<T> {
    if (!this.isBackendAvailable && this.checkBackendStatus.value) {
      this.notifyFallback(endpoint, 'backend-unavailable')
      return mockData
    }

    try {
      const payload = await this.fetchJSON<T>(endpoint)
      const normalized = normalize ? normalize(payload) : payload

      if (Array.isArray(normalized) && normalized.length === 0) {
        this.notifyFallback(endpoint, 'empty-backend-data')
        return mockData
      }

      return normalized
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback(endpoint, 'backend-unavailable')
        return mockData
      }

      throw err
    }
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
    if (!this.isBackendAvailable && this.checkBackendStatus.value) {
      this.notifyFallback(`/menu/${id}`, 'backend-unavailable')
      return mockMenuItems.find(i => i.id === id) || null
    }

    try {
      const payload = await this.fetchJSON<BackendMenuItem>(`/menu/${id}`)
      return this.normalizeMenuItem(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback(`/menu/${id}`, 'backend-unavailable')
        return mockMenuItems.find(i => i.id === id) || null
      }
      throw err
    }
  }

  async createOrder(orderData: Partial<Order>): Promise<Order> {
    if (!this.isBackendAvailable && this.checkBackendStatus.value) {
      this.notifyFallback('/orders', 'backend-unavailable')
      return {
        id: Math.random().toString(36).substr(2, 9),
        orderNumber: `FG-${Date.now()}`,
        customerName: orderData.customerName || 'Guest',
        customerPhone: orderData.customerPhone || '',
        orderType: orderData.orderType || 'pickup',
        status: 'pending',
        totalAmount: orderData.totalAmount || 0,
        items: orderData.items || [],
        createdAt: new Date().toISOString(),
      }
    }

    try {
      const payload = await this.fetchJSON<BackendOrder>('/orders', {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(orderData),
      })
      return this.normalizeOrder(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback('/orders', 'backend-unavailable')
        return {
          id: Math.random().toString(36).substr(2, 9),
          orderNumber: `FG-${Date.now()}`,
          customerName: orderData.customerName || 'Guest',
          customerPhone: orderData.customerPhone || '',
          orderType: orderData.orderType || 'pickup',
          status: 'pending',
          totalAmount: orderData.totalAmount || 0,
          items: orderData.items || [],
          createdAt: new Date().toISOString(),
        }
      }

      throw err
    }
  }

  async getOrder(orderNumber: string): Promise<Order | null> {
    if (!this.isBackendAvailable && this.checkBackendStatus.value) {
      this.notifyFallback('/orders/track/:orderNumber', 'backend-unavailable')
      return mockOrders[0] || null
    }

    try {
      const payload = await this.fetchJSON<BackendOrder>(`/orders/track/${encodeURIComponent(orderNumber)}`)
      return this.normalizeOrder(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback('/orders/track/:orderNumber', 'backend-unavailable')
        return mockOrders[0] || null
      }

      throw err
    }
  }

  async submitCateringRequest(req: CateringRequestPayload): Promise<CateringRequest> {
    if (!this.isBackendAvailable && this.checkBackendStatus.value) {
      this.notifyFallback('/catering/requests', 'backend-unavailable')
      return { ...req, id: Math.random().toString(36).substr(2, 9), status: 'pending', createdAt: new Date().toISOString() }
    }

    try {
      const payload = await this.fetchJSON<BackendCateringRequest>('/catering/requests', {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(req),
      })
      return this.normalizeCateringRequest(payload)
    } catch (err) {
      if (this.isBackendUnavailableError(err)) {
        this.markBackendUnavailable()
        this.notifyFallback('/catering/requests', 'backend-unavailable')
        return { ...req, id: Math.random().toString(36).substr(2, 9), status: 'pending', createdAt: new Date().toISOString() }
      }

      throw err
    }
  }

  async getActivePromotions(): Promise<any[]> {
    return this.fetchWithFallback('/promotions/active', mockPromotions)
  }
}

let apiInstance: ApiService | null = null

export function useApi() {
  if (!apiInstance) {
    const config = useRuntimeConfig()
    let notifyMockUsage: ToastNotify | undefined

    if (process.client) {
      const { $toast } = useNuxtApp()
      const toast = $toast as { warning?: (message: string) => void } | undefined
      notifyMockUsage = (message: string) => {
        if (toast && typeof toast.warning === 'function') {
          toast.warning(message)
          return
        }
        console.warn(message)
      }
    }

    apiInstance = new ApiService(config.public.apiBaseUrl as string, notifyMockUsage)
  }

  return apiInstance
}

