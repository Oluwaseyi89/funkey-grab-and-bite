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

  // Menu
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

  // Orders
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

  // Catering
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

  // Promotions
  async getActivePromotions(): Promise<any[]> {
    return this.fetchWithFallback('/promotions/active', [])
  }
}

// Singleton instance factory
let apiInstance: ApiService | null = null

export function useApi() {
  if (!apiInstance && process.client) {
    // Client-side: use runtime config
    const config = useRuntimeConfig()
    apiInstance = new ApiService(config.public.apiBaseUrl as string)
  } else if (!apiInstance) {
    // Server-side or fallback
    apiInstance = new ApiService()
  }
  return apiInstance
}

// Backward compatibility - export singleton
export const api = new ApiService()






















// // utils/api.ts
// import type { Ref } from 'vue'
// import { ref } from 'vue'
// import type { MenuCategory, MenuItem } from '~/types/menu'
// import type { Order, CateringRequest } from '~/types/order'
// import { mockCategories, mockMenuItems, mockOrders } from './mockData'
// import { useRuntimeConfig } from 'nuxt/app' // Nuxt auto-import

// export class ApiService {
//   private isBackendAvailable = true
//   private checkBackendStatus: Ref<boolean> = ref(false)
//   private baseURL: string
//   private headers = {
//     'Content-Type': 'application/json',
//   }

//   constructor() {
//     const config = useRuntimeConfig()
//     this.baseURL = config.public.apiBaseUrl as string
//     this.checkBackend()
//   }

//   private async checkBackend() {
//     try {
//       const res = await fetch(`${this.baseURL}/health`)
//       this.isBackendAvailable = res.ok
//     } catch {
//       this.isBackendAvailable = false
//     } finally {
//       this.checkBackendStatus.value = true
//     }
//   }

//   private async fetchWithFallback<T>(endpoint: string, mockData: T): Promise<T> {
//     if (!this.isBackendAvailable && this.checkBackendStatus.value) {
//       console.warn(`Using mock data for ${endpoint}`)
//       return mockData
//     }

//     try {
//       const res = await fetch(`${this.baseURL}${endpoint}`)
//       if (!res.ok) throw new Error(`HTTP ${res.status}`)
//       return (await res.json()) as T
//     } catch (err) {
//       console.error(`API error at ${endpoint}:`, err)
//       return mockData
//     }
//   }

//   // Menu
//   async getMenuCategories(): Promise<MenuCategory[]> {
//     return this.fetchWithFallback('/menu/categories', mockCategories)
//   }

//   async getMenuItems(categoryId?: string): Promise<MenuItem[]> {
//     const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
//     return categoryId ? items.filter(i => i.categoryId === categoryId) : items
//   }

//   async getMenuItem(id: string): Promise<MenuItem | null> {
//     const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
//     return items.find(i => i.id === id) || null
//   }

//   // Orders
//   async createOrder(orderData: Partial<Order>): Promise<Order> {
//     if (!this.isBackendAvailable) {
//       return {
//         id: Math.random().toString(36).substr(2, 9),
//         orderNumber: `FG-${Date.now()}`,
//         customerName: orderData.customerName || 'Guest',
//         customerPhone: orderData.customerPhone || '',
//         orderType: orderData.orderType || 'pickup',
//         status: 'pending',
//         totalAmount: orderData.totalAmount || 0,
//         items: orderData.items || [],
//         createdAt: new Date().toISOString(),
//       }
//     }
//     const res = await fetch(`${this.baseURL}/orders`, {
//       method: 'POST',
//       headers: this.headers,
//       body: JSON.stringify(orderData),
//     })
//     return (await res.json()) as Order
//   }

//   async getOrder(orderNumber: string): Promise<Order | null> {
//     return this.fetchWithFallback(`/orders/${orderNumber}`, mockOrders[0] || null)
//   }

//   // Catering
//   async submitCateringRequest(req: CateringRequest): Promise<CateringRequest> {
//     if (!this.isBackendAvailable) {
//       return { ...req, id: Math.random().toString(36).substr(2, 9), status: 'pending', createdAt: new Date().toISOString() }
//     }
//     const res = await fetch(`${this.baseURL}/catering`, {
//       method: 'POST',
//       headers: this.headers,
//       body: JSON.stringify(req),
//     })
//     return (await res.json()) as CateringRequest
//   }

