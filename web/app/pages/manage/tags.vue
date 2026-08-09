<script setup lang="ts">
import { ManageClientBoundary } from "~/utils/manageComponents";
import {
  CollectionPanel,
  type CollectionPanelMessages,
  type CollectionPanelState,
} from "@yueli/ui/collection/pattern";
import { z } from "zod";
import type { NavigationTag, NavigationTagsResponse } from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "标签管理 · 月离导航" });

const { call } = useApi();
const { can } = useMe();
const canManageStructure = computed(() => can("nav.structure.manage"));
const search = ref("");
const q = ref("");
const page = ref(1);
const size = ref(30);
const panelOpen = ref(false);
const current = ref<NavigationTag>();
const schema = z.object({
  target: z
    .string()
    .trim()
    .min(1, "请输入标签名称")
    .max(80, "标签名称不能超过 80 个字符"),
});
const renameForm = reactive({ target: "" });
const saving = ref(false);
const operationError = ref("");
const deleteOpen = ref(false);
let timer: ReturnType<typeof setTimeout> | undefined;
watch(search, (value) => {
  if (timer) clearTimeout(timer);
  timer = setTimeout(() => {
    q.value = value.trim();
    page.value = 1;
  }, 250);
});
onScopeDispose(() => {
  if (timer) clearTimeout(timer);
});

const { data, pending, error, refresh } = await useAsyncData(
  "nav-tags",
  () =>
    call<NavigationTagsResponse>("/api/v1/admin/nav/tags", {
      query: { q: q.value || undefined },
    }),
  { server: false, watch: [q], default: () => ({ tags: [] }) },
);
const tags = computed(() => data.value?.tags ?? []);
const totalPages = computed(() =>
  Math.max(1, Math.ceil(tags.value.length / size.value)),
);
const pagedTags = computed(() =>
  tags.value.slice((page.value - 1) * size.value, page.value * size.value),
);
const pageSizes = [15, 30, 50] as const;
const panelState = computed<CollectionPanelState>(() =>
  error.value ? "error" : pending.value ? "loading" : "ready",
);
const messages = computed<CollectionPanelMessages>(() => ({
  searchPlaceholder: "搜索标签名称…",
  searchAction: "搜索",
  filtersAction: "筛选",
  activeFilters: (count) => `筛选（${count}）`,
  clearFilters: "清除筛选",
  selectPage: "选择当前页标签",
  selectItem: (label) => `选择标签：${label}`,
  bulkRegion: "标签批量操作",
  selected: (count) => `已选择 ${count} 个标签`,
  selectAllResults: "选择全部结果",
  clearSelection: "取消选择",
  emptyTitle: q.value ? "没有匹配的标签" : "还没有标签",
  emptyDescription: q.value
    ? "请调整搜索词后重试。"
    : "为站点添加标签后会显示在这里。",
  errorTitle: "标签列表加载失败",
  retry: "重新加载",
  showing: (first, last, count) => `显示 ${first}–${last}，共 ${count} 个`,
  pageSize: "每页",
  pageSizeControl: "每页标签数量",
  pageSizeOption: (value) => `${value} 个`,
}));
const tagKey = (tag: NavigationTag) => tag.name;
const tagLabel = (tag: NavigationTag) => tag.name;
watch(
  totalPages,
  (lastPage) => {
    if (page.value > lastPage) page.value = lastPage;
  },
  { flush: "sync" },
);

function submitSearch(value: string) {
  if (timer) clearTimeout(timer);
  search.value = value;
  q.value = value.trim();
  page.value = 1;
}

