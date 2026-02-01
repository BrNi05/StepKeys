<template>
  <div class="h-12 px-4 flex items-center justify-between border-b border-gray-800 bg-gray-900">
    <!-- Top bar, left side -->
    <div class="flex items-center gap-3">
      <img src="/icon.ico" class="h-6 w-6" alt="StepKeys Icon" />
      <span class="font-semibold tracking-wide">StepKeys</span>
    </div>

    <!-- Top bar, right side -->
    <div class="flex items-center gap-4">
      <!-- Enabled Toggle -->
      <div class="flex items-center gap-2">
        <span class="text-sm opacity-70">Enabled</span>

        <ToggleSwitch
          :model-value="enabled"
          active-color="bg-green-500/40"
          inactive-color="bg-gray-700"
          @update:model-value="toggleEnabledState"
        />
      </div>

      <!-- Start on boot Toggle -->
      <div class="flex items-center gap-2">
        <span class="text-sm opacity-70">Start on boot</span>

        <ToggleSwitch
          :model-value="boot"
          active-color="bg-green-500/40"
          inactive-color="bg-gray-700"
          @update:model-value="toggleBootState"
        />
      </div>

      <!-- GitHub icon -->
      <a href="https://github.com/BrNi05/StepKeys" target="_blank" class="opacity-100">
        <img src="/github.png" alt="GitHub" class="h-6 w-6" />
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import { getEnabled, toggleEnabled, getBoot, toggleBoot } from '@/api/client';
  import ToggleSwitch from '@/components/generic/ToggleSwitch.vue';

  const enabled = ref(false);
  const boot = ref(false);

  const toggleEnabledState = async () => {
    enabled.value = (await toggleEnabled()).data.value; // silent failure can happen
  };

  const toggleBootState = async () => {
    boot.value = (await toggleBoot()).data.value; // silent failure can happen
  };

  const setupWebSocket = () => {
    const ws = new WebSocket(`ws://${globalThis.location.host}/ws/settings`);

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.event === 'enabled') enabled.value = data.value;
        if (data.event === 'boot') boot.value = data.value;
      } catch (err) {
        console.error('Failed to parse WS message:', err);
      }
    };

    ws.onclose = () => {
      console.warn('Settings WebSocket closed, reconnecting in 2s...');
      setTimeout(setupWebSocket, 2000);
    };
  };

  onMounted(async () => {
    enabled.value = (await getEnabled()).data.value;
    boot.value = (await getBoot()).data.value;

    setupWebSocket(); // keep GUI in sync
  });
</script>
