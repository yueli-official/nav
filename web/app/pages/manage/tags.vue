<script setup lang="ts">
import {
  ManageCollectionFooter,
  ManageCollectionToolbar,
  ManageEmpty,
  ManageHeader,
  SkeletonList,
} from "@platform/manage/components";
import { z } from "zod";
import type { NavigationTag, NavigationTagsResponse } from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "标签管理 · 月离导航" });

const { call } = useApi();
const { isAdmin } = useAuth();
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
watch(
  totalPages,
  (lastPage) => {
    if (page.value > lastPage) page.value = lastPage;
  },
  { flush: "sync" },
);

function openEdit(tag: NavigationTag) {
  if (!isAdmin.value) return;
  current.value = tag;
  renameForm.target = tag.name;
  operationError.value = "";
  panelOpen.value = true;
}
async function rename() {
  if (
    !isAdmin.value ||
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
  if (!isAdmin.value) return;
  current.value = tag;
  operationError.value = "";
  deleteOpen.value = true;
}
async function remove() {
  if (!isAdmin.value || !current.value || saving.value) return;
  saving.value = true;
  operationError.value = "";
  try {
    await call("/api/v1/admin/nav/tags/delete", {
      method: "POST",
      body: { name: current.value.name },
    });
    deleteOpen.value = false;
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
  <div class="space-y-6">
    <ManageHeader title="标签管理">
      <template #subtitle
        >标签从站点内容中产生；在这里统一重命名、合并重复词或解除关联。</template
      >
    </ManageHeader>
    <UAlert
      v-if="!isAdmin"
      color="warning"
      icon="i-tabler-lock"
      title="当前账号没有管理权限"
      description="标签对当前账号只读。"
    />

    <ManageCollectionToolbar
      v-model:search="search"
      search-placeholder="搜索标签名称…"
    />
    <UAlert
      v-if="error"
      color="error"
      icon="i-tabler-alert-circle"
      title="标签加载失败"
      description="请检查 Nav API 与数据库状态。"
    />
    <SkeletonList v-else-if="pending" :rows="6" />
    <ManageEmpty
      v-else-if="!pagedTags.length"
      icon="i-tabler-hash-off"
      :text="q ? '没有匹配的标签' : '还没有标签；为站点添加标签后会显示在这里'"
    />

    <div
      v-else
      class="divide-y divide-default overflow-hidden rounded-xl border border-default bg-default"
    >
      <article
        v-for="tag in pagedTags"
        :key="tag.name"
        class="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-3 px-4 py-3"
      >
        <button
          type="button"
          class="flex min-w-0 items-center gap-3 text-left disabled:cursor-default"
          :disabled="!isAdmin"
          @click="openEdit(tag)"
        >
          <span
            class="grid size-10 shrink-0 place-items-center rounded-lg bg-primary/10 text-primary"
            ><UIcon name="i-tabler-hash" class="size-5"
          /></span>
          <span class="min-w-0"
            ><span
              class="block truncate text-sm font-medium text-highlighted"
              >{{ tag.name }}</span
            ><span class="mt-0.5 block text-xs text-muted"
              >关联 {{ tag.linkCount }} 个站点</span
            ></span
          >
        </button>
        <div class="flex items-center gap-1">
          <UTooltip text="重命名或合并"
            ><UButton
              icon="i-tabler-pencil"
              color="neutral"
              variant="ghost"
              size="sm"
              square
              :aria-label="`编辑标签：${tag.name}`"
              :disabled="!isAdmin"
              @click="openEdit(tag)"
          /></UTooltip>
          <UTooltip text="删除标签"
            ><UButton
              icon="i-tabler-trash"
              color="error"
              variant="ghost"
              size="sm"
              square
              :aria-label="`删除标签：${tag.name}`"
              :disabled="!isAdmin"
              @click="confirmDelete(tag)"
          /></UTooltip>
        </div>
      </article>
    </div>

    <ManageCollectionFooter
      v-if="tags.length"
      v-model:page="page"
      v-model:size="size"
      :total="tags.length"
      :total-pages="totalPages"
      label="标签统计与分页"
    >
      <template #selection
        ><span class="text-xs">共 {{ tags.length }} 个标签</span></template
      >
    </ManageCollectionFooter>

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
  </div>
</template>
