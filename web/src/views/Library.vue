<template>
  <Page title="Library" subtitle="Browse and manage your media" :content-class="'space-y-3'">
    <!-- <LibraryItem v-for="item in query.data || []" :key="item.id" :item="item" /> -->
  </Page>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { getV1LibraryOptions } from '@/client/@tanstack/vue-query.gen'

import Page from '@/components/Page.vue'
import LibraryItem from '@/components/library/LibraryItem.vue'

const query = useQuery(getV1LibraryOptions())

onMounted(async () => {
  const { data } = await query.refetch()
  console.log(data)
})
</script>

<style scoped></style>
