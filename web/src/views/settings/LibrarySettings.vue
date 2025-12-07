<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Button from 'primevue/button'
import Message from 'primevue/message'
import Skeleton from 'primevue/skeleton'
import Select from 'primevue/select'
import {
  getV1Libraries,
  postV1Libraries,
  putV1LibrariesById,
  deleteV1LibrariesById,
  postV1LibrariesByIdScan,
} from '@/client/sdk.gen'

type Library = {
  id?: string
  name?: string
  type?: 'movie' | 'series' | string
  root_path?: string
  enabled?: boolean
  default?: boolean
}

const libraries = ref<Library[]>([])
const libsLoading = ref(false)
const libsError = ref<string | null>(null)

async function loadLibraries() {
  libsLoading.value = true
  libsError.value = null
  try {
    const res = await getV1Libraries<true>({ throwOnError: true })
    libraries.value = (res.data as Library[]) ?? []
  } catch {
    libsError.value = 'Failed to load libraries'
  } finally {
    libsLoading.value = false
  }
}

onMounted(loadLibraries)

// Create form
const newLib = ref<Required<Omit<Library, 'id'>>>({
  name: '',
  type: 'movie',
  root_path: '',
  enabled: true,
  default: false,
})
const isCreating = ref(false)
async function createLibrary() {
  if (!newLib.value.name || !newLib.value.root_path) return
  isCreating.value = true
  try {
    const res = await postV1Libraries<true>({
      throwOnError: true,
      body: {
        name: newLib.value.name,
        type: newLib.value.type,
        root_path: newLib.value.root_path,
        enabled: newLib.value.enabled,
        default: newLib.value.default,
      },
    })
    libraries.value = [...libraries.value, res.data as Library]
    newLib.value = { name: '', type: 'movie', root_path: '', enabled: true, default: false }
  } finally {
    isCreating.value = false
  }
}

// Edit helpers
const editingId = ref<string | null>(null)
const editBuf = ref<Required<Omit<Library, 'id'>>>({
  name: '',
  type: 'movie',
  root_path: '',
  enabled: true,
  default: false,
})

function startEdit(lib: Library) {
  editingId.value = lib.id ?? null
  editBuf.value = {
    name: lib.name ?? '',
    type: (lib.type as 'movie' | 'series') ?? 'movie',
    root_path: lib.root_path ?? '',
    enabled: lib.enabled ?? true,
    default: lib.default ?? false,
  }
}

async function saveEdit(id: string) {
  await putV1LibrariesById<true>({
    throwOnError: true,
    path: { id },
    body: {
      name: editBuf.value.name,
      type: editBuf.value.type,
      root_path: editBuf.value.root_path,
      enabled: editBuf.value.enabled,
      default: editBuf.value.default,
    },
  })
  libraries.value = libraries.value.map((l) => (l.id === id ? { ...l, ...editBuf.value } : l))
  editingId.value = null
}

async function removeLib(id: string) {
  await deleteV1LibrariesById<true>({ throwOnError: true, path: { id } })
  libraries.value = libraries.value.filter((l) => l.id !== id)
}

async function scanLib(lib: Library) {
  if (!lib.id) return
  await postV1LibrariesByIdScan({ throwOnError: true, path: { id: lib.id } })
}
</script>

<template>
  <div>
    <Message v-if="libsError" severity="error">{{ libsError }}</Message>

    <div class="grid gap-4 md:grid-cols-2">
      <Card>
        <template #title>Add Library</template>
        <template #content>
          <div class="flex flex-col gap-3">
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Name</label>
              <InputText v-model="newLib.name" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Type</label>
              <Select
                v-model="newLib.type"
                :options="[
                  { label: 'Movies', value: 'movie' },
                  { label: 'Series', value: 'series' },
                ]"
                optionLabel="label"
                optionValue="value"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Root Path</label>
              <InputText v-model="newLib.root_path" placeholder="/mnt/media/Movies" />
            </div>
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-color">Enabled</div>
              <ToggleSwitch v-model="newLib.enabled" />
            </div>
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-color">Default</div>
              <ToggleSwitch v-model="newLib.default" />
            </div>
            <div class="flex justify-end">
              <Button label="Create" :loading="isCreating" @click="createLibrary" />
            </div>
          </div>
        </template>
      </Card>

      <Card>
        <template #title>Libraries</template>
        <template #content>
          <div v-if="libsLoading" class="space-y-2">
            <Skeleton height="2.5rem" />
            <Skeleton height="2.5rem" />
          </div>
          <div v-else class="space-y-3">
            <div
              v-for="lib in libraries"
              :key="lib.id"
              class="border rounded p-3 flex flex-col gap-2"
            >
              <div class="flex items-center justify-between">
                <div class="font-medium">
                  {{ lib.name }}
                  <span class="text-sm text-muted-color">({{ lib.type }})</span>
                </div>
                <div class="flex items-center gap-2">
                  <Button size="small" label="Scan" @click="scanLib(lib)" />
                  <Button size="small" label="Edit" @click="startEdit(lib)" />
                  <Button
                    size="small"
                    label="Delete"
                    severity="danger"
                    @click="lib.id && removeLib(lib.id)"
                  />
                </div>
              </div>
              <div class="text-sm text-muted-color">{{ lib.root_path }}</div>
              <div class="flex gap-4 text-xs">
                <div>
                  Enabled:
                  <span :class="lib.enabled ? 'text-green-500' : 'text-red-500'">{{
                    lib.enabled ? 'Yes' : 'No'
                  }}</span>
                </div>
                <div>
                  Default:
                  <span :class="lib.default ? 'text-green-500' : 'text-muted-color'">{{
                    lib.default ? 'Yes' : 'No'
                  }}</span>
                </div>
              </div>

              <div v-if="editingId === lib.id" class="mt-2 border-t pt-3 space-y-3">
                <div class="grid gap-3 md:grid-cols-2">
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Name</label>
                    <InputText v-model="editBuf.name" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Type</label>
                    <Select
                      v-model="editBuf.type"
                      :options="[
                        { label: 'Movies', value: 'movie' },
                        { label: 'Series', value: 'series' },
                      ]"
                      optionLabel="label"
                      optionValue="value"
                    />
                  </div>
                  <div class="md:col-span-2 flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Root Path</label>
                    <InputText v-model="editBuf.root_path" />
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
                  <Button size="small" label="Save" @click="lib.id && saveEdit(lib.id)" />
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>
