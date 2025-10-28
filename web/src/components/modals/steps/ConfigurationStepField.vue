<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import MultiSelect from 'primevue/multiselect'
import Password from 'primevue/password'
import Message from 'primevue/message'
import Checkbox from 'primevue/checkbox'
import { type ModelIndexerDefinition, type ModelIndexerField } from '@/client/types.gen'
import { postV1IndexerActionByNameMutation } from '@/client/@tanstack/vue-query.gen'

const model = computed({
  get: () => props.field.value,
  set: (value) => {
    emit('value-change', props.field.name, value)
  },
})

const emit = defineEmits<{
  (e: 'value-change', fieldName: string, value: unknown): void
}>()

const props = defineProps<{
  selectedIndexer: ModelIndexerDefinition
  field: ModelIndexerField
}>()

const options = computed(() => {
  if (props.field.selectOptions) {
    return props.field.selectOptions.map((option) => ({
      label: option.name,
      value: option.value,
      hint: option.hint,
    }))
  } else if (isAsyncAction.value && actionMutation.data?.value) {
    // @ts-expect-error shit ain't got types
    return actionMutation.data.value.options || []
  }

  return []
})

const isMultiSelect = computed(() => {
  // Check if the field should support multiple selections
  // This could be based on field name or other properties
  return (
    props.field.name.includes('searchTypes') ||
    props.field.name.includes('languagesOnly') ||
    Array.isArray(model.value)
  )
})

const hasHelpText = computed(() => {
  return props.field.helpText || props.field.helpTextWarning
})

const helpSeverity = computed(() => {
  return props.field.helpTextWarning ? 'warning' : 'secondary'
})

const fieldId = computed(() => `field-${props.field.name}`)

const isAsyncAction = computed(() => {
  return props.field.selectOptionsProviderAction
})

const selectOptionLabel = computed(() => (isAsyncAction.value ? 'name' : 'label'))

const actionMutation = useMutation({
  ...postV1IndexerActionByNameMutation(),
  onSuccess: (data) => {
    console.log('Action performed successfully', data)
  },
  onError: (error) => {
    console.error('Failed to perform action:', error)
  },
})

const performAction = () => {
  if (props.field.selectOptionsProviderAction) {
    actionMutation.mutate({
      path: { name: props.field.selectOptionsProviderAction },
      // @ts-expect-error todo: fix the type here
      body: props.selectedIndexer,
    })
  }
}

onMounted(() => {
  performAction()
})
</script>

<template>
  <div
    class="field-container"
    :class="{ hidden: field.hidden === 'hidden', 'advanced-field': field.advanced }"
  >
    <!-- Text Input Fields -->
    <template v-if="field.type === 'textbox'">
      <label :for="fieldId" class="block text-sm font-medium mb-2">
        {{ field.label }}
        <span v-if="field.advanced" class="text-xs text-gray-500 ml-1">(Advanced)</span>
      </label>
      <InputText
        :id="fieldId"
        v-model="model as string"
        :placeholder="`Enter ${field.label}`"
        variant="filled"
        class="w-full"
        :class="{ 'p-invalid': field.helpTextWarning }"
      />
    </template>

    <!-- Password Input Fields -->
    <template v-else-if="field.type === 'password'">
      <label :for="fieldId" class="block text-sm font-medium mb-2">
        {{ field.label }}
        <span v-if="field.advanced" class="text-xs text-gray-500 ml-1">(Advanced)</span>
      </label>
      <Password
        :id="fieldId"
        v-model="model as string"
        :placeholder="`Enter ${field.label}`"
        variant="filled"
        class="w-full"
        :class="{ 'p-invalid': field.helpTextWarning }"
        :feedback="false"
        toggleMask
      />
    </template>

    <!-- Number Input Fields -->
    <template v-else-if="field.type === 'number'">
      <label :for="fieldId" class="block text-sm font-medium mb-2">
        {{ field.label }}
        <span v-if="field.unit" class="text-xs text-gray-500 ml-1">({{ field.unit }})</span>
        <span v-if="field.advanced" class="text-xs text-gray-500 ml-1">(Advanced)</span>
      </label>
      <InputNumber
        :id="fieldId"
        v-model="model as number"
        :placeholder="`Enter ${field.label}`"
        variant="filled"
        class="w-full"
        :class="{ 'p-invalid': field.helpTextWarning }"
        :minFractionDigits="field.isFloat ? 2 : 0"
        :maxFractionDigits="field.isFloat ? 2 : 0"
        :useGrouping="false"
      />
    </template>

    <!-- Select Fields -->
    <template v-else-if="field.type === 'select'">
      <label :for="fieldId" class="block text-sm font-medium mb-2">
        {{ field.label }}
        <span v-if="field.advanced" class="text-xs text-gray-500 ml-1">(Advanced)</span>
      </label>
      <MultiSelect
        v-if="isMultiSelect"
        :id="fieldId"
        v-model="model"
        :options="options"
        :optionLabel="selectOptionLabel"
        optionValue="value"
        :placeholder="`Select ${field.label}`"
        variant="filled"
        class="w-full"
        :class="{ 'p-invalid': field.helpTextWarning }"
        :maxSelectedLabels="3"
        selectedItemsLabel="{0} items selected"
      />
      <Select
        v-else
        :id="fieldId"
        v-model="model"
        :options="options"
        :optionLabel="selectOptionLabel"
        optionValue="value"
        :placeholder="`Select ${field.label}`"
        variant="filled"
        class="w-full"
        :class="{ 'p-invalid': field.helpTextWarning }"
      />
    </template>

    <!-- Checkbox Fields -->
    <template v-else-if="field.type === 'checkbox'">
      <div class="flex items-center gap-2">
        <Checkbox
          :id="fieldId"
          v-model="model"
          :binary="true"
          :class="{ 'p-invalid': field.helpTextWarning }"
        />
        <label :for="fieldId" class="text-sm font-medium cursor-pointer">
          {{ field.label }}
          <span v-if="field.advanced" class="text-xs text-gray-500 ml-1">(Advanced)</span>
        </label>
      </div>
    </template>

    <!-- Fallback for unknown field types -->
    <template v-else>
      <div class="p-4 border border-yellow-200 bg-yellow-50 rounded">
        <p class="text-sm text-yellow-800"><strong>Unknown field type:</strong> {{ field.type }}</p>
        <p class="text-sm text-yellow-700 mt-1">Field: {{ field.name }} ({{ field.label }})</p>
      </div>
    </template>

    <Message v-if="hasHelpText" class="mt-1" size="small" :severity="helpSeverity" variant="simple"
      >{{ field.helpText || field.helpTextWarning }}
      <a
        v-if="field.helpLink"
        :href="field.helpLink"
        target="_blank"
        class="ml-1 text-blue-600 hover:underline"
      >
        Learn more
      </a>
    </Message>
  </div>
</template>
