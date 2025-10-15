<template>
  <div class="poster-wrap" :class="sizeClass">
    <img
      class="poster"
      :class="{ 'is-loaded': !isLoading }"
      :src="posterPath"
      :alt="item.title"
      loading="lazy"
      decoding="async"
      @load="onLoad"
      @error="onError"
    />
    <Skeleton v-if="isLoading" class="poster-skeleton" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import Skeleton from 'primevue/skeleton'
import {
  type ModelMovie,
  type ModelMovieRail,
  type ModelSeries,
  type ModelSeriesRail,
} from '@/client/types.gen'

// Poster size configuration
const POSTER_SIZES = {
  small: {
    width: '8rem', // ~128px
    class: 'poster--sm',
    tmdbSize: 'w154',
  },
  medium: {
    width: '11rem', // ~176px
    class: 'poster--md',
    tmdbSize: 'w185',
  },
  large: {
    width: '16rem', // ~256px
    class: 'poster--lg',
    tmdbSize: 'w342',
  },
} as const

type PosterSize = keyof typeof POSTER_SIZES

const props = withDefaults(
  defineProps<{
    item: ModelMovie | ModelSeries | ModelMovieRail | ModelSeriesRail
    size?: PosterSize
  }>(),
  {
    size: 'medium',
  },
)

const posterPath = computed(() => {
  const sizeConfig = POSTER_SIZES[props.size]
  return `https://image.tmdb.org/t/p/${sizeConfig.tmdbSize}/${props.item.posterPath}`
})

const sizeClass = computed(() => {
  return POSTER_SIZES[props.size].class
})

const isLoading = ref(true)
const onLoad = () => {
  isLoading.value = false
}
const onError = () => {
  isLoading.value = false
}
</script>

<style scoped></style>
<style scoped>
.poster-wrap {
  display: block;
  width: 100%;
  aspect-ratio: 2 / 3; /* common movie/TV poster ratio */
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  background-color: #111827; /* neutral placeholder while loading */
}

.poster {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity 150ms ease;
}

.poster.is-loaded {
  opacity: 1;
}

.poster-skeleton {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.poster--sm {
  --poster-width: 8rem; /* ~128px */
  width: var(--poster-width);
  max-width: var(--poster-width);
}

.poster--md {
  --poster-width: 11rem; /* ~176px */
  width: var(--poster-width);
  max-width: var(--poster-width);
}

.poster--lg {
  --poster-width: 16rem; /* ~256px */
  width: var(--poster-width);
  max-width: var(--poster-width);
}
</style>
