<!-- eslint-disable @typescript-eslint/no-explicit-any -->

<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, ref } from 'vue'
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
  data: T[]
  columns: TableColumn<T>[]
  actions?: TableAction<T>[]
  loading?: boolean
  emptyMessage?: string
  selectionMode?: 'single' | 'multiple' | 'checkbox' | 'radiobutton'
  selectable?: boolean
  paginator?: boolean
  rows?: number
  sortField?: string
  sortOrder?: 1 | -1
  filterable?: boolean
  searchable?: boolean
  searchPlaceholder?: string
}

const props = withDefaults(defineProps<Props>(), {
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
})

const emit = defineEmits<{
  selectionChange: [selection: T | T[] | null]
  rowSelect: [row: T]
  rowUnselect: [row: T]
}>()

const selectedRows = ref<T | T[] | null>(null)
const globalFilter = ref('')

const filteredData = computed(() => {
  if (!props.searchable || !globalFilter.value) {
    return props.data
  }

  const filter = globalFilter.value.toLowerCase()
  return props.data.filter((row) => {
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
      :loading="loading"
      :selection-mode="selectable ? selectionMode : undefined"
      :selection="selectedRows"
      :paginator="paginator"
      :rows="rows"
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
          <component :is="column.render ? 'div' : 'span'" v-html="renderCell(column, data)" />
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
              :severity="getActionSeverity(action, data)"
              :variant="getActionVariant(action, data)"
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
