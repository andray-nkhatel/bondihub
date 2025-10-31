import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { NotificationService } from '../services/notification.service';

export const roleGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  const notificationService = inject(NotificationService);

  // Check if user is authenticated
  if (!authService.isAuthenticated()) {
    notificationService.showWarning('Please log in to access this page');
    router.navigate(['/login'], { queryParams: { returnUrl: state.url } });
    return false;
  }

  // Get required roles from route data
  const requiredRoles = route.data?.['roles'] as string[];
  
  if (!requiredRoles || requiredRoles.length === 0) {
    return true; // No role requirements
  }

  // Check if user has any of the required roles
  const hasRequiredRole = authService.hasAnyRole(requiredRoles);

  if (!hasRequiredRole) {
    notificationService.showError('You do not have permission to access this page');
    router.navigate(['/403']);
    return false;
  }

  return true;
};
