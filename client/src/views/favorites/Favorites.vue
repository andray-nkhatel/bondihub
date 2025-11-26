<template>
  <div class="favorites-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">
        <i class="pi pi-heart-fill text-red-500 mr-2"></i>My Favorites
      </h1>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <!-- Favorites Grid -->
    <div v-else-if="favorites.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card v-for="fav in favorites" :key="fav.id" class="favorite-card">
        <template #header>
          <div class="house-image-container">
            <img v-if="fav.house?.images?.[0]" :src="fav.house.images[0].image_url" :alt="fav.house?.title" class="house-image" />
            <div v-else class="house-image-placeholder">
              <i class="pi pi-home text-4xl text-surface-400"></i>
            </div>
            <Tag :severity="getStatusSeverity(fav.house?.status)" :value="fav.house?.status" class="status-tag" />
          </div>
        </template>
        <template #title>{{ fav.house?.title || 'Property' }}</template>
        <template #subtitle>
          <div class="flex items-center gap-2 text-surface-500">
            <i class="pi pi-map-marker"></i>
            <span class="truncate">{{ fav.house?.address }}</span>
          </div>
        </template>
        <template #content>
          <div class="flex items-center gap-4 text-sm text-surface-600 dark:text-surface-300 mb-3">
            <span><i class="pi pi-th-large mr-1"></i>{{ fav.house?.bedrooms }} bed</span>
            <span><i class="pi pi-inbox mr-1"></i>{{ fav.house?.bathrooms }} bath</span>
          </div>
          <div class="text-2xl font-bold text-primary">
            K{{ formatNumber(fav.house?.monthly_rent) }}<span class="text-sm font-normal text-surface-500">/month</span>
          </div>
        </template>
        <template #footer>
          <div class="flex gap-2">
            <Button label="View" icon="pi pi-eye" outlined size="small" @click="viewHouse(fav.house)" class="flex-1" />
            <Button icon="pi pi-heart-fill" severity="danger" size="small" @click="removeFavorite(fav.house?.id)" v-tooltip="'Remove'" />
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <i class="pi pi-heart text-6xl text-surface-300 mb-4"></i>
      <p class="text-xl text-surface-500 mb-4">No favorites yet</p>
      <p class="text-surface-400 mb-4">Browse properties and add them to your favorites</p>
      <Button label="Browse Properties" icon="pi pi-home" @click="$router.push('/app/houses')" />
    </div>

    <Toast />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import { favoriteService } from '@/service/api.service';

const router = useRouter();
const toast = useToast();

const favorites = ref([]);
const loading = ref(false);

const loadFavorites = async () => {
  loading.value = true;
  try {
    const response = await favoriteService.getAll();
    favorites.value = response.favorites || [];
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to load favorites', life: 3000 });
  } finally {
    loading.value = false;
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'available': return 'success';
    case 'occupied': return 'info';
    case 'maintenance': return 'warn';
    default: return 'secondary';
  }
};

const formatNumber = (num) => new Intl.NumberFormat('en-ZM').format(num || 0);

const viewHouse = (house) => {
  if (house?.id) {
    router.push(`/app/houses?view=${house.id}`);
  }
};

const removeFavorite = async (houseId) => {
  if (!houseId) return;
  try {
    await favoriteService.remove(houseId);
    favorites.value = favorites.value.filter(f => f.house?.id !== houseId);
    toast.add({ severity: 'info', summary: 'Removed', detail: 'Removed from favorites', life: 2000 });
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to remove', life: 3000 });
  }
};

onMounted(() => {
  loadFavorites();
});
</script>

<style scoped>
.favorite-card {
  transition: transform 0.2s, box-shadow 0.2s;
}
.favorite-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}
.house-image-container {
  position: relative;
  height: 180px;
  overflow: hidden;
}
.house-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.house-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--p-surface-100);
}
.status-tag {
  position: absolute;
  top: 12px;
  left: 12px;
}
</style>

