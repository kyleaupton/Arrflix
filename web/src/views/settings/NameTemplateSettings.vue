<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Button from 'primevue/button'
import Message from 'primevue/message'
import Skeleton from 'primevue/skeleton'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import {
  getV1NameTemplates,
  postV1NameTemplates,
  putV1NameTemplatesById,
  deleteV1NameTemplatesById,
} from '@/client/sdk.gen'

type NameTemplate = {
  id?: string
  name?: string
  type?: 'movie' | 'series' | string
  template?: string
  default?: boolean
}

const templates = ref<NameTemplate[]>([])
const templatesLoading = ref(false)
const templatesError = ref<string | null>(null)

const movieVariables = ['{Title}', '{Year}', '{Quality}', '{Resolution}', '{Extension}']
const seriesVariables = [
  '{Title}',
  '{Season}',
  '{Episode}',
  '{EpisodeTitle}',
  '{Quality}',
  '{Resolution}',
  '{Extension}',
]

const availableVariables = computed(() => {
  return newTemplate.value.type === 'series' ? seriesVariables : movieVariables
})

async function loadTemplates() {
  templatesLoading.value = true
  templatesError.value = null
  try {
    const res = await getV1NameTemplates<true>({ throwOnError: true })
    templates.value = (res.data as NameTemplate[]) ?? []
  } catch {
    templatesError.value = 'Failed to load name templates'
  } finally {
    templatesLoading.value = false
  }
}

onMounted(loadTemplates)

// Create form
const newTemplate = ref<Required<Omit<NameTemplate, 'id'>>>({
  name: '',
  type: 'movie',
  template: '',
  default: false,
})
const isCreating = ref(false)
async function createTemplate() {
  if (!newTemplate.value.name || !newTemplate.value.template) return
  isCreating.value = true
  try {
    const res = await postV1NameTemplates<true>({
      throwOnError: true,
      body: {
        name: newTemplate.value.name,
        type: newTemplate.value.type,
        template: newTemplate.value.template,
        default: newTemplate.value.default,
      },
    })
    templates.value = [...templates.value, res.data as NameTemplate]
    newTemplate.value = { name: '', type: 'movie', template: '', default: false }
    await loadTemplates() // Reload to refresh default flags
  } catch (err: any) {
    templatesError.value = err.message || 'Failed to create template'
  } finally {
    isCreating.value = false
  }
}

// Edit helpers
const editingId = ref<string | null>(null)
const editBuf = ref<Required<Omit<NameTemplate, 'id'>>>({
  name: '',
  type: 'movie',
  template: '',
  default: false,
})

const editAvailableVariables = computed(() => {
  return editBuf.value.type === 'series' ? seriesVariables : movieVariables
})

function startEdit(template: NameTemplate) {
  editingId.value = template.id ?? null
  editBuf.value = {
    name: template.name ?? '',
    type: (template.type as 'movie' | 'series') ?? 'movie',
    template: template.template ?? '',
    default: template.default ?? false,
  }
}

async function saveEdit(id: string) {
  try {
    await putV1NameTemplatesById<true>({
      throwOnError: true,
      path: { id },
      body: {
        name: editBuf.value.name,
        type: editBuf.value.type,
        template: editBuf.value.template,
        default: editBuf.value.default,
      },
    })
    templates.value = templates.value.map((t) => (t.id === id ? { ...t, ...editBuf.value } : t))
    editingId.value = null
    await loadTemplates() // Reload to refresh default flags
  } catch (err: any) {
    templatesError.value = err.message || 'Failed to update template'
  }
}

async function removeTemplate(id: string) {
  try {
    await deleteV1NameTemplatesById<true>({ throwOnError: true, path: { id } })
    templates.value = templates.value.filter((t) => t.id !== id)
  } catch (err: any) {
    templatesError.value = err.message || 'Failed to delete template'
  }
}
</script>

