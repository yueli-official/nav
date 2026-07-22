<script setup lang="ts">
import { PageHeader } from '@yueli/ui/dashboard/pattern'
import {
  ManageActiveFilters,
  ManageClientBoundary,
  ManageCollectionFooter,
  ManageCollectionToolbar,
  ManageEmpty,
  ManageLifecycleTabs,
  ManagePageSelection,
  ManageRowShell,
  ManageSortDirectionButton,
  ManageTaxonomyChips,
  SkeletonList,
} from "@platform/manage/components";
import {
  manageCollectionQueryFingerprint,
  serializeManageCollectionQuery,
  type ManageCollectionDefinition,
} from "@platform/manage/collection";
import { useManageCollectionState } from "@platform/manage/use-manage-collection-state";
import { useManageSelection } from "@platform/manage/use-manage-selection";
import { rel } from "@platform/ui/date";
import type {
  AdminNavigationLink,
  AdminNavigationResponse,
  NavigationCategory,
  NavigationLifecycleCounts,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "站点管理 · 月离导航" });

const ALL = "__all__";
const emptyCounts: NavigationLifecycleCounts = {
  all: 0,
  published: 0,
  draft: 0,
  archived: 0,
};
const collectionDefinition = {
  resourceKind: "navigation-link",
  statuses: ["", "published", "draft", "archived"],
  views: ["list"],
  sortKeys: ["updated", "title", "published", "default"],
  pageSizes: [15, 30, 60],
  defaultStatus: "",
  defaultView: "list",
  defaultSort: "updated",
  defaultDirection: "desc",
  defaultPageSize: 15,
  pagination: "server",
  selection: "page",
  filters: ["category", "group", "tag"],
  quickEditFields: ["title", "url", "status", "category", "group", "tags"],
  bulkActions: ["publish", "draft", "archive", "delete"],
} as const satisfies ManageCollectionDefinition;

const { isAdmin, user } = useAuth();
const { call } = useApi();
const route = useRoute();
const router = useRouter();
const {
  status,
  searchInput,
  q,
  page,
  size,
  sort,
  direction,
  state: collectionState,
  filterModel,
} = useManageCollectionState({
  definition: collectionDefinition,
  routeQuery: computed(() => route.query),
  replaceQuery: (query) => router.replace({ query }),
});

const categoryId = filterModel("category", ALL);
const groupId = filterModel("group", ALL);
const tag = filterModel("tag", ALL);
const editorOpen = ref(false);
const editingLink = ref<AdminNavigationLink>();
const { data, pending, error, refresh } = await useAsyncData(
  "admin-navigation-links",
  () =>
    call<AdminNavigationResponse>("/api/v1/admin/nav/links", {
      query: {
        q: q.value || undefined,
        status: status.value || undefined,
        categoryId: categoryId.value === ALL ? undefined : categoryId.value,
        groupId: groupId.value === ALL ? undefined : groupId.value,
        tag: tag.value === ALL ? undefined : tag.value,
        sort: sort.value,
        direction: direction.value,
        page: page.value,
        size: size.value,
      },
    }),
  {
    server: false,
    watch: [q, status, categoryId, groupId, tag, sort, direction, page, size],
    default: () => ({
      links: [],
      categories: [],
      tags: [],
      counts: emptyCounts,
      total: 0,
      page: 1,
      size: 15,
    }),
  },
);

