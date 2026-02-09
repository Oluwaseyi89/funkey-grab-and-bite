// src/api/axiosConfig.ts
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling and response formatting
api.interceptors.response.use(
  (response) => {
    // Your Go API returns different formats:
    // 1. {success: true, data: {...}, error: null} - for handlers using handlers.Success()
    // 2. Direct data {...} - for handlers returning direct JSON
    
    // Check if it's the success wrapper format
    if (response.data && typeof response.data === 'object' && 'success' in response.data) {
      if (response.data.success) {
        // Extract data from success wrapper
        return {
          ...response,
          data: response.data.data || response.data // Fallback to whole response if no data field
        };
      } else {
        // Handle API error in success wrapper format
        const error = new Error(response.data.error?.message || 'API request failed');
        (error as any).code = response.data.error?.code;
        (error as any).details = response.data.error?.details;
        return Promise.reject(error);
      }
    }
    
    // Direct response format (no wrapper)
    return response;
  },
  (error) => {
    // Handle HTTP errors
    if (axios.isAxiosError(error)) {
      // Unauthorized - redirect to login
      if (error.response?.status === 401) {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user');
        window.location.href = '/login';
        return Promise.reject(new Error('Session expired. Please login again.'));
      }
      
      // Forbidden - admin access required
      if (error.response?.status === 403) {
        return Promise.reject(new Error('Admin access required'));
      }
      
      // Extract error message from your Go API response format
      const apiError = error.response?.data;
      let errorMessage = 'Request failed';
      
      if (apiError?.error?.message) {
        // Format: {error: {message: "...", details: "..."}}
        errorMessage = apiError.error.message;
      } else if (apiError?.error) {
        // Format: {error: "error message"}
        errorMessage = apiError.error;
      } else if (apiError?.message) {
        // Format: {message: "error message"}
        errorMessage = apiError.message;
      } else if (error.response?.data) {
        // Direct error string
        errorMessage = error.response.data;
      }
      
      const formattedError = new Error(errorMessage);
      (formattedError as any).status = error.response?.status;
      (formattedError as any).code = apiError?.error?.code;
      return Promise.reject(formattedError);
    }
    
    // Network errors
    if (error.message === 'Network Error') {
      return Promise.reject(new Error('Network error. Please check your connection.'));
    }
    
    return Promise.reject(error);
  }
);

export default api;

// Helper function for consistent API calls
export async function apiRequest<T>(promise: Promise<any>): Promise<T> {
  try {
    const response = await promise;
    return response.data;
  } catch (error: any) {
    // Already formatted by interceptor
    throw error;
  }
}
