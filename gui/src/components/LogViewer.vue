<template>
  <div class="flex flex-col rounded-xl border border-gray-800 bg-gray-900 overflow-hidden h-full">
    <!-- Title Bar -->
    <div class="h-10 px-4 flex items-center border-b border-gray-800 bg-gray-900/80">
      <span class="text-sm font-semibold tracking-wide text-gray-200"> Log Viewer </span>
    </div>

    <!-- Content -->
    <div ref="logContainer" class="relative flex-1 overflow-auto font-mono text-xs" @scroll="handleScroll">
      <!-- Top spacer -->
      <div class="h-2"></div>

      <!-- Logs -->
      <div v-for="(line, index) in logs" :key="index" class="mb-2 pl-3 pr-2">
        <!-- Timestamp -->
        <span class="text-yellow-400 font-mono">
          {{ line.slice(0, 19) }}
        </span>

        <!-- Message -->
        <span class="ml-2 text-gray-100 font-mono">
          {{ line.slice(20) }}
        </span>
      </div>

      <!-- Bottom spacer -->
      <div class="h-2"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, nextTick } from 'vue';
  import axios from 'axios';

  const logs = ref<string[]>([]);
  const logContainer = ref<HTMLElement | null>(null);
  const isScrolled = ref(false);

  // Scroll logic
  const isAtBottom = () => {
    const el = logContainer.value;
    if (!el) return true;
    return el.scrollHeight - el.scrollTop - el.clientHeight < 20;
  };

  const scrollToBottom = () => {
    const el = logContainer.value;
    if (!el) return;

    requestAnimationFrame(() => {
      el.scrollTop = el.scrollHeight - el.clientHeight;
    });
  };

  const handleScroll = () => {
    const el = logContainer.value;
    if (!el) return;
    isScrolled.value = el.scrollTop > 5;
  };

  // Initial logs fetch
  const fetchInitialLogs = async () => {
    try {
      logs.value = (await axios.get('/api/logs')).data.split('\n').filter((line: string) => line.trim() !== '');

      nextTick(() => scrollToBottom());
    } catch (err) {
      console.error('Failed to fetch logs:', err);
    }
  };

  const setupWebSocket = () => {
    const ws = new WebSocket(`ws://${globalThis.location.host}/ws/logs`);

    ws.onmessage = (event) => {
      const message = event.data as string;

      const shouldScroll = isAtBottom();
      logs.value.push(message);

      if (shouldScroll) {
        nextTick(() => scrollToBottom());
      }
    };

    ws.onclose = () => {
      console.warn('Log WebSocket closed. Reconnecting in 2s...');
      setTimeout(setupWebSocket, 2000);
    };
  };

  onMounted(() => {
    fetchInitialLogs();
    setupWebSocket();
  });
</script>
