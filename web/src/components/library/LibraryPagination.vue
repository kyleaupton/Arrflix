<template>
  <div class="flex items-center justify-between py-4">
    <div class="text-sm text-muted-foreground">
      Showing {{ startItem }}-{{ endItem }} of {{ pagination.total }} items
    </div>
    <Pagination
      v-slot="{ page }"
      :default-page="pagination.page"
      :total="pagination.total"
      :items-per-page="pagination.pageSize"
      :sibling-count="1"
      show-edges
      @update:page="$emit('update:page', $event)"
    >
      <PaginationContent class="flex items-center gap-1">
        <PaginationFirst />
        <PaginationPrevious />

        <template v-for="(item, index) in page.items" :key="index">
          <PaginationItem v-if="item.type === 'page'" :value="item.value" as-child>
            <Button
              class="w-9 h-9 p-0"
              :variant="item.value === pagination.page ? 'default' : 'outline'"
            >
              {{ item.value }}
            </Button>
          </PaginationItem>
          <PaginationEllipsis v-else :index="index" />
        </template>

        <PaginationNext />
        <PaginationLast />
      </PaginationContent>
    </Pagination>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationFirst,
  PaginationItem,
  PaginationLast,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import type { ModelPagination } from '@/client/types.gen'

const props = defineProps<{
  pagination: ModelPagination
}>()

defineEmits<{
  (e: 'update:page', page: number): void
}>()

const startItem = computed(() => {
  return (props.pagination.page - 1) * props.pagination.pageSize + 1
})

const endItem = computed(() => {
  return Math.min(
    props.pagination.page * props.pagination.pageSize,
    props.pagination.total
  )
})
</script>

