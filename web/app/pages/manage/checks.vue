<script setup lang="ts">
import { ManageClientBoundary } from "~/utils/manageComponents";
import {
  CollectionPanel,
  type CollectionControl,
  type CollectionControlValue,
  type CollectionPanelMessages,
  type CollectionPanelState,
} from "@yueli/ui/collection/pattern";
import type {
  NavigationChecksResponse,
  NavigationHealthCounts,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "站点检查 · 月离导航" });

const ALL = "__all__";
const emptyCounts: NavigationHealthCounts = {
  all: 0,
  unchecked: 0,
  healthy: 0,
  redirected: 0,
  broken: 0,
  timeout: 0,
  error: 0,
};
const { call } = useApi();
const { can } = useMe();
const canRunChecks = computed(() => can("nav.health_check.run"));
const search = ref("");
const q = ref("");
const health = ref("");
const page = ref(1);
const size = ref(15);
const selected = ref(new Set<string>());
const checking = ref(false);
const checkMessage = ref("");
let searchTimer: ReturnType<typeof setTimeout> | undefined;

watch(search, (value) => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    q.value = value.trim();
    page.value = 1;
  }, 250);
});
onScopeDispose(() => {
  clearTimeout(searchTimer);
});
watch(health, () => {
  page.value = 1;
  selected.value = new Set();
});
watch([q, page, size], () => {
  selected.value = new Set();
});

const { data, pending, error, refresh } = await useAsyncData(
  "nav-health-checks",
  () =>
    call<NavigationChecksResponse>("/api/v1/admin/nav/checks", {
      query: {
        q: q.value || undefined,
        health: health.value || undefined,
        page: page.value,
        size: size.value,
      },
    }),
  {
    server: false,
    watch: [q, health, page, size],
    default: () => ({
      links: [],
      counts: emptyCounts,
      total: 0,
      page: 1,
      size: 15,
    }),
  },
);
const links = computed(() => data.value?.links ?? []);
const counts = computed(() => data.value?.counts ?? emptyCounts);
const total = computed(() => data.value?.total ?? 0);
const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / size.value)),
);
const issueCount = computed(
  () =>
    counts.value.redirected +
    counts.value.broken +
    counts.value.timeout +
    counts.value.error,
);
const healthOptions = computed(() => [
  { label: `全部 · ${counts.value.all}`, value: ALL },
  { label: `待处理 · ${issueCount.value}`, value: "issue" },
  { label: `未检查 · ${counts.value.unchecked}`, value: "unchecked" },
  { label: `正常 · ${counts.value.healthy}`, value: "healthy" },
  { label: `重定向 · ${counts.value.redirected}`, value: "redirected" },
  { label: `失效 · ${counts.value.broken}`, value: "broken" },
]);
const controls = computed<CollectionControl[]>(() => [
  {
    kind: "select",
    id: "health",
    label: "检查状态",
    value: health.value || ALL,
    options: healthOptions.value,
    icon: "i-tabler-heartbeat",
    class: "w-40",
  },
]);
const activeFilterCount = computed(() => (health.value ? 1 : 0));
const pageSizes = [15, 30, 50] as const;
const panelState = computed<CollectionPanelState>(() =>
  error.value ? "error" : pending.value ? "loading" : "ready",
);
const messages = computed<CollectionPanelMessages>(() => ({
  searchPlaceholder: "搜索站点名称、地址或描述…",
  searchAction: "搜索",
  filtersAction: "筛选",
  activeFilters: (count) => `筛选（${count}）`,
  clearFilters: "清除筛选",
  selectPage: "选择当前页站点",
  selectItem: (label) => `选择站点：${label}`,
  bulkRegion: "链接检查批量操作",
  selected: (count) => `已选择 ${count} 个站点`,
  selectAllResults: "选择全部结果",
  clearSelection: "取消选择",
  emptyTitle:
    q.value || health.value ? "没有匹配的检查结果" : "还没有可检查的站点",
  emptyDescription:
    q.value || health.value
      ? "请调整搜索或检查状态后重试。"
      : "添加站点后即可在这里运行检查。",
  errorTitle: "检查列表加载失败",
  retry: "重新加载",
  showing: (first, last, count) => `显示 ${first}–${last}，共 ${count} 个`,
  pageSize: "每页",
  pageSizeControl: "每页检查数量",
  pageSizeOption: (value) => `${value} 个`,
}));
const pageIds = computed(() => links.value.map((link) => link.id));
const isPageSelected = computed(
  () =>
    pageIds.value.length > 0 &&
    pageIds.value.every((id) => selected.value.has(id)),
);
const isPageIndeterminate = computed(
  () =>
    !isPageSelected.value && pageIds.value.some((id) => selected.value.has(id)),
);
const selectedCount = computed(() => selected.value.size);
const checkButtonLabel = computed(() =>
  q.value || health.value
    ? `检查筛选结果 ${total.value} 项`
    : `检查全部 ${total.value} 项`,
);
const linkKey = (link: NavigationChecksResponse["links"][number]) => link.id;
const linkLabel = (link: NavigationChecksResponse["links"][number]) =>
  link.title;
