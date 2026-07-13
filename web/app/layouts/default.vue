<script setup lang="ts">
const route = useRoute();
const mainWidth = computed(
  () =>
    PAGE_WIDTHS[(route.meta.width as PageWidth) ?? "full"] ?? PAGE_WIDTHS.full,
);
</script>

<template>
  <div class="platform-app-shell flex min-h-dvh flex-col text-default">
    <AppHeader :width-class="mainWidth" />
    <main
      id="public-main"
      tabindex="-1"
      class="mx-auto w-full flex-1 px-4 pb-16 pt-6 outline-none sm:px-6 sm:pt-8 lg:px-8"
      :class="mainWidth"
    >
      <slot />
    </main>
    <footer class="border-t border-default bg-default/80">
      <div
        class="mx-auto flex w-full flex-col gap-2 px-4 py-6 text-sm text-muted sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8"
        :class="mainWidth"
      >
        <p>月离导航，持续整理值得回访的互联网入口。</p>
        <p>内容以公开链接为准，请遵循目标站点的使用规则。</p>
      </div>
    </footer>
  </div>
</template>
