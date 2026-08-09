<script setup lang="ts">
import type { NavigationResult } from "~/types/navigation";

const {
  entries,
  placeholder,
  loading = false,
  failed = false,
} = defineProps<{
  entries: NavigationResult[];
  placeholder: string;
  loading?: boolean;
  failed?: boolean;
}>();

const { open, closeSearch } = useNavigationSearch();
const query = ref("");
const searchInput = useTemplateRef<{ focus: () => void }>("search-input");
const results = computed(() => searchNavigation(entries, query.value));
const recommendations = computed(() => featuredNavigation(entries, 6));
const displayedEntries = computed(() =>
  query.value.trim() ? results.value.slice(0, 10) : recommendations.value,
);

watch(open, async (isOpen) => {
  if (!isOpen) return;
  query.value = "";
  await nextTick();
  searchInput.value?.focus();
});

function selectEntry(entry: NavigationResult) {
  recordNavigationClick(entry.item.id);
  closeSearch();
}

function openFirstResult() {
  const first = displayedEntries.value[0];
  if (!first || !import.meta.client) return;
  window.open(first.item.url, "_blank", "noopener,noreferrer");
  selectEntry(first);
}
</script>

<template>
  <UModal
    v-model:open="open"
    title="搜索导航"
    description="搜索月离导航已经收录的工具、文档、社区和站点。"
    :ui="{
      content: 'overflow-hidden sm:max-w-2xl',
      body: 'p-0 sm:p-0',
    }"
  >
    <template #body>
      <div class="border-b border-default p-4">
        <NavigationSearch
          ref="search-input"
          v-model="query"
          :placeholder="placeholder"
          :result-count="query.trim() ? results.length : undefined"
          @submit="openFirstResult"
        />
      </div>

      <div class="max-h-[min(62dvh,34rem)] overflow-y-auto p-3 sm:p-4">
        <div
          class="mb-2 flex items-center justify-between gap-3 px-2 text-xs text-muted"
        >
          <p class="font-medium text-toned">
            {{ query.trim() ? "搜索结果" : "推荐上手" }}
          </p>
          <p v-if="query.trim() && results.length > 10">
            显示前 10 项，共 {{ results.length }} 项
          </p>
        </div>

        <div
          v-if="loading"
          class="space-y-2"
          role="status"
          aria-busy="true"
          aria-label="正在加载导航数据"
        >
          <div
            v-for="index in 5"
            :key="index"
            class="flex items-center gap-3 rounded-xl border border-default p-3"
          >
            <USkeleton class="size-10 shrink-0 rounded-lg" />
            <div class="min-w-0 flex-1 space-y-2">
              <USkeleton class="h-4 w-36" />
              <USkeleton class="h-3 w-3/4" />
            </div>
          </div>
        </div>

        <UAlert
          v-else-if="failed"
          color="error"
          variant="soft"
          icon="i-tabler-alert-circle"
          title="搜索数据暂时不可用"
          description="请关闭搜索后稍后重试。"
        />

        <ul v-else-if="displayedEntries.length" class="space-y-1">
          <li v-for="entry in displayedEntries" :key="entry.item.id">
            <NuxtLink
              :to="entry.item.url"
              external
              target="_blank"
              rel="noopener noreferrer"
              class="group flex min-w-0 items-center gap-3 rounded-xl px-3 py-3 transition duration-200 hover:bg-elevated focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-primary"
              @click="selectEntry(entry)"
            >
              <span
                class="grid size-10 shrink-0 place-items-center overflow-hidden rounded-lg bg-primary/10 text-primary ring-1 ring-primary/15"
              >
                <NavigationFavicon
                  :id="entry.item.id"
                  :title="entry.item.title"
                  :revision="entry.item.faviconRevision"
                />
              </span>
              <span class="min-w-0 flex-1">
                <span class="flex min-w-0 items-center gap-2">
                  <span
                    class="line-clamp-1 font-display font-semibold text-highlighted group-hover:text-primary"
                    >{{ entry.item.title }}</span
                  >
                  <span class="hidden shrink-0 text-xs text-dimmed sm:inline">
                    {{ entry.domain }}
                  </span>
                </span>
                <span class="mt-0.5 line-clamp-1 text-sm text-muted">
                  {{ entry.item.description }}
                </span>
              </span>
              <span
                class="hidden max-w-32 shrink-0 items-center gap-1.5 text-xs text-muted sm:flex"
              >
                <UIcon
                  :name="entry.categoryIcon"
                  class="size-4"
                  aria-hidden="true"
                />
                <span class="line-clamp-1">{{ entry.categoryTitle }}</span>
              </span>
              <UIcon
                name="i-tabler-arrow-up-right"
                class="size-4 shrink-0 text-dimmed transition group-hover:-translate-y-0.5 group-hover:translate-x-0.5 group-hover:text-primary"
                aria-hidden="true"
              />
            </NuxtLink>
          </li>
        </ul>

        <div v-else class="px-6 py-12 text-center" role="status">
          <span
            class="mx-auto grid size-12 place-items-center rounded-xl bg-elevated text-dimmed"
          >
            <UIcon
              name="i-tabler-search-off"
              class="size-6"
              aria-hidden="true"
            />
          </span>
          <h3 class="mt-4 font-display text-lg font-semibold text-highlighted">
            没有找到匹配站点
          </h3>
          <p class="mt-1 text-sm text-muted">
            试试更短的关键词、英文名称、标签或域名。
          </p>
        </div>
      </div>

      <div
        class="flex items-center justify-between gap-3 border-t border-default bg-elevated/45 px-4 py-3 text-xs text-muted"
      >
        <span>Enter 打开第一项</span>
        <span>Esc 关闭</span>
      </div>
    </template>
  </UModal>
</template>
