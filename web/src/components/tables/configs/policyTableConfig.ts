import { type TableColumn, type TableAction } from '../DataTable.vue'
import { type DbgenPolicy } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'

export const policyColumns: TableColumn<DbgenPolicy>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    filterable: true,
  },
  {
    key: 'description',
    label: 'Description',
    sortable: true,
    filterable: true,
  },
  {
    key: 'priority',
    label: 'Priority',
    sortable: true,
    filterable: true,
    width: '100px',
    align: 'center',
  },
  {
    key: 'enabled',
    label: 'Status',
    sortable: true,
    width: '120px',
    align: 'center',
    render: (value: boolean) => {
      return value
        ? '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">Enabled</span>'
        : '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300">Disabled</span>'
    },
  },
]

export const createPolicyActions = (
  onEdit: (policy: DbgenPolicy) => void,
  onRule: (policy: DbgenPolicy) => void,
  onActions: (policy: DbgenPolicy) => void,
  onDelete: (policy: DbgenPolicy) => void,
): TableAction<DbgenPolicy>[] => [
  {
    key: 'rule',
    label: 'Rule',
    icon: PrimeIcons.COG,
    severity: 'secondary',
    variant: 'text',
    command: onRule,
  },
  {
    key: 'actions',
    label: 'Actions',
    icon: PrimeIcons.LIST,
    severity: 'secondary',
    variant: 'text',
    command: onActions,
  },
  {
    key: 'edit',
    label: 'Edit',
    icon: PrimeIcons.PENCIL,
    severity: 'primary',
    variant: 'text',
    command: onEdit,
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

