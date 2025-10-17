<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Steps from 'primevue/steps'
import { PrimeIcons } from '@/icons'
import { type ModelIndexerDefinition, type ModelIndexerInput } from '@/client/types.gen'
import { postV1IndexerMutation } from '@/client/@tanstack/vue-query.gen'
import SelectIndexerTypeStep from './steps/SelectIndexerTypeStep.vue'
import ConfigurationStep from './steps/ConfigurationStep.vue'
import ReviewStep from './steps/ReviewStep.vue'

const emit = defineEmits<{
  close: []
  'indexer-added': []
}>()

// Form state
const currentStep = ref(0)
const selectedIndexerType = ref<ModelIndexerDefinition | null>(null)
const saveData = ref<ModelIndexerInput | undefined>(undefined)

const createIndexerMutation = useMutation({
  ...postV1IndexerMutation(),
  onSuccess: () => {
    emit('indexer-added')
    closeModal()
  },
  onError: (error) => {
    console.error('Failed to create indexer:', error)
  },
})

// Computed properties
const steps = computed(() => [
  { label: 'Select Type' },
  { label: 'Configuration' },
  { label: 'Review' },
])

const isLastStep = computed(() => currentStep.value === steps.value.length - 1)
const isFirstStep = computed(() => currentStep.value === 0)

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 0:
      return selectedIndexerType.value !== null
    case 1:
      return true // Configuration step validation can be added later
    case 2:
      return true
    default:
      return false
  }
})

// Methods
const nextStep = () => {
  if (canProceed.value && !isLastStep.value) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (!isFirstStep.value) {
    currentStep.value--
  }
}

const selectIndexerType = (indexer: ModelIndexerDefinition) => {
  selectedIndexerType.value = indexer
}

const createIndexer = () => {
  if (canProceed.value && selectedIndexerType.value && saveData.value) {
    createIndexerMutation.mutate({
      body: saveData.value,
    })
  }
}

const closeModal = () => {
  emit('close')
}

const handleUpdateVisible = (visible: boolean) => {
  if (!visible) {
    closeModal()
  }
}
</script>

<template>
  <Dialog
    :visible="true"
    :modal="true"
    :closable="true"
    :dismissable-mask="true"
    header="Add New Indexer"
    class="w-full max-w-4xl"
    @update:visible="handleUpdateVisible"
  >
    <div class="add-indexer-modal overflow-hidden h-full">
      <!-- Progress Steps -->
      <div class="mb-6">
        <Steps :model="steps" :active-index="currentStep" />
      </div>

      <!-- Step Content -->
      <div class="max-h-[calc(100vh*0.6)] overflow-y-auto">
        <!-- Step 1: Select Indexer Type -->
        <div v-if="currentStep === 0" class="step-1 h-full overflow-hidden">
          <SelectIndexerTypeStep
            ref="selectStepRef"
            :selected-indexer="selectedIndexerType"
            @indexer-selected="selectIndexerType"
          />
        </div>

        <!-- Step 2: Configuration -->
        <div v-if="currentStep === 1 && selectedIndexerType" class="step-2 overflow-hidden">
          <ConfigurationStep v-model="saveData" :selected-indexer="selectedIndexerType" />
        </div>

        <!-- Step 4: Review -->
        <div v-if="currentStep === 3" class="step-4">
          <ReviewStep :selected-indexer="selectedIndexerType" :save-data="saveData" />
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="flex justify-between items-center pt-6">
        <Button
          :label="isFirstStep ? 'Cancel' : 'Previous'"
          :icon="isFirstStep ? PrimeIcons.TIMES : PrimeIcons.ANGLE_LEFT"
          severity="secondary"
          variant="outlined"
          @click="isFirstStep ? closeModal() : prevStep()"
        />

        <div class="flex gap-2">
          <Button
            v-if="!isLastStep"
            label="Next"
            :icon="PrimeIcons.ANGLE_RIGHT"
            icon-pos="right"
            :disabled="!canProceed"
            @click="nextStep"
          />
          <Button
            v-else
            label="Create Indexer"
            :icon="PrimeIcons.PLUS"
            :loading="createIndexerMutation.isPending.value"
            :disabled="!canProceed"
            @click="createIndexer"
          />
        </div>
      </div>
    </div>
  </Dialog>
</template>

<style scoped></style>
