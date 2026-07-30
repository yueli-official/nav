<script setup lang="ts">
const { title } = defineProps<{
  id: string;
  title: string;
  imageClass?: string;
}>();
</script>

<template>
  <span class="font-display text-sm font-bold uppercase" aria-hidden="true">{{
    title.slice(0, 1)
  }}</span>
</template>
