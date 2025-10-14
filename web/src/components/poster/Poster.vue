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

const props = withDefaults(
  defineProps<{
    item: ModelMovie | ModelSeries | ModelMovieRail | ModelSeriesRail
    size?: 'small' | 'medium' | 'large'
  }>(),
  {
    size: 'medium',
  },
)

const posterPath = computed(() => {
  return `https://image.tmdb.org/t/p/w342/${props.item.posterPath}`
})

const sizeClass = computed(() => {
  if (props.size === 'small') return 'poster--sm'
  if (props.size === 'large') return 'poster--lg'
  return 'poster--md'
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
  max-width: 154px; /* roughly TMDB w154 */
}

.poster--md {
  max-width: 342px; /* roughly TMDB w342 */
}

.poster--lg {
  max-width: 500px; /* roughly TMDB w500 */
}
</style>
