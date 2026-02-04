<template>
  <section v-if="hasProviders" class="space-y-4">
    <h2 class="text-xl font-semibold">Where to Watch</h2>

    <div class="space-y-4">
      <!-- Streaming providers -->
      <div v-if="providers?.flatrate?.length" class="space-y-2">
        <h3 class="text-sm font-medium text-muted-foreground">Stream</h3>
        <div class="flex flex-wrap gap-2">
          <ProviderLogo
            v-for="provider in providers.flatrate"
            :key="provider.providerId"
            :provider="provider"
          />
        </div>
      </div>

      <!-- Rent providers -->
      <div v-if="providers?.rent?.length" class="space-y-2">
        <h3 class="text-sm font-medium text-muted-foreground">Rent</h3>
        <div class="flex flex-wrap gap-2">
          <ProviderLogo
            v-for="provider in providers.rent"
            :key="provider.providerId"
            :provider="provider"
          />
        </div>
      </div>

      <!-- Buy providers -->
      <div v-if="providers?.buy?.length" class="space-y-2">
        <h3 class="text-sm font-medium text-muted-foreground">Buy</h3>
        <div class="flex flex-wrap gap-2">
          <ProviderLogo
            v-for="provider in providers.buy"
            :key="provider.providerId"
            :provider="provider"
          />
        </div>
      </div>
    </div>

    <!-- JustWatch attribution -->
    <p class="text-xs text-muted-foreground">
      Data provided by
      <a
        v-if="providers?.link"
        :href="providers.link"
        target="_blank"
        rel="noopener noreferrer"
        class="underline hover:text-foreground transition-colors"
      >
        JustWatch
      </a>
      <span v-else>JustWatch</span>
    </p>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ModelWatchProviders } from '@/client/types.gen'
import ProviderLogo from './ProviderLogo.vue'

const props = defineProps<{
  providers?: ModelWatchProviders
}>()

const hasProviders = computed(() => {
  if (!props.providers) return false
  return (
    (props.providers.flatrate?.length ?? 0) > 0 ||
    (props.providers.rent?.length ?? 0) > 0 ||
    (props.providers.buy?.length ?? 0) > 0
  )
})
</script>
