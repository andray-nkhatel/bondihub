<template>
    <div class="user-management">
      <!-- Header -->
      <div class="flex justify-content-between align-items-center mb-4">
        <h1 class="text-2xl font-bold m-0">User Management</h1>
        <Button 
          style="display: flex; margin-left: auto;"
          label="Add User" 
          icon="pi pi-plus" 
          @click="showAddDialog = true"
          class="p-button-primary"
        />
      </div>
  
      <!-- Filters and Search -->
      <Card class="mb-4">
        <template #content>
          <div class="grid">
            <div class="col-12 md:col-4 mb-2">
              <span class="p-input-icon-left w-full">
                <!-- <i class="pi pi-search" /> -->
                <InputText 
                  v-model="searchTerm" 
                  placeholder="Search users..." 
                  class="w-full"
                  @input="onSearch"
                />
              </span>
            </div>
            <div class="col-12 md:col-4 mb-2">
              <Dropdown
                v-model="selectedRole"
                :options="roleOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Filter by Role"
                class="w-full"
                @change="onRoleFilter"
              />
            </div>
            <div class="col-12 md:col-4 mb-2">
              <Dropdown
                v-model="selectedStatus"
                :options="statusOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Filter by Status"
                class="w-full"
                @change="onStatusFilter"
              />
            </div>
            <div class="col-12 md:col-2">
              <Button 
                label="Clear" 
                icon="pi pi-filter-slash" 
                @click="clearFilters"
                class="p-button-outlined w-full"
              />
            </div>
          </div>
        </template>
      </Card>
  
      <!-- Users DataTable -->
      <Card>
        <template #content>
          <DataTable 
            :value="filteredUsers" 
            :paginator="true" 
            :rows="10"
            :rowsPerPageOptions="[5, 10, 20, 50]"
            :loading="loading"
            stripedRows
            showGridlines
            sortMode="multiple"
            removableSort
            class="p-datatable-sm"
          >
            <Column field="id" header="ID" sortable style="width: 80px">
              <template #body="{ data }">
                <Badge :value="data.id" class="p-badge-secondary" />
              </template>
            </Column>
  
            <Column field="username" header="Username" sortable>
              <template #body="{ data }">
                <div class="flex align-items-center">
                  <Avatar 
                    :label="data.username.charAt(0).toUpperCase()" 
                    class="mr-2" 
                    size="small"
                    :style="{ backgroundColor: getAvatarColor(data.username) }"
                  />
                  <span class="font-medium">{{ data.username }}</span>
                </div>
              </template>
            </Column>
  
            <Column field="fullName" header="Full Name" sortable />
  
            <Column field="email" header="Email" sortable>
              <template #body="{ data }">
                <a :href="`mailto:${data.email}`" class="text-primary no-underline">
                  {{ data.email }}
                </a>
              </template>
            </Column>
  
            <Column field="role" header="Role" sortable>
              <template #body="{ data }">
                <Tag 
                  :value="roleValueToLabel(data.role)" 
                  :severity="getRoleSeverity(roleValueToLabel(data.role))"
                  class="text-sm"
                />
              </template>
            </Column>
  
            <Column field="isActive" header="Status" sortable>
              <template #body="{ data }">
                <Tag 
                  :value="data.isActive ? 'Active' : 'Inactive'" 
                  :severity="data.isActive ? 'success' : 'danger'"
                  :icon="data.isActive ? 'pi pi-check' : 'pi pi-times'"
                />
              </template>
            </Column>
  
            <Column field="createdAt" header="Created" sortable>
              <template #body="{ data }">
                <span :title="formatDate(data.createdAt)">
                  {{ formatRelativeDate(data.createdAt) }}
                </span>
              </template>
            </Column>
  
            <Column field="lastLoginAt" header="Last Login" sortable>
              <template #body="{ data }">
                <span :title="formatDate(data.lastLoginAt)">
                  {{ formatRelativeDate(data.lastLoginAt) }}
                </span>
              </template>
            </Column>
  
            <Column header="Actions" style="width: 120px">
              <template #body="{ data }">
                <div class="flex gap-2">
                  <Button 
                    icon="pi pi-pencil" 
                    size="small"
                    @click="editUser(data)"
                    v-tooltip.top="'Edit User'"
                  />
                  <Button 
                    icon="pi pi-eye" 
                    size="small"
                    severity="info"
                    @click="viewUser(data)"
                    v-tooltip.top="'View Details'"
                  />
                  <Button 
                    icon="pi pi-key" 
                    size="small"
                    severity="warning"
                    @click="resetUserPassword(data)"
                    v-tooltip.top="'Reset Password'"
                  />
                  <Button 
                    :icon="data.isActive ? 'pi pi-ban' : 'pi pi-check'" 
                    size="small"
                    :severity="data.isActive ? 'warning' : 'success'"
                    @click="deleteUser(data)"
                    :v-tooltip.top="data.isActive ? 'Deactivate User' : 'Activate User'"
                  />
                </div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>


      <!-- Password Reset Dialog -->
