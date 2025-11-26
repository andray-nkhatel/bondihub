<template>
  <div class="dashboard">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">Dashboard</h1>
      <span class="text-surface-500">Welcome back, {{ userName }}!</span>
    </div>

    <!-- Admin Dashboard Stats -->
    <div v-if="isAdmin && adminStats" class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <Panel class="stat-card">
        <template #header><span class="font-semibold text-blue-600"><i class="pi pi-users mr-2"></i>Total Users</span></template>
        <div class="text-center">
          <p class="text-4xl font-bold text-blue-600">{{ adminStats.total_users || 0 }}</p>
        </div>
      </Panel>
      <Panel class="stat-card">
        <template #header><span class="font-semibold text-green-600"><i class="pi pi-home mr-2"></i>Properties</span></template>
        <div class="text-center">
          <p class="text-4xl font-bold text-green-600">{{ adminStats.total_houses || 0 }}</p>
        </div>
      </Panel>
      <Panel class="stat-card">
        <template #header><span class="font-semibold text-purple-600"><i class="pi pi-file mr-2"></i>Active Rentals</span></template>
        <div class="text-center">
          <p class="text-4xl font-bold text-purple-600">{{ adminStats.active_rentals || 0 }}</p>
        </div>
      </Panel>
      <Panel class="stat-card">
        <template #header><span class="font-semibold text-orange-600"><i class="pi pi-dollar mr-2"></i>Revenue</span></template>
        <div class="text-center">
          <p class="text-4xl font-bold text-orange-600">K{{ formatNumber(adminStats.total_revenue) }}</p>
        </div>
      </Panel>
    </div>

    <!-- Quick Actions -->
    <Panel class="mb-6">
      <template #header>
        <span class="font-semibold">Quick Actions</span>
      </template>
      <div class="flex flex-wrap gap-3">
        <Button label="Browse Properties" icon="pi pi-home" @click="$router.push('/app/houses')" />
        <Button v-if="isLandlord" label="Add Property" icon="pi pi-plus" severity="success" @click="$router.push('/app/houses?action=create')" />
        <Button label="My Rentals" icon="pi pi-file" outlined @click="$router.push('/app/rentals')" />
        <Button v-if="isTenant" label="My Favorites" icon="pi pi-heart" outlined severity="danger" @click="$router.push('/app/favorites')" />
        <Button label="Payments" icon="pi pi-credit-card" outlined severity="info" @click="$router.push('/app/payments')" />
        <Button label="Maintenance" icon="pi pi-wrench" outlined severity="warn" @click="$router.push('/app/maintenance')" />
      </div>
    </Panel>

    <!-- Recent Activity / Stats based on role -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent Properties (for tenants) or My Properties (for landlords) -->
      <Panel>
        <template #header>
          <span class="font-semibold">{{ isLandlord ? 'My Properties' : 'Recent Properties' }}</span>
        </template>
        <div v-if="loadingHouses" class="flex justify-center py-4">
          <ProgressSpinner style="width: 30px; height: 30px;" />
        </div>
        <div v-else-if="recentHouses.length > 0">
          <div v-for="house in recentHouses.slice(0, 5)" :key="house.id" 
            class="flex items-center gap-3 p-3 hover:bg-surface-100 dark:hover:bg-surface-700 rounded-lg cursor-pointer transition-colors"
            @click="$router.push(`/app/houses?view=${house.id}`)">
            <img v-if="house.images?.[0]" :src="house.images[0].image_url" class="w-12 h-12 rounded object-cover" />
            <i v-else class="pi pi-home text-2xl text-surface-400 w-12 text-center"></i>
            <div class="flex-1">
              <p class="font-medium">{{ house.title }}</p>
              <p class="text-sm text-surface-500">K{{ formatNumber(house.monthly_rent) }}/month</p>
            </div>
            <Tag :severity="getStatusSeverity(house.status)" :value="house.status" />
          </div>
        </div>
        <div v-else class="text-center py-4 text-surface-500">
          <i class="pi pi-home text-3xl mb-2"></i>
          <p>No properties found</p>
        </div>
      </Panel>

      <!-- Notifications -->
      <Panel>
        <template #header>
          <div class="flex justify-between items-center w-full">
            <span class="font-semibold">Recent Notifications</span>
            <Badge v-if="unreadNotifications > 0" :value="unreadNotifications" severity="danger" />
          </div>
        </template>
        <div v-if="loadingNotifications" class="flex justify-center py-4">
          <ProgressSpinner style="width: 30px; height: 30px;" />
        </div>
        <div v-else-if="notifications.length > 0">
          <div v-for="notif in notifications.slice(0, 5)" :key="notif.id" 
            :class="['p-3 rounded-lg mb-2 cursor-pointer transition-colors', notif.is_read ? 'bg-surface-50 dark:bg-surface-800' : 'bg-blue-50 dark:bg-blue-900/20']"
            @click="$router.push('/app/notifications')">
            <p class="font-medium text-sm">{{ notif.title }}</p>
            <p class="text-xs text-surface-500 mt-1">{{ formatRelativeDate(notif.created_at) }}</p>
          </div>
        </div>
        <div v-else class="text-center py-4 text-surface-500">
          <i class="pi pi-bell text-3xl mb-2"></i>
          <p>No notifications</p>
        </div>
        <Button v-if="notifications.length > 0" label="View All" link class="mt-2" @click="$router.push('/app/notifications')" />
      </Panel>
    </div>

    <!-- Maintenance Requests (for landlords/tenants) -->
    <Panel class="mt-6" v-if="!isAdmin">
      <template #header>
        <span class="font-semibold">Recent Maintenance Requests</span>
      </template>
      <div v-if="loadingMaintenance" class="flex justify-center py-4">
        <ProgressSpinner style="width: 30px; height: 30px;" />
      </div>
      <DataTable v-else-if="maintenanceRequests.length > 0" :value="maintenanceRequests.slice(0, 5)" responsiveLayout="scroll" class="p-datatable-sm">
        <Column field="title" header="Title" />
        <Column field="house.title" header="Property">
          <template #body="{ data }">{{ data.house?.title || 'N/A' }}</template>
        </Column>
        <Column field="priority" header="Priority">
          <template #body="{ data }">
            <Tag :severity="getPrioritySeverity(data.priority)" :value="data.priority" />
          </template>
        </Column>
        <Column field="status" header="Status">
          <template #body="{ data }">
            <Tag :severity="getMaintenanceStatusSeverity(data.status)" :value="data.status" />
          </template>
        </Column>
      </DataTable>
      <div v-else class="text-center py-4 text-surface-500">
        <i class="pi pi-wrench text-3xl mb-2"></i>
        <p>No maintenance requests</p>
      </div>
      <Button v-if="maintenanceRequests.length > 0" label="View All" link class="mt-2" @click="$router.push('/app/maintenance')" />
    </Panel>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import { houseService, notificationService, maintenanceService, adminService } from '@/service/api.service';

