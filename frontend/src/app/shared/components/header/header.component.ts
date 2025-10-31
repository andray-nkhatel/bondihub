import { Component, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { Subject, takeUntil } from 'rxjs';
import { AuthService } from '../../../core/services/auth.service';
import { User } from '../../../core/models/user.model';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  template: `
    <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center">
            <a routerLink="/" class="flex items-center space-x-2">
              <img src="assets/logo.svg" alt="BondiHub" class="h-8 w-8">
              <span class="text-xl font-bold text-primary-600">BondiHub</span>
            </a>
          </div>

          <!-- Desktop Navigation -->
          <nav class="hidden md:flex items-center space-x-8">
            <a routerLink="/houses" 
               class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
              Browse Houses
            </a>
            <a routerLink="/about" 
               class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
              About
            </a>
            <a routerLink="/contact" 
               class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
              Contact
            </a>
          </nav>

          <!-- User Menu -->
          <div class="flex items-center space-x-4">
            <ng-container *ngIf="!isAuthenticated; else authenticatedMenu">
              <!-- Guest Menu -->
              <a routerLink="/login" 
                 class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
                Login
              </a>
              <a routerLink="/register" 
                 class="btn-primary">
                Sign Up
              </a>
            </ng-container>

            <ng-template #authenticatedMenu>
              <!-- Notifications -->
              <button class="relative p-2 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-full">
                <i class="pi pi-bell text-lg"></i>
                <span class="absolute -top-1 -right-1 h-4 w-4 bg-danger-500 text-white text-xs rounded-full flex items-center justify-center">
                  3
                </span>
              </button>

              <!-- User Dropdown -->
              <div class="relative" #userDropdown>
                <button (click)="toggleUserMenu()" 
                        class="flex items-center space-x-2 text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-primary-500">
                  <img [src]="user?.profile_image || 'assets/default-avatar.png'" 
                       [alt]="user?.full_name" 
                       class="h-8 w-8 rounded-full object-cover">
                  <span class="hidden md:block text-gray-700 font-medium">{{ user?.full_name }}</span>
                  <i class="pi pi-chevron-down text-gray-400"></i>
                </button>

                <!-- Dropdown Menu -->
                <div *ngIf="showUserMenu" 
                     class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200">
                  <a routerLink="/profile" 
                     (click)="closeUserMenu()"
                     class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    <i class="pi pi-user mr-2"></i>Profile
                  </a>
                  <a routerLink="/dashboard" 
                     (click)="closeUserMenu()"
                     class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    <i class="pi pi-home mr-2"></i>Dashboard
                  </a>
                  <a routerLink="/notifications" 
                     (click)="closeUserMenu()"
                     class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    <i class="pi pi-bell mr-2"></i>Notifications
                  </a>
                  <div class="border-t border-gray-100"></div>
                  <button (click)="onLogout()" 
                          class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    <i class="pi pi-sign-out mr-2"></i>Logout
                  </button>
                </div>
              </div>
            </ng-template>

            <!-- Mobile Menu Button -->
            <button (click)="toggleMobileMenu()" 
                    class="md:hidden p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-primary-500">
              <i class="pi pi-bars text-lg"></i>
            </button>
          </div>
        </div>

        <!-- Mobile Menu -->
        <div *ngIf="showMobileMenu" class="md:hidden">
          <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-white border-t border-gray-200">
            <a routerLink="/houses" 
               (click)="closeMobileMenu()"
               class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-primary-600 hover:bg-gray-50">
              Browse Houses
            </a>
            <a routerLink="/about" 
               (click)="closeMobileMenu()"
               class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-primary-600 hover:bg-gray-50">
              About
            </a>
            <a routerLink="/contact" 
               (click)="closeMobileMenu()"
               class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-primary-600 hover:bg-gray-50">
              Contact
            </a>
          </div>
        </div>
      </div>
    </header>
  `,
  styles: [`
    .dropdown-enter {
      animation: dropdownEnter 0.2s ease-out;
    }

    @keyframes dropdownEnter {
      from {
        opacity: 0;
        transform: translateY(-10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
  `]
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() isAuthenticated = false;
  @Input() user: User | null = null;
  @Output() logout = new EventEmitter<void>();
  @Output() toggleSidebar = new EventEmitter<void>();

  showUserMenu = false;
  showMobileMenu = false;
  
  private destroy$ = new Subject<void>();

  constructor(private authService: AuthService) {}

  ngOnInit(): void {
    // Listen for clicks outside to close dropdown
    document.addEventListener('click', this.handleClickOutside.bind(this));
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    document.removeEventListener('click', this.handleClickOutside.bind(this));
  }

  toggleUserMenu(): void {
    this.showUserMenu = !this.showUserMenu;
  }

  closeUserMenu(): void {
    this.showUserMenu = false;
  }

  toggleMobileMenu(): void {
    this.showMobileMenu = !this.showMobileMenu;
  }

  closeMobileMenu(): void {
    this.showMobileMenu = false;
  }

  onLogout(): void {
    this.closeUserMenu();
    this.logout.emit();
  }

  private handleClickOutside(event: Event): void {
    const target = event.target as HTMLElement;
    const userDropdown = document.querySelector('[#userDropdown]');
    
    if (userDropdown && !userDropdown.contains(target)) {
      this.closeUserMenu();
    }
  }
}