const isSelected = (id: string | number) => selected.value.has(String(id));

const healthMeta = {
  unchecked: {
    label: "未检查",
    color: "neutral",
    icon: "i-tabler-clock-question",
  },
  healthy: { label: "正常", color: "success", icon: "i-tabler-circle-check" },
  redirected: { label: "重定向", color: "warning", icon: "i-tabler-route" },
  broken: { label: "失效", color: "error", icon: "i-tabler-link-off" },
  timeout: {
    label: "超时",
    color: "warning",
    icon: "i-tabler-clock-exclamation",
  },
  error: { label: "异常", color: "error", icon: "i-tabler-alert-triangle" },
} as const;

function togglePage(value: boolean) {
  const next = new Set(selected.value);
  pageIds.value.forEach((id) => {
    if (value) next.add(id);
    else next.delete(id);
  });
  selected.value = next;
}
function toggle(id: string, value: boolean) {
  const next = new Set(selected.value);
  if (value) next.add(id);
  else next.delete(id);
  selected.value = next;
}
function submitSearch(value: string) {
  clearTimeout(searchTimer);
  search.value = value;
  q.value = value.trim();
  page.value = 1;
}
function changeControl(id: string, value: CollectionControlValue) {
  if (id === "health" && typeof value === "string") {
    health.value = value === ALL ? "" : value;
  }
}
function clearFilters() {
  health.value = "";
}
function changePageSize(value: number) {
  size.value = value;
  page.value = 1;
}
async function runChecks(scope: "filtered" | "selected", ids: string[] = []) {
  const selectedIDs = ids.length ? ids : [...selected.value];
  if (
    !canRunChecks.value ||
    !total.value ||
    checking.value ||
    (scope === "selected" && !selectedIDs.length)
  )
    return;
  checking.value = true;
  checkMessage.value = "";
  try {
    const result = await call<{ checked: number }>(
      "/api/v1/admin/nav/checks/run",
      {
        method: "POST",
        body:
          scope === "selected"
            ? { scope: "selected", ids: selectedIDs }
            : {
                scope: "filtered",
                q: q.value || undefined,
                health: health.value || undefined,
              },
      },
    );
    checkMessage.value = `已完成 ${result.checked} 个站点检查`;
    selected.value = new Set();
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    checkMessage.value = apiError.data?.message || "检查失败，请稍后重试。";
  } finally {
    checking.value = false;
  }
}
</script>

