<script setup lang="ts">
import type {
  NavigationItem,
  NavigationResponse,
  NavigationResult,
} from "~/types/navigation";

definePageMeta({ width: "full" });

const ALL_GROUPS = "__all__";
const SEARCH_PAGE_SIZE = 24;
const route = useRoute();
const router = useRouter();

function queryValue(value: unknown) {
  if (Array.isArray(value)) return typeof value[0] === "string" ? value[0] : "";
  return typeof value === "string" ? value : "";
}

const searchQuery = ref(queryValue(route.query.q));
const selectedCategoryId = ref(queryValue(route.query.category));
const selectedGroupId = ref(queryValue(route.query.group) || ALL_GROUPS);
const visibleSearchCount = ref(SEARCH_PAGE_SIZE);

const { data, error, status, refresh } = await useFetch<NavigationResponse>(
  "/api/navigation",
  {
    key: "navigation-catalog",
  },
);

const categories = computed(() => data.value?.categories ?? []);
const allEntries = computed(() =>
  data.value ? flattenNavigation(data.value) : [],
);
const entryMap = computed(
  () => new Map(allEntries.value.map((entry) => [entry.item.id, entry])),
);
const selectedCategory = computed(
  () =>
    categories.value.find(
      (category) => category.id === selectedCategoryId.value,
    ) ?? categories.value[0],
);
const featuredEntries = computed(() =>
  allEntries.value.filter((entry) => entry.item.featured).slice(0, 5),
);
const searchResults = computed(() =>
  searchNavigation(allEntries.value, searchQuery.value),
);
const visibleSearchResults = computed(() =>
  searchResults.value.slice(0, visibleSearchCount.value),
);
const visibleGroups = computed(() => {
  const category = selectedCategory.value;
  if (!category) return [];
  if (selectedGroupId.value === ALL_GROUPS) return category.groups;
  return category.groups.filter((group) => group.id === selectedGroupId.value);
});

watch(
  categories,
  (nextCategories) => {
    if (!nextCategories.length) return;
    if (
      !nextCategories.some(
        (category) => category.id === selectedCategoryId.value,
      )
    ) {
      selectedCategoryId.value = nextCategories[0]!.id;
    }
  },
  { immediate: true },
);

watch(selectedCategoryId, () => {
  const groupExists = selectedCategory.value?.groups.some(
    (group) => group.id === selectedGroupId.value,
  );
  if (!groupExists) selectedGroupId.value = ALL_GROUPS;
});

watch(searchQuery, () => {
  visibleSearchCount.value = SEARCH_PAGE_SIZE;
});

watch(
  () => route.fullPath,
  () => {
    const nextQuery = queryValue(route.query.q);
    const nextCategory = queryValue(route.query.category);
    const nextGroup = queryValue(route.query.group) || ALL_GROUPS;
    const resolvedCategory = nextCategory || categories.value[0]?.id || "";
    if (searchQuery.value !== nextQuery) searchQuery.value = nextQuery;
    if (selectedCategoryId.value !== resolvedCategory)
      selectedCategoryId.value = resolvedCategory;
    if (selectedGroupId.value !== nextGroup) selectedGroupId.value = nextGroup;
  },
);

watch([searchQuery, selectedCategoryId, selectedGroupId], async () => {
  if (!import.meta.client) return;

  const nextQuery: Record<string, string> = {};
  if (searchQuery.value.trim()) nextQuery.q = searchQuery.value.trim();
  if (
    selectedCategoryId.value &&
    selectedCategoryId.value !== categories.value[0]?.id
  ) {
    nextQuery.category = selectedCategoryId.value;
  }
  if (selectedGroupId.value !== ALL_GROUPS)
    nextQuery.group = selectedGroupId.value;

  const currentQuery = {
    q: queryValue(route.query.q),
    category: queryValue(route.query.category),
    group: queryValue(route.query.group),
  };
  if (
    currentQuery.q === (nextQuery.q ?? "") &&
    currentQuery.category === (nextQuery.category ?? "") &&
    currentQuery.group === (nextQuery.group ?? "")
  )
    return;

  await router.replace({ query: nextQuery });
});

function entryFor(item: NavigationItem): NavigationResult {
  const entry = entryMap.value.get(item.id);
  if (!entry) throw new Error(`Missing navigation entry: ${item.id}`);
  return entry;
}

function loadMoreSearchResults() {
  visibleSearchCount.value += SEARCH_PAGE_SIZE;
}

async function retryLoad() {
  await refresh();
}

function clearFilters() {
  searchQuery.value = "";
  selectedGroupId.value = ALL_GROUPS;
}

