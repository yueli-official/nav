<script setup lang="ts">
import type {
  NavigationItem,
  NavigationResponse,
  NavigationResult,
} from "~/types/navigation";

definePageMeta({ width: "full" });

const ALL_GROUPS = "__all__";
const SEARCH_PAGE_SIZE = 24;
const GROUP_PREVIEW_SIZE = 4;
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
  { key: "navigation-catalog" },
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
const searchResults = computed(() =>
  searchNavigation(allEntries.value, searchQuery.value),
);
const visibleSearchResults = computed(() =>
  searchResults.value.slice(0, visibleSearchCount.value),
);
const visibleGroups = computed(() => {
  const category = selectedCategory.value;
  if (!category) return [];
  return selectedGroupId.value === ALL_GROUPS
    ? category.groups
    : category.groups.filter((group) => group.id === selectedGroupId.value);
});

watch(
  categories,
  (next) => {
    if (
      next.length &&
      !next.some((category) => category.id === selectedCategoryId.value)
    )
      selectedCategoryId.value = next[0]!.id;
  },
  { immediate: true },
);
watch(selectedCategoryId, () => {
  if (
    !selectedCategory.value?.groups.some(
      (group) => group.id === selectedGroupId.value,
    )
  )
    selectedGroupId.value = ALL_GROUPS;
});
watch(searchQuery, () => {
  visibleSearchCount.value = SEARCH_PAGE_SIZE;
});
watch(
  () => route.fullPath,
  () => {
    const nextQuery = queryValue(route.query.q);
    const nextCategory =
      queryValue(route.query.category) || categories.value[0]?.id || "";
    const nextGroup = queryValue(route.query.group) || ALL_GROUPS;
    if (searchQuery.value !== nextQuery) searchQuery.value = nextQuery;
    if (selectedCategoryId.value !== nextCategory)
      selectedCategoryId.value = nextCategory;
    if (selectedGroupId.value !== nextGroup) selectedGroupId.value = nextGroup;
  },
);
watch([searchQuery, selectedCategoryId, selectedGroupId], async () => {
  if (!import.meta.client) return;
  const query: Record<string, string> = {};
  if (searchQuery.value.trim()) query.q = searchQuery.value.trim();
  if (
    selectedCategoryId.value &&
    selectedCategoryId.value !== categories.value[0]?.id
  )
    query.category = selectedCategoryId.value;
  if (selectedGroupId.value !== ALL_GROUPS) query.group = selectedGroupId.value;
  if (JSON.stringify(query) !== JSON.stringify(route.query))
    await router.replace({ query });
});

function entryFor(item: NavigationItem): NavigationResult {
  const entry = entryMap.value.get(item.id);
  if (!entry) throw new Error(`Missing navigation entry: ${item.id}`);
  return entry;
}
function previewItems(items: NavigationItem[]) {
  return items.slice(0, GROUP_PREVIEW_SIZE);
}
function selectCategory(id: string) {
  selectedCategoryId.value = id;
  selectedGroupId.value = ALL_GROUPS;
}
function clearSearch() {
  searchQuery.value = "";
}
function loadMore() {
  visibleSearchCount.value += SEARCH_PAGE_SIZE;
}
async function retryLoad() {
  await refresh();
}

useSeoMeta({
  title: () => data.value?.site.name,
  description: () => data.value?.site.description,
  ogTitle: () => data.value?.site.title,
  ogDescription: () => data.value?.site.description,
});
</script>

