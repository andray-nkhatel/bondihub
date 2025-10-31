import { Injectable } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { 
  Notification, 
  NotificationFilters, 
  NotificationListResponse,
  NotificationStats 
} from '../models/notification.model';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  constructor(
    private toastr: ToastrService,
    private apiService: ApiService
  ) {}

  // Toast notifications
  showSuccess(message: string, title?: string): void {
    this.toastr.success(message, title || 'Success');
  }

  showError(message: string, title?: string): void {
    this.toastr.error(message, title || 'Error');
  }

  showWarning(message: string, title?: string): void {
    this.toastr.warning(message, title || 'Warning');
  }

  showInfo(message: string, title?: string): void {
    this.toastr.info(message, title || 'Information');
  }

  // API notifications
  getNotifications(filters?: NotificationFilters): Observable<NotificationListResponse> {
    return this.apiService.get<NotificationListResponse>('/notifications', filters);
  }

  getNotification(id: string): Observable<Notification> {
    return this.apiService.get<Notification>(`/notifications/${id}`);
  }

  markAsRead(id: string): Observable<Notification> {
    return this.apiService.put<Notification>(`/notifications/${id}/read`, {});
  }

  markAllAsRead(): Observable<any> {
    return this.apiService.put('/notifications/read-all', {});
  }

  deleteNotification(id: string): Observable<any> {
    return this.apiService.delete(`/notifications/${id}`);
  }

  getNotificationStats(): Observable<NotificationStats> {
    return this.apiService.get<NotificationStats>('/notifications/stats');
  }

  // Utility methods
  formatNotificationMessage(notification: Notification): string {
    return notification.message;
  }

  getNotificationIcon(type: string): string {
    const icons: { [key: string]: string } = {
      payment: 'pi pi-credit-card',
      maintenance: 'pi pi-wrench',
      agreement: 'pi pi-file-text',
      review: 'pi pi-star',
      general: 'pi pi-bell'
    };
    return icons[type] || 'pi pi-bell';
  }

  getNotificationColor(type: string): string {
    const colors: { [key: string]: string } = {
      payment: 'text-green-600',
      maintenance: 'text-orange-600',
      agreement: 'text-blue-600',
      review: 'text-yellow-600',
      general: 'text-gray-600'
    };
    return colors[type] || 'text-gray-600';
  }

  // Real-time notifications (WebSocket implementation would go here)
  connectToNotifications(): void {
    // Implementation for WebSocket connection
    console.log('Connecting to real-time notifications...');
  }

  disconnectFromNotifications(): void {
    // Implementation for WebSocket disconnection
    console.log('Disconnecting from real-time notifications...');
  }
}
