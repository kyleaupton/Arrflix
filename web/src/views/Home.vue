<template>
  <div class="flex flex-col gap-6">
    <div v-if="isLoading" class="space-y-6">
      <div class="space-y-2">
        <Skeleton class="h-8 w-64" />
        <Skeleton class="h-4 w-48" />
      </div>
      <div class="flex gap-3 overflow-x-auto pb-4">
        <Skeleton v-for="i in 5" :key="i" class="h-64 w-44 flex-shrink-0" />
      </div>
    </div>

    <div v-else-if="isError" class="flex flex-col items-center justify-center py-12 text-center">
      <p class="text-destructive">Failed to load content</p>
      <p class="text-sm text-muted-foreground mt-2">Please try again later</p>
    </div>

    <div
      v-else-if="!data || data.length === 0"
      class="flex flex-col items-center justify-center py-12 text-center"
    >
      <p class="text-muted-foreground">No content available</p>
    </div>

    <div v-else class="flex flex-col gap-6">
      <Rail v-for="rail in data" :key="rail.id" :rail="rail" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { getV1HomeOptions } from '@/client/@tanstack/vue-query.gen'
import { Skeleton } from '@/components/ui/skeleton'
import Rail from '@/components/rails/Rail.vue'

const { isLoading, isError, data } = useQuery(getV1HomeOptions())
</script>
