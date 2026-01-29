<template>
  <div class="h-full w-full flex flex-col bg-gray-950 text-gray-100 overflow-hidden">
    <TopBar />

    <div class="flex-1 flex gap-3 p-3 overflow-hidden">
      <PedalEditor />

      <div class="w-[32%]">
        <LogViewer />
      </div>
    </div>

    <BottomBar />
  </div>
</template>

<script setup lang="ts">
  import TopBar from '@/components/TopBar.vue';
  import BottomBar from '@/components/BottomBar.vue';
  import PedalEditor from '@/components/PedalEditor.vue';
  import LogViewer from '@/components/LogViewer.vue';
</script>
