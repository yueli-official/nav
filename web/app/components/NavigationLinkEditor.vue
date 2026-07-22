<script setup lang="ts">
import { createPlatformNotifier } from "@platform/ui/feedback";
import { useActionFeedback } from "@yueli/ui/feedback";
import { ActionFeedbackButton } from "@yueli/ui/feedback/pattern";
import { z } from "zod";
import type {
  AdminNavigationLink,
  NavigationCategory,
  NavigationItemKind,
  NavigationStatus,
} from "~/types/navigation";

const { categories, link } = defineProps<{
  categories: NavigationCategory[];
  link?: AdminNavigationLink;
}>();
const emit = defineEmits<{
  saved: [link: AdminNavigationLink];
  deleted: [id: string];
}>();
const open = defineModel<boolean>("open", { required: true });
const { call } = useApi();
const toast = createPlatformNotifier(useToast());

const kindItems = [
  { label: "官方站点", value: "official" },
  { label: "在线工具", value: "tool" },
  { label: "社区", value: "community" },
  { label: "学习平台", value: "learning" },
  { label: "资源", value: "resource" },
  { label: "参考资料", value: "reference" },
  { label: "研究机构", value: "research" },
];
const statusItems = [
  { label: "已发布", value: "published" },
  { label: "草稿", value: "draft" },
  { label: "已归档", value: "archived" },
];
const categoryItems = computed(() =>
  categories.map((category) => ({
    label: category.title,
    value: category.id,
  })),
);
const form = reactive({
  categoryId: "",
  groupId: "",
  title: "",
  url: "",
  description: "",
  tags: [] as string[],
  keywords: [] as string[],
  kind: "tool" as NavigationItemKind,
  featured: false,
  status: "published" as NavigationStatus,
  sortOrder: 0,
});
const groupItems = computed(() => {
  const category = categories.find((item) => item.id === form.categoryId);
  return (category?.groups ?? []).map((group) => ({
    label: group.title,
    value: group.id,
  }));
});
const schema = z.object({
  categoryId: z.string().min(1, "请选择分类"),
  groupId: z.string().min(1, "请选择主题"),
  title: z.string().trim().min(1, "请输入名称").max(200),
  url: z.string().trim().url("请输入完整网址"),
  description: z.string().trim().min(1, "请输入简介").max(500),
});
const saveError = ref("");
const deleteOpen = ref(false);
const deleting = ref(false);
const deleteError = ref("");
const {
  status: saveStatus,
  pending: markSaving,
  success: markSaved,
  error: markSaveError,
  reset: resetSave,
} = useActionFeedback({ resetMs: 1600 });
const saving = computed(() => saveStatus.value === "pending");

watch(
  () => [open.value, link] as const,
  ([isOpen]) => {
    if (!isOpen) return;
    saveError.value = "";
    resetSave();
    const firstCategory = categories[0];
    Object.assign(form, {
      categoryId: link?.categoryId ?? firstCategory?.id ?? "",
      groupId: link?.groupId ?? firstCategory?.groups[0]?.id ?? "",
      title: link?.title ?? "",
      url: link?.url ?? "",
      description: link?.description ?? "",
      tags: [...(link?.tags ?? [])],
      keywords: [...(link?.keywords ?? [])],
      kind: link?.kind ?? "tool",
      featured: link?.featured ?? false,
      status: link?.status ?? "published",
      sortOrder: link?.sortOrder ?? 0,
    });
  },
  { immediate: true },
);

watch(
  () => form.categoryId,
  () => {
    if (!groupItems.value.some((item) => item.value === form.groupId)) {
      form.groupId = groupItems.value[0]?.value ?? "";
    }
  },
);

async function save() {
  if (saving.value) return;
  saveError.value = "";
  markSaving();
  try {
    const path = link
      ? `/api/v1/admin/nav/links/${link.id}`
      : "/api/v1/admin/nav/links";
    const response = await call<{ link: AdminNavigationLink }>(path, {
      method: link ? "PATCH" : "POST",
      body: { ...form },
    });
    emit("saved", response.link);
    markSaved();
    await new Promise((resolve) => setTimeout(resolve, 700));
    open.value = false;
  } catch (error) {
    const apiError = error as { data?: { message?: string } };
    saveError.value = apiError.data?.message || "请检查输入后重试。";
    markSaveError();
    toast.add({
      title: "保存失败",
      description: saveError.value,
      color: "error",
      icon: "i-tabler-alert-circle",
    });
  }
}

