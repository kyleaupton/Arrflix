<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import ToggleSwitch from 'primevue/toggleswitch'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Message from 'primevue/message'
import { PrimeIcons } from '@/icons'
import {
  getV1PoliciesOptions,
  postV1PoliciesMutation,
  putV1PoliciesByIdMutation,
  deleteV1PoliciesByIdMutation,
  getV1PoliciesByIdRuleOptions,
  postV1PoliciesByIdRuleMutation,
  putV1PoliciesByIdRuleMutation,
  deleteV1PoliciesByIdRuleMutation,
  getV1PoliciesByIdActionsOptions,
  postV1PoliciesByIdActionsMutation,
  putV1PoliciesByIdActionsByActionIdMutation,
  deleteV1PoliciesByIdActionsByActionIdMutation,
  getV1LibrariesOptions,
  getV1NameTemplatesOptions,
  getV1DownloadersOptions,
} from '@/client/@tanstack/vue-query.gen'
import { type DbgenPolicy, type DbgenRule, type DbgenAction } from '@/client/types.gen'
import DataTable from '@/components/tables/DataTable.vue'
import { policyColumns, createPolicyActions } from '@/components/tables/configs/policyTableConfig'
import policyOptions from '@/config/policyOptions.json'

const queryClient = useQueryClient()

// Data queries
const { data: policies, isLoading, refetch } = useQuery(getV1PoliciesOptions())
const { data: libraries } = useQuery(getV1LibrariesOptions())
const { data: nameTemplates } = useQuery(getV1NameTemplatesOptions())
const { data: downloaders } = useQuery(getV1DownloadersOptions())

// Mutations
const createPolicyMutation = useMutation(postV1PoliciesMutation())
const updatePolicyMutation = useMutation(putV1PoliciesByIdMutation())
const deletePolicyMutation = useMutation(deleteV1PoliciesByIdMutation())
const createRuleMutation = useMutation(postV1PoliciesByIdRuleMutation())
const updateRuleMutation = useMutation(putV1PoliciesByIdRuleMutation())
const deleteRuleMutation = useMutation(deleteV1PoliciesByIdRuleMutation())
const createActionMutation = useMutation(postV1PoliciesByIdActionsMutation())
const updateActionMutation = useMutation(putV1PoliciesByIdActionsByActionIdMutation())
const deleteActionMutation = useMutation(deleteV1PoliciesByIdActionsByActionIdMutation())

// Modal states
const showPolicyModal = ref(false)
const showRuleModal = ref(false)
const showActionsModal = ref(false)
const editingPolicy = ref<DbgenPolicy | null>(null)
const editingRule = ref<DbgenRule | null>(null)
const editingActions = ref<DbgenAction[]>([])

// Policy form
const policyForm = ref({
  name: '',
  description: '',
  enabled: true,
  priority: 0,
})

// Rule form
const ruleForm = ref({
  left_operand: '',
  operator: '',
  right_operand: '',
})

// Action form
const actionForm = ref({
  type: '',
  value: '',
  order: 0,
})

// Computed values
const actionTypeOptions = computed(() => policyOptions.actionTypes)
const operatorOptions = computed(() => policyOptions.operators)
const torrentFieldOptions = computed(() => policyOptions.torrentFields)

const libraryOptions = computed(() => {
  return libraries.value?.map((lib) => ({
    label: `${lib.name} (${lib.type})`,
    value: lib.id,
  })) || []
})

const nameTemplateOptions = computed(() => {
  return nameTemplates.value?.map((nt) => ({
    label: `${nt.name} (${nt.type})`,
    value: nt.id,
  })) || []
})

