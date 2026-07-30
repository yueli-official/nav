<script setup lang="ts">
const { rows = 6 } = defineProps<{ rows?: number }>();
</script>

<template>
  <ClientOnly>
    <slot />
    <template #fallback>
      <slot name="fallback">
        <SkeletonList :rows="rows" />
      </slot>
    </template>
  </ClientOnly>
</template>
