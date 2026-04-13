import type { Ref } from 'vue'
import { ref } from 'vue'
import type { MenuCategory, MenuItem } from '~/types/menu'
import type { Order, CateringRequest } from '~/types/order'
import { mockCategories, mockMenuItems, mockOrders } from './mockData'
import { useRuntimeConfig } from 'nuxt/app'

export class ApiService {
  private isBackendAvailable = true
  private checkBackendStatus: Ref<boolean> = ref(false)
  private baseURL: string
  private headers = {
    'Content-Type': 'application/json',
  }

  constructor(baseURL: string = '') {
    this.baseURL = baseURL || 'http://localhost:3000/api'
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

  private async fetchWithFallback<T>(endpoint: string, mockData: T): Promise<T> {
    if (!this.isBackendAvailable && this.checkBackendStatus.value) {
      console.warn(`Using mock data for ${endpoint}`)
      return mockData
    }

    try {
      const res = await fetch(`${this.baseURL}${endpoint}`)
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      return (await res.json()) as T
    } catch (err) {
      console.error(`API error at ${endpoint}:`, err)
      return mockData
    }
  }

  async getMenuCategories(): Promise<MenuCategory[]> {
    return this.fetchWithFallback('/menu/categories', mockCategories)
  }

  async getMenuItems(categoryId?: string): Promise<MenuItem[]> {
    const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
    return categoryId ? items.filter(i => i.categoryId === categoryId) : items
  }

  async getMenuItem(id: string): Promise<MenuItem | null> {
    const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
    return items.find(i => i.id === id) || null
  }

  async createOrder(orderData: Partial<Order>): Promise<Order> {
    if (!this.isBackendAvailable) {
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
    const res = await fetch(`${this.baseURL}/orders`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(orderData),
    })
    return (await res.json()) as Order
  }

  async getOrder(orderNumber: string): Promise<Order | null> {
    return this.fetchWithFallback(`/orders/${orderNumber}`, mockOrders[0] || null)
  }

  async submitCateringRequest(req: CateringRequest): Promise<CateringRequest> {
    if (!this.isBackendAvailable) {
      return { ...req, id: Math.random().toString(36).substr(2, 9), status: 'pending', createdAt: new Date().toISOString() }
    }
    const res = await fetch(`${this.baseURL}/catering`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(req),
    })
    return (await res.json()) as CateringRequest
  }

  async getActivePromotions(): Promise<any[]> {
    return this.fetchWithFallback('/promotions/active', [])
  }
}

let apiInstance: ApiService | null = null

export function useApi() {
  if (!apiInstance && process.client) {
    const config = useRuntimeConfig()
    apiInstance = new ApiService(config.public.apiBaseUrl as string)
  } else if (!apiInstance) {
    apiInstance = new ApiService()
  }
  return apiInstance
}

export const api = new ApiService()

