import React, { Suspense, lazy } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/layout/Layout';
import Login from '../pages/Login';

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




const LoadingFallback = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
  </div>
);

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
        
        <Route path="/login" element={<Login />} />
        
        
        <Route path="/" element={<PrivateRoute><Layout /></PrivateRoute>}>
          
          <Route index element={<Navigate to="/dashboard" />} />
          <Route path="dashboard" element={<Dashboard />} />
          
          
          <Route path="orders" element={<OrderList />} />
          <Route path="orders/:id" element={<OrderDetails />} />
          <Route path="orders/calendar" element={<OrderCalendar />} />

          
          
          <Route path="menu" element={<MenuList />} />
          <Route path="menu/new" element={<MenuForm />} />
          <Route path="menu/:id/edit" element={<MenuForm />} />
          <Route path="menu/categories" element={<Categories />} />
          
          
          <Route path="customers" element={<CustomerList />} />
          <Route path="customers/:id" element={<CustomerDetails />} />
          
          
          <Route path="catering" element={<CateringList />} />
          <Route path="catering/:id" element={<CateringDetails />} />


          
          
          <Route path="inventory" element={<InventoryList />} />
          <Route path="inventory/alerts" element={<InventoryAlerts />} />

          
          
          <Route path="promotions" element={<PromotionList />} />
          <Route path="/promotions/:id/edit" element={<PromotionForm />} />

          
          
          <Route path="reports" element={<SalesReport />} />
          
          
          <Route path="settings" element={<GeneralSettings />} />
          <Route path="settings/profile" element={<ProfileSettings />} />

          <Route path="settings/admins" element={<AdminUsers />} />
          
          
          <Route path="*" element={<Navigate to="/dashboard" />} />
        </Route>
      </Routes>
    </Suspense>
  );
};

export default AppRoutes;

