import { type TableColumn, type TableAction } from '../DataTable.vue'
import { type JackettIndexerConfig } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'

export const indexerColumns: TableColumn<JackettIndexerConfig>[] = [
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
    key: 'type',
    label: 'Type',
    sortable: true,
    filterable: true,
    width: '120px',
  },
  {
    key: 'enabled',
    label: 'Status',
    sortable: true,
    width: '100px',
    align: 'center',
    render: (value: boolean) => {
      return value
        ? '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">Enabled</span>'
        : '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">Disabled</span>'
    },
  },
  {
    key: 'configured',
    label: 'Configured',
    sortable: true,
    width: '100px',
    align: 'center',
    render: (value: boolean) => {
      return value
        ? '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">Yes</span>'
        : '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">No</span>'
    },
  },
]

export const createIndexerActions = (
  onEdit: (indexer: JackettIndexerConfig) => void,
  onToggle: (indexer: JackettIndexerConfig) => void,
  onDelete: (indexer: JackettIndexerConfig) => void,
): TableAction<JackettIndexerConfig>[] => [
  {
    key: 'edit',
    label: 'Edit',
    icon: PrimeIcons.PENCIL,
    severity: 'primary',
    variant: 'text',
    command: onEdit,
  },
  {
    key: 'toggle',
    label: 'Toggle',
    icon: PrimeIcons.POWER_OFF,
    severity: 'secondary',
    variant: 'text',
    command: onToggle,
  },
  {
    key: 'delete',
    label: 'Delete',
    icon: PrimeIcons.TRASH,
    severity: 'danger',
    variant: 'text',
    command: onDelete,
  },
]
