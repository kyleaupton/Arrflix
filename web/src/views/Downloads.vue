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
          <Button label="Refresh" severity="secondary" @click="jobs.refresh()" />
        </div>
      </div>

      <DataTable :value="jobs.jobsSorted" :loading="jobs.isLoading" dataKey="id" class="text-sm">
        <Column header="Title">
          <template #body="{ data }">
            <div class="font-medium">{{ data.candidate_title }}</div>
            <div class="text-xs text-muted-foreground">{{ data.protocol }} • {{ data.id }}</div>
          </template>
        </Column>

        <Column header="Status" style="width: 140px">
          <template #body="{ data }">
            <Tag :value="data.status" />
          </template>
        </Column>

        <Column header="Progress" style="width: 220px">
          <template #body="{ data }">
            <ProgressBar
              :value="Math.round(((data.progress ?? 0) as number) * 100)"
              :showValue="true"
            />
          </template>
        </Column>

        <Column header="Imported To">
          <template #body="{ data }">
            <span class="text-xs">{{ data.import_dest_path || '-' }}</span>
          </template>
        </Column>

        <Column header="Error">
          <template #body="{ data }">
            <span class="text-xs text-destructive">{{ data.last_error || '' }}</span>
          </template>
        </Column>

        <Column header="" style="width: 120px">
          <template #body="{ data }">
            <Button
              label="Cancel"
              size="small"
              severity="danger"
              :disabled="
                data.status === 'cancelled' ||
                data.status === 'imported' ||
                data.status === 'failed'
              "
              @click="jobs.cancelJob(data.id)"
            />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ProgressBar from 'primevue/progressbar'
import { useEventsStore } from '@/stores/events'
import { useDownloadJobsStore } from '@/stores/downloadJobs'

const events = useEventsStore()
const jobs = useDownloadJobsStore()

onMounted(async () => {
  jobs.connectLive()
  await jobs.refresh()
})
</script>

<style scoped></style>
