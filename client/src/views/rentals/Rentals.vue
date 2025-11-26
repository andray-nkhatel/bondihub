<template>
  <div class="rentals-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">Rental Agreements</h1>
      <Button v-if="canCreateRental" label="New Agreement" icon="pi pi-plus" @click="openCreateDialog" />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <!-- Rentals Table -->
    <DataTable v-else :value="rentals" :paginator="true" :rows="10" dataKey="id"
      :rowHover="true" responsiveLayout="scroll" class="p-datatable-sm">
      
      <Column field="house.title" header="Property" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <img v-if="data.house?.images?.[0]" :src="data.house.images[0].image_url" class="w-10 h-10 rounded object-cover" />
            <i v-else class="pi pi-home text-2xl text-surface-400"></i>
            <span>{{ data.house?.title || 'N/A' }}</span>
          </div>
        </template>
      </Column>
      
      <Column field="tenant.full_name" header="Tenant" sortable />
      
      <Column field="start_date" header="Start Date" sortable>
        <template #body="{ data }">{{ formatDate(data.start_date) }}</template>
      </Column>
      
      <Column field="end_date" header="End Date" sortable>
        <template #body="{ data }">{{ formatDate(data.end_date) }}</template>
      </Column>
      
      <Column field="monthly_rent" header="Monthly Rent" sortable>
        <template #body="{ data }">K{{ formatNumber(data.monthly_rent) }}</template>
      </Column>
      
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :severity="getStatusSeverity(data.status)" :value="data.status" />
        </template>
      </Column>
      
      <Column header="Actions" style="width: 150px">
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-eye" outlined size="small" @click="viewRental(data)" />
            <Button v-if="data.status === 'active'" icon="pi pi-times" outlined size="small" severity="danger" 
              @click="confirmTerminate(data)" v-tooltip="'Terminate'" />
          </div>
        </template>
      </Column>
      
      <template #empty>
        <div class="text-center py-8">
          <i class="pi pi-file text-4xl text-surface-300 mb-4"></i>
          <p class="text-surface-500">No rental agreements found</p>
        </div>
      </template>
    </DataTable>

    <!-- Create/View Dialog -->
    <Dialog v-model:visible="rentalDialog" :header="dialogMode === 'view' ? 'Rental Details' : 'New Rental Agreement'" 
      :modal="true" :style="{ width: '600px' }">
      
      <div v-if="dialogMode === 'view' && selectedRental" class="rental-detail">
        <div class="grid grid-cols-2 gap-4">
          <div><strong>Property:</strong> {{ selectedRental.house?.title }}</div>
          <div><strong>Tenant:</strong> {{ selectedRental.tenant?.full_name }}</div>
          <div><strong>Start Date:</strong> {{ formatDate(selectedRental.start_date) }}</div>
          <div><strong>End Date:</strong> {{ formatDate(selectedRental.end_date) }}</div>
          <div><strong>Monthly Rent:</strong> K{{ formatNumber(selectedRental.monthly_rent) }}</div>
          <div><strong>Deposit:</strong> K{{ formatNumber(selectedRental.deposit_amount) }}</div>
          <div><strong>Status:</strong> <Tag :severity="getStatusSeverity(selectedRental.status)" :value="selectedRental.status" /></div>
          <div class="col-span-2"><strong>Terms:</strong><br>{{ selectedRental.terms || 'No additional terms' }}</div>
        </div>
      </div>

      <div v-else class="rental-form">
        <div class="grid grid-cols-2 gap-4">
          <div class="field col-span-2">
            <label class="block mb-2 font-medium">Property *</label>
            <Select v-model="rentalForm.house_id" :options="availableHouses" optionLabel="title" optionValue="id" 
              placeholder="Select property" class="w-full" filter />
          </div>
          
          <div class="field col-span-2">
            <label class="block mb-2 font-medium">Tenant Email *</label>
            <InputText v-model="rentalForm.tenant_email" class="w-full" placeholder="tenant@email.com" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Start Date *</label>
            <DatePicker v-model="rentalForm.start_date" class="w-full" dateFormat="yy-mm-dd" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">End Date *</label>
            <DatePicker v-model="rentalForm.end_date" class="w-full" dateFormat="yy-mm-dd" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Monthly Rent (ZMW) *</label>
            <InputNumber v-model="rentalForm.monthly_rent" mode="currency" currency="ZMW" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Deposit (ZMW)</label>
            <InputNumber v-model="rentalForm.deposit_amount" mode="currency" currency="ZMW" class="w-full" />
          </div>
          
          <div class="field col-span-2">
            <label class="block mb-2 font-medium">Terms & Conditions</label>
            <Textarea v-model="rentalForm.terms" rows="4" class="w-full" />
          </div>
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" @click="rentalDialog = false" outlined />
        <Button v-if="dialogMode === 'create'" label="Create" icon="pi pi-check" @click="createRental" :loading="saving" />
      </template>
    </Dialog>

    <ConfirmDialog />
    <Toast />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import { useConfirm } from 'primevue/useconfirm';