const links = computed(() => data.value?.links ?? []);
const categories = computed<NavigationCategory[]>(
  () => data.value?.categories ?? [],
);
const total = computed(() => data.value?.total ?? 0);
const counts = computed(() => data.value?.counts ?? emptyCounts);
const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / size.value)),
);
const selectedCategory = computed(() =>
  categories.value.find((item) => item.id === categoryId.value),
);
const categoryItems = computed(() => [
  { label: "全部分类", value: ALL },
  ...categories.value.map((item) => ({ label: item.title, value: item.id })),
]);
const groupItems = computed(() => [
  { label: "全部主题", value: ALL },
  ...(
    selectedCategory.value?.groups ??
    categories.value.flatMap((item) => item.groups)
  ).map((item) => ({ label: item.title, value: item.id })),
]);
const tagItems = computed(() => [
  { label: "全部标签", value: ALL },
  ...(data.value?.tags ?? []).map((item) => ({
    label: `${item.name} · ${item.linkCount}`,
    value: item.name,
  })),
]);
const sortItems = [
  { label: "最近更新", value: "updated" },
  { label: "标题", value: "title" },
  { label: "发布日期", value: "published" },
  { label: "手动顺序", value: "default" },
];
const tabs = computed(() => [
  { key: "", label: "全部", count: counts.value.all },
  { key: "published", label: "已发布", count: counts.value.published },
  { key: "draft", label: "草稿", count: counts.value.draft },
  { key: "archived", label: "归档", count: counts.value.archived },
]);
watch(categoryId, () => {
  if (!groupItems.value.some((item) => item.value === groupId.value))
    groupId.value = ALL;
});
watch(
  totalPages,
  (lastPage) => {
    if (page.value > lastPage) page.value = lastPage;
  },
  { flush: "sync" },
);

const activeFilters = computed(() => [
  ...(categoryId.value !== ALL
    ? [
        {
          key: "category",
          label: `分类：${categoryItems.value.find((item) => item.value === categoryId.value)?.label}`,
        },
      ]
    : []),
  ...(groupId.value !== ALL
    ? [
        {
          key: "group",
          label: `主题：${groupItems.value.find((item) => item.value === groupId.value)?.label}`,
        },
      ]
    : []),
  ...(tag.value !== ALL ? [{ key: "tag", label: `标签：${tag.value}` }] : []),
]);
function removeFilter(key: string) {
  if (key === "category") categoryId.value = ALL;
  if (key === "group") groupId.value = ALL;
  if (key === "tag") tag.value = ALL;
}
function clearFilters() {
  categoryId.value = ALL;
  groupId.value = ALL;
  tag.value = ALL;
}

function categoryLabel(link: AdminNavigationLink) {
  const category = categories.value.find((item) => item.id === link.categoryId);
  const group = category?.groups.find((item) => item.id === link.groupId);
  return [category?.title, group?.title].filter(Boolean).join(" / ");
}
function openCreate() {
  if (!isAdmin.value) return;
  editingLink.value = undefined;
  editorOpen.value = true;
}
function openEdit(link: AdminNavigationLink) {
  if (!isAdmin.value) return;
  editingLink.value = link;
  editorOpen.value = true;
}
const selectionResetKey = computed(() =>
  manageCollectionQueryFingerprint(
    serializeManageCollectionQuery(collectionState.value, collectionDefinition),
  ),
);
const {
  selectedIds,
  selectionCount,
  isPageSelected,
  isPageIndeterminate,
  isSelected,
  toggleOne,
  togglePage,
  replace: replaceSelection,
  clear: clearSelection,
} = useManageSelection({
  visibleIds: computed(() => links.value.map((link) => link.id)),
  filteredTotal: total,
  resetKey: selectionResetKey,
});
const batchItems = [
  { label: "发布", value: "publish" },
  { label: "转为草稿", value: "draft" },
  { label: "归档", value: "archive" },
  { label: "删除", value: "delete" },
];
const batchAction = ref<string>();
const batchBusy = ref(false);
const batchMessage = ref("");
const batchDeleteOpen = ref(false);
async function applyBatch() {
  if (
    !isAdmin.value ||
    !batchAction.value ||
    !selectedIds.value.length ||
    batchBusy.value
  )
    return;
  if (batchAction.value === "delete") {
    batchDeleteOpen.value = true;
    return;
  }
  await executeBatch();
}
async function executeBatch() {
  if (!batchAction.value || !selectedIds.value.length || batchBusy.value)
    return;
  batchBusy.value = true;
  batchMessage.value = "";
  try {
    const result = await call<{ changed: number; failedIds: string[] }>(
      "/api/v1/admin/nav/links/bulk",
      {
        method: "POST",
        body: { ids: selectedIds.value, action: batchAction.value },
      },
    );
    batchMessage.value = result.failedIds.length
      ? `已处理 ${result.changed} 个，${result.failedIds.length} 个失败；失败项已保留选择。`
      : `已处理 ${result.changed} 个站点`;
    if (result.failedIds.length) replaceSelection(result.failedIds);
    else clearSelection();
    batchDeleteOpen.value = false;
    batchAction.value = undefined;
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    batchMessage.value = apiError.data?.message || "批量操作失败，选择已保留。";
  } finally {
    batchBusy.value = false;
  }
}
</script>