<template>
  <div>
    <Message v-if="templatesError" severity="error" @close="templatesError = null">{{
      templatesError
    }}</Message>

    <div class="grid gap-4 md:grid-cols-2">
      <Card>
        <template #title>Add Name Template</template>
        <template #content>
          <div class="flex flex-col gap-3">
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Name</label>
              <InputText v-model="newTemplate.name" placeholder="My Movie Template" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Type</label>
              <Select
                v-model="newTemplate.type"
                :options="[
                  { label: 'Movies', value: 'movie' },
                  { label: 'Series', value: 'series' },
                ]"
                optionLabel="label"
                optionValue="value"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">Template</label>
              <Textarea
                v-model="newTemplate.template"
                :placeholder="
                  newTemplate.type === 'series'
                    ? '{Title} S{Season}E{Episode} - {EpisodeTitle}'
                    : '{Title} ({Year})'
                "
                rows="3"
              />
              <div class="text-xs text-muted-color mt-1">
                Available variables:
                <span class="font-mono">{{ availableVariables.join(', ') }}</span>
              </div>
            </div>
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-color">Default</div>
              <ToggleSwitch v-model="newTemplate.default" />
            </div>
            <div class="flex justify-end">
              <Button label="Create" :loading="isCreating" @click="createTemplate" />
            </div>
          </div>
        </template>
      </Card>

      <Card>
        <template #title>Name Templates</template>
        <template #content>
          <div v-if="templatesLoading" class="space-y-2">
            <Skeleton height="2.5rem" />
            <Skeleton height="2.5rem" />
          </div>
          <div v-else class="space-y-3">
            <div
              v-for="template in templates"
              :key="template.id"
              class="border rounded p-3 flex flex-col gap-2"
            >
              <div class="flex items-center justify-between">
                <div class="font-medium">
                  {{ template.name }}
                  <span class="text-sm text-muted-color">({{ template.type }})</span>
                  <span
                    v-if="template.default"
                    class="ml-2 text-xs px-2 py-0.5 bg-green-500 text-white rounded"
                  >
                    Default
                  </span>
                </div>
                <div class="flex items-center gap-2">
                  <Button size="small" label="Edit" @click="startEdit(template)" />
                  <Button
                    size="small"
                    label="Delete"
                    severity="danger"
                    @click="template.id && removeTemplate(template.id)"
                  />
                </div>
              </div>
              <div class="text-sm font-mono text-muted-color bg-gray-100 dark:bg-gray-800 p-2 rounded">
                {{ template.template }}
              </div>
              <div class="flex gap-4 text-xs">
                <div>
                  Default:
                  <span :class="template.default ? 'text-green-500' : 'text-muted-color'">{{
                    template.default ? 'Yes' : 'No'
                  }}</span>
                </div>
              </div>

              <div v-if="editingId === template.id" class="mt-2 border-t pt-3 space-y-3">
                <div class="grid gap-3 md:grid-cols-2">
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Name</label>
                    <InputText v-model="editBuf.name" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Type</label>
                    <Select
                      v-model="editBuf.type"
                      :options="[
                        { label: 'Movies', value: 'movie' },
                        { label: 'Series', value: 'series' },
                      ]"
                      optionLabel="label"
                      optionValue="value"
                    />
                  </div>
                  <div class="md:col-span-2 flex flex-col gap-1">
                    <label class="text-sm text-muted-color">Template</label>
                    <Textarea v-model="editBuf.template" rows="3" />
                    <div class="text-xs text-muted-color mt-1">
                      Available variables:
                      <span class="font-mono">{{ editAvailableVariables.join(', ') }}</span>
                    </div>
                  </div>
                  <div class="md:col-span-2 flex items-center justify-between">
                    <div class="text-sm text-muted-color">Default</div>
                    <ToggleSwitch v-model="editBuf.default" />
                  </div>
                </div>
                <div class="flex justify-end gap-2">
                  <Button
                    size="small"
                    label="Cancel"
                    severity="secondary"
                    @click="editingId = null"
                  />
                  <Button size="small" label="Save" @click="template.id && saveEdit(template.id)" />
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

