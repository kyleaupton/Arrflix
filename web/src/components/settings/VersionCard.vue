<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getV1Version, getV1Update } from '@/client/sdk.gen'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'

type BuildInfo = {
  version: string
  commit?: string
  buildDate?: string
  components?: Record<string, string>
}

type UpdateInfo = {
  status: 'up_to_date' | 'update_available' | 'unknown'
  reason?: string
  current: {
    version: string
    commit?: string
  }
  latest?: {
    version: string
    tag: string
    url: string
    publishedAt?: string
    notes?: string
    commit?: string
    ref?: string
  }
}

const isLoadingVersion = ref(true)
const isLoadingUpdate = ref(false)
const versionError = ref<string | null>(null)
const updateError = ref<string | null>(null)

const buildInfo = ref<BuildInfo | null>(null)
const updateInfo = ref<UpdateInfo | null>(null)

async function loadVersion() {
  isLoadingVersion.value = true
  versionError.value = null
  try {
    const res = await getV1Version<true>({ throwOnError: true })
    buildInfo.value = res.data as BuildInfo
  } catch (err) {
    versionError.value = 'Failed to load version information'
    console.error(err)
  } finally {
    isLoadingVersion.value = false
  }
}

async function checkUpdate() {
  isLoadingUpdate.value = true
  updateError.value = null
  try {
    const res = await getV1Update<true>({ throwOnError: true })
    updateInfo.value = res.data as UpdateInfo
  } catch (err) {
    updateError.value = 'Failed to check for updates'
    console.error(err)
  } finally {
    isLoadingUpdate.value = false
  }
}

onMounted(async () => {
  await loadVersion()
  await checkUpdate()
})

const statusBadgeVariant = computed(() => {
  if (!updateInfo.value) return 'secondary'
  switch (updateInfo.value.status) {
    case 'up_to_date':
      return 'outline'
    case 'update_available':
      return 'default'
    case 'unknown':
      return 'secondary'
    default:
      return 'secondary'
  }
})

const statusText = computed(() => {
  if (!updateInfo.value) return 'Unknown'
  switch (updateInfo.value.status) {
    case 'up_to_date':
      return 'Up to date'
    case 'update_available':
      return 'Update available'
    case 'unknown':
      return 'Unknown'
    default:
      return 'Unknown'
  }
})

const formattedBuildDate = computed(() => {
  if (!buildInfo.value?.buildDate) return null
  try {
    return new Date(buildInfo.value.buildDate).toLocaleString()
  } catch {
    return buildInfo.value.buildDate
  }
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Version</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="versionError" class="text-sm text-destructive">
        {{ versionError }}
      </div>
      <div v-else-if="isLoadingVersion" class="space-y-2">
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-3/4" />
      </div>
      <div v-else class="space-y-4">
        <!-- Current Version -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Current version</span>
            <span class="font-mono text-sm font-medium">{{ buildInfo?.version }}</span>
          </div>
          <div v-if="buildInfo?.commit" class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Commit</span>
            <span class="font-mono text-sm">{{ buildInfo.commit }}</span>
          </div>
          <div v-if="formattedBuildDate" class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Build date</span>
            <span class="text-sm">{{ formattedBuildDate }}</span>
          </div>
          <div v-if="buildInfo?.components?.prowlarr" class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Prowlarr</span>
            <span class="font-mono text-sm">{{ buildInfo.components.prowlarr }}</span>
          </div>
        </div>

        <!-- Update Status -->
        <div class="border-t pt-4">
          <div class="flex items-center justify-between mb-3">
            <span class="text-sm text-muted-foreground">Update status</span>
            <Badge :variant="statusBadgeVariant">{{ statusText }}</Badge>
          </div>

          <div v-if="updateError" class="text-sm text-destructive mb-3">
            {{ updateError }}
          </div>

          <!-- Update Available Info -->
          <div
            v-if="updateInfo?.status === 'update_available' && updateInfo.latest"
            class="space-y-2 mb-3"
          >
            <div class="flex items-center justify-between">
              <span class="text-sm text-muted-foreground">Latest version</span>
              <span class="font-mono text-sm font-medium">{{ updateInfo.latest.version }}</span>
            </div>
            <div v-if="updateInfo.latest.publishedAt" class="flex items-center justify-between">
              <span class="text-sm text-muted-foreground">Published</span>
              <span class="text-sm">{{
                new Date(updateInfo.latest.publishedAt).toLocaleDateString()
              }}</span>
            </div>
            <div v-if="updateInfo.latest.notes" class="mt-3">
              <div class="text-sm text-muted-foreground mb-2">Release notes</div>
              <div
                class="text-sm bg-muted p-3 rounded-md max-h-48 overflow-y-auto whitespace-pre-wrap"
              >
                {{ updateInfo.latest.notes }}
              </div>
            </div>
            <Button
              v-if="updateInfo.latest.url"
              as="a"
              :href="updateInfo.latest.url"
              target="_blank"
              variant="outline"
              size="sm"
              class="w-full mt-2"
            >
              View on GitHub
            </Button>
          </div>

          <!-- Unknown Status Reason -->
          <div
            v-if="updateInfo?.status === 'unknown' && updateInfo.reason"
            class="text-sm text-muted-foreground mb-3"
          >
            {{ updateInfo.reason.replace(/_/g, ' ') }}
          </div>

          <Button
            @click="checkUpdate"
            :disabled="isLoadingUpdate"
            variant="outline"
            size="sm"
            class="w-full"
          >
            {{ isLoadingUpdate ? 'Checking...' : 'Check for updates' }}
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
