<template>
  <div class="payments-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">Payments</h1>
      <Button v-if="canMakePayment" label="Make Payment" icon="pi pi-credit-card" @click="openPaymentDialog" />
    </div>

    <!-- Stats Cards -->
    <div v-if="stats" class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <Panel>
        <template #header><span class="font-semibold">Total Payments</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-primary">{{ stats.total_count || 0 }}</p>
        </div>
      </Panel>
      <Panel>
        <template #header><span class="font-semibold">Total Amount</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-green-500">K{{ formatNumber(stats.total_amount) }}</p>
        </div>
      </Panel>
      <Panel>
        <template #header><span class="font-semibold">Pending</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-yellow-500">{{ stats.pending_count || 0 }}</p>
        </div>
      </Panel>
      <Panel>
        <template #header><span class="font-semibold">Completed</span></template>
        <div class="text-center">
          <p class="text-3xl font-bold text-blue-500">{{ stats.completed_count || 0 }}</p>
        </div>
      </Panel>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <!-- Payments Table -->
    <DataTable v-else :value="payments" :paginator="true" :rows="10" dataKey="id" responsiveLayout="scroll">
      <Column field="reference" header="Reference" sortable />
      
      <Column field="rental_agreement.house.title" header="Property">
        <template #body="{ data }">{{ data.rental_agreement?.house?.title || 'N/A' }}</template>
      </Column>
      
      <Column field="amount" header="Amount" sortable>
        <template #body="{ data }">
          <span class="font-semibold">K{{ formatNumber(data.amount) }}</span>
        </template>
      </Column>
      
      <Column field="payment_method" header="Method" sortable>
        <template #body="{ data }">
          <Tag :value="data.payment_method" severity="info" />
        </template>
      </Column>
      
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :severity="getStatusSeverity(data.status)" :value="data.status" />
        </template>
      </Column>
      
      <Column field="payment_date" header="Date" sortable>
        <template #body="{ data }">{{ formatDate(data.payment_date || data.created_at) }}</template>
      </Column>
      
      <Column header="Actions" style="width: 100px">
        <template #body="{ data }">
          <Button icon="pi pi-eye" outlined size="small" @click="viewPayment(data)" />
        </template>
      </Column>
      
      <template #empty>
        <div class="text-center py-8">
          <i class="pi pi-wallet text-4xl text-surface-300 mb-4"></i>
          <p class="text-surface-500">No payments found</p>
        </div>
      </template>
    </DataTable>

    <!-- Payment Dialog -->
    <Dialog v-model:visible="paymentDialog" :header="dialogMode === 'view' ? 'Payment Details' : 'Make Payment'" 
      :modal="true" :style="{ width: '500px' }">
      
      <div v-if="dialogMode === 'view' && selectedPayment" class="payment-detail">
        <div class="grid grid-cols-2 gap-4">
          <div><strong>Reference:</strong> {{ selectedPayment.reference }}</div>
          <div><strong>Amount:</strong> K{{ formatNumber(selectedPayment.amount) }}</div>
          <div><strong>Method:</strong> {{ selectedPayment.payment_method }}</div>
          <div><strong>Status:</strong> <Tag :severity="getStatusSeverity(selectedPayment.status)" :value="selectedPayment.status" /></div>
          <div><strong>Date:</strong> {{ formatDate(selectedPayment.payment_date || selectedPayment.created_at) }}</div>
          <div><strong>Property:</strong> {{ selectedPayment.rental_agreement?.house?.title || 'N/A' }}</div>
          <div class="col-span-2" v-if="selectedPayment.notes"><strong>Notes:</strong><br>{{ selectedPayment.notes }}</div>
        </div>
      </div>

      <div v-else class="payment-form">
        <div class="grid gap-4">
          <div class="field">
            <label class="block mb-2 font-medium">Rental Agreement *</label>
            <Select v-model="paymentForm.rental_agreement_id" :options="activeRentals" optionLabel="label" optionValue="id" 
              placeholder="Select rental" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Amount (ZMW) *</label>
            <InputNumber v-model="paymentForm.amount" mode="currency" currency="ZMW" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Payment Method *</label>
            <Select v-model="paymentForm.payment_method" :options="paymentMethods" class="w-full" />
          </div>
          
          <div class="field">
            <label class="block mb-2 font-medium">Notes</label>
            <Textarea v-model="paymentForm.notes" rows="3" class="w-full" />
          </div>
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" @click="paymentDialog = false" outlined />
        <Button v-if="dialogMode === 'create'" label="Submit Payment" icon="pi pi-check" @click="processPayment" :loading="saving" />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import { useToast } from 'primevue/usetoast';
import { paymentService, rentalService } from '@/service/api.service';

const store = useStore();
const toast = useToast();

const payments = ref([]);
const stats = ref(null);
const activeRentals = ref([]);
const loading = ref(false);
const saving = ref(false);
const paymentDialog = ref(false);
const dialogMode = ref('view');
const selectedPayment = ref(null);

const paymentForm = ref({
  rental_agreement_id: null,
  amount: 0,
  payment_method: 'mobile_money',
  notes: ''
});

const paymentMethods = ['mobile_money', 'bank_transfer', 'cash', 'card'];

const user = computed(() => store.getters['auth/user']);
const userRole = computed(() => user.value?.role || user.value?.roles?.[0]);
const canMakePayment = computed(() => ['tenant', 'admin'].includes(userRole.value));

const loadPayments = async () => {
  loading.value = true;
  try {
    const response = await paymentService.getAll();
    payments.value = response.payments || [];
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to load payments', life: 3000 });
  } finally {
    loading.value = false;
  }
};

const loadStats = async () => {
  try {
    stats.value = await paymentService.getStats();
  } catch (error) {
    console.error('Failed to load stats:', error);
  }
};

const loadActiveRentals = async () => {
  try {
    const response = await rentalService.getAll({ status: 'active' });
    activeRentals.value = (response.rental_agreements || response.rentals || []).map(r => ({
      id: r.id,
      label: `${r.house?.title || 'Property'} - K${r.monthly_rent}/month`
    }));
  } catch (error) {
    console.error('Failed to load rentals:', error);
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'completed': return 'success';
    case 'pending': return 'warn';
    case 'failed': return 'danger';
    default: return 'info';
  }
};

const formatNumber = (num) => new Intl.NumberFormat('en-ZM').format(num || 0);
const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A';

const openPaymentDialog = () => {
  dialogMode.value = 'create';
  paymentForm.value = {
    rental_agreement_id: null,
    amount: 0,
    payment_method: 'mobile_money',
    notes: ''
  };
  loadActiveRentals();
  paymentDialog.value = true;
};

const viewPayment = (payment) => {
  dialogMode.value = 'view';
  selectedPayment.value = payment;
  paymentDialog.value = true;
};

const processPayment = async () => {
  saving.value = true;
  try {
    await paymentService.process(paymentForm.value);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Payment submitted successfully', life: 3000 });
    paymentDialog.value = false;
    loadPayments();
    loadStats();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Payment failed', life: 3000 });
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  loadPayments();
  loadStats();
});
</script>

