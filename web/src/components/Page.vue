<template>
  <div class="flex flex-col gap-3 p-4 md:p-6 h-full">
    <header v-if="hasHeader" class="flex items-start justify-between gap-3">
      <slot name="header">
        <div class="flex flex-col gap-1 min-w-0">
          <div class="flex items-center gap-2 min-w-0">
            <slot name="icon">
              <span v-if="icon" :class="[icon, 'text-primary']" />
            </slot>
            <slot name="title">
              <h1 v-if="title" class="text-2xl font-semibold truncate">{{ title }}</h1>
            </slot>
          </div>
          <slot name="subtitle">
            <p v-if="subtitle" class="text-muted-color truncate">{{ subtitle }}</p>
          </slot>
        </div>
      </slot>

      <div class="flex items-center gap-2 shrink-0">
        <slot name="actions" />
      </div>
    </header>

    <div v-if="$slots.toolbar" class="rounded border border-surface bg-emphasis px-3 py-2">
      <slot name="toolbar" />
    </div>

    <div class="flex-1">
      <div v-if="loading" class="rounded bg-emphasis p-6 animate-pulse min-h-32" />
      <template v-else>
        <div v-if="isEmpty">
          <slot name="empty">
            <div class="text-center text-muted-color py-12">No data to display</div>
          </slot>
        </div>
        <div v-else :class="contentClass">
          <slot />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

interface PageProps {
  title?: string
  subtitle?: string
  icon?: string
  loading?: boolean
  isEmpty?: boolean
  contentClass?: string
  error?: string
}

const props = defineProps<PageProps>()
const slots = useSlots()

const hasHeader = computed(() => {
  return (
    !!props.title ||
    !!props.subtitle ||
    !!props.icon ||
    !!slots.header ||
    !!slots.title ||
    !!slots.subtitle ||
    !!slots.actions ||
    !!slots.icon
  )
})

const {
  title,
  subtitle,
  icon,
  loading = false,
  isEmpty = false,
  contentClass = 'space-y-4',
} = props
</script>

<style scoped></style>
