<template>
  <Page>
    <div v-if="isLoading">Loading...</div>
    <div v-else-if="isError">Error</div>
    <div v-else-if="data">
      <h1>{{ data.title }}</h1>
      <p>{{ data.releaseDate }}</p>
    </div>
  </Page>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { getV1SeriesByIdOptions } from '@/client/@tanstack/vue-query.gen'
import Page from '@/components/Page.vue'

const route = useRoute()

const id = computed(() => {
  const castAttept = Number(Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)
  if (isNaN(castAttept)) {
    throw new Error('Invalid series ID')
  }

  return castAttept
})

const { isLoading, isError, data } = useQuery({
  ...getV1SeriesByIdOptions({ path: { id: id.value } }),
})
</script>

<style scoped></style>
