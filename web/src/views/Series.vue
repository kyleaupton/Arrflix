<template>
  <div class="flex flex-col gap-6">
    <div v-if="isLoading" class="space-y-4">
      <Skeleton class="h-96 w-full rounded-lg" />
    </div>
    <div v-else-if="isError" class="flex flex-col items-center justify-center py-12 text-center">
      <p class="text-destructive">Failed to load series</p>
      <p class="text-sm text-muted-foreground mt-2">Please try again later</p>
    </div>
    <template v-else-if="data">
      <MediaHero
        class="mb-1"
        :title="data.title"
        :subtitle="seriesSubTitle"
        :overview="data.overview"
        :backdrop-url="backdropUrl"
        :chips="seriesChips"
      >
        <template #poster>
          <Poster :item="data" size="large" :clickable="false" :is-downloading="isDownloading" />
        </template>
        <template #actions>
          <Button @click="searchForSeriesCandidates">
            <Download class="mr-2 size-4" />
            Snag Series
          </Button>
        </template>
      </MediaHero>

      <div v-if="data.seasons?.length" class="space-y-4">
        <h2 class="text-xl font-semibold">Seasons</h2>
        <div class="space-y-2">
          <Collapsible
            v-for="season in sortedSeasons"
            :key="season.seasonNumber"
            v-model:open="openSeasons[season.seasonNumber]"
            class="border rounded-lg overflow-hidden"
          >
            <div
              class="flex items-center justify-between p-4 bg-muted/30 hover:bg-muted/50 transition-colors"
            >
              <CollapsibleTrigger class="flex items-center gap-4 flex-1 text-left">
                <ChevronRight
                  :class="[
                    'size-4 transition-transform',
                    openSeasons[season.seasonNumber] ? 'rotate-90' : '',
                  ]"
                />
                <div>
                  <h3 class="font-medium">Season {{ season.seasonNumber }}</h3>
                  <p v-if="season.airDate" class="text-xs text-muted-foreground">
                    {{ season.airDate }}
                  </p>
                </div>
              </CollapsibleTrigger>
              <div class="flex items-center gap-2">
                <Button
                  size="sm"
                  variant="ghost"
                  @click.stop="searchForSeasonCandidates(season.seasonNumber)"
                >
                  <Download class="size-4 mr-2" />
                  Search
                </Button>
              </div>
            </div>
            <CollapsibleContent>
              <div class="p-4 bg-background border-t space-y-4">
                <p v-if="season.overview" class="text-sm text-muted-foreground">
                  {{ season.overview }}
                </p>
                <div class="space-y-3">
                  <div
                    v-for="episode in season.episodes"
                    :key="episode.episodeNumber"
                    class="flex items-start gap-4 p-3 rounded-md hover:bg-muted/20 border border-transparent hover:border-border transition-all"
                  >
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center gap-2 mb-1">
                        <span class="text-xs font-mono text-muted-foreground">
                          E{{ episode.episodeNumber.toString().padStart(2, '0') }}
                        </span>
                        <h4 class="font-medium text-sm truncate">
                          {{ episode.title || 'Episode ' + episode.episodeNumber }}
                        </h4>
                        <Badge
                          v-if="episode.available"
                          variant="secondary"
                          class="h-5 text-[10px] px-1.5 uppercase tracking-wider"
                        >
                          {{ episode.file?.status || 'Available' }}
                        </Badge>
                      </div>
                      <p v-if="episode.overview" class="text-xs text-muted-foreground line-clamp-2">
                        {{ episode.overview }}
                      </p>
                      <p v-if="episode.airDate" class="text-[10px] text-muted-foreground/60 mt-1">
                        Aired: {{ episode.airDate }}
                      </p>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      class="h-8 text-xs shrink-0"
                      @click="
                        searchForEpisodeCandidates(season.seasonNumber, episode.episodeNumber)
                      "
                    >
                      <Search class="size-3 mr-1.5" />
                      Snag
                    </Button>
                  </div>
                </div>
              </div>
            </CollapsibleContent>
          </Collapsible>
        </div>
      </div>

      <RailCast v-if="data.credits?.cast?.length" title="Cast" :cast="data.credits.cast" />
      <RailVideos v-if="data.videos?.length" title="Videos" :videos="data.videos" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { Download, ChevronRight, Search } from 'lucide-vue-next'
import { getV1SeriesByIdOptions } from '@/client/@tanstack/vue-query.gen'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import MediaHero from '@/components/media/MediaHero.vue'
import Poster from '@/components/poster/Poster.vue'
import RailCast from '@/components/rails/RailCast.vue'
import RailVideos from '@/components/rails/RailVideos.vue'
import { useModal } from '@/composables/useModal'
import { useDownloadJobsStore } from '@/stores/downloadJobs'
import DownloadCandidatesDialog from '@/components/download-candidates/DownloadCandidatesDialog.vue'

const route = useRoute()
const modal = useModal()
const downloadJobs = useDownloadJobsStore()

const openSeasons = ref<Record<number, boolean>>({})

const id = computed(() => {
  const castAttept = Number(Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)
  if (isNaN(castAttept)) {
    throw new Error('Invalid series ID')
  }

  return castAttept
})

const { isLoading, isError, data } = useQuery(
  computed(() => getV1SeriesByIdOptions({ path: { id: id.value } })),
)

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

const backdropUrl = computed(() =>
  data.value?.backdropPath
    ? `https://image.tmdb.org/t/p/w1280/${data.value.backdropPath}`
    : undefined,
)

const seriesChips = computed(() => {
  const chips: string[] = []
  if (data.value?.genres?.length) {
    chips.push(...data.value.genres.slice(0, 3).map((g) => g.name))
  }
  if (data.value?.status) {
    chips.push(data.value.status)
  }
  return chips
})

const sortedSeasons = computed(() => {
  if (!data.value?.seasons) return []
  return [...data.value.seasons].sort((a, b) => b.seasonNumber - a.seasonNumber)
})

const isDownloading = computed(() => {
  // Simple check for now: if any episode is in downloading status
  return (
    data.value?.seasons?.some((s) => s.episodes?.some((e) => e.file?.status === 'downloading')) ??
    false
  )
})

onMounted(() => {
  downloadJobs.connectLive()
})

const searchForSeriesCandidates = () => {
  modal.open(DownloadCandidatesDialog, {
    props: {
      class: 'max-w-[90vw] sm:max-w-4xl lg:max-w-6xl',
      seriesId: id.value,
    },
  })
}

const searchForSeasonCandidates = (seasonNumber: number) => {
  modal.open(DownloadCandidatesDialog, {
    props: {
      class: 'max-w-[90vw] sm:max-w-4xl lg:max-w-6xl',
      seriesId: id.value,
      season: seasonNumber,
    },
  })
}

const searchForEpisodeCandidates = (seasonNumber: number, episodeNumber: number) => {
  modal.open(DownloadCandidatesDialog, {
    props: {
      class: 'max-w-[90vw] sm:max-w-4xl lg:max-w-6xl',
      seriesId: id.value,
      season: seasonNumber,
      episode: episodeNumber,
    },
  })
}
</script>

<style scoped></style>
