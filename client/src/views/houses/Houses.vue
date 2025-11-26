<template>
  <div class="houses-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">Properties</h1>
      <Button v-if="canManageHouses" label="Add Property" icon="pi pi-plus" @click="openCreateDialog" />
    </div>

    <!-- Filters -->
    <Panel class="mb-6">
      <template #header>
        <span class="font-semibold">Search & Filter</span>
      </template>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <InputText v-model="filters.search" placeholder="Search properties..." @input="debouncedSearch" />
        <Select v-model="filters.status" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="Status" showClear @change="loadHouses" />
        <Select v-model="filters.houseType" :options="houseTypeOptions" optionLabel="label" optionValue="value" placeholder="Type" showClear @change="loadHouses" />
        <InputNumber v-model="filters.maxRent" placeholder="Max Rent" mode="currency" currency="ZMW" @input="debouncedSearch" />
      </div>
    </Panel>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <!-- Houses Grid -->
    <div v-else-if="houses.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card v-for="house in houses" :key="house.id" class="house-card">
        <template #header>
          <div class="house-image-container">
            <img v-if="house.images?.length" :src="house.images[0].image_url" :alt="house.title" class="house-image" />
            <div v-else class="house-image-placeholder">
              <i class="pi pi-home text-4xl text-surface-400"></i>
            </div>
            <Tag :severity="getStatusSeverity(house.status)" :value="house.status" class="status-tag" />
            <Button v-if="canFavorite" :icon="isFavorite(house.id) ? 'pi pi-heart-fill' : 'pi pi-heart'" 
              class="favorite-btn" rounded text 
              :severity="isFavorite(house.id) ? 'danger' : 'secondary'"
              @click.stop="toggleFavorite(house.id)" />
          </div>
        </template>
        <template #title>
          <div class="flex justify-between items-start">
            <span class="text-lg font-semibold truncate">{{ house.title }}</span>
          </div>
        </template>
        <template #subtitle>
          <div class="flex items-center gap-2 text-surface-500">
            <i class="pi pi-map-marker"></i>
            <span class="truncate">{{ house.address }}</span>
          </div>
        </template>
        <template #content>
          <div class="house-details">
            <div class="flex items-center gap-4 text-sm text-surface-600 dark:text-surface-300 mb-3">
              <span><i class="pi pi-th-large mr-1"></i>{{ house.bedrooms }} bed</span>
              <span><i class="pi pi-inbox mr-1"></i>{{ house.bathrooms }} bath</span>
              <span><i class="pi pi-arrows-alt mr-1"></i>{{ house.area }} m²</span>
            </div>
            <div class="text-2xl font-bold text-primary">
              K{{ formatNumber(house.monthly_rent) }}<span class="text-sm font-normal text-surface-500">/month</span>
            </div>
          </div>
        </template>
        <template #footer>
          <div class="flex gap-2">
            <Button label="View" icon="pi pi-eye" outlined size="small" @click="viewHouse(house)" class="flex-1" />
            <Button v-if="canEdit(house)" icon="pi pi-pencil" outlined size="small" severity="secondary" @click="editHouse(house)" />
            <Button v-if="canEdit(house)" icon="pi pi-trash" outlined size="small" severity="danger" @click="confirmDelete(house)" />
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <i class="pi pi-home text-6xl text-surface-300 mb-4"></i>
      <p class="text-xl text-surface-500">No properties found</p>
      <Button v-if="canManageHouses" label="Add Your First Property" icon="pi pi-plus" class="mt-4" @click="openCreateDialog" />
    </div>

    <!-- Pagination -->
    <Paginator v-if="totalRecords > pageSize" :rows="pageSize" :totalRecords="totalRecords" :first="(currentPage - 1) * pageSize" @page="onPageChange" class="mt-6" />

    <!-- View/Edit Dialog -->
    <Dialog v-model:visible="houseDialog" :header="dialogMode === 'view' ? 'Property Details' : (dialogMode === 'edit' ? 'Edit Property' : 'Add Property')" 
      :modal="true" :style="{ width: '800px' }" :closable="true">
      
      <div v-if="dialogMode === 'view' && selectedHouse" class="house-detail-view">
        <!-- Image Gallery -->
        <div v-if="selectedHouse.images?.length" class="mb-6">
          <Galleria :value="selectedHouse.images" :numVisible="5" containerStyle="max-width: 100%">
            <template #item="slotProps">
              <img :src="slotProps.item.image_url" :alt="selectedHouse.title" style="width: 100%; max-height: 400px; object-fit: cover;" />
            </template>
            <template #thumbnail="slotProps">
              <img :src="slotProps.item.image_url" :alt="selectedHouse.title" style="width: 80px; height: 60px; object-fit: cover;" />
            </template>
          </Galleria>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div><strong>Title:</strong> {{ selectedHouse.title }}</div>
          <div><strong>Type:</strong> {{ selectedHouse.house_type }}</div>
          <div><strong>Status:</strong> <Tag :severity="getStatusSeverity(selectedHouse.status)" :value="selectedHouse.status" /></div>
          <div><strong>Monthly Rent:</strong> K{{ formatNumber(selectedHouse.monthly_rent) }}</div>
          <div><strong>Bedrooms:</strong> {{ selectedHouse.bedrooms }}</div>
          <div><strong>Bathrooms:</strong> {{ selectedHouse.bathrooms }}</div>
          <div><strong>Area:</strong> {{ selectedHouse.area }} m²</div>
          <div><strong>Address:</strong> {{ selectedHouse.address }}</div>
          <div class="col-span-2"><strong>Description:</strong><br>{{ selectedHouse.description }}</div>
        </div>

        <!-- Reviews Section -->
        <Divider />
        <h3 class="text-lg font-semibold mb-4">Reviews</h3>
        <div v-if="houseReviews.length > 0">
          <div v-for="review in houseReviews" :key="review.id" class="mb-4 p-4 bg-surface-100 dark:bg-surface-800 rounded-lg">
            <div class="flex justify-between items-center mb-2">
              <Rating :modelValue="review.rating" :readonly="true" :cancel="false" />
              <span class="text-sm text-surface-500">{{ formatDate(review.created_at) }}</span>
            </div>
            <p>{{ review.comment }}</p>
          </div>
        </div>
        <p v-else class="text-surface-500">No reviews yet</p>
      </div>

      <!-- Create/Edit Form -->
      <div v-else class="house-form">
        <div class="grid grid-cols-2 gap-4">
          <div class="field col-span-2">
            <label for="title" class="block mb-2 font-medium">Title *</label>
            <InputText id="title" v-model="houseForm.title" class="w-full" :class="{'p-invalid': errors.title}" />
            <small v-if="errors.title" class="p-error">{{ errors.title }}</small>
          </div>
          
          <div class="field">
            <label for="houseType" class="block mb-2 font-medium">Type *</label>
            <Select id="houseType" v-model="houseForm.house_type" :options="houseTypeOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>
          
          <div class="field">
            <label for="monthlyRent" class="block mb-2 font-medium">Monthly Rent (ZMW) *</label>
            <InputNumber id="monthlyRent" v-model="houseForm.monthly_rent" mode="currency" currency="ZMW" class="w-full" />
          </div>
          
          <div class="field">
            <label for="bedrooms" class="block mb-2 font-medium">Bedrooms</label>
            <InputNumber id="bedrooms" v-model="houseForm.bedrooms" :min="0" class="w-full" />
          </div>
          
          <div class="field">
            <label for="bathrooms" class="block mb-2 font-medium">Bathrooms</label>
            <InputNumber id="bathrooms" v-model="houseForm.bathrooms" :min="0" class="w-full" />
          </div>
          
          <div class="field">
            <label for="area" class="block mb-2 font-medium">Area (m²)</label>
            <InputNumber id="area" v-model="houseForm.area" :min="0" class="w-full" />
          </div>
          
          <div class="field">
            <label for="status" class="block mb-2 font-medium">Status</label>
            <Select id="status" v-model="houseForm.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>
          
          <div class="field col-span-2">
            <label for="address" class="block mb-2 font-medium">Address *</label>
            <InputText id="address" v-model="houseForm.address" class="w-full" />
          </div>
          
          <div class="field col-span-2">
            <label for="description" class="block mb-2 font-medium">Description</label>
            <Textarea id="description" v-model="houseForm.description" rows="4" class="w-full" />
          </div>

          <!-- Image Upload -->
          <div v-if="dialogMode === 'edit'" class="field col-span-2">
            <label class="block mb-2 font-medium">Images</label>
            <FileUpload mode="basic" accept="image/*" :maxFileSize="5000000" @select="onImageSelect" chooseLabel="Add Image" />
            <div v-if="selectedHouse?.images?.length" class="flex gap-2 mt-4 flex-wrap">
              <div v-for="img in selectedHouse.images" :key="img.id" class="relative">
                <img :src="img.image_url" class="w-24 h-24 object-cover rounded" />
                <Button icon="pi pi-times" rounded text severity="danger" size="small" 
                  class="absolute -top-2 -right-2" @click="deleteImage(img.id)" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <Button v-if="dialogMode === 'view'" label="Close" icon="pi pi-times" @click="houseDialog = false" outlined />
        <template v-else>
          <Button label="Cancel" icon="pi pi-times" @click="houseDialog = false" outlined />
          <Button :label="dialogMode === 'edit' ? 'Update' : 'Create'" icon="pi pi-check" @click="saveHouse" :loading="saving" />
        </template>
      </template>
    </Dialog>

    <!-- Delete Confirmation -->
    <ConfirmDialog />
    <Toast />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import { useConfirm } from 'primevue/useconfirm';
