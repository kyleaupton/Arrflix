import { type TableColumn, type TableAction } from '../DataTable.vue'
import { type ModelDownloadCandidate } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'

// Helper function to format file size
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

// Helper function to format age
const formatAge = (ageHours: number): string => {
  if (ageHours < 1) {
    return '< 1 hour'
  } else if (ageHours < 24) {
    return `${Math.floor(ageHours)} hours`
  } else {
    const days = Math.floor(ageHours / 24)
    return `${days} day${days !== 1 ? 's' : ''}`
  }
}

export const downloadCandidateColumns: TableColumn<ModelDownloadCandidate>[] = [
  {
    key: 'title',
    label: 'Title',
    sortable: true,
    filterable: true,
  },
  {
    key: 'indexer',
    label: 'Indexer',
    sortable: true,
    filterable: true,
    width: '150px',
  },
  {
    key: 'protocol',
    label: 'Protocol',
    sortable: true,
    filterable: true,
    width: '100px',
    render: (value: string) => {
      const protocol = value.toLowerCase()
      const colorClass = protocol === 'torrent' ? 'text-blue-500' : 'text-green-500'
      return `<span class="font-medium ${colorClass}">${value}</span>`
    },
  },
  {
    key: 'size',
    label: 'Size',
    sortable: true,
    width: '100px',
    render: (value: number) => formatFileSize(value),
  },
  {
    key: 'seeders',
    label: 'Seeders',
    sortable: true,
    width: '100px',
    align: 'center',
    render: (value: number) => {
      const colorClass = value > 10 ? 'text-green-600' : value > 0 ? 'text-yellow-600' : 'text-red-600'
      return `<span class="font-semibold ${colorClass}">${value}</span>`
    },
  },
  {
    key: 'peers',
    label: 'Peers',
    sortable: true,
    width: '100px',
    align: 'center',
  },
  {
    key: 'ageHours',
    label: 'Age',
    sortable: true,
    width: '120px',
    render: (value: number) => formatAge(value),
  },
  {
    key: 'grabs',
    label: 'Grabs',
    sortable: true,
    width: '100px',
    align: 'center',
  },
]

export const createDownloadCandidateActions = (
  onEnqueue: (candidate: ModelDownloadCandidate) => void,
): TableAction<ModelDownloadCandidate>[] => [
  {
    key: 'enqueue',
    label: 'Enqueue Download',
    icon: PrimeIcons.DOWNLOAD,
    severity: 'primary',
    variant: 'text',
    command: onEnqueue,
  },
]

