import React, { createContext, useContext, useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '../stores/authStore';

interface AuthContextType {
  isAuthenticated: boolean;
  user: any;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  user: null,
  loading: true,
});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, user } = useAuthStore();
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const token = localStorage.getItem('admin_token');
    
    if (!token && !location.pathname.includes('/login')) {
      navigate('/login');
    } else if (token && location.pathname === '/login') {
      navigate('/dashboard');
    }
    
    setTimeout(() => setLoading(false), 500); // Simulate loading
  }, [location.pathname, navigate]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, loading: false }}>
      {children}
    </AuthContext.Provider>
  );
};




















    
    


