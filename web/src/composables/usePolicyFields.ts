import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { getV1PoliciesFieldsOptions } from '@/client/@tanstack/vue-query.gen'
import type { ModelFieldDefinition } from '@/client/types.gen'

export function usePolicyFields() {
  const { data: fields, isLoading, error } = useQuery(getV1PoliciesFieldsOptions())

  // Group fields by category (candidate vs quality)
  const candidateFields = computed(() => {
    return fields.value?.filter((f) => f.path.startsWith('candidate.')) || []
  })

  const qualityFields = computed(() => {
    return fields.value?.filter((f) => f.path.startsWith('quality.')) || []
  })

  // Get field definition by path
  const getFieldByPath = (path: string): ModelFieldDefinition | undefined => {
    return fields.value?.find((f) => f.path === path)
  }

  // Get valid operators for a field
  const getValidOperators = (field: ModelFieldDefinition | undefined) => {
    if (!field) return []
    return field.operators || []
  }

  return {
    fields,
    candidateFields,
    qualityFields,
    isLoading,
    error,
    getFieldByPath,
    getValidOperators,
  }
}

