<template>
  <div class="p-3 border border-gray-700 rounded-lg flex flex-col gap-2 bg-gray-800">
    <div class="flex items-center gap-2">
      <span class="font-semibold text-yellow-400">Pedal {{ id }}</span>
      <button
        class="ml-auto px-2 py-1 rounded border transition-colors duration-200 bg-red-900/10 border-red-500/50 text-red-400 hover:bg-red-900/20 hover:border-red-400 hover:text-red-300 text-xs"
        @click="$emit('remove')"
      >
        Remove
      </button>
    </div>

    <div class="flex gap-3">
      <div class="flex flex-col text-xs text-gray-400">
        <span class="mb-1">Mode</span>
        <select
          v-model="model.mode"
          class="bg-gray-700 rounded px-2 py-1 text-sm text-gray-100 outline-none border border-transparent focus:border-gray-600 cursor-pointer"
        >
          <option value="sequence">Sequence</option>
          <option value="combo">Combo</option>
        </select>
      </div>

      <div class="flex flex-col text-xs text-gray-400">
        <span class="mb-1">Behaviour</span>
        <select
          v-model="model.behaviour"
          class="bg-gray-700 rounded px-2 py-1 text-sm text-gray-100 outline-none border border-transparent focus:border-gray-600 cursor-pointer"
        >
          <option value="oneshot">Oneshot</option>
          <option value="toggle">Toggle</option>
          <option value="hold">Hold</option>
        </select>
      </div>
    </div>

    <div class="flex flex-col gap-1 text-xs text-gray-400">
      <span class="select-none">Keys</span>
      <div class="flex flex-wrap gap-1 items-center bg-gray-900/20 p-1 rounded-md border border-gray-700">
        <span
          v-for="(key, idx) in model.keys"
          :key="`${id}-${key}-${idx}`"
          class="flex items-center gap-1 px-2 py-1 bg-gray-700 rounded text-sm text-gray-100 font-medium border border-gray-600 transition-colors"
        >
          {{ key }}
          <button
            type="button"
            class="w-5 h-5 flex items-center justify-center rounded-full text-red-500 bg-red-500/10 hover:bg-red-500/20 transition-colors appearance-none outline-none leading-none text-lg pb-0.5"
            @click.stop="removeKey(idx)"
          >
            ×
          </button>
        </span>

        <input
          v-model="query"
          class="bg-transparent px-2 py-1 text-sm flex-1 min-w-[100px] text-white placeholder-gray-500 outline-none"
          placeholder="Type key..."
          @keydown.enter.prevent="handleEnterOrSpace"
          @keydown.space.prevent="handleEnterOrSpace"
          @keydown.backspace="handleBackspace"
          @keydown.down.prevent="moveHighlight(1)"
          @keydown.up.prevent="moveHighlight(-1)"
          @keydown.esc="query = ''"
        />
      </div>
    </div>

    <div
      v-if="suggestions.length"
      ref="scrollContainer"
      class="bg-gray-700 mt-1 rounded shadow-xl border border-gray-600 max-h-32 overflow-auto z-10"
    >
      <div
        v-for="(s, i) in suggestions"
        :key="s"
        class="px-2 py-1.5 cursor-pointer text-sm border-b border-gray-600/50 last:border-0 transition-colors"
        :class="highlightedIdx === i ? 'bg-blue-600 text-white' : 'hover:bg-gray-600 text-white'"
        @mouseenter="highlightedIdx = i"
        @click="addKey(s)"
      >
        {{ s }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, watch, nextTick } from 'vue';
  import type { PedalAction } from '@/interfaces/pedal';

  const props = defineProps<{ id: string; validKeys: string[] }>();
  const model = defineModel<PedalAction>({ required: true });
  defineEmits(['remove']); // emitted when the pedal entry remove button is clicked

  const query = ref('');
  const highlightedIdx = ref(0);
  const scrollContainer = ref<HTMLElement | null>(null);

  // Compute suggestions based on current query and existing keys
  // Do not suggest (and allow) keys that are already assigned to this pedal
  const suggestions = computed(() => {
    const input = query.value.trim(); // original, case-sensitive input
    if (!input) return [];

    const lowerInput = input.toLowerCase();

    return props.validKeys
      .filter((k) => k.toLowerCase().startsWith(lowerInput) && !model.value.keys.includes(k))
      .sort((a, b) => {
        // Case-sensitive sorting matches
        const aStartsExact = a.startsWith(input);
        const bStartsExact = b.startsWith(input);

        if (aStartsExact && !bStartsExact) return -1;
        if (!aStartsExact && bStartsExact) return 1;

        // Sort by length
        if (a.length !== b.length) {
          return a.length - b.length;
        }

        // Sort by alphabetical order
        return a.localeCompare(b);
      });
  });

  // Reset highlighting when the suggestion list changes
  watch(suggestions, () => {
    highlightedIdx.value = 0;
  });

  const moveHighlight = (delta: number) => {
    const len = suggestions.value.length;
    if (len === 0) return;
    highlightedIdx.value = (highlightedIdx.value + delta + len) % len;

    // Auto-scroll logic
    nextTick(() => {
      const el = scrollContainer.value?.children[highlightedIdx.value] as HTMLElement;
      if (el) el.scrollIntoView({ block: 'nearest' });
    });
  };

  const handleEnterOrSpace = () => {
    // If a suggestion is highlighted and valid, use it
    if (suggestions.value.length > 0 && highlightedIdx.value >= 0) {
      addKey(suggestions.value[highlightedIdx.value]!);
    } else {
      // Otherwise try to add the raw input
      addKey(query.value);
    }
  };

  // Delete last key if query is empty
  const handleBackspace = () => {
    if (query.value.length === 0 && model.value.keys.length > 0) {
      removeKey(model.value.keys.length - 1);
    }
  };

  // Add (keyboard) key by assigning a new array with the new key appended
  const addKey = (keyName: string) => {
    const key = keyName.trim();
    if (!key || !props.validKeys.includes(key) || model.value.keys.includes(key)) return;

    model.value = { ...model.value, keys: [...model.value.keys, key] };
    query.value = '';
  };

  // Remove (keyboard) key by assigning a new array without the key at idx
  const removeKey = (idx: number) => {
    model.value = { ...model.value, keys: model.value.keys.filter((_, i) => i !== idx) };
  };
</script>
