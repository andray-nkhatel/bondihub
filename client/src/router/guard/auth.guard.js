import store from '@/store';

/**
 * Authentication guard for Vue Router
 * 
 * Usage in router:
 * {
 *   path: '/admin',
 *   component: AdminDashboard,
 *   meta: { 
 *     requiresAuth: true,
 *     roles: ['Admin'],
 *     permissions: ['view_dashboard']
 *   }
 * }
 */
export const authGuard = async (to, from, next) => {
  // Check if route requires authentication
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  
  if (!requiresAuth) {
    return next();
  }
  
  // Make sure token is valid and refresh if needed
  await store.dispatch('auth/checkTokenAndRefreshIfNeeded');
  
  // Check if user is authenticated
  const isAuthenticated = store.getters['auth/isAuthenticated'];
  
  if (!isAuthenticated) {
    // Redirect to login page with return url
    return next({
      path: '/auth/login',
      query: { redirect: to.fullPath }
    });
  }
  
  // Check for required roles
  const requiredRoles = to.matched.reduce((roles, record) => {
    if (record.meta.roles) {
      return [...roles, ...record.meta.roles];
    }
    return roles;
  }, []);
  
  if (requiredRoles.length > 0) {
    const hasRequiredRole = store.getters['auth/hasAnyRole'](requiredRoles);
    if (!hasRequiredRole) {
      return next({ path: '/unauthorized' });
    }
  }
  
  // Check for required permissions
  const requiredPermissions = to.matched.reduce((permissions, record) => {
    if (record.meta.permissions) {
      return [...permissions, ...record.meta.permissions];
    }
    return permissions;
  }, []);
  
  if (requiredPermissions.length > 0) {
    const hasRequiredPermission = store.getters['auth/hasAnyPermission'](requiredPermissions);
    if (!hasRequiredPermission) {
      return next({ path: '/unauthorized' });
    }
  }
  
  // All checks passed, proceed to route
  next();
};
