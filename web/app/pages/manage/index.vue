<script setup lang="ts">
import {
  ManageClientBoundary,
  ManageTaxonomyChips,
} from "~/utils/manageComponents";
import {
  createCollectionRouteQueryCodec,
  createJsonCollectionQueryPolicy,
  type CollectionControl,
  type CollectionControlValue,
  type CollectionPanelMessages,
  type CollectionPanelState,
  type CollectionWorkflow,
} from "@yueli/ui/collection";
import { useVueCollectionWorkflow } from "@yueli/ui/collection/vue";
import { createVueRouterCollectionQuerySync } from "@yueli/ui/collection/vue-router";
import { CollectionPanel } from "@yueli/ui/collection/pattern";
import { rel } from "~/utils/date";
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
const { can } = useMe();
const canSubmit = computed(() => can("nav.link.submit"));
const canUpdate = computed(() => can("nav.link.update"));
const canModerate = computed(() => can("nav.link.moderate"));
const { call } = useApi();
const router = useRouter();
const route = useRoute();
type LinkStatus = "" | "published" | "draft" | "archived";
type LinkSort = "updated" | "title" | "published" | "default";
type LinkDirection = "asc" | "desc";
interface LinkCollectionQuery {
  q: string;
  status: LinkStatus;
  page: number;
  size: number;
  sort: LinkSort;
  direction: LinkDirection;
  category: string;
  group: string;
  tag: string;
}
const statuses = ["", "published", "draft", "archived"] as const;
const sorts = ["updated", "title", "published", "default"] as const;
const directions = ["asc", "desc"] as const;
const pageSizes = [15, 30, 60] as const;
const defaultQuery: LinkCollectionQuery = {
  q: "",
  status: "",
  page: 1,
  size: 15,
  sort: "updated",
  direction: "desc",
  category: ALL,
  group: ALL,
  tag: ALL,
};
const categories = ref<NavigationCategory[]>([]);
const tags = ref<AdminNavigationResponse["tags"]>([]);
const counts = ref<NavigationLifecycleCounts>(emptyCounts);
async function loadLinks(
  nextQuery: Readonly<LinkCollectionQuery>,
  activeWorkflow: CollectionWorkflow<
    AdminNavigationLink,
    string,
    LinkCollectionQuery
  >,
) {
  const token = activeWorkflow.beginLoad();
  try {
    const data = await call<AdminNavigationResponse>(
      "/api/v1/admin/nav/links",
      {
        query: {
          q: nextQuery.q || undefined,
          status: nextQuery.status || undefined,
          categoryId:
            nextQuery.category === ALL ? undefined : nextQuery.category,
          groupId: nextQuery.group === ALL ? undefined : nextQuery.group,
          tag: nextQuery.tag === ALL ? undefined : nextQuery.tag,
          sort: nextQuery.sort,
          direction: nextQuery.direction,
          page: nextQuery.page,
          size: nextQuery.size,
        },
      },
    );
    const lastPage = Math.max(1, Math.ceil(data.total / nextQuery.size));
    if (nextQuery.page > lastPage) {
      activeWorkflow.setQuery({ ...nextQuery, page: lastPage });
      return;
    }
    categories.value = data.categories ?? [];
    tags.value = data.tags ?? [];
    counts.value = data.counts ?? emptyCounts;
    activeWorkflow.resolveLoad(token, {
      items: data.links ?? [],
      total: data.total ?? 0,
    });
  } catch {
    activeWorkflow.rejectLoad(token, {
      key: "nav.links.collection.load_failed",
    });
  }
}
const querySync = createVueRouterCollectionQuerySync({
  router,
  codec: createCollectionRouteQueryCodec({
    q: { kind: "string", default: defaultQuery.q, maxLength: 200 },
    status: { kind: "enum", values: statuses, default: defaultQuery.status },
    page: { kind: "positive-integer", default: defaultQuery.page },
    size: {
      kind: "positive-integer",
      values: pageSizes,
      default: defaultQuery.size,
    },
    sort: { kind: "enum", values: sorts, default: defaultQuery.sort },
    direction: {
      kind: "enum",
      values: directions,
      default: defaultQuery.direction,
    },
    category: {
      kind: "string",
      default: defaultQuery.category,
      maxLength: 200,
    },
    group: { kind: "string", default: defaultQuery.group, maxLength: 200 },
    tag: { kind: "string", default: defaultQuery.tag, maxLength: 200 },
  }),
});
const {
  snapshot: linkCollection,
  workflow: linkWorkflow,
  reload: refresh,
} = useVueCollectionWorkflow({
  initialQuery: defaultQuery,
  queryPolicy: createJsonCollectionQueryPolicy<LinkCollectionQuery>(),
  keyOf: (link: AdminNavigationLink) => link.id,
  querySync,
  dataQueryKey: (query) => JSON.stringify(query),
  load: loadLinks,
});
const collectionQuery = computed(() => linkCollection.value.query);
function updateCollectionQuery(
  patch: Partial<LinkCollectionQuery>,
  resetPage = true,
) {
  linkWorkflow.setQuery({
    ...collectionQuery.value,
    ...patch,
    ...(resetPage ? { page: 1 } : {}),
  });
}
const status = computed({
  get: () => collectionQuery.value.status,
  set: (value: LinkStatus) => updateCollectionQuery({ status: value }),
});
const page = computed({
  get: () => collectionQuery.value.page,
  set: (value: number) => updateCollectionQuery({ page: value }, false),
});
const size = computed({
  get: () => collectionQuery.value.size,
  set: (value: number) => updateCollectionQuery({ size: value }),
});
const sort = computed({
  get: () => collectionQuery.value.sort,
  set: (value: LinkSort) => updateCollectionQuery({ sort: value }),
});
const direction = computed({
  get: () => collectionQuery.value.direction,
  set: (value: LinkDirection) => updateCollectionQuery({ direction: value }),
});
const categoryId = computed({
  get: () => collectionQuery.value.category,
  set: (value: string) =>
    updateCollectionQuery({ category: value, group: ALL }),
});
const groupId = computed({
  get: () => collectionQuery.value.group,
  set: (value: string) => updateCollectionQuery({ group: value }),
});
const tag = computed({
  get: () => collectionQuery.value.tag,
  set: (value: string) => updateCollectionQuery({ tag: value }),
});
const searchInput = ref(collectionQuery.value.q);
let searchTimer: ReturnType<typeof setTimeout> | undefined;
watch(searchInput, (value) => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(
    () => updateCollectionQuery({ q: value.trim() }),
    300,
  );
});
watch(
  () => collectionQuery.value.q,
  (value) => {
    if (searchInput.value !== value) searchInput.value = value;
  },
);
onScopeDispose(() => {
  if (searchTimer) clearTimeout(searchTimer);
});
function submitSearch(value: string) {
  if (searchTimer) clearTimeout(searchTimer);
  searchInput.value = value;
  updateCollectionQuery({ q: value.trim() });
}
const editorOpen = ref(false);
const editingLink = ref<AdminNavigationLink>();
const links = computed(() => linkCollection.value.items);
const total = computed(() => linkCollection.value.total);
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
  ...tags.value.map((item) => ({
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
watch(categoryId, () => {
  if (!groupItems.value.some((item) => item.value === groupId.value))
    groupId.value = ALL;
});
function clearFilters() {
  updateCollectionQuery({ status: "", category: ALL, group: ALL, tag: ALL });
}

function categoryLabel(link: AdminNavigationLink) {
  const category = categories.value.find((item) => item.id === link.categoryId);
  const group = category?.groups.find((item) => item.id === link.groupId);
  return [category?.title, group?.title].filter(Boolean).join(" / ");
}
function openCreate() {
  if (!canSubmit.value) return;
  editingLink.value = undefined;
  editorOpen.value = true;
}
function openEdit(link: AdminNavigationLink) {
  if (!canUpdate.value) return;
  editingLink.value = link;
  editorOpen.value = true;
}
watch(
  [() => route.query.action, canSubmit],
  ([action, allowed]) => {
    if (action === "create" && allowed && !editorOpen.value) openCreate();
  },
  { immediate: true },
);
const selectedIds = computed<readonly string[]>(() =>
  linkCollection.value.selection.mode === "keys"
    ? linkCollection.value.selection.keys
    : [],
);
const selectionCount = computed(() => linkCollection.value.selection.count);
const isPageSelected = computed(() => linkCollection.value.isPageSelected);
const isPageIndeterminate = computed(
  () => linkCollection.value.isPageIndeterminate,
);
function toggleOne(id: string) {
  linkWorkflow.toggleKey(id);
}
function togglePage(selected: boolean) {
  linkWorkflow.togglePage(selected);
}
function clearSelection() {
  linkWorkflow.clearSelection();
}
function replaceSelection(ids: readonly string[]) {
  linkWorkflow.clearSelection();
  for (const id of ids) linkWorkflow.toggleKey(id);
}
const statusItems = computed(() => [
  { value: ALL, label: `全部 · ${counts.value.all}` },
  { value: "published", label: `已发布 · ${counts.value.published}` },
  { value: "draft", label: `草稿 · ${counts.value.draft}` },
  { value: "archived", label: `归档 · ${counts.value.archived}` },
]);
const controls = computed<CollectionControl[]>(() => [
  {
    kind: "select",
    id: "status",
    label: "状态",
    value: status.value || ALL,
    options: statusItems.value,
    class: "w-32",
  },
  {
    kind: "select",
    id: "category",
    label: "分类",
    value: categoryId.value,
    options: categoryItems.value,
    searchPlaceholder: "搜索分类…",
    icon: "i-tabler-folders",
    class: "w-36",
  },
  {
    kind: "select",
    id: "group",
    label: "主题",
    value: groupId.value,
    options: groupItems.value,
    searchPlaceholder: "搜索主题…",
    icon: "i-tabler-layout-list",
    class: "w-36",
  },
  {
    kind: "select",
    id: "tag",
    label: "标签",
    value: tag.value,
    options: tagItems.value,
    searchPlaceholder: "搜索标签…",
    icon: "i-tabler-hash",
    class: "w-36",
  },
  {
    kind: "select",
    id: "sort",
    label: "排序字段",
    value: sort.value,
    options: sortItems,
    icon: "i-tabler-arrows-sort",
    class: "w-32",
  },
  {
    kind: "direction",
    id: "direction",
    label: "排序方向",
    value: direction.value,
    ascendingLabel: "切换为倒序",
    descendingLabel: "切换为正序",
  },
]);
const activeFilterCount = computed(
  () =>
    [
      status.value !== "",
      categoryId.value !== ALL,
      groupId.value !== ALL,
      tag.value !== ALL,
    ].filter(Boolean).length,
);
function changeControl(id: string, value: CollectionControlValue) {
  if (typeof value !== "string") return;
  if (id === "status") {
    if (value === ALL) status.value = "";
    else if (statuses.includes(value as LinkStatus))
      status.value = value as LinkStatus;
  }
  if (id === "category") categoryId.value = value;
  if (id === "group") groupId.value = value;
  if (id === "tag") tag.value = value;
  if (id === "sort" && sorts.includes(value as LinkSort))
    sort.value = value as LinkSort;
  if (id === "direction" && directions.includes(value as LinkDirection))
    direction.value = value as LinkDirection;
}
const messages: CollectionPanelMessages = {
  searchPlaceholder: "搜索名称、网址或简介…",
  searchAction: "搜索",
  filtersAction: "筛选",
  activeFilters: (count) => `筛选（${count}）`,
  clearFilters: "清除筛选",
  selectPage: "选择当前页站点",
  selectItem: (label) => `选择站点：${label}`,
  bulkRegion: "站点批量操作",
  selected: (count) => `已选择 ${count} 个站点`,
  selectAllResults: "选择全部结果",
  clearSelection: "取消选择",
  emptyTitle: "没有匹配的站点",
  emptyDescription: "请调整搜索或筛选条件后重试。",
  errorTitle: "站点列表加载失败",
  retry: "重新加载",
  showing: (first, last, count) => `显示 ${first}–${last}，共 ${count} 个`,
  pageSize: "每页",
  pageSizeControl: "每页站点数量",
  pageSizeOption: (value) => `${value} 个`,
};
const panelState = computed<CollectionPanelState>(() =>
  linkCollection.value.issue
    ? "error"
    : ["idle", "loading", "refreshing"].includes(linkCollection.value.loadState)
      ? "loading"
      : "ready",
);
const linkKey = (link: AdminNavigationLink) => link.id;
const linkLabel = (link: AdminNavigationLink) => link.title;
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
function dismissBatchMessage() {
  batchMessage.value = "";
}
async function applyBatch() {
  if (
    !canModerate.value ||
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
  <YAdminPage
    id="links"
    title="站点链接"
    description="搜索、筛选并维护导航入口；审核动作与内容维护相互独立。"
    icon="i-tabler-world-www"
    main-id="manage-main"
    body-class="mx-auto flex min-h-0 w-full max-w-screen-2xl flex-col gap-4 !overflow-hidden"
  >
      <template #actions>
        <UButton
          icon="i-tabler-plus"
          label="添加站点"
          :disabled="!canSubmit"
          @click="openCreate"
        />
      </template>

    <UAlert
      v-if="!canUpdate && !canSubmit"
      color="warning"
      icon="i-tabler-lock"
      title="当前账号没有管理权限"
      description="当前账号没有可用的链接维护范围。"
    />

    <ManageClientBoundary :rows="6">
      <div
        class="flex min-h-0 flex-1 flex-col gap-3"
        :inert="batchBusy"
        :aria-busy="batchBusy"
      >
        <CollectionPanel
          data-link-list-panel
          class="flex min-h-0 flex-1 flex-col [&>[aria-live=polite]]:min-h-0 [&>[aria-live=polite]]:flex-1 [&>[aria-live=polite]]:overflow-y-auto [&>[aria-live=polite]]:overscroll-contain [&>footer]:shrink-0"
          v-model:search="searchInput"
          :items="links"
          :item-key="linkKey"
          :item-label="linkLabel"
          :controls="controls"
          :messages="messages"
          :state="panelState"
          error-message="请确认 Nav API、数据库和登录状态正常。"
          :total="total"
          :page="page"
          :page-size="size"
          :page-sizes="pageSizes"
          :active-filter-count="activeFilterCount"
          :selection-count="selectionCount"
          :page-selected="isPageSelected"
          :page-indeterminate="isPageIndeterminate"
          :is-selected="linkWorkflow.isSelected"
          :is-item-selectable="() => canModerate && !batchBusy"
          label="站点列表"
          selectable
          @search="submitSearch"
          @control-change="changeControl"
          @clear-filters="clearFilters"
          @retry="refresh"
          @toggle-page="togglePage"
          @toggle-item="toggleOne"
          @clear-selection="clearSelection"
          @page-change="page = $event"
          @page-size-change="size = $event"
        >
          <template #columns>
            <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-3">
              <span>站点、网址与分类</span
              ><span class="hidden w-36 text-right md:block"
                >更新时间与操作</span
              >
            </div>
          </template>
          <template #bulk-actions>
            <USelect
              v-model="batchAction"
              :items="batchItems"
              value-key="value"
              placeholder="批量操作"
              size="xs"
              class="w-28"
              :disabled="!canModerate"
            />
            <UButton
              label="应用"
              size="xs"
              color="primary"
              variant="soft"
              :disabled="!batchAction || !canModerate"
              :loading="batchBusy"
              @click="applyBatch"
            />
          </template>
          <template #item="{ item: link }">
            <div
              class="grid min-w-0 gap-3 sm:grid-cols-[3rem_minmax(0,1fr)] md:grid-cols-[3rem_minmax(0,1fr)_9rem_auto] md:items-center"
            >
              <span
                class="hidden size-12 place-items-center rounded-lg bg-primary/10 font-display font-semibold text-primary sm:grid"
                >{{ link.title.slice(0, 1) }}</span
              >
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <button
                    type="button"
                    class="truncate text-left text-sm font-semibold text-highlighted hover:text-primary disabled:cursor-default"
                    :disabled="!canUpdate"
                    @click="openEdit(link)"
                  >
                    {{ link.title }}</button
                  ><UBadge
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
                  }}</span
                  ><ManageTaxonomyChips
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
              <div class="text-xs md:text-right">
                <ClientOnly
                  ><p class="text-muted">
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
                  <template #fallback
                    ><p class="text-dimmed">…</p></template
                  ></ClientOnly
                >
              </div>
              <div class="flex justify-end gap-1">
                <UTooltip text="打开站点"
                  ><UButton
                    :to="link.url"
                    external
                    target="_blank"
                    icon="i-tabler-external-link"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    square
                    :aria-label="`打开 ${link.title}`" /></UTooltip
                ><UTooltip text="编辑"
                  ><UButton
                    icon="i-tabler-pencil"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    square
                    :aria-label="`编辑 ${link.title}`"
                    :disabled="!canUpdate"
                    @click="openEdit(link)"
                /></UTooltip>
              </div>
            </div>
          </template>
        </CollectionPanel>
        <div
          v-if="batchMessage"
          class="flex items-center justify-between gap-2 rounded-lg border border-default bg-elevated px-3 py-2.5 text-xs text-default"
          role="status"
        >
          <span>{{ batchMessage }}</span
          ><UButton
            icon="i-tabler-x"
            color="neutral"
            variant="ghost"
            size="xs"
            square
            aria-label="关闭批量结果"
            @click="dismissBatchMessage"
          />
        </div>
      </div>
    </ManageClientBoundary>

    <NavigationLinkEditor
      v-model:open="editorOpen"
      :categories="categories"
      :link="editingLink"
      :can-moderate="canModerate"
      :can-delete="canModerate"
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
  </YAdminPage>
</template>
