<template>
  <BaseDialog title="Download Candidates">
    <DownloadCandidates
      :movie-id="movieId"
      :series-id="seriesId"
      :season="season"
      :episode="episode"
      @download-enqueued="handleDownloadEnqueued"
    />
  </BaseDialog>
</template>

<script setup lang="ts">
import BaseDialog from '@/components/modals/BaseDialog.vue'
import DownloadCandidates from './DownloadCandidates.vue'

defineProps<{
  movieId?: number
  seriesId?: number
  season?: number
  episode?: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const handleDownloadEnqueued = () => {
  emit('close')
}
</script>
