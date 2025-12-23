<template>
  <div class="settings-layout">
    <router-view />
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.settings-layout {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
