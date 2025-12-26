<template>
  <div class="download-candidates-container flex flex-col h-full">
    <div class="flex-1 overflow-auto">
      <DownloadCandidateList
        v-if="!selectedCandidate"
        :movie-id="movieId"
        :series-id="seriesId"
        :season="season"
        :episode="episode"
        @enqueue="handlePreview"
      />
      <DownloadCandidatePreview
        v-else
        :movie-id="movieId"
        :series-id="seriesId"
        :season="season"
        :episode="episode"
        :candidate="selectedCandidate"
      />
    </div>

    <div v-if="selectedCandidate" class="flex flex-col">
      <Separator class="my-4" />

      <div class="flex justify-end gap-2">
        <Button variant="secondary" @click="handleCancel"> Cancel </Button>
        <Button
          :disabled="enqueueMovieMutation.isPending.value || enqueueSeriesMutation.isPending.value"
          @click="handleEnqueue"
        >
          {{
            enqueueMovieMutation.isPending.value || enqueueSeriesMutation.isPending.value
              ? 'Enqueuing...'
              : 'Enqueue'
          }}
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
import {
  postV1MovieByIdCandidateDownloadMutation,
  postV1SeriesByIdCandidateDownloadMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type ModelDownloadCandidate } from '@/client/types.gen'
import DownloadCandidateList from './DownloadCandidatesList.vue'
import DownloadCandidatePreview from './DownloadCandidatePreview.vue'
import { useModal } from '@/composables/useModal'

const modal = useModal()

const props = defineProps<{
  movieId?: number
  seriesId?: number
  season?: number
  episode?: number
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

// Enqueue movie mutation
const enqueueMovieMutation = useMutation({
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

// Enqueue series mutation
const enqueueSeriesMutation = useMutation({
  ...postV1SeriesByIdCandidateDownloadMutation(),
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

  if (props.movieId) {
    enqueueMovieMutation.mutate({
      path: { id: props.movieId },
      body: {
        indexerId: selectedCandidate.value.indexerId,
        guid: selectedCandidate.value.guid,
      },
    })
  } else if (props.seriesId) {
    enqueueSeriesMutation.mutate({
      path: { id: props.seriesId },
      body: {
        indexerId: selectedCandidate.value.indexerId,
        guid: selectedCandidate.value.guid,
        season: props.season,
        episode: props.episode,
      },
    })
  }
}
</script>

<style scoped>
.download-candidates-container {
  min-height: 500px;
}
</style>

<style scoped></style>
