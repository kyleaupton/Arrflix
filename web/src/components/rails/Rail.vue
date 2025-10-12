<template>
  <div>
    <h1>{{ rail.title }}</h1>

    <div class="flex gap-2 overflow-x-auto">
      <template v-if="rail.type === 'movie'">
        <Poster v-for="movie in rail.movies" :key="movie.tmdbId" :item="movie" />
      </template>

      <template v-else-if="rail.type === 'series'">
        <Poster v-for="series in rail.series" :key="series.tmdbId" :item="series" />
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { type ModelRail } from '@/client/types.gen'
import Poster from '@/components/poster/Poster.vue'

defineProps<{
  rail: ModelRail
}>()
</script>

<style scoped></style>
