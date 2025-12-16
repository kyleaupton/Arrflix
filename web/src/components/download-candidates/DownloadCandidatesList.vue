<template>
  <div class="download-candidates">
    <div class="mb-4">
      <!-- <h3 class="text-lg font-semibold mb-2">Download Candidates</h3> -->
      <!-- <p class="text-sm text-muted-color">
        Search results for download candidates. Select a candidate to enqueue it for download.
      </p> -->
    </div>

    <DataTable
      :query-options="queryOptions"
      :columns="downloadCandidateColumns"
      :actions="candidateActions"
      searchable
      search-placeholder="Search candidates..."
      paginator
      :rows="20"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { getV1MovieByIdCandidatesOptions } from '@/client/@tanstack/vue-query.gen'
import { type ModelDownloadCandidate } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import {
  downloadCandidateColumns,
  createDownloadCandidateActions,
} from '@/components/tables/configs/downloadCandidateTableConfig'

const emit = defineEmits<{
  (e: 'enqueue', candidate: ModelDownloadCandidate): void
}>()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const dialogRef = inject('dialogRef') as any

const movieId = computed(() => {
  const id = dialogRef.value?.data?.movieId
  if (!id) {
    throw new Error('Movie ID is required')
  }
  return id
})

// Query options for fetching download candidates
const queryOptions = computed(() =>
  getV1MovieByIdCandidatesOptions({
    path: { id: movieId.value },
  }),
)

// Handle enqueue action
const handleEnqueue = (candidate: ModelDownloadCandidate) => {
  console.log('enqueue', candidate)
  emit('enqueue', candidate)
}

// Create actions
const candidateActions = createDownloadCandidateActions(handleEnqueue)
</script>

<style scoped>
.download-candidates {
  width: 100%;
  height: 100%;
}
</style>
