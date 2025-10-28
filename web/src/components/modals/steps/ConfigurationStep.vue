<script setup lang="ts">
import { onMounted } from 'vue'
import {
  type ModelIndexerDefinition,
  type ModelIndexerInput,
  type ModelProtocol,
} from '@/client/types.gen'
import { cloneDeep } from '@/utils'
import ConfigurationStepField from './ConfigurationStepField.vue'

const model = defineModel<ModelIndexerInput | undefined>(undefined)

const props = defineProps<{
  selectedIndexer: ModelIndexerDefinition
}>()

onMounted(() => {
  const copy = cloneDeep(props.selectedIndexer)

  model.value = {
    enable: copy.enable,
    redirect: copy.redirect,
    priority: copy.priority,
    appProfileId: 1,
    configContract: copy.configContract,
    implementation: copy.implementation,
    name: copy.name,
    protocol: copy.protocol as ModelProtocol,
    // tags: copy.tags,
    fields: copy.fields,
  }
})

const handleValueChange = (fieldName: string, value: unknown) => {
  if (model.value) {
    const index = model.value.fields.findIndex((field) => field.name === fieldName)
    if (index !== -1) {
      model.value.fields[index]!.value = value
    }
  }
}
</script>

<template>
  <div class="configuration-step h-full overflow-y-auto">
    <h3 class="text-lg font-semibold mb-4">Configuration</h3>
    <p class="text-muted-color mb-6">
      Configure the specific settings for {{ selectedIndexer?.name }}.
    </p>

    <!-- Configuration Fields -->
    <div class="space-y-4">
      <ConfigurationStepField
        v-for="field in selectedIndexer.fields"
        :key="field.name"
        :field="field"
        :selected-indexer="selectedIndexer"
        @value-change="handleValueChange"
      />
    </div>
  </div>
</template>