<Dialog 
  v-model:visible="showPasswordResetDialog" 
  header="Reset User Password"
  :modal="true" 
  class="p-fluid" 
  style="width: 400px"
>
  <div v-if="userToResetPassword" class="mb-4">
    <div class="flex align-items-center mb-3">
      <Avatar 
        :label="userToResetPassword.username.charAt(0).toUpperCase()" 
        class="mr-3" 
        :style="{ backgroundColor: getAvatarColor(userToResetPassword.username) }"
      />
      <div>
        <h4 class="m-0">{{ userToResetPassword.fullName }}</h4>
        <p class="text-600 m-0">@{{ userToResetPassword.username }}</p>
      </div>
    </div>
    
    <Message severity="warn" :closable="false">
      <p class="m-0">This will generate a new password for the user. The old password will be invalidated immediately.</p>
    </Message>
  </div>

  <div class="grid">
    <div class="col-12">
      <label for="newPassword" class="block text-900 font-medium mb-2">New Password</label>
      <div class="p-inputgroup">
        <InputText 
          id="newPassword"
          v-model="passwordResetForm.newPassword" 
          :type="showPassword ? 'text' : 'password'"
          placeholder="Enter new password"
          :class="{ 'p-invalid': passwordResetSubmitted && (!passwordResetForm.newPassword || passwordResetForm.newPassword.length < 6) }"
          class="w-full"
        />
        <Button 
          :icon="showPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"
          @click="showPassword = !showPassword"
          class="p-button-outlined"
          type="button"
        />
      </div>
      <small v-show="passwordResetSubmitted && !passwordResetForm.newPassword" class="p-error">
        New password is required.
      </small>
      <small v-show="passwordResetSubmitted && passwordResetForm.newPassword && passwordResetForm.newPassword.length < 6" class="p-error">
        Password must be at least 6 characters long.
      </small>
    </div>

    <div class="col-12 mt-2">
      <label for="confirmPassword" class="block text-900 font-medium mb-2">Confirm Password</label>
      <InputText 
        id="confirmPassword"
        v-model="passwordResetForm.confirmPassword" 
        :type="showConfirmPassword ? 'text' : 'password'"
        placeholder="Confirm new password"
        :class="{ 'p-invalid': passwordResetSubmitted && (!passwordResetForm.confirmPassword || passwordResetForm.newPassword !== passwordResetForm.confirmPassword) }"
        class="w-full"
      />
      <small v-show="passwordResetSubmitted && !passwordResetForm.confirmPassword" class="p-error">
        Please confirm the password.
      </small>
      <small v-show="passwordResetSubmitted && passwordResetForm.confirmPassword && passwordResetForm.newPassword !== passwordResetForm.confirmPassword" class="p-error">
        Passwords do not match.
      </small>
    </div>

    <div class="col-12 mt-3">
      <div class="field-checkbox">
        <Checkbox 
          id="notifyUser" 
          v-model="passwordResetForm.notifyUser" 
          :binary="true" 
        />
        <label class="ml-2" for="notifyUser">Send password reset notification to user</label>
      </div>
    </div>
  </div>

  <template #footer>
    <Button 
      label="Cancel" 
      icon="pi pi-times" 
      @click="hidePasswordResetDialog" 
      class="p-button-text" 
    />
    <Button 
      label="Reset Password" 
      icon="pi pi-key" 
      @click="confirmPasswordReset" 
      class="p-button-danger"
      :loading="passwordResetLoading"
    />
  </template>
