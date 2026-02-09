import React, { Suspense, lazy } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/layout/Layout';
import Login from '../pages/Login';

// Lazy load pages for better performance
const Dashboard = lazy(() => import('../pages/Dashboard'));
const OrderList = lazy(() => import('../pages/Orders/OrderList'));
const OrderDetails = lazy(() => import('../pages/Orders/OrderDetails'));
const MenuList = lazy(() => import('../pages/Menu/MenuList'));
const MenuForm = lazy(() => import('../pages/Menu/MenuForm'));
const Categories = lazy(() => import('../pages/Menu/Categories'));
const CustomerList = lazy(() => import('../pages/Customers/CustomerList'));
const CustomerDetails = lazy(() => import('../pages/Customers/CustomerDetails'));
const CateringList = lazy(() => import('../pages/Catering/CateringList'));
const InventoryList = lazy(() => import('../pages/Inventory/InventoryList'));
const InventoryAlerts = lazy(() => import('../pages/Inventory/InventoryAlerts'));
const PromotionList = lazy(() => import('../pages/Promotions/PromotionList'));
const PromotionForm = lazy(() => import('../pages/Promotions/PromotionForm'));
const SalesReport = lazy(() => import('../pages/Reports/SalesReport'));
const Analytics = lazy(() => import('../pages/Reports/Analytics'));
const GeneralSettings = lazy(() => import('../pages/Settings/GeneralSettings'));
const ProfileSettings = lazy(() => import('../pages/Settings/ProfileSettings'));
const AdminUsers = lazy(() => import('../pages/Settings/AdminUsers'));
const OrderCalendar = lazy(() => import('../pages/Orders/OrderCalendar'));
const CateringDetails = lazy(() => import('../pages/Catering/CateringDetails'));




// Loading fallback component
const LoadingFallback = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
  </div>
);

// Private route wrapper
const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, loading } = useAuth();
  
  if (loading) {
    return <LoadingFallback />;
  }
  
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

const AppRoutes: React.FC = () => {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<Login />} />
        
        {/* Protected routes */}
        <Route path="/" element={<PrivateRoute><Layout /></PrivateRoute>}>
          {/* Dashboard */}
          <Route index element={<Navigate to="/dashboard" />} />
          <Route path="dashboard" element={<Dashboard />} />
          
          {/* Orders */}
          <Route path="orders" element={<OrderList />} />
          <Route path="orders/:id" element={<OrderDetails />} />
          <Route path="orders/calendar" element={<OrderCalendar />} />

          
          {/* Menu */}
          <Route path="menu" element={<MenuList />} />
          <Route path="menu/new" element={<MenuForm />} />
          <Route path="menu/:id/edit" element={<MenuForm />} />
          <Route path="menu/categories" element={<Categories />} />
          
          {/* Customers */}
          <Route path="customers" element={<CustomerList />} />
          <Route path="customers/:id" element={<CustomerDetails />} />
          
          {/* Catering */}
          <Route path="catering" element={<CateringList />} />
          <Route path="catering/:id" element={<CateringDetails />} />


          
          {/* Inventory */}
          <Route path="inventory" element={<InventoryList />} />
          <Route path="inventory/alerts" element={<InventoryAlerts />} />

          
          {/* Promotions */}
          <Route path="promotions" element={<PromotionList />} />
          <Route path="/promotions/:id/edit" element={<PromotionForm />} />

          
          {/* Reports */}
          <Route path="reports" element={<SalesReport />} />
          
          {/* Settings */}
          <Route path="settings" element={<GeneralSettings />} />
          <Route path="settings/profile" element={<ProfileSettings />} />

          <Route path="settings/admins" element={<AdminUsers />} />
          
          {/* Catch all route */}
          <Route path="*" element={<Navigate to="/dashboard" />} />
        </Route>
      </Routes>
    </Suspense>
  );
};

export default AppRoutes;

