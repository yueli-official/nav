<script setup lang="ts">
const {
  id,
  title,
  imageClass = "size-6",
} = defineProps<{
  id: string;
  title: string;
  imageClass?: string;
}>();
const failed = ref(false);
</script>

<template>
  <img
    v-if="!failed"
    :src="`/api/favicon/${encodeURIComponent(id)}`"
    :alt="`${title} 图标`"
    :class="['object-contain', imageClass]"
    loading="lazy"
    @error="failed = true"
  />
  <span
    v-else
    class="font-display text-sm font-bold uppercase"
    aria-hidden="true"
    >{{ title.slice(0, 1) }}</span
  >
</template>