//   // Promotions
//   async getActivePromotions(): Promise<any[]> {
//     return this.fetchWithFallback('/promotions/active', [])
//   }
// }

// export const api = new ApiService()





















// import type { Ref } from 'vue'
// import { ref } from 'vue'
// import type { MenuCategory, MenuItem } from '~/types/menu'
// import type { Order, CateringRequest } from '~/types/order'
// import { mockCategories, mockMenuItems, mockOrders } from './mockData'
// import { useRuntimeConfig } from 'nuxt/app' // <-- use Nuxt auto-import

// export class ApiService {
//   private isBackendAvailable: boolean = true
//   private checkBackendStatus: Ref<boolean> = ref(false)
//   private apiBaseUrl: string

//   constructor() {
//     // Initialize runtime config at constructor time
//     const config = useRuntimeConfig()
//     this.apiBaseUrl = config.public.apiBaseUrl as string
//     this.checkBackend()
//   }

//   private async checkBackend() {
//     try {
//       const response = await fetch(`${this.apiBaseUrl}/health`)
//       this.isBackendAvailable = response.ok
//     } catch {
//       this.isBackendAvailable = false
//     }
//     this.checkBackendStatus.value = true
//   }

//   private async fetchWithFallback<T>(endpoint: string, mockData: T): Promise<T> {
//     if (!this.isBackendAvailable && this.checkBackendStatus.value) {
//       console.warn(`Using mock data for ${endpoint}`)
//       return mockData
//     }

//     try {
//       const response = await fetch(`${this.apiBaseUrl}${endpoint}`, {
//         headers: {
//           'Content-Type': 'application/json',
//         },
//       })
//       if (!response.ok) throw new Error(`HTTP ${response.status}`)
//       return (await response.json()) as T
//     } catch (error) {
//       console.error(`API Error (${endpoint}):`, error)
//       return mockData
//     }
//   }

//   // --- Menu endpoints ---
//   async getMenuCategories(): Promise<MenuCategory[]> {
//     return this.fetchWithFallback('/menu/categories', mockCategories)
//   }

//   async getMenuItems(categoryId?: string): Promise<MenuItem[]> {
//     const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
//     return categoryId ? items.filter(item => item.categoryId === categoryId) : items
//   }

//   async getMenuItem(id: string): Promise<MenuItem | null> {
//     const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
//     return items.find(item => item.id === id) || null
//   }

//   // --- Orders ---
//   async createOrder(orderData: Partial<Order>): Promise<Order> {
//     if (!this.isBackendAvailable) {
//       return {
//         id: Math.random().toString(36).substr(2, 9),
//         orderNumber: `FG-${Date.now()}`,
//         customerName: orderData.customerName || 'Guest',
//         customerPhone: orderData.customerPhone || '',
//         orderType: orderData.orderType || 'pickup',
//         status: 'pending',
//         totalAmount: orderData.totalAmount || 0,
//         items: orderData.items || [],
//         createdAt: new Date().toISOString(),
//       }
//     }

//     const response = await fetch(`${this.apiBaseUrl}/orders`, {
//       method: 'POST',
//       headers: { 'Content-Type': 'application/json' },
//       body: JSON.stringify(orderData),
//     })
//     return (await response.json()) as Order
//   }

//   async getOrder(orderNumber: string): Promise<Order | null> {
//     return this.fetchWithFallback(`/orders/${orderNumber}`, mockOrders[0] || null)
//   }

//   // --- Catering ---
//   async submitCateringRequest(request: CateringRequest): Promise<CateringRequest> {
//     if (!this.isBackendAvailable) {
//       return {
//         ...request,
//         id: Math.random().toString(36).substr(2, 9),
//         status: 'pending',
//         createdAt: new Date().toISOString(),
//       }
//     }

//     const response = await fetch(`${this.apiBaseUrl}/catering`, {
//       method: 'POST',
//       headers: { 'Content-Type': 'application/json' },
//       body: JSON.stringify(request),
//     })
//     return (await response.json()) as CateringRequest
//   }

