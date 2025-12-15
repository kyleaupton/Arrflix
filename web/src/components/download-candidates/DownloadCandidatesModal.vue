<template>
  <div class="download-candidates-modal h-full flex flex-col">
    <DownloadCandidates />
  </div>
</template>

<script setup lang="ts">
import DownloadCandidates from './DownloadCandidates.vue'
</script>

<style scoped>
.download-candidates-modal {
  min-height: 500px;
  max-height: 80vh;
}
</style>
