export interface PaginatedResponse<T> {
    data: T[];
    pagination: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
    };
  }
  
  export interface ApiResponse<T = any> {
    success: boolean;
    data?: T;
    error?: {
      code: string;
      message: string;
      details?: string;
    };
    meta?: any;
  }
  
  export interface BusinessHours {
    day: string;
    openTime: string;
    closeTime: string;
    isOpen: boolean;
  }
  
  export interface BusinessSettings {
    id: number;
    businessName: string;
    phoneNumber: string;
    email: string;
    address: string;
    openingHours: BusinessHours[];
    deliveryFee: number;
    minOrderAmount: number;
    taxRate: number;
    isDeliveryOpen: boolean;
    isPickupOpen: boolean;
    createdAt: string;
    updatedAt: string;
  }
  
  export interface OpeningHours {
    day: string;
    openTime: string;
    closeTime: string;
    isOpen: boolean;
  }