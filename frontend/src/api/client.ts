import axios from 'axios';

// Create a core axios instance
export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  timeout: 60000, // 60 seconds timeout since OCR processing can take time
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor for global response/error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);
