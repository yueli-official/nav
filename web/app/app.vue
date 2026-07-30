<script setup lang="ts">
const { brand: siteBrand } = useSiteRuntime();

useHead({
  titleTemplate: (title) =>
    title ? `${title} · ${siteBrand.value}` : siteBrand.value,
});
</script>

<template>
  <UApp :toaster="{ duration: 4500 }">
    <NuxtLoadingIndicator color="var(--ui-primary)" :height="3" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>