const downloaderOptions = computed(() => {
  return downloaders.value?.map((d) => ({
    label: `${d.name} (${d.type}/${d.protocol})`,
    value: d.id,
  })) || []
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

// Policy handlers
const handleAddPolicy = () => {
  editingPolicy.value = null
  policyForm.value = { name: '', description: '', enabled: true, priority: 0 }
  showPolicyModal.value = true
}

const handleEditPolicy = async (policy: DbgenPolicy) => {
  editingPolicy.value = policy
  policyForm.value = {
    name: policy.name || '',
    description: policy.description || '',
    enabled: policy.enabled ?? true,
    priority: policy.priority || 0,
  }
  showPolicyModal.value = true
}

const handleSavePolicy = async () => {
  try {
    if (editingPolicy.value?.id) {
      await updatePolicyMutation.mutateAsync({
        path: { id: editingPolicy.value.id },
        body: {
          name: policyForm.value.name,
          description: policyForm.value.description || undefined,
          enabled: policyForm.value.enabled,
          priority: policyForm.value.priority,
        },
      })
    } else {
      await createPolicyMutation.mutateAsync({
        body: {
          name: policyForm.value.name,
          description: policyForm.value.description || undefined,
          enabled: policyForm.value.enabled,
          priority: policyForm.value.priority,
        },
      })
    }
    showPolicyModal.value = false
    refetch()
  } catch (error) {
    console.error('Failed to save policy:', error)
  }
}

const handleDeletePolicy = async (policy: DbgenPolicy) => {
  if (!policy.id || !confirm(`Are you sure you want to delete "${policy.name}"?`)) return
  try {
    await deletePolicyMutation.mutateAsync({ path: { id: policy.id } })
    refetch()
  } catch (error) {
    console.error('Failed to delete policy:', error)
  }
}

// Rule handlers
const handleEditRule = async (policy: DbgenPolicy) => {
  if (!policy.id) return
  editingPolicy.value = policy
  try {
    const { data: rule } = await queryClient.fetchQuery(
      getV1PoliciesByIdRuleOptions({ path: { id: policy.id } })
    )
    if (rule) {
      editingRule.value = rule
      ruleForm.value = {
        left_operand: rule.left_operand || '',
        operator: rule.operator || '',
        right_operand: rule.right_operand || '',
      }
    } else {
      editingRule.value = null
      ruleForm.value = { left_operand: '', operator: '', right_operand: '' }
    }
    showRuleModal.value = true
  } catch (error) {
    // No rule exists yet
    editingRule.value = null
    ruleForm.value = { left_operand: '', operator: '', right_operand: '' }
    showRuleModal.value = true
  }
}

const handleSaveRule = async () => {
  if (!editingPolicy.value?.id) return
  try {
    if (editingRule.value?.id) {
      await updateRuleMutation.mutateAsync({
        path: { id: editingPolicy.value.id },
        body: ruleForm.value,
      })
    } else {
      await createRuleMutation.mutateAsync({
        path: { id: editingPolicy.value.id },
        body: ruleForm.value,
      })
    }
    showRuleModal.value = false
    refetch()
  } catch (error) {
    console.error('Failed to save rule:', error)
  }
}

const handleDeleteRule = async (policy: DbgenPolicy) => {
  if (!policy.id || !confirm('Are you sure you want to delete this rule?')) return
  try {
    await deleteRuleMutation.mutateAsync({ path: { id: policy.id } })
    refetch()
  } catch (error) {
    console.error('Failed to delete rule:', error)
  }
}

// Actions handlers
const handleEditActions = async (policy: DbgenPolicy) => {
  if (!policy.id) return
  editingPolicy.value = policy
  try {
    const result = await queryClient.fetchQuery(
      getV1PoliciesByIdActionsOptions({ path: { id: policy.id } })
    )
    editingActions.value = (result as DbgenAction[]) || []
    showActionsModal.value = true
  } catch (error) {
    editingActions.value = []
    showActionsModal.value = true
  }
}

const handleAddAction = () => {
  actionForm.value = { type: '', value: '', order: editingActions.value.length }
  editingActions.value.push({
    id: `new-${Date.now()}`,
    type: '',
    value: '',
    order: editingActions.value.length,
  } as any)
}

const handleRemoveAction = (index: number) => {
  editingActions.value.splice(index, 1)
  // Reorder
  editingActions.value.forEach((action, idx) => {
    action.order = idx
  })
}

const handleSaveActions = async () => {
  if (!editingPolicy.value?.id) return
  try {
    // Delete all existing actions first
    const existingActions = editingActions.value.filter((a) => a.id && !a.id.toString().startsWith('new-'))
    for (const action of existingActions) {
      if (action.id) {
        try {
          await deleteActionMutation.mutateAsync({
            path: { id: editingPolicy.value.id, actionId: action.id },
          })
        } catch (error) {
          // Ignore errors
        }
      }
    }

    // Create new actions
    for (const action of editingActions.value) {
      await createActionMutation.mutateAsync({
        path: { id: editingPolicy.value.id },
        body: {
          type: action.type,
          value: action.value,
          order: action.order,
        },
      })
    }
    showActionsModal.value = false
    refetch()
  } catch (error) {
    console.error('Failed to save actions:', error)
  }
}

const policyActions = createPolicyActions(handleEditPolicy, handleEditRule, handleEditActions, handleDeletePolicy)
</script>

<template>
  <div class="policies-settings">
    <div class="card">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold mb-2">Policies</h3>
            <p class="text-muted-color">Configure policies to automatically handle torrent downloads.</p>
          </div>
          <Button
            label="Add Policy"
            :icon="PrimeIcons.PLUS"
            severity="primary"
            @click="handleAddPolicy"
          />
        </div>

        <DataTable
          :data="policies || []"
          :columns="policyColumns"
          :actions="policyActions"
          :loading="isLoading"
          empty-message="No policies configured"
          searchable
          search-placeholder="Search policies..."
          paginator
          :rows="10"
        />
      </div>
    </div>

    <!-- Policy Modal -->
    <Dialog
      v-model:visible="showPolicyModal"
      :header="editingPolicy ? 'Edit Policy' : 'Add Policy'"
      :modal="true"
      :style="{ width: '600px' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Name</label>
          <InputText v-model="policyForm.name" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Description</label>
          <Textarea v-model="policyForm.description" rows="3" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Priority</label>
          <InputNumber v-model="policyForm.priority" :min="0" />
          <span class="text-xs text-muted-color">Higher priority policies are evaluated first</span>
        </div>
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium">Enabled</label>
          <ToggleSwitch v-model="policyForm.enabled" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showPolicyModal = false" />
        <Button
          label="Save"
          :loading="createPolicyMutation.isPending.value || updatePolicyMutation.isPending.value"
          @click="handleSavePolicy"
        />
      </template>
    </Dialog>

    <!-- Rule Modal -->
    <Dialog
      v-model:visible="showRuleModal"
      header="Edit Rule"
      :modal="true"
      :style="{ width: '700px' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Left Operand</label>
          <Select
            v-model="ruleForm.left_operand"
            :options="torrentFieldOptions"
            option-label="label"
            option-value="value"
            placeholder="Select field or enter value"
            :filter="true"
            :show-clear="true"
          />
          <InputText
            v-if="!ruleForm.left_operand || !torrentFieldOptions.find((f) => f.value === ruleForm.left_operand)"
            v-model="ruleForm.left_operand"
            placeholder="Or enter custom value"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Operator</label>
          <Select
            v-model="ruleForm.operator"
            :options="operatorOptions"
            option-label="label"
            option-value="value"
            placeholder="Select operator"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Right Operand</label>
          <InputText v-model="ruleForm.right_operand" placeholder="Enter value" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showRuleModal = false" />
        <Button
          v-if="editingRule"
          label="Delete"
          severity="danger"
          @click="editingPolicy && handleDeleteRule(editingPolicy)"
        />
        <Button label="Save" @click="handleSaveRule" />
      </template>
    </Dialog>

    <!-- Actions Modal -->
    <Dialog
      v-model:visible="showActionsModal"
      header="Edit Actions"
      :modal="true"
      :style="{ width: '800px' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex justify-between items-center">
          <p class="text-sm text-muted-color">Actions are executed in order when the policy matches.</p>
          <Button label="Add Action" size="small" @click="handleAddAction" />
        </div>
        <div v-for="(action, index) in editingActions" :key="index" class="border rounded p-4 flex flex-col gap-3">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">Action {{ index + 1 }}</span>
            <Button
              label="Remove"
              size="small"
              severity="danger"
              variant="text"
              @click="handleRemoveAction(index)"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <label class="text-sm font-medium">Type</label>
              <Select
                v-model="action.type"
                :options="actionTypeOptions"
                option-label="label"
                option-value="value"
                placeholder="Select action type"
                @update:model-value="(val) => { editingActions[index].type = val; editingActions[index].value = '' }"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm font-medium">Value</label>
              <Select
                v-if="action.type && getActionValueOptions(action.type).length > 0"
                v-model="action.value"
                :options="getActionValueOptions(action.type)"
                option-label="label"
                option-value="value"
                placeholder="Select value"
              />
              <InputText
                v-else-if="action.type === 'stop_processing'"
                value="N/A"
                disabled
              />
              <InputText
                v-else
                v-model="action.value"
                placeholder="Enter value"
              />
            </div>
          </div>
        </div>
        <div v-if="editingActions.length === 0" class="text-center py-8 text-muted-color">
          No actions configured. Click "Add Action" to add one.
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showActionsModal = false" />
        <Button label="Save" @click="handleSaveActions" />
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.policies-settings {
  max-width: 100%;
}
</style>

