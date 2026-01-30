import { type TableColumn, type TableAction } from '../types'
import type { DownloadJob } from '@/stores/downloadJobs'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { h } from 'vue'

// Status color mapping for unified import_status
const statusConfig: Record<string, { label: string; class: string }> = {
  download_pending: { label: 'Downloading', class: 'bg-blue-500 text-white' },
  download_failed: { label: 'Download Failed', class: 'bg-red-500 text-white' },
  download_cancelled: { label: 'Cancelled', class: 'bg-gray-600 text-white' },
  awaiting_import: { label: 'Awaiting Import', class: 'bg-yellow-500 text-white' },
  importing: { label: 'Importing', class: 'bg-yellow-600 text-white' },
  partial_failure: { label: 'Partial Failure', class: 'bg-orange-500 text-white' },
  import_failed: { label: 'Import Failed', class: 'bg-red-500 text-white' },
  fully_imported: { label: 'Imported', class: 'bg-green-500 text-white' },
  unknown: { label: 'Unknown', class: 'bg-gray-500 text-white' },
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
        h('div', { class: 'text-xs text-muted-foreground' }, `${row.protocol} • ${row.id.slice(0, 8)}`),
      ])
    },
  },
  {
    key: 'import_status',
    label: 'Status',
    sortable: true,
    filterable: true,
    width: '150px',
    render: (_value: string, row: DownloadJob) => {
      const config = statusConfig[row.import_status] || statusConfig['unknown']!
      return h(
        Badge,
        {
          class: `${config.class} border-transparent`,
        },
        () => config.label,
      )
    },
  },
  {
    key: 'progress',
    label: 'Progress',
    sortable: true,
    width: '220px',
    render: (_value: number | null | undefined, row: DownloadJob) => {
      // During download phase, show download progress
      if (row.import_status === 'download_pending') {
        const progressValue = Math.round((row.progress ?? 0) * 100)
        return h('div', { class: 'flex items-center gap-2' }, [
          h(Progress, { modelValue: progressValue, class: 'flex-1' }),
          h('span', { class: 'text-xs text-muted-foreground min-w-[3ch]' }, `${progressValue}%`),
        ])
      }

      // After download, show import progress
      const total = row.total_import_tasks || 0
      if (total === 0) {
        return h('span', { class: 'text-xs text-muted-foreground' }, '-')
      }

      const completed = row.completed_imports || 0
      const importProgress = Math.round((completed / total) * 100)
      return h('div', { class: 'flex items-center gap-2' }, [
        h(Progress, { modelValue: importProgress, class: 'flex-1' }),
        h(
          'span',
          { class: 'text-xs text-muted-foreground whitespace-nowrap' },
          `${completed}/${total} files`,
        ),
      ])
    },
  },
]

export const createDownloadJobActions = (
  onViewDetails: (job: DownloadJob) => void,
  onReimport: (job: DownloadJob, all: boolean) => void,
  onCancel: (job: DownloadJob) => void,
): TableAction<DownloadJob>[] => [
  {
    key: 'view_details',
    label: 'View Details',
    command: onViewDetails,
  },
  {
    key: 'reimport_failed',
    label: 'Re-import Failed',
    visible: (row: DownloadJob) => {
      return ['partial_failure', 'import_failed'].includes(row.import_status)
    },
    command: (job: DownloadJob) => onReimport(job, false),
  },
  {
    key: 'reimport_all',
    label: 'Re-import All',
    visible: (row: DownloadJob) => {
      return ['partial_failure', 'import_failed', 'fully_imported'].includes(row.import_status)
    },
    command: (job: DownloadJob) => onReimport(job, true),
  },
  {
    key: 'cancel',
    label: 'Cancel',
    severity: 'danger',
    visible: (row: DownloadJob) => {
      return row.import_status === 'download_pending'
    },
    command: onCancel,
  },
]
