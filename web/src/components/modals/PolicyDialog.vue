<script setup lang="ts">
import { ref, inject, watch, computed } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { postV1PoliciesMutation, putV1PoliciesByIdMutation } from '@/client/@tanstack/vue-query.gen'
import { type DbgenPolicy } from '@/client/types.gen'
import BaseDialog from './BaseDialog.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'

interface Props {
  policy?: DbgenPolicy | null
}

const props = defineProps<Props>()

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }

const createPolicyMutation = useMutation(postV1PoliciesMutation())
const updatePolicyMutation = useMutation(putV1PoliciesByIdMutation())

const policyForm = ref({
  name: '',
  description: '',
  enabled: true,
  priority: 0,
})

// Initialize form when policy changes
watch(
  () => props.policy,
  (policy) => {
    if (policy) {
      policyForm.value = {
        name: policy.name || '',
        description: policy.description || '',
        enabled: policy.enabled ?? true,
        priority: policy.priority || 0,
      }
    } else {
      policyForm.value = { name: '', description: '', enabled: true, priority: 0 }
    }
  },
  { immediate: true },
)

const handleSave = async () => {
  try {
    if (props.policy?.id) {
      await updatePolicyMutation.mutateAsync({
        path: { id: props.policy.id },
        body: {
          name: policyForm.value.name,
          description: policyForm.value.description || '',
          enabled: policyForm.value.enabled,
          priority: policyForm.value.priority,
        },
      })
    } else {
      await createPolicyMutation.mutateAsync({
        body: {
          name: policyForm.value.name,
          description: policyForm.value.description || '',
          enabled: policyForm.value.enabled,
          priority: policyForm.value.priority,
        },
      })
    }
    dialogRef.value.close({ saved: true })
  } catch (error) {
    console.error('Failed to save policy:', error)
  }
}

const handleCancel = () => {
  dialogRef.value.close()
}

const isLoading = computed(
  () => createPolicyMutation.isPending.value || updatePolicyMutation.isPending.value,
)
</script>

<template>
  <BaseDialog :title="policy ? 'Edit Policy' : 'Add Policy'">
    <div class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="policy-name">Name</Label>
        <Input id="policy-name" v-model="policyForm.name" />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="policy-description">Description</Label>
        <Textarea id="policy-description" v-model="policyForm.description" rows="3" />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="policy-priority">Priority</Label>
        <Input id="policy-priority" type="number" v-model.number="policyForm.priority" :min="0" />
        <span class="text-xs text-muted-foreground"
          >Higher priority policies are evaluated first</span
        >
      </div>
      <div class="flex items-center justify-between">
        <Label for="policy-enabled">Enabled</Label>
        <Switch id="policy-enabled" v-model="policyForm.enabled" />
      </div>
    </div>
    <template #footer>
      <Button variant="outline" @click="handleCancel">Cancel</Button>
      <Button :disabled="isLoading" @click="handleSave">Save</Button>
    </template>
  </BaseDialog>
</template>
