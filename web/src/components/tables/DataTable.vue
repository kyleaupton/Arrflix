<!-- eslint-disable @typescript-eslint/no-explicit-any -->

<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, ref, onMounted, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import { PrimeIcons } from '@/icons'

export interface TableColumn<T = any> {
  key: keyof T | string
  label: string
  sortable?: boolean
  filterable?: boolean
  width?: string
  align?: 'left' | 'center' | 'right'
  render?: (value: any, row: T) => any
}

export interface TableAction<T = any> {
  key: string
  label: string
  icon?: string
  severity?: 'primary' | 'secondary' | 'success' | 'info' | 'warning' | 'danger'
  variant?: 'text' | 'outlined' | 'filled'
  disabled?: (row: T) => boolean
  visible?: (row: T) => boolean
  command: (row: T) => void
}

interface Props {
  data?: T[]
  columns: TableColumn<T>[]
  actions?: TableAction<T>[]
  loading?: boolean
  emptyMessage?: string
  selectionMode?: 'single' | 'multiple' | 'checkbox' | 'radiobutton'
  selectable?: boolean
  paginator?: boolean
  rows?: number
  scrollable?: boolean
  scrollHeight?: string
  virtualScrollerOptions?: any
  sortField?: string
  sortOrder?: 1 | -1
  filterable?: boolean
  searchable?: boolean
  searchPlaceholder?: string
  // Async loading support (legacy)
  asyncData?: () => Promise<T[]>
  autoLoad?: boolean
  // TanStack Query support - pass query options to use TanStack Query for data fetching
  // This takes precedence over asyncData when both are provided
  queryOptions?: any
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
  loading: false,
  emptyMessage: 'No data available',
  selectionMode: 'single',
  selectable: false,
  paginator: true,
  rows: 10,
  sortField: '',
  sortOrder: 1,
  filterable: true,
  searchable: true,
  searchPlaceholder: 'Search...',
  autoLoad: true,
})

const emit = defineEmits<{
  selectionChange: [selection: T | T[] | null]
  rowSelect: [row: T]
  rowUnselect: [row: T]
  'data-loaded': [data: T[]]
  'load-error': [error: Error]
  'query-success': [data: T[]]
  'query-error': [error: Error]
}>()

const selectedRows = ref<T | T[] | null>(null)
const globalFilter = ref('')

// Async loading state (legacy)
const legacyAsyncData = ref<T[]>([])
const asyncLoading = ref(false)
const asyncError = ref<Error | null>(null)

// TanStack Query integration
const queryResult = props.queryOptions ? useQuery(props.queryOptions) : null

// Determine which data source to use
// Priority: TanStack Query > asyncData > static data
const tableData = computed(() => {
  if (props.queryOptions && queryResult?.data?.value) {
    return queryResult.data.value
  }
  return props.asyncData ? legacyAsyncData.value : props.data || []
})

const isLoading = computed(() => {
  if (props.queryOptions && queryResult?.isLoading?.value) {
    return queryResult.isLoading.value
  }
  return props.asyncData ? asyncLoading.value : props.loading
})

const queryError = computed(() => {
  if (props.queryOptions && queryResult?.error?.value) {
    return queryResult.error.value
  }
  return asyncError.value
})

const filteredData = computed(() => {
  const data = tableData.value

  if (!props.searchable || !globalFilter.value) {
    return data
  }

  const filter = globalFilter.value.toLowerCase()
  return data.filter((row) => {
    return props.columns.some((column) => {
      const value = getNestedValue(row, column.key as string)
      return String(value).toLowerCase().includes(filter)
    })
  })
})

const getNestedValue = (obj: any, path: string) => {
  return path.split('.').reduce((current, key) => current?.[key], obj)
}

const handleSelectionChange = (selection: T | T[] | null) => {
  selectedRows.value = selection
  emit('selectionChange', selection)
}

const renderCell = (column: TableColumn<T>, row: T) => {
  const value = getNestedValue(row, column.key as string)

  if (column.render) {
    return column.render(value, row)
  }

  return value
}

const getActionSeverity = (action: TableAction<T>) => {
  return action.severity || 'secondary'
}

const getActionVariant = (action: TableAction<T>) => {
  return action.variant || 'text'
}

const isActionDisabled = (action: TableAction<T>, row: T) => {
  return action.disabled ? action.disabled(row) : false
}

const isActionVisible = (action: TableAction<T>, row: T) => {
  return action.visible ? action.visible(row) : true
}

// Async loading functions
const loadAsyncData = async () => {
  if (!props.asyncData) return

  try {
    asyncLoading.value = true
    asyncError.value = null
    const data = await props.asyncData()
    legacyAsyncData.value = data
    emit('data-loaded', data)
  } catch (error) {
    asyncError.value = error as Error
    emit('load-error', error as Error)
  } finally {
    asyncLoading.value = false
  }
}

// Expose load function for manual triggering
const loadData = () => {
  if (props.asyncData) {
    loadAsyncData()
  }
}

// Auto-load on mount if enabled
onMounted(() => {
  if (props.autoLoad && props.asyncData) {
    loadAsyncData()
  }
})

// Watch for asyncData changes
watch(
  () => props.asyncData,
  (newAsyncData) => {
    if (newAsyncData && props.autoLoad) {
      loadAsyncData()
    }
  },
)

// Watch for TanStack Query state changes
watch(
  () => queryResult?.data?.value,
  (newData) => {
    if (newData && props.queryOptions) {
      emit('query-success', newData)
    }
  },
)

watch(
  () => queryResult?.error?.value,
  (newError) => {
    if (newError && props.queryOptions) {
      emit('query-error', newError)
    }
  },
)

// Expose methods for parent components
defineExpose({
  loadData,
  asyncData: legacyAsyncData.value,
  asyncLoading: asyncLoading.value,
  asyncError: asyncError.value,
  // TanStack Query methods
  refetch: queryResult?.refetch,
  queryResult: queryResult,
  queryError: queryError.value,
})
</script>

<template>
  <div class="data-table-container">
    <!-- Search Bar -->
    <div v-if="searchable" class="mb-4">
      <!-- <div class="p-input-icon-left w-full max-w-md">
        <i :class="PrimeIcons.SEARCH" class="text-muted-color" />
        <input
          v-model="globalFilter"
          type="text"
          :placeholder="searchPlaceholder"
          class="w-full p-inputtext p-component"
        />
      </div> -->

      <IconField class="search-field">
        <InputIcon :class="PrimeIcons.SEARCH" />
        <InputText
          v-model="globalFilter"
          class="w-full"
          :placeholder="searchPlaceholder"
          variant="filled"
          size="small"
        />
      </IconField>
    </div>

    <!-- Data Table -->
    <DataTable
      :value="filteredData"
      :loading="isLoading"
      :selection-mode="selectable ? selectionMode : undefined"
      :selection="selectedRows"
      :paginator="paginator"
      :rows="rows"
      :scrollable="scrollable"
      :scroll-height="scrollHeight"
      :virtual-scroller-options="virtualScrollerOptions"
      :sort-field="sortField"
      :sort-order="sortOrder"
      :global-filter-fields="searchable ? columns.map((col) => col.key as string) : undefined"
      data-key="id"
      class="p-datatable-sm"
      @selection-change="handleSelectionChange"
      @row-select="emit('rowSelect', $event.data)"
      @row-unselect="emit('rowUnselect', $event.data)"
    >
      <!-- Selection Column -->
      <Column v-if="selectable" selection-mode="multiple" header-style="width: 3rem" />

      <!-- Data Columns -->
      <Column
        v-for="column in columns"
        :key="column.key as string"
        :field="column.key as string"
        :header="column.label"
        :sortable="column.sortable"
        :style="{ width: column.width, textAlign: column.align || 'left' }"
      >
        <template #body="{ data }">
          <span v-if="!column.render">{{ renderCell(column, data) }}</span>
          <div v-else v-html="renderCell(column, data)"></div>
        </template>
      </Column>

      <!-- Actions Column -->
      <Column
        v-if="actions && actions.length > 0"
        header="Actions"
        header-style="width: 8rem"
        body-style="text-align: center"
      >
        <template #body="{ data }">
          <div class="flex gap-1 justify-center">
            <Button
              v-for="action in actions"
              :key="action.key"
              :label="action.label"
              :icon="action.icon"
              :severity="getActionSeverity(action)"
              :variant="getActionVariant(action)"
              :disabled="isActionDisabled(action, data)"
              size="small"
              v-show="isActionVisible(action, data)"
              @click="action.command(data)"
            />
          </div>
        </template>
      </Column>

      <!-- Empty State -->
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          {{ emptyMessage }}
        </div>
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.data-table-container {
  width: 100%;
}

:deep(.p-datatable) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.p-datatable-header) {
  background: var(--p-card-background);
  border-bottom: 1px solid var(--p-border-color);
}

:deep(.p-datatable-tbody > tr) {
  transition: background-color 0.2s ease;
}

:deep(.p-datatable-tbody > tr:hover) {
  background: var(--p-emphasis-background);
}

:deep(.p-datatable-tbody > tr.p-highlight) {
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}
</style>
