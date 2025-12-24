<script setup lang="ts">
import { ref, computed, inject } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { ChevronLeft, ChevronRight, Plus, X } from 'lucide-vue-next'
import { type ModelIndexerDefinition } from '@/client/types.gen'
import { postV1IndexerMutation } from '@/client/@tanstack/vue-query.gen'
import BaseDialog from './BaseDialog.vue'
import { Button } from '@/components/ui/button'
import SelectIndexerTypeStep from './steps/SelectIndexerTypeStep.vue'
import ConfigurationStep from './steps/ConfigurationStep.vue'
import ReviewStep from './steps/ReviewStep.vue'

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }

// Form state
const currentStep = ref(0)
const selectedIndexerType = ref<ModelIndexerDefinition | null>(null)
const saveData = ref<ModelIndexerDefinition | undefined>(undefined)

const createIndexerMutation = useMutation({
  ...postV1IndexerMutation(),
  onSuccess: () => {
    dialogRef.value.close({ indexerAdded: true })
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
  dialogRef.value.close()
}
</script>

<template>
  <BaseDialog title="Add New Indexer">
    <!-- Progress Steps -->
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <div v-for="(step, index) in steps" :key="index" class="flex items-center flex-1">
          <div class="flex items-center flex-1">
            <!-- Step Circle -->
            <div
              class="flex items-center justify-center size-8 rounded-full border-2 transition-colors"
              :class="
                index < currentStep
                  ? 'bg-primary border-primary text-primary-foreground'
                  : index === currentStep
                    ? 'bg-primary border-primary text-primary-foreground'
                    : 'bg-background border-muted-foreground text-muted-foreground'
              "
            >
              <span v-if="index < currentStep" class="text-sm font-semibold">✓</span>
              <span v-else class="text-sm font-semibold">{{ index + 1 }}</span>
            </div>
            <!-- Step Label -->
            <div class="ml-3 flex-1">
              <div
                class="text-sm font-medium"
                :class="index <= currentStep ? 'text-foreground' : 'text-muted-foreground'"
              >
                {{ step.label }}
              </div>
            </div>
          </div>
          <!-- Connector Line -->
          <div
            v-if="index < steps.length - 1"
            class="flex-1 h-0.5 mx-4 transition-colors"
            :class="index < currentStep ? 'bg-primary' : 'bg-muted'"
          />
        </div>
      </div>
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

      <!-- Step 3: Review -->
      <div v-if="currentStep === 2" class="step-3">
        <ReviewStep :selected-indexer="selectedIndexerType" :save-data="saveData" />
      </div>
    </div>

    <template #footer>
      <div class="flex justify-between items-center w-full">
        <Button variant="outline" @click="isFirstStep ? closeModal() : prevStep()">
          <X v-if="isFirstStep" class="mr-2 size-4" />
          <ChevronLeft v-else class="mr-2 size-4" />
          {{ isFirstStep ? 'Cancel' : 'Previous' }}
        </Button>

        <div class="flex gap-2">
          <Button v-if="!isLastStep" :disabled="!canProceed" @click="nextStep">
            Next
            <ChevronRight class="ml-2 size-4" />
          </Button>
          <Button
            v-else
            :disabled="!canProceed || createIndexerMutation.isPending.value"
            @click="createIndexer"
          >
            <Plus class="mr-2 size-4" />
            Create Indexer
          </Button>
        </div>
      </div>
    </template>
  </BaseDialog>
</template>
