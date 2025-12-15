<template>
  <div class="downloader-reference">
    <div v-if="isLoading" class="inline-flex items-center gap-1 text-muted-color">
      <i class="pi pi-spin pi-spinner text-sm"></i>
      <span class="text-sm">Loading...</span>
    </div>
    <div v-else-if="error" class="inline-flex items-center gap-1 text-red-400">
      <i class="pi pi-exclamation-triangle text-sm"></i>
      <span class="text-sm">Error</span>
    </div>
    <div v-else-if="downloader" class="inline-flex items-center gap-2">
      <span class="font-semibold">{{ downloader.name }}</span>
      <Badge :value="downloader.type" severity="secondary" size="small" />
      <Badge :value="downloader.protocol" severity="info" size="small" />
      <Badge v-if="downloader.default" value="Default" severity="success" size="small" />
    </div>
    <span v-else class="text-muted-color text-sm">Unknown</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import Badge from 'primevue/badge'
import { getV1DownloadersByIdOptions } from '@/client/@tanstack/vue-query.gen'

const props = defineProps<{
  downloaderId: string
}>()

const {
  data: downloader,
  isLoading,
  error,
} = useQuery(
  computed(() =>
    getV1DownloadersByIdOptions({
      path: { id: props.downloaderId },
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } as any),
  ),
)
</script>

<style scoped>
.downloader-reference {
  display: inline-flex;
  align-items: center;
}
</style>

