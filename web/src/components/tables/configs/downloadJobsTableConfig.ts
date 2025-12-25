import { type TableColumn, type TableAction } from '../types'
import type { DbgenDownloadJob } from '@/client/types.gen'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Button } from '@/components/ui/button'
import { h } from 'vue'

type DownloadJob = DbgenDownloadJob & {
  progress?: number | null
  last_error?: string | null
  import_dest_path?: string | null
}

export const downloadJobColumns: TableColumn<DownloadJob>[] = [
  {
    key: 'candidate_title',
    label: 'Title',
    sortable: true,
    filterable: true,
    render: (value: string, row: DownloadJob) => {
      return h('div', { class: 'flex flex-col gap-1' }, [
        h('div', { class: 'font-medium' }, value || ''),
        h('div', { class: 'text-xs text-muted-foreground' }, `${row.protocol} • ${row.id}`),
      ])
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
        created: 'bg-gray-500 text-white',
        enqueued: 'bg-blue-500 text-white',
        downloading: 'bg-blue-600 text-white',
        completed: 'bg-green-500 text-white',
        importing: 'bg-yellow-500 text-white',
        imported: 'bg-green-600 text-white',
        failed: 'bg-red-500 text-white',
        cancelled: 'bg-gray-600 text-white',
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
  {
    key: 'progress',
    label: 'Progress',
    sortable: true,
    width: '220px',
    render: (value: number | null | undefined, row: DownloadJob) => {
      const progressValue = Math.round(((value ?? 0) as number) * 100)
      return h('div', { class: 'flex items-center gap-2' }, [
        h(Progress, { modelValue: progressValue, class: 'flex-1' }),
        h('span', { class: 'text-xs text-muted-foreground min-w-[3ch]' }, `${progressValue}%`),
      ])
    },
  },
  {
    key: 'import_dest_path',
    label: 'Imported To',
    sortable: true,
    filterable: true,
    render: (value: string | null | undefined) => {
      return h('span', { class: 'text-xs' }, value || '-')
    },
  },
  {
    key: 'last_error',
    label: 'Error',
    sortable: true,
    filterable: true,
    render: (value: string | null | undefined) => {
      return h('span', { class: 'text-xs text-destructive' }, value || '')
    },
  },
]

export const createDownloadJobActions = (
  onCancel: (job: DownloadJob) => void,
): TableAction<DownloadJob>[] => [
  {
    key: 'cancel',
    label: 'Cancel',
    severity: 'danger',
    visible: (row: DownloadJob) => {
      return !['cancelled', 'imported', 'failed'].includes(row.status)
    },
    command: onCancel,
  },
]

