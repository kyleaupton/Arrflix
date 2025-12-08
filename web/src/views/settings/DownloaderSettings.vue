<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Button from 'primevue/button'
import Message from 'primevue/message'
import Skeleton from 'primevue/skeleton'
import Select from 'primevue/select'
import Password from 'primevue/password'
import {
  getV1Downloaders,
  postV1Downloaders,
  putV1DownloadersById,
  deleteV1DownloadersById,
  postV1DownloadersByIdTest,
} from '@/client/sdk.gen'

type Downloader = {
  id?: string
  name?: string
  type?: 'qbittorrent' | string
  protocol?: 'torrent' | 'usenet' | string
  url?: string
  username?: string | null
  password?: string | null
  config_json?: Record<string, any> | null
  enabled?: boolean
  default?: boolean
}

const downloaders = ref<Downloader[]>([])
const downloadersLoading = ref(false)
const downloadersError = ref<string | null>(null)

async function loadDownloaders() {
  downloadersLoading.value = true
  downloadersError.value = null
  try {
    const res = await getV1Downloaders<true>({ throwOnError: true })
    downloaders.value = (res.data as Downloader[]) ?? []
  } catch {
    downloadersError.value = 'Failed to load downloaders'
  } finally {
    downloadersLoading.value = false
  }
}

onMounted(loadDownloaders)

