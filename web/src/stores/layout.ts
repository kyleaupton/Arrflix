import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useLayoutStore = defineStore('layout', () => {
  const screenWidth = ref(window.innerWidth)

  return {
    screenWidth,
  }
})
