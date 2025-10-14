<template>
  <Page>
    <div v-if="isLoading">Loading...</div>
    <div v-else-if="isError">Error</div>
    <div v-else-if="data">
      <MediaHero
        :title="data.title"
        :subtitle="releaseYear"
        :overview="data.overview"
        :poster-url="posterUrl"
        :backdrop-url="backdropUrl"
        :chips="movieChips"
      />

      <!-- TODO: sections like cast, recommendations, similar, etc. -->
    </div>
  </Page>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { getV1MovieByIdOptions } from '@/client/@tanstack/vue-query.gen'
import Page from '@/components/Page.vue'
import MediaHero from '@/components/media/MediaHero.vue'

const route = useRoute()

const id = computed(() => {
  const castAttept = Number(Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)
  if (isNaN(castAttept)) {
    throw new Error('Invalid movie ID')
  }

  return castAttept
})

const { isLoading, isError, data } = useQuery({
  ...getV1MovieByIdOptions({ path: { id: id.value } }),
})

const releaseYear = computed(() =>
  data.value?.releaseDate ? new Date(data.value.releaseDate).getFullYear().toString() : '',
)

// Image URLs: backend returns posterPath/backdropPath; map them to TMDB URLs
const posterUrl = computed(() =>
  data.value?.posterPath ? `https://image.tmdb.org/t/p/w342/${data.value.posterPath}` : undefined,
)
const backdropUrl = computed(() =>
  data.value?.backdropPath
    ? `https://image.tmdb.org/t/p/w1280/${data.value.backdropPath}`
    : undefined,
)

const movieChips = computed(() => {
  const chips: string[] = []
  // if (data.value?.runtimeMinutes) chips.push(`${data.value.runtimeMinutes}m`)
  // if (data.value?.certification) chips.push(data.value.certification)
  // if (data.value?.genres?.length) chips.push(...data.value.genres.slice(0, 3))
  return chips
})
</script>

<style scoped></style>
