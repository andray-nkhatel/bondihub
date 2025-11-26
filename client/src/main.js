import '@/assets/styles.scss';
import Aura from '@primeuix/themes/aura';
import PrimeVue from 'primevue/config';
import ConfirmationService from 'primevue/confirmationservice';
import ToastService from 'primevue/toastservice';
import { createApp } from 'vue';
import App from './App.vue';
import RbacDirectives from './directives/rbac';
import router from './router';
import store from './store';

const CustomAura = {
    ...Aura,
};

// Update the primary color mapping (customize as needed)
if (CustomAura.semantic && CustomAura.semantic.primary) {
    CustomAura.semantic.primary = {
      "50": "{blue.50}",
      "100": "{blue.100}",
      "200": "{blue.200}",
      "300": "{blue.300}",
      "400": "{blue.400}",
      "500": "{blue.500}",
      "600": "{blue.600}",
      "700": "{blue.700}",
      "800": "{blue.800}",
      "900": "{blue.900}",
      "950": "{blue.950}"
    };
}

const app = createApp(App);

// Setup store and interceptors on app start
store.dispatch('auth/setupInterceptors');

// Register RBAC directives
app.use(RbacDirectives);
app.use(router);
app.use(PrimeVue, {
    theme: {
        preset: CustomAura,
        options: {
            darkModeSelector: '.app-dark'
        }
    }
});
app.use(ToastService);
app.use(ConfirmationService);

app.use(store).use(router).mount('#app');
