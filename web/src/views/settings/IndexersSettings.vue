<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { Plus } from 'lucide-vue-next'
import { getV1IndexersConfiguredOptions } from '@/client/@tanstack/vue-query.gen'
import { type ModelIndexerDefinition } from '@/client/types.gen'
import {
  indexerColumns,
  createIndexerActions,
} from '@/components/tables/configs/indexerTableConfig'
import DataTable from '@/components/tables/DataTable.vue'
import AddIndexerModal from '@/components/modals/AddIndexerModal.vue'
import { useModal } from '@/composables/useModal'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const { data: indexers, isLoading, error, refetch } = useQuery(getV1IndexersConfiguredOptions())
const modal = useModal()

const handleEdit = (indexer: ModelIndexerDefinition) => {
  console.log('Edit indexer:', indexer)
  // TODO: Implement edit functionality
}

const handleToggle = (indexer: ModelIndexerDefinition) => {
  console.log('Toggle indexer:', indexer)
  // TODO: Implement toggle functionality
}

const handleDelete = (indexer: ModelIndexerDefinition) => {
  console.log('Delete indexer:', indexer)
  // TODO: Implement delete functionality
}

const handleAddIndexer = () => {
  modal.open(AddIndexerModal, {
    props: {
      class: 'max-w-[90vw] sm:max-w-4xl lg:max-w-6xl',
    },
    onClose: (result) => {
      if ((result?.data as { indexerAdded?: boolean })?.indexerAdded) {
        refetch()
      }
    },
  })
}

const indexerActions = createIndexerActions(handleEdit, handleToggle, handleDelete)
</script>

<template>
  <div class="flex flex-col gap-6">
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-xl font-semibold mb-2">Indexers</CardTitle>
            <p class="text-sm text-muted-foreground">
              Configure your media indexers and search providers.
            </p>
          </div>
          <Button @click="handleAddIndexer">
            <Plus class="mr-2 size-4" />
            Add Indexer
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="isLoading" class="space-y-3">
          <Skeleton class="h-12 w-full" />
          <Skeleton class="h-12 w-full" />
          <Skeleton class="h-12 w-full" />
        </div>
        <DataTable
          v-else
          :data="(indexers || []) as unknown as ModelIndexerDefinition[]"
          :columns="indexerColumns"
          :actions="indexerActions"
          :loading="isLoading"
          :empty-message="error ? `Error: ${error.message}` : 'No indexers configured'"
          searchable
          search-placeholder="Search indexers..."
          paginator
          :rows="10"
        />
      </CardContent>
    </Card>
  </div>
</template>
