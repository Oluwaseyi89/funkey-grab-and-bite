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










// import type { AdminUser } from "./menu.types";

// export interface User {
//     id: number;
//     phone: string;
//     email?: string;
//     fullName: string;
//     isVerified: boolean;
//     isActive: boolean;
//     lastLogin?: string;
//     createdAt: string;
//     updatedAt?: string;
//   }
  
//   export interface AuthResponse {
//     user: User;
//     token: string;
//   }
  
//   export interface LoginCredentials {
//     phone: string;
//     password: string;
//   }
  
//   export interface RegisterCredentials {
//     phone: string;
//     email?: string;
//     fullName: string;
//     password: string;
//   }
  
//   export interface AdminLoginCredentials {
//     email: string;
//     password: string;
//   }
  
//   export interface AdminLoginResponse {
//     admin: AdminUser;
//     token: string;
//   }