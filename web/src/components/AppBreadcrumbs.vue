<script setup lang="ts">
import { RouterLink } from 'vue-router'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Separator } from '@/components/ui/separator'
import { useBreadcrumbs } from '@/composables/useBreadcrumbs'

const { shouldShowBreadcrumbs, breadcrumbItems } = useBreadcrumbs()
</script>

<template>
  <template v-if="shouldShowBreadcrumbs">
    <Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
    <Breadcrumb>
      <BreadcrumbList>
        <template v-for="item in breadcrumbItems" :key="item.path">
          <BreadcrumbItem v-if="!item.isLast" class="hidden md:block">
            <BreadcrumbLink :as="RouterLink" :to="item.path">
              {{ item.label }}
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator v-if="!item.isLast" class="hidden md:block" />
          <BreadcrumbItem v-if="item.isLast">
            <BreadcrumbPage>{{ item.label }}</BreadcrumbPage>
          </BreadcrumbItem>
        </template>
      </BreadcrumbList>
    </Breadcrumb>
  </template>
</template>
