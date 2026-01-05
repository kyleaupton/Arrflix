<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Command {
  id: string
  title: string
  description: string
  icon: string
}

interface Props {
  commands: Command[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  select: [commandId: string]
  close: []
}>()

const selectedIndex = ref(0)

const selectCommand = (index: number, event?: MouseEvent) => {
  // Stop event propagation to prevent closing the dialog
  if (event) {
    event.stopPropagation()
    event.preventDefault()
  }
  const command = props.commands[index]
  if (command) {
    emit('select', command.id)
  }
}

function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    event.preventDefault()
    emit('close')
    return
  }

  if (event.key === 'ArrowUp') {
    event.preventDefault()
    selectedIndex.value = Math.max(0, selectedIndex.value - 1)
  }

  if (event.key === 'ArrowDown') {
    event.preventDefault()
    selectedIndex.value = Math.min(props.commands.length - 1, selectedIndex.value + 1)
  }

  if (event.key === 'Enter') {
    event.preventDefault()
    selectCommand(selectedIndex.value)
  }

  // Close on backspace if we're filtering (user deleting the /)
  if (event.key === 'Backspace') {
    emit('close')
  }
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.slash-commands')) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
  document.addEventListener('mousedown', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
  document.removeEventListener('mousedown', handleClickOutside)
})
</script>

<template>
  <div
    class="slash-commands w-64 rounded-md border bg-popover p-2 text-popover-foreground shadow-md"
    @click.stop
    @mousedown.stop
  >
    <div
      v-for="(command, index) in commands"
      :key="command.id"
      :class="[
        'flex items-center gap-3 rounded-sm px-2 py-2 cursor-pointer',
        index === selectedIndex ? 'bg-accent text-accent-foreground' : 'hover:bg-accent/50',
      ]"
      @click="selectCommand(index, $event)"
      @mousedown.stop
      @mouseenter="selectedIndex = index"
    >
      <span class="text-xl">{{ command.icon }}</span>
      <div class="flex flex-col gap-0.5">
        <div class="font-medium text-sm">{{ command.title }}</div>
        <div class="text-xs text-muted-foreground">{{ command.description }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slash-commands {
  animation: fadeIn 0.15s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
