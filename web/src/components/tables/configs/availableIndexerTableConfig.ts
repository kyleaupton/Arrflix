import { type TableColumn, type TableAction } from '../DataTable.vue'
import { type ModelIndexerDefinition } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'

export const availableIndexerColumns: TableColumn<ModelIndexerDefinition>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    filterable: true,
    width: '200px',
  },
  {
    key: 'description',
    label: 'Description',
    sortable: true,
    filterable: true,
    width: '300px',
  },
  {
    key: 'language',
    label: 'Language',
    sortable: true,
    filterable: true,
    width: '120px',
    align: 'center',
    render: (value: string) => {
      return `<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">${value}</span>`
    },
  },
]

export const createAvailableIndexerActions = (
  onSelect: (indexer: ModelIndexerDefinition) => void,
): TableAction<ModelIndexerDefinition>[] => [
  {
    key: 'select',
    label: 'Select',
    icon: PrimeIcons.CHECK,
    severity: 'primary',
    variant: 'outlined',
    command: onSelect,
  },
]
