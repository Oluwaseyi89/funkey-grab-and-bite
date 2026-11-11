export interface AdminUser {
    id: number;
    username: string;
    email: string;
    role: string;
    isActive: boolean;
    lastLogin?: string;
    createdAt: string;
  }
  
