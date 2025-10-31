import { Component, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router, NavigationEnd } from '@angular/router';
import { Subject, takeUntil, filter } from 'rxjs';
import { User } from '../../../core/models/user.model';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <aside [class]="sidebarClasses">
      <div class="flex flex-col h-full">
        <!-- Logo -->
        <div class="flex items-center justify-between p-4 border-b border-gray-200">
          <div class="flex items-center space-x-2" *ngIf="!isCollapsed">
            <img src="assets/logo.svg" alt="BondiHub" class="h-6 w-6">
            <span class="text-lg font-bold text-primary-600">BondiHub</span>
          </div>
          <img *ngIf="isCollapsed" src="assets/logo.svg" alt="BondiHub" class="h-6 w-6">
          
          <!-- Close button for mobile -->
          <button *ngIf="isMobile" 
                  (click)="onClose()"
                  class="p-1 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <i class="pi pi-times text-lg"></i>
          </button>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 p-4 space-y-2">
          <!-- Dashboard -->
          <a routerLink="/dashboard" 
             [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
             class="nav-item">
            <i class="pi pi-home text-lg"></i>
            <span *ngIf="!isCollapsed">Dashboard</span>
          </a>

          <!-- Houses -->
          <a routerLink="/houses" 
             [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
             class="nav-item">
            <i class="pi pi-building text-lg"></i>
            <span *ngIf="!isCollapsed">Browse Houses</span>
          </a>

          <!-- Role-specific navigation -->
          <ng-container *ngIf="isAuthenticated && user">
            <!-- Landlord/Admin Navigation -->
            <ng-container *ngIf="user.role === 'landlord' || user.role === 'admin'">
              <div class="pt-4" *ngIf="!isCollapsed">
                <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                  Property Management
                </h3>
              </div>
              
              <a routerLink="/landlord/houses" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-building text-lg"></i>
                <span *ngIf="!isCollapsed">My Properties</span>
              </a>
              
              <a routerLink="/landlord/agreements" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-file-text text-lg"></i>
                <span *ngIf="!isCollapsed">Agreements</span>
              </a>
              
              <a routerLink="/landlord/payments" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-credit-card text-lg"></i>
                <span *ngIf="!isCollapsed">Payments</span>
              </a>
              
              <a routerLink="/landlord/maintenance" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-wrench text-lg"></i>
                <span *ngIf="!isCollapsed">Maintenance</span>
              </a>
            </ng-container>

            <!-- Tenant Navigation -->
            <ng-container *ngIf="user.role === 'tenant' || user.role === 'admin'">
              <div class="pt-4" *ngIf="!isCollapsed">
                <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                  Tenant Services
                </h3>
              </div>
              
              <a routerLink="/tenant/favorites" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-heart text-lg"></i>
                <span *ngIf="!isCollapsed">Favorites</span>
              </a>
              
              <a routerLink="/tenant/agreements" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-file-text text-lg"></i>
                <span *ngIf="!isCollapsed">My Agreements</span>
              </a>
              
              <a routerLink="/tenant/payments" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-credit-card text-lg"></i>
                <span *ngIf="!isCollapsed">Payments</span>
              </a>
              
              <a routerLink="/tenant/maintenance" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-wrench text-lg"></i>
                <span *ngIf="!isCollapsed">Maintenance</span>
              </a>
              
              <a routerLink="/tenant/reviews" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-star text-lg"></i>
                <span *ngIf="!isCollapsed">My Reviews</span>
              </a>
            </ng-container>

            <!-- Admin Navigation -->
            <ng-container *ngIf="user.role === 'admin'">
              <div class="pt-4" *ngIf="!isCollapsed">
                <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                  Administration
                </h3>
              </div>
              
              <a routerLink="/admin/users" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-users text-lg"></i>
                <span *ngIf="!isCollapsed">Users</span>
              </a>
              
              <a routerLink="/admin/houses" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-building text-lg"></i>
                <span *ngIf="!isCollapsed">All Houses</span>
              </a>
              
              <a routerLink="/admin/payments" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-credit-card text-lg"></i>
                <span *ngIf="!isCollapsed">All Payments</span>
              </a>
              
              <a routerLink="/admin/reports" 
                 [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
                 class="nav-item">
                <i class="pi pi-chart-bar text-lg"></i>
                <span *ngIf="!isCollapsed">Reports</span>
              </a>
            </ng-container>
          </ng-container>

          <!-- Common Navigation -->
          <div class="pt-4" *ngIf="!isCollapsed && isAuthenticated">
            <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              Account
            </h3>
          </div>
          
          <a routerLink="/profile" 
             [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
             class="nav-item">
            <i class="pi pi-user text-lg"></i>
            <span *ngIf="!isCollapsed">Profile</span>
          </a>
          
          <a routerLink="/notifications" 
             [routerLinkActive]="'bg-primary-50 text-primary-700 border-r-2 border-primary-600'"
             class="nav-item">
            <i class="pi pi-bell text-lg"></i>
            <span *ngIf="!isCollapsed">Notifications</span>
            <span *ngIf="!isCollapsed" class="ml-auto bg-danger-500 text-white text-xs rounded-full px-2 py-1">
              3
            </span>
          </a>
        </nav>

        <!-- User Info (Desktop) -->
        <div *ngIf="isAuthenticated && user && !isCollapsed" 
             class="p-4 border-t border-gray-200">
          <div class="flex items-center space-x-3">
            <img [src]="user.profile_image || 'assets/default-avatar.png'" 
                 [alt]="user.full_name" 
                 class="h-8 w-8 rounded-full object-cover">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate">
                {{ user.full_name }}
              </p>
              <p class="text-xs text-gray-500 capitalize">
                {{ user.role }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </aside>
  `,
  styles: [`
    .nav-item {
      @apply flex items-center space-x-3 px-3 py-2 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-100 transition-colors duration-200;
    }
    
    .sidebar-collapsed {
      @apply w-16;
    }
    
    .sidebar-expanded {
      @apply w-64;
    }
    
    .sidebar-mobile {
      @apply fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg transform transition-transform duration-300 ease-in-out;
    }
  `]
})
export class SidebarComponent implements OnInit, OnDestroy {
  @Input() isAuthenticated = false;
  @Input() user: User | null = null;
  @Input() isCollapsed = false;
  @Input() isMobile = false;
  @Output() close = new EventEmitter<void>();

  private destroy$ = new Subject<void>();

  constructor(private router: Router) {}

  ngOnInit(): void {
    // Listen for route changes to close mobile sidebar
    this.router.events
      .pipe(
        filter(event => event instanceof NavigationEnd),
        takeUntil(this.destroy$)
      )
      .subscribe(() => {
        if (this.isMobile) {
          this.onClose();
        }
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  get sidebarClasses(): string {
    if (this.isMobile) {
      return 'sidebar-mobile';
    }
    return this.isCollapsed ? 'sidebar-collapsed' : 'sidebar-expanded';
  }

  onClose(): void {
    this.close.emit();
  }
}
