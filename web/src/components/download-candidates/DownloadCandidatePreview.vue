<template>
  <div class="preview-container">
    <div v-if="isLoading" class="flex items-center justify-center p-8">
      <i class="pi pi-spin pi-spinner text-2xl"></i>
      <span class="ml-2">Loading preview...</span>
    </div>

    <div v-else-if="error" class="p-4 bg-surface border border-red-500/30 rounded">
      <div class="text-red-400 font-semibold mb-1">Error loading preview</div>
      <div class="text-red-300 text-sm">{{ error }}</div>
    </div>

    <div v-else-if="trace" class="preview-content">
      <!-- Candidate Info -->
      <Card class="mb-4">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-file"></i>
            <span>{{ candidate.title }}</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="font-semibold">Size:</span>
              <span class="ml-2">{{ formatSize(candidate.size) }}</span>
            </div>
            <div>
              <span class="font-semibold">Seeders:</span>
              <span class="ml-2">{{ candidate.seeders }}</span>
            </div>
            <div>
              <span class="font-semibold">Indexer:</span>
              <span class="ml-2">{{ candidate.indexer }}</span>
            </div>
            <div>
              <span class="font-semibold">Categories:</span>
              <span class="ml-2">{{ candidate.categories.join(', ') || 'None' }}</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Final Plan -->
      <Card class="final-plan-card mb-4">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-check-circle text-primary"></i>
            <span>Final Download Plan</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <i class="pi pi-download text-primary"></i>
              <div>
                <div class="text-sm text-muted-color">Downloader</div>
                <div class="font-semibold">
                  <DownloaderReference
                    v-if="trace.finalPlan.downloaderId"
                    :downloader-id="trace.finalPlan.downloaderId"
                  />
                  <span v-else class="text-muted-color">Not set</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <i class="pi pi-folder text-primary"></i>
              <div>
                <div class="text-sm text-muted-color">Library</div>
                <div class="font-semibold">
                  <LibraryReference
                    v-if="trace.finalPlan.libraryId"
                    :library-id="trace.finalPlan.libraryId"
                  />
                  <span v-else class="text-muted-color">Not set</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <i class="pi pi-file-edit text-primary"></i>
              <div>
                <div class="text-sm text-muted-color">Name Template</div>
                <div class="font-semibold">
                  <NameTemplateReference
                    v-if="trace.finalPlan.nameTemplateId"
                    :name-template-id="trace.finalPlan.nameTemplateId"
                  />
                  <span v-else class="text-muted-color">Not set</span>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <Divider />

      <!-- Policy Evaluation -->
      <div>
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <i class="pi pi-list-check"></i>
          Policy Evaluation
        </h3>

        <div class="space-y-3">
          <Card
            v-for="policy in trace.policies"
            :key="policy.policyId"
            :class="{
              'policy-matched': policy.matched,
              'policy-unmatched': !policy.matched,
            }"
          >
            <template #content>
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-2">
                    <i
                      :class="
                        policy.matched
                          ? 'pi pi-check-circle text-green-400'
                          : 'pi pi-times-circle text-muted-color'
                      "
                      class="text-lg"
                    />
                    <span class="font-semibold">{{ policy.policyName }}</span>
                    <Badge
                      :value="`Priority: ${policy.priority}`"
                      severity="secondary"
                      size="small"
                    />
                    <Badge v-if="policy.matched" value="Matched" severity="success" size="small" />
                  </div>

                  <div v-if="policy.ruleEvaluated" class="text-sm text-muted-color mb-2">
                    <span class="font-medium">Rule:</span>
                    <code class="ml-2 px-2 py-1 bg-surface rounded">
                      {{ policy.ruleEvaluated.leftOperand }}
                      {{ policy.ruleEvaluated.operator }}
                      {{ policy.ruleEvaluated.rightOperand }}
                    </code>
                  </div>

                  <div v-if="policy.matched && policy.actionsApplied.length > 0" class="mt-3">
                    <div class="text-sm font-medium mb-2">Actions Applied:</div>
                    <ul class="list-disc list-inside space-y-1 text-sm text-content-color">
                      <li v-for="action in policy.actionsApplied" :key="action.order">
                        <span class="font-medium">{{ formatActionType(action.type) }}:</span>
                        <span class="ml-1">{{ action.value }}</span>
                      </li>
                    </ul>
                  </div>

                  <div v-if="policy.stoppedProcessing" class="mt-2">
                    <Badge value="Processing Stopped" severity="warning" />
                  </div>
                </div>
              </div>
            </template>
          </Card>

          <div v-if="trace.policies.length === 0" class="text-center py-8 text-muted-color">
            No policies evaluated
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, watch, ref } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import Card from 'primevue/card'
import Badge from 'primevue/badge'
import Divider from 'primevue/divider'
import { postV1MovieByIdCandidatePreviewMutation } from '@/client/@tanstack/vue-query.gen'
import { type ModelDownloadCandidate, type ModelEvaluationTrace } from '@/client/types.gen'
import LibraryReference from '@/components/references/LibraryReference.vue'
import DownloaderReference from '@/components/references/DownloaderReference.vue'
import NameTemplateReference from '@/components/references/NameTemplateReference.vue'

const props = defineProps<{
  candidate: ModelDownloadCandidate
}>()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const dialogRef = inject('dialogRef') as any

const movieId = computed(() => {
  const id = dialogRef.value?.data?.movieId
  if (!id) {
    throw new Error('Movie ID is required')
  }
  return id
})

// Local state for trace
const trace = ref<ModelEvaluationTrace | undefined>(undefined)
const error = ref<string | undefined>(undefined)

// Preview mutation
const previewMutation = useMutation({
  ...postV1MovieByIdCandidatePreviewMutation(),
  onSuccess: (data) => {
    trace.value = data
    error.value = undefined
  },
  onError: (err) => {
    error.value = err?.message || 'Failed to load preview'
    trace.value = undefined
  },
})

const isLoading = computed(() => previewMutation.isPending.value)

// Trigger preview when candidate changes
watch(
  () => props.candidate,
  () => {
    if (props.candidate) {
      trace.value = undefined
      error.value = undefined
      previewMutation.mutate({
        path: { id: movieId.value },
        body: {
          indexerId: props.candidate.indexerId,
          guid: props.candidate.guid,
        },
      })
    }
  },
  { immediate: true },
)

// Helper functions
const formatSize = (bytes: number): string => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(2)} ${units[unitIndex]}`
}

const formatActionType = (type: string): string => {
  const actionMap: Record<string, string> = {
    set_downloader: 'Set Downloader',
    set_library: 'Set Library',
    set_name_template: 'Set Name Template',
    stop_processing: 'Stop Processing',
  }
  return actionMap[type] || type
}
</script>

<style scoped>
.preview-container {
  width: 100%;
  max-height: 70vh;
  overflow-y: auto;
}

.preview-content {
  padding: 0.5rem;
}

.policy-matched {
  background-color: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.policy-unmatched {
  background-color: var(--p-surface-ground);
  border: 1px solid var(--p-surface-border);
  opacity: 0.6;
}

.final-plan-card {
  border: 1px solid color-mix(in srgb, var(--p-primary-color) 30%, transparent);
}

code {
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
}
</style>

<style>
.final-plan-card .p-card-body {
  background: linear-gradient(
    135deg,
    color-mix(in srgb, var(--p-primary-color) 15%, var(--p-surface-ground)),
    color-mix(in srgb, var(--p-primary-color) 8%, var(--p-surface-ground))
  );
}
</style>