<template>
  <YAdminPage
    id="checks"
    title="链接检查"
    description="验证站点可访问性；异常结果集中到待处理队列，保留单项重试与批量检查语义。"
    icon="i-tabler-heartbeat"
    main-id="manage-main"
    body-class="mx-auto flex min-h-0 w-full max-w-screen-2xl flex-col gap-4 !overflow-hidden"
  >
    <template #actions>
      <UButton
        icon="i-tabler-heartbeat"
        :label="checkButtonLabel"
        :loading="checking"
        :disabled="!canRunChecks || !total"
        @click="runChecks('filtered')"
      />
    </template>

    <ManageClientBoundary :rows="8">
      <div
        class="flex min-h-0 flex-1 flex-col gap-3"
        :inert="checking"
        :aria-busy="checking"
      >
        <UAlert
          v-if="checkMessage"
          color="neutral"
          variant="subtle"
          icon="i-tabler-info-circle"
          :description="checkMessage"
        />

        <CollectionPanel
          data-check-list-panel
          class="flex min-h-0 flex-1 flex-col [&>[aria-live=polite]]:min-h-0 [&>[aria-live=polite]]:flex-1 [&>[aria-live=polite]]:overflow-y-auto [&>[aria-live=polite]]:overscroll-contain [&>footer]:shrink-0"
          v-model:search="search"
          :items="links"
          :item-key="linkKey"
          :item-label="linkLabel"
          :controls="controls"
          :messages="messages"
          :state="panelState"
          error-message="请检查 Nav API、数据库和登录状态。"
          :total="total"
          :page="page"
          :page-size="size"
          :page-sizes="pageSizes"
          :active-filter-count="activeFilterCount"
          :selection-count="selectedCount"
          :page-selected="isPageSelected"
          :page-indeterminate="isPageIndeterminate"
          :is-selected="isSelected"
          :is-item-selectable="() => canRunChecks && !checking"
          label="链接检查列表"
          selectable
          @search="submitSearch"
          @control-change="changeControl"
          @clear-filters="clearFilters"
          @retry="refresh"
          @toggle-page="togglePage"
          @toggle-item="toggle"
          @clear-selection="selected = new Set()"
          @page-change="page = $event"
          @page-size-change="changePageSize"
        >
          <template #view>
            <UButton
              icon="i-tabler-refresh"
              color="neutral"
              variant="outline"
              size="xs"
              aria-label="刷新检查结果"
              @click="refresh()"
            />
          </template>

          <template #columns>
            <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-3">
              <span>站点与地址</span>
              <span class="hidden w-60 text-right md:block"
                >检查结果与操作</span
              >
            </div>
          </template>

          <template #bulk-actions>
            <UButton
              :label="`检查已选 ${selectedCount} 项`"
              icon="i-tabler-heartbeat"
              size="xs"
              color="primary"
              variant="soft"
              :loading="checking"
              :disabled="!canRunChecks"
              @click="runChecks('selected')"
            />
          </template>

          <template #item="{ item: link }">
            <div
              class="grid min-w-0 gap-3 md:grid-cols-[minmax(0,1fr)_minmax(12rem,18rem)_auto] md:items-center"
            >
              <div class="min-w-0">
                <div class="flex min-w-0 items-center gap-2">
                  <span
                    class="grid size-7 shrink-0 place-items-center overflow-hidden rounded-md bg-primary/10 text-primary"
                  >
                    <NavigationFavicon
                      :id="link.id"
                      :title="link.title"
                      image-class="size-5"
                    />
                  </span>
                  <p class="truncate text-sm font-medium text-highlighted">
                    {{ link.title }}
                  </p>
                </div>
                <p class="mt-1 truncate text-xs text-muted">{{ link.url }}</p>
                <p
                  v-if="link.healthError"
                  class="mt-1 line-clamp-1 text-xs text-error"
                >
                  {{ link.healthError }}
                </p>
              </div>
              <div
                class="flex flex-wrap items-center gap-2 text-xs text-muted md:justify-end"
              >
                <span v-if="link.healthHttpStatus" class="tabular-nums"
                  >HTTP {{ link.healthHttpStatus }}</span
                >
                <span v-if="link.healthLatencyMs" class="tabular-nums"
                  >{{ link.healthLatencyMs }} ms</span
                >
                <UBadge
                  v-bind="healthMeta[link.healthStatus || 'unchecked']"
                  variant="subtle"
                  size="sm"
                />
              </div>
              <UButton
                icon="i-tabler-refresh"
                color="neutral"
                variant="ghost"
                size="sm"
                square
                class="justify-self-end"
                :aria-label="`重新检查 ${link.title}`"
                :disabled="checking || !canRunChecks"
                @click="runChecks('selected', [link.id])"
              />
            </div>
          </template>
        </CollectionPanel>
      </div>
    </ManageClientBoundary>
  </YAdminPage>
</template>
