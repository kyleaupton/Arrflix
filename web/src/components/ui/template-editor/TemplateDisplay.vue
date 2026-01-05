<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'

interface Props {
  /** The template string to display */
  template: string
  /** For series: array of [show, season, episode] templates */
  seriesTemplates?: [string, string, string]
  /** Whether this is a series template */
  isSeries?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isSeries: false,
})

interface ParsedSegment {
  type: 'text' | 'variable'
  value: string
  func?: string
}

/**
 * Parse a template string into segments of text and variables
 */
function parseTemplate(template: string): ParsedSegment[] {
  const segments: ParsedSegment[] = []
  // Match {{func .Variable}} or {{.Variable}}
  const regex = /\{\{(\w+\s+)?(\.[^}]+)\}\}/g
  let lastIndex = 0
  let match

  while ((match = regex.exec(template)) !== null) {
    // Add text before this match
    if (match.index > lastIndex) {
      segments.push({
        type: 'text',
        value: template.slice(lastIndex, match.index),
      })
    }

    // Add the variable
    const func = match[1]?.trim()
    const variable = match[2] || ''
    segments.push({
      type: 'variable',
      value: variable,
      func: func || undefined,
    })

    lastIndex = regex.lastIndex
  }

  // Add remaining text
  if (lastIndex < template.length) {
    segments.push({
      type: 'text',
      value: template.slice(lastIndex),
    })
  }

  return segments
}

/**
 * Format a variable for display: { Title } or { func Title }
 */
function formatVariable(segment: ParsedSegment): string {
  const displayValue = segment.value.startsWith('.') ? segment.value.slice(1) : segment.value
  if (segment.func) {
    return `${segment.func} ${displayValue}`
  }
  return displayValue
}

const parsedTemplate = computed(() => parseTemplate(props.template))

const parsedSeriesTemplates = computed(() => {
  if (!props.seriesTemplates) return null
  return props.seriesTemplates.map((t) => parseTemplate(t))
})
</script>

<template>
  <div class="inline-flex items-center flex-wrap gap-0.5 font-mono text-sm">
    <!-- Series: show all three templates on one line with / separators -->
    <template v-if="isSeries && parsedSeriesTemplates">
      <!-- Show template -->
      <template v-for="(segment, idx) in parsedSeriesTemplates[0]" :key="`show-${idx}`">
        <Badge v-if="segment.type === 'variable'" variant="default" class="font-mono text-xs">
          {{ formatVariable(segment) }}
        </Badge>
        <span v-else class="whitespace-pre">{{ segment.value }}</span>
      </template>
      <span class="mx-1 text-muted-foreground">/</span>

      <!-- Season template -->
      <template v-for="(segment, idx) in parsedSeriesTemplates[1]" :key="`season-${idx}`">
        <Badge v-if="segment.type === 'variable'" variant="default" class="font-mono text-xs">
          {{ formatVariable(segment) }}
        </Badge>
        <span v-else class="whitespace-pre">{{ segment.value }}</span>
      </template>
      <span class="mx-1 text-muted-foreground">/</span>

      <!-- Episode template -->
      <template v-for="(segment, idx) in parsedSeriesTemplates[2]" :key="`episode-${idx}`">
        <Badge v-if="segment.type === 'variable'" variant="default" class="font-mono text-xs">
          {{ formatVariable(segment) }}
        </Badge>
        <span v-else class="whitespace-pre">{{ segment.value }}</span>
      </template>
    </template>

    <!-- Single template (movies or fallback) -->
    <template v-else>
      <template v-for="(segment, idx) in parsedTemplate" :key="idx">
        <Badge v-if="segment.type === 'variable'" variant="default" class="font-mono text-xs">
          {{ formatVariable(segment) }}
        </Badge>
        <span v-else class="whitespace-pre">{{ segment.value }}</span>
      </template>
    </template>
  </div>
</template>
