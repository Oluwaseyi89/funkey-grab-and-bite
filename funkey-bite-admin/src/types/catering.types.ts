export type CateringStatus = 'pending' | 'confirmed' | 'declined' | 'completed';

export interface CateringRequest {
  id: number;
  userId?: number;
  eventName?: string;
  contactName: string;
  contactPhone: string;
  contactEmail?: string;
  eventDate: string;
  eventTime?: string;
  guestCount: number;
  eventType: string;
  package?: string;
  budget?: number;
  specialRequests?: string;
  status: CateringStatus;
  createdAt: string;
}

export interface CateringPackage {
  id: string;
  name: string;
  description: string;
  pricePerPerson: number;
  minGuests: number;
  maxGuests?: number;
  includes: string[];
}