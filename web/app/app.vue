<script setup lang="ts">
import { FeedbackToastRegion } from "@yueli/ui/feedback/pattern";

const { brand: siteBrand } = useSiteRuntime();
const { call } = useApi();
const { data: siteBranding } = await useAsyncData("nav-site-branding", () =>
  call<{ faviconUrl: string }>("/api/v1/site/branding"),
);

useHead(() => ({
  htmlAttrs: { lang: "zh-CN" },
  titleTemplate: (title) =>
    title ? `${title} · ${siteBrand.value}` : siteBrand.value,
  link: [{ rel: "icon", href: siteBranding.value?.faviconUrl || "/favicon.svg" }],
}));
</script>

<template>
  <UApp :toaster="null">
    <NuxtLoadingIndicator color="var(--ui-primary)" :height="3" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <FeedbackToastRegion />
  </UApp>
</template>
