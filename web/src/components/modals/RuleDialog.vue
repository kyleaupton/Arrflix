<script setup lang="ts">
import { ref, inject, watch, computed } from 'vue'
import { useMutation, useQueryClient, useQuery } from '@tanstack/vue-query'
import {
  getV1PoliciesByIdRuleOptions,
  postV1PoliciesByIdRuleMutation,
  putV1PoliciesByIdRuleMutation,
  deleteV1PoliciesByIdRuleMutation,
} from '@/client/@tanstack/vue-query.gen'
import { type DbgenPolicy, type DbgenRule } from '@/client/types.gen'
import { usePolicyFields } from '@/composables/usePolicyFields'
import { getV1IndexersConfiguredOptions } from '@/client/@tanstack/vue-query.gen'
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
import { useModal } from '@/composables/useModal'

interface Props {
  policy: DbgenPolicy
}

const props = defineProps<Props>()

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }
const queryClient = useQueryClient()
const modal = useModal()

const { fields: fieldDefinitions, getFieldByPath, getValidOperators } = usePolicyFields()
const { data: indexers } = useQuery(getV1IndexersConfiguredOptions())

const createRuleMutation = useMutation(postV1PoliciesByIdRuleMutation())
const updateRuleMutation = useMutation(putV1PoliciesByIdRuleMutation())
const deleteRuleMutation = useMutation(deleteV1PoliciesByIdRuleMutation())

const editingRule = ref<DbgenRule | null>(null)
const ruleForm = ref({
  left_operand: '',
  operator: '',
  right_operand: '',
})

// Field options for left operand dropdown
const fieldOptions = computed(() => {
  return (
    fieldDefinitions.value?.map((f) => ({
      label: f.label,
      value: f.path,
    })) || []
  )
})

// Selected field definition
const selectedField = computed(() => {
  if (!ruleForm.value.left_operand) return undefined
  return getFieldByPath(ruleForm.value.left_operand)
})

// Valid operators for selected field
const validOperators = computed(() => {
  const ops = getValidOperators(selectedField.value)
  return policyOptions.operators.filter((op) => ops.includes(op.value))
})

// Right operand options (for enum/dynamic fields)
const rightOperandOptions = computed(() => {
  const field = selectedField.value
  if (!field) return []

  if (field.type === 'enum' && field.enumValues) {
    return field.enumValues.map((ev) => ({
      label: ev.label,
      value: ev.value,
    }))
  }

  if (field.type === 'dynamic' && field.dynamicSource === '/api/v1/indexers/configured') {
    return (
      indexers.value?.map((idx) => ({
        label: idx.name || 'Unknown',
        value: idx.name || '',
      })) || []
    )
  }

  if (field.type === 'boolean') {
    return [
      { label: 'True', value: 'true' },
      { label: 'False', value: 'false' },
    ]
  }

  return []
})

// Load rule when component mounts
watch(
  () => props.policy,
  async (policy) => {
    if (!policy?.id) return
    try {
      const rule = await queryClient.fetchQuery(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        getV1PoliciesByIdRuleOptions({ path: { id: String(policy.id) } } as any),
      )
      if (rule) {
        editingRule.value = rule as DbgenRule
        ruleForm.value = {
          left_operand: rule.left_operand ?? '',
          operator: rule.operator ?? '',
          right_operand: rule.right_operand ?? '',
        }
      } else {
        editingRule.value = null
        ruleForm.value = { left_operand: '', operator: '', right_operand: '' }
      }
    } catch {
      // No rule exists yet
      editingRule.value = null
      ruleForm.value = { left_operand: '', operator: '', right_operand: '' }
    }
  },
  { immediate: true },
)

const handleSave = async () => {
  if (!props.policy?.id) return
  try {
    if (editingRule.value?.id) {
      await updateRuleMutation.mutateAsync({
        path: { id: props.policy.id },
        body: ruleForm.value,
      })
    } else {
      await createRuleMutation.mutateAsync({
        path: { id: props.policy.id },
        body: ruleForm.value,
      })
    }
    dialogRef.value.close({ saved: true })
  } catch (error) {
    console.error('Failed to save rule:', error)
  }
}

const handleDelete = async () => {
  if (!props.policy?.id || !editingRule.value?.id) return
  const confirmed = await modal.confirm({
    title: 'Delete Rule',
    message: 'Are you sure you want to delete this rule?',
    severity: 'danger',
  })
  if (!confirmed) return
  try {
    await deleteRuleMutation.mutateAsync({ path: { id: String(props.policy.id) } })
    dialogRef.value.close({ deleted: true })
  } catch (err) {
    console.error('Failed to delete rule:', err)
  }
}

const handleCancel = () => {
  dialogRef.value.close()
}

const handleLeftOperandChange = (value: string) => {
  ruleForm.value.left_operand = value
  ruleForm.value.operator = ''
  ruleForm.value.right_operand = ''
}

const isLoading = computed(
  () => createRuleMutation.isPending.value || updateRuleMutation.isPending.value,
)
</script>

<template>
  <BaseDialog title="Edit Rule">
    <div class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="rule-left-operand">Left Operand</Label>
        <Select v-model="ruleForm.left_operand" @update:model-value="handleLeftOperandChange">
          <SelectTrigger id="rule-left-operand" class="w-full">
            <SelectValue placeholder="Select field" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="option in fieldOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="flex flex-col gap-2">
        <Label for="rule-operator">Operator</Label>
        <Select v-model="ruleForm.operator" :disabled="!selectedField">
          <SelectTrigger id="rule-operator" class="w-full">
            <SelectValue placeholder="Select operator" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="op in validOperators" :key="op.value" :value="op.value">
              {{ op.label }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="flex flex-col gap-2">
        <Label for="rule-right-operand">Right Operand</Label>
        <!-- Enum or Dynamic dropdown -->
        <Select
          v-if="
            selectedField &&
            (selectedField.type === 'enum' ||
              selectedField.type === 'dynamic' ||
              selectedField.type === 'boolean')
          "
          v-model="ruleForm.right_operand"
        >
          <SelectTrigger id="rule-right-operand" class="w-full">
            <SelectValue placeholder="Select value" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="option in rightOperandOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </SelectItem>
          </SelectContent>
        </Select>
        <!-- Number input -->
        <Input
          v-else-if="selectedField && selectedField.type === 'number'"
          id="rule-right-operand"
          type="number"
          :model-value="Number(ruleForm.right_operand) || 0"
          @update:model-value="(val) => (ruleForm.right_operand = String(val ?? 0))"
          placeholder="Enter number"
        />
        <!-- Text input -->
        <Input
          v-else
          id="rule-right-operand"
          v-model="ruleForm.right_operand"
          placeholder="Enter value"
          :disabled="!selectedField"
        />
      </div>
    </div>
    <template #footer>
      <Button variant="outline" @click="handleCancel">Cancel</Button>
      <Button v-if="editingRule" variant="destructive" @click="handleDelete">Delete</Button>
      <Button :disabled="isLoading" @click="handleSave">Save</Button>
    </template>
  </BaseDialog>
</template>

