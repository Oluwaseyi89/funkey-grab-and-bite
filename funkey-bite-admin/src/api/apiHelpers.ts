import api from './axiosConfig';

export interface ApiHelperResponse<T> {
  data: T;
  meta?: any;
  envelope?: any;
}

export const apiGet = async <T>(url: string, params?: any): Promise<T> => {
  const response = await api.get(url, { params });
  return response.data;
};

export const apiGetResponse = async <T>(url: string, params?: any): Promise<ApiHelperResponse<T>> => {
  const response = await api.get(url, { params });
  return {
    data: response.data,
    meta: (response as any).meta,
    envelope: (response as any).envelope,
  };
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

export const apiUpload = async <T>(url: string, formData: FormData): Promise<T> => {
  const response = await api.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};
