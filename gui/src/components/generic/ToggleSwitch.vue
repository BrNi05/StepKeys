<template>
  <div
    class="w-11 h-6 rounded-full cursor-pointer transition-colors duration-200 flex items-center px-0.5"
    :class="modelValue ? activeColor : inactiveColor"
    @click="emit('update:modelValue', !modelValue)"
  >
    <div
      class="w-5 h-5 rounded-full bg-white transition-transform duration-200 shadow"
      :class="modelValue ? 'translate-x-6' : 'translate-x-0'"
    ></div>
  </div>
</template>

<script setup lang="ts">
  defineProps({
    modelValue: {
      type: Boolean,
      required: true,
    },
    activeColor: {
      type: String,
      default: 'bg-green-600/40',
    },
    inactiveColor: {
      type: String,
      default: 'bg-gray-700',
    },
  });

  const emit = defineEmits(['update:modelValue']);
</script>
