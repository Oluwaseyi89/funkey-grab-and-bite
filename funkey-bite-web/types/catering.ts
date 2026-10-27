export type CateringStatus = 'pending' | 'quoted' | 'confirmed' | 'declined' | 'completed' | 'cancelled'
export type CateringRequestType = 'package' | 'custom'
export type ServiceStyle = 'drop_off' | 'buffet' | 'full_service'

export interface CateringPackageItem {
  menuItemId: number
  name: string
  isAvailable: boolean
}

// An admin-defined catering offer from GET /catering/packages.
export interface CateringPackage {
  id: number
  name: string
  description: string
  pricePerPerson: number
  minGuests: number
  maxGuests?: number
  serviceStyle: ServiceStyle
  minNoticeDays: number
  inclusions: string[]
  // How many of menuItems the customer picks; 0 means a fixed menu.
  dishChoices: number
  menuItems: CateringPackageItem[]
}

// What the customer submits; the API validates every field again.
export interface CateringRequestPayload {
  requestType: CateringRequestType
  packageId?: number
  contactName: string
  contactPhone: string
  contactEmail: string
  organization?: string
  eventType: string
  eventTypeOther?: string
  eventName?: string
  eventDate: string
  eventTime: string
  durationHours: number
  guestCount: number
  venueName?: string
  venueAddress: string
  venueArea: string
  venueState: string
  venueNotes?: string
  serviceStyle?: ServiceStyle
  selectedItemIds: number[]
  menuNotes?: string
  dietaryRequirements: string[]
  dietaryNotes?: string
  budget?: number
  specialRequests?: string
}

// The customer's view of a request.
export interface CateringRequest {
  id: number
  reference: string
  requestType: CateringRequestType
  packageName?: string
  contactName: string
  contactPhone: string
  contactEmail?: string
  eventType: string
  eventTypeOther?: string
  eventName?: string
  eventDate: string
  eventTime?: string
  durationHours?: number
  guestCount: number
  venueName?: string
  venueAddress?: string
  venueArea?: string
  venueState?: string
  serviceStyle?: ServiceStyle
  selectedItems: { menuItemId: number; name: string }[]
  menuNotes?: string
  dietaryRequirements: string[]
  estimatedTotal?: number
  quotedAmount?: number
  status: CateringStatus
  createdAt: string
}

export const SERVICE_STYLES: { id: ServiceStyle; name: string; description: string }[] = [
  { id: 'drop_off', name: 'Drop-off', description: 'Food delivered ready to serve; you handle serving.' },
  { id: 'buffet', name: 'Buffet with staff', description: 'We set up and staff a buffet for your guests.' },
  { id: 'full_service', name: 'Full service', description: 'Servers, setup and cleanup handled by us.' },
]

export const EVENT_TYPES: { id: string; name: string }[] = [
  { id: 'party', name: 'Party / Owambe' },
  { id: 'wedding', name: 'Wedding' },
  { id: 'birthday', name: 'Birthday' },
  { id: 'naming', name: 'Naming ceremony' },
  { id: 'burial', name: 'Burial / Remembrance' },
  { id: 'church', name: 'Church event' },
  { id: 'corporate', name: 'Corporate / Office' },
  { id: 'graduation', name: 'Graduation' },
  { id: 'other', name: 'Other' },
]

export const DIETARY_OPTIONS: { id: string; name: string }[] = [
  { id: 'vegetarian', name: 'Vegetarian' },
  { id: 'vegan', name: 'Vegan' },
  { id: 'halal', name: 'Halal' },
  { id: 'no_pork', name: 'No pork' },
  { id: 'nut_free', name: 'Nut-free' },
  { id: 'gluten_free', name: 'Gluten-free' },
  { id: 'dairy_free', name: 'Dairy-free' },
]

export const CATERING_STATUS_LABELS: Record<CateringStatus, string> = {
  pending: 'Received — we will be in touch',
  quoted: 'Quote sent',
  confirmed: 'Confirmed',
  declined: 'Declined',
  completed: 'Completed',
  cancelled: 'Cancelled',
}

// Custom requests have no package to set a notice period; matches the API.
export const CUSTOM_NOTICE_DAYS = 3
