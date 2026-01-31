<script setup lang="ts">
import { computed } from 'vue'
import { Check, X, Minus } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

export type CircularProgressState = 'progress' | 'indeterminate' | 'success' | 'error' | 'cancelled'
export type CircularProgressSize = 'sm' | 'md' | 'lg'

const props = withDefaults(
  defineProps<{
    /** Progress value 0-100. Only used when state is 'progress'. */
    value?: number
    /** Visual state of the component */
    state?: CircularProgressState
    /** Size variant */
    size?: CircularProgressSize
    /** Optional CSS class */
    class?: string
  }>(),
  {
    value: 0,
    state: 'progress',
    size: 'md',
  },
)

const sizes = {
  sm: { dim: 24, stroke: 2.5, icon: 12, radius: 10 },
  md: { dim: 32, stroke: 3, icon: 16, radius: 13 },
  lg: { dim: 40, stroke: 3.5, icon: 20, radius: 16 },
}

const config = computed(() => sizes[props.size])
const viewBox = computed(() => `0 0 ${config.value.dim} ${config.value.dim}`)
const center = computed(() => config.value.dim / 2)
const circumference = computed(() => 2 * Math.PI * config.value.radius)

const offset = computed(() => {
  if (props.state !== 'progress') return 0
  const clampedValue = Math.min(100, Math.max(0, props.value))
  return circumference.value * (1 - clampedValue / 100)
})

const ringClass = computed(() => {
  switch (props.state) {
    case 'success':
      return 'text-green-500'
    case 'error':
      return 'text-destructive'
    case 'cancelled':
      return 'text-muted-foreground'
    default:
      return 'text-primary'
  }
})

const showFullRing = computed(() => ['success', 'error', 'cancelled'].includes(props.state))

const isIndeterminate = computed(() => props.state === 'indeterminate')

// For indeterminate state, show a partial arc that spins
const indeterminateDashArray = computed(() => {
  const arcLength = circumference.value * 0.25 // 25% of the circle
  return `${arcLength} ${circumference.value - arcLength}`
})
</script>

<template>
  <svg
    :viewBox="viewBox"
    :width="config.dim"
    :height="config.dim"
    :class="cn('shrink-0', props.class)"
  >
    <!-- Background track -->
    <circle
      :cx="center"
      :cy="center"
      :r="config.radius"
      fill="none"
      :stroke-width="config.stroke"
      class="stroke-muted-foreground/20"
    />

    <!-- Progress/state ring -->
    <circle
      :cx="center"
      :cy="center"
      :r="config.radius"
      fill="none"
      :stroke-width="config.stroke"
      stroke-linecap="round"
      :class="['stroke-current', ringClass, isIndeterminate && 'animate-spin-progress']"
      :style="!isIndeterminate ? { transition: 'stroke-dashoffset 0.3s ease-in-out' } : undefined"
      :stroke-dasharray="isIndeterminate ? indeterminateDashArray : (showFullRing ? undefined : circumference)"
      :stroke-dashoffset="isIndeterminate ? 0 : (showFullRing ? 0 : offset)"
      :transform="`rotate(-90 ${center} ${center})`"
      :transform-origin="`${center}px ${center}px`"
    />

    <!-- Center icon for terminal states -->
    <foreignObject
      v-if="state === 'success' || state === 'error' || state === 'cancelled'"
      :x="center - config.icon / 2"
      :y="center - config.icon / 2"
      :width="config.icon"
      :height="config.icon"
    >
      <div class="flex items-center justify-center w-full h-full">
        <Check v-if="state === 'success'" :size="config.icon" class="text-green-500" />
        <X v-else-if="state === 'error'" :size="config.icon" class="text-destructive" />
        <Minus v-else-if="state === 'cancelled'" :size="config.icon" class="text-muted-foreground" />
      </div>
    </foreignObject>
  </svg>
</template>

<style scoped>
@keyframes spin-progress {
  from {
    transform: rotate(-90deg);
  }
  to {
    transform: rotate(270deg);
  }
}

.animate-spin-progress {
  animation: spin-progress 1.5s linear infinite;
}
</style>
