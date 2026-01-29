<template>
  <div class="h-11 px-4 flex items-center justify-between border-t border-gray-800 bg-gray-900 text-sm">
    <!-- Bottom bar, left side -->
    <div class="opacity-70">
      <span class="font-bold">Serial port:</span> <span class="ml-1">{{ serial }}</span>
    </div>

    <!-- Bottom bar, right side -->
    <div class="flex items-center gap-3">
      <!-- Update Button -->
      <button
        class="px-3 py-1 rounded border transition-all duration-200"
        :class="
          updateAvailable
            ? 'bg-yellow-500/20 border-yellow-500 text-yellow-400 animate-pulse hover:bg-yellow-500/30'
            : 'bg-gray-800 border-gray-700 text-gray-400 hover:bg-gray-700'
        "
        :style="updateAvailable ? 'animation-duration: 3s;' : ''"
        @click="handleUpdateAction"
      >
        {{ updateAvailable ? 'Update available' : 'Check updates' }}
      </button>

      <!-- Quit -->
      <button
        class="px-2 py-1 rounded border transition-colors duration-200 appearance-none outline-none bg-red-900/10 border-red-500/50 text-red-400 hover:bg-red-900/20 hover:border-red-400 hover:text-red-300"
        @click="quit"
      >
        ⏻
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import { getSerial, getUpdate, quitApp } from '@/api/client';

  const serial = ref('');
  const updateAvailable = ref(false);

  const handleUpdateAction = async () => {
    if (updateAvailable.value) {
      // Update available
      window.open('https://github.com/BrNi05/StepKeys/releases', '_blank');
    } else {
      // No update on startup, check again
      updateAvailable.value = (await getUpdate(true)).data.value;
    }
  };

  const quit = async () => {
    await quitApp();
    window.close();
  };

  // Set initial values
  onMounted(async () => {
    serial.value = (await getSerial()).data.value;
    updateAvailable.value = (await getUpdate(false)).data.value;
  });
</script>