import { useToast } from 'primevue/usetoast';
import { houseService, reviewService, favoriteService } from '@/service/api.service';
import { debounce } from 'lodash-es';

const store = useStore();
const confirm = useConfirm();
const toast = useToast();

// State
const houses = ref([]);
const loading = ref(false);
const saving = ref(false);
const totalRecords = ref(0);
const currentPage = ref(1);
const pageSize = ref(12);
const favorites = ref([]);
const houseReviews = ref([]);

// Dialog state
const houseDialog = ref(false);
const dialogMode = ref('view'); // 'view', 'create', 'edit'
const selectedHouse = ref(null);

// Form
const houseForm = ref({
  title: '',
  house_type: 'apartment',
  monthly_rent: 0,
  bedrooms: 1,
  bathrooms: 1,
  area: 0,
  address: '',
  description: '',
  status: 'available'
});
const errors = ref({});

// Filters
const filters = ref({
  search: '',
  status: null,
  houseType: null,
  maxRent: null
});

// Options
const statusOptions = [
  { label: 'Available', value: 'available' },
  { label: 'Occupied', value: 'occupied' },
  { label: 'Maintenance', value: 'maintenance' }
];

const houseTypeOptions = [
  { label: 'Apartment', value: 'apartment' },
  { label: 'House', value: 'house' },
  { label: 'Studio', value: 'studio' },
  { label: 'Townhouse', value: 'townhouse' },
  { label: 'Commercial', value: 'commercial' }
];

