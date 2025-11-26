<template>
  <div class="notifications-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">
        <i class="pi pi-bell mr-2"></i>Notifications
        <Badge v-if="unreadCount > 0" :value="unreadCount" severity="danger" class="ml-2" />
      </h1>
      <Button v-if="notifications.length > 0" label="Mark All Read" icon="pi pi-check-circle" outlined @click="markAllRead" />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <!-- Notifications List -->
    <div v-else-if="notifications.length > 0" class="space-y-3">
      <div v-for="notif in notifications" :key="notif.id" 
        :class="['notification-item p-4 rounded-lg border cursor-pointer transition-all', notif.is_read ? 'bg-surface-50 dark:bg-surface-800' : 'bg-blue-50 dark:bg-blue-900/20 border-blue-200']"
        @click="viewNotification(notif)">
        <div class="flex items-start gap-4">
          <div :class="['notification-icon p-3 rounded-full', getIconClass(notif.type)]">
            <i :class="['pi', getIcon(notif.type)]"></i>
          </div>
          <div class="flex-1">
            <div class="flex justify-between items-start">
              <h3 class="font-semibold text-surface-900 dark:text-surface-0">{{ notif.title }}</h3>
              <span class="text-sm text-surface-500">{{ formatDate(notif.created_at) }}</span>
            </div>
            <p class="text-surface-600 dark:text-surface-300 mt-1">{{ notif.message }}</p>
          </div>
          <div class="flex gap-2">
            <Button v-if="!notif.is_read" icon="pi pi-check" text size="small" @click.stop="markAsRead(notif.id)" v-tooltip="'Mark as read'" />
            <Button icon="pi pi-trash" text size="small" severity="danger" @click.stop="deleteNotification(notif.id)" v-tooltip="'Delete'" />
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <i class="pi pi-bell-slash text-6xl text-surface-300 mb-4"></i>
      <p class="text-xl text-surface-500">No notifications</p>
      <p class="text-surface-400">You're all caught up!</p>
    </div>

    <Toast />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { notificationService } from '@/service/api.service';

const toast = useToast();

const notifications = ref([]);
const loading = ref(false);

const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length);

const loadNotifications = async () => {
  loading.value = true;
  try {
    const response = await notificationService.getAll();
    notifications.value = response.notifications || [];
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to load notifications', life: 3000 });
  } finally {
    loading.value = false;
  }
};

const getIcon = (type) => {
  switch (type) {
    case 'payment': return 'pi-credit-card';
    case 'rental': return 'pi-file';
    case 'maintenance': return 'pi-wrench';
    case 'review': return 'pi-star';
    case 'house': return 'pi-home';
    default: return 'pi-info-circle';
  }
};

const getIconClass = (type) => {
  switch (type) {
    case 'payment': return 'bg-green-100 text-green-600';
    case 'rental': return 'bg-blue-100 text-blue-600';
    case 'maintenance': return 'bg-yellow-100 text-yellow-600';
    case 'review': return 'bg-purple-100 text-purple-600';
    case 'house': return 'bg-cyan-100 text-cyan-600';
    default: return 'bg-surface-100 text-surface-600';
  }
};

const formatDate = (date) => {
  if (!date) return '';
  const d = new Date(date);
  const now = new Date();
  const diff = now - d;
  
  if (diff < 60000) return 'Just now';
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}d ago`;
  return d.toLocaleDateString();
};

const viewNotification = async (notif) => {
  if (!notif.is_read) {
    await markAsRead(notif.id);
  }
};

const markAsRead = async (id) => {
  try {
    await notificationService.markAsRead(id);
    const notif = notifications.value.find(n => n.id === id);
    if (notif) notif.is_read = true;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to mark as read', life: 3000 });
  }
};

const markAllRead = async () => {
  try {
    await notificationService.markAllAsRead();
    notifications.value.forEach(n => n.is_read = true);
    toast.add({ severity: 'success', summary: 'Success', detail: 'All notifications marked as read', life: 2000 });
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to mark all as read', life: 3000 });
  }
};

const deleteNotification = async (id) => {
  try {
    await notificationService.delete(id);
    notifications.value = notifications.value.filter(n => n.id !== id);
    toast.add({ severity: 'info', summary: 'Deleted', detail: 'Notification removed', life: 2000 });
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.userMessage || 'Failed to delete', life: 3000 });
  }
};

onMounted(() => {
  loadNotifications();
});
</script>

<style scoped>
.notification-item:hover {
  transform: translateX(4px);
}
</style>

