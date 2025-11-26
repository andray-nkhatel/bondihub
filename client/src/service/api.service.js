// src/services/api.service.js
import axios from 'axios';

// Get API base URL from environment variable
function getApiBaseUrl() {
  const envUrl = import.meta.env.VITE_API_BASE_URL;
  
  if (!envUrl) {
    console.error('❌ VITE_API_BASE_URL is not set in environment variables!');
    console.error('💡 Create a .env file with: VITE_API_BASE_URL=http://localhost:8080/api/v1');
    throw new Error('VITE_API_BASE_URL environment variable is required');
  }
  
  if (import.meta.env.DEV) {
    console.log('✅ Using API base URL:', envUrl);
  }
  
  return envUrl;
}

// Create axios instance
const apiClient = axios.create({
  baseURL: getApiBaseUrl(),
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 30000,
  transformResponse: [
    function (data, headers) {
      const contentType = headers?.['content-type'] || headers?.['Content-Type'] || '';
      
      if (contentType.includes('application/json') || contentType.includes('text/json')) {
        try {
          if (typeof data === 'string') {
            return JSON.parse(data);
          }
          return data;
        } catch (e) {
          console.warn('Failed to parse JSON response:', e);
          return data;
        }
      }
      
      return data;
    }
  ]
});

// Log the base URL for debugging (only in development)
if (import.meta.env.DEV) {
  console.log('🔍 API Configuration:');
  console.log('  Base URL:', apiClient.defaults.baseURL);
  console.log('  Environment mode:', import.meta.env.MODE);
}

