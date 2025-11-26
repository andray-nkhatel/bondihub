<template>
  <div class="maintenance-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">Maintenance Requests</h1>
      <Button v-if="canCreateRequest" label="New Request" icon="pi pi-plus" @click="openCreateDialog" />
    </div>

    <!-- Stats -->
    <div v-if="stats" class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <Panel>
        <template #header><span class="font-semibold">Total Requests</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-primary">{{ stats.total || 0 }}</p>
        </div>
      </Panel>
      <Panel>
        <template #header><span class="font-semibold">Pending</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-yellow-500">{{ stats.pending || 0 }}</p>
        </div>
      </Panel>
      <Panel>
        <template #header><span class="font-semibold">In Progress</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-blue-500">{{ stats.in_progress || 0 }}</p>
        </div>
      </Panel>
      <Panel>
        <template #header><span class="font-semibold">Completed</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-green-500">{{ stats.completed || 0 }}</p>
        </div>
      </Panel>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <!-- Requests Table -->
    <DataTable v-else :value="requests" :paginator="true" :rows="10" dataKey="id" responsiveLayout="scroll">
      <Column field="title" header="Title" sortable />
      
      <Column field="house.title" header="Property">
        <template #body="{ data }">{{ data.house?.title || 'N/A' }}</template>
      </Column>
      
      <Column field="priority" header="Priority" sortable>
        <template #body="{ data }">
          <Tag :severity="getPrioritySeverity(data.priority)" :value="data.priority" />
        </template>
      </Column>
      
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :severity="getStatusSeverity(data.status)" :value="data.status" />
        </template>
      </Column>
      
      <Column field="created_at" header="Created" sortable>
        <template #body="{ data }">{{ formatDate(data.created_at) }}</template>
      </Column>
      
      <Column header="Actions" style="width: 150px">
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-eye" outlined size="small" @click="viewRequest(data)" />
            <Button v-if="canUpdateRequest" icon="pi pi-pencil" outlined size="small" severity="secondary" @click="editRequest(data)" />
          </div>
        </template>
      </Column>
      
      <template #empty>
        <div class="text-center py-8">
          <i class="pi pi-wrench text-4xl text-surface-300 mb-4"></i>
          <p class="text-surface-500">No maintenance requests found</p>
        </div>
      </template>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="requestDialog" :header="dialogMode === 'view' ? 'Request Details' : (dialogMode === 'edit' ? 'Update Request' : 'New Request')" 
      :modal="true" :style="{ width: '600px' }">
      
      <div v-if="dialogMode === 'view' && selectedRequest" class="request-detail">
        <div class="grid grid-cols-2 gap-4">
          <div class="col-span-2"><strong>Title:</strong> {{ selectedRequest.title }}</div>
          <div><strong>Property:</strong> {{ selectedRequest.house?.title }}</div>
          <div><strong>Priority:</strong> <Tag :severity="getPrioritySeverity(selectedRequest.priority)" :value="selectedRequest.priority" /></div>
          <div><strong>Status:</strong> <Tag :severity="getStatusSeverity(selectedRequest.status)" :value="selectedRequest.status" /></div>
          <div><strong>Created:</strong> {{ formatDate(selectedRequest.created_at) }}</div>
          <div class="col-span-2"><strong>Description:</strong><br>{{ selectedRequest.description }}</div>
          <div v-if="selectedRequest.resolution" class="col-span-2"><strong>Resolution:</strong><br>{{ selectedRequest.resolution }}</div>
        </div>
      </div>

      <div v-else class="request-form">
        <div class="grid gap-4">
          <div v-if="dialogMode === 'create'" class="field">
            <label class="block mb-2 font-medium">Property *</label>
            <Select v-model="requestForm.house_id" :options="userHouses" optionLabel="title" optionValue="id" 
              placeholder="Select property" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Title *</label>
            <InputText v-model="requestForm.title" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Priority</label>
            <Select v-model="requestForm.priority" :options="priorityOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>
          
          <div v-if="dialogMode === 'edit' && canUpdateRequest" class="field">
            <label class="block mb-2 font-medium">Status</label>
            <Select v-model="requestForm.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Description *</label>
            <Textarea v-model="requestForm.description" rows="4" class="w-full" />
          </div>
          
          <div v-if="dialogMode === 'edit' && canUpdateRequest" class="field">
            <label class="block mb-2 font-medium">Resolution Notes</label>
            <Textarea v-model="requestForm.resolution" rows="3" class="w-full" />
          </div>
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" @click="requestDialog = false" outlined />
        <Button v-if="dialogMode !== 'view'" :label="dialogMode === 'edit' ? 'Update' : 'Submit'" icon="pi pi-check" @click="saveRequest" :loading="saving" />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import { useToast } from 'primevue/usetoast';
