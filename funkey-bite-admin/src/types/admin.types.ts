export interface AdminUser {
    id: number;
    username: string;
    email: string;
    role: string;
    isActive: boolean;
    lastLogin?: string;
    createdAt: string;
  }
  
  export interface AdminStats {
    totalOrders: number;
    totalRevenue: number;
    pendingOrders: number;
    activeCatering: number;
    newCustomers: number;
    popularItems: MenuItemStats[];
  }
  
  export interface MenuItemStats {
    id: number;
    name: string;
    totalSold: number;
    revenue: number;
  }
  
  export interface SalesReport {
    date: string;
    totalOrders: number;
    totalRevenue: number;
    averageOrder: number;
  }
  
  export interface DashboardStats {
    todayRevenue: number;
    todayOrders: number;
    pendingOrders: number;
    newCustomers: number;
    revenueChange: number;
    ordersChange: number;
    customersChange: number;
  }