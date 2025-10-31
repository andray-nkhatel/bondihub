import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet],
  template: `
    <div class="min-h-screen bg-gray-50">
      <header class="bg-white shadow-sm border-b">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between items-center h-16">
            <div class="flex items-center">
              <h1 class="text-2xl font-bold text-primary-600">BondiHub</h1>
              <span class="ml-2 text-sm text-gray-500">Zambia</span>
            </div>
            <nav class="hidden md:flex space-x-8">
              <a href="#" class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium">Houses</a>
              <a href="#" class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium">About</a>
              <a href="#" class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium">Contact</a>
              <a href="#" class="bg-primary-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-primary-700">Login</a>
            </nav>
          </div>
        </div>
      </header>

      <main>
        <router-outlet></router-outlet>
      </main>

      <footer class="bg-gray-800 text-white">
        <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
          <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <h3 class="text-lg font-semibold mb-4">BondiHub</h3>
              <p class="text-gray-300 text-sm">Your trusted partner for house renting in Zambia.</p>
            </div>
            <div>
              <h4 class="text-md font-semibold mb-4">Quick Links</h4>
              <ul class="space-y-2 text-sm">
                <li><a href="#" class="text-gray-300 hover:text-white">Find Houses</a></li>
                <li><a href="#" class="text-gray-300 hover:text-white">List Property</a></li>
                <li><a href="#" class="text-gray-300 hover:text-white">How it Works</a></li>
              </ul>
            </div>
            <div>
              <h4 class="text-md font-semibold mb-4">Support</h4>
              <ul class="space-y-2 text-sm">
                <li><a href="#" class="text-gray-300 hover:text-white">Help Center</a></li>
                <li><a href="#" class="text-gray-300 hover:text-white">Contact Us</a></li>
                <li><a href="#" class="text-gray-300 hover:text-white">Privacy Policy</a></li>
              </ul>
            </div>
            <div>
              <h4 class="text-md font-semibold mb-4">Connect</h4>
              <ul class="space-y-2 text-sm">
                <li><a href="#" class="text-gray-300 hover:text-white">Facebook</a></li>
                <li><a href="#" class="text-gray-300 hover:text-white">Twitter</a></li>
                <li><a href="#" class="text-gray-300 hover:text-white">Instagram</a></li>
              </ul>
            </div>
          </div>
          <div class="mt-8 pt-8 border-t border-gray-700">
            <p class="text-center text-gray-300 text-sm">&copy; 2024 BondiHub. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  `,
  styles: []
})
export class AppComponent {
  title = 'BondiHub - House Renting Service for Zambia';
}