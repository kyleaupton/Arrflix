<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getV1Version } from '@/client/sdk.gen'
import type { ServiceVersionInfo } from '@/client/types.gen'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { CheckCircle2, ArrowUpCircle, ExternalLink, RefreshCw } from 'lucide-vue-next'

const isLoading = ref(true)
const isRefreshing = ref(false)
const error = ref<string | null>(null)
const data = ref<ServiceVersionInfo | null>(null)

async function load(refreshing = false) {
  if (refreshing) {
    isRefreshing.value = true
  } else {
    isLoading.value = true
  }
  error.value = null
  try {
    const res = await getV1Version<true>({ throwOnError: true })
    data.value = res.data as ServiceVersionInfo
  } catch {
    error.value = 'Failed to load version information'
  } finally {
    isLoading.value = false
    isRefreshing.value = false
  }
}

onMounted(() => load())

const isDev = computed(() => {
  return data.value?.update.status === 'unknown' && data.value?.update.reason === 'dev_build'
})

const isUpdateAvailable = computed(() => {
  return data.value?.update.status === 'update_available'
})

const isUpToDate = computed(() => {
  return data.value?.update.status === 'up_to_date'
})

const formattedBuildDate = computed(() => {
  if (!data.value?.buildDate) return null
  try {
    return new Date(data.value.buildDate).toLocaleDateString()
  } catch {
    return data.value.buildDate
  }
})

const secondaryInfo = computed(() => {
  const parts: string[] = []
  if (data.value?.commit) parts.push(data.value.commit)
  if (formattedBuildDate.value) parts.push(formattedBuildDate.value)
  return parts.join(' \u00b7 ')
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Version</CardTitle>
    </CardHeader>
    <CardContent>
      <!-- Loading -->
      <div v-if="isLoading" class="space-y-3">
        <Skeleton class="h-8 w-24" />
        <Skeleton class="h-4 w-48" />
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-sm text-destructive">
        {{ error }}
      </div>

      <!-- Loaded -->
      <div v-else-if="data" class="space-y-4">
        <!-- Hero version -->
        <div>
          <div class="text-2xl font-semibold font-mono tracking-tight">
            {{ data.version }}
          </div>
          <div v-if="secondaryInfo" class="text-xs text-muted-foreground font-mono mt-1">
            {{ secondaryInfo }}
          </div>
        </div>

        <!-- Up to date -->
        <div v-if="isUpToDate" class="flex items-center justify-between">
          <div class="flex items-center gap-1.5 text-sm text-emerald-600 dark:text-emerald-400">
            <CheckCircle2 class="size-4" />
            <span>Up to date</span>
          </div>
          <Button
            variant="ghost"
            size="sm"
            class="h-7 px-2 text-xs text-muted-foreground"
            :disabled="isRefreshing"
            @click="load(true)"
          >
            <RefreshCw class="size-3 mr-1" :class="{ 'animate-spin': isRefreshing }" />
            Check again
          </Button>
        </div>

        <!-- Update available -->
        <div v-else-if="isUpdateAvailable && data.update.latest" class="space-y-3">
          <div
            class="rounded-lg border border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/50 p-3 space-y-2"
          >
            <div class="flex items-center gap-1.5">
              <ArrowUpCircle class="size-4 text-blue-600 dark:text-blue-400 shrink-0" />
              <span class="text-sm font-medium text-blue-900 dark:text-blue-100">
                {{ data.update.latest.version }} available
              </span>
            </div>
            <div
              v-if="data.update.latest.notes"
              class="text-xs text-blue-800 dark:text-blue-200/80 max-h-32 overflow-y-auto whitespace-pre-wrap pl-[22px]"
            >
              {{ data.update.latest.notes }}
            </div>
            <div class="flex items-center gap-2 pl-[22px]">
              <Button
                v-if="data.update.latest.url"
                as="a"
                :href="data.update.latest.url"
                target="_blank"
                variant="outline"
                size="sm"
                class="h-7 px-2 text-xs"
              >
                View on GitHub
                <ExternalLink class="size-3 ml-1" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                class="h-7 px-2 text-xs text-muted-foreground"
                :disabled="isRefreshing"
                @click="load(true)"
              >
                <RefreshCw class="size-3 mr-1" :class="{ 'animate-spin': isRefreshing }" />
                Refresh
              </Button>
            </div>
          </div>
        </div>

        <!-- Dev build -->
        <div v-else-if="isDev" class="text-xs text-muted-foreground">
          Development build
        </div>

        <!-- Other unknown -->
        <div v-else class="flex items-center justify-between">
          <span class="text-xs text-muted-foreground"> Unable to determine update status </span>
          <Button
            variant="ghost"
            size="sm"
            class="h-7 px-2 text-xs text-muted-foreground"
            :disabled="isRefreshing"
            @click="load(true)"
          >
            <RefreshCw class="size-3 mr-1" :class="{ 'animate-spin': isRefreshing }" />
            Check again
          </Button>
        </div>

        <!-- Prowlarr footnote -->
        <div
          v-if="data.components?.prowlarr"
          class="text-xs text-muted-foreground border-t pt-3"
        >
          Prowlarr {{ data.components.prowlarr }}
        </div>
      </div>
    </CardContent>
  </Card>
</template>