import { maintenanceService, houseService } from '@/service/api.service';

const store = useStore();
const toast = useToast();

const requests = ref([]);
const stats = ref(null);
const userHouses = ref([]);
const loading = ref(false);
const saving = ref(false);
const requestDialog = ref(false);
const dialogMode = ref('view');
const selectedRequest = ref(null);

const requestForm = ref({
  house_id: null,
  title: '',
  description: '',
  priority: 'medium',
  status: 'pending',
  resolution: ''
});

const priorityOptions = [
  { label: 'Low', value: 'low' },
  { label: 'Medium', value: 'medium' },
  { label: 'High', value: 'high' },
  { label: 'Urgent', value: 'urgent' }
];

const statusOptions = [
  { label: 'Pending', value: 'pending' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' }
];

const user = computed(() => store.getters['auth/user']);
const userRole = computed(() => user.value?.role || user.value?.roles?.[0]);
const canCreateRequest = computed(() => ['tenant', 'admin'].includes(userRole.value));
const canUpdateRequest = computed(() => ['landlord', 'admin'].includes(userRole.value));

const loadRequests = async () => {
  loading.value = true;
  try {
    const response = await maintenanceService.getAll();
    requests.value = response.maintenance_requests || response.requests || [];
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to load requests', life: 3000 });
  } finally {
    loading.value = false;
  }
};

const loadStats = async () => {
  try {
    stats.value = await maintenanceService.getStats();
  } catch (error) {
    console.error('Failed to load stats:', error);
  }
};

const loadUserHouses = async () => {
  try {
    const response = await houseService.getAll();
    userHouses.value = response.houses || [];
  } catch (error) {
    console.error('Failed to load houses:', error);
  }
};

const getPrioritySeverity = (priority) => {
  switch (priority) {
    case 'urgent': return 'danger';
    case 'high': return 'warn';
    case 'medium': return 'info';
    case 'low': return 'secondary';
    default: return 'info';
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'completed': return 'success';
    case 'in_progress': return 'info';
    case 'pending': return 'warn';
    case 'cancelled': return 'secondary';
    default: return 'info';
  }
};

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A';

const openCreateDialog = () => {
  dialogMode.value = 'create';
  requestForm.value = {
    house_id: null,
    title: '',
    description: '',
    priority: 'medium',
    status: 'pending',
    resolution: ''
  };
  loadUserHouses();
  requestDialog.value = true;
};

const viewRequest = (request) => {
  dialogMode.value = 'view';
  selectedRequest.value = request;
  requestDialog.value = true;
};

const editRequest = (request) => {
  dialogMode.value = 'edit';
  selectedRequest.value = request;
  requestForm.value = { ...request };
  requestDialog.value = true;
};

const saveRequest = async () => {
  saving.value = true;
  try {
    if (dialogMode.value === 'create') {
      await maintenanceService.create(requestForm.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Request submitted', life: 3000 });
    } else {
      await maintenanceService.update(selectedRequest.value.id, requestForm.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Request updated', life: 3000 });
    }
    requestDialog.value = false;
    loadRequests();
    loadStats();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to save request', life: 3000 });
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  loadRequests();
  loadStats();
});
</script>

