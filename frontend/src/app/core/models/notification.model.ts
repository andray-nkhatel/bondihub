export interface Notification {
  id: string;
  user_id: string;
  title: string;
  message: string;
  is_read: boolean;
  type: NotificationType;
  created_at: string;
}

export type NotificationType = 'payment' | 'maintenance' | 'agreement' | 'review' | 'general';

export interface NotificationFilters {
  unread_only?: boolean;
  type?: string;
  page?: number;
  limit?: number;
}

export interface NotificationListResponse {
  notifications: Notification[];
  pagination: PaginationInfo;
}

export interface NotificationStats {
  total_notifications: number;
  unread_notifications: number;
  notifications_by_type: NotificationTypeStats[];
}

export interface NotificationTypeStats {
  type: string;
  count: number;
}

// Import related types
import { PaginationInfo } from './house.model';
