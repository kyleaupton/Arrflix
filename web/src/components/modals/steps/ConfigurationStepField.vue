<script setup lang="ts">
import { computed } from 'vue'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Message from 'primevue/message'
import Checkbox from 'primevue/checkbox'
import { type JackettIndexerConfigField } from '@/client/types.gen'

const model = defineModel<string>()

const props = defineProps<{
  field: JackettIndexerConfigField
}>()

const options = computed(() => {
  if (props.field.type === 'inputselect') {
    return Object.entries(props.field.options as Record<string, string>).map(([key, value]) => ({
      label: value,
      value: key,
    }))
  }
  return []
})
</script>

<template>
  <template v-if="field.type === 'inputstring'">
    <label class="block text-sm font-medium mb-2">{{ field.name }}</label>
    <InputText
      v-model="model"
      :placeholder="`Enter ${field.name}`"
      variant="filled"
      class="w-full"
    />
  </template>

  <template v-else-if="field.type === 'inputselect'">
    <label class="block text-sm font-medium mb-2">{{ field.name }}</label>
    <Select
      v-model="model"
      :options="options"
      optionLabel="label"
      optionValue="value"
      :placeholder="`Select ${field.name}`"
      variant="filled"
      class="w-full"
    />
  </template>

  <div v-else-if="field.type === 'inputbool'" class="flex items-center gap-2">
    <Checkbox v-model="model" :inputId="field.id" :binary="true" />
    <label :for="field.id" class="text-sm font-medium">{{ field.name }}</label>
  </div>

  <template v-else-if="field.type === 'displayinfo'">
    <Message>
      <div v-html="field.value" />
    </Message>
  </template>
</template>

<style scoped></style>
