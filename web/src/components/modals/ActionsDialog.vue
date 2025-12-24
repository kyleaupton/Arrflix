<script setup lang="ts">
import { ref, inject, watch, computed } from 'vue'
import { useMutation, useQueryClient, useQuery } from '@tanstack/vue-query'
import {
  getV1PoliciesByIdActionsOptions,
  postV1PoliciesByIdActionsMutation,
  deleteV1PoliciesByIdActionsByActionIdMutation,
  getV1LibrariesOptions,
  getV1NameTemplatesOptions,
  getV1DownloadersOptions,
} from '@/client/@tanstack/vue-query.gen'
import { type DbgenPolicy, type DbgenAction } from '@/client/types.gen'
import policyOptions from '@/config/policyOptions.json'
import BaseDialog from './BaseDialog.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Card, CardContent } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Trash2 } from 'lucide-vue-next'

interface Props {
  policy: DbgenPolicy
}

const props = defineProps<Props>()

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }
const queryClient = useQueryClient()

const createActionMutation = useMutation(postV1PoliciesByIdActionsMutation())
const deleteActionMutation = useMutation(deleteV1PoliciesByIdActionsByActionIdMutation())

const { data: libraries } = useQuery(getV1LibrariesOptions())
const { data: nameTemplates } = useQuery(getV1NameTemplatesOptions())
const { data: downloaders } = useQuery(getV1DownloadersOptions())

const editingActions = ref<DbgenAction[]>([])

const actionTypeOptions = computed(() => policyOptions.actionTypes)

const libraryOptions = computed(() => {
  return (
    libraries.value?.map((lib) => ({
      label: `${lib.name} (${lib.type})`,
      value: lib.id,
    })) || []
  )
})

const nameTemplateOptions = computed(() => {
  return (
    nameTemplates.value?.map((nt) => ({
      label: `${nt.name} (${nt.type})`,
      value: nt.id,
    })) || []
  )
})

const downloaderOptions = computed(() => {
  return (
    downloaders.value?.map((d) => ({
      label: `${d.name} (${d.type}/${d.protocol})`,
      value: d.id,
    })) || []
  )
})

const getActionValueOptions = (actionType: string) => {
  switch (actionType) {
    case 'set_downloader':
      return downloaderOptions.value
    case 'set_library':
      return libraryOptions.value
    case 'set_name_template':
      return nameTemplateOptions.value
    default:
      return []
  }
}

// Load actions when component mounts
watch(
  () => props.policy,
  async (policy) => {
    if (!policy?.id) return
    try {
      const result = await queryClient.fetchQuery(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        getV1PoliciesByIdActionsOptions({ path: { id: String(policy.id) } } as any),
      )
      editingActions.value = (result as DbgenAction[]) || []
    } catch {
      editingActions.value = []
    }
  },
  { immediate: true },
)

const handleAddAction = () => {
  editingActions.value.push({
    id: `new-${Date.now()}`,
    type: '',
    value: '',
    order: editingActions.value.length,
  } as unknown as DbgenAction)
}

const handleRemoveAction = (index: number) => {
  editingActions.value.splice(index, 1)
  // Reorder
  editingActions.value.forEach((action, idx) => {
    action.order = idx
  })
}

const handleActionTypeChange = (index: number, value: string) => {
  if (editingActions.value[index]) {
    editingActions.value[index].type = value
    editingActions.value[index].value = ''
  }
}

const handleSave = async () => {
  if (!props.policy?.id) return
  try {
    // Delete all existing actions first
    const existingActions = editingActions.value.filter(
      (a) => a.id && !String(a.id).startsWith('new-'),
    )
    for (const action of existingActions) {
      if (action.id && props.policy?.id) {
        try {
          await deleteActionMutation.mutateAsync({
            path: { id: String(props.policy.id), actionId: String(action.id) },
          })
        } catch {
          // Ignore errors
        }
      }
    }

    // Create new actions
    for (const action of editingActions.value) {
      await createActionMutation.mutateAsync({
        path: { id: props.policy.id },
        body: {
          type: action.type,
          value: action.value,
          order: action.order,
        },
      })
    }
    dialogRef.value.close({ saved: true })
  } catch (err) {
    console.error('Failed to save actions:', err)
  }
}

const handleCancel = () => {
  dialogRef.value.close()
}

const isLoading = computed(() => createActionMutation.isPending.value)
</script>

<template>
  <BaseDialog title="Edit Actions">
    <div class="flex flex-col gap-4">
      <div class="flex justify-between items-center">
        <p class="text-sm text-muted-foreground">
          Actions are executed in order when the policy matches.
        </p>
        <Button size="sm" @click="handleAddAction">Add Action</Button>
      </div>
      <div v-if="editingActions.length === 0" class="text-center py-8 text-muted-foreground">
        No actions configured. Click "Add Action" to add one.
      </div>
      <div v-else class="space-y-3">
        <Card
          v-for="(action, index) in editingActions"
          :key="index"
          class="border rounded-lg"
        >
          <CardContent class="p-4">
            <div class="flex flex-col gap-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">Action {{ index + 1 }}</span>
                <Button
                  size="sm"
                  variant="ghost"
                  @click="handleRemoveAction(index)"
                >
                  <Trash2 class="size-4" />
                </Button>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div class="flex flex-col gap-2">
                  <Label :for="`action-type-${index}`">Type</Label>
                  <Select
                    :model-value="action.type"
                    @update:model-value="(val) => handleActionTypeChange(index, val as string)"
                  >
                    <SelectTrigger :id="`action-type-${index}`" class="w-full">
                      <SelectValue placeholder="Select action type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="option in actionTypeOptions"
                        :key="option.value"
                        :value="option.value"
                      >
                        {{ option.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div class="flex flex-col gap-2">
                  <Label :for="`action-value-${index}`">Value</Label>
                  <Select
                    v-if="action.type && getActionValueOptions(action.type).length > 0"
                    v-model="action.value"
                  >
                    <SelectTrigger :id="`action-value-${index}`" class="w-full">
                      <SelectValue placeholder="Select value" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="option in getActionValueOptions(action.type)"
                        :key="option.value"
                        :value="String(option.value)"
                      >
                        {{ option.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <Input
                    v-else-if="action.type === 'stop_processing'"
                    :id="`action-value-${index}`"
                    model-value="N/A"
                    disabled
                  />
                  <Input
                    v-else
                    :id="`action-value-${index}`"
                    v-model="action.value"
                    placeholder="Enter value"
                  />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
    <template #footer>
      <Button variant="outline" @click="handleCancel">Cancel</Button>
      <Button :disabled="isLoading" @click="handleSave">Save</Button>
    </template>
  </BaseDialog>
</template>

