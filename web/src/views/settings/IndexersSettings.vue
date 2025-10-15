<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import Button from 'primevue/button'
import { getV1IndexersConfiguredOptions } from '@/client/@tanstack/vue-query.gen'
import { type JackettIndexerConfig } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'
import {
  indexerColumns,
  createIndexerActions,
} from '@/components/tables/configs/indexerTableConfig'
import DataTable from '@/components/tables/DataTable.vue'
import AddIndexerModal from '@/components/modals/AddIndexerModal.vue'

const { data: indexers, isLoading, error, refetch } = useQuery(getV1IndexersConfiguredOptions())

// Modal state
const showAddModal = ref(false)

const handleEdit = (indexer: JackettIndexerConfig) => {
  console.log('Edit indexer:', indexer)
  // TODO: Implement edit functionality
}

const handleToggle = (indexer: JackettIndexerConfig) => {
  console.log('Toggle indexer:', indexer)
  // TODO: Implement toggle functionality
}

const handleDelete = (indexer: JackettIndexerConfig) => {
  console.log('Delete indexer:', indexer)
  // TODO: Implement delete functionality
}

const handleAddIndexer = () => {
  showAddModal.value = true
}

const handleIndexerAdded = (newIndexer: JackettIndexerConfig) => {
  console.log('Indexer added:', newIndexer)
  // Refetch the indexers list to show the new one
  refetch()
}

const indexerActions = createIndexerActions(handleEdit, handleToggle, handleDelete)
</script>

<template>
  <div class="indexers-settings">
    <div class="card">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold mb-2">Indexers</h3>
            <p class="text-muted-color">Configure your media indexers and search providers.</p>
          </div>
          <Button
            label="Add Indexer"
            :icon="PrimeIcons.PLUS"
            severity="primary"
            @click="handleAddIndexer"
          />
        </div>

        <DataTable
          :data="indexers || []"
          :columns="indexerColumns"
          :actions="indexerActions"
          :loading="isLoading"
          :empty-message="error ? `Error: ${error.message}` : 'No indexers configured'"
          searchable
          search-placeholder="Search indexers..."
          paginator
          :rows="10"
        />
      </div>
    </div>

    <!-- Add Indexer Modal -->
    <AddIndexerModal v-model:visible="showAddModal" @indexer-added="handleIndexerAdded" />
  </div>
</template>

<style scoped>
.indexers-settings {
  max-width: 100%;
}
</style>
