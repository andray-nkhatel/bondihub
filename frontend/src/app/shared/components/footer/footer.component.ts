import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <footer class="bg-gray-900 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <!-- Company Info -->
          <div class="space-y-4">
            <div class="flex items-center space-x-2">
              <img src="assets/logo-white.svg" alt="BondiHub" class="h-8 w-8">
              <span class="text-xl font-bold">{{ environment.appName }}</span>
            </div>
            <p class="text-gray-300 text-sm">
              Find your perfect home in Zambia with our comprehensive house renting platform. 
              Mobile money payments, local currency support, and mobile-first design.
            </p>
            <div class="flex space-x-4">
              <a [href]="environment.socialMedia.facebook" target="_blank" 
                 class="text-gray-400 hover:text-white transition-colors">
                <i class="pi pi-facebook text-lg"></i>
              </a>
              <a [href]="environment.socialMedia.twitter" target="_blank" 
                 class="text-gray-400 hover:text-white transition-colors">
                <i class="pi pi-twitter text-lg"></i>
              </a>
              <a [href]="environment.socialMedia.instagram" target="_blank" 
                 class="text-gray-400 hover:text-white transition-colors">
                <i class="pi pi-instagram text-lg"></i>
              </a>
              <a [href]="environment.socialMedia.linkedin" target="_blank" 
                 class="text-gray-400 hover:text-white transition-colors">
                <i class="pi pi-linkedin text-lg"></i>
              </a>
            </div>
          </div>

          <!-- Quick Links -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Quick Links</h3>
            <ul class="space-y-2">
              <li>
                <a routerLink="/houses" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Browse Houses
                </a>
              </li>
              <li>
                <a routerLink="/about" class="text-gray-300 hover:text-white transition-colors text-sm">
                  About Us
                </a>
              </li>
              <li>
                <a routerLink="/contact" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Contact
                </a>
              </li>
              <li>
                <a routerLink="/privacy" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Privacy Policy
                </a>
              </li>
              <li>
                <a routerLink="/terms" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Terms of Service
                </a>
              </li>
            </ul>
          </div>

          <!-- For Landlords -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">For Landlords</h3>
            <ul class="space-y-2">
              <li>
                <a routerLink="/landlord" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Landlord Dashboard
                </a>
              </li>
              <li>
                <a routerLink="/landlord/houses" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Manage Properties
                </a>
              </li>
              <li>
                <a routerLink="/landlord/agreements" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Rental Agreements
                </a>
              </li>
              <li>
                <a routerLink="/landlord/payments" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Payment Management
                </a>
              </li>
              <li>
                <a routerLink="/pricing" class="text-gray-300 hover:text-white transition-colors text-sm">
                  Pricing Plans
                </a>
              </li>
            </ul>
          </div>

          <!-- Contact Info -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Contact Info</h3>
            <div class="space-y-3">
              <div class="flex items-center space-x-3">
                <i class="pi pi-envelope text-gray-400"></i>
                <a [href]="'mailto:' + environment.supportEmail" 
                   class="text-gray-300 hover:text-white transition-colors text-sm">
                  {{ environment.supportEmail }}
                </a>
              </div>
              <div class="flex items-center space-x-3">
                <i class="pi pi-phone text-gray-400"></i>
                <a [href]="'tel:' + environment.supportPhone" 
                   class="text-gray-300 hover:text-white transition-colors text-sm">
                  {{ environment.supportPhone }}
                </a>
              </div>
              <div class="flex items-center space-x-3">
                <i class="pi pi-map-marker text-gray-400"></i>
                <span class="text-gray-300 text-sm">
                  Lusaka, Zambia
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Bottom Section -->
        <div class="mt-8 pt-8 border-t border-gray-800">
          <div class="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <div class="text-gray-400 text-sm">
              © {{ currentYear }} {{ environment.appName }}. All rights reserved.
            </div>
            <div class="flex items-center space-x-6 text-sm">
              <span class="text-gray-400">Version {{ environment.appVersion }}</span>
              <div class="flex items-center space-x-2">
                <div class="w-2 h-2 bg-green-500 rounded-full"></div>
                <span class="text-gray-400">All systems operational</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  `,
  styles: [`
    :host {
      display: block;
    }
  `]
})
export class FooterComponent {
  environment = environment;
  currentYear = new Date().getFullYear();
}
