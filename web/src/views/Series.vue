<template>
  <Page>
    <div v-if="isLoading">Loading...</div>
    <div v-else-if="isError">Error</div>
    <div v-else-if="data">
      <MediaHero
        :title="data.title"
        :subtitle="seriesSubTitle"
        :overview="data.overview"
        :poster-url="posterUrl"
        :backdrop-url="backdropUrl"
        :chips="seriesChips"
      />

      <!-- TODO: seasons list, cast, recommendations, similar -->
    </div>
  </Page>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { getV1SeriesByIdOptions } from '@/client/@tanstack/vue-query.gen'
import Page from '@/components/Page.vue'
import MediaHero from '@/components/media/MediaHero.vue'

const route = useRoute()

const id = computed(() => {
  const castAttept = Number(Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)
  if (isNaN(castAttept)) {
    throw new Error('Invalid series ID')
  }

  return castAttept
})

const { isLoading, isError, data } = useQuery({
  ...getV1SeriesByIdOptions({ path: { id: id.value } }),
})

const firstAirYear = computed(() =>
  data.value?.firstAirDate ? new Date(data.value.firstAirDate).getFullYear().toString() : '',
)
const lastAirYear = computed(() =>
  data.value?.lastAirDate ? new Date(data.value.lastAirDate).getFullYear().toString() : '',
)
const seriesSubTitle = computed(() => {
  const first = firstAirYear.value
  const last = lastAirYear.value
  if (first && last && first !== last) return `${first} - ${last}`
  return first
})

const posterUrl = computed(() =>
  data.value?.posterPath ? `https://image.tmdb.org/t/p/w342/${data.value.posterPath}` : undefined,
)
const backdropUrl = computed(() =>
  data.value?.backdropPath
    ? `https://image.tmdb.org/t/p/w1280/${data.value.backdropPath}`
    : undefined,
)

const seriesChips = computed(() => {
  const chips: string[] = []
  // if (data.value?.numberOfSeasons) chips.push(`${data.value.numberOfSeasons} seasons`)
  // if (data.value?.status) chips.push(data.value.status)
  // if (data.value?.genres?.length) chips.push(...data.value.genres.slice(0, 3))
  return chips
})
</script>

<style scoped></style>