<template>
  <div class="space-y-6">
    <PageHeader title="站点管理">
      <template #subtitle>
        已登录：<span class="text-default">{{
          user?.name || user?.email
        }}</span>
        · 搜索、筛选并批量维护导航入口。
      </template>
      <template #actions>
        <UButton
          icon="i-tabler-plus"
          label="添加站点"
          :disabled="!isAdmin"
          @click="openCreate"
        />
      </template>
    </PageHeader>

    <UAlert
      v-if="!isAdmin"
      color="warning"
      icon="i-tabler-lock"
      title="当前账号没有管理权限"
      description="仅 Nav 管理员可以修改站点内容。"
    />

    <ManageClientBoundary :rows="6">
      <div class="space-y-5" :inert="batchBusy" :aria-busy="batchBusy">
        <ManageLifecycleTabs v-model="status" :items="tabs" />

        <ManageCollectionToolbar
          v-model:search="searchInput"
          search-placeholder="搜索名称、网址或简介…"
          :filter-count="activeFilters.length"
        >
          <template #filters>
            <USelectMenu
              v-model="categoryId"
              :items="categoryItems"
              value-key="value"
              :search-input="{ placeholder: '搜索分类…' }"
              icon="i-tabler-folders"
              size="sm"
              aria-label="筛选分类"
            />
            <USelectMenu
              v-model="groupId"
              :items="groupItems"
              value-key="value"
              :search-input="{ placeholder: '搜索主题…' }"
              icon="i-tabler-layout-list"
              size="sm"
              aria-label="筛选主题"
            />
            <USelectMenu
              v-model="tag"
              :items="tagItems"
              value-key="value"
              :search-input="{ placeholder: '搜索标签…' }"
              icon="i-tabler-hash"
              size="sm"
              aria-label="筛选标签"
            />
            <USelect
              v-model="sort"
              :items="sortItems"
              value-key="value"
              icon="i-tabler-arrows-sort"
              size="sm"
              aria-label="排序字段"
            />
            <ManageSortDirectionButton v-model="direction" />
          </template>
        </ManageCollectionToolbar>

        <ManageActiveFilters
          :items="activeFilters"
          @remove="removeFilter"
          @clear="clearFilters"
        />

        <UAlert
          v-if="error && !links.length"
          color="error"
          icon="i-tabler-alert-circle"
          title="站点列表加载失败"
          description="请确认 Nav API、数据库和登录状态正常。"
        >
          <template #actions
            ><UButton
              label="重试"
              color="error"
              variant="soft"
              size="sm"
              @click="() => refresh()"
          /></template>
        </UAlert>
        <SkeletonList v-else-if="pending" :rows="6" />
        <ManageEmpty
          v-else-if="!links.length"
          icon="i-tabler-world-off"
          :text="
            q || activeFilters.length
              ? '没有匹配的站点'
              : '这个状态下还没有站点'
          "
        />

        <div
          v-else
          class="overflow-hidden rounded-xl border border-default bg-default"
        >
          <ManageRowShell
            v-for="link in links"
            :key="link.id"
            :selected="isSelected(link.id)"
            :selection-disabled="batchBusy || !isAdmin"
            :selection-label="`选择站点：${link.title}`"
            @select="toggleOne(link.id)"
          >
            <template #media>
              <span
                class="grid size-12 place-items-center rounded-lg bg-primary/10 font-display font-semibold text-primary"
                >{{ link.title.slice(0, 1) }}</span
              >
            </template>
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <button
                  type="button"
                  class="truncate text-left text-sm font-semibold text-highlighted hover:text-primary disabled:cursor-default disabled:hover:text-highlighted"
                  :disabled="!isAdmin"
                  @click="openEdit(link)"
                >
                  {{ link.title }}
                </button>
                <UBadge
                  v-if="link.featured"
                  label="精选"
                  color="primary"
                  variant="subtle"
                  size="sm"
                />
              </div>
              <p class="mt-0.5 truncate text-xs text-muted">{{ link.url }}</p>
              <p class="mt-1 line-clamp-1 text-sm text-toned">
                {{ link.description }}
              </p>
              <div class="mt-1.5 flex flex-wrap items-center gap-1.5">
                <span class="text-xs text-muted">{{
                  categoryLabel(link)
                }}</span>
                <ManageTaxonomyChips
                  :items="
                    link.tags.map((item) => ({
                      key: item,
                      label: item,
                      kind: 'tag',
                    }))
                  "
                />
              </div>
            </div>
            <template #meta>
              <div class="text-xs md:w-36 md:text-right">
                <ClientOnly>
                  <p class="text-muted">
                    {{
                      link.updatedAt
                        ? `更新 ${rel(link.updatedAt)}`
                        : "尚未更新"
                    }}
                  </p>
                  <p class="mt-1 text-dimmed">
                    {{
                      link.publishedAt
                        ? `发布 ${rel(link.publishedAt)}`
                        : "尚未发布"
                    }}
                  </p>
                  <template #fallback><p class="text-dimmed">…</p></template>
                </ClientOnly>
              </div>
            </template>
            <template #actions>
              <UTooltip text="打开站点"
                ><UButton
                  :to="link.url"
                  external
                  target="_blank"
                  icon="i-tabler-external-link"
                  color="neutral"
                  variant="ghost"
                  size="sm"
                  square
                  :aria-label="`打开 ${link.title}`"
              /></UTooltip>
              <UTooltip text="编辑"
                ><UButton
                  icon="i-tabler-pencil"
                  color="neutral"
                  variant="ghost"
                  size="sm"
                  square
                  :aria-label="`编辑 ${link.title}`"
                  :disabled="!isAdmin"
                  @click="openEdit(link)"
              /></UTooltip>
            </template>
          </ManageRowShell>
        </div>

        <ManageCollectionFooter
          v-if="total > 0 || links.length"
          v-model:page="page"
          v-model:size="size"
          :total="total"
          :total-pages="totalPages"
          :page-size-options="[15, 30, 60]"
          label="站点选择、批量操作与分页"
        >
          <template #selection>
            <ManagePageSelection
              :model-value="isPageSelected"
              :indeterminate="isPageIndeterminate"
              :disabled="batchBusy || !isAdmin"
              label="选择当前页站点"
              @update:model-value="togglePage"
            />
            <span
              v-if="batchMessage"
              class="rounded-lg bg-elevated px-2.5 py-1.5 text-xs text-default"
              >{{ batchMessage }}</span
            >
            <template v-if="selectionCount">
              <span class="text-sm text-default"
                >已选 {{ selectionCount }}</span
              >
              <USelect
                v-model="batchAction"
                :items="batchItems"
                value-key="value"
                placeholder="批量操作"
                size="sm"
                class="w-28"
                :disabled="batchBusy || !isAdmin"
                aria-label="批量操作"
              />
              <UButton
                label="应用"
                size="sm"
                color="primary"
                variant="soft"
                :disabled="!batchAction"
                :loading="batchBusy"
                @click="applyBatch"
              />
              <UButton
                label="取消"
                size="sm"
                color="neutral"
                variant="ghost"
                :disabled="batchBusy"
                @click="clearSelection"
              />
            </template>
            <span v-else class="text-xs">共 {{ total }} 个站点</span>
          </template>
        </ManageCollectionFooter>
      </div>
    </ManageClientBoundary>

    <NavigationLinkEditor
      v-model:open="editorOpen"
      :categories="categories"
      :link="editingLink"
      @saved="() => refresh()"
      @deleted="() => refresh()"
    />

    <UModal
      v-model:open="batchDeleteOpen"
      title="批量删除站点"
      description="删除后会立即从公开导航中移除，且不能撤销。"
    >
      <template #body
        ><p class="text-sm leading-6 text-toned">
          确定删除已选择的 {{ selectionCount }} 个站点吗？
        </p></template
      >
      <template #footer
        ><div class="flex w-full justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            :disabled="batchBusy"
            @click="
              () => {
                batchDeleteOpen = false;
              }
            "
          /><UButton
            label="确认删除"
            icon="i-tabler-trash"
            color="error"
            :loading="batchBusy"
            @click="executeBatch"
          /></div
      ></template>
    </UModal>
  </div>
</template>
