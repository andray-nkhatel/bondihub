export interface User {
  id: string;
  full_name: string;
  email: string;
  phone: string;
  role: UserRole;
  is_active: boolean;
  is_verified: boolean;
  profile_image?: string;
  subscription_plan: SubscriptionPlan;
  plan_expiry_date?: string;
  created_at: string;
  updated_at: string;
}

export type UserRole = 'landlord' | 'tenant' | 'agent' | 'admin';

export type SubscriptionPlan = 'basic' | 'premium' | 'enterprise';

export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  full_name: string;
  email: string;
  password: string;
  phone: string;
  role: UserRole;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface UpdateProfileRequest {
  full_name?: string;
  phone?: string;
  profile_image?: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}
