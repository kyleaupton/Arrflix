import { computed } from 'vue'
import { useBreakpoints } from '@vueuse/core'

const breakpoints = useBreakpoints({
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
})

export function useGridColumns(rowCount = 3) {
  const columns = computed(() => {
    if (breakpoints.xl.value) return 6
    if (breakpoints.lg.value) return 5
    if (breakpoints.md.value) return 4
    if (breakpoints.sm.value) return 3
    return 2
  })

  const pageSize = computed(() => columns.value * rowCount)

  return { columns, pageSize }
}
