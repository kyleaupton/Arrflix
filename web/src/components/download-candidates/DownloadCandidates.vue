<template>
  <div class="download-candidates-container flex flex-col h-full">
    <div class="flex-1 overflow-auto">
      <DownloadCandidateList v-if="!selectedCandidate" @enqueue="handlePreview" />
      <DownloadCandidatePreview v-else :candidate="selectedCandidate" />
    </div>

    <div v-if="selectedCandidate" class="flex flex-col">
      <Divider />

      <div class="flex justify-end gap-2">
        <Button label="Cancel" severity="secondary" @click="handleCancel" />
        <Button label="Enqueue" :loading="enqueueMutation.isPending.value" @click="handleEnqueue" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Divider from 'primevue/divider'
import { postV1MovieByIdEnqueueCandidateMutation } from '@/client/@tanstack/vue-query.gen'
import { type ModelDownloadCandidate } from '@/client/types.gen'
import DownloadCandidateList from './DownloadCandidatesList.vue'
import DownloadCandidatePreview from './DownloadCandidatePreview.vue'
import { useModal } from '@/composables/useModal'

const modal = useModal()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const dialogRef = inject('dialogRef') as any

const movieId = computed(() => {
  const id = dialogRef.value?.data?.movieId
  if (!id) {
    throw new Error('Movie ID is required')
  }
  return id
})

const selectedCandidate = ref<ModelDownloadCandidate | null>(null)

const handlePreview = (candidate: ModelDownloadCandidate) => {
  selectedCandidate.value = candidate
}

const handleCancel = () => {
  selectedCandidate.value = null
}

// Enqueue mutation
const enqueueMutation = useMutation({
  ...postV1MovieByIdEnqueueCandidateMutation(),
  onSuccess: () => {
    modal.alert({
      title: 'Download Enqueued',
      message: 'The download has been successfully enqueued.',
      severity: 'success',
    })
    dialogRef.value?.close()
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
    path: { id: movieId.value },
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
