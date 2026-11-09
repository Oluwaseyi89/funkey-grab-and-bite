// A registered customer.
export interface Customer {
  id: number
  phone: string
  email?: string
  fullName: string
  isVerified: boolean
}

export interface AuthSession {
  user: Customer
  accessToken: string
  isNewUser?: boolean
}

export interface RegistrationDetails {
  fullName: string
  phone: string
  email: string
  password: string
}

// Registering a number that already has guest orders sends a code to the
// email on those orders instead of signing in straight away.
export type RegistrationResult =
  | { kind: 'signedIn'; session: AuthSession }
  | { kind: 'verify'; sentTo: string }

export interface PaymentRefresh {
  checked: boolean
  message: string
}
