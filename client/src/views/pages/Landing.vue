<script setup>
import { houseService } from '@/service/api.service';
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';

const router = useRouter();
const store = useStore();

const houses = ref([]);
const loading = ref(true);
const searchQuery = ref('');
const selectedType = ref(null);

const isAuthenticated = computed(() => store.getters['auth/isAuthenticated']);

const propertyTypes = [
    { label: 'All', value: null, icon: 'pi pi-th-large' },
    { label: 'Apartments', value: 'apartment', icon: 'pi pi-building' },
    { label: 'Houses', value: 'house', icon: 'pi pi-home' },
    { label: 'Studios', value: 'studio', icon: 'pi pi-box' },
    { label: 'Townhouses', value: 'townhouse', icon: 'pi pi-server' },
];

const filteredHouses = computed(() => {
    let result = houses.value;
    if (selectedType.value) {
        result = result.filter(h => h.house_type === selectedType.value);
    }
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase();
        result = result.filter(h => 
            h.title?.toLowerCase().includes(query) || 
            h.address?.toLowerCase().includes(query)
        );
    }
    return result;
});

const loadHouses = async () => {
    loading.value = true;
    try {
        const response = await houseService.getAll({ status: 'available', limit: 50 });
        houses.value = response.houses || [];
    } catch (error) {
        console.error('Failed to load houses:', error);
        houses.value = [];
    } finally {
        loading.value = false;
    }
};

const formatPrice = (price) => {
    return new Intl.NumberFormat('en-ZM').format(price);
};

const goToLogin = () => router.push('/auth/login');
const goToRegister = () => router.push('/auth/register');
const goToDashboard = () => router.push('/app/dashboard');

const viewProperty = (house) => {
    if (isAuthenticated.value) {
        router.push('/app/houses');
    } else {
        router.push('/auth/login?redirect=/app/houses');
    }
};

onMounted(() => {
    loadHouses();
});
</script>

