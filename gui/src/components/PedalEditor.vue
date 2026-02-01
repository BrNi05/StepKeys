<template>
  <div class="flex-1 flex flex-col rounded-xl border border-gray-800 bg-gray-900 overflow-hidden">
    <!-- Pedal Entries -->
    <div class="flex-1 overflow-auto p-4 space-y-4">
      <PedalEntry
        v-for="(_, pedalID) in pedals"
        :id="String(pedalID)"
        :key="pedalID"
        v-model="pedals[pedalID]!"
        :valid-keys="validKeys"
        @remove="removePedal(String(pedalID))"
      />
    </div>

    <!-- Pedal Editor menu bar -->
    <div class="shrink-0 flex items-center gap-3 p-3 border-t border-gray-800 bg-gray-900/90">
      <button
        class="px-4 py-2 rounded-lg text-red-400 border border-red-500/50 bg-red-900/10 hover:bg-red-900/20 hover:border-red-400 hover:text-red-300 font-medium transition-colors duration-200"
        @click="loadPedals"
      >
        Reset
      </button>

      <div class="ml-auto"></div>

      <div class="ml-3 flex gap-2 items-center">
        <input
          v-model="profileName"
          type="text"
          placeholder="Profile name"
          class="px-3 py-2 rounded-lg bg-gray-800 border border-gray-700 text-gray-100 focus:outline-none focus:ring-2 focus:ring-yellow-400"
        />

        <button
          class="px-4 py-2 rounded-lg text-yellow-400 border border-yellow-500/50 bg-yellow-900/10 hover:bg-yellow-400/10 hover:border-yellow-400 hover:text-yellow-300 font-medium transition-colors duration-200"
          @click="saveConfig"
        >
          Save Config
        </button>

        <button
          class="px-4 py-2 rounded-lg text-yellow-400 border border-yellow-500/50 bg-yellow-900/10 hover:bg-yellow-400/10 hover:border-yellow-400 hover:text-yellow-300 font-medium transition-colors duration-200"
          @click="loadConfigFromFile"
        >
          Load Config
        </button>
      </div>

      <button
        class="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-500 text-white font-medium transition-colors"
        @click="addPedal"
      >
        + Add Pedal
      </button>

      <button
        class="px-4 py-2 rounded-lg bg-green-600 hover:bg-green-500 text-white font-medium transition-colors"
        @click="applyPedals"
      >
        Apply
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import { getPedals, setPedals, getValidKeys } from '@/api/client';
  import PedalEntry from './generic/PedalEntry.vue';
  import type { PedalAction } from '@/interfaces/pedal';

  // Current pedals map
  const pedals = ref<Record<string, PedalAction>>({});

  // Fetched list of valid keyboard keys
  const validKeys = ref<string[]>([]);

  // Profile name for saving/loading configs
  const profileName = ref('');

  // Load pedals from backend and overwrite current state
  const loadPedals = async () => {
    try {
      const res = await getPedals();
      pedals.value = res.data;
    } catch (err) {
      console.error('Failed to load pedals:', err);
    }
  };

  // Send current pedal config to the backend
  const applyPedals = async () => {
    const { ok, message } = await setPedals(pedals.value);
    if (!ok) console.error(`Failed to save: ${message}`);
  };

  // Add a new pedal with the lowest available ID
  const addPedal = () => {
    globalThis.getSelection()?.removeAllRanges();

    const used = Object.keys(pedals.value)
      .map(Number)
      .filter(Number.isInteger)
      .sort((a, b) => a - b);

    let nextID = 0;
    for (const id of used) {
      if (id === nextID) nextID++;
      else break;
    }

    // Set a logical default
    pedals.value[String(nextID)] = {
      mode: 'combo',
      behaviour: 'oneshot',
      keys: [],
    };
  };

  const removePedal = (id: string) => {
    const newMap = { ...pedals.value };
    delete newMap[id];
    pedals.value = newMap;
  };

  const setupWebSocket = () => {
    const ws = new WebSocket(`ws://${location.host}/ws/pedals`);

    // Trigger an API call on WS ping
    ws.onmessage = () => loadPedals();

    ws.onclose = () => {
      console.warn('Settings WebSocket closed, reconnecting in 2s...');
      setTimeout(setupWebSocket, 2000);
    };
  };

  onMounted(async () => {
    const [keysRes] = await Promise.all([getValidKeys(), loadPedals()]);
    validKeys.value = keysRes.data;
    setupWebSocket();
  });

  // Save current pedal config to a JSON file
  const saveConfig = () => {
    const dataToSave = {
      profileName: profileName.value,
      pedals: pedals.value,
    };
    const json = JSON.stringify(dataToSave, null, 2);

    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);

    const a = document.createElement('a');
    a.href = url;
    a.download = profileName.value ? `${profileName.value}.json` : 'pedals.json';
    a.click();

    URL.revokeObjectURL(url);
  };

  // Load pedal config from a JSON file
  const loadConfigFromFile = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'application/json';

    input.onchange = async () => {
      if (!input.files?.length || !input.files[0]) {
        console.error('No file selected');
        return;
      }

      try {
        const file = input.files[0];
        const text = await file.text();
        pedals.value = JSON.parse(text);

        profileName.value = file.name.replace(/\.[^/.]+$/, '');
      } catch {
        console.error('Invalid config file. Loading failed.');
      }
    };

    input.click();
  };
</script>
