import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { getV1PoliciesFieldsOptions } from '@/client/@tanstack/vue-query.gen'
import type { ModelFieldDefinition } from '@/client/types.gen'

export interface TemplateVariable {
  /** Template syntax path, e.g., ".Title", ".Quality.Resolution" */
  path: string
  /** Human-readable label */
  label: string
  /** Namespace group: Media, Quality, Candidate, MediaInfo */
  namespace: string
  /** Value type for display hints */
  valueType: string
  /** Whether this is only available post-download */
  postDownloadOnly: boolean
}

export interface TemplateFunction {
  name: string
  label: string
  description: string
}

/** Available template functions */
export const TEMPLATE_FUNCTIONS: TemplateFunction[] = [
  {
    name: 'clean',
    label: 'clean',
    description: 'Sanitizes value and treats "unknown" as empty',
  },
  {
    name: 'sanitize',
    label: 'sanitize',
    description: 'Removes path-unsafe characters',
  },
]

/**
 * Converts an API field path to Go template syntax
 * e.g., "media.title" -> ".Media.Title"
 */
function toTemplatePath(apiPath: string): string {
  return (
    '.' +
    apiPath
      .split('.')
      .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
      .join('.')
  )
}

/**
 * Extracts the namespace from an API field path
 * e.g., "media.title" -> "Media"
 */
function extractNamespace(apiPath: string): string {
  const ns = apiPath.split('.')[0]
  return ns.charAt(0).toUpperCase() + ns.slice(1)
}

/**
 * Composable for accessing template variables from the API
 */
export function useTemplateVariables(options?: { mediaType?: 'movie' | 'series' }) {
  const { data: fields, isLoading, error } = useQuery(getV1PoliciesFieldsOptions())

  /** All variables transformed for template use */
  const allVariables = computed<TemplateVariable[]>(() => {
    if (!fields.value) return []

    return fields.value.map((field: ModelFieldDefinition) => ({
      path: toTemplatePath(field.path),
      label: field.label,
      namespace: extractNamespace(field.path),
      valueType: field.valueType || 'string',
      postDownloadOnly: field.path.startsWith('mediainfo.'),
    }))
  })

  /** Top-level shortcut variables (Title, Year, Season, Episode, EpisodeTitle) */
  const shortcutVariables = computed<TemplateVariable[]>(() => {
    const shortcuts = [
      { path: '.Title', label: 'Title', namespace: 'Shortcuts', valueType: 'string' },
      { path: '.Year', label: 'Year', namespace: 'Shortcuts', valueType: 'string' },
    ]

    if (options?.mediaType === 'series') {
      shortcuts.push(
        { path: '.Season', label: 'Season', namespace: 'Shortcuts', valueType: 'string' },
        { path: '.Episode', label: 'Episode', namespace: 'Shortcuts', valueType: 'string' },
        { path: '.EpisodeTitle', label: 'Episode Title', namespace: 'Shortcuts', valueType: 'string' },
      )
    }

    return shortcuts.map((s) => ({ ...s, postDownloadOnly: false }))
  })

  /** Variables grouped by namespace */
  const variablesByNamespace = computed(() => {
    const groups: Record<string, TemplateVariable[]> = {
      Shortcuts: shortcutVariables.value,
    }

    for (const variable of allVariables.value) {
      // Skip series-only fields for movie templates
      if (options?.mediaType === 'movie') {
        if (
          variable.path === '.Media.Season' ||
          variable.path === '.Media.Episode' ||
          variable.path === '.Media.EpisodeTitle'
        ) {
          continue
        }
      }

      if (!groups[variable.namespace]) {
        groups[variable.namespace] = []
      }
      groups[variable.namespace].push(variable)
    }

    return groups
  })

  /** Flat list of commonly used variables (shortcuts + quality) */
  const commonVariables = computed<TemplateVariable[]>(() => {
    const common = [...shortcutVariables.value]

    // Add commonly used quality fields
    const qualityPaths = ['.Quality.Resolution', '.Quality.Source', '.Quality.Full']
    for (const path of qualityPaths) {
      const variable = allVariables.value.find((v) => v.path === path)
      if (variable) {
        common.push(variable)
      }
    }

    return common
  })

  /** Search variables by label or path */
  const searchVariables = (query: string): TemplateVariable[] => {
    const q = query.toLowerCase()
    const all = [...shortcutVariables.value, ...allVariables.value]

    // Deduplicate by path
    const seen = new Set<string>()
    const unique: TemplateVariable[] = []
    for (const v of all) {
      if (!seen.has(v.path)) {
        seen.add(v.path)
        unique.push(v)
      }
    }

    return unique.filter(
      (v) => v.label.toLowerCase().includes(q) || v.path.toLowerCase().includes(q),
    )
  }

  /** Get variable by template path */
  const getVariableByPath = (path: string): TemplateVariable | undefined => {
    return (
      shortcutVariables.value.find((v) => v.path === path) ||
      allVariables.value.find((v) => v.path === path)
    )
  }

  return {
    allVariables,
    shortcutVariables,
    variablesByNamespace,
    commonVariables,
    searchVariables,
    getVariableByPath,
    functions: TEMPLATE_FUNCTIONS,
    isLoading,
    error,
  }
}


