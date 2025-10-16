<script setup lang="ts">
import { type ModelIndexerDefinition } from '@/client/types.gen'
import ConfigurationStepField from './ConfigurationStepField.vue'

const model = defineModel<Record<string, unknown>>({ required: true })

defineProps<{
  selectedIndexer: ModelIndexerDefinition
}>()
</script>

<template>
  <div class="configuration-step">
    <h3 class="text-lg font-semibold mb-4">Configuration</h3>
    <p class="text-muted-color mb-6">
      Configure the specific settings for {{ selectedIndexer?.name }}.
    </p>

    <!-- Configuration Fields -->
    <div class="space-y-4">
      <ConfigurationStepField
        v-for="field in selectedIndexer.fields"
        v-model="model[field.name]"
        :key="field.name"
        :field="field"
      />
    </div>
  </div>
</template>
