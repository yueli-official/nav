<script setup lang="ts">
import type { NavigationResponse } from "~/types/navigation";
definePageMeta({ width: "wide" });
const route = useRoute();
if (typeof route.query.category === "string" && route.query.category) {
  await navigateTo({ path: `/categories/${encodeURIComponent(route.query.category)}`, query: typeof route.query.group === "string" ? { group: route.query.group } : {} }, { redirectCode: 301, replace: true });
}
const { data, error, status, refresh } = await useFetch<NavigationResponse>("/api/navigation", { key: "navigation-catalog" });
const entries = computed(() => data.value ? flattenNavigation(data.value) : []);
const query = ref("");
const results = computed(() => searchNavigation(entries.value, query.value));
useSeoMeta({ title: () => data.value?.site.name, description: () => data.value?.site.description });
</script>

<template>
  <div class="nav-home min-w-0 space-y-7">
    <section class="nav-discovery mx-auto max-w-3xl space-y-4 py-5 sm:py-8" aria-labelledby="home-title">
      <p class="nav-eyebrow"><span aria-hidden="true" />为好奇心，留一个入口</p>
      <h1 id="home-title" class="nav-discovery-title font-display text-2xl font-semibold text-highlighted sm:text-3xl">发现值得收藏的站点</h1>
      <p class="text-sm text-muted">少一点漫无目的，多一点恰好有用。</p>
      <NavigationSearch v-model="query" :placeholder="data?.site.searchPlaceholder || '搜索名称、标签或域名'" :result-count="query.trim() ? results.length : undefined" />
    </section>
    <UAlert v-if="error" title="导航数据暂时不可用" description="请稍后重试。" color="error" variant="soft">
      <template #actions><UButton label="重试" color="neutral" variant="soft" @click="refresh()" /></template>
    </UAlert>
    <NavigationSkeleton v-else-if="status === 'pending'" />
    <section v-else-if="query.trim()" aria-labelledby="search-results-title" class="space-y-4">
      <h2 id="search-results-title" class="text-lg font-semibold text-highlighted">搜索结果</h2>
      <div v-if="results.length" class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
        <NavigationLinkCard v-for="entry in results" :key="entry.item.id" :entry="entry" />
      </div>
      <p v-else class="py-8 text-center text-sm text-muted">没有匹配的站点，试试其他名称或关键词。</p>
    </section>
    <template v-else-if="data">
      <NavigationFeatured :entries="entries" />
      <div v-if="!entries.length" class="rounded-xl border border-default p-6 text-sm text-muted">还没有收录站点。</div>
    </template>
    <section v-if="data && !error" id="categories" class="scroll-mt-32 space-y-4" aria-labelledby="categories-title">
      <div class="flex items-end justify-between gap-3"><h2 id="categories-title" class="font-display text-xl font-semibold text-highlighted">分类</h2><span class="text-xs text-muted">从兴趣出发</span></div>
    <div v-if="data.categories.length" class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <NuxtLink v-for="category in data.categories" :key="category.id" :to="`/categories/${encodeURIComponent(category.id)}`" data-spotlight class="nav-category-card nav-spotlight group flex min-w-0 flex-col gap-4 rounded-xl border border-default bg-default p-5 transition focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary">
        <div class="flex min-w-0 items-center gap-3">
          <span class="nav-category-icon grid size-10 shrink-0 place-items-center rounded-lg"><UIcon :name="category.icon" class="size-5" /></span>
          <h3 class="min-w-0 flex-1 font-display text-lg font-semibold text-highlighted group-hover:text-primary">{{ category.title }}</h3>
          <UIcon name="i-tabler-arrow-right" class="size-4 shrink-0 text-muted" />
        </div>
        <p class="text-sm leading-6 text-muted">{{ category.description }}</p>
        <div class="mt-auto flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted">
          <span>{{ category.groups.length }} 个主题</span>
          <span>{{ category.groups.reduce((sum, group) => sum + group.linkCount, 0) }} 个站点</span>
        </div>
      </NuxtLink>
    </div>
    <p v-else class="py-8 text-center text-sm text-muted">还没有分类。</p>
    </section>

  </div>
</template>
