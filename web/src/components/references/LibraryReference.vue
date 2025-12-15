<template>
  <div class="library-reference">
    <div v-if="isLoading" class="inline-flex items-center gap-1 text-muted-color">
      <i class="pi pi-spin pi-spinner text-sm"></i>
      <span class="text-sm">Loading...</span>
    </div>
    <div v-else-if="error" class="inline-flex items-center gap-1 text-red-400">
      <i class="pi pi-exclamation-triangle text-sm"></i>
      <span class="text-sm">Error</span>
    </div>
    <div v-else-if="library" class="inline-flex items-center gap-2">
      <span class="font-semibold">{{ library.name }}</span>
      <Badge
        :value="
          library.type === 'movie' ? 'Movie' : library.type === 'series' ? 'Series' : library.type
        "
        severity="secondary"
        size="small"
      />
      <Badge v-if="library.default" value="Default" severity="info" size="small" />
    </div>
    <span v-else class="text-muted-color text-sm">Unknown</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import Badge from 'primevue/badge'
import { getV1LibrariesByIdOptions } from '@/client/@tanstack/vue-query.gen'

const props = defineProps<{
  libraryId: string
}>()

const {
  data: library,
  isLoading,
  error,
} = useQuery(
  computed(() =>
    getV1LibrariesByIdOptions({
      path: { id: props.libraryId },
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } as any),
  ),
)
</script>

<style scoped>
.library-reference {
  display: inline-flex;
  align-items: center;
}
</style>
