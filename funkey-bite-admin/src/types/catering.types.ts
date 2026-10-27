// A request moves pending → quoted (quote sent) → accepted (customer agreed;
// invoice issued) → confirmed (deposit paid) → paid (in full; receipt issued)
// → completed. It can be declined or cancelled along the way.
export type CateringStatus =
  | 'pending'
  | 'quoted'
  | 'accepted'
  | 'confirmed'
  | 'paid'
  | 'declined'
  | 'completed'
  | 'cancelled';
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
  acceptedQuoteId?: number;
  invoiceNumber?: string;
  contractedAt?: string;
  paidInFullAt?: string;
  status: CateringStatus;
  createdAt: string;
  updatedAt?: string;
}

// Only notes are edited directly; prices and statuses go through quotes,
// payments and status actions.
export interface CateringRequestUpdate {
  adminNotes?: string;
}

export type CateringQuoteStatus = 'sent' | 'accepted' | 'superseded' | 'withdrawn';
export type CateringPaymentMethod = 'cash' | 'bank_transfer' | 'paystack_transfer' | 'card';

export interface CateringQuoteLine {
  description: string;
  amount: number;
}

export interface CateringQuote {
  id: number;
  requestId: number;
  version: number;
  amount: number;
  lineItems: CateringQuoteLine[];
  depositPercent: number;
  depositAmount: number;
  terms: string;
  message?: string;
  validUntil: string;
  status: CateringQuoteStatus;
  isExpired: boolean;
  acceptedAt?: string;
  acceptedVia?: 'email_link' | 'status_page';
  createdAt: string;
}

export interface CateringQuoteInput {
  amount: number;
  lineItems: CateringQuoteLine[];
  depositPercent?: number;
  terms?: string;
  message?: string;
  validDays?: number;
}

export interface CateringPayment {
  id: number;
  requestId: number;
  amount: number;
  method: CateringPaymentMethod;
  reference?: string;
  note?: string;
  recordedBy?: number;
  paidAt: string;
  createdAt: string;
}

export interface CateringPaymentInput {
  amount: number;
  method: 'cash' | 'bank_transfer';
  reference?: string;
  note?: string;
  paidOn?: string;
}

export interface CateringInvoice {
  id: number;
  requestId: number;
  quoteId: number;
  invoiceNumber: string;
  version: number;
  reason: 'accepted' | 'payment';
  total: number;
  depositDue: number;
  amountPaid: number;
  balance: number;
  issuedAt: string;
}

export interface CateringReceipt {
  id: number;
  requestId: number;
  receiptNumber: string;
  totalPaid: number;
  issuedAt: string;
}

export interface CateringStatusChange {
  id: number;
  requestId: number;
  fromStatus?: CateringStatus;
  toStatus: CateringStatus;
  actorType: 'admin' | 'customer' | 'system';
  actorId?: number;
  note?: string;
  createdAt: string;
}

export interface CateringAccount {
  total: number;
  depositDue: number;
  amountPaid: number;
  balance: number;
  depositRemaining: number;
}

export interface CateringPaymentCharge {
  reference: string;
  amount: number;
  accountNumber?: string;
  accountName?: string;
  bankName?: string;
  expiresAt?: string;
  status: 'pending' | 'paid';
  createdAt: string;
}

// A request with everything that has happened to it.
export interface CateringContract {
  request: CateringRequest;
  quote?: CateringQuote;
  quotes?: CateringQuote[];
  account?: CateringAccount;
  payments: CateringPayment[];
  invoice?: CateringInvoice;
  invoices?: CateringInvoice[];
  receipt?: CateringReceipt;
  pendingCharge?: CateringPaymentCharge;
  history?: CateringStatusChange[];
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