async function remove() {
  if (!link || deleting.value) return;
  deleting.value = true;
  deleteError.value = "";
  try {
    await call(`/api/v1/admin/nav/links/${link.id}`, { method: "DELETE" });
    emit("deleted", link.id);
    deleteOpen.value = false;
    open.value = false;
  } catch (error) {
    const apiError = error as { data?: { message?: string } };
    deleteError.value = apiError.data?.message || "删除失败，请稍后重试。";
    toast.add({
      title: "删除失败",
      description: deleteError.value,
      color: "error",
      icon: "i-tabler-alert-circle",
    });
  } finally {
    deleting.value = false;
  }
}

function openDelete() {
  deleteError.value = "";
  deleteOpen.value = true;
}

function closeDelete() {
  deleteOpen.value = false;
}
</script>

<template>
  <div class="contents">
    <UModal
      v-model:open="open"
      :title="link ? '编辑站点' : '添加站点'"
      description="修改后会直接写入导航站数据库。"
      :ui="{ content: 'sm:max-w-2xl', footer: 'justify-end' }"
    >
      <template #body>
        <UForm :schema="schema" :state="form" class="space-y-5" @submit="save">
          <div class="grid gap-4 sm:grid-cols-2">
            <UFormField name="categoryId" label="分类" required>
              <USelect
                v-model="form.categoryId"
                :items="categoryItems"
                value-key="value"
                class="w-full"
              />
            </UFormField>
            <UFormField name="groupId" label="主题" required>
              <USelect
                v-model="form.groupId"
                :items="groupItems"
                value-key="value"
                class="w-full"
              />
            </UFormField>
          </div>

          <UFormField name="title" label="名称" required>
            <UInput
              v-model="form.title"
              class="w-full"
              placeholder="例如 MDN Web Docs"
              autofocus
            />
          </UFormField>
          <UFormField name="url" label="网址" required>
            <UInput
              v-model="form.url"
              class="w-full"
              type="url"
              placeholder="https://example.com/"
            />
          </UFormField>
          <UFormField name="description" label="简介" required>
            <UTextarea
              v-model="form.description"
              class="w-full"
              :rows="3"
              autoresize
              :maxrows="5"
            />
          </UFormField>

          <div class="grid gap-4 sm:grid-cols-2">
            <UFormField label="标签" hint="最多 6 个">
              <UInputTags
                v-model="form.tags"
                class="w-full"
                :max="6"
                placeholder="输入后回车"
              />
            </UFormField>
            <UFormField label="搜索关键词" hint="最多 12 个">
              <UInputTags
                v-model="form.keywords"
                class="w-full"
                :max="12"
                placeholder="输入后回车"
              />
            </UFormField>
          </div>

          <div class="grid gap-4 sm:grid-cols-3">
            <UFormField label="类型">
              <USelect
                v-model="form.kind"
                :items="kindItems"
                value-key="value"
                class="w-full"
              />
            </UFormField>
            <UFormField label="状态">
              <USelect
                v-model="form.status"
                :items="statusItems"
                value-key="value"
                class="w-full"
              />
            </UFormField>
            <UFormField label="排序">
              <UInputNumber v-model="form.sortOrder" class="w-full" :min="0" />
            </UFormField>
          </div>

          <USwitch v-model="form.featured" label="加入首页精选" />

          <UAlert
            v-if="saveError"
            color="error"
            icon="i-tabler-alert-circle"
            title="未能保存站点"
            :description="saveError"
          />

          <div
            class="flex flex-col-reverse gap-2 border-t border-default pt-5 sm:flex-row sm:items-center"
          >
            <UButton
              v-if="link"
              type="button"
              label="删除站点"
              icon="i-tabler-trash"
              color="error"
              variant="ghost"
              :disabled="saving"
              @click="openDelete"
            />
            <div class="ml-auto flex justify-end gap-2">
              <UButton
                label="取消"
                color="neutral"
                variant="outline"
                :disabled="saving"
                @click="
                  () => {
                    open = false;
                  }
                "
              />
              <ActionFeedbackButton
                type="submit"
                :status="saveStatus"
                idle-label="保存"
                pending-label="保存中"
                :success-label="link ? '已更新' : '已添加'"
                error-label="重试保存"
              />
            </div>
          </div>
        </UForm>
      </template>
    </UModal>

    <UModal
      v-model:open="deleteOpen"
      title="删除站点"
      description="删除后会立即从公开导航中移除，且无法撤销。"
    >
      <template #body>
        <p class="text-sm leading-6 text-toned">
          确定删除
          <strong class="text-highlighted">{{ link?.title }}</strong>
          吗？
        </p>
        <UAlert
          v-if="deleteError"
          class="mt-4"
          color="error"
          icon="i-tabler-alert-circle"
          title="未能删除站点"
          :description="deleteError"
          role="alert"
        />
      </template>
      <template #footer>
        <div class="flex w-full justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            :disabled="deleting"
            @click="closeDelete"
          />
          <UButton
            label="确认删除"
            icon="i-tabler-trash"
            color="error"
            :loading="deleting"
            @click="remove"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>
