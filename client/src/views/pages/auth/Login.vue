<script setup>
import FloatingConfigurator from '@/components/FloatingConfigurator.vue';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';

const store = useStore();
const router = useRouter();
const toast = useToast();

const email = ref('');
const password = ref('');
const loading = ref(false);

const login = async () => {
    if (!email.value || !password.value) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: 'Please enter email and password',
            life: 3000
        });
        return;
    }

    loading.value = true;

    try {
        await store.dispatch('auth/login', {
            email: email.value,
            password: password.value
        });

        const queryRedirect = router.currentRoute.value.query.redirect;
        const redirectPath = queryRedirect || '/app/dashboard';

        toast.add({
            severity: 'success',
            summary: 'Welcome!',
            detail: 'You have been logged in successfully.',
            life: 2000
        });

        router.push(redirectPath);
        password.value = '';
    } catch (error) {
        console.error('Login error:', error);

        if (error.isCorsError) {
            toast.add({
                severity: 'error',
                summary: 'Connection Error',
                detail: 'Cannot connect to the server. Please try again later.',
                life: 5000
            });
        } else if (error.response && error.response.status === 401) {
            toast.add({
                severity: 'error',
                summary: 'Login Failed',
                detail: 'Invalid email or password.',
                life: 3000
            });
        } else {
            let errorMessage = 'Invalid credentials';
            if (error.response?.data?.message) {
                errorMessage = error.response.data.message;
            } else if (error.userMessage) {
                errorMessage = error.userMessage;
            } else if (error.message) {
                errorMessage = error.message;
            }
            toast.add({
                severity: 'error',
                summary: 'Login Failed',
                detail: errorMessage,
                life: 3000
            });
        }
        password.value = '';
    } finally {
        loading.value = false;
    }
};

const registerNavigation = () => {
    router.push('/auth/register');
};
</script>

<template>
    <Toast position="top-center" />
    <FloatingConfigurator />
    <div class="bg-surface-200 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden">
        <div class="flex flex-col items-center justify-center w-full">
            <div class="w-full max-w-lg py-20 px-4 sm:px-8" style="border-radius: 53px">
                <div class="text-center mb-8">
                    <div class="flex items-center justify-center gap-3 mb-4">
                        <i class="pi pi-home text-4xl text-primary"></i>
                        <span class="text-surface-900 dark:text-surface-0 text-3xl font-bold">BondiHub</span>
                    </div>
                    <span class="text-muted-color font-medium">Sign in to continue</span>
                </div>
                <form @submit.prevent="login">
                    <div>
                        <label for="email" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">Email</label>
                        <InputText id="email" type="email" placeholder="Enter your email" class="w-full mb-6" v-model="email" />
                        
                        <label for="password" class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">Password</label>
                        <Password id="password" v-model="password" placeholder="Enter your password" :toggleMask="true" class="w-full mb-4" fluid :feedback="false"></Password>
                        
                        <div class="flex items-center justify-end mt-2 mb-6">
                            <span @click="registerNavigation" class="font-medium no-underline cursor-pointer text-primary hover:underline">
                                Don't have an account? Register
                            </span>
                        </div>
                        
                        <Button type="submit" label="Sign In" class="w-full" :loading="loading" icon="pi pi-sign-in"></Button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>