<template>
  <div class="min-w-0 space-y-7 overflow-x-hidden">
    <header id="search" class="scroll-mt-24 border-b border-default pb-7">
      <div
        class="grid gap-5 lg:grid-cols-[minmax(0,1fr)_minmax(22rem,0.72fr)] lg:items-end"
      >
        <div class="min-w-0">
          <p
            class="mb-2 flex items-center gap-2 text-sm font-medium text-primary"
          >
            <UIcon name="i-tabler-compass" class="size-4" />精选互联网入口
          </p>
          <h1
            class="font-display text-balance text-3xl font-bold tracking-tight text-highlighted sm:text-4xl"
          >
            {{ data?.site.title }}
          </h1>
          <p
            class="mt-3 max-w-[66ch] text-sm leading-6 text-muted sm:text-base"
          >
            {{ data?.site.description }}
          </p>
        </div>
        <NavigationSearch
          v-model="searchQuery"
          :placeholder="data?.site.searchPlaceholder ?? ''"
          :result-count="searchQuery ? searchResults.length : undefined"
        />
      </div>
    </header>

    <div
      v-if="error"
      class="rounded-xl border border-error/35 bg-error/5 p-6"
      role="alert"
    >
      <div class="flex items-start gap-4">
        <span
          class="grid size-10 shrink-0 place-items-center rounded-lg bg-error/10 text-error"
          ><UIcon name="i-tabler-alert-triangle" class="size-5"
        /></span>
        <div>
          <h2 class="font-display text-lg font-semibold text-highlighted">
            导航数据加载失败
          </h2>
          <p class="mt-1 text-sm text-muted">
            请稍后重试，或检查 Nav API 与数据库状态。
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

    <div
      v-else-if="data"
      id="catalog"
      class="grid min-w-0 scroll-mt-24 gap-7 lg:grid-cols-[15rem_minmax(0,1fr)] xl:grid-cols-[16rem_minmax(0,1fr)]"
    >
      <aside
        class="min-w-0 lg:sticky lg:top-23 lg:self-start"
        aria-label="导航分类"
      >
        <div class="lg:hidden">
          <USelectMenu
            v-model="selectedCategoryId"
            :items="
              categories.map((category) => ({
                label: category.title,
                value: category.id,
                icon: category.icon,
              }))
            "
            value-key="value"
            :search-input="false"
            class="w-full"
            aria-label="选择分类"
          />
        </div>

        <div
          class="hidden overflow-hidden rounded-xl border border-default bg-default lg:block"
        >
          <div class="border-b border-default px-4 py-3">
            <p
              class="text-xs font-semibold uppercase tracking-[0.12em] text-dimmed"
            >
              分类
            </p>
          </div>
          <nav class="p-2">
            <button
              v-for="category in categories"
              :key="category.id"
              type="button"
              class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm transition focus-visible:outline-2 focus-visible:outline-primary"
              :class="
                selectedCategory?.id === category.id
                  ? 'bg-primary/10 font-medium text-primary'
                  : 'text-muted hover:bg-elevated hover:text-default'
              "
              :aria-current="
                selectedCategory?.id === category.id ? 'page' : undefined
              "
              @click="selectCategory(category.id)"
            >
              <UIcon :name="category.icon" class="size-5 shrink-0" /><span
                class="min-w-0 flex-1 truncate"
                >{{ category.title }}</span
              ><span class="text-xs tabular-nums text-dimmed">{{
                category.groups.reduce((sum, group) => sum + group.linkCount, 0)
              }}</span>
            </button>
          </nav>
          <template v-if="selectedCategory?.groups.length">
            <div class="border-t border-default px-4 py-3">
              <p
                class="text-xs font-semibold uppercase tracking-[0.12em] text-dimmed"
              >
                当前主题
              </p>
            </div>
            <nav class="space-y-0.5 p-2 pt-0">
              <button
                type="button"
                class="flex w-full items-center justify-between rounded-md px-3 py-2 text-left text-xs"
                :class="
                  selectedGroupId === ALL_GROUPS
                    ? 'bg-elevated font-medium text-default'
                    : 'text-muted hover:text-default'
                "
                @click="selectedGroupId = ALL_GROUPS"
              >
                <span>全部主题</span
                ><span>{{
                  selectedCategory.groups.reduce(
                    (sum, group) => sum + group.linkCount,
                    0,
                  )
                }}</span>
              </button>
              <button
                v-for="group in selectedCategory.groups"
                :key="group.id"
                type="button"
                class="flex w-full items-center justify-between rounded-md px-3 py-2 text-left text-xs"
                :class="
                  selectedGroupId === group.id
                    ? 'bg-elevated font-medium text-default'
                    : 'text-muted hover:text-default'
                "
                @click="selectedGroupId = group.id"
              >
                <span class="truncate">{{ group.title }}</span
                ><span class="ml-2 tabular-nums text-dimmed">{{
                  group.linkCount
                }}</span>
              </button>
            </nav>
          </template>
        </div>
      </aside>

      <div class="min-w-0">
        <section
          v-if="searchQuery"
          class="space-y-5"
          aria-labelledby="search-results-title"
        >
          <div
            class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
          >
            <div>
              <h2
                id="search-results-title"
                class="font-display text-2xl font-semibold text-highlighted"
              >
                搜索结果
              </h2>
              <p class="mt-1 text-sm text-muted">
                检索名称、描述、域名、分类与标签，共
                {{ searchResults.length }} 个结果。
              </p>
            </div>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-tabler-x"
              label="清空搜索"
              @click="clearSearch"
            />
          </div>
          <div
            v-if="visibleSearchResults.length"
            class="grid gap-3 sm:grid-cols-2 2xl:grid-cols-3"
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
            class="rounded-xl border border-dashed border-default bg-elevated/35 px-6 py-12 text-center"
          >
            <UIcon
              name="i-tabler-search-off"
              class="mx-auto size-8 text-dimmed"
            />
            <h3
              class="mt-3 font-display text-lg font-semibold text-highlighted"
            >
              没有找到匹配站点
            </h3>
            <p class="mt-1 text-sm text-muted">
              试试更短的关键词、英文名称或域名。
            </p>
            <UButton
              class="mt-4"
              color="neutral"
              variant="soft"
              label="清空搜索"
              @click="clearSearch"
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
              @click="loadMore"
            />
          </div>
        </section>

        <section
          v-else-if="selectedCategory"
          class="space-y-7"
          :aria-labelledby="`${selectedCategory.id}-title`"
        >
          <div
            class="flex min-w-0 items-start gap-3 border-b border-default pb-5"
          >
            <span
              class="grid size-11 shrink-0 place-items-center rounded-xl bg-primary/10 text-primary"
              ><UIcon :name="selectedCategory.icon" class="size-6"
            /></span>
            <div class="min-w-0">
              <h2
                :id="`${selectedCategory.id}-title`"
                class="font-display text-2xl font-semibold tracking-tight text-highlighted"
              >
                {{ selectedCategory.title }}
              </h2>
              <p class="mt-1 max-w-[66ch] text-sm leading-6 text-muted">
                {{ selectedCategory.description }}
              </p>
            </div>
          </div>

          <div class="min-w-0 overflow-x-auto overflow-y-hidden pb-1 lg:hidden">
            <div class="flex min-w-max gap-2">
              <UButton
                :color="selectedGroupId === ALL_GROUPS ? 'primary' : 'neutral'"
                :variant="selectedGroupId === ALL_GROUPS ? 'soft' : 'ghost'"
                label="全部主题"
                @click="
                  () => {
                    selectedGroupId = ALL_GROUPS;
                  }
                "
              /><UButton
                v-for="group in selectedCategory.groups"
                :key="group.id"
                :color="selectedGroupId === group.id ? 'primary' : 'neutral'"
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

          <div class="space-y-9">
            <section
              v-for="group in visibleGroups"
              :key="group.id"
              :id="group.id"
              class="scroll-mt-24 space-y-3"
            >
              <div
                class="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between"
              >
                <div>
                  <h3
                    class="font-display text-lg font-semibold text-highlighted"
                  >
                    {{ group.title }}
                  </h3>
                  <p class="mt-0.5 text-sm text-muted">
                    {{ group.description }}
                  </p>
                </div>
                <div class="flex shrink-0 items-center gap-2">
                  <span class="text-xs text-dimmed"
                    >{{ group.linkCount }} 个站点</span
                  >
                  <UButton
                    v-if="group.linkCount > GROUP_PREVIEW_SIZE"
                    :to="`/topics/${group.id}`"
                    label="查看全部"
                    trailing-icon="i-tabler-arrow-right"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                  />
                </div>
              </div>
              <div class="grid gap-3 sm:grid-cols-2 2xl:grid-cols-3">
                <NavigationLinkCard
                  v-for="item in previewItems(group.items)"
                  :key="item.id"
                  :entry="entryFor(item)"
                />
              </div>
            </section>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>
