<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useMutation, useQuery } from '@tanstack/vue-query'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Steps from 'primevue/steps'
import { PrimeIcons } from '@/icons'
import { type HandlersIndexerCreateRequest, type JackettIndexerConfig } from '@/client/types.gen'
import {
  getV1IndexersUnconfiguredOptions,
  // postV1IndexersOptions,
  postV1IndexersMutation,
} from '@/client/@tanstack/vue-query.gen'
import DataTable from '@/components/tables/DataTable.vue'
import {
  availableIndexerColumns,
  createAvailableIndexerActions,
} from '@/components/tables/configs/availableIndexerTableConfig'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'indexer-added': [indexer: JackettIndexerConfig]
}>()

// Form state
const currentStep = ref(0)
const selectedIndexerType = ref<JackettIndexerConfig | null>(null)
const formData = ref<Partial<HandlersIndexerCreateRequest>>({
  name: '',
  description: '',
  enabled: true,
  type: '',
  fields: {},
})

// API queries
const { data: availableIndexers, isLoading: loadingIndexers } = useQuery(
  getV1IndexersUnconfiguredOptions(),
)

const createIndexerMutation = useMutation({
  ...postV1IndexersMutation(),
  onSuccess: (data) => {
    emit('indexer-added', data)
    closeModal()
  },
  onError: (error) => {
    console.error('Failed to create indexer:', error)
  },
})

// Computed properties
const steps = computed(() => [
  { label: 'Select Type' },
  { label: 'Basic Info' },
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
      return formData.value.name && formData.value.description
    case 2:
      return true // Configuration step validation can be added later
    case 3:
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

const selectIndexerType = (indexer: JackettIndexerConfig) => {
  selectedIndexerType.value = indexer
  formData.value.type = indexer.type || ''
  formData.value.name = indexer.name || ''
  formData.value.description = indexer.description || ''
}

// Create actions for the available indexers table
const availableIndexerActions = createAvailableIndexerActions(selectIndexerType)

const createIndexer = () => {
  if (canProceed.value) {
    createIndexerMutation.mutate({
      body: formData.value as HandlersIndexerCreateRequest,
    })
  }
}

const closeModal = () => {
  emit('update:visible', false)
  // Reset form
  currentStep.value = 0
  selectedIndexerType.value = null
  formData.value = {
    name: '',
    description: '',
    enabled: true,
    type: '',
    fields: {},
  }
}

// Watch for modal visibility changes
watch(
  () => props.visible,
  (visible) => {
    if (!visible) {
      closeModal()
    }
  },
)
</script>

<template>
  <Dialog
    :visible="visible"
    :modal="true"
    :closable="true"
    :dismissable-mask="true"
    header="Add New Indexer"
    class="w-full max-w-4xl"
    @update:visible="emit('update:visible', $event)"
  >
    <div class="add-indexer-modal">
      <!-- Progress Steps -->
      <div class="mb-6">
        <Steps :model="steps" :active-index="currentStep" />
      </div>

      <!-- Step Content -->
      <div class="step-content min-h-96">
        <!-- Step 1: Select Indexer Type -->
        <div v-if="currentStep === 0" class="step-1">
          <h3 class="text-lg font-semibold mb-4">Select Indexer Type</h3>
          <p class="text-muted-color mb-6">Choose from available indexer types to configure.</p>

          <div v-if="loadingIndexers" class="text-center py-8">
            <i :class="PrimeIcons.SPINNER" class="pi pi-spin text-2xl text-muted-color" />
            <p class="mt-2 text-muted-color">Loading available indexers...</p>
          </div>

          <div v-else-if="availableIndexers && availableIndexers.length > 0">
            <DataTable
              :data="availableIndexers"
              :columns="availableIndexerColumns"
              :actions="availableIndexerActions"
              :loading="loadingIndexers"
              empty-message="No unconfigured indexers available"
              searchable
              search-placeholder="Search available indexers..."
              paginator
              :rows="8"
              selectable
              selection-mode="single"
              @selection-change="
                (selection) => {
                  if (selection && !Array.isArray(selection)) {
                    selectIndexerType(selection)
                  }
                }
              "
            />
          </div>

          <div v-else class="text-center py-8 text-muted-color">
            <i :class="PrimeIcons.INFO_CIRCLE" class="text-4xl mb-4" />
            <p>No unconfigured indexers available.</p>
          </div>
        </div>

        <!-- Step 2: Basic Information -->
        <div v-if="currentStep === 1" class="step-2">
          <h3 class="text-lg font-semibold mb-4">Basic Information</h3>
          <p class="text-muted-color mb-6">Configure the basic settings for your indexer.</p>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium mb-2">Name</label>
                <input
                  v-model="formData.name"
                  type="text"
                  class="w-full p-inputtext p-component"
                  placeholder="Enter indexer name"
                />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2">Description</label>
                <textarea
                  v-model="formData.description"
                  class="w-full p-inputtextarea p-component"
                  rows="3"
                  placeholder="Enter indexer description"
                />
              </div>
            </div>

            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium mb-2">Type</label>
                <input
                  :value="formData.type"
                  type="text"
                  class="w-full p-inputtext p-component"
                  disabled
                />
              </div>

              <div class="flex items-center space-x-2">
                <input v-model="formData.enabled" type="checkbox" class="p-checkbox" id="enabled" />
                <label for="enabled" class="text-sm font-medium">Enable this indexer</label>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 3: Configuration -->
        <div v-if="currentStep === 2" class="step-3">
          <h3 class="text-lg font-semibold mb-4">Configuration</h3>
          <p class="text-muted-color mb-6">
            Configure the specific settings for {{ selectedIndexerType?.name }}.
          </p>

          <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div class="flex items-start">
              <i
                :class="PrimeIcons.EXCLAMATION_TRIANGLE"
                class="text-yellow-600 text-xl mr-3 mt-0.5"
              />
              <div>
                <h4 class="font-semibold text-yellow-800">Configuration Fields</h4>
                <p class="text-yellow-700 text-sm mt-1">
                  Configuration fields will be dynamically generated based on the selected indexer
                  type. This step will be implemented once the API provides field definitions.
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 4: Review -->
        <div v-if="currentStep === 3" class="step-4">
          <h3 class="text-lg font-semibold mb-4">Review & Create</h3>
          <p class="text-muted-color mb-6">Review your indexer configuration before creating.</p>

          <div class="bg-gray-50 rounded-lg p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 class="font-semibold mb-3">Basic Information</h4>
                <dl class="space-y-2">
                  <div>
                    <dt class="text-sm font-medium text-muted-color">Name</dt>
                    <dd class="text-sm">{{ formData.name }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-muted-color">Description</dt>
                    <dd class="text-sm">{{ formData.description }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-muted-color">Type</dt>
                    <dd class="text-sm">{{ formData.type }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-muted-color">Status</dt>
                    <dd class="text-sm">
                      <span
                        class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                        :class="
                          formData.enabled
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        "
                      >
                        {{ formData.enabled ? 'Enabled' : 'Disabled' }}
                      </span>
                    </dd>
                  </div>
                </dl>
              </div>

              <div>
                <h4 class="font-semibold mb-3">Configuration</h4>
                <p class="text-sm text-muted-color">
                  Configuration fields will be displayed here once implemented.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="flex justify-between items-center mt-8 pt-6 border-t border-surface">
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

<style scoped>
.add-indexer-modal {
  min-height: 500px;
}

.step-content {
  min-height: 400px;
}

:deep(.p-steps .p-steps-item.p-highlight .p-steps-number) {
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}

:deep(.p-dialog-content) {
  padding: 2rem;
}
</style>
