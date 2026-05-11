import api from './axiosConfig';

export interface ApiHelperResponse<T> {
  data: T;
  meta?: any;
  envelope?: any;
}

type ErrorWithStatus = {
  status?: number;
  code?: string;
  message?: string;
};

export const isBackendUnavailableError = (error: unknown): boolean => {
  const candidate = (error || {}) as ErrorWithStatus;
  const status = candidate.status;
  const code = candidate.code;
  const message = (candidate.message || '').toLowerCase();

  if (typeof status === 'number' && status >= 500) {
    return true;
  }

  if (code === 'ERR_NETWORK' || code === 'ECONNABORTED' || code === 'ETIMEDOUT') {
    return true;
  }

  return (
    message.includes('network error') ||
    message.includes('failed to fetch') ||
    message.includes('timeout') ||
    message.includes('service unavailable') ||
    message.includes('unavailable')
  );
};

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
