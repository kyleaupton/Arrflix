import { type TableColumn, type TableAction } from '../DataTable.vue'
import { type JackettIndexerConfig } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'

export const availableIndexerColumns: TableColumn<JackettIndexerConfig>[] = [
  {
    key: 'title',
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
    key: 'type',
    label: 'Type',
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
  onSelect: (indexer: JackettIndexerConfig) => void,
): TableAction<JackettIndexerConfig>[] => [
  {
    key: 'select',
    label: 'Select',
    icon: PrimeIcons.CHECK,
    severity: 'primary',
    variant: 'outlined',
    command: onSelect,
  },
]