// Computed
const user = computed(() => store.getters['auth/user']);
const userRole = computed(() => user.value?.role || user.value?.roles?.[0]);
const canManageHouses = computed(() => ['landlord', 'admin'].includes(userRole.value));
const canFavorite = computed(() => ['tenant', 'admin'].includes(userRole.value));

const canEdit = (house) => {
  if (userRole.value === 'admin') return true;
  if (userRole.value === 'landlord' && house.landlord_id === user.value?.id) return true;
  return false;
};

const isFavorite = (houseId) => favorites.value.includes(houseId);

// Methods
const loadHouses = async () => {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      search: filters.value.search || undefined,
      status: filters.value.status || undefined,
      house_type: filters.value.houseType || undefined,
      max_rent: filters.value.maxRent || undefined
    };
    
    const response = await houseService.getAll(params);
    houses.value = response.houses || [];
    totalRecords.value = response.total || houses.value.length;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to load properties', life: 3000 });
  } finally {
    loading.value = false;
  }
};

const loadFavorites = async () => {
  if (!canFavorite.value) return;
  try {
    const response = await favoriteService.getAll();
    favorites.value = (response.favorites || []).map(f => f.house_id);
  } catch (error) {
    console.error('Failed to load favorites:', error);
  }
};

const debouncedSearch = debounce(() => {
  currentPage.value = 1;
  loadHouses();
}, 300);