useSeoMeta({
  title: () => data.value?.site.name ?? "月离导航",
  description: () => data.value?.site.description ?? "精选互联网工作台",
  ogTitle: () => data.value?.site.title ?? "月离导航",
  ogDescription: () => data.value?.site.description ?? "精选互联网工作台",
});
</script>

<template>
  <div class="min-w-0 space-y-12 overflow-x-hidden">
    <section
      class="nav-hero-surface overflow-hidden rounded-2xl border border-default px-5 py-8 sm:px-8 sm:py-10 lg:px-12 lg:py-12"
    >
      <div
        class="grid min-w-0 gap-8 lg:grid-cols-[minmax(0,1fr)_minmax(320px,0.72fr)] lg:items-end"
      >
        <div class="min-w-0">
          <p
            class="mb-3 inline-flex items-center gap-2 text-sm font-medium text-primary"
          >
            <UIcon name="i-tabler-compass" class="size-4" />
            面向创作与开发的精选索引
          </p>
          <h1
            class="font-display max-w-3xl text-balance text-4xl font-bold leading-[1.08] tracking-tight text-highlighted sm:text-5xl lg:text-6xl"
          >
            {{ data?.site.title ?? "把常用互联网，整理成工作台" }}
          </h1>
          <p
            class="mt-4 max-w-[64ch] text-base leading-7 text-toned sm:text-lg"
          >
            {{
              data?.site.description ??
              "按任务浏览，也可以直接搜索名称、标签和域名。"
            }}
          </p>
        </div>

        <dl
          v-if="data"
          class="grid grid-cols-3 divide-x divide-default rounded-xl border border-default bg-default/82 py-4 backdrop-blur"
        >
          <div class="px-3 text-center sm:px-5">
            <dt class="text-xs text-muted">分类</dt>
            <dd
              class="font-display mt-1 text-2xl font-semibold text-highlighted"
            >
              {{ data.stats.categoryCount }}
            </dd>
          </div>
          <div class="px-3 text-center sm:px-5">
            <dt class="text-xs text-muted">主题</dt>
            <dd
              class="font-display mt-1 text-2xl font-semibold text-highlighted"
            >
              {{ data.stats.groupCount }}
            </dd>
          </div>
          <div class="px-3 text-center sm:px-5">
            <dt class="text-xs text-muted">站点</dt>
            <dd
              class="font-display mt-1 text-2xl font-semibold text-highlighted"
            >
              {{ data.stats.linkCount }}
            </dd>
          </div>
        </dl>
      </div>

      <div class="mt-8 max-w-4xl">
        <NavigationSearch
          v-model="searchQuery"
          :placeholder="
            data?.site.searchPlaceholder ?? '搜索工具、文档、社区或关键词'
          "
          :result-count="searchQuery ? searchResults.length : undefined"
        />
      </div>
    </section>

    <section
      v-if="featuredEntries.length && !searchQuery"
      id="featured"
      class="scroll-mt-24 space-y-5"
      aria-labelledby="featured-title"
    >
      <div>
        <h2
          id="featured-title"
          class="font-display text-2xl font-semibold tracking-tight text-highlighted"
        >
          常用精选
        </h2>
        <p class="mt-1 text-sm leading-6 text-muted">
          跨分类保留少量高频入口，减少重复查找。
        </p>
      </div>
      <div class="grid gap-3 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,1.95fr)]">
        <NavigationLinkCard
          v-if="featuredEntries[0]"
          :entry="featuredEntries[0]"
        />
        <div class="grid gap-3 sm:grid-cols-2">
          <NavigationLinkCard
            v-for="entry in featuredEntries.slice(1)"
            :key="entry.item.id"
            :entry="entry"
          />
        </div>
      </div>
    </section>

    <section id="catalog" class="scroll-mt-24 space-y-8">
      <div
        v-if="error"
        class="rounded-2xl border border-error/35 bg-error/5 p-6"
        role="alert"
      >
        <div class="flex items-start gap-4">
          <span
            class="grid size-10 shrink-0 place-items-center rounded-lg bg-error/10 text-error"
          >
            <UIcon name="i-tabler-alert-triangle" class="size-5" />
          </span>
          <div class="min-w-0">
            <h2 class="font-display text-lg font-semibold text-highlighted">
              导航数据加载失败
            </h2>
            <p class="mt-1 text-sm leading-6 text-muted">
              请稍后重试。如果问题持续存在，需要检查导航数据格式。
            </p>
            <UButton
              class="mt-4"
              color="neutral"
              variant="soft"
              icon="i-tabler-refresh"
              label="重新加载"
              @click="retryLoad"
            />
          </div>
        </div>
      </div>

      <NavigationSkeleton v-else-if="status === 'pending'" />

      <template v-else-if="data">
        <div v-if="searchQuery" class="space-y-6">
          <div
            class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
          >
            <div>
              <h2
                class="font-display text-2xl font-semibold tracking-tight text-highlighted"
              >
                搜索结果
              </h2>
              <p class="mt-1 text-sm text-muted">
                同时检索名称、描述、标签、域名和所属分类。
              </p>
            </div>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-tabler-x"
              label="清空搜索"
              @click="clearFilters"
            />
          </div>

          <div
            v-if="visibleSearchResults.length"
            class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"
          >
            <NavigationLinkCard
              v-for="entry in visibleSearchResults"
              :key="entry.item.id"
              :entry="entry"
              show-context
            />
          </div>

          <div
            v-else
            class="rounded-2xl border border-dashed border-default bg-elevated/35 px-6 py-12 text-center"
          >
            <span
              class="mx-auto grid size-12 place-items-center rounded-xl bg-default text-muted ring-1 ring-default"
            >
              <UIcon name="i-tabler-search-off" class="size-6" />
            </span>
            <h2
              class="font-display mt-4 text-xl font-semibold text-highlighted"
            >
              没有找到匹配站点
            </h2>
            <p class="mx-auto mt-2 max-w-md text-sm leading-6 text-muted">
              试试更短的关键词、英文名称，或者清空搜索后按分类浏览。
            </p>
            <UButton
              class="mt-5"
              color="neutral"
              variant="soft"
              label="返回分类"
              @click="clearFilters"
            />
          </div>

          <div
            v-if="visibleSearchCount < searchResults.length"
            class="flex justify-center"
          >
            <UButton
              color="neutral"
              variant="soft"
              icon="i-tabler-arrow-down"
              label="加载更多"
              @click="loadMoreSearchResults"
            />
          </div>
        </div>

        <div v-else class="space-y-8">
          <CategoryRail v-model="selectedCategoryId" :categories="categories" />

          <template v-if="selectedCategory">
            <div
              class="grid gap-5 border-b border-default pb-7 lg:grid-cols-[minmax(0,1fr)_minmax(320px,0.8fr)] lg:items-end"
            >
              <div class="min-w-0">
                <div class="flex items-center gap-3">
                  <span
                    class="grid size-11 shrink-0 place-items-center rounded-xl bg-primary/10 text-primary"
                  >
                    <UIcon :name="selectedCategory.icon" class="size-6" />
                  </span>
                  <h2
                    class="font-display text-3xl font-semibold tracking-tight text-highlighted"
                  >
                    {{ selectedCategory.title }}
                  </h2>
                </div>
                <p class="mt-3 max-w-[62ch] text-base leading-7 text-muted">
                  {{ selectedCategory.description }}
                </p>
              </div>

              <div
                class="min-w-0 overflow-x-auto overflow-y-hidden pb-1 lg:justify-self-end"
              >
                <div class="flex min-w-max gap-2">
                  <UButton
                    :color="
                      selectedGroupId === ALL_GROUPS ? 'primary' : 'neutral'
                    "
                    :variant="selectedGroupId === ALL_GROUPS ? 'soft' : 'ghost'"
                    label="全部主题"
                    @click="
                      () => {
                        selectedGroupId = ALL_GROUPS;
                      }
                    "
                  />
                  <UButton
                    v-for="group in selectedCategory.groups"
                    :key="group.id"
                    :color="
                      selectedGroupId === group.id ? 'primary' : 'neutral'
                    "
                    :variant="selectedGroupId === group.id ? 'soft' : 'ghost'"
                    :label="group.title"
                    @click="
                      () => {
                        selectedGroupId = group.id;
                      }
                    "
                  />
                </div>
              </div>
            </div>

            <div class="space-y-12">
              <section
                v-for="group in visibleGroups"
                :key="group.id"
                :id="group.id"
                class="scroll-mt-24 space-y-4"
                :aria-labelledby="`${group.id}-title`"
              >
                <div>
                  <div class="flex items-center gap-3">
                    <h3
                      :id="`${group.id}-title`"
                      class="font-display text-xl font-semibold text-highlighted"
                    >
                      {{ group.title }}
                    </h3>
                    <span class="text-sm text-muted"
                      >{{ group.items.length }} 个站点</span
                    >
                  </div>
                  <p class="mt-1 max-w-[62ch] text-sm leading-6 text-muted">
                    {{ group.description }}
                  </p>
                </div>
                <div
                  class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"
                >
                  <NavigationLinkCard
                    v-for="item in group.items"
                    :key="item.id"
                    :entry="entryFor(item)"
                  />
                </div>
              </section>
            </div>
          </template>
        </div>
      </template>
    </section>
  </div>
</template>
