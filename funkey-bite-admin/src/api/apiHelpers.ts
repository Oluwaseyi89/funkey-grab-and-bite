// src/api/apiHelpers.ts
import api from './axiosConfig';

export const apiGet = async <T>(url: string, params?: any): Promise<T> => {
  const response = await api.get(url, { params });
  return response.data;
};

export const apiPost = async <T>(url: string, data?: any): Promise<T> => {
  const response = await api.post(url, data);
  return response.data;
};

export const apiPut = async <T>(url: string, data?: any): Promise<T> => {
  const response = await api.put(url, data);
  return response.data;
};

export const apiPatch = async <T>(url: string, data?: any): Promise<T> => {
  const response = await api.patch(url, data);
  return response.data;
};

export const apiDelete = async <T>(url: string): Promise<T> => {
  const response = await api.delete(url);
  return response.data;
};

// Helper for file uploads
export const apiUpload = async <T>(url: string, formData: FormData): Promise<T> => {
  const response = await api.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};
