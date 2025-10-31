import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="bg-gradient-to-r from-primary-600 to-primary-800 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
        <div class="text-center">
          <h1 class="text-4xl md:text-6xl font-bold mb-6">
            Find Your Perfect Home in Zambia
          </h1>
          <p class="text-xl md:text-2xl mb-8 text-primary-100">
            Discover amazing houses for rent across Zambia with our trusted platform
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <button class="bg-white text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors">
              Browse Houses
            </button>
            <button class="border-2 border-white text-white px-8 py-3 rounded-lg font-semibold hover:bg-white hover:text-primary-600 transition-colors">
              List Your Property
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="py-16 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <h2 class="text-3xl font-bold text-gray-900 mb-4">Why Choose BondiHub?</h2>
          <p class="text-lg text-gray-600">We make house hunting simple and secure</p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="text-center p-6">
            <div class="bg-primary-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Verified Properties</h3>
            <p class="text-gray-600">All our properties are verified and inspected to ensure quality and safety.</p>
          </div>
          
          <div class="text-center p-6">
            <div class="bg-primary-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Secure Payments</h3>
            <p class="text-gray-600">Pay safely with MTN MoMo, Airtel Money, or bank transfers.</p>
          </div>
          
          <div class="text-center p-6">
            <div class="bg-primary-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M12 2.25a9.75 9.75 0 100 19.5 9.75 9.75 0 000-19.5z"></path>
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">24/7 Support</h3>
            <p class="text-gray-600">Get help whenever you need it with our round-the-clock customer support.</p>
          </div>
        </div>
      </div>
    </div>

    <div class="py-16 bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <h2 class="text-3xl font-bold text-gray-900 mb-4">Popular Locations</h2>
          <p class="text-lg text-gray-600">Find houses in Zambia's most popular cities</p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Lusaka</h3>
            <p class="text-gray-600 text-sm">Capital city with modern amenities</p>
          </div>
          <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Ndola</h3>
            <p class="text-gray-600 text-sm">Copperbelt's commercial hub</p>
          </div>
          <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Kitwe</h3>
            <p class="text-gray-600 text-sm">Mining and industrial center</p>
          </div>
          <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Livingstone</h3>
            <p class="text-gray-600 text-sm">Tourist destination near Victoria Falls</p>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: []
})
export class HomeComponent {
}