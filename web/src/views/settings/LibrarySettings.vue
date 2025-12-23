<script setup lang="ts">
import { ref } from 'vue'
import { useQuery, useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Select from 'primevue/select'
import Message from 'primevue/message'
import { PrimeIcons } from '@/icons'
import {
  getV1LibrariesOptions,
  postV1LibrariesMutation,
  putV1LibrariesByIdMutation,
  deleteV1LibrariesByIdMutation,
  postV1LibrariesByIdScanMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type HandlersLibrarySwagger } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import {
  libraryColumns,
  createLibraryActions,
} from '@/components/tables/configs/libraryTableConfig'
import { useModal } from '@/composables/useModal'
import Page from '@/components/Page.vue'

// Data queries
const { data: libraries, isLoading, refetch } = useQuery(getV1LibrariesOptions())
const modal = useModal()

// Mutations
const createLibraryMutation = useMutation(postV1LibrariesMutation())
const updateLibraryMutation = useMutation(putV1LibrariesByIdMutation())
const deleteLibraryMutation = useMutation(deleteV1LibrariesByIdMutation())
const scanLibraryMutation = useMutation(postV1LibrariesByIdScanMutation())

// Modal state
const showLibraryModal = ref(false)
const editingLibrary = ref<HandlersLibrarySwagger | null>(null)
const libraryError = ref<string | null>(null)
const scanningId = ref<string | null>(null)

// Library form
const libraryForm = ref({
  name: '',
  type: 'movie' as 'movie' | 'series',
  root_path: '',
  enabled: true,
  default: false,
})

// Handlers
const handleAddLibrary = () => {
  editingLibrary.value = null
  libraryForm.value = { name: '', type: 'movie', root_path: '', enabled: true, default: false }
  libraryError.value = null
  showLibraryModal.value = true
}

const handleEditLibrary = (library: HandlersLibrarySwagger) => {
  editingLibrary.value = library
  libraryForm.value = {
    name: library.name || '',
    type: (library.type as 'movie' | 'series') || 'movie',
    root_path: library.root_path || '',
    enabled: library.enabled ?? true,
    default: library.default || false,
  }
  libraryError.value = null
  showLibraryModal.value = true
}

const handleSaveLibrary = async () => {
  if (!libraryForm.value.name || !libraryForm.value.root_path) {
    libraryError.value = 'Name and root path are required'
    return
  }

  try {
    if (editingLibrary.value?.id) {
      await updateLibraryMutation.mutateAsync({
        path: { id: editingLibrary.value.id },
        body: {
          name: libraryForm.value.name,
          type: libraryForm.value.type,
          root_path: libraryForm.value.root_path,
          enabled: libraryForm.value.enabled,
          default: libraryForm.value.default,
        },
      })
    } else {
      await createLibraryMutation.mutateAsync({
        body: {
          name: libraryForm.value.name,
          type: libraryForm.value.type,
          root_path: libraryForm.value.root_path,
          enabled: libraryForm.value.enabled,
          default: libraryForm.value.default,
        },
      })
    }
    showLibraryModal.value = false
    refetch()
  } catch (err: any) {
    libraryError.value = err.message || 'Failed to save library'
  }
}

const handleDeleteLibrary = async (library: HandlersLibrarySwagger) => {
  if (!library.id) return
  const confirmed = await modal.confirm({
    title: 'Delete Library',
    message: `Are you sure you want to delete "${library.name}"?`,
    severity: 'danger',
  })
  if (!confirmed) return
  try {
    await deleteLibraryMutation.mutateAsync({ path: { id: library.id } })
    refetch()
  } catch (err: any) {
    libraryError.value = err.message || 'Failed to delete library'
  }
}

const handleScanLibrary = async (library: HandlersLibrarySwagger) => {
  if (!library.id) return
  scanningId.value = library.id
  try {
    await scanLibraryMutation.mutateAsync({ path: { id: library.id } })
    libraryError.value = null
  } catch (err: any) {
    libraryError.value = err.message || 'Failed to scan library'
  } finally {
    scanningId.value = null
  }
}

const libraryActions = createLibraryActions(
  handleScanLibrary,
  handleEditLibrary,
  handleDeleteLibrary,
)
</script>

<template>
  <Page title="Library Settings">
    <div class="card">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <Button
            label="Add Library"
            :icon="PrimeIcons.PLUS"
            severity="primary"
            raised
            @click="handleAddLibrary"
          />
        </div>

        <Message v-if="libraryError" severity="error" @close="libraryError = null">{{
          libraryError
        }}</Message>

        <DataTable
          :data="libraries || []"
          :columns="libraryColumns"
          :actions="libraryActions"
          :loading="isLoading"
          empty-message="No libraries configured"
          searchable
          search-placeholder="Search libraries..."
          paginator
          :rows="10"
        />
      </div>
    </div>

    <!-- Library Modal -->
    <Dialog
      v-model:visible="showLibraryModal"
      :header="editingLibrary ? 'Edit Library' : 'Add Library'"
      :modal="true"
      :style="{ width: '600px' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Name</label>
          <InputText v-model="libraryForm.name" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Type</label>
          <Select
            v-model="libraryForm.type"
            :options="[
              { label: 'Movies', value: 'movie' },
              { label: 'Series', value: 'series' },
            ]"
            optionLabel="label"
            optionValue="value"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Root Path</label>
          <InputText v-model="libraryForm.root_path" placeholder="/mnt/media/Movies" />
        </div>
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium">Enabled</label>
          <ToggleSwitch v-model="libraryForm.enabled" />
        </div>
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium">Default</label>
          <ToggleSwitch v-model="libraryForm.default" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showLibraryModal = false" />
        <Button
          label="Save"
          :loading="createLibraryMutation.isPending.value || updateLibraryMutation.isPending.value"
          @click="handleSaveLibrary"
        />
      </template>
    </Dialog>
  </Page>
</template>

<style scoped></style>