//   // --- Promotions ---
//   async getActivePromotions(): Promise<any[]> {
//     return this.fetchWithFallback('/promotions/active', [])
//   }
// }

// export const api = new ApiService()























// import type { Ref } from 'vue'
// import { ref } from 'vue'
// import type { MenuCategory, MenuItem } from '~/types/menu'
// import type { Order, CateringRequest } from '~/types/order'
// import { mockCategories, mockMenuItems, mockOrders } from './mockData'
// import { useRuntimeConfig } from 'nuxt/app' 

// const API_CONFIG = {
//   baseURL: useRuntimeConfig().public.apiBaseUrl,
//   timeout: 10000,
//   headers: {
//     'Content-Type': 'application/json',
//   },
// }

// class ApiService {
//   private isBackendAvailable: boolean = true
//   private checkBackendStatus: Ref<boolean> = ref(false)

//   constructor() {
//     this.checkBackend()
//   }

//   private async checkBackend() {
//     try {
//       const response = await fetch(`${API_CONFIG.baseURL}/health`)
//       this.isBackendAvailable = response.ok
//     } catch {
//       this.isBackendAvailable = false
//     }
//     this.checkBackendStatus.value = true
//   }

//   private async fetchWithFallback<T>(endpoint: string, mockData: T): Promise<T> {
//     if (!this.isBackendAvailable && this.checkBackendStatus.value) {
//       console.warn(`Using mock data for ${endpoint}`)
//       return mockData
//     }

//     try {
//       const response = await fetch(`${API_CONFIG.baseURL}${endpoint}`)
//       if (!response.ok) throw new Error(`HTTP ${response.status}`)
//       return (await response.json()) as T
//     } catch (error) {
//       console.error(`API Error (${endpoint}):`, error)
//       return mockData
//     }
//   }

//   // Menu endpoints
//   async getMenuCategories(): Promise<MenuCategory[]> {
//     return this.fetchWithFallback('/menu/categories', mockCategories)
//   }

//   async getMenuItems(categoryId?: string): Promise<MenuItem[]> {
//     const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
//     return categoryId ? items.filter(item => item.categoryId === categoryId) : items
//   }

//   async getMenuItem(id: string): Promise<MenuItem | null> {
//     const items = await this.fetchWithFallback('/menu/items', mockMenuItems)
//     return items.find(item => item.id === id) || null
//   }

//   // Order endpoints
//   async createOrder(orderData: Partial<Order>): Promise<Order> {
//     if (!this.isBackendAvailable) {
//       const mockOrder: Order = {
//         id: Math.random().toString(36).substr(2, 9),
//         orderNumber: `FG-${Date.now()}`,
//         customerName: orderData.customerName || 'Guest',
//         customerPhone: orderData.customerPhone || '',
//         orderType: orderData.orderType || 'pickup',
//         status: 'pending',
//         totalAmount: orderData.totalAmount || 0,
//         items: orderData.items || [],
//         createdAt: new Date().toISOString(),
//       }
//       return mockOrder
//     }

//     const response = await fetch(`${API_CONFIG.baseURL}/orders`, {
//       method: 'POST',
//       headers: API_CONFIG.headers,
//       body: JSON.stringify(orderData),
//     })
//     return (await response.json()) as Order
//   }

//   async getOrder(orderNumber: string): Promise<Order | null> {
//     return this.fetchWithFallback(`/orders/${orderNumber}`, mockOrders[0] || null)
//   }

//   // Catering endpoints
//   async submitCateringRequest(request: CateringRequest): Promise<CateringRequest> {
//     if (!this.isBackendAvailable) {
//       return {
//         ...request,
//         id: Math.random().toString(36).substr(2, 9),
//         status: 'pending',
//         createdAt: new Date().toISOString(),
//       }
//     }

//     const response = await fetch(`${API_CONFIG.baseURL}/catering`, {
//       method: 'POST',
//       headers: API_CONFIG.headers,
//       body: JSON.stringify(request),
//     })
//     return (await response.json()) as CateringRequest
//   }

//   // Promotions
//   async getActivePromotions(): Promise<any[]> {
//     return this.fetchWithFallback('/promotions/active', [])
//   }
// }

// export const api = new ApiService()
