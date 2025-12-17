<script setup lang="ts">
import { ref } from 'vue'
import { useQuery, useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Select from 'primevue/select'
import Password from 'primevue/password'
import Message from 'primevue/message'
import { PrimeIcons } from '@/icons'
import {
  getV1DownloadersOptions,
  postV1DownloadersMutation,
  putV1DownloadersByIdMutation,
  deleteV1DownloadersByIdMutation,
  postV1DownloadersByIdTestMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type DbgenDownloader } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import { downloaderColumns, createDownloaderActions } from '@/components/tables/configs/downloaderTableConfig'
import { useModal } from '@/composables/useModal'

// Data queries
const { data: downloaders, isLoading, refetch } = useQuery(getV1DownloadersOptions())
const modal = useModal()

// Mutations
const createDownloaderMutation = useMutation(postV1DownloadersMutation())
const updateDownloaderMutation = useMutation(putV1DownloadersByIdMutation())
const deleteDownloaderMutation = useMutation(deleteV1DownloadersByIdMutation())
const testDownloaderMutation = useMutation(postV1DownloadersByIdTestMutation())

// Modal state
const showDownloaderModal = ref(false)
const editingDownloader = ref<DbgenDownloader | null>(null)
const downloaderError = ref<string | null>(null)
const testingId = ref<string | null>(null)

// Downloader form
const downloaderForm = ref({
  name: '',
  type: 'qbittorrent',
  protocol: 'torrent' as 'torrent' | 'usenet',
  url: '',
  username: '',
  password: '',
  config_json: undefined as Record<string, any> | undefined,
  enabled: true,
  default: false,
})

// Handlers
const handleAddDownloader = () => {
  editingDownloader.value = null
  downloaderForm.value = {
    name: '',
    type: 'qbittorrent',
    protocol: 'torrent',
    url: '',
    username: '',
    password: '',
    config_json: undefined,
    enabled: true,
    default: false,
  }
  downloaderError.value = null
  showDownloaderModal.value = true
}

const handleEditDownloader = (downloader: DbgenDownloader) => {
  editingDownloader.value = downloader
  downloaderForm.value = {
    name: downloader.name || '',
    type: downloader.type || 'qbittorrent',
    protocol: (downloader.protocol as 'torrent' | 'usenet') || 'torrent',
    url: downloader.url || '',
    username: downloader.username || '',
    password: '', // Don't show existing password
    config_json: undefined,
    enabled: downloader.enabled ?? true,
    default: downloader.default || false,
  }
  downloaderError.value = null
  showDownloaderModal.value = true
}

const handleSaveDownloader = async () => {
  if (!downloaderForm.value.name || !downloaderForm.value.url) {
    downloaderError.value = 'Name and URL are required'
    return
  }

  try {
    if (editingDownloader.value?.id) {
      await updateDownloaderMutation.mutateAsync({
        path: { id: editingDownloader.value.id },
        body: {
          name: downloaderForm.value.name,
          type: downloaderForm.value.type,
          protocol: downloaderForm.value.protocol,
          url: downloaderForm.value.url,
          username: downloaderForm.value.username || undefined,
          password: downloaderForm.value.password || undefined,
          config_json: downloaderForm.value.config_json,
          enabled: downloaderForm.value.enabled,
          default: downloaderForm.value.default,
        },
      })
    } else {
      await createDownloaderMutation.mutateAsync({
        body: {
          name: downloaderForm.value.name,
          type: downloaderForm.value.type,
          protocol: downloaderForm.value.protocol,
          url: downloaderForm.value.url,
          username: downloaderForm.value.username || undefined,
          password: downloaderForm.value.password || undefined,
          config_json: downloaderForm.value.config_json,
          enabled: downloaderForm.value.enabled,
          default: downloaderForm.value.default,
        },
      })
    }
    showDownloaderModal.value = false
    refetch()
  } catch (err: any) {
    downloaderError.value = err.message || 'Failed to save downloader'
  }
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
  } catch (err: any) {
    downloaderError.value = err.message || 'Failed to delete downloader'
  }
}

