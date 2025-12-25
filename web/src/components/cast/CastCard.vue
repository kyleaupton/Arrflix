<template>
  <router-link
    :to="{ path: `/person/${member.tmdbId}` }"
    class="cast-card flex flex-col items-center gap-2 w-32 sm:w-36"
  >
    <div
      class="cast-photo-wrapper relative w-full aspect-[2/3] rounded-lg overflow-hidden bg-muted"
    >
      <img
        v-if="profileImageUrl"
        :src="profileImageUrl"
        :alt="member.name"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
        <User class="size-12" />
      </div>
    </div>
    <div class="text-center w-full">
      <p class="font-medium text-sm truncate">{{ member.name }}</p>
      <p class="text-xs text-muted-foreground truncate">{{ member.character }}</p>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { User } from 'lucide-vue-next'
import { type ModelCastMember } from '@/client/types.gen'

const props = defineProps<{
  member: ModelCastMember
}>()

const profileImageUrl = computed(() => {
  if (!props.member.profilePath) return undefined
  return `https://image.tmdb.org/t/p/w185/${props.member.profilePath}`
})
</script>
