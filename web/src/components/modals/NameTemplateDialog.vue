<script setup lang="ts">
import { ref, inject, watch, computed } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import {
  postV1NameTemplatesMutation,
  putV1NameTemplatesByIdMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type HandlersNameTemplateSwagger } from '@/client/types.gen'
import BaseDialog from './BaseDialog.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface Props {
  template?: HandlersNameTemplateSwagger | null
}

const props = defineProps<Props>()

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }

const createTemplateMutation = useMutation(postV1NameTemplatesMutation())
const updateTemplateMutation = useMutation(putV1NameTemplatesByIdMutation())

const templateForm = ref({
  name: '',
  type: 'movie' as 'movie' | 'series',
  template: '',
  default: false,
})

const templateError = ref<string | null>(null)

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

const typeOptions = [
  { label: 'Movies', value: 'movie' },
  { label: 'Series', value: 'series' },
]

// Initialize form when template changes
watch(
  () => props.template,
  (template) => {
    if (template) {
      templateForm.value = {
        name: template.name || '',
        type: (template.type as 'movie' | 'series') || 'movie',
        template: template.template || '',
        default: template.default || false,
      }
    } else {
      templateForm.value = { name: '', type: 'movie', template: '', default: false }
    }
    templateError.value = null
  },
  { immediate: true },
)

const handleSave = async () => {
  if (!templateForm.value.name || !templateForm.value.template) {
    templateError.value = 'Name and template are required'
    return
  }

  try {
    if (props.template?.id) {
      await updateTemplateMutation.mutateAsync({
        path: { id: props.template.id },
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
    templateError.value = null
    dialogRef.value.close({ saved: true })
  } catch (err) {
    const error = err as { message?: string }
    templateError.value = error.message || 'Failed to save template'
  }
}

const handleCancel = () => {
  dialogRef.value.close()
}

const isLoading = computed(
  () => createTemplateMutation.isPending.value || updateTemplateMutation.isPending.value,
)
</script>

<template>
  <BaseDialog :title="template ? 'Edit Name Template' : 'Add Name Template'">
    <div class="flex flex-col gap-4">
      <div
        v-if="templateError"
        class="p-4 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive text-sm"
      >
        {{ templateError }}
      </div>

      <div class="flex flex-col gap-2">
        <Label for="template-name">Name</Label>
        <Input id="template-name" v-model="templateForm.name" placeholder="My Movie Template" />
      </div>

      <div class="flex flex-col gap-2">
        <Label for="template-type">Type</Label>
        <Select v-model="templateForm.type">
          <SelectTrigger id="template-type" class="w-full">
            <SelectValue placeholder="Select type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="option in typeOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex flex-col gap-2">
        <Label for="template-template">Template</Label>
        <Textarea
          id="template-template"
          v-model="templateForm.template"
          :placeholder="
            templateForm.type === 'series'
              ? '{Title} S{Season}E{Episode} - {EpisodeTitle}'
              : '{Title} ({Year})'
          "
          rows="3"
        />
        <div class="text-xs text-muted-foreground mt-1">
          Available variables:
          <span class="font-mono">{{ availableVariables.join(', ') }}</span>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <Label for="template-default">Default</Label>
        <Switch id="template-default" v-model:checked="templateForm.default" />
      </div>
    </div>

    <template #footer>
      <Button variant="outline" @click="handleCancel">Cancel</Button>
      <Button :disabled="isLoading" @click="handleSave">Save</Button>
    </template>
  </BaseDialog>
</template>

