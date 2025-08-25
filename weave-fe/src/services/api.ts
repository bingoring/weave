import axios, { AxiosResponse } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/v1/api';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Unauthorized - clear token and redirect to login
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_data');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Generic response interface
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

export interface PaginatedApiResponse<T = any> extends ApiResponse<T> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

// API helper functions
export const apiRequest = async <T = any>(
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH',
  url: string,
  data?: any,
  config?: any
): Promise<T> => {
  try {
    let response: AxiosResponse<ApiResponse<T>>;
    
    switch (method) {
      case 'GET':
        response = await api.get(url, config);
        break;
      case 'POST':
        response = await api.post(url, data, config);
        break;
      case 'PUT':
        response = await api.put(url, data, config);
        break;
      case 'DELETE':
        response = await api.delete(url, config);
        break;
      case 'PATCH':
        response = await api.patch(url, data, config);
        break;
      default:
        throw new Error(`Unsupported method: ${method}`);
    }

    if (!response.data.success) {
      throw new Error(response.data.error || 'API request failed');
    }

    return response.data.data as T;
  } catch (error: any) {
    if (error.response?.data?.error) {
      throw new Error(error.response.data.error);
    }
    throw error;
  }
};

export const paginatedApiRequest = async <T = any>(
  method: 'GET',
  url: string,
  params?: any
): Promise<{ data: T; pagination: any }> => {
  try {
    const response: AxiosResponse<PaginatedApiResponse<T>> = await api.get(url, { params });

    if (!response.data.success) {
      throw new Error(response.data.error || 'API request failed');
    }

    return {
      data: response.data.data as T,
      pagination: response.data.pagination,
    };
  } catch (error: any) {
    if (error.response?.data?.error) {
      throw new Error(error.response.data.error);
    }
    throw error;
  }
};

export default api;