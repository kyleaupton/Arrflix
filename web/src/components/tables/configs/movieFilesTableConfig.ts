import { type TableColumn } from '../types'
import type { ModelFileInfo } from '@/client/types.gen'
import { Badge } from '@/components/ui/badge'
import LibraryReference from '@/components/references/LibraryReference.vue'
import { h } from 'vue'

export const movieFilesColumns: TableColumn<ModelFileInfo>[] = [
  {
    key: 'path',
    label: 'Path',
    sortable: true,
    filterable: true,
    render: (value: string) => {
      return h('span', { class: 'font-mono text-sm' }, value || '')
    },
  },
  {
    key: 'libraryId',
    label: 'Library',
    sortable: true,
    filterable: true,
    width: '200px',
    render: (value: string) => {
      return h(LibraryReference, { libraryId: value })
    },
  },
  {
    key: 'status',
    label: 'Status',
    sortable: true,
    filterable: true,
    width: '140px',
    render: (value: string) => {
      const statusColors: Record<string, string> = {
        available: 'bg-green-500 text-white',
        downloading: 'bg-blue-500 text-white',
        importing: 'bg-yellow-500 text-white',
        missing: 'bg-gray-500 text-white',
        failed: 'bg-red-500 text-white',
        deleted: 'bg-gray-600 text-white',
      }
      const colorClass = statusColors[value] || 'bg-gray-500 text-white'
      return h(
        Badge,
        {
          class: `${colorClass} capitalize border-transparent`,
        },
        () => value,
      )
    },
  },
]

