<template>
  <component :is="to ? 'router-link' : 'div'" :to="to" class="poster-wrap" :class="sizeClass">
    <img
      class="poster"
      :class="{ 'is-loaded': !isLoading, 'cursor-pointer': clickable }"
      :src="posterPath"
      :alt="item.title"
      loading="lazy"
      decoding="async"
      @load="onLoad"
      @error="onError"
    />
    <Skeleton v-if="isLoading" class="poster-skeleton" />
    <!-- Show download status badge when downloading, otherwise show library badge -->
    <div v-if="isDownloading" class="status-badge download-badge">
      <Loader2 class="size-4 animate-spin" aria-hidden="true" />
      <span>Downloading</span>
    </div>
    <div v-else-if="showLibraryBadge" class="status-badge library-badge">
      <CheckCircle2 class="size-4" aria-hidden="true" />
      <span>In library</span>
    </div>
  </component>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { CheckCircle2, Loader2 } from 'lucide-vue-next'
import { Skeleton } from '@/components/ui/skeleton'
import {
  type ModelLibraryItem,
  type ModelMovieDetail,
  type ModelMovieRail,
  type ModelSeriesDetail,
  type ModelSeriesRail,
} from '@/client/types.gen'

/**
 * "poster_sizes": [
      "w92",
      "w154",
      "w185",
      "w342",
      "w500",
      "w780",
      "original"
    ],
 */

// Poster size configuration
const POSTER_SIZES = {
  small: {
    width: '8rem', // ~128px
    class: 'poster--sm',
    tmdbSize: 'w185',
  },
  medium: {
    width: '11rem', // ~176px
    class: 'poster--md',
    tmdbSize: 'w500',
  },
  large: {
    width: '16rem', // ~256px
    class: 'poster--lg',
    tmdbSize: 'w500',
  },
} as const

type PosterSize = keyof typeof POSTER_SIZES

const props = withDefaults(
  defineProps<{
    item: ModelMovieDetail | ModelSeriesDetail | ModelMovieRail | ModelSeriesRail | ModelLibraryItem
    size?: PosterSize
    to?: { path: string } | string
    clickable?: boolean
    isDownloading?: boolean
    responsive?: boolean
  }>(),
  {
    size: 'medium',
    clickable: true,
    isDownloading: false,
    responsive: false,
  },
)

const isInLibrary = computed(() => {
  if ('files' in props.item && props.item.files) {
    return props.item.files.some((file) => file.status === 'available')
  } else if ('isInLibrary' in props.item) {
    return props.item.isInLibrary
  }

  return false
})

const showLibraryBadge = computed(() => !props.isDownloading && isInLibrary.value)

const posterPath = computed(() => {
  const sizeConfig = POSTER_SIZES[props.size]
  return `https://image.tmdb.org/t/p/${sizeConfig.tmdbSize}/${props.item.posterPath}`
})

const sizeClass = computed(() => {
  if (props.responsive) return ''
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

<style scoped>
.poster-wrap {
  display: block;
  width: 100%;
  aspect-ratio: 2 / 3; /* common movie/TV poster ratio */
  position: relative;
  border-radius: 8px !important;
  overflow: hidden;
  background-color: #111827; /* neutral placeholder while loading */
  text-decoration: none;
  color: inherit;
}

.poster {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity 150ms ease;
  border-radius: 8px !important;
}

.poster.is-loaded {
  opacity: 1;
}

.status-badge {
  position: absolute;
  top: 0.5rem;
  left: 0.5rem;
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.3rem 0.55rem;
  border-radius: 9999px;
  background: rgba(0, 0, 0, 0.7);
  color: #e5e7eb;
  font-weight: 600;
  font-size: 0.75rem;
  letter-spacing: 0.01em;
  border: 1px solid rgba(255, 255, 255, 0.15);
  z-index: 5;
}

.status-badge svg {
  flex-shrink: 0;
}

.library-badge svg {
  color: #22c55e; /* emerald-500 */
}

.download-badge svg {
  color: #3b82f6; /* blue-500 */
}

.poster-skeleton {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  border-radius: 8px !important;
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