// Create form
const newDownloader = ref<Required<Omit<Downloader, 'id' | 'config_json' | 'username' | 'password'>> & {
  username?: string
  password?: string
  config_json?: Record<string, any>
}>({
  name: '',
  type: 'qbittorrent',
  protocol: 'torrent',
  url: '',
  username: '',
  password: '',
  config_json: undefined,
  enabled: true,
  default: false,
})
const isCreating = ref(false)
async function createDownloader() {
  if (!newDownloader.value.name || !newDownloader.value.url) return
  isCreating.value = true
  try {
    const res = await postV1Downloaders<true>({
      throwOnError: true,
      body: {
        name: newDownloader.value.name,
        type: newDownloader.value.type,
        protocol: newDownloader.value.protocol,
        url: newDownloader.value.url,
        username: newDownloader.value.username || undefined,
        password: newDownloader.value.password || undefined,
        config_json: newDownloader.value.config_json,
        enabled: newDownloader.value.enabled,
        default: newDownloader.value.default,
      },
    })
    downloaders.value = [...downloaders.value, res.data as Downloader]
    newDownloader.value = {
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
    await loadDownloaders() // Reload to refresh default flags
  } catch (err: any) {
    downloadersError.value = err.message || 'Failed to create downloader'
  } finally {
    isCreating.value = false
  }
}

// Edit helpers
const editingId = ref<string | null>(null)
const editBuf = ref<Required<Omit<Downloader, 'id' | 'config_json' | 'username' | 'password'>> & {
  username?: string
  password?: string
  config_json?: Record<string, any>
}>({
  name: '',
  type: 'qbittorrent',
  protocol: 'torrent',
  url: '',
  username: '',
  password: '',
  config_json: undefined,
  enabled: true,
  default: false,
})

function startEdit(downloader: Downloader) {
  editingId.value = downloader.id ?? null
  editBuf.value = {
    name: downloader.name ?? '',
    type: (downloader.type as 'qbittorrent') ?? 'qbittorrent',
    protocol: (downloader.protocol as 'torrent' | 'usenet') ?? 'torrent',
    url: downloader.url ?? '',
    username: downloader.username ?? '',
    password: '', // Don't show existing password
    config_json: downloader.config_json ?? undefined,
    enabled: downloader.enabled ?? true,
    default: downloader.default ?? false,
  }
}

async function saveEdit(id: string) {
  try {
    await putV1DownloadersById<true>({
      throwOnError: true,
      path: { id },
      body: {
        name: editBuf.value.name,
        type: editBuf.value.type,
        protocol: editBuf.value.protocol,
        url: editBuf.value.url,
        username: editBuf.value.username || undefined,
        password: editBuf.value.password || undefined,
        config_json: editBuf.value.config_json,
        enabled: editBuf.value.enabled,
        default: editBuf.value.default,
      },
    })
    downloaders.value = downloaders.value.map((d) => (d.id === id ? { ...d, ...editBuf.value } : d))
    editingId.value = null
    await loadDownloaders() // Reload to refresh default flags
  } catch (err: any) {
    downloadersError.value = err.message || 'Failed to update downloader'
  }
}

async function removeDownloader(id: string) {
  try {
    await deleteV1DownloadersById<true>({ throwOnError: true, path: { id } })
    downloaders.value = downloaders.value.filter((d) => d.id !== id)
  } catch (err: any) {
    downloadersError.value = err.message || 'Failed to delete downloader'
  }
}

const testingId = ref<string | null>(null)
async function testDownloader(downloader: Downloader) {
  if (!downloader.id) return
  testingId.value = downloader.id
  try {
    await postV1DownloadersByIdTest({ throwOnError: true, path: { id: downloader.id } })
    downloadersError.value = null
    // Show success message (could use a toast here)
  } catch (err: any) {
    downloadersError.value = err.message || 'Connection test failed'
  } finally {
    testingId.value = null
  }
}
</script>

<template>
  <div>
    <Message v-if="downloadersError" severity="error" @close="downloadersError = null">{{
      downloadersError
    }}</Message>

    <div class="grid gap-4 md:grid-cols-2">
      <Card>
        <template #title>Add Downloader</template>
        <template #content>
          <div class="flex flex-col gap-3">
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Name</label>
              <InputText v-model="newDownloader.name" placeholder="My qBittorrent" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Type</label>
              <Select
                v-model="newDownloader.type"
                :options="[{ label: 'qBittorrent', value: 'qbittorrent' }]"
                optionLabel="label"
                optionValue="value"
                disabled
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Protocol</label>
              <Select
                v-model="newDownloader.protocol"
                :options="[
                  { label: 'Torrent', value: 'torrent' },
                  { label: 'Usenet', value: 'usenet' },
                ]"
                optionLabel="label"
                optionValue="value"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">URL</label>
              <InputText v-model="newDownloader.url" placeholder="http://localhost:8080" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Username (optional)</label>
              <InputText v-model="newDownloader.username" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Password (optional)</label>
              <Password v-model="newDownloader.password" :feedback="false" toggleMask />
            </div>
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-color">Enabled</div>
              <ToggleSwitch v-model="newDownloader.enabled" />
            </div>
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-color">Default</div>
              <ToggleSwitch v-model="newDownloader.default" />
            </div>
            <div class="flex justify-end">
              <Button label="Create" :loading="isCreating" @click="createDownloader" />
            </div>
          </div>
        </template>
      </Card>

      <Card>
        <template #title>Downloaders</template>
        <template #content>
          <div v-if="downloadersLoading" class="space-y-2">
            <Skeleton height="2.5rem" />
            <Skeleton height="2.5rem" />
          </div>
          <div v-else class="space-y-3">
            <div
              v-for="downloader in downloaders"
              :key="downloader.id"
              class="border rounded p-3 flex flex-col gap-2"
            >
              <div class="flex items-center justify-between">
                <div class="font-medium">
                  {{ downloader.name }}
                  <span class="text-sm text-muted-color">({{ downloader.type }}/{{ downloader.protocol }})</span>
                  <span
                    v-if="downloader.default"
                    class="ml-2 text-xs px-2 py-0.5 bg-green-500 text-white rounded"
                  >
                    Default
                  </span>
                </div>
                <div class="flex items-center gap-2">
                  <Button
                    size="small"
                    label="Test"
                    :loading="testingId === downloader.id"
                    @click="testDownloader(downloader)"
                  />
                  <Button size="small" label="Edit" @click="startEdit(downloader)" />
                  <Button
                    size="small"
                    label="Delete"
                    severity="danger"
                    @click="downloader.id && removeDownloader(downloader.id)"
                  />
                </div>
              </div>
              <div class="text-sm text-muted-color">{{ downloader.url }}</div>
              <div class="flex gap-4 text-xs">
                <div>
                  Enabled:
                  <span :class="downloader.enabled ? 'text-green-500' : 'text-red-500'">{{
                    downloader.enabled ? 'Yes' : 'No'
                  }}</span>
                </div>
                <div>
                  Default:
                  <span :class="downloader.default ? 'text-green-500' : 'text-muted-color'">{{
                    downloader.default ? 'Yes' : 'No'
                  }}</span>
                </div>
              </div>

              <div v-if="editingId === downloader.id" class="mt-2 border-t pt-3 space-y-3">
                <div class="grid gap-3 md:grid-cols-2">
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Name</label>
                    <InputText v-model="editBuf.name" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Type</label>
                    <Select
                      v-model="editBuf.type"
                      :options="[{ label: 'qBittorrent', value: 'qbittorrent' }]"
                      optionLabel="label"
                      optionValue="value"
                      disabled
                    />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Protocol</label>
                    <Select
                      v-model="editBuf.protocol"
                      :options="[
                        { label: 'Torrent', value: 'torrent' },
                        { label: 'Usenet', value: 'usenet' },
                      ]"
                      optionLabel="label"
                      optionValue="value"
                    />
                  </div>
                  <div class="md:col-span-2 flex flex-col gap-1">
                    <label class="text-sm text-muted-color">URL</label>
                    <InputText v-model="editBuf.url" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Username (optional)</label>
                    <InputText v-model="editBuf.username" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Password (optional)</label>
                    <Password v-model="editBuf.password" :feedback="false" toggleMask placeholder="Leave blank to keep current" />
                  </div>
                  <div class="md:col-span-2 flex items-center justify-between">
                    <div class="text-sm text-muted-color">Enabled</div>
                    <ToggleSwitch v-model="editBuf.enabled" />
                  </div>
                  <div class="md:col-span-2 flex items-center justify-between">
                    <div class="text-sm text-muted-color">Default</div>
                    <ToggleSwitch v-model="editBuf.default" />
                  </div>
                </div>
                <div class="flex justify-end gap-2">
                  <Button
                    size="small"
                    label="Cancel"
                    severity="secondary"
                    @click="editingId = null"
                  />
                  <Button size="small" label="Save" @click="downloader.id && saveEdit(downloader.id)" />
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

