export interface House {
  id: string;
  landlord_id: string;
  title: string;
  description: string;
  address: string;
  monthly_rent: number;
  status: HouseStatus;
  house_type: HouseType;
  latitude: number;
  longitude: number;
  bedrooms: number;
  bathrooms: number;
  area: number;
  is_featured: boolean;
  featured_until?: string;
  created_at: string;
  updated_at: string;
  landlord?: User;
  images?: HouseImage[];
  reviews?: Review[];
  average_rating?: number;
}

export type HouseStatus = 'available' | 'occupied' | 'maintenance';

export type HouseType = 'apartment' | 'house' | 'studio' | 'townhouse' | 'commercial';

export interface HouseImage {
  id: string;
  house_id: string;
  image_url: string;
  is_primary: boolean;
  created_at: string;
}

export interface CreateHouseRequest {
  title: string;
  description: string;
  address: string;
  monthly_rent: number;
  house_type: HouseType;
  latitude: number;
  longitude: number;
  bedrooms: number;
  bathrooms: number;
  area: number;
  is_featured: boolean;
}

export interface UpdateHouseRequest {
  title?: string;
  description?: string;
  address?: string;
  monthly_rent?: number;
  status?: HouseStatus;
  house_type?: HouseType;
  latitude?: number;
  longitude?: number;
  bedrooms?: number;
  bathrooms?: number;
  area?: number;
  is_featured?: boolean;
}

export interface HouseFilters {
  house_type?: string;
  status?: string;
  min_rent?: number;
  max_rent?: number;
  bedrooms?: number;
  bathrooms?: number;
  featured?: boolean;
  search?: string;
  page?: number;
  limit?: number;
}

export interface HouseListResponse {
  houses: House[];
  pagination: PaginationInfo;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

// Import User and Review types
import { User } from './user.model';
import { Review } from './review.model';
