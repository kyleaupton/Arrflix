<template>
  <RouterLink v-if="to" :to="to" class="poster-link" @click="test">
    <div
      class="poster-wrap"
      :class="{
        'poster--sm w-36': size === 'small',
        'poster--md w-48': size === 'medium',
        'poster--lg w-60': size === 'large',
      }"
    >
      <img
        class="poster"
        :class="{ 'is-loaded': !isLoading }"
        :src="posterPath"
        :alt="item.title"
        draggable="false"
        loading="lazy"
        decoding="async"
        @dragstart.prevent
        @load="onLoad"
        @error="onError"
      />
      <Skeleton v-if="isLoading" class="poster-skeleton" />
    </div>
  </RouterLink>
  <div
    v-else
    class="poster-wrap"
    :class="{
      'poster--sm w-36': size === 'small',
      'poster--md w-48': size === 'medium',
      'poster--lg w-60': size === 'large',
    }"
  >
    <img
      class="poster"
      :class="{ 'is-loaded': !isLoading }"
      :src="posterPath"
      :alt="item.title"
      draggable="false"
      loading="lazy"
      decoding="async"
      @dragstart.prevent
      @load="onLoad"
      @error="onError"
    />
    <Skeleton v-if="isLoading" class="poster-skeleton" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import Skeleton from 'primevue/skeleton'
import { type ModelMovie, type ModelSeries } from '@/client/types.gen'
import type { RouteLocationRaw } from 'vue-router'

const test = () => {
  console.log('test')
}

const props = withDefaults(
  defineProps<{
    item: ModelMovie | ModelSeries
    size?: 'small' | 'medium' | 'large'
    to?: RouteLocationRaw
  }>(),
  {
    size: 'medium',
  },
)

const posterPath = computed(() => {
  return `https://image.tmdb.org/t/p/${tmdbSize.value}/${props.item.posterPath}`
})

const isLoading = ref(true)
const onLoad = () => {
  isLoading.value = false
}
const onError = () => {
  isLoading.value = false
}

const tmdbSize = computed(() => {
  if (props.size === 'small') return 'w154'
  if (props.size === 'large') return 'w500'
  return 'w342'
})
</script>

<style scoped></style>
<style scoped>
.poster-wrap {
  display: inline-block;
  flex: 0 0 auto;
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
  -webkit-user-drag: none; /* Safari */
  user-select: none;
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
</style>