const handleTestDownloader = async (downloader: DbgenDownloader) => {
  if (!downloader.id) return
  testingId.value = downloader.id
  try {
    await testDownloaderMutation.mutateAsync({ path: { id: downloader.id } })
    downloaderError.value = null
  } catch (err: any) {
    downloaderError.value = err.message || 'Connection test failed'
  } finally {
    testingId.value = null
  }
}

const handleTestDownloaderFromForm = async () => {
  // For new downloaders, we need to save first, then test
  if (!editingDownloader.value?.id) {
    downloaderError.value = 'Please save the downloader first before testing'
    return
  }
  
  testingId.value = editingDownloader.value.id
  downloaderError.value = null
  try {
    const result = await testDownloaderMutation.mutateAsync({ path: { id: editingDownloader.value.id } })
    if (result.success) {
      downloaderError.value = null
      // Show success message using alert
      await modal.alert({
        title: 'Connection Test Successful',
        message: result.message || `Connection test passed. Version: ${result.version || 'unknown'}`,
        severity: 'success',
      })
    } else {
      downloaderError.value = result.error || 'Connection test failed'
    }
  } catch (err: any) {
    downloaderError.value = err.message || 'Connection test failed'
  } finally {
    testingId.value = null
  }
}

const downloaderActions = createDownloaderActions(handleTestDownloader, handleEditDownloader, handleDeleteDownloader)
</script>

<template>
  <div class="downloaders-settings">
    <div class="card">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold mb-2">Downloaders</h3>
            <p class="text-muted-color">Configure download clients for managing torrents and usenet downloads.</p>
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

    <!-- Downloader Modal -->
    <Dialog
      v-model:visible="showDownloaderModal"
      :header="editingDownloader ? 'Edit Downloader' : 'Add Downloader'"
      :modal="true"
      :style="{ width: '600px' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Name</label>
          <InputText v-model="downloaderForm.name" placeholder="My qBittorrent" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Type</label>
          <Select
            v-model="downloaderForm.type"
            :options="[{ label: 'qBittorrent', value: 'qbittorrent' }]"
            optionLabel="label"
            optionValue="value"
            disabled
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Protocol</label>
          <Select
            v-model="downloaderForm.protocol"
            :options="[
              { label: 'Torrent', value: 'torrent' },
              { label: 'Usenet', value: 'usenet' },
            ]"
            optionLabel="label"
            optionValue="value"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">URL</label>
          <InputText v-model="downloaderForm.url" placeholder="http://localhost:8080" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Username (optional)</label>
          <InputText v-model="downloaderForm.username" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Password (optional)</label>
          <Password
            v-model="downloaderForm.password"
            :feedback="false"
            toggleMask
            :placeholder="editingDownloader ? 'Leave blank to keep current' : ''"
          />
        </div>
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium">Enabled</label>
          <ToggleSwitch v-model="downloaderForm.enabled" />
        </div>
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium">Default</label>
          <ToggleSwitch v-model="downloaderForm.default" />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-between w-full">
          <Button
            v-if="editingDownloader?.id"
            label="Test Connection"
            :icon="PrimeIcons.CHECK"
            severity="secondary"
            :loading="testingId === editingDownloader?.id"
            :disabled="!editingDownloader?.id || testDownloaderMutation.isPending.value"
            @click="handleTestDownloaderFromForm"
          />
          <div class="flex gap-2 ml-auto">
            <Button label="Cancel" severity="secondary" @click="showDownloaderModal = false" />
            <Button
              label="Save"
              :loading="createDownloaderMutation.isPending.value || updateDownloaderMutation.isPending.value"
              @click="handleSaveDownloader"
            />
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.downloaders-settings {
  max-width: 100%;
}
</style>
