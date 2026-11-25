<script setup lang="ts">
import type { NavigationResponse } from "~/types/navigation";

const route = useRoute();
const mainWidth = computed(
  () =>
    PAGE_WIDTHS[(route.meta.width as PageWidth) ?? "full"] ?? PAGE_WIDTHS.full,
);
const { data } = await useFetch<NavigationResponse>(
  "/api/navigation",
  {
    key: "navigation-catalog",
  },
);
const { openSearch } = useNavigationSearch();

function moveSpotlight(event: PointerEvent) {
  if (event.pointerType !== "mouse" || window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
  const card = (event.target as Element).closest<HTMLElement>("[data-spotlight]");
  if (!card) return;
  const bounds = card.getBoundingClientRect();
  card.style.setProperty("--pointer-x", `${event.clientX - bounds.left}px`);
  card.style.setProperty("--pointer-y", `${event.clientY - bounds.top}px`);
}

defineShortcuts({
  "/": { usingInput: false, handler: openSearch },
  meta_k: { usingInput: true, handler: openSearch },
  ctrl_k: { usingInput: true, handler: openSearch },
});
</script>

<template>
  <div class="nav-public yueli-app-shell flex min-h-dvh flex-col text-default" @pointermove="moveSpotlight">
    <AppHeader
      :width-class="mainWidth"
      :brand-name="data?.site.name"
      :brand-logo-url="data?.site.logoUrl"
    />
    <NavigationSearchDialog :placeholder="data?.site.searchPlaceholder" />
    <main
      id="public-main"
      tabindex="-1"
      class="mx-auto w-full flex-1 px-4 pb-16 pt-6 outline-none sm:px-6 sm:pt-8 lg:px-8"
      :class="mainWidth"
    >
      <slot />
    </main>
    <footer
      v-if="data?.site.footerTagline"
      class="border-t border-default bg-default/80 dark:bg-[var(--yueli-surface-region)]"
    >
      <div
        class="mx-auto flex w-full flex-col gap-2 px-4 py-6 text-sm text-muted sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8"
        :class="mainWidth"
      >
        <p>{{ data?.site.footerTagline }}</p>
      </div>
    </footer>
  </div>
</template>

<style src="~/assets/css/public.css"></style>
