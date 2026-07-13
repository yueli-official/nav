<script setup lang="ts">
import {
  ActionFeedbackButton,
  ManageCollectionToolbar,
  ManageEmpty,
  ManageHeader,
  SkeletonList,
} from "@platform/manage/components";
import { useActionFeedback } from "@platform/manage/use-action-feedback";
import { createPlatformNotifier } from "@platform/ui/feedback";
import type {
  AdminNavigationLink,
  AdminNavigationResponse,
  NavigationCategory,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "站点管理 · 月离导航" });

const { isAdmin, user } = useAuth();
const { call } = useApi();
const toast = createPlatformNotifier(useToast());
const query = ref("");
const categoryId = ref("");
const status = ref("");
const editorOpen = ref(false);
const editingLink = ref<AdminNavigationLink>();
const deleteOpen = ref(false);
const deleteTarget = ref<AdminNavigationLink>();
const deleteErrorText = ref("");
const {
  status: deleteStatus,
  pending: markDeleting,
  success: markDeleted,
  error: markDeleteError,
  reset: resetDelete,
} = useActionFeedback({ resetMs: 1600 });
const deleting = computed(() => deleteStatus.value === "pending");

const { data, pending, error, refresh } = await useAsyncData(
  "admin-navigation-links",
  () =>
    call<AdminNavigationResponse>("/api/v1/admin/nav/links", {
      query: {
        q: query.value || undefined,
        categoryId: categoryId.value || undefined,
        status: status.value || undefined,
      },
    }),
  {
    server: false,
    watch: [query, categoryId, status],
    default: () => ({ links: [], categories: [] }),
  },
);

const links = computed(() => data.value?.links ?? []);
const categories = computed<NavigationCategory[]>(
  () => data.value?.categories ?? [],
);
const categoryItems = computed(() => [
  { label: "全部分类", value: "" },
  ...categories.value.map((category) => ({
    label: category.title,
    value: category.id,
  })),
]);
const statusItems = [
  { label: "全部状态", value: "" },
  { label: "已发布", value: "published" },
  { label: "草稿", value: "draft" },
  { label: "已归档", value: "archived" },
];

function categoryLabel(link: AdminNavigationLink) {
  const category = categories.value.find((item) => item.id === link.categoryId);
  const group = category?.groups.find((item) => item.id === link.groupId);
  return [category?.title, group?.title].filter(Boolean).join(" / ");
}

function openCreate() {
  editingLink.value = undefined;
  editorOpen.value = true;
}

function openEdit(link: AdminNavigationLink) {
  editingLink.value = link;
  editorOpen.value = true;
}

function confirmDelete(link: AdminNavigationLink) {
  deleteTarget.value = link;
  deleteErrorText.value = "";
  resetDelete();
  deleteOpen.value = true;
}

async function removeLink() {
  if (!deleteTarget.value || deleting.value) return;
  deleteErrorText.value = "";
  markDeleting();
  try {
    await call(`/api/v1/admin/nav/links/${deleteTarget.value.id}`, {
      method: "DELETE",
    });
    await refresh();
    markDeleted();
    deleteOpen.value = false;
  } catch (deleteError) {
    const apiError = deleteError as { data?: { message?: string } };
    deleteErrorText.value = apiError.data?.message || "请稍后重试。";
    markDeleteError();
    toast.add({
      title: "删除失败",
      description: deleteErrorText.value,
      color: "error",
      icon: "i-tabler-alert-circle",
    });
  }
}

async function onSaved() {
  await refresh();
}
</script>

