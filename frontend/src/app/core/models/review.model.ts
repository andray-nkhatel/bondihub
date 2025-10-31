export interface Review {
  id: string;
  tenant_id: string;
  house_id: string;
  rating: number;
  comment: string;
  created_at: string;
  updated_at: string;
  tenant?: User;
  house?: House;
}

export interface CreateReviewRequest {
  house_id: string;
  rating: number;
  comment: string;
}

export interface UpdateReviewRequest {
  rating?: number;
  comment?: string;
}

export interface ReviewFilters {
  page?: number;
  limit?: number;
}

export interface ReviewListResponse {
  reviews: Review[];
  pagination: PaginationInfo;
  average_rating: number;
  rating_distribution: RatingDistribution[];
}

export interface RatingDistribution {
  rating: number;
  count: number;
}

// Import related types
import { User } from './user.model';
import { House } from './house.model';
import { PaginationInfo } from './house.model';
