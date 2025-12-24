<template>
  <BaseDialog title="Download Candidates">
    <DownloadCandidates :movie-id="movieId" @download-enqueued="handleDownloadEnqueued" />
  </BaseDialog>
</template>

<script setup lang="ts">
import BaseDialog from '@/components/modals/BaseDialog.vue'
import DownloadCandidates from './DownloadCandidates.vue'

defineProps<{
  movieId: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const handleDownloadEnqueued = () => {
  emit('close')
}
</script>

<style scoped>
.download-candidates-modal {
  min-height: 500px;
  max-height: 80vh;
}
</style>
