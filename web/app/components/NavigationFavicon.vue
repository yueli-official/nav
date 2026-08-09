<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    id: string;
    title: string;
    revision?: string;
    eager?: boolean;
    imageClass?: string;
  }>(),
  {
    imageClass: "size-6",
  },
);
const failed = ref(false);
const loaded = ref(false);
const fallback = computed(
  () => Array.from(props.title.trim())[0]?.toLocaleUpperCase() || "?",
);
const faviconSrc = computed(() =>
  props.revision
    ? `/api/favicon/${encodeURIComponent(props.id)}?v=${encodeURIComponent(props.revision)}`
    : "",
);

watch(
  () => [props.id, props.revision],
  () => {
    failed.value = false;
    loaded.value = false;
  },
);
</script>

<template>
  <span class="relative grid size-full place-items-center" aria-hidden="true">
    <span
      v-if="!loaded"
      data-navigation-favicon-fallback
      class="font-display text-sm font-bold uppercase"
      >{{ fallback }}</span
    >
    <img
      v-if="faviconSrc && !failed"
      data-navigation-favicon-image
      :src="faviconSrc"
      alt=""
      :class="['absolute object-contain', props.imageClass, !loaded && 'opacity-0']"
      :loading="props.eager ? 'eager' : 'lazy'"
      decoding="async"
      draggable="false"
      @load="loaded = true"
      @error="failed = true"
    />
  </span>
</template>
