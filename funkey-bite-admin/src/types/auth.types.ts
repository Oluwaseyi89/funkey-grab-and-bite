import type { AdminUser } from "./admin.types";

export interface User {
    id: number;
    phone: string;
    email?: string;
    fullName: string;
    isVerified: boolean;
    isActive: boolean;
    lastLogin?: string;
    createdAt: string;
    updatedAt?: string;
  }

  /** A customer as the admin sees them. Cancelled orders are left out; totalSpent counts paid orders only. */
  export interface Customer extends User {
    orderCount: number;
    totalSpent: number;
    firstOrderAt?: string;
    lastOrderAt?: string;
    favoriteCategory?: string;
  }

  export type CustomerStatusFilter = 'active' | 'inactive' | 'verified' | 'unverified';
  
  export interface AuthResponse {
    user: User;
    token: string;
  }
  
  export interface LoginCredentials {
    phone: string;
    password: string;
  }
  
  export interface RegisterCredentials {
    phone: string;
    email?: string;
    fullName: string;
    password: string;
  }
  
  export interface AdminLoginCredentials {
    email: string;
    password: string;
  }
  
  export interface AdminLoginResponse {
    admin: AdminUser;
    token: string;
  }











  
  
  
  
  
