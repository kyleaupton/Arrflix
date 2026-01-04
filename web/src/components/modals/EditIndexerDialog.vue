<script setup lang="ts">
import { ref, inject, computed, watch } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { Save } from 'lucide-vue-next'
import { type ModelIndexerOutput, type ModelIndexerInput } from '@/client/types.gen'
import { postV1IndexerMutation } from '@/client/@tanstack/vue-query.gen'
import { Button } from '@/components/ui/button'
import ConfigurationStep from './steps/ConfigurationStep.vue'
import BaseDialog from './BaseDialog.vue'

interface Props {
  indexer: ModelIndexerOutput
}

const props = defineProps<Props>()

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }

// Form state
const saveData = ref<ModelIndexerInput | undefined>(undefined)
const indexerError = ref<string | null>(null)

const updateIndexerMutation = useMutation({
  ...postV1IndexerMutation(),
  onSuccess: () => {
    indexerError.value = null
    dialogRef.value.close({ indexerUpdated: true })
  },
  onError: (error) => {
    console.error('Failed to update indexer:', error)
    const err = error as { message?: string }
    indexerError.value = err.message || 'Failed to update indexer'
  },
})

// Watch for changes to the indexer prop and reset error
watch(
  () => props.indexer,
  () => {
    indexerError.value = null
  },
  { immediate: true },
)

const handleSave = () => {
  if (!saveData.value) {
    indexerError.value = 'Configuration data is required'
    return
  }

  // Include the ID to indicate this is an update operation
  const updatePayload: ModelIndexerInput = {
    ...saveData.value,
    id: props.indexer.id,
  }

  updateIndexerMutation.mutate({
    body: updatePayload,
  })
}

const handleCancel = () => {
  dialogRef.value.close()
}

const canSave = computed(() => {
  return saveData.value !== undefined
})
</script>

<template>
  <BaseDialog :title="`Edit ${indexer.name}`">
    <div class="flex flex-col gap-4">
      <div
        v-if="indexerError"
        class="p-4 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive text-sm"
      >
        {{ indexerError }}
      </div>

      <div class="max-h-[calc(100vh*0.6)] overflow-y-auto">
        <ConfigurationStep v-model="saveData" :selected-indexer="indexer" />
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 w-full">
        <Button variant="outline" @click="handleCancel">Cancel</Button>
        <Button
          :disabled="!canSave || updateIndexerMutation.isPending.value"
          @click="handleSave"
        >
          <Save class="mr-2 size-4" />
          Save Changes
        </Button>
      </div>
    </template>
  </BaseDialog>
</template>