<template>
    <div class="min-h-screen bg-surface-0 dark:bg-surface-900">
        <!-- Header -->
        <header class="sticky top-0 z-50 bg-surface-0 dark:bg-surface-900 border-b border-surface-200 dark:border-surface-700">
            <div class="px-6 py-4 flex items-center justify-between max-w-screen-2xl mx-auto">
                <div class="flex items-center gap-2 cursor-pointer" @click="router.push('/')">
                    <!-- <i class="pi pi-home text-2xl text-primary"></i> -->
                    <span class="text-xl font-bold text-surface-900 dark:text-surface-0">BondiHub</span>
                </div>
                
                <!-- Search Bar -->
                <div class="hidden md:flex flex-1 max-w-xl mx-8">
                    <IconField class="w-full">
                        <InputIcon class="pi pi-search" />
                        <InputText v-model="searchQuery" placeholder="Search by location or property name..." class="w-full" />
                    </IconField>
                </div>
                
                <div class="flex items-center gap-3">
                    <template v-if="isAuthenticated">
                        <Button label="Dashboard" @click="goToDashboard" rounded />
                    </template>
                    <template v-else>
                        <Button label="Sign In" @click="goToLogin" text rounded />
                        <Button label="Sign Up" @click="goToRegister" rounded />
                    </template>
                </div>
            </div>
            
            <!-- Mobile Search -->
            <div class="md:hidden px-6 pb-4">
                <IconField class="w-full">
                    <InputIcon class="pi pi-search" />
                    <InputText v-model="searchQuery" placeholder="Search properties..." class="w-full" />
                </IconField>
            </div>
        </header>

        <!-- Property Type Filter -->
        <div class="border-b border-surface-200 dark:border-surface-700 bg-surface-0 dark:bg-surface-900">
            <div class="max-w-screen-2xl mx-auto px-6 py-4">
                <div class="flex gap-6 overflow-x-auto pb-2 justify-center">
                    <button 
                        v-for="type in propertyTypes" 
                        :key="type.value"
                        @click="selectedType = type.value"
                        class="flex flex-col items-center gap-2 min-w-[80px] py-2 px-4 rounded-lg transition-all"
                        :class="selectedType === type.value 
                            ? 'text-primary border-b-2 border-primary' 
                            : 'text-surface-500 hover:text-surface-900 dark:hover:text-surface-0'"
                    >
                        <i :class="type.icon" class="text-xl"></i>
                        <span class="text-sm font-medium whitespace-nowrap">{{ type.label }}</span>
                    </button>
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <main class="max-w-screen-2xl mx-auto px-6 py-8">
            <!-- Loading State -->
            <div v-if="loading" class="flex justify-center py-20">
                <ProgressSpinner />
            </div>

            <!-- Property Grid -->
            <div v-else-if="filteredHouses.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                <div 
                    v-for="house in filteredHouses" 
                    :key="house.id" 
                    class="group cursor-pointer"
                    @click="viewProperty(house)"
                >
                    <!-- Image -->
                    <div class="relative aspect-square rounded-xl overflow-hidden mb-3">
                        <img 
                            v-if="house.images?.length" 
                            :src="house.images[0].image_url" 
                            :alt="house.title" 
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                        />
                        <div v-else class="w-full h-full bg-surface-200 dark:bg-surface-700 flex items-center justify-center">
                            <i class="pi pi-home text-5xl text-surface-400"></i>
                        </div>
                        <!-- Status Badge -->
                        <Tag 
                            :value="house.status" 
                            :severity="house.status === 'available' ? 'success' : 'secondary'"
                            class="absolute top-3 left-3"
                        />
                        <!-- Favorite Button -->
                        <Button 
                            icon="pi pi-heart" 
                            rounded 
                            text 
                            class="absolute top-3 right-3 bg-surface-0/80 hover:bg-surface-0"
                            @click.stop
                        />
                    </div>
                    
                    <!-- Details -->
                    <div class="space-y-1">
                        <div class="flex justify-between items-start">
                            <h3 class="font-semibold text-surface-900 dark:text-surface-0 truncate">{{ house.title }}</h3>
                        </div>
                        <p class="text-surface-500 text-sm flex items-center gap-1">
                            <i class="pi pi-map-marker text-xs"></i>
                            {{ house.address }}
                        </p>
                        <div class="flex items-center gap-3 text-sm text-surface-600 dark:text-surface-400">
                            <span><i class="pi pi-th-large mr-1"></i>{{ house.bedrooms }} bed</span>
                            <span><i class="pi pi-inbox mr-1"></i>{{ house.bathrooms }} bath</span>
                            <span v-if="house.area"><i class="pi pi-arrows-alt mr-1"></i>{{ house.area }}m²</span>
                        </div>
                        <p class="text-surface-900 dark:text-surface-0 font-semibold">
                            K{{ formatPrice(house.monthly_rent) }} <span class="font-normal text-surface-500">/ month</span>
                        </p>
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div v-else class="text-center py-20">
                <i class="pi pi-home text-6xl text-surface-300 mb-4"></i>
                <h2 class="text-2xl font-semibold text-surface-900 dark:text-surface-0 mb-2">No properties found</h2>
                <p class="text-surface-500 mb-6">Try adjusting your search or filters</p>
                <Button label="Clear Filters" @click="selectedType = null; searchQuery = ''" outlined />
            </div>
        </main>

        <!-- Footer -->
        <footer class="border-t border-surface-200 dark:border-surface-700 mt-auto">
            <div class="max-w-screen-2xl mx-auto px-6 py-8">
                <div class="flex flex-col md:flex-row justify-between items-center gap-4">
                    <div class="flex items-center gap-2">
                        <!-- <i class="pi pi-home text-primary"></i> -->
                        <span class="font-semibold text-surface-900 dark:text-surface-0">BondiHub</span>
                    </div>
                    <p class="text-surface-500 text-sm">&copy; {{ new Date().getFullYear() }} BondiHub. All rights reserved.</p>
                </div>
            </div>
        </footer>
    </div>
</template>