</Dialog>
  
      <!-- Add/Edit User Dialog -->
      <Dialog 
        v-model:visible="showAddDialog" 
        :header="editingUser ? 'Edit User' : 'Add New User'"
        :modal="true" 
        class="p-fluid" 
        style="width: 400px"
      >
        <div class="grid">
          <div class="col-12">
            <label for="username" class="block text-900 font-medium mb-2">Username</label>
            <InputText 
              class="w-full"
              id="username"
              v-model="userForm.username" 
              :class="{ 'p-invalid': submitted && !userForm.username }" 
            />
            <small v-show="submitted && !userForm.username" class="p-error">Username is required.</small>
          </div>
  
          <div class="col-12 mt-2">
            <label for="fullName" class="block text-900 font-medium mb-2">Full Name</label>
            <InputText 
            class="w-full"
              id="fullName"
              v-model="userForm.fullName" 
              :class="{ 'p-invalid': submitted && !userForm.fullName }" 
            />
            <small v-show="submitted && !userForm.fullName" class="p-error">Full Name is required.</small>
          </div>
  
          <div class="col-12 mt-2">
            <label for="email" class="block text-900 font-medium mb-2">Email</label>
            <InputText 
              class="w-full"
              id="email"
              v-model="userForm.email" 
              type="email"
              :class="{ 'p-invalid': submitted && (!userForm.email || !isValidEmail(userForm.email)) }" 
            />
            <small v-show="submitted && !userForm.email" class="p-error">Email is required.</small>
            <small v-show="submitted && userForm.email && !isValidEmail(userForm.email)" class="p-error">Please enter a valid email.</small>
          </div>

          <div v-if="!editingUser" class="col-12 mt-2">
          <label for="password" class="block text-900 font-medium mb-2">Password</label>
          <InputText 
            class="w-full"
            id="password"
            v-model="userForm.password" 
            type="password"
            :class="{ 'p-invalid': submitted && (!userForm.password || userForm.password.length < 6) }" 
          />
          <small v-show="submitted && !userForm.password" class="p-error">Password is required.</small>
          <small v-show="submitted && userForm.password && userForm.password.length < 6" class="p-error">Password must be at least 6 characters long.</small>
        </div>
  
  
          <div class="col-12 mt-2">
            <label for="role" class="block text-900 font-medium mb-2">Role</label>
            <Dropdown
              id="role"
              class="w-full"
              v-model="userForm.role"
              :options="roleOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Select Role"
              :class="{ 'p-invalid': submitted && !userForm.role }"
            />
            <small v-show="submitted && !userForm.role" class="p-error">Role is required.</small>
          </div>
  
          <div class="col-12 mt-5">
            <div class="field-checkbox">
              <Checkbox 
                id="isActive" 
                v-model="userForm.isActive" 
                :binary="true" 
              />
              <label class="ml-2" for="isActive">Active User</label>
            </div>
          </div>
        </div>
  
        <template #footer>
          <Button 
            label="Cancel" 
            icon="pi pi-times" 
            @click="hideDialog" 
            class="p-button-text" 
          />
          <Button 
            :label="editingUser ? 'Update' : 'Save'" 
            icon="pi pi-check" 
            @click="saveUser" 
            autofocus 
          />
        </template>
      </Dialog>
  
      <!-- User Details Dialog -->
      <Dialog 
        v-model:visible="showDetailsDialog" 
        header="User Details"
        :modal="true" 
        style="width: 600px"
      >
        <div v-if="selectedUser" class="grid">
          <div class="col-12 text-center mb-4">
            <Avatar 
              :label="selectedUser.username.charAt(0).toUpperCase()" 
              size="xlarge"
              :style="{ backgroundColor: getAvatarColor(selectedUser.username) }"
            />
            <h3 class="mt-3 mb-1">{{ selectedUser.fullName }}</h3>
            <p class="text-600 m-0">@{{ selectedUser.username }}</p>
          </div>
  
          <div class="col-6">
            <strong>ID:</strong>
            <p>{{ selectedUser.id }}</p>
          </div>
  
          <div class="col-6">
            <strong>Email:</strong>
            <p>{{ selectedUser.email }}</p>
          </div>
  
          <div class="col-6">
            <strong>Role:</strong>
            <div class="flex flex-wrap gap-1 mt-1">
              <Tag 
                :value="roleValueToLabel(selectedUser.role)"
                :severity="getRoleSeverity(roleValueToLabel(selectedUser.role))"
              />
            </div>
          </div>
  
          <div class="col-6">
            <strong>Status:</strong>
            <p>
              <Tag 
                :value="selectedUser.isActive ? 'Active' : 'Inactive'" 
                :severity="selectedUser.isActive ? 'success' : 'danger'"
                :icon="selectedUser.isActive ? 'pi pi-check' : 'pi pi-times'"
              />
            </p>
          </div>
  
          <div class="col-6">
            <strong>Created:</strong>
            <p>{{ formatDate(selectedUser.createdAt) }}</p> 
         
          </div>
  
          <div class="col-6">
            <strong>Last Login:</strong>
            <p>{{ formatDate(selectedUser.lastLoginAt) }}</p>
          </div>
        </div>
  
        <template #footer>
          <Button 
            label="Close" 
            icon="pi pi-times" 
            @click="showDetailsDialog = false" 
            class="p-button-text" 
          />
          <Button 
            label="Edit User" 
            icon="pi pi-pencil" 
            @click="editUserFromDetails" 
          />
        </template>
      </Dialog>
    </div>
    <ConfirmDialog></ConfirmDialog>
    <Toast></Toast>
  </template>
  
  <script>
