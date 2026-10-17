<script setup lang="ts">
import type {
  NavigationItem,
  NavigationResponse,
  NavigationResult,
} from "~/types/navigation";

definePageMeta({ width: "full", key: (route) => String(route.params.categoryId) });

const ALL_GROUPS = "__all__";
const GROUP_PREVIEW_SIZE = 4;

const route = useRoute();
const router = useRouter();


function queryValue(value: unknown) {
  if (Array.isArray(value)) return typeof value[0] === "string" ? value[0] : "";
  return typeof value === "string" ? value : "";
}

const selectedCategoryId = computed(() => String(route.params.categoryId));
const selectedGroupId = computed({
  get: () => {
    const group = queryValue(route.query.group);
    return selectedCategory.value?.groups.some((entry) => entry.id === group) ? group : ALL_GROUPS;
  },
  set: (group: string) => {
    void router.push({ query: group === ALL_GROUPS ? {} : { group } });
  },
});
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
    ),
);
const visibleGroups = computed(() => {
  const category = selectedCategory.value;
  if (!category) return [];
  return selectedGroupId.value === ALL_GROUPS
    ? category.groups
    : category.groups.filter((group) => group.id === selectedGroupId.value);
});

if (data.value && !selectedCategory.value) throw createError({ statusCode: 404, statusMessage: "Category not found" });

function entryFor(item: NavigationItem): NavigationResult {
  const entry = entryMap.value.get(item.id);
  if (!entry) throw new Error(`Missing navigation entry: ${item.id}`);
  return entry;
}
function previewItems(items: NavigationItem[]) {
  return items.slice(0, GROUP_PREVIEW_SIZE);
}
async function retryLoad() {
  await refresh();
}

useSeoMeta({
  title: () => `${selectedCategory.value?.title} · ${data.value?.site.name}`,
  description: () => selectedCategory.value?.description,
  ogTitle: () => data.value?.site.title,
  ogDescription: () => data.value?.site.description,
});
</script>

<template>
  <div class="min-w-0 space-y-7 overflow-x-clip">
    <nav class="flex items-center gap-2 text-sm text-muted" aria-label="面包屑"><NuxtLink to="/#categories" class="hover:text-primary">分类</NuxtLink><UIcon name="i-tabler-chevron-right" class="size-4" /><span>{{ selectedCategory?.title }}</span></nav>
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
            请稍后重试。
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
      class="grid min-w-0 scroll-mt-24 gap-7 lg:grid-cols-[13rem_minmax(0,1fr)]"
    >
      <aside v-if="selectedCategory" class="min-w-0" aria-label="分类主题">
        <USelectMenu
          v-model="selectedGroupId"
          :items="[{ label: '全部主题', value: ALL_GROUPS }, ...selectedCategory.groups.map((group) => ({ label: group.title, value: group.id }))]"
          value-key="value"
          :search-input="false"
          class="w-full lg:hidden"
          aria-label="选择主题"
        />
        <div class="hidden rounded-xl border border-default bg-default lg:sticky lg:top-23 lg:block dark:bg-[var(--yueli-surface-region)]">
          <div class="border-b border-default px-4 py-3 text-xs font-semibold text-dimmed">主题</div>
          <nav class="space-y-1 p-2" aria-label="主题导航">
            <NuxtLink
              :to="{ path: route.path }"
              class="flex items-center rounded-lg px-3 py-2.5 text-sm transition focus-visible:outline-2 focus-visible:outline-primary"
              :class="selectedGroupId === ALL_GROUPS ? 'bg-primary/10 font-medium text-primary' : 'text-muted hover:bg-elevated hover:text-default'"
              :aria-current="selectedGroupId === ALL_GROUPS ? 'page' : undefined"
            >全部主题</NuxtLink>
            <NuxtLink
              v-for="group in selectedCategory.groups"
              :key="group.id"
              :to="{ path: route.path, query: { group: group.id } }"
              class="flex min-w-0 items-center gap-2 rounded-lg px-3 py-2.5 text-sm transition focus-visible:outline-2 focus-visible:outline-primary"
              :class="selectedGroupId === group.id ? 'bg-primary/10 font-medium text-primary' : 'text-muted hover:bg-elevated hover:text-default'"
              :aria-current="selectedGroupId === group.id ? 'page' : undefined"
            ><span class="min-w-0 flex-1 break-words">{{ group.title }}</span><span class="shrink-0 text-xs tabular-nums text-dimmed" aria-hidden="true">{{ group.linkCount }}</span></NuxtLink>
          </nav>
        </div>
      </aside>

      <div class="min-w-0">
        <section
          v-if="selectedCategory"
          class="space-y-7"
          :aria-labelledby="`${selectedCategory.id}-title`"
        >
          <div class="space-y-4 border-b border-default pb-5">
            <div class="flex min-w-0 items-start gap-3">
              <span
                class="grid size-11 shrink-0 place-items-center rounded-xl bg-primary/10 text-primary"
                ><UIcon :name="selectedCategory.icon" class="size-6"
              /></span>
              <div class="min-w-0">
                <h1
                  :id="`${selectedCategory.id}-title`"
                  class="font-display text-2xl font-semibold tracking-tight text-highlighted"
                >
                  {{ selectedCategory.title }}
                </h1>
                <p class="mt-1 max-w-[66ch] text-sm leading-6 text-muted">
                  {{ selectedCategory.description }}
                </p>
              </div>
            </div>
          </div>

          <div class="space-y-9">
            <p v-if="!visibleGroups.length" class="py-6 text-sm text-muted">这个分类还没有主题。</p>
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
                  <h2
                    class="font-display text-lg font-semibold text-highlighted"
                  >
                    {{ group.title }}
                  </h2>
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
