<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold">Library</h1>
        <p class="text-sm text-muted-foreground">Browse and manage your media</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            v-model="searchInput"
            placeholder="Search titles..."
            class="pl-9 w-64"
            @input="onSearchInput"
          />
        </div>
        <ViewToggle v-model="viewMode" />
      </div>
    </div>

    <!-- Type Filter Tabs -->
    <Tabs v-model="typeFilter" @update:model-value="onTypeFilterChange">
      <TabsList>
        <TabsTrigger value="">All</TabsTrigger>
        <TabsTrigger value="movie">Movies</TabsTrigger>
        <TabsTrigger value="series">Series</TabsTrigger>
      </TabsList>
    </Tabs>

    <!-- Loading State -->
    <div v-if="isLoading" class="space-y-6">
      <div v-if="viewMode === 'grid'" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
        <Skeleton v-for="i in 12" :key="i" class="aspect-[2/3] rounded-lg" />
      </div>
      <div v-else class="space-y-2">
        <Skeleton class="h-10 w-full" />
        <Skeleton v-for="i in 5" :key="i" class="h-16 w-full" />
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="isError" class="flex flex-col items-center justify-center py-12 text-center">
      <p class="text-destructive">Failed to load library</p>
      <p class="text-sm text-muted-foreground mt-2">{{ error?.message || 'Please try again later' }}</p>
      <Button variant="outline" class="mt-4" @click="refetch">
        Try Again
      </Button>
    </div>

    <!-- Empty State -->
    <div v-else-if="!data?.data || data.data.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
      <Film class="h-12 w-12 text-muted-foreground mb-4" />
      <p class="text-lg font-medium">No media found</p>
      <p class="text-sm text-muted-foreground mt-1">
        {{ searchQuery ? 'Try adjusting your search or filters' : 'Your library is empty. Add some media to get started!' }}
      </p>
    </div>

    <!-- Content -->
    <template v-else>
      <!-- Grid View with Flexbox -->
      <div
        v-if="viewMode === 'grid'"
        class="flex flex-wrap gap-3 max-w-[1600px]"
      >
        <Poster
          v-for="item in data.data"
          :key="item.id"
          :item="item"
          :to="getItemRoute(item)"
          size="medium"
        />
      </div>

      <!-- Table View -->
      <div v-else class="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead class="w-16">Poster</TableHead>
              <TableHead>Title</TableHead>
              <TableHead class="w-24">Type</TableHead>
              <TableHead class="w-20">Year</TableHead>
              <TableHead class="w-40">Added</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="item in data.data"
              :key="item.id"
              class="cursor-pointer hover:bg-muted/50"
              @click="navigateToItem(item)"
            >
              <TableCell>
                <img
                  v-if="item.posterPath"
                  :src="`https://image.tmdb.org/t/p/w92${item.posterPath}`"
                  :alt="item.title"
                  class="w-10 h-15 object-cover rounded"
                />
                <div v-else class="w-10 h-15 bg-muted rounded flex items-center justify-center">
                  <Film class="h-4 w-4 text-muted-foreground" />
                </div>
              </TableCell>
              <TableCell class="font-medium">{{ item.title }}</TableCell>
              <TableCell>
                <Badge :variant="item.type === 'movie' ? 'default' : 'secondary'">
                  {{ item.type === 'movie' ? 'Movie' : 'Series' }}
                </Badge>
              </TableCell>
              <TableCell>{{ item.year || '-' }}</TableCell>
              <TableCell class="text-muted-foreground">
                {{ formatDate(item.createdAt) }}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>

      <!-- Pagination -->
      <LibraryPagination
        v-if="data.pagination.totalPages > 1"
        :pagination="data.pagination"
        @update:page="onPageChange"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useDebounceFn, useLocalStorage } from '@vueuse/core'
import { Search, Film } from 'lucide-vue-next'
import { getV1LibraryOptions, getV1LibraryQueryKey } from '@/client/@tanstack/vue-query.gen'
import type { ModelLibraryItem } from '@/client/types.gen'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import Poster from '@/components/poster/Poster.vue'
import ViewToggle from '@/components/library/ViewToggle.vue'
import LibraryPagination from '@/components/library/LibraryPagination.vue'

const router = useRouter()

// State
const page = ref(1)
const pageSize = ref(20)
const typeFilter = ref('')
const searchQuery = ref('')
const searchInput = ref('')
const viewMode = useLocalStorage<'grid' | 'table'>('library-view-mode', 'grid')

// Debounced search
const debouncedSearch = useDebounceFn((value: string) => {
  searchQuery.value = value
  page.value = 1 // Reset to first page on search
}, 300)

const onSearchInput = () => {
  debouncedSearch(searchInput.value)
}

const onTypeFilterChange = () => {
  page.value = 1 // Reset to first page on filter change
}

const onPageChange = (newPage: number) => {
  page.value = newPage
}

// Query
const queryParams = computed(() => ({
  query: {
    page: page.value,
    pageSize: pageSize.value,
    type: typeFilter.value || undefined,
    search: searchQuery.value || undefined,
  },
}))

const { isLoading, isError, error, data, refetch } = useQuery({
  ...getV1LibraryOptions(queryParams.value),
  queryKey: computed(() => getV1LibraryQueryKey(queryParams.value)),
})

// Watch for query param changes to refetch
watch([page, pageSize, typeFilter, searchQuery], () => {
  // Query will auto-refetch due to reactive queryKey
})

// Navigation
const getItemRoute = (item: ModelLibraryItem) => {
  return { path: `/${item.type}/${item.tmdbId}` }
}

const navigateToItem = (item: ModelLibraryItem) => {
  router.push(getItemRoute(item))
}

// Formatting
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
