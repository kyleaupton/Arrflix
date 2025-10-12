<template>
  <header class="app-header">
    <Button class="menu-btn" icon="pi pi-bars" variant="text" @click="emit('toggle-sidebar')" />

    <div class="search-wrap">
      <InputText
        v-model="query"
        placeholder="Search for movies, TV shows, etc."
        class="search-input"
        @keyup.enter="emit('search', query)"
      />
    </div>

    <div class="avatar-wrap">
      <Avatar label="S" shape="circle" />
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import InputText from 'primevue/inputtext'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'

const query = ref('')

const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
  (e: 'search', query: string): void
}>()
</script>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--p-surface-0);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.menu-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border: none;
  background: transparent;
  font-size: 1.25rem;
  cursor: pointer;
}

@media (min-width: 1024px) {
  .menu-btn {
    display: none;
  }
}

.search-wrap {
  flex: 1 1 auto;
}

.search-input {
  width: 100%;
}

.avatar-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
