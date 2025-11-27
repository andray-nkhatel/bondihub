import AppLayout from '@/layout/AppLayout.vue';
import { createRouter, createWebHashHistory } from 'vue-router';
import { authGuard } from './guard/auth.guard';

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        // Public home/landing page
        {
            path: '/',
            name: 'landing',
            component: () => import('@/views/pages/Landing.vue')
        },
        {
            path: '/auth/login',
            name: 'login',
            component: () => import('@/views/pages/auth/Login.vue')
        },
        {
            path: '/auth/register',
            name: 'register',
            component: () => import('@/views/pages/auth/Register.vue')
        },
        // All authenticated routes under /app
        {
            path: '/app',
            component: AppLayout,
            redirect: '/app/dashboard',
            children: [
                {
                    path: 'dashboard',
                    name: 'dashboard',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/dashboard/Dashboard.vue')
                },
                {   
                    path: 'profile',
                    name: 'profile',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/users/Profile.vue')
                },
                {
                    path: 'users',
                    name: 'users',
                    meta: {
                        requiresAuth: true,
                        roles: ['admin']
                    },
                    component: () => import('@/views/users/Users.vue')
                },
                // Houses / Properties
                {
                    path: 'houses',
                    name: 'houses',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/houses/Houses.vue')
                },
                // Rentals
                {
                    path: 'rentals',
                    name: 'rentals',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/rentals/Rentals.vue')
                },
                // Payments
                {
                    path: 'payments',
                    name: 'payments',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/payments/Payments.vue')
                },
                // Maintenance
                {
                    path: 'maintenance',
                    name: 'maintenance',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/maintenance/Maintenance.vue')
                },
                // Favorites (tenants)
                {
                    path: 'favorites',
                    name: 'favorites',
                    meta: {
                        requiresAuth: true,
                        roles: ['tenant', 'admin']
                    },
                    component: () => import('@/views/favorites/Favorites.vue')
                },
                // Notifications
                {
                    path: 'notifications',
                    name: 'notifications',
                    meta: {
                        requiresAuth: true,
                    },
                    component: () => import('@/views/notifications/Notifications.vue')
                }
            ]
        }
    ]
});

router.beforeEach(authGuard);
export default router;