const onPageChange = (event) => {
  currentPage.value = event.page + 1;
  loadHouses();
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'available': return 'success';
    case 'occupied': return 'info';
    case 'maintenance': return 'warn';
    default: return 'secondary';
  }
};

const formatNumber = (num) => {
  return new Intl.NumberFormat('en-ZM').format(num);
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};

const openCreateDialog = () => {
  dialogMode.value = 'create';
  houseForm.value = {
    title: '',
    house_type: 'apartment',
    monthly_rent: 0,
    bedrooms: 1,
    bathrooms: 1,
    area: 0,
    address: '',
    description: '',
    status: 'available'
  };
  errors.value = {};
  houseDialog.value = true;
};

const viewHouse = async (house) => {
  dialogMode.value = 'view';
  selectedHouse.value = house;
  houseDialog.value = true;
  
  // Load reviews
  try {
    const response = await reviewService.getByHouse(house.id);
    houseReviews.value = response.reviews || [];
  } catch (error) {
    houseReviews.value = [];
  }
};

const editHouse = (house) => {
  dialogMode.value = 'edit';
  selectedHouse.value = house;
  houseForm.value = { ...house };
  errors.value = {};
  houseDialog.value = true;
};

const saveHouse = async () => {
  errors.value = {};
  
  if (!houseForm.value.title) {
    errors.value.title = 'Title is required';
    return;
  }
  
  saving.value = true;
  try {
    if (dialogMode.value === 'create') {
      await houseService.create(houseForm.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Property created successfully', life: 3000 });
    } else {
      await houseService.update(selectedHouse.value.id, houseForm.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Property updated successfully', life: 3000 });
    }
    houseDialog.value = false;
    loadHouses();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to save property', life: 3000 });
  } finally {
    saving.value = false;
  }
};

const confirmDelete = (house) => {
  confirm.require({
    message: `Are you sure you want to delete "${house.title}"?`,
    header: 'Delete Confirmation',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await houseService.delete(house.id);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Property deleted', life: 3000 });
        loadHouses();
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to delete property', life: 3000 });
      }
    }
  });
};

const toggleFavorite = async (houseId) => {
  try {
    if (isFavorite(houseId)) {
      await favoriteService.remove(houseId);
      favorites.value = favorites.value.filter(id => id !== houseId);
      toast.add({ severity: 'info', summary: 'Removed', detail: 'Removed from favorites', life: 2000 });
    } else {
      await favoriteService.add(houseId);
      favorites.value.push(houseId);
      toast.add({ severity: 'success', summary: 'Added', detail: 'Added to favorites', life: 2000 });
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to update favorites', life: 3000 });
  }
};

const onImageSelect = async (event) => {
  if (!selectedHouse.value) return;
  
  const file = event.files[0];
  try {
    await houseService.uploadImage(selectedHouse.value.id, file);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Image uploaded', life: 3000 });
    // Reload house to get updated images
    const updated = await houseService.getById(selectedHouse.value.id);
    selectedHouse.value = updated.house;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to upload image', life: 3000 });
  }
};

const deleteImage = async (imageId) => {
  try {
    await houseService.deleteImage(imageId);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Image deleted', life: 3000 });
    selectedHouse.value.images = selectedHouse.value.images.filter(img => img.id !== imageId);
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to delete image', life: 3000 });
  }
};

onMounted(() => {
  loadHouses();
  loadFavorites();
});
</script>

<style scoped>
.house-card {
  transition: transform 0.2s, box-shadow 0.2s;
}
.house-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}
.house-image-container {
  position: relative;
  height: 200px;
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
.favorite-btn {
  position: absolute;
  top: 8px;
  right: 8px;
}
</style>