function openEdit(tag: NavigationTag) {
  if (!canManageStructure.value) return;
  current.value = tag;
  renameForm.target = tag.name;
  operationError.value = "";
  panelOpen.value = true;
}
async function rename() {
  if (
    !canManageStructure.value ||
    !current.value ||
    !renameForm.target.trim() ||
    renameForm.target.trim() === current.value.name ||
    saving.value
  )
    return;
  saving.value = true;
  operationError.value = "";
  try {
    await call("/api/v1/admin/nav/tags/rename", {
      method: "POST",
      body: { source: current.value.name, target: renameForm.target.trim() },
    });
    panelOpen.value = false;
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    operationError.value = apiError.data?.message || "重命名失败，请稍后重试。";
  } finally {
    saving.value = false;
  }
}
function confirmDelete(tag: NavigationTag) {
  if (!canManageStructure.value) return;
  current.value = tag;
  operationError.value = "";
  deleteOpen.value = true;
}
async function remove() {
  if (!canManageStructure.value || !current.value || saving.value) return;
  saving.value = true;
  operationError.value = "";
  try {
    await call("/api/v1/admin/nav/tags/delete", {
      method: "POST",
      body: { name: current.value.name },
    });
    deleteOpen.value = false;
    panelOpen.value = false;
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    operationError.value = apiError.data?.message || "删除失败，请稍后重试。";
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <YAdminPage
    id="tags"
    title="标签治理"
    description="统一重命名、合并重复词或解除关联；变更会作用于所有关联链接。"
    icon="i-tabler-hash"
    main-id="manage-main"
    body-class="mx-auto flex min-h-0 w-full max-w-screen-2xl flex-col gap-4 !overflow-hidden"
  >
    <UAlert
      v-if="!canManageStructure"
      color="warning"
      icon="i-tabler-lock"
      title="当前账号没有管理权限"
      description="标签对当前账号只读。"
    />

    <ManageClientBoundary :rows="6">
      <div class="flex min-h-0 flex-1 flex-col gap-3">
        <CollectionPanel
          data-tag-list-panel
          class="flex min-h-0 flex-1 flex-col [&>[aria-live=polite]]:min-h-0 [&>[aria-live=polite]]:flex-1 [&>[aria-live=polite]]:overflow-y-auto [&>[aria-live=polite]]:overscroll-contain [&>footer]:shrink-0"
          v-model:search="search"
          :items="pagedTags"
          :item-key="tagKey"
          :item-label="tagLabel"
          :messages="messages"
          :state="panelState"
          error-message="请检查 Nav API、数据库和登录状态。"
          :total="tags.length"
          :page="page"
          :page-size="size"
          :page-sizes="pageSizes"
          label="标签列表"
          @search="submitSearch"
          @retry="refresh"
          @page-change="page = $event"
          @page-size-change="size = $event"
        >
          <template #columns>
            <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-3">
              <span>标签与关联站点</span>
              <span class="w-16 text-right">操作</span>
            </div>
          </template>

          <template #item="{ item: tag }">
            <div
              class="grid min-w-0 grid-cols-[minmax(0,1fr)_auto] items-center gap-3"
            >
              <button
                type="button"
                class="flex min-w-0 items-center gap-3 text-left disabled:cursor-default"
                :disabled="!canManageStructure"
                @click="openEdit(tag)"
              >
                <span
                  class="grid size-10 shrink-0 place-items-center rounded-lg bg-primary/10 text-primary"
                >
                  <UIcon name="i-tabler-hash" class="size-5" />
                </span>
                <span class="min-w-0">
                  <span
                    class="block truncate text-sm font-medium text-highlighted"
                    >{{ tag.name }}</span
                  >
                  <span class="mt-0.5 block text-xs text-muted"
                    >关联 {{ tag.linkCount }} 个站点</span
                  >
                </span>
              </button>
              <div class="flex w-16 items-center justify-end">
                <UTooltip text="重命名或合并">
                  <UButton
                    icon="i-tabler-pencil"
                    color="neutral"
                    variant="ghost"
                    size="sm"
                    square
                    :aria-label="`编辑标签：${tag.name}`"
                    :disabled="!canManageStructure"
                    @click="openEdit(tag)"
                  />
                </UTooltip>
              </div>
            </div>
          </template>
        </CollectionPanel>
      </div>
    </ManageClientBoundary>

    <USlideover v-model:open="panelOpen" title="重命名或合并标签">
      <template #body>
        <UForm
          :schema="schema"
          :state="renameForm"
          class="space-y-5"
          @submit="rename"
        >
          <UAlert
            color="neutral"
            variant="subtle"
            icon="i-tabler-arrows-join"
            title="同名即合并"
            :description="`将「${current?.name}」改为已有标签名时，所有关联会合并并自动去重。`"
          />
          <UAlert
            v-if="operationError"
            color="error"
            variant="subtle"
            icon="i-tabler-alert-circle"
            title="操作失败"
            :description="operationError"
          />
          <UFormField name="target" label="新标签名称" required
            ><UInput v-model="renameForm.target" class="w-full" autofocus
          /></UFormField>
          <UButton
            type="submit"
            block
            label="保存"
            icon="i-tabler-device-floppy"
            :loading="saving"
            :disabled="renameForm.target.trim() === current?.name"
          />
        </UForm>
        <div class="mt-8 border-t border-default pt-5">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-medium text-highlighted">危险操作</p>
              <p class="mt-1 text-xs text-muted">
                删除只解除站点关联，不会删除任何站点。
              </p>
            </div>
            <UButton
              type="button"
              label="删除标签"
              icon="i-tabler-trash"
              color="error"
              variant="soft"
              :disabled="!canManageStructure || !current"
              @click="current && confirmDelete(current)"
            />
          </div>
        </div>
      </template>
    </USlideover>

    <UModal
      v-model:open="deleteOpen"
      title="删除标签"
      description="删除只会解除站点关联，不会删除站点。"
    >
      <template #body
        ><p class="text-sm text-toned">
          确定删除「{{ current?.name }}」？当前关联
          {{ current?.linkCount }} 个站点。
        </p>
        <UAlert
          v-if="operationError"
          class="mt-4"
          color="error"
          icon="i-tabler-alert-circle"
          title="操作失败"
          :description="operationError"
      /></template>
      <template #footer
        ><div class="flex w-full justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            @click="
              () => {
                deleteOpen = false;
              }
            "
          /><UButton
            label="确认删除"
            color="error"
            icon="i-tabler-trash"
            :loading="saving"
            @click="remove"
          /></div
      ></template>
    </UModal>
  </YAdminPage>
</template>
