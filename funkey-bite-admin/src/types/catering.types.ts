// A request moves pending → quoted → confirmed → completed; it can be
// declined by the business or cancelled by the customer.
export type CateringStatus = 'pending' | 'quoted' | 'confirmed' | 'declined' | 'completed' | 'cancelled';
export type CateringRequestType = 'package' | 'custom';
export type ServiceStyle = 'drop_off' | 'buffet' | 'full_service';

export interface CateringSelectedItem {
  menuItemId: number;
  name: string;
}

export interface CateringRequest {
  id: number;
  reference: string;
  requestType: CateringRequestType;
  userId?: number;
  packageId?: number;
  packageName?: string;
  contactName: string;
  contactPhone: string;
  contactEmail?: string;
  organization?: string;
  eventType: string;
  eventTypeOther?: string;
  eventName?: string;
  eventDate: string;
  eventTime?: string;
  durationHours?: number;
  guestCount: number;
  venueName?: string;
  venueAddress?: string;
  venueArea?: string;
  venueState?: string;
  venueNotes?: string;
  serviceStyle?: ServiceStyle;
  selectedItems: CateringSelectedItem[];
  menuNotes?: string;
  dietaryRequirements: string[];
  dietaryNotes?: string;
  budget?: number;
  specialRequests?: string;
  estimatedTotal?: number;
  quotedAmount?: number;
  adminNotes?: string;
  status: CateringStatus;
  createdAt: string;
  updatedAt?: string;
}

export interface CateringRequestUpdate {
  status?: CateringStatus;
  quotedAmount?: number;
  adminNotes?: string;
}

export interface CateringPackageItem {
  menuItemId: number;
  name: string;
  isAvailable: boolean;
}

export interface CateringPackage {
  id: number;
  name: string;
  description: string;
  pricePerPerson: number;
  minGuests: number;
  maxGuests?: number;
  serviceStyle: ServiceStyle;
  minNoticeDays: number;
  inclusions: string[];
  dishChoices: number;
  menuItems: CateringPackageItem[];
  isActive: boolean;
  displayOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface CateringPackageInput {
  name: string;
  description: string;
  pricePerPerson: number;
  minGuests: number;
  maxGuests?: number;
  serviceStyle: ServiceStyle;
  minNoticeDays: number;
  inclusions: string[];
  dishChoices: number;
  menuItemIds: number[];
  isActive: boolean;
  displayOrder: number;
}
