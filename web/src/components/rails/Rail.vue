<template>
  <div class="rail space-y-2">
    <div class="rail-header flex items-center justify-between">
      <h2 class="text-xl font-semibold">{{ rail.title }}</h2>
      <div class="flex items-center gap-2">
        <Button
          variant="outline"
          size="icon-sm"
          :disabled="!canScrollPrev"
          aria-label="Scroll left"
          @click="scrollByPage(-1)"
        >
          <ChevronLeft class="size-4" />
        </Button>
        <Button
          variant="outline"
          size="icon-sm"
          :disabled="!canScrollNext"
          aria-label="Scroll right"
          @click="scrollByPage(1)"
        >
          <ChevronRight class="size-4" />
        </Button>
      </div>
    </div>

    <div class="rail-body relative">
      <div
        ref="scroller"
        class="scroller flex gap-3 overflow-x-auto overflow-y-hidden pb-4"
        @scroll="onScroll"
      >
        <template v-if="rail.type === 'movie'">
          <div v-for="movie in rail.movies" :key="movie.tmdbId" class="flex-shrink-0">
            <Poster :item="movie" :to="{ path: `/movie/${movie.tmdbId}` }" />
          </div>
        </template>

        <template v-else-if="rail.type === 'series'">
          <div v-for="series in rail.series" :key="series.tmdbId" class="flex-shrink-0">
            <Poster :item="series" :to="{ path: `/series/${series.tmdbId}` }" />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { type ModelRail } from '@/client/types.gen'
import { Button } from '@/components/ui/button'
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
.scroller {
  scrollbar-width: none; /* Firefox */
}

.scroller::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}
</style>
