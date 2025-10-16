<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { PrimeIcons } from '@/icons'
import { type JackettIndexerDetails } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import {
  availableIndexerColumns,
  createAvailableIndexerActions,
} from '@/components/tables/configs/availableIndexerTableConfig'
import { getV1IndexersUnconfiguredOptions } from '@/client/@tanstack/vue-query.gen'

defineProps<{
  selectedIndexer: JackettIndexerDetails | null
}>()

const emit = defineEmits<{
  'indexer-selected': [indexer: JackettIndexerDetails]
}>()

const queryOptions = getV1IndexersUnconfiguredOptions()

// Create actions for the available indexers table
const availableIndexerActions = createAvailableIndexerActions((indexer: JackettIndexerDetails) => {
  emit('indexer-selected', indexer)
})
</script>

<template>
  <div class="select-indexer-step">
    <h3 class="text-lg font-semibold mb-4">Select Indexer Type</h3>
    <p class="text-muted-color mb-6">Choose from available indexer types to configure.</p>

    <DataTable
      ref="dataTableRef"
      :query-options="queryOptions"
      :columns="availableIndexerColumns"
      :actions="availableIndexerActions"
      :auto-load="false"
      empty-message="No unconfigured indexers available"
      searchable
      search-placeholder="Search available indexers..."
      paginator
      :rows="15"
      selectable
      selection-mode="single"
      @selection-change="
        (selection) => {
          if (selection && !Array.isArray(selection)) {
            emit('indexer-selected', selection)
          }
        }
      "
      @data-loaded="(data) => console.log('Loaded indexers:', data.length)"
      @load-error="(error) => console.error('Failed to load indexers:', error)"
    />
  </div>
</template>
