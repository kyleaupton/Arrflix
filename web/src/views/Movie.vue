<template>
  <div class="flex flex-col gap-6">
    <div v-if="isLoading" class="space-y-4">
      <Skeleton class="h-96 w-full rounded-lg" />
    </div>
    <div v-else-if="isError" class="flex flex-col items-center justify-center py-12 text-center">
      <p class="text-destructive">Failed to load movie</p>
      <p class="text-sm text-muted-foreground mt-2">Please try again later</p>
    </div>
    <template v-else-if="data">
      <MediaHero
        class="mb-1"
        :title="data.title"
        :subtitle="releaseYear"
        :overview="data.overview"
        :backdrop-url="backdropUrl"
        :chips="movieChips"
        :trailer-url="trailerUrl"
      >
        <template #poster>
          <Poster :item="data" size="large" :clickable="false" />
        </template>
        <template #actions>
          <Button @click="searchForDownloadCandidates">
            <Download class="mr-2 size-4" />
            Snag
          </Button>
        </template>
      </MediaHero>

      <div v-if="data.files?.length" class="space-y-4">
        <h2 class="text-xl font-semibold">Local Files</h2>
        <DataTable
          :data="data.files"
          :columns="movieFilesColumns"
          :loading="false"
          empty-message="No files found"
          :searchable="false"
          search-placeholder="Search files..."
          paginator
          :rows="10"
        />
      </div>

      <!-- <Card>
        <CardHeader>
          <CardTitle>Local Files</CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable
            :data="data.files"
            :columns="movieFilesColumns"
            :loading="false"
            empty-message="No files found"
            searchable
            search-placeholder="Search files..."
            paginator
            :rows="10"
          />
        </CardContent>
      </Card> -->

      <RailCast v-if="data.credits?.cast?.length" title="Cast" :cast="data.credits.cast" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { Download } from 'lucide-vue-next'
import { getV1MovieByIdOptions } from '@/client/@tanstack/vue-query.gen'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import MediaHero from '@/components/media/MediaHero.vue'
import Poster from '@/components/poster/Poster.vue'
import RailCast from '@/components/rails/RailCast.vue'
import DataTable from '@/components/tables/DataTable.vue'
import { movieFilesColumns } from '@/components/tables/configs/movieFilesTableConfig'
import { useModal } from '@/composables/useModal'
import DownloadCandidatesDialog from '@/components/download-candidates/DownloadCandidatesDialog.vue'

const route = useRoute()
const modal = useModal()

const id = computed(() => {
  const castAttept = Number(Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)
  if (isNaN(castAttept)) {
    throw new Error('Invalid movie ID')
  }

  return castAttept
})

const trailerUrl = computed(() => {
  const trailer = data.value?.videos?.find((v) => v.isOfficialTrailer)
  if (!trailer) return undefined

  switch (trailer.site) {
    case 'YouTube':
      return `https://www.youtube.com/watch?v=${trailer.key}`
    case 'Vimeo':
      return `https://www.vimeo.com/watch?v=${trailer.key}`
    default:
      console.warn(`Unknown trailer site: ${trailer.site}`)
      return undefined
  }
})

const { isLoading, isError, data } = useQuery({
  ...getV1MovieByIdOptions({ path: { id: id.value } }),
})

const releaseYear = computed(() =>
  data.value?.releaseDate ? new Date(data.value.releaseDate).getFullYear().toString() : '',
)

const backdropUrl = computed(() =>
  data.value?.backdropPath
    ? `https://image.tmdb.org/t/p/w1280/${data.value.backdropPath}`
    : undefined,
)

const movieChips = computed(() => {
  const chips: string[] = []
  if (data.value?.genres?.length) {
    chips.push(...data.value.genres.slice(0, 4).map((g) => g.name))
  }
  return chips
})

const searchForDownloadCandidates = () => {
  modal.open(DownloadCandidatesDialog, {
    props: {
      class: 'max-w-[90vw] sm:max-w-4xl lg:max-w-6xl',
      movieId: id.value,
    },
  })
}
</script>

<style scoped></style>