const store = useStore();

// State
const recentHouses = ref([]);
const notifications = ref([]);
const maintenanceRequests = ref([]);
const adminStats = ref(null);
const loadingHouses = ref(false);
const loadingNotifications = ref(false);
const loadingMaintenance = ref(false);

// Computed
const user = computed(() => store.getters['auth/user']);
const userName = computed(() => user.value?.full_name || user.value?.fullName || user.value?.email || 'User');
const userRole = computed(() => user.value?.role || user.value?.roles?.[0]);
const isAdmin = computed(() => userRole.value === 'admin');
const isLandlord = computed(() => userRole.value === 'landlord');
const isTenant = computed(() => userRole.value === 'tenant');
const unreadNotifications = computed(() => notifications.value.filter(n => !n.is_read).length);

// Methods
const formatNumber = (num) => new Intl.NumberFormat('en-ZM').format(num || 0);

const formatRelativeDate = (date) => {
  if (!date) return '';
  const d = new Date(date);
  const now = new Date();
  const diff = now - d;
  
  if (diff < 60000) return 'Just now';
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
  return d.toLocaleDateString();
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'available': return 'success';
    case 'occupied': return 'info';
    case 'maintenance': return 'warn';
    default: return 'secondary';
  }
};

const getPrioritySeverity = (priority) => {
  switch (priority) {
    case 'urgent': return 'danger';
    case 'high': return 'warn';
    case 'medium': return 'info';
    default: return 'secondary';
  }
};

const getMaintenanceStatusSeverity = (status) => {
  switch (status) {
    case 'completed': return 'success';
    case 'in_progress': return 'info';
    case 'pending': return 'warn';
    default: return 'secondary';
  }
};

const loadDashboardData = async () => {
  // Load houses
  loadingHouses.value = true;
  try {
    const response = await houseService.getAll({ limit: 5 });
    recentHouses.value = response.houses || [];
  } catch (error) {
    console.error('Failed to load houses:', error);
  } finally {
    loadingHouses.value = false;
  }

  // Load notifications
  loadingNotifications.value = true;
  try {
    const response = await notificationService.getAll({ limit: 5 });
    notifications.value = response.notifications || [];
  } catch (error) {
    console.error('Failed to load notifications:', error);
  } finally {
    loadingNotifications.value = false;
  }

  // Load maintenance
  if (!isAdmin.value) {
    loadingMaintenance.value = true;
    try {
      const response = await maintenanceService.getAll({ limit: 5 });
      maintenanceRequests.value = response.maintenance_requests || response.requests || [];
    } catch (error) {
      console.error('Failed to load maintenance:', error);
    } finally {
      loadingMaintenance.value = false;
    }
  }

  // Load admin stats
  if (isAdmin.value) {
    try {
      adminStats.value = await adminService.getDashboardStats();
    } catch (error) {
      console.error('Failed to load admin stats:', error);
    }
  }
};

onMounted(() => {
  loadDashboardData();
});
</script>

<style scoped>
.stat-card {
  transition: transform 0.2s;
}
.stat-card:hover {
  transform: translateY(-2px);
}
</style>