<template>
  <div class="space-y-6">
    <ManageHeader title="站点管理">
      <template #subtitle>
        已登录：<span class="text-default">{{
          user?.name || user?.email
        }}</span>
        · 新增、编辑和整理公开导航入口。
      </template>
      <template #actions>
        <UButton
          icon="i-tabler-plus"
          label="添加站点"
          :disabled="!isAdmin"
          @click="openCreate"
        />
      </template>
    </ManageHeader>

    <UAlert
      v-if="!isAdmin"
      color="warning"
      icon="i-tabler-lock"
      title="当前账号没有管理权限"
      description="站点内容仅允许 Catalog 中配置的 Nav 管理员修改。"
    />

    <ManageCollectionToolbar
      v-model:search="query"
      search-placeholder="搜索名称、网址或简介"
    >
      <template #filters>
        <USelect
          v-model="categoryId"
          :items="categoryItems"
          value-key="value"
          class="w-36"
          aria-label="筛选分类"
        />
        <USelect
          v-model="status"
          :items="statusItems"
          value-key="value"
          class="w-32"
          aria-label="筛选状态"
        />
      </template>
    </ManageCollectionToolbar>

    <UAlert
      v-if="error && !links.length"
      color="error"
      icon="i-tabler-alert-circle"
      title="站点列表加载失败"
      description="请确认 Nav API、数据库和登录状态正常。"
    >
      <template #actions>
        <UButton
          label="重试"
          color="error"
          variant="soft"
          size="sm"
          @click="() => refresh()"
        />
      </template>
    </UAlert>

    <SkeletonList v-else-if="pending" :rows="6" />
    <ManageEmpty
      v-else-if="!links.length"
      icon="i-tabler-world-off"
      :text="query ? '没有匹配的站点' : '还没有站点，点击右上角添加'"
    />

    <div
      v-else
      class="overflow-hidden rounded-xl border border-default bg-default"
    >
      <article
        v-for="link in links"
        :key="link.id"
        class="grid gap-4 border-b border-default p-4 last:border-b-0 md:grid-cols-[minmax(0,1fr)_auto] md:items-center"
      >
        <div class="flex min-w-0 items-start gap-3">
          <span
            class="grid size-10 shrink-0 place-items-center rounded-lg bg-primary/10 font-semibold text-primary"
          >
            {{ link.title.slice(0, 1) }}
          </span>
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2">
              <h2 class="font-display truncate font-semibold text-highlighted">
                {{ link.title }}
              </h2>
              <UBadge
                :label="
                  link.status === 'published'
                    ? '已发布'
                    : link.status === 'draft'
                      ? '草稿'
                      : '已归档'
                "
                :color="link.status === 'published' ? 'success' : 'neutral'"
                variant="soft"
                size="sm"
              />
              <UBadge
                v-if="link.featured"
                label="精选"
                color="primary"
                variant="subtle"
                size="sm"
              />
            </div>
            <p class="mt-1 truncate text-sm text-muted">{{ link.url }}</p>
            <p class="mt-1 line-clamp-1 text-sm text-toned">
              {{ link.description }}
            </p>
            <p class="mt-2 text-xs text-muted">{{ categoryLabel(link) }}</p>
          </div>
        </div>

        <div class="flex items-center justify-end gap-1">
          <UTooltip text="打开站点">
            <UButton
              :to="link.url"
              external
              target="_blank"
              icon="i-tabler-external-link"
              color="neutral"
              variant="ghost"
              square
              :aria-label="`打开 ${link.title}`"
            />
          </UTooltip>
          <UTooltip text="编辑">
            <UButton
              icon="i-tabler-pencil"
              color="neutral"
              variant="ghost"
              square
              :aria-label="`编辑 ${link.title}`"
              :disabled="!isAdmin"
              @click="openEdit(link)"
            />
          </UTooltip>
          <UTooltip text="删除">
            <UButton
              icon="i-tabler-trash"
              color="error"
              variant="ghost"
              square
              :aria-label="`删除 ${link.title}`"
              :disabled="!isAdmin"
              @click="confirmDelete(link)"
            />
          </UTooltip>
        </div>
      </article>
    </div>

    <NavigationLinkEditor
      v-model:open="editorOpen"
      :categories="categories"
      :link="editingLink"
      @saved="onSaved"
    />

    <UModal
      v-model:open="deleteOpen"
      title="删除站点"
      description="删除后会立即从公开导航中移除。"
    >
      <template #body>
        <p class="text-sm leading-6 text-toned">
          确定删除
          <strong class="text-highlighted">{{ deleteTarget?.title }}</strong>
          吗？这个操作不能撤销。
        </p>
        <UAlert
          v-if="deleteErrorText"
          class="mt-4"
          color="error"
          icon="i-tabler-alert-circle"
          title="未能删除站点"
          :description="deleteErrorText"
        />
      </template>
      <template #footer>
        <div class="flex w-full justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            :disabled="deleting"
            @click="
              () => {
                deleteOpen = false;
              }
            "
          />
          <ActionFeedbackButton
            :status="deleteStatus"
            idle-label="确认删除"
            pending-label="删除中"
            success-label="已删除"
            error-label="重试删除"
            color="error"
            idle-icon="i-tabler-trash"
            @click="removeLink"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>
