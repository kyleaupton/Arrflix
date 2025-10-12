<template>
  <Page>
    <div v-if="isLoading">Loading...</div>
    <div v-else-if="isError">Error</div>
    <div v-else-if="data">
      <Rail v-for="rail in data" :key="rail.id" :rail="rail" />
    </div>
  </Page>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { getV1HomeOptions } from '@/client/@tanstack/vue-query.gen'
import Page from '@/components/Page.vue'
import Rail from '@/components/rails/Rail.vue'

const { isLoading, isError, data } = useQuery(getV1HomeOptions())
</script>

<style scoped></style>
