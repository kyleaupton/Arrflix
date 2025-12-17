<script setup lang="ts">
import { ref, inject, computed } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Select from 'primevue/select'
import Password from 'primevue/password'
import { PrimeIcons } from '@/icons'
import {
  postV1DownloadersMutation,
  putV1DownloadersByIdMutation,
  postV1DownloadersByIdTestMutation,
} from '@/client/@tanstack/vue-query.gen'
import { useModal } from '@/composables/useModal'

/* eslint-disable-next-line @typescript-eslint/no-explicit-any */
const dialogRef = inject('dialogRef') as any
const modal = useModal()

const data = computed(() => dialogRef.value?.data || {})

// Mutations
const createDownloaderMutation = useMutation(postV1DownloadersMutation())
const updateDownloaderMutation = useMutation(putV1DownloadersByIdMutation())
const testDownloaderMutation = useMutation(postV1DownloadersByIdTestMutation())

// Form state
const downloaderForm = ref({
  name: data.value.downloader?.name || '',
  type: data.value.downloader?.type || 'qbittorrent',
  protocol: (data.value.downloader?.protocol as 'torrent' | 'usenet') || 'torrent',
  url: data.value.downloader?.url || '',
  username: data.value.downloader?.username || '',
  password: data.value.downloader?.password || '',
  config_json: {} as Record<string, unknown>,
  enabled: data.value.downloader?.enabled ?? true,
  default: data.value.downloader?.default || false,
})

const downloaderError = ref<string | null>(null)
const testingId = ref<string | null>(null)

// Handlers
const handleSaveDownloader = async () => {
  if (!downloaderForm.value.name || !downloaderForm.value.url) {
    downloaderError.value = 'Name and URL are required'
    return
  }

  try {
    if (data.value.downloader?.id) {
      await updateDownloaderMutation.mutateAsync({
        path: { id: data.value.downloader.id },
        body: {
          name: downloaderForm.value.name,
          type: downloaderForm.value.type,
          protocol: downloaderForm.value.protocol,
          url: downloaderForm.value.url,
          username: downloaderForm.value.username || '',
          password: downloaderForm.value.password || '',
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
          username: downloaderForm.value.username || '',
          password: downloaderForm.value.password || '',
          config_json: downloaderForm.value.config_json,
          enabled: downloaderForm.value.enabled,
          default: downloaderForm.value.default,
        },
      })
    }
    dialogRef.value.close({ saved: true })
  } catch (err: unknown) {
    const error = err as { message?: string }
    downloaderError.value = error.message || 'Failed to save downloader'
  }
}

const handleTestDownloader = async () => {
  if (!data.value.downloader?.id) {
    downloaderError.value = 'Please save the downloader first before testing'
    return
  }

  testingId.value = data.value.downloader.id
  downloaderError.value = null
  try {
    const result = await testDownloaderMutation.mutateAsync({
      path: { id: data.value.downloader.id },
    })
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

const handleCancel = () => {
  dialogRef.value.close()
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div
      v-if="downloaderError"
      class="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded text-red-700 dark:text-red-300 text-sm"
    >
      {{ downloaderError }}
    </div>

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
        :placeholder="data.downloader ? 'Leave blank to keep current' : ''"
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

    <div class="flex items-center justify-between w-full pt-2">
      <Button
        v-if="data.downloader?.id"
        label="Test Connection"
        :icon="PrimeIcons.CHECK"
        severity="secondary"
        :loading="testingId === data.downloader?.id"
        :disabled="!data.downloader?.id || testDownloaderMutation.isPending.value"
        @click="handleTestDownloader"
      />
      <div class="flex gap-2 ml-auto">
        <Button label="Cancel" severity="secondary" @click="handleCancel" />
        <Button
          label="Save"
          :loading="
            createDownloaderMutation.isPending.value || updateDownloaderMutation.isPending.value
          "
          @click="handleSaveDownloader"
        />
      </div>
    </div>
  </div>
</template>

<style scoped></style>
