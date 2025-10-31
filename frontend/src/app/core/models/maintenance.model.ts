export interface MaintenanceRequest {
  id: string;
  tenant_id: string;
  house_id: string;
  title: string;
  description: string;
  status: MaintenanceRequestStatus;
  priority: string;
  reported_at: string;
  resolved_at?: string;
  created_at: string;
  updated_at: string;
  tenant?: User;
  house?: House;
}

export type MaintenanceRequestStatus = 'pending' | 'in_progress' | 'resolved' | 'cancelled';

export interface CreateMaintenanceRequest {
  house_id: string;
  title: string;
  description: string;
  priority: string;
}

export interface UpdateMaintenanceRequest {
  status?: MaintenanceRequestStatus;
}

export interface MaintenanceRequestFilters {
  status?: string;
  priority?: string;
  page?: number;
  limit?: number;
}

export interface MaintenanceRequestListResponse {
  maintenance_requests: MaintenanceRequest[];
  pagination: PaginationInfo;
}

export interface MaintenanceStats {
  total_requests: number;
  requests_by_status: StatusStats[];
  requests_by_priority: PriorityStats[];
  avg_resolution_days: number;
}

export interface StatusStats {
  status: string;
  count: number;
}

export interface PriorityStats {
  priority: string;
  count: number;
}

// Import related types
import { User } from './user.model';
import { House } from './house.model';
import { PaginationInfo } from './house.model';