import { useToast } from 'primevue/usetoast';
import { rentalService, houseService } from '@/service/api.service';

const store = useStore();
const confirm = useConfirm();
const toast = useToast();

const rentals = ref([]);
const availableHouses = ref([]);
const loading = ref(false);
const saving = ref(false);
const rentalDialog = ref(false);
const dialogMode = ref('view');
const selectedRental = ref(null);

const rentalForm = ref({
  house_id: null,
  tenant_email: '',
  start_date: null,
  end_date: null,
  monthly_rent: 0,
  deposit_amount: 0,
  terms: ''
});

const user = computed(() => store.getters['auth/user']);
const userRole = computed(() => user.value?.role || user.value?.roles?.[0]);
const canCreateRental = computed(() => ['landlord', 'admin'].includes(userRole.value));

const loadRentals = async () => {
  loading.value = true;
  try {
    const response = await rentalService.getAll();
    rentals.value = response.rental_agreements || response.rentals || [];
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to load rentals', life: 3000 });
  } finally {
    loading.value = false;
  }
};

const loadAvailableHouses = async () => {
  try {
    const response = await houseService.getAll({ status: 'available' });
    availableHouses.value = response.houses || [];
  } catch (error) {
    console.error('Failed to load houses:', error);
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'active': return 'success';
    case 'pending': return 'warn';
    case 'terminated': return 'danger';
    case 'expired': return 'secondary';
    default: return 'info';
  }
};

const formatNumber = (num) => new Intl.NumberFormat('en-ZM').format(num || 0);
const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A';

const openCreateDialog = () => {
  dialogMode.value = 'create';
  rentalForm.value = {
    house_id: null,
    tenant_email: '',
    start_date: null,
    end_date: null,
    monthly_rent: 0,
    deposit_amount: 0,
    terms: ''
  };
  loadAvailableHouses();
  rentalDialog.value = true;
};

const viewRental = (rental) => {
  dialogMode.value = 'view';
  selectedRental.value = rental;
  rentalDialog.value = true;
};

const createRental = async () => {
  saving.value = true;
  try {
    const data = {
      ...rentalForm.value,
      start_date: rentalForm.value.start_date?.toISOString?.() || rentalForm.value.start_date,
      end_date: rentalForm.value.end_date?.toISOString?.() || rentalForm.value.end_date
    };
    await rentalService.create(data);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Rental agreement created', life: 3000 });
    rentalDialog.value = false;
    loadRentals();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to create rental', life: 3000 });
  } finally {
    saving.value = false;
  }
};

const confirmTerminate = (rental) => {
  confirm.require({
    message: 'Are you sure you want to terminate this rental agreement?',
    header: 'Terminate Agreement',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await rentalService.terminate(rental.id);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Rental terminated', life: 3000 });
        loadRentals();
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to terminate', life: 3000 });
      }
    }
  });
};

onMounted(() => {
  loadRentals();
});
</script>