// Add request interceptor to attach token
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token && !config.url.includes('/auth/login') && !config.url.includes('/auth/register')) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    if (import.meta.env.DEV) {
      console.log(`🚀 ${config.method.toUpperCase()} ${config.url}`, config.data);
    }
    
    return config;
  },
  error => {
    console.error('❌ Request Error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
apiClient.interceptors.response.use(
  response => {
    if (import.meta.env.DEV) {
      console.log(`✅ ${response.status} ${response.config.url}`, response.data);
    }
    return response;
  },
  error => {
    if (error.response) {
      const { status, data, headers } = error.response;    
      const url = error.config?.url || '';
      
      const contentType = headers?.['content-type'] || headers?.['Content-Type'] || '';
      const isNonJsonResponse = typeof data === 'string' || 
                               contentType.includes('text/html') || 
                               contentType.includes('text/xml');
      
      if (import.meta.env.DEV) {
        console.error(`❌ ${status} ${url}`, isNonJsonResponse ? `[Non-JSON response]` : data);
      }    
      
      // Handle unauthorized - redirect to login
      if (status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/#/auth/login';
        return Promise.reject(new Error('Session expired. Please log in again.'));
      }    
      
      // Handle forbidden
      if (status === 403) {
        const backendMessage = (typeof data === 'string') ? data : (data?.message || data?.title);
        return Promise.reject(new Error(backendMessage || 'Access denied. Insufficient permissions.'));
      }    
      
      // Handle server errors (500+)
      if (status >= 500) {
        let errorMessage = `Server error (${status}).`;
        if (data?.message) {
          errorMessage = data.message;
        }
        error.userMessage = errorMessage;
        return Promise.reject(error);
      }    
      
      const message = data?.message || data?.title || data?.error || 'An error occurred';
      error.userMessage = message;
      return Promise.reject(error);
    }   
    
    // Network errors
    if (error.request) {
      const isCorsError = !error.response && 
                         (error.message?.includes('CORS') || 
                          error.code === 'ERR_NETWORK');
      
      if (isCorsError) {
        const frontendOrigin = window.location.origin;
        const backendUrl = apiClient.defaults.baseURL;
        
        let errorMessage = `CORS Configuration Error: The backend API at ${backendUrl} is not configured to allow requests from ${frontendOrigin}.`;
        
        console.error('🚫 CORS Error:', { frontendOrigin, backendUrl });
        
        const corsError = new Error(errorMessage);
        corsError.isCorsError = true;
        return Promise.reject(corsError);
      }
      
      return Promise.reject(new Error('Network error. Please check your connection.'));
    }
     
    return Promise.reject(new Error(error.message || 'An unexpected error occurred'));
  }
);

// ==================== AUTH SERVICE ====================
export const authService = {
  async login(credentials) {
    const response = await apiClient.post('/auth/login', credentials);
    const responseData = response.data?.data || response.data;

    if (responseData.token) {
      localStorage.setItem('token', responseData.token);
      
      const userData = { ...responseData.user };
      if (userData.role && !userData.roles) {
        userData.roles = [userData.role];
      } else if (!userData.roles) {
        userData.roles = [];
      }
      
      localStorage.setItem('user', JSON.stringify(userData));
    }
    
    return responseData;
  },

  async register(userData) {
    const response = await apiClient.post('/auth/register', userData);
    const responseData = response.data?.data || response.data;

    if (responseData.token) {
      localStorage.setItem('token', responseData.token);
      
      const user = { ...responseData.user };
      if (user.role && !user.roles) {
        user.roles = [user.role];
      } else if (!user.roles) {
        user.roles = [];
      }
      
      localStorage.setItem('user', JSON.stringify(user));
    }
    
    return responseData;
  },

  async getProfile() {
    const response = await apiClient.get('/auth/profile');
    return response.data?.data || response.data;
  },

  async updateProfile(profileData) {
    const response = await apiClient.put('/auth/profile', profileData);
    return response.data?.data || response.data;
  },

  async changePassword(passwordData) {
    const response = await apiClient.put('/auth/change-password', passwordData);
    return response.data?.data || response.data;
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('roles');
    localStorage.removeItem('permissions');
    window.location.href = '/#/auth/login';
  }
};

// ==================== HOUSE SERVICE ====================
export const houseService = {
  async getAll(params = {}) {
    const response = await apiClient.get('/houses', { params });
    return response.data?.data || response.data;
  },

  async getById(id) {
    const response = await apiClient.get(`/houses/${id}`);
    return response.data?.data || response.data;
  },

  async create(houseData) {
    const response = await apiClient.post('/houses', houseData);
    return response.data?.data || response.data;
  },

  async update(id, houseData) {
    const response = await apiClient.put(`/houses/${id}`, houseData);
    return response.data?.data || response.data;
  },

  async delete(id) {
    const response = await apiClient.delete(`/houses/${id}`);
    return response.data?.data || response.data;
  },

  async uploadImage(houseId, imageFile) {
    const formData = new FormData();
    formData.append('image', imageFile);
    
    const response = await apiClient.post(`/houses/${houseId}/images`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data?.data || response.data;
  },

  async deleteImage(imageId) {
    const response = await apiClient.delete(`/houses/images/${imageId}`);
    return response.data?.data || response.data;
  }
};

// ==================== RENTAL SERVICE ====================
export const rentalService = {
  async getAll(params = {}) {
    const response = await apiClient.get('/rentals', { params });
    return response.data?.data || response.data;
  },

  async getById(id) {
    const response = await apiClient.get(`/rentals/${id}`);
    return response.data?.data || response.data;
  },

  async create(rentalData) {
    const response = await apiClient.post('/rentals', rentalData);
    return response.data?.data || response.data;
  },

  async update(id, rentalData) {
    const response = await apiClient.put(`/rentals/${id}`, rentalData);
    return response.data?.data || response.data;
  },

  async terminate(id) {
    const response = await apiClient.put(`/rentals/${id}/terminate`);
    return response.data?.data || response.data;
  }
};

// ==================== PAYMENT SERVICE ====================
export const paymentService = {
  async getAll(params = {}) {
    const response = await apiClient.get('/payments', { params });
    return response.data?.data || response.data;
  },

  async getById(id) {
    const response = await apiClient.get(`/payments/${id}`);
    return response.data?.data || response.data;
  },

  async process(paymentData) {
    const response = await apiClient.post('/payments', paymentData);
    return response.data?.data || response.data;
  },

  async getStats() {
    const response = await apiClient.get('/payments/stats');
    return response.data?.data || response.data;
  }
};

// ==================== REVIEW SERVICE ====================
export const reviewService = {
  async getByHouse(houseId, params = {}) {
    const response = await apiClient.get(`/houses/${houseId}/reviews`, { params });
    return response.data?.data || response.data;
  },

  async getMyReviews() {
    const response = await apiClient.get('/reviews/my');
    return response.data?.data || response.data;
  },

  async create(reviewData) {
    const response = await apiClient.post('/reviews', reviewData);
    return response.data?.data || response.data;
  },

  async update(id, reviewData) {
    const response = await apiClient.put(`/reviews/${id}`, reviewData);
    return response.data?.data || response.data;
  },

  async delete(id) {
    const response = await apiClient.delete(`/reviews/${id}`);
    return response.data?.data || response.data;
  }
};

// ==================== MAINTENANCE SERVICE ====================
export const maintenanceService = {
  async getAll(params = {}) {
    const response = await apiClient.get('/maintenance', { params });
    return response.data?.data || response.data;
  },

  async getById(id) {
    const response = await apiClient.get(`/maintenance/${id}`);
    return response.data?.data || response.data;
  },

  async create(maintenanceData) {
    const response = await apiClient.post('/maintenance', maintenanceData);
    return response.data?.data || response.data;
  },

  async update(id, maintenanceData) {
    const response = await apiClient.put(`/maintenance/${id}`, maintenanceData);
    return response.data?.data || response.data;
  },

  async getStats() {
    const response = await apiClient.get('/maintenance/stats');
    return response.data?.data || response.data;
  }
};

// ==================== FAVORITE SERVICE ====================
export const favoriteService = {
  async getAll() {
    const response = await apiClient.get('/favorites');
    return response.data?.data || response.data;
  },

  async add(houseId) {
    const response = await apiClient.post(`/favorites/${houseId}`);
    return response.data?.data || response.data;
  },

  async remove(houseId) {
    const response = await apiClient.delete(`/favorites/${houseId}`);
    return response.data?.data || response.data;
  },

  async check(houseId) {
    const response = await apiClient.get(`/favorites/${houseId}/check`);
    return response.data?.data || response.data;
  }
};

// ==================== NOTIFICATION SERVICE ====================
export const notificationService = {
  async getAll(params = {}) {
    const response = await apiClient.get('/notifications', { params });
    return response.data?.data || response.data;
  },

  async getById(id) {
    const response = await apiClient.get(`/notifications/${id}`);
    return response.data?.data || response.data;
  },

  async markAsRead(id) {
    const response = await apiClient.put(`/notifications/${id}/read`);
    return response.data?.data || response.data;
  },

  async markAllAsRead() {
    const response = await apiClient.put('/notifications/read-all');
    return response.data?.data || response.data;
  },

  async delete(id) {
    const response = await apiClient.delete(`/notifications/${id}`);
    return response.data?.data || response.data;
  },

  async getStats() {
    const response = await apiClient.get('/notifications/stats');
    return response.data?.data || response.data;
  }
};

// ==================== ADMIN SERVICE ====================
export const adminService = {
  async getDashboardStats() {
    const response = await apiClient.get('/admin/dashboard');
    return response.data?.data || response.data;
  },

  async getUsers(params = {}) {
    const response = await apiClient.get('/admin/users', { params });
    return response.data?.data || response.data;
  },

  async updateUserStatus(userId, statusData) {
    const response = await apiClient.put(`/admin/users/${userId}/status`, statusData);
    return response.data?.data || response.data;
  },

  async getReports(params = {}) {
    const response = await apiClient.get('/admin/reports', { params });
    return response.data?.data || response.data;
  }
};

// ==================== UTILITY FUNCTIONS ====================
export const getCurrentUser = () => {
  try {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  } catch {
    return null;
  }
};

export const isAuthenticated = () => {
  const token = localStorage.getItem('token');
  const user = getCurrentUser();
  return !!(token && user);
};

export const hasRole = (role) => {
  const user = getCurrentUser();
  if (!user) return false;
  
  if (user.roles && Array.isArray(user.roles)) {
    return user.roles.includes(role);
  }
  
  return user.role === role;
};

export const hasAnyRole = (roles) => {
  const user = getCurrentUser();
  if (!user || !Array.isArray(roles)) return false;
  
  if (user.roles && Array.isArray(user.roles)) {
    return user.roles.some(userRole => roles.includes(userRole));
  }
  
  return roles.includes(user.role);
};

export default apiClient;
