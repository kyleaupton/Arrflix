<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import { PrimeIcons } from '@/icons'
import {
  getV1NameTemplatesOptions,
  postV1NameTemplatesMutation,
  putV1NameTemplatesByIdMutation,
  deleteV1NameTemplatesByIdMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type HandlersNameTemplateSwagger } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import {
  nameTemplateColumns,
  createNameTemplateActions,
} from '@/components/tables/configs/nameTemplateTableConfig'
import { useModal } from '@/composables/useModal'

const queryClient = useQueryClient()
const modal = useModal()

// Data queries
const { data: templates, isLoading, refetch } = useQuery(getV1NameTemplatesOptions())

// Mutations
const createTemplateMutation = useMutation(postV1NameTemplatesMutation())
const updateTemplateMutation = useMutation(putV1NameTemplatesByIdMutation())
const deleteTemplateMutation = useMutation(deleteV1NameTemplatesByIdMutation())

// Modal state
const showTemplateModal = ref(false)
const editingTemplate = ref<HandlersNameTemplateSwagger | null>(null)
const templateError = ref<string | null>(null)

// Template form
const templateForm = ref({
  name: '',
  type: 'movie' as 'movie' | 'series',
  template: '',
  default: false,
})

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
  return templateForm.value.type === 'series' ? seriesVariables : movieVariables
})

// Handlers
const handleAddTemplate = () => {
  editingTemplate.value = null
  templateForm.value = { name: '', type: 'movie', template: '', default: false }
  templateError.value = null
  showTemplateModal.value = true
}

const handleEditTemplate = (template: HandlersNameTemplateSwagger) => {
  editingTemplate.value = template
  templateForm.value = {
    name: template.name || '',
    type: (template.type as 'movie' | 'series') || 'movie',
    template: template.template || '',
    default: template.default || false,
  }
  templateError.value = null
  showTemplateModal.value = true
}

const handleSaveTemplate = async () => {
  if (!templateForm.value.name || !templateForm.value.template) {
    templateError.value = 'Name and template are required'
    return
  }

  try {
    if (editingTemplate.value?.id) {
      await updateTemplateMutation.mutateAsync({
        path: { id: editingTemplate.value.id },
        body: {
          name: templateForm.value.name,
          type: templateForm.value.type,
          template: templateForm.value.template,
          default: templateForm.value.default,
        },
      })
    } else {
      await createTemplateMutation.mutateAsync({
        body: {
          name: templateForm.value.name,
          type: templateForm.value.type,
          template: templateForm.value.template,
          default: templateForm.value.default,
        },
      })
    }
    showTemplateModal.value = false
    refetch()
  } catch (err: any) {
    templateError.value = err.message || 'Failed to save template'
  }
}

const handleDeleteTemplate = async (template: HandlersNameTemplateSwagger) => {
  if (!template.id) return
  const confirmed = await modal.confirm({
    title: 'Delete Template',
    message: `Are you sure you want to delete "${template.name}"?`,
    severity: 'danger',
  })
  if (!confirmed) return
  try {
    await deleteTemplateMutation.mutateAsync({ path: { id: template.id } })
    refetch()
  } catch (err: any) {
    templateError.value = err.message || 'Failed to delete template'
  }
}

const templateActions = createNameTemplateActions(handleEditTemplate, handleDeleteTemplate)
</script>

<template>
  <div class="name-templates-settings">
    <div class="card">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold mb-2">Name Templates</h3>
            <p class="text-muted-color">Configure templates for naming downloaded media files.</p>
          </div>
          <Button
            label="Add Template"
            :icon="PrimeIcons.PLUS"
            severity="primary"
            @click="handleAddTemplate"
          />
        </div>

        <Message v-if="templateError" severity="error" @close="templateError = null">{{
          templateError
        }}</Message>

        <DataTable
          :data="templates || []"
          :columns="nameTemplateColumns"
          :actions="templateActions"
          :loading="isLoading"
          empty-message="No name templates configured"
          searchable
          search-placeholder="Search templates..."
          paginator
          :rows="10"
        />
      </div>
    </div>

    <!-- Template Modal -->
    <Dialog
      v-model:visible="showTemplateModal"
      :header="editingTemplate ? 'Edit Name Template' : 'Add Name Template'"
      :modal="true"
      :style="{ width: '600px' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Name</label>
          <InputText v-model="templateForm.name" placeholder="My Movie Template" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Type</label>
          <Select
            v-model="templateForm.type"
            :options="[
              { label: 'Movies', value: 'movie' },
              { label: 'Series', value: 'series' },
            ]"
            optionLabel="label"
            optionValue="value"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Template</label>
          <Textarea
            v-model="templateForm.template"
            :placeholder="
              templateForm.type === 'series'
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
          <label class="text-sm font-medium">Default</label>
          <ToggleSwitch v-model="templateForm.default" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showTemplateModal = false" />
        <Button
          label="Save"
          :loading="
            createTemplateMutation.isPending.value || updateTemplateMutation.isPending.value
          "
          @click="handleSaveTemplate"
        />
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.name-templates-settings {
  max-width: 100%;
}
</style>
