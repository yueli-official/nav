<script setup lang="ts">
import { portalSites, searchPortalEntries } from "~/data/portal";
import type { NavigationResponse } from "~/types/navigation";

definePageMeta({ width: "wide" });

const route = useRoute();
if (typeof route.query.category === "string" && route.query.category) {
  await navigateTo({
    path: `/categories/${encodeURIComponent(route.query.category)}`,
    query: typeof route.query.group === "string" ? { group: route.query.group } : {},
  }, { redirectCode: 301, replace: true });
}

const { data } = useNuxtData<NavigationResponse>("navigation-catalog");
const { brand } = useSiteRuntime();
const query = ref("");
const resultList = useTemplateRef<HTMLUListElement>("result-list");
const searching = computed(() => Boolean(query.value.trim()));
const results = computed(() => searchPortalEntries(query.value));

useHead(() => ({
  title: data.value?.site.name || brand.value,
  titleTemplate: "%s",
}));
useSeoMeta({ description: "了解月离站群各个网站的用途，直达工具、教程、文档和常用功能。" });
</script>

<template>
  <div class="mx-auto min-w-0 max-w-5xl space-y-7">
    <section class="nav-discovery mx-auto max-w-2xl space-y-3 py-4 sm:py-7" aria-labelledby="home-title">
      <h1 id="home-title" class="nav-discovery-title font-display text-2xl font-semibold text-highlighted sm:text-3xl">
        月离工具与资源
      </h1>
      <NavigationSearch
        v-model="query"
        :placeholder="data?.site.searchPlaceholder || '搜索网站、工具或教程'"
        :result-count="searching ? results.length : undefined"
        @submit="resultList?.querySelector('a')?.focus()"
      />
    </section>

    <section v-if="searching" id="categories" class="scroll-mt-24 space-y-3" aria-labelledby="results-title">
      <h2 id="results-title" class="font-display text-lg font-semibold text-highlighted">搜索结果</h2>
      <ul v-if="results.length" ref="result-list" class="grid gap-2 sm:grid-cols-2">
        <li v-for="entry in results" :key="entry.url" class="rounded-xl border border-default bg-default p-1">
          <PortalEntryLink :entry="entry" />
        </li>
      </ul>
      <p v-else role="status" class="rounded-xl border border-default bg-default px-4 py-6 text-sm text-muted">
        没有匹配的入口，试试网站名、软件名或更短的关键词。
      </p>
    </section>

    <section v-else id="categories" class="scroll-mt-24 space-y-4" aria-labelledby="sites-title">
      <h2 id="sites-title" class="font-display text-lg font-semibold text-highlighted">月离站群</h2>
      <div class="grid gap-4 md:grid-cols-2">
        <article
          v-for="site in portalSites"
          :key="site.id"
          :aria-labelledby="`portal-${site.id}`"
          class="flex min-w-0 flex-col rounded-2xl border border-default bg-default"
        >
          <div class="flex-1 p-4 sm:p-5">
            <a :href="site.url" class="group flex items-center gap-3 rounded-lg focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-primary">
              <img :src="site.logo" alt="" width="40" height="40" class="size-10 shrink-0 rounded-lg object-contain" />
              <div class="min-w-0 flex-1">
                <h3 :id="`portal-${site.id}`" class="font-display text-lg font-semibold text-highlighted group-hover:text-primary">{{ site.name }}</h3>
                <span class="block truncate text-xs text-muted">{{ site.hostname }}</span>
              </div>
              <UIcon name="i-tabler-arrow-up-right" class="size-4 shrink-0 text-muted group-hover:text-primary" aria-hidden="true" />
            </a>
            <p class="mt-3 text-sm leading-6 text-toned">{{ site.description }}</p>
          </div>
          <ul class="flex flex-wrap gap-2 border-t border-default px-4 py-3 sm:px-5" aria-label="快速入口">
            <li v-for="link in site.links" :key="link.path">
              <UButton
                :href="link.url"
                :label="link.title"
                size="sm"
                color="neutral"
                variant="soft"
              />
            </li>
          </ul>
        </article>
      </div>
    </section>
  </div>
</template>
