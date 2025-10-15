<template>
  <div class="rail">
    <div class="rail-header flex items-center justify-between mb-2">
      <h1 class="text-xl font-semibold">{{ rail.title }}</h1>
      <div class="flex items-center gap-2">
        <Button
          :icon="PrimeIcons.CHEVRON_LEFT"
          :disabled="!canScrollPrev"
          severity="secondary"
          variant="outlined"
          size="small"
          rounded
          aria-label="Scroll left"
          @click="scrollByPage(-1)"
        />
        <Button
          :icon="PrimeIcons.CHEVRON_RIGHT"
          :disabled="!canScrollNext"
          severity="secondary"
          variant="outlined"
          size="small"
          rounded
          aria-label="Scroll right"
          @click="scrollByPage(1)"
        />
      </div>
    </div>

    <div class="rail-body relative">
      <div
        ref="scroller"
        class="scroller flex gap-3 overflow-x-auto overflow-y-hidden snap-x snap-mandatory scroll-smooth"
        tabindex="0"
        @wheel="onWheel"
        @keydown="onKeydown"
        @scroll="onScroll"
      >
        <template v-if="rail.type === 'movie'">
          <div v-for="movie in rail.movies" :key="movie.tmdbId" class="snap-start flex-shrink-0">
            <Poster :item="movie" :to="{ path: `/movie/${movie.tmdbId}` }" />
          </div>
        </template>

        <template v-else-if="rail.type === 'series'">
          <div v-for="series in rail.series" :key="series.tmdbId" class="snap-start flex-shrink-0">
            <Poster :item="series" :to="{ path: `/series/${series.tmdbId}` }" />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import { type ModelRail } from '@/client/types.gen'
import { PrimeIcons } from '@/icons'
import Poster from '@/components/poster/Poster.vue'
import { ref, onMounted, onBeforeUnmount } from 'vue'

defineProps<{
  rail: ModelRail
}>()

const scroller = ref<HTMLDivElement | null>(null)
const canScrollPrev = ref(false)
const canScrollNext = ref(false)
let resizeObserver: ResizeObserver | null = null
let rafId: number | null = null

const updateScrollState = () => {
  const el = scroller.value
  if (!el) return
  const maxScrollLeft = el.scrollWidth - el.clientWidth - 1
  canScrollPrev.value = el.scrollLeft > 0
  canScrollNext.value = el.scrollLeft < maxScrollLeft
}

const onScroll = () => {
  if (rafId != null) cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(updateScrollState)
}

const onWheel = (e: WheelEvent) => {
  // Only handle horizontal scroll if user is explicitly scrolling horizontally
  // Don't interfere with vertical page scrolling
  const el = scroller.value
  if (!el) return

  // Only prevent default and scroll horizontally if the user is actually scrolling horizontally
  if (Math.abs(e.deltaX) > Math.abs(e.deltaY)) {
    e.preventDefault()
    el.scrollLeft += e.deltaX
  }
  // Let vertical scrolling pass through normally
}

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'ArrowLeft') {
    e.preventDefault()
    scrollByPage(-1)
  } else if (e.key === 'ArrowRight') {
    e.preventDefault()
    scrollByPage(1)
  }
}

const scrollByPage = (direction: number) => {
  const el = scroller.value
  if (!el) return
  const page = Math.max(1, el.clientWidth - 64)
  el.scrollBy({ left: direction * page, behavior: 'smooth' })
}

onMounted(() => {
  updateScrollState()
  resizeObserver = new ResizeObserver(updateScrollState)
  if (scroller.value) {
    resizeObserver.observe(scroller.value)
  }
})

onBeforeUnmount(() => {
  if (resizeObserver && scroller.value) {
    resizeObserver.unobserve(scroller.value)
    resizeObserver.disconnect()
    resizeObserver = null
  }
})
</script>

<style scoped>
.rail-body {
  position: relative;
}

.scroller {
  scrollbar-width: none; /* Firefox */
}

.scroller::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}
</style>
