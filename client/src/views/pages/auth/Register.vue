<script setup>
import FloatingConfigurator from '@/components/FloatingConfigurator.vue';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';

const store = useStore();
const router = useRouter();
const toast = useToast();

const fullName = ref('');
const email = ref('');
const phone = ref('');
const password = ref('');
const confirmPassword = ref('');
const role = ref('tenant');
const loading = ref(false);

const roleOptions = [
  { label: 'Tenant - Looking for a place to rent', value: 'tenant' },
  { label: 'Landlord - I have properties to rent', value: 'landlord' }
];

const register = async () => {
  // Validate form
  if (!fullName.value || !email.value || !password.value || !phone.value) {
    toast.add({ 
      severity: 'error', 
      summary: 'Error', 
      detail: 'Please fill all required fields', 
      life: 5000 
    });
    return;
  }
  
  // Validate email format
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email.value)) {
    toast.add({ 
      severity: 'error', 
      summary: 'Error', 
      detail: 'Please enter a valid email address', 
      life: 5000 
    });
    return;
  }
  
  // Validate password length
  if (password.value.length < 6) {
    toast.add({ 
      severity: 'error', 
      summary: 'Error', 
      detail: 'Password must be at least 6 characters', 
      life: 5000 
    });
    return;
  }
  
  // Validate password confirmation
  if (password.value !== confirmPassword.value) {
    toast.add({ 
      severity: 'error', 
      summary: 'Error', 
      detail: 'Passwords do not match', 
      life: 5000 
    });
    return;
  }
  
  loading.value = true;
  
  try {
    const registrationData = {
      full_name: fullName.value,
      email: email.value,
      phone: phone.value,
      password: password.value,
      role: role.value
    };
    
    await store.dispatch('auth/register', registrationData);
    
    toast.add({ 
      severity: 'success', 
      summary: 'Success', 
      detail: 'Registration successful! Redirecting...', 
      life: 2000 
    });
    
    // Redirect to dashboard after successful registration
    setTimeout(() => {
      router.push('/app/dashboard');
    }, 1000);
  } catch (error) {
    console.error('Registration error:', error);
    
    let errorMessage = 'Registration failed';
    
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message;
    } else if (error.userMessage) {
      errorMessage = error.userMessage;
    } else if (error.message) {
      errorMessage = error.message;
    }
    
    toast.add({
      severity: 'error',
      summary: 'Registration Failed',
      detail: errorMessage,
      life: 5000
    });
  } finally {
    loading.value = false;
  }
};

const loginNavigation = () => {
  router.push('/auth/login');
};
</script>

<template>
  <Toast position="top-center" />
  <FloatingConfigurator />
  <div class="bg-surface-200 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden">
    <div class="flex flex-col items-center justify-center">
      <div class="w-full py-12 px-8 sm:px-20" style="border-radius: 53px">
        <div class="text-center mb-8">
          <div class="flex items-center justify-center gap-3 mb-4">
            <i class="pi pi-home text-4xl text-primary"></i>
            <span class="text-surface-900 dark:text-surface-0 text-3xl font-bold">BondiHub</span>
          </div>
          <span class="text-muted-color font-medium">Create an account to get started</span>
        </div>
        <form @submit.prevent="register">
          <div>
            <!-- Full Name field -->
            <label for="fullName" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">Full Name *</label>
            <InputText id="fullName" type="text" placeholder="Enter your full name" class="w-full md:w-[30rem] mb-4" v-model="fullName" />
            
            <!-- Email field -->
            <label for="email" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">Email *</label>
            <InputText id="email" type="email" placeholder="Enter email address" class="w-full md:w-[30rem] mb-4" v-model="email" />
            
            <!-- Phone field -->
            <label for="phone" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">Phone *</label>
            <InputText id="phone" type="tel" placeholder="+260123456789" class="w-full md:w-[30rem] mb-4" v-model="phone" />
            
            <!-- Role selection -->
            <label for="role" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">I am a *</label>
            <Select id="role" v-model="role" :options="roleOptions" optionLabel="label" optionValue="value" 
              placeholder="Select your role" class="w-full md:w-[30rem] mb-4" />
            
            <!-- Password field -->
            <label for="password" class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">Password *</label>
            <Password id="password" v-model="password" placeholder="Enter password" :toggleMask="true" class="mb-4" fluid :feedback="true"></Password>
            
            <!-- Confirm Password field -->
            <label for="confirmPassword" class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">Confirm Password *</label>
            <Password id="confirmPassword" v-model="confirmPassword" placeholder="Confirm password" :toggleMask="true" class="mb-6" fluid :feedback="false"></Password>
            
            <div class="flex items-center justify-end mt-2 mb-6">
              <span @click="loginNavigation" class="font-medium no-underline ml-2 text-right cursor-pointer text-primary">Already have an account? Log in</span>
            </div>
            
            <Button type="submit" label="Create Account" class="w-full" :loading="loading" icon="pi pi-user-plus"></Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
