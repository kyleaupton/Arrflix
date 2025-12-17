<script setup lang="ts">
import { ref } from 'vue'
import { useQuery, useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Message from 'primevue/message'
import { PrimeIcons } from '@/icons'
import {
  getV1DownloadersOptions,
  deleteV1DownloadersByIdMutation,
  postV1DownloadersByIdTestMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type DbgenDownloader } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import {
  downloaderColumns,
  createDownloaderActions,
} from '@/components/tables/configs/downloaderTableConfig'
import { useModal } from '@/composables/useModal'
import DownloaderUpsertDialog from './DownloaderUpsertDialog.vue'

// Data queries
const { data: downloaders, isLoading, refetch } = useQuery(getV1DownloadersOptions())
const modal = useModal()

// Mutations
const deleteDownloaderMutation = useMutation(deleteV1DownloadersByIdMutation())
const testDownloaderMutation = useMutation(postV1DownloadersByIdTestMutation())

// Error state for table-level errors
const downloaderError = ref<string | null>(null)
const testingId = ref<string | null>(null)

// Handlers
const handleAddDownloader = () => {
  modal.open(DownloaderUpsertDialog, {
    props: {
      header: 'Add Downloader',
      modal: true,
      style: { width: '600px' },
    },
    onClose: (result) => {
      const data = result?.data as { saved?: boolean } | undefined
      if (data?.saved) {
        refetch()
      }
    },
  })
}

const handleEditDownloader = (downloader: DbgenDownloader) => {
  modal.open(DownloaderUpsertDialog, {
    props: {
      header: 'Edit Downloader',
      modal: true,
      style: { width: '600px' },
    },
    data: {
      downloader,
    },
    onClose: (result) => {
      const data = result?.data as { saved?: boolean } | undefined
      if (data?.saved) {
        refetch()
      }
    },
  })
}

const handleDeleteDownloader = async (downloader: DbgenDownloader) => {
  if (!downloader.id) return
  const confirmed = await modal.confirm({
    title: 'Delete Downloader',
    message: `Are you sure you want to delete "${downloader.name}"?`,
    severity: 'danger',
  })
  if (!confirmed) return
  try {
    await deleteDownloaderMutation.mutateAsync({ path: { id: downloader.id } })
    refetch()
  } catch (err: unknown) {
    const error = err as { message?: string }
    downloaderError.value = error.message || 'Failed to delete downloader'
  }
}

const handleTestDownloader = async (downloader: DbgenDownloader) => {
  if (!downloader.id) return
  testingId.value = downloader.id
  try {
    const result = await testDownloaderMutation.mutateAsync({ path: { id: downloader.id } })
    if (result.success) {
      downloaderError.value = null
      await modal.alert({
        title: 'Connection Test Successful',
        message:
          result.message || `Connection test passed. Version: ${result.version || 'unknown'}`,
        severity: 'success',
      })
    } else {
      downloaderError.value = result.error || 'Connection test failed'
    }
  } catch (err: unknown) {
    const error = err as { message?: string }
    downloaderError.value = error.message || 'Connection test failed'
  } finally {
    testingId.value = null
  }
}

const downloaderActions = createDownloaderActions(
  handleTestDownloader,
  handleEditDownloader,
  handleDeleteDownloader,
)
</script>

<template>
  <div class="downloaders-settings">
    <div class="card">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold mb-2">Downloaders</h3>
            <p class="text-muted-color">
              Configure download clients for managing torrents and usenet downloads.
            </p>
          </div>
          <Button
            label="Add Downloader"
            :icon="PrimeIcons.PLUS"
            severity="primary"
            @click="handleAddDownloader"
          />
        </div>

        <Message v-if="downloaderError" severity="error" @close="downloaderError = null">{{
          downloaderError
        }}</Message>

        <DataTable
          :data="downloaders || []"
          :columns="downloaderColumns"
          :actions="downloaderActions"
          :loading="isLoading"
          empty-message="No downloaders configured"
          searchable
          search-placeholder="Search downloaders..."
          paginator
          :rows="10"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.downloaders-settings {
  max-width: 100%;
}
</style>