import { userService } from '@/service/api.service';
import { formatDateTime, formatDate } from '@/service/dateUtils';
import MultiSelect from 'primevue/multiselect';
import { useConfirm } from 'primevue/useconfirm';
import { useToast } from 'primevue/usetoast';
import { computed, onMounted, ref } from 'vue';

export default {
  name: 'User',
  components: {
    MultiSelect
  },
  setup() {
    // Reactive data
    const users = ref([])
    const filteredUsers = ref([])
    const loading = ref(false)
    const searchTerm = ref('')
    const selectedRole = ref(null)
    const selectedStatus = ref(null)
    const confirm = useConfirm();
    const toast = useToast();
  
    
    // Dialog states
    const showAddDialog = ref(false)
    const showDetailsDialog = ref(false)
    const submitted = ref(false)
    const editingUser = ref(false)
    const selectedUser = ref(null)

    const showPasswordResetDialog = ref(false)
    const passwordResetSubmitted = ref(false)
    const passwordResetLoading = ref(false)
    const userToResetPassword = ref(null)
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)

    // Form data (add to existing form data)
const passwordResetForm = ref({
  newPassword: '',
  confirmPassword: '',
  notifyUser: true
})

// Methods (add to existing methods)
const resetUserPassword = (user) => {
  userToResetPassword.value = user
  passwordResetForm.value = {
    newPassword: '',
    confirmPassword: '',
    notifyUser: true
  }
  showPasswordResetDialog.value = true
  passwordResetSubmitted.value = false
  showPassword.value = false
  showConfirmPassword.value = false
}

const hidePasswordResetDialog = () => {
  showPasswordResetDialog.value = false
  passwordResetSubmitted.value = false
  passwordResetLoading.value = false
  userToResetPassword.value = null
  passwordResetForm.value = {
    newPassword: '',
    confirmPassword: '',
    notifyUser: true
  }
  showPassword.value = false
  showConfirmPassword.value = false
}

const confirmPasswordReset = () => {
  passwordResetSubmitted.value = true
  
  // Validation
  if (!passwordResetForm.value.newPassword || 
      passwordResetForm.value.newPassword.length < 6 ||
      !passwordResetForm.value.confirmPassword ||
      passwordResetForm.value.newPassword !== passwordResetForm.value.confirmPassword) {
    return
  }

  // Show confirmation dialog
  confirm.require({
    message: `Are you sure you want to reset the password for "${userToResetPassword.value.fullName}"? This action cannot be undone and will immediately invalidate their current password.`,
    header: 'Reset Password Confirmation',
    icon: 'pi pi-exclamation-triangle',
    rejectClass: 'p-button-secondary p-button-outlined',
    rejectLabel: 'Cancel',
    acceptClass: 'p-button-danger',
    acceptLabel: 'Reset Password',
    accept: async () => {
      await executePasswordReset()
    }
  })
}

const executePasswordReset = async () => {
  passwordResetLoading.value = true
  
  try {
    // Validate inputs
    if (!userToResetPassword.value?.id) {
      throw new Error('No user selected for password reset')
    }

    const newPassword = passwordResetForm.value.newPassword?.trim()
    if (!newPassword) {
      throw new Error('Password cannot be empty')
    }

    if (typeof newPassword !== 'string') {
      throw new Error('New password must be a string')
    }

    // Validate password strength on frontend too
    if (newPassword.length < 8) {
      throw new Error('Password must be at least 8 characters long')
    }

    // Match your backend DTO structure exactly
    const resetData = {
      newPassword: newPassword
    }

    console.log('🔍 Resetting password for user:', userToResetPassword.value.fullName)
    console.log('🔍 User ID:', userToResetPassword.value.id)
    console.log('🔍 Reset data structure:', resetData)
    
    const response = await userService.resetPassword(userToResetPassword.value.id, resetData)
    
    console.log('✅ Password reset response:', response)
    console.log('✅ Password reset successfully for user:', userToResetPassword.value.fullName)
    
    // Show success toast
    toast.add({
      severity: 'success',
      summary: 'Password Reset Successfully',
      detail: `Password has been reset for "${userToResetPassword.value.fullName}".`,
      life: 4000
    })
    
    hidePasswordResetDialog()
    
    // Optional: Clear the form
    passwordResetForm.value.newPassword = ''
    
  } catch (error) {
    console.error('❌ Error resetting password:', error)
    
    let errorMessage = 'Failed to reset password'
    
    // Handle different types of errors
    if (error.response?.status === 404) {
      errorMessage = 'User not found'
    } else if (error.response?.status === 400) {
      // Handle validation errors
      const validationErrors = error.response?.data?.errors
      if (validationErrors) {
        const passwordErrors = validationErrors.NewPassword || validationErrors.newPassword
        if (passwordErrors && passwordErrors.length > 0) {
          errorMessage = passwordErrors[0]
        } else {
          errorMessage = error.response?.data?.message || 'Invalid request data'
        }
      } else {
        errorMessage = error.response?.data?.message || 'Invalid request data'
      }
    } else if (error.response?.status === 403) {
      errorMessage = 'You do not have permission to reset passwords'
    } else if (error.response?.status === 401) {
      errorMessage = 'You are not authorized. Please log in again.'
    } else if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }
    
    console.error('❌ Error details:', {
      status: error.response?.status,
      data: error.response?.data,
      message: error.message
    })
    
    // Show error toast
    toast.add({
      severity: 'error',
      summary: 'Password Reset Failed',
      detail: errorMessage,
      life: 6000
    })
  } finally {
    passwordResetLoading.value = false
  }
}

     
      
      // Form data
      const userForm = ref({
        username: '',
        fullName: '',
        email: '',
        role: '',
        password: '',
        isActive: true
      })
  
      // Options
      const roleOptions = ref([
        { label: 'Admin', value: 'Admin' },
        { label: 'Teacher', value: 'Teacher' },
        { label: 'Staff', value: 'Staff' }
      ])
  
      const statusOptions = ref([
        { label: 'Active', value: true },
        { label: 'Inactive', value: false }
      ])
  
      // Computed
      const computedFilteredUsers = computed(() => {
        let result = [...users.value]
        
        // Search filter
        if (searchTerm.value) {
          const search = searchTerm.value.toLowerCase()
          result = result.filter(user => 
            user.username.toLowerCase().includes(search) ||
            user.fullName.toLowerCase().includes(search) ||
            user.email.toLowerCase().includes(search) ||
            (user.role && user.role.toLowerCase().includes(search))
          )
        }
        
        // Role filter
        if (selectedRole.value !== null && selectedRole.value !== '') {
          result = result.filter(user => user.role === selectedRole.value)
        }
        
        // Status filter
        if (selectedStatus.value !== null) {
          result = result.filter(user => user.isActive === selectedStatus.value)
        }
        
        return result
      })
  
      // Methods
      const fetchUsers = async () => {
        loading.value = true
        try {
          // Replace with your actual API call
          const response = await userService.getAll();
          users.value = response.map(u => ({
            ...u,
            // role is already a string
          }))
          filteredUsers.value = [...users.value]
        } catch (error) {
          console.error('Error fetching users:', error)
        } finally {
          loading.value = false
        }
      }
  
      const onSearch = () => {
        filteredUsers.value = computedFilteredUsers.value
      }
  
      const onRoleFilter = () => {
        filteredUsers.value = computedFilteredUsers.value
      }
  
      const onStatusFilter = () => {
        filteredUsers.value = computedFilteredUsers.value
      }
  
      const clearFilters = () => {
        searchTerm.value = ''
        selectedRole.value = null
        selectedStatus.value = null
        filteredUsers.value = [...users.value]
      }
  
      const editUser = (user) => {
        userForm.value = { ...user }
        editingUser.value = true
        showAddDialog.value = true
      }
  
      const viewUser = (user) => {
        selectedUser.value = user
        showDetailsDialog.value = true
      }
  
      const editUserFromDetails = () => {
        showDetailsDialog.value = false
        editUser(selectedUser.value)
      }
  
      const deleteUser = async (user) => {
        confirm.require({
        message: `Are you sure you want to delete user "${user.fullName}"? This action cannot be undone.`,
        header: 'Delete User Confirmation',
        icon: 'pi pi-exclamation-triangle',
        rejectClass: 'p-button-secondary p-button-outlined',
        rejectLabel: 'Cancel',
        acceptClass: 'p-button-danger',
        acceptLabel: 'Delete',
        accept: async () => {
            try {
                // API call to delete user
                await userService.delete(user.id);
                
                // Remove user from local arrays
                const index = users.value.findIndex(u => u.id === user.id);
                if (index !== -1) {
                    users.value.splice(index, 1);
                }
                
                // Update filtered users
                filteredUsers.value = computedFilteredUsers.value;
                
                console.log('✅ User deleted successfully:', user.fullName);
                
                // Show success toast
                toast.add({
                    severity: 'success',
                    summary: 'User Deleted',
                    detail: `User "${user.fullName}" has been deleted successfully.`,
                    life: 3000
                });
                
            } catch (error) {
                console.error('❌ Error deleting user:', error);
                
                // Show error toast
                const errorMessage = error.response?.data?.message || 'Failed to delete user';
                toast.add({
                    severity: 'error',
                    summary: 'Delete Failed',
                    detail: errorMessage,
                    life: 5000
                });
            }
        },
        reject: () => {
            // Optional: Show cancelled message
            toast.add({
                severity: 'info',
                summary: 'Cancelled',
                detail: 'User deletion cancelled',
                life: 2000
            });
        }
    });
};
  
      
      const saveUser = async () => {
    submitted.value = true
    if (!userForm.value.username || !userForm.value.fullName ||
        !userForm.value.email || !userForm.value.role ||
        !isValidEmail(userForm.value.email)) {
        console.log('❌ Frontend validation failed');
        return
    }
    // Only validate password for NEW users
    if (!editingUser.value && (!userForm.value.password || userForm.value.password.length < 6)) {
        alert('Password is required and must be at least 6 characters long');
        submitted.value = false;
        return;
    }
    const userData = {
        username: userForm.value.username.trim(),
        fullName: userForm.value.fullName.trim(),
        email: userForm.value.email.trim().toLowerCase(),
        role: userForm.value.role ? userForm.value.role.trim() : '',
    };
    if (!editingUser.value) {
        userData.password = userForm.value.password;
    } else {
        userData.isActive = userForm.value.isActive;
    }
    // Validate role value before sending
    const allowedRoles = ['Admin', 'Teacher', 'Staff'];
    if (!allowedRoles.includes(userData.role)) {
        alert('Role must be one of: Admin, Teacher, Staff');
        submitted.value = false;
        return;
    }
    try {
        let savedUser;
        if (editingUser.value) {
            savedUser = await userService.update(userForm.value.id, userData);
            const index = users.value.findIndex(u => u.id === userForm.value.id);
            if (index !== -1) {
                users.value[index] = savedUser;
            }
            toast.add({
                    severity: 'success',
                    summary: 'User Successfully Updated',
                    detail: `User "${savedUser.fullName}" has been updated successfully.`,
                    life: 3000
                });
        } else {
            // Send plain object, not wrapped
            savedUser = await userService.create({
              username: userForm.value.username.trim(),
              fullName: userForm.value.fullName.trim(),
              email: userForm.value.email.trim().toLowerCase(),
              role: userForm.value.role ? userForm.value.role.trim() : '',
              password: userForm.value.password
            });
            users.value.push(savedUser);
             toast.add({
                    severity: 'success',
                    summary: 'User Successfully Created',
                    detail: `User "${savedUser.fullName}" has been created successfully.`,
                    life: 3000
                });
        }
        filteredUsers.value = computedFilteredUsers.value;
        hideDialog();
        console.log('✅ User saved successfully:', savedUser);
    } catch (error) {
        console.error('❌ Failed to save user:', error);
        const data = error.response?.data;
        if (data) {
            if (data.message) {
                alert('Backend error: ' + data.message);
            } else if (data.errors) {
                alert('Validation errors: ' + JSON.stringify(data.errors));
            } else {
                alert('Unknown backend error: ' + JSON.stringify(data));
            }
        }
        submitted.value = false;
    }
}

      const hideDialog = () => {
        showAddDialog.value = false
        submitted.value = false
        editingUser.value = false
        userForm.value = {
          username: '',
          fullName: '',
          email: '',
          role: '',
          password:'',
          isActive: true
        }
      }
  
      const isValidEmail = (email) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        return emailRegex.test(email)
      }
  
      // Use date utility with Zambia timezone
      const formatDateLocal = formatDateTime

      const formatRelativeDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  if (isNaN(date.getTime())) return 'Invalid date'

  const now = new Date()
  const diffTime = now - date
  const diffSeconds = Math.floor(diffTime / 1000)
  const diffMinutes = Math.floor(diffTime / (1000 * 60))
  const diffHours = Math.floor(diffTime / (1000 * 60 * 60))
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

  // Handle future dates
  if (diffTime < 0) {
    const absDiffSeconds = Math.abs(diffSeconds)
    const absDiffMinutes = Math.abs(diffMinutes)
    const absDiffHours = Math.abs(diffHours)
    const absDiffDays = Math.abs(diffDays)
    
    if (absDiffSeconds < 60) return 'in a few seconds'
    if (absDiffMinutes < 60) return `in ${absDiffMinutes} minute${absDiffMinutes === 1 ? '' : 's'}`
    if (absDiffHours < 24) return `in ${absDiffHours} hour${absDiffHours === 1 ? '' : 's'}`
    if (absDiffDays === 1) return 'tomorrow'
    if (absDiffDays < 7) return `in ${absDiffDays} days`
    return `in ${Math.floor(absDiffDays / 7)} weeks`
  }

      // Past dates
      if (diffSeconds < 10) return 'just now'
      if (diffSeconds < 60) return `${diffSeconds} seconds ago`
      if (diffMinutes === 1) return '1 minute ago'
      if (diffMinutes < 60) return `${diffMinutes} minutes ago`
      if (diffHours === 1) return '1 hour ago'
      if (diffHours < 24) return `${diffHours} hours ago`
      if (diffDays === 1) return 'yesterday'
      if (diffDays < 7) return `${diffDays} days ago`
      if (diffDays < 30) {
        const weeks = Math.floor(diffDays / 7)
        return `${weeks} week${weeks === 1 ? '' : 's'} ago`
      }
      if (diffDays < 365) {
        const months = Math.floor(diffDays / 30)
        return `${months} month${months === 1 ? '' : 's'} ago`
      }
      const years = Math.floor(diffDays / 365)
      return `${years} year${years === 1 ? '' : 's'} ago`
    }

    const getRoleSeverity = (role) => {
      const severityMap = {
        admin: 'danger',
        manager: 'warning',
        user: 'info',
        guest: 'secondary'
      }
      return severityMap[role?.toLowerCase()] || 'info'
    }

    const getAvatarColor = (username) => {
      const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8']
      const index = username.charCodeAt(0) % colors.length
      return colors[index]
    }

    // Add a helper to map role value to label
    const roleValueToLabel = (value) => {
      const found = roleOptions.value.find(opt => opt.value === value)
      return found ? found.label : value
    }

    // Lifecycle
    onMounted(() => {
      fetchUsers()
    })

    return {
      users,
      filteredUsers,
      loading,
      searchTerm,
      selectedRole,
      selectedStatus,
      showAddDialog,
      showDetailsDialog,
      submitted,
      editingUser,
      selectedUser,
      userForm,
      roleOptions,
      statusOptions,
      onSearch,
      onRoleFilter,
      onStatusFilter,
      clearFilters,
      editUser,
      viewUser,
      editUserFromDetails,
      deleteUser,
      saveUser,
      hideDialog,
      isValidEmail,
      formatDate: formatDateLocal,
      formatRelativeDate,
      getRoleSeverity,
      getAvatarColor,
      showPasswordResetDialog,
      passwordResetSubmitted,
      passwordResetLoading,
      userToResetPassword,
      showPassword,
      showConfirmPassword,
      passwordResetForm,
      resetUserPassword,
      hidePasswordResetDialog,
      confirmPasswordReset,
      executePasswordReset,
      roleValueToLabel
    }
  }
}
</script>

<style scoped>
.user-management {
  padding: 1rem;
}

.p-datatable .p-datatable-tbody > tr > td {
  padding: 0.75rem;
}

.p-avatar {
  color: white;
  font-weight: bold;
}

.text-primary {
  color: var(--primary-color) !important;
}

.no-underline {
  text-decoration: none;
}

.no-underline:hover {
  text-decoration: underline;
}
</style>