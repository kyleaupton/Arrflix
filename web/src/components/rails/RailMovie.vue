<template>
  <Rail :title="rail.title">
    <div v-for="movie in rail.movies" :key="movie.tmdbId" class="flex-shrink-0">
      <Poster :item="movie" :to="{ path: `/movie/${movie.tmdbId}` }" :is-downloading="movie.isDownloading" />
    </div>
  </Rail>
</template>

<script setup lang="ts">
import { type ModelRail } from '@/client/types.gen'
import Rail from './Rail.vue'
import Poster from '@/components/poster/Poster.vue'

defineProps<{
  rail: ModelRail & { type: 'movie' }
}>()
</script>

