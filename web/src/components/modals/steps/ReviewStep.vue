<script setup lang="ts">
import { type ModelIndexerDefinition } from '@/client/types.gen'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

defineProps<{
  selectedIndexer: ModelIndexerDefinition | null
  saveData: unknown
}>()
</script>

<template>
  <div class="review-step h-full overflow-y-auto">
    <div class="space-y-1 mb-6">
      <h3 class="text-lg font-semibold">Review & Create</h3>
      <p class="text-sm text-muted-foreground">
        Review your indexer configuration before creating.
      </p>
    </div>

    <Card>
      <CardContent class="pt-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-4">
            <CardHeader class="px-0 pt-0">
              <CardTitle class="text-base">Basic Information</CardTitle>
            </CardHeader>
            <dl class="space-y-3">
              <div class="space-y-1">
                <dt class="text-sm font-medium text-muted-foreground">Name</dt>
                <dd class="text-sm">{{ selectedIndexer?.name || '—' }}</dd>
              </div>
              <div class="space-y-1">
                <dt class="text-sm font-medium text-muted-foreground">Description</dt>
                <dd class="text-sm">{{ selectedIndexer?.description || '—' }}</dd>
              </div>
              <div class="space-y-1">
                <dt class="text-sm font-medium text-muted-foreground">Language</dt>
                <dd class="text-sm">{{ selectedIndexer?.language || '—' }}</dd>
              </div>
            </dl>
          </div>

          <!-- <div class="space-y-4">
            <CardHeader class="px-0 pt-0">
              <CardTitle class="text-base">Configuration</CardTitle>
            </CardHeader>
            <div
              v-if="
                indexerConfig && indexerConfig.fields && Object.keys(indexerConfig.fields).length > 0
              "
              class="space-y-3"
            >
              <div
                v-for="(field, key) in indexerConfig.fields"
                :key="key"
                class="space-y-1"
              >
                <dt class="text-sm font-medium text-muted-foreground">{{ key }}</dt>
                <dd class="text-sm">{{ formData.fields?.[key] || 'Not set' }}</dd>
              </div>
            </div>
            <p v-else class="text-sm text-muted-foreground">No additional configuration required.</p>
          </div> -->
        </div>
      </CardContent>
    </Card>
  </div>
</template>
