import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
});

const isLoginEndpoint = (url?: string) => {
  if (!url) return false;
  return url.includes('/api/v1/admin/auth/login');
};

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_token');
    if (token && !isLoginEndpoint(config.url)) {
      config.headers.Authorization = `Bearer ${token}`;
    } else if (config.headers?.Authorization) {
      delete config.headers.Authorization;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => {
    
    if (response.data && typeof response.data === 'object' && 'success' in response.data) {
      if (response.data.success) {
        return {
          ...response,
          data: response.data.data || response.data // Fallback to whole response if no data field
        };
      } else {
        const error = new Error(response.data.error?.message || 'API request failed');
        (error as any).code = response.data.error?.code;
        (error as any).details = response.data.error?.details;
        return Promise.reject(error);
      }
    }
    
    return response;
  },
  (error) => {
    if (axios.isAxiosError(error)) {
      if (error.response?.status === 401) {
        const requestUrl = error.config?.url;
        const apiError = error.response?.data;

        if (isLoginEndpoint(requestUrl)) {
          let loginErrorMessage = 'Invalid email or password';
          if (apiError?.error?.message) {
            loginErrorMessage = apiError.error.message;
          } else if (apiError?.error) {
            loginErrorMessage = apiError.error;
          } else if (apiError?.message) {
            loginErrorMessage = apiError.message;
          }
          return Promise.reject(new Error(loginErrorMessage));
        }

        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user');
        localStorage.removeItem('admin-auth-storage');
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
        return Promise.reject(new Error('Session expired. Please login again.'));
      }
      
      if (error.response?.status === 403) {
        return Promise.reject(new Error('Admin access required'));
      }
      
      const apiError = error.response?.data;
      let errorMessage = 'Request failed';
      
      if (apiError?.error?.message) {
        errorMessage = apiError.error.message;
      } else if (apiError?.error) {
        errorMessage = apiError.error;
      } else if (apiError?.message) {
        errorMessage = apiError.message;
      } else if (error.response?.data) {
        errorMessage = error.response.data;
      }
      
      const formattedError = new Error(errorMessage);
      (formattedError as any).status = error.response?.status;
      (formattedError as any).code = apiError?.error?.code;
      return Promise.reject(formattedError);
    }
    
    if (error.message === 'Network Error') {
      return Promise.reject(new Error('Network error. Please check your connection.'));
    }
    
    return Promise.reject(error);
  }
);

export default api;

export async function apiRequest<T>(promise: Promise<any>): Promise<T> {
  try {
    const response = await promise;
    return response.data;
  } catch (error: any) {
    throw error;
  }
}
