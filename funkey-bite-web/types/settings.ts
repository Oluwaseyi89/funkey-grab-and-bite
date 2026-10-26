// Times are 24-hour "HH:MM" in the business timezone. A closeTime at or before
// openTime means the business closes after midnight.
export interface OpeningHours {
  day: string
  openTime: string
  closeTime: string
  isOpen: boolean
}

// Public business settings from GET /settings.
export interface BusinessSettings {
  businessName: string
  phoneNumber: string
  email: string
  address: string
  city: string
  state: string
  postalCode: string
  country: string
  openingHours: OpeningHours[]
  deliveryFee: number
  minOrderAmount: number
  taxRate: number
  isDeliveryOpen: boolean
  isPickupOpen: boolean
  isOpenNow: boolean
}
