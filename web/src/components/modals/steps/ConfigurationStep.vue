<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { PrimeIcons } from '@/icons'
import { type JackettIndexerDetails } from '@/client/types.gen'
import { getV1IndexersByIdConfigOptions } from '@/client/@tanstack/vue-query.gen'
import ConfigurationStepField from './ConfigurationStepField.vue'

const props = defineProps<{
  selectedIndexer: JackettIndexerDetails
}>()

const emit = defineEmits<{
  'update:save-data': [data: Record<string, unknown>]
}>()

const formData = ref<Record<string, unknown>>({})

const {
  data: indexerConfig,
  isLoading: loadingConfig,
  error: configError,
} = useQuery(getV1IndexersByIdConfigOptions({ path: { id: props.selectedIndexer.id ?? '' } }))

const updateField = (key: string, value: unknown) => {
  formData.value = { ...formData.value, [key]: value }
  emit('update:save-data', formData.value)
}
</script>

<template>
  <div class="configuration-step">
    <h3 class="text-lg font-semibold mb-4">Configuration</h3>
    <p class="text-muted-color mb-6">
      Configure the specific settings for {{ selectedIndexer?.title }}.
    </p>

    <!-- Loading State -->
    <div v-if="loadingConfig" class="text-center py-8">
      <i :class="PrimeIcons.SPINNER" class="pi pi-spin text-2xl text-muted-color" />
      <p class="mt-2 text-muted-color">Loading configuration fields...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="configError" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex items-start">
        <i :class="PrimeIcons.EXCLAMATION_TRIANGLE" class="text-red-600 text-xl mr-3 mt-0.5" />
        <div>
          <h4 class="font-semibold text-red-800">Configuration Error</h4>
          <p class="text-red-700 text-sm mt-1">
            Failed to load configuration fields: {{ configError.message }}
          </p>
        </div>
      </div>
    </div>

    <!-- Configuration Fields -->
    <div v-else-if="indexerConfig" class="space-y-4">
      <ConfigurationStepField
        v-for="field in indexerConfig"
        :key="field.id"
        :field="field"
        :model-value="
          (formData as Record<string, unknown>)[field.id as string] as string | undefined
        "
        @update:model-value="(value) => updateField(field.id as string, value)"
      />
    </div>

    <!-- No Configuration -->
    <div v-else class="bg-blue-50 border border-blue-200 rounded-lg p-4">
      <div class="flex items-start">
        <i :class="PrimeIcons.INFO_CIRCLE" class="text-blue-600 text-xl mr-3 mt-0.5" />
        <div>
          <h4 class="font-semibold text-blue-800">No Configuration Required</h4>
          <p class="text-blue-700 text-sm mt-1">
            This indexer doesn't require any additional configuration fields.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.configuration-step {
  /* Height is now managed by parent scroll container */
}
</style>
