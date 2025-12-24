<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getV1Settings, patchV1Settings } from '@/client/sdk.gen'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { Skeleton } from '@/components/ui/skeleton'

type SettingsMap = Record<string, unknown>

const isLoading = ref(true)
const error = ref<string | null>(null)
const settings = ref<SettingsMap>({})
const isSaving = ref(false)

async function loadSettings() {
  isLoading.value = true
  error.value = null
  try {
    const res = await getV1Settings<true>({ throwOnError: true })
    settings.value = res.data as SettingsMap
  } catch {
    error.value = 'Failed to load settings'
  } finally {
    isLoading.value = false
  }
}

async function saveSetting(key: string, value: unknown) {
  isSaving.value = true
  try {
    await patchV1Settings<true>({ throwOnError: true, body: { key, value } })
    // Optimistically update local state
    settings.value = { ...settings.value, [key]: value }
  } finally {
    isSaving.value = false
  }
}

onMounted(loadSettings)

const siteTitle = computed({
  get: () => String(settings.value['site.title'] ?? ''),
  set: (v: string) => saveSetting('site.title', v),
})

const allowSignups = computed({
  get: () => Boolean(settings.value['auth.allow_signups'] ?? false),
  set: (v: boolean) => saveSetting('auth.allow_signups', v),
})

const maxPerUser = computed({
  get: () => Number(settings.value['requests.max_per_user'] ?? 0),
  set: (v: number) => saveSetting('requests.max_per_user', v),
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold">General Settings</h1>
    </div>
    <div
      v-if="error"
      class="p-4 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive"
    >
      {{ error }}
    </div>
    <div v-if="isLoading" class="space-y-3">
      <Skeleton class="h-24 w-full" />
      <Skeleton class="h-24 w-full" />
      <Skeleton class="h-24 w-full" />
    </div>
    <div v-else class="grid gap-4 md:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Site</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-2">
            <Label for="site-title" class="text-sm text-muted-foreground">Site title</Label>
            <Input
              id="site-title"
              :model-value="siteTitle"
              @update:model-value="siteTitle = String($event)"
              :disabled="isSaving"
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Authentication</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <div class="font-medium">Allow signups</div>
              <div class="text-sm text-muted-foreground">Permit new user self-registration</div>
            </div>
            <Checkbox
              :checked="allowSignups"
              @update:checked="allowSignups = $event"
              :disabled="isSaving"
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Requests</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-2">
            <Label for="max-per-user" class="text-sm text-muted-foreground">Max per user</Label>
            <Input
              id="max-per-user"
              type="number"
              :model-value="maxPerUser"
              @update:model-value="maxPerUser = Number($event)"
              :min="0"
              :step="1"
              :disabled="isSaving"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
