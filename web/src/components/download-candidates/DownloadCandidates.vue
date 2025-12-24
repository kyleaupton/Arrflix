<template>
  <div class="download-candidates-container flex flex-col h-full">
    <div class="flex-1 overflow-auto">
      <DownloadCandidateList v-if="!selectedCandidate" @enqueue="handlePreview" />
      <DownloadCandidatePreview v-else :candidate="selectedCandidate" />
    </div>

    <div v-if="selectedCandidate" class="flex flex-col">
      <Separator />

      <div class="flex justify-end gap-2">
        <Button variant="secondary" @click="handleCancel"> Cancel </Button>
        <Button :disabled="enqueueMutation.isPending.value" @click="handleEnqueue">
          {{ enqueueMutation.isPending.value ? 'Enqueuing...' : 'Enqueue' }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { postV1MovieByIdCandidateDownloadMutation } from '@/client/@tanstack/vue-query.gen'
import { type ModelDownloadCandidate } from '@/client/types.gen'
import DownloadCandidateList from './DownloadCandidatesList.vue'
import DownloadCandidatePreview from './DownloadCandidatePreview.vue'
import { useModal } from '@/composables/useModal'

const modal = useModal()

const props = defineProps<{
  movieId: number
}>()

const emit = defineEmits<{
  (e: 'download-enqueued'): void
}>()

const selectedCandidate = ref<ModelDownloadCandidate | null>(null)

const handlePreview = (candidate: ModelDownloadCandidate) => {
  selectedCandidate.value = candidate
}

const handleCancel = () => {
  selectedCandidate.value = null
}

// Enqueue mutation
const enqueueMutation = useMutation({
  ...postV1MovieByIdCandidateDownloadMutation(),
  onSuccess: () => {
    modal.alert({
      title: 'Download Enqueued',
      message: 'The download has been successfully enqueued.',
      severity: 'success',
    })
    emit('download-enqueued')
  },
  onError: (error) => {
    modal.alert({
      title: 'Enqueue Failed',
      message: error?.message || 'Failed to enqueue download candidate',
      severity: 'error',
    })
  },
})

const handleEnqueue = () => {
  if (!selectedCandidate.value) return

  enqueueMutation.mutate({
    path: { id: props.movieId },
    body: {
      indexerId: selectedCandidate.value.indexerId,
      guid: selectedCandidate.value.guid,
    },
  })
}
</script>

<style scoped>
.download-candidates-container {
  min-height: 500px;
}
</style>

<style scoped></style>
