<script setup lang="ts">
import { ref } from 'vue'
import { useQuery, useMutation } from '@tanstack/vue-query'
import { Plus } from 'lucide-vue-next'
import {
  getV1LibrariesOptions,
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
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import LibraryDialog from '@/components/modals/LibraryDialog.vue'

// Data queries
const { data: libraries, isLoading, refetch } = useQuery(getV1LibrariesOptions())
const modal = useModal()

// Mutations
const deleteLibraryMutation = useMutation(deleteV1LibrariesByIdMutation())
const scanLibraryMutation = useMutation(postV1LibrariesByIdScanMutation())

// State
const libraryError = ref<string | null>(null)
const scanningId = ref<string | null>(null)

// Handlers
const handleAddLibrary = () => {
  modal.open(LibraryDialog, {
    props: {
      library: null,
    },
    onClose: () => {
      refetch()
    },
  })
}

const handleEditLibrary = (library: HandlersLibrarySwagger) => {
  modal.open(LibraryDialog, {
    props: {
      library,
    },
    onClose: () => {
      refetch()
    },
  })
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
  } catch (err) {
    libraryError.value = err instanceof Error ? err.message : 'Failed to delete library'
  }
}

const handleScanLibrary = async (library: HandlersLibrarySwagger) => {
  if (!library.id) return
  scanningId.value = library.id
  try {
    await scanLibraryMutation.mutateAsync({ path: { id: library.id } })
    libraryError.value = null
  } catch (err) {
    libraryError.value = err instanceof Error ? err.message : 'Failed to scan library'
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
  <div class="flex flex-col gap-6">
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-xl font-semibold mb-2">Library Settings</CardTitle>
            <p class="text-sm text-muted-foreground">
              Configure libraries to organize your media content.
            </p>
          </div>
          <Button @click="handleAddLibrary">
            <Plus class="mr-2 size-4" />
            Add Library
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div
          v-if="libraryError"
          class="p-4 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive mb-4"
        >
          {{ libraryError }}
        </div>
        <div v-if="isLoading" class="space-y-3">
          <Skeleton class="h-12 w-full" />
          <Skeleton class="h-12 w-full" />
          <Skeleton class="h-12 w-full" />
        </div>
        <DataTable
          v-else
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
      </CardContent>
    </Card>
  </div>
</template>
