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
    /** Street line; the rest of the postal address has its own fields. */
    address: string;
    city: string;
    state: string;
    postalCode: string;
    country: string;
    openingHours: BusinessHours[];
    deliveryFee: number;
    minOrderAmount: number;
    taxRate: number;
    isDeliveryOpen: boolean;
    isPickupOpen: boolean;
    /** Customers may pay on delivery / at pickup instead of online. */
    payOnDeliveryEnabled: boolean;
    /** IANA timezone the opening hours are in, e.g. "Africa/Lagos". */
    timezone: string;
    createdAt: string;
    updatedAt: string;
  }
  
  export interface OpeningHours {
    day: string;
    openTime: string;
    closeTime: string;
    isOpen: boolean;
  }