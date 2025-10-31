import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, of, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { ApiService } from './api.service';
import { 
  User, 
  AuthState, 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse,
  UpdateProfileRequest,
  ChangePasswordRequest
} from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private authStateSubject = new BehaviorSubject<AuthState>({
    isAuthenticated: false,
    user: null,
    token: null
  });

  public authState$ = this.authStateSubject.asObservable();

  constructor(private apiService: ApiService) {
    this.initializeAuth();
  }

  private initializeAuth(): void {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    
    if (token && userStr) {
      try {
        const user = JSON.parse(userStr);
        this.authStateSubject.next({
          isAuthenticated: true,
          user,
          token
        });
      } catch (error) {
        this.clearAuthData();
      }
    }
  }

  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.apiService.post<AuthResponse>('/auth/login', credentials).pipe(
      tap(response => {
        this.setAuthData(response.user, response.token);
      }),
      catchError(error => {
        this.clearAuthData();
        return throwError(() => error);
      })
    );
  }

  register(userData: RegisterRequest): Observable<AuthResponse> {
    return this.apiService.post<AuthResponse>('/auth/register', userData).pipe(
      tap(response => {
        this.setAuthData(response.user, response.token);
      }),
      catchError(error => {
        this.clearAuthData();
        return throwError(() => error);
      })
    );
  }

  logout(): Observable<any> {
    return this.apiService.post('/auth/logout', {}).pipe(
      tap(() => {
        this.clearAuthData();
      }),
      catchError(error => {
        // Clear auth data even if logout request fails
        this.clearAuthData();
        return of(null);
      })
    );
  }

  getProfile(): Observable<User> {
    return this.apiService.get<User>('/auth/profile');
  }

  updateProfile(profileData: UpdateProfileRequest): Observable<User> {
    return this.apiService.put<User>('/auth/profile', profileData).pipe(
      tap(user => {
        this.updateUserInState(user);
      })
    );
  }

  changePassword(passwordData: ChangePasswordRequest): Observable<any> {
    return this.apiService.put('/auth/change-password', passwordData);
  }

  checkAuthStatus(): Observable<AuthState> {
    const token = localStorage.getItem('token');
    
    if (!token) {
      return of({
        isAuthenticated: false,
        user: null,
        token: null
      });
    }

    return this.getProfile().pipe(
      map(user => {
        const authState = {
          isAuthenticated: true,
          user,
          token
        };
        this.authStateSubject.next(authState);
        return authState;
      }),
      catchError(error => {
        this.clearAuthData();
        return of({
          isAuthenticated: false,
          user: null,
          token: null
        });
      })
    );
  }

  isAuthenticated(): boolean {
    return this.authStateSubject.value.isAuthenticated;
  }

  getCurrentUser(): User | null {
    return this.authStateSubject.value.user;
  }

  getToken(): string | null {
    return this.authStateSubject.value.token;
  }

  hasRole(role: string): boolean {
    const user = this.getCurrentUser();
    return user ? user.role === role : false;
  }

  hasAnyRole(roles: string[]): boolean {
    const user = this.getCurrentUser();
    return user ? roles.includes(user.role) : false;
  }

  isLandlord(): boolean {
    return this.hasRole('landlord');
  }

  isTenant(): boolean {
    return this.hasRole('tenant');
  }

  isAdmin(): boolean {
    return this.hasRole('admin');
  }

  isAgent(): boolean {
    return this.hasRole('agent');
  }

  private setAuthData(user: User, token: string): void {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    
    this.authStateSubject.next({
      isAuthenticated: true,
      user,
      token
    });
  }

  private clearAuthData(): void {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    this.authStateSubject.next({
      isAuthenticated: false,
      user: null,
      token: null
    });
  }

  private updateUserInState(user: User): void {
    const currentState = this.authStateSubject.value;
    this.authStateSubject.next({
      ...currentState,
      user
    });
    localStorage.setItem('user', JSON.stringify(user));
  }

  // Utility methods for role-based access
  canManageHouses(): boolean {
    return this.hasAnyRole(['landlord', 'admin']);
  }

  canMakePayments(): boolean {
    return this.hasAnyRole(['tenant', 'admin']);
  }

  canViewAdminPanel(): boolean {
    return this.hasRole('admin');
  }

  canCreateAgreements(): boolean {
    return this.hasAnyRole(['landlord', 'admin']);
  }

  canRequestMaintenance(): boolean {
    return this.hasAnyRole(['tenant', 'admin']);
  }
}
