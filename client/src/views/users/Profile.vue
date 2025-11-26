<script setup>
import { useToast } from 'primevue/usetoast';
import { computed, onMounted, ref } from 'vue';
import { authService } from '../../service/api.service';
import { formatDate as formatDateUtil, formatDateTime as formatDateTimeUtil } from '@/service/dateUtils';

// Utilities
const toast = useToast();

// Reactive data
const loading = ref(false);
const profileData = ref({
  id: null,
  username: '',
  fullName: '',
  email: '',
  role: '',
  isActive: true,
  createdAt: '',
  lastLoginAt: ''
});

const formattedCreatedAt = computed(() => {
  if (!profileData.value.createdAt) return 'N/A';
  return formatDateUtil(profileData.value.createdAt, {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
});

const formattedLastLogin = computed(() => {
  if (!profileData.value.lastLoginAt) return 'Never';
  return formatDateTimeUtil(profileData.value.lastLoginAt, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
});

const roleIcon = computed(() => {
  switch (profileData.value.role) {
    case 'Admin': return 'pi-crown';
    case 'Teacher': return 'pi-graduation-cap';
    case 'Staff': return 'pi-users';
    default: return 'pi-user';
  }
});

const roleColor = computed(() => {
  switch (profileData.value.role) {
    case 'Admin': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
    case 'Teacher': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
    case 'Staff': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
  }
});

// Methods
const loadProfile = async () => {
  loading.value = true;
  try {
    const data = await authService.getProfile();
    profileData.value = data;
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load profile data',
      life: 3000
    });
  } finally {
    loading.value = false;
  }
};

const refreshProfile = async () => {
  await loadProfile();
  toast.add({
    severity: 'info',
    summary: 'Refreshed',
    detail: 'Profile data refreshed',
    life: 2000
  });
};

const profileTableRows = computed(() => [
  { label: 'Full Name', value: profileData.value.fullName || 'Not provided' },
  { label: 'Username', value: profileData.value.username },
  { label: 'Email Address', value: profileData.value.email || 'Not provided' },
  { label: 'Role', value: profileData.value.role },
  { label: 'Account Created', value: formattedCreatedAt.value },
  { label: 'Last Login', value: formattedLastLogin.value }
]);

// Lifecycle
onMounted(() => {
  loadProfile();
});
</script>

<template>
  <div class="relative">
    <!-- Loading Overlay -->
    <div v-if="loading && !profileData.id" class="absolute inset-0 z-50 flex items-center justify-center bg-black bg-opacity-30">
      <ProgressSpinner style="width: 60px; height: 60px" strokeWidth="4" />
    </div>

    <!-- Profile Header -->
    <div class="flex flex-col items-center justify-center gap-3 py-8 mb-6 bg-slate-50 dark:bg-slate-800 rounded-xl shadow-md w-full px-4 text-gray-900 dark:text-gray-100">
      <Avatar
        :label="profileData.fullName?.charAt(0) || 'U'"
        size="xxlarge"
        class="text-white font-bold shadow-lg border-4 border-white mb-2"
        :style="{ backgroundColor: profileData.role === 'Admin' ? '#ef4444' : profileData.role === 'Teacher' ? '#3b82f6' : '#10b981' }"
      />
      <div class="flex flex-col sm:flex-row items-center gap-2 mt-2 w-full justify-center text-center">
        <h2 class="text-2xl sm:text-3xl font-bold text-900 dark:text-gray-100 m-0">{{ profileData.fullName }}</h2>
        <Tag :class="roleColor" :icon="`pi ${roleIcon}`" class="ml-0 sm:ml-2 mt-1 sm:mt-0">{{ profileData.role }}</Tag>
        <Tag :severity="profileData.isActive ? 'success' : 'danger'" :value="profileData.isActive ? 'Active' : 'Inactive'" class="ml-0 sm:ml-2 mt-1 sm:mt-0" />
      </div>
      <p class="text-600 dark:text-gray-300 mt-1 mb-0 text-sm sm:text-base">View your account information</p>
      <Button 
        icon="pi pi-refresh" 
        outlined 
        @click="refreshProfile"
        :loading="loading"
        v-tooltip.top="'Refresh Profile'"
        class="mt-2 w-full sm:w-auto"
      />
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Profile Details Card -->
      <Panel class="shadow-md bg-slate-50 dark:bg-slate-800 text-gray-900 dark:text-gray-100">
        <template #header>
          <div class="flex items-center gap-2">
            <i class="pi pi-user text-primary"></i>
            <span class="font-semibold">Personal Details</span>
          </div>
        </template>
        <div class="p-2">
          <DataTable :value="profileTableRows" class="p-datatable-sm w-full text-gray-900 dark:text-gray-100">
            <Column field="label" header="Field" style="width: 40%">
              <template #body="{ data }">
                <span class="flex items-center gap-2">
                  <i v-if="data.label === 'Full Name'" class="pi pi-id-card text-blue-400"></i>
                  <i v-else-if="data.label === 'Username'" class="pi pi-user text-green-400"></i>
                  <i v-else-if="data.label === 'Email Address'" class="pi pi-envelope text-purple-400"></i>
                  <i v-else-if="data.label === 'Role'" class="pi pi-shield text-orange-400"></i>
                  <i v-else-if="data.label === 'Account Created'" class="pi pi-calendar-plus text-cyan-400"></i>
                  <i v-else-if="data.label === 'Last Login'" class="pi pi-clock text-pink-400"></i>
                  {{ data.label }}
                </span>
              </template>
            </Column>
            <Column field="value" header="Value"></Column>
          </DataTable>
        </div>
      </Panel>

      <!-- Quick Stats Card -->
      <Panel class="shadow-md bg-slate-50 dark:bg-slate-800 text-gray-900 dark:text-gray-100">
        <template #header>
          <div class="flex items-center gap-2">
            <i class="pi pi-info-circle text-primary"></i>
            <span class="font-semibold">Account Information</span>
          </div>
        </template>
        <div class="flex flex-col gap-4 p-2">
          <div class="flex items-center justify-between p-3 rounded-lg bg-slate-50 hover:bg-slate-100 dark:bg-slate-800 dark:hover:bg-slate-700 transition">
            <div class="flex items-center gap-2">
              <i class="pi pi-id-card text-primary"></i>
              <span class="font-medium">User ID</span>
            </div>
            <span class="font-bold">#{{ profileData.id }}</span>
          </div>
          <div class="flex items-center justify-between p-3 rounded-lg bg-slate-50 hover:bg-slate-100 dark:bg-slate-800 dark:hover:bg-slate-700 transition">
            <div class="flex items-center gap-2">
              <i class="pi pi-shield text-primary"></i>
              <span class="font-medium">Account Status</span>
            </div>
            <Tag :severity="profileData.isActive ? 'success' : 'danger'" :value="profileData.isActive ? 'Active' : 'Inactive'" />
          </div>
          <div class="flex items-center justify-between p-3 rounded-lg bg-slate-50 hover:bg-slate-100 dark:bg-slate-800 dark:hover:bg-slate-700 transition">
            <div class="flex items-center gap-2">
              <i :class="`pi ${roleIcon} text-primary`"></i>
              <span class="font-medium">Access Level</span>
            </div>
            <span class="font-bold">{{ profileData.role }}</span>
          </div>
        </div>
      </Panel>
    </div>

    <!-- Session Information Card -->
    <Panel class="mt-6 shadow-md bg-slate-50 dark:bg-slate-800 text-gray-900 dark:text-gray-100">
      <template #header>
        <div class="flex items-center gap-2">
          <i class="pi pi-clock text-primary"></i>
          <span class="font-semibold">Session Information</span>
        </div>
      </template>
      <div class="flex flex-col md:flex-row gap-6 p-4 items-center justify-center">
        <div class="flex-1 text-center">
          <i class="pi pi-clock text-primary text-3xl mb-2"></i>
          <p class="text-600 dark:text-gray-300 text-sm m-0">Last Login</p>
          <p class="font-bold text-900 dark:text-gray-100 mt-1">{{ formattedLastLogin }}</p>
        </div>
        <Divider layout="vertical" class="hidden md:block" />
        <div class="flex-1 text-center">
          <i class="pi pi-calendar text-primary text-3xl mb-2"></i>
          <p class="text-600 dark:text-gray-300 text-sm m-0">Member Since</p>
          <p class="font-bold text-900 dark:text-gray-100 mt-1">{{ formattedCreatedAt }}</p>
        </div>
      </div>
    </Panel>
  </div>
</template>

<style scoped>
.field {
  margin-bottom: 1rem;
}

.field label {
  display: block;
  margin-bottom: 0.25rem;
}
</style>