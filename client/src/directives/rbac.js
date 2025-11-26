import store from '@/store';

/**
 * Vue directive for role-based access control
 * 
 * Usage:
 * <button v-role="'admin'">Admin Action</button>
 * <div v-roles="['admin', 'manager']">Admin or Manager Content</div>
 * <button v-permission="'delete_users'">Delete User</button>
 */

// Single role check
export const vRole = {
  mounted(el, binding) {
    const role = binding.value;
    if (!role) return;
    
    const hasRole = store.getters['auth/hasRole'](role);
    if (!hasRole) {
      el.style.display = 'none';
    }
  }
};

// Multiple roles check (any of the roles)
export const vRoles = {
  mounted(el, binding) {
    const roles = binding.value;
    if (!roles || !Array.isArray(roles) || roles.length === 0) return;
    
    const hasAnyRole = store.getters['auth/hasAnyRole'](roles);
    if (!hasAnyRole) {
      el.style.display = 'none';
    }
  }
};

// Permission check
export const vPermission = {
  mounted(el, binding) {
    const permission = binding.value;
    if (!permission) return;
    
    const hasPermission = store.getters['auth/hasPermission'](permission);
    if (!hasPermission) {
      el.style.display = 'none';
    }
  }
};

// Multiple permissions check (any of the permissions)
export const vPermissions = {
  mounted(el, binding) {
    const permissions = binding.value;
    if (!permissions || !Array.isArray(permissions) || permissions.length === 0) return;
    
    const hasAnyPermission = store.getters['auth/hasAnyPermission'](permissions);
    if (!hasAnyPermission) {
      el.style.display = 'none';
    }
  }
};

// Register directives
export default {
  install(app) {
    app.directive('role', vRole);
    app.directive('roles', vRoles);
    app.directive('permission', vPermission);
    app.directive('permissions', vPermissions);
  }
};
