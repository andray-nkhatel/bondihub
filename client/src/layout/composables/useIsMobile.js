import { onMounted, onUnmounted, ref } from 'vue';

export function useIsMobile(breakpoint = 768) {
  const isMobile = ref(window.innerWidth < breakpoint);

  function handleResize() {
    isMobile.value = window.innerWidth < breakpoint;
  }

  onMounted(() => {
    window.addEventListener('resize', handleResize);
    handleResize();
  });

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
  });

  return isMobile;
}