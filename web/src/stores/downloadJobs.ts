import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { DbgenDownloadJob } from '@/client/types.gen'
import { getV1DownloadJobs, deleteV1DownloadJobsById } from '@/client/sdk.gen'
import { useEventsStore } from '@/stores/events'

type DownloadJob = DbgenDownloadJob & {
  // Backend nullable fields can come through as null; tolerate at runtime
  progress?: number | null
  last_error?: string | null
  import_dest_path?: string | null
}

export const useDownloadJobsStore = defineStore('downloadJobs', () => {
  const jobsById = ref<Record<string, DownloadJob>>({})
  const isLoading = ref(false)

  const jobsSorted = computed(() => {
    return Object.values(jobsById.value).sort((a, b) => {
      return new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
    })
  })

  function upsert(job: DownloadJob) {
    jobsById.value = { ...jobsById.value, [job.id]: job }
  }

  function replaceAll(list: DownloadJob[]) {
    if (!list) return
    const next: Record<string, DownloadJob> = {}
    for (const j of list) next[j.id] = j
    jobsById.value = next
  }

  async function refresh() {
    isLoading.value = true
    try {
      const res = await getV1DownloadJobs({ throwOnError: true })
      replaceAll(res.data as unknown as DownloadJob[])
    } finally {
      isLoading.value = false
    }
  }

  function connectLive() {
    const events = useEventsStore()
    events.connect(['download_jobs_snapshot', 'download_job_updated', 'ping', 'ready'])

    events.on('download_jobs_snapshot', (data) => {
      replaceAll((data as unknown as DownloadJob[]) ?? [])
    })
    events.on('download_job_updated', (data) => {
      if (!data) return
      upsert(data as unknown as DownloadJob)
    })
  }

  async function cancelJob(id: string) {
    const res = await deleteV1DownloadJobsById({
      throwOnError: true,
      path: { id },
    })
    // Optimistically update local state; SSE may also deliver an update later.
    upsert(res.data as unknown as DownloadJob)
  }

  return {
    jobsById,
    jobsSorted,
    isLoading,
    refresh,
    connectLive,
    cancelJob,
  }
})
