<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import ToggleSwitch from 'primevue/toggleswitch'
import Message from 'primevue/message'
import Skeleton from 'primevue/skeleton'
import { getV1Settings, patchV1Settings } from '@/client/sdk.gen'

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
    <Message v-if="error" severity="error">{{ error }}</Message>
    <div v-if="isLoading" class="space-y-3">
      <Skeleton height="3rem" />
      <Skeleton height="3rem" />
      <Skeleton height="3rem" />
    </div>
    <div v-else class="grid gap-4 md:grid-cols-2">
      <Card>
        <template #title>Site</template>
        <template #content>
          <div class="flex flex-col gap-2">
            <label class="text-sm text-muted-color">Site title</label>
            <InputText :value="siteTitle" @update:value="siteTitle = $event" :disabled="isSaving" />
          </div>
        </template>
      </Card>

      <Card>
        <template #title>Authentication</template>
        <template #content>
          <div class="flex items-center justify-between">
            <div>
              <div class="font-medium">Allow signups</div>
              <div class="text-sm text-muted-color">Permit new user self-registration</div>
            </div>
            <ToggleSwitch
              :model-value="allowSignups"
              @update:model-value="allowSignups = $event"
              :disabled="isSaving"
            />
          </div>
        </template>
      </Card>

      <Card>
        <template #title>Requests</template>
        <template #content>
          <div class="flex flex-col gap-2">
            <label class="text-sm text-muted-color">Max per user</label>
            <InputNumber
              :value="maxPerUser"
              @update:value="maxPerUser = Number($event)"
              :min="0"
              :step="1"
              :disabled="isSaving"
            />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>
