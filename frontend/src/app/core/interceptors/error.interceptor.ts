import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { NotificationService } from '../services/notification.service';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const notificationService = inject(NotificationService);
  const authService = inject(AuthService);
  const router = inject(Router);

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      let errorMessage = 'An unexpected error occurred';

      if (error.error instanceof ErrorEvent) {
        // Client-side error
        errorMessage = error.error.message;
      } else {
        // Server-side error
        switch (error.status) {
          case 400:
            errorMessage = error.error?.message || 'Bad request';
            break;
          case 401:
            errorMessage = 'Unauthorized access';
            authService.logout().subscribe();
            router.navigate(['/login']);
            break;
          case 403:
            errorMessage = 'Access forbidden';
            router.navigate(['/403']);
            break;
          case 404:
            errorMessage = 'Resource not found';
            break;
          case 409:
            errorMessage = error.error?.message || 'Conflict occurred';
            break;
          case 422:
            errorMessage = 'Validation error';
            if (error.error?.errors) {
              const validationErrors = Object.values(error.error.errors).flat();
              errorMessage = validationErrors.join(', ');
            }
            break;
          case 429:
            errorMessage = 'Too many requests. Please try again later';
            break;
          case 500:
            errorMessage = 'Internal server error';
            break;
          case 503:
            errorMessage = 'Service temporarily unavailable';
            break;
          default:
            errorMessage = error.error?.message || `Error ${error.status}: ${error.statusText}`;
        }
      }

      // Show error notification
      notificationService.showError(errorMessage);

      return throwError(() => error);
    })
  );
};
