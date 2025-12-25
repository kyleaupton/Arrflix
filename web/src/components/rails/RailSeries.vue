<template>
  <Rail :title="rail.title">
    <div v-for="series in rail.series" :key="series.tmdbId" class="flex-shrink-0">
      <Poster :item="series" :to="{ path: `/series/${series.tmdbId}` }" />
    </div>
  </Rail>
</template>

<script setup lang="ts">
import { type ModelRail } from '@/client/types.gen'
import Rail from './Rail.vue'
import Poster from '@/components/poster/Poster.vue'

defineProps<{
  rail: ModelRail & { type: 'series' }
}>()
</script>

