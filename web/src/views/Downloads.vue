<script setup lang="ts">
import { onMounted } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import { useEventsStore } from '@/stores/events'
import { useDownloadJobsStore } from '@/stores/downloadJobs'
import DataTable from '@/components/tables/DataTable.vue'
import {
  downloadJobColumns,
  createDownloadJobActions,
} from '@/components/tables/configs/downloadJobsTableConfig'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'

const events = useEventsStore()
const jobs = useDownloadJobsStore()

const handleCancelJob = (job: ReturnType<typeof useDownloadJobsStore>['jobsSorted'][number]) => {
  jobs.cancelJob(job.id)
}

const jobActions = createDownloadJobActions(handleCancelJob)

onMounted(async () => {
  jobs.connectLive()
  await jobs.refresh()
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold">Downloads</h1>
    </div>
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <div class="text-sm text-muted-foreground">
          <span class="font-semibold">Live:</span>
          <span class="ml-2">{{ events.status }}</span>
          <span v-if="events.lastError" class="ml-3 text-destructive">{{ events.lastError }}</span>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="outline" @click="jobs.refresh()">
            <RefreshCw class="mr-2 size-4" />
            Refresh
          </Button>
        </div>
      </div>

      <div v-if="jobs.isLoading" class="space-y-3">
        <Skeleton class="h-12 w-full" />
        <Skeleton class="h-12 w-full" />
        <Skeleton class="h-12 w-full" />
      </div>
      <DataTable
        v-else
        :data="jobs.jobsSorted"
        :columns="downloadJobColumns"
        :actions="jobActions"
        :loading="jobs.isLoading"
        empty-message="No download jobs"
        searchable
        search-placeholder="Search downloads..."
        paginator
        :rows="10"
      />
    </div>
  </div>
</template>
