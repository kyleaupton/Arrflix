import { type TableColumn, type TableAction } from '../types'
import { type HandlersNameTemplateSwagger } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'

export const nameTemplateColumns: TableColumn<HandlersNameTemplateSwagger>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    filterable: true,
  },
  {
    key: 'type',
    label: 'Type',
    sortable: true,
    filterable: true,
    width: '120px',
    render: (value: string) => {
      const label = value === 'movie' ? 'Movie' : value === 'series' ? 'Series' : value
      return `<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">${label}</span>`
    },
  },
  {
    key: 'template',
    label: 'Template',
    sortable: true,
    filterable: true,
    render: (value: string, row: HandlersNameTemplateSwagger) => {
      if (row.type === 'series' && (row.series_show_template || row.series_season_template)) {
        const parts = [
          row.series_show_template || '...',
          row.series_season_template || '...',
          value || '...',
        ]
        return `<div class="flex flex-col gap-0.5">
          <span class="font-mono text-xs opacity-70">${parts[0]} /</span>
          <span class="font-mono text-xs opacity-70">${parts[1]} /</span>
          <span class="font-mono text-sm">${parts[2]}</span>
        </div>`
      }
      return `<span class="font-mono text-sm">${value || ''}</span>`
    },
  },
  {
    key: 'default',
    label: 'Default',
    sortable: true,
    width: '100px',
    align: 'center',
    render: (value: boolean) => {
      return value
        ? '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">Yes</span>'
        : '<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300">No</span>'
    },
  },
]

export const createNameTemplateActions = (
  onEdit: (template: HandlersNameTemplateSwagger) => void,
  onDelete: (template: HandlersNameTemplateSwagger) => void,
): TableAction<HandlersNameTemplateSwagger>[] => [
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
