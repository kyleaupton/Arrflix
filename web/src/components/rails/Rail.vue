<template>
  <div class="rail">
    <div class="rail-header flex items-center justify-between mb-2">
      <h1 class="text-xl font-semibold">{{ rail.title }}</h1>
    </div>

    <div class="rail-body relative group">
      <button
        v-if="canScrollPrev"
        class="nav-btn left"
        type="button"
        aria-label="Scroll left"
        @click="scrollByPage(-1)"
      >
        <span aria-hidden="true">‹</span>
      </button>

      <div
        ref="scroller"
        class="scroller flex gap-3 overflow-x-auto overflow-y-hidden snap-x snap-mandatory scroll-smooth"
        tabindex="0"
        @wheel="onWheel"
        @keydown="onKeydown"
        @scroll="onScroll"
        @pointerdown="onPointerDown"
      >
        <template v-if="rail.type === 'movie'">
          <div v-for="movie in rail.movies" :key="movie.tmdbId" class="snap-start">
            <Poster :item="movie" />
          </div>
        </template>

        <template v-else-if="rail.type === 'series'">
          <div v-for="series in rail.series" :key="series.tmdbId" class="snap-start">
            <Poster :item="series" />
          </div>
        </template>
      </div>

      <button
        v-if="canScrollNext"
        class="nav-btn right"
        type="button"
        aria-label="Scroll right"
        @click="scrollByPage(1)"
      >
        <span aria-hidden="true">›</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { type ModelRail } from '@/client/types.gen'
import Poster from '@/components/poster/Poster.vue'
import { ref, onMounted, onBeforeUnmount } from 'vue'

defineProps<{
  rail: ModelRail
}>()

const scroller = ref<HTMLDivElement | null>(null)
const canScrollPrev = ref(false)
const canScrollNext = ref(false)
let resizeObserver: ResizeObserver | null = null

let isPointerDown = false
let pointerStartX = 0
let scrollStartLeft = 0
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
  const el = scroller.value
  if (!el) return
  if (Math.abs(e.deltaY) > Math.abs(e.deltaX)) {
    e.preventDefault()
    el.scrollLeft += e.deltaY
  }
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

const onPointerDown = (e: PointerEvent) => {
  const el = scroller.value
  if (!el) return
  isPointerDown = true
  pointerStartX = e.clientX
  scrollStartLeft = el.scrollLeft
  el.setPointerCapture(e.pointerId)
  el.classList.add('is-dragging')
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp, { once: true })
}

const onPointerMove = (e: PointerEvent) => {
  if (!isPointerDown) return
  const el = scroller.value
  if (!el) return
  const dx = e.clientX - pointerStartX
  el.scrollLeft = scrollStartLeft - dx
}

const onPointerUp = (e: PointerEvent) => {
  const el = scroller.value
  isPointerDown = false
  if (el) {
    el.releasePointerCapture(e.pointerId)
    el.classList.remove('is-dragging')
  }
  window.removeEventListener('pointermove', onPointerMove)
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
  window.removeEventListener('pointermove', onPointerMove)
})
</script>

<style scoped>
.rail-body {
  position: relative;
}

.scroller {
  scrollbar-width: none; /* Firefox */
  cursor: grab;
}

.scroller::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.scroller.is-dragging {
  cursor: grabbing;
}

.nav-btn {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 48px;
  display: grid;
  place-items: center;
  color: white;
  background: linear-gradient(to right, rgba(0, 0, 0, 0.6), transparent);
  border: none;
  opacity: 0;
  z-index: 2;
  transition: opacity 150ms ease;
}

.nav-btn.left {
  left: 0;
}

.nav-btn.right {
  right: 0;
  background: linear-gradient(to left, rgba(0, 0, 0, 0.6), transparent);
}

.rail-body:hover .nav-btn,
.rail-body:focus-within .nav-btn {
  opacity: 1;
}

.nav-btn:disabled {
  opacity: 0;
  pointer-events: none;
}
</style>
