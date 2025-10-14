<template>
  <section class="media-hero relative overflow-hidden rounded-lg">
    <div class="backdrop" :class="{ 'has-image': !!backdropUrl }">
      <img v-if="backdropUrl" :src="backdropUrl" alt="" aria-hidden="true" />
      <div class="backdrop-overlay" />
    </div>

    <div class="content relative p-4 sm:p-6 md:p-8">
      <div class="flex gap-4 md:gap-6 items-start">
        <div v-if="posterUrl" class="poster shadow-lg">
          <img :src="posterUrl" :alt="title" loading="eager" decoding="async" />
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-start justify-between gap-3">
            <h1 class="title text-2xl sm:text-3xl md:text-4xl font-semibold truncate">
              {{ title }}
            </h1>
            <div class="actions shrink-0">
              <slot name="actions" />
            </div>
          </div>
          <p v-if="subtitle" class="subtitle text-sm opacity-80 mt-1">{{ subtitle }}</p>

          <div v-if="chips && chips.length" class="chips mt-3 flex flex-wrap gap-2">
            <span v-for="(chip, i) in chips" :key="i" class="chip">{{ chip }}</span>
          </div>

          <p
            v-if="overview"
            class="overview mt-4 max-w-prose text-sm md:text-base leading-relaxed opacity-90"
          >
            {{ overview }}
          </p>

          <slot />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  overview?: string
  posterUrl?: string
  backdropUrl?: string
  chips?: string[]
}>()
</script>

<style scoped>
.media-hero {
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.35), rgba(0, 0, 0, 0.35));
}

.backdrop {
  position: absolute;
  inset: 0;
}

.backdrop img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: blur(6px);
  transform: scale(1.03);
}

.backdrop-overlay {
  position: absolute;
  inset: 0;
  /* Dark-ish, readable, preserves image detail */
  background: linear-gradient(
    180deg,
    rgba(0, 0, 0, 0.5),
    rgba(0, 0, 0, 0.25) 45%,
    rgba(0, 0, 0, 0.6)
  );
}

.poster {
  flex: 0 0 auto;
  width: 8rem; /* 128px */
  aspect-ratio: 2 / 3;
  border-radius: 12px;
  overflow: hidden;
  background: #0f172a;
}

@media (min-width: 768px) {
  .poster {
    width: 10rem;
  }
}

.poster img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.title {
  color: var(--p-content-color);
}
.subtitle {
  color: var(--p-content-color);
}
.overview {
  color: var(--p-content-color);
}

.chip {
  display: inline-block;
  padding: 2px 8px;
  font-size: 12px;
  border-radius: 9999px;
  background: rgba(0, 0, 0, 0.35);
  color: var(--p-content-color);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
</style>
