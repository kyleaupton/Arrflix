<script setup lang="ts">
import { type ModelIndexerDefinition } from '@/client/types.gen'

defineProps<{
  selectedIndexer: ModelIndexerDefinition | null
  saveData: unknown
}>()
</script>

<template>
  <div class="review-step">
    <h3 class="text-lg font-semibold mb-4">Review & Create</h3>
    <p class="text-muted-color mb-6">Review your indexer configuration before creating.</p>

    <div class="bg-gray-50 rounded-lg p-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h4 class="font-semibold mb-3">Basic Information</h4>
          <dl class="space-y-2">
            <div>
              <dt class="text-sm font-medium text-muted-color">Name</dt>
              <dd class="text-sm">{{ selectedIndexer?.name }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-muted-color">Description</dt>
              <dd class="text-sm">{{ selectedIndexer?.description }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-muted-color">Language</dt>
              <dd class="text-sm">{{ selectedIndexer?.language }}</dd>
            </div>
          </dl>
        </div>

        <!-- <div>
          <h4 class="font-semibold mb-3">Configuration</h4>
          <div
            v-if="
              indexerConfig && indexerConfig.fields && Object.keys(indexerConfig.fields).length > 0
            "
            class="space-y-2"
          >
            <div
              v-for="(field, key) in indexerConfig.fields"
              :key="key"
              class="flex justify-between"
            >
              <span class="text-sm font-medium text-muted-color">{{ key }}:</span>
              <span class="text-sm">{{ formData.fields?.[key] || 'Not set' }}</span>
            </div>
          </div>
          <p v-else class="text-sm text-muted-color">No additional configuration required.</p>
        </div> -->
      </div>
    </div>
  </div>
</template>

<style scoped>
.review-step {
  /* Height is now managed by parent scroll container */
}
</style>
