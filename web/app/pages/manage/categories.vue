<script setup lang="ts">
import {
  ManageClientBoundary,
  ManageCollectionToolbar,
  ManageEmpty,
  ManageHeader,
  ManageIconPicker,
  SkeletonList,
} from "@platform/manage/components";
import { z } from "zod";
import type {
  NavigationCategory,
  NavigationGroup,
  NavigationStructureResponse,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "分类与主题 · 月离导航" });

const { isAdmin } = useAuth();
const { call } = useApi();
const search = ref("");
const panelOpen = ref(false);
const entityKind = ref<"category" | "group">("category");
const currentId = ref("");
const saving = ref(false);
const saveError = ref("");
const deleteOpen = ref(false);
const deleteKind = ref<"category" | "group">("category");
const deleteId = ref("");
const deleteName = ref("");
const deleting = ref(false);
const deleteError = ref("");
const form = reactive({
  categoryId: "",
  title: "",
  description: "",
  icon: "i-tabler-folder",
  sortOrder: 0,
});
const schema = computed(() =>
  z.object({
    categoryId:
      entityKind.value === "group"
        ? z.string().trim().min(1, "请选择所属分类")
        : z.string(),
    title: z
      .string()
      .trim()
      .min(1, "请输入名称")
      .max(120, "名称不能超过 120 个字符"),
    description: z.string().trim().max(500, "描述不能超过 500 个字符"),
    icon:
      entityKind.value === "category"
        ? z.string().trim().min(1, "请输入图标名称")
        : z.string(),
    sortOrder: z.number().min(0, "排序不能小于 0"),
  }),
);

const { data, pending, error, refresh } = await useAsyncData(
  "nav-structure",
  () => call<NavigationStructureResponse>("/api/v1/admin/nav/structure"),
  { server: false, default: () => ({ categories: [] }) },
);
const categories = computed(() => data.value?.categories ?? []);
const categoryItems = computed(() =>
  categories.value.map((item) => ({ label: item.title, value: item.id })),
);
const visibleCategories = computed(() => {
  const query = search.value.trim().toLocaleLowerCase("zh-CN");
  if (!query) return categories.value;
  return categories.value
    .map((category) => ({
      ...category,
      groups: category.groups.filter((group) =>
        [group.title, group.description]
          .join(" ")
          .toLocaleLowerCase("zh-CN")
          .includes(query),
      ),
    }))
    .filter(
      (category) =>
        [category.title, category.description]
          .join(" ")
          .toLocaleLowerCase("zh-CN")
          .includes(query) || category.groups.length,
    );
});

function openCategory(category?: NavigationCategory) {
  if (!isAdmin.value) return;
  entityKind.value = "category";
  currentId.value = category?.id ?? "";
  Object.assign(form, {
    categoryId: "",
    title: category?.title ?? "",
    description: category?.description ?? "",
    icon: category?.icon ?? "i-tabler-folder",
    sortOrder: category?.sortOrder ?? 0,
  });
  saveError.value = "";
  panelOpen.value = true;
}
function openGroup(category: NavigationCategory, group?: NavigationGroup) {
  if (!isAdmin.value) return;
  entityKind.value = "group";
  currentId.value = group?.id ?? "";
  Object.assign(form, {
    categoryId: group?.categoryId ?? category.id,
    title: group?.title ?? "",
    description: group?.description ?? "",
    icon: "",
    sortOrder: group?.sortOrder ?? category.groups.length,
  });
  saveError.value = "";
  panelOpen.value = true;
}
async function save() {
  if (!isAdmin.value || saving.value) return;
  saving.value = true;
  saveError.value = "";
  const base = entityKind.value === "category" ? "categories" : "groups";
  try {
    await call(
      `/api/v1/admin/nav/${base}${currentId.value ? `/${currentId.value}` : ""}`,
      {
        method: currentId.value ? "PATCH" : "POST",
        body:
          entityKind.value === "category"
            ? {
                title: form.title,
                description: form.description,
                icon: form.icon || "i-tabler-folder",
                sortOrder: form.sortOrder,
              }
            : {
                categoryId: form.categoryId,
                title: form.title,
                description: form.description,
                sortOrder: form.sortOrder,
              },
      },
    );
    panelOpen.value = false;
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    saveError.value = apiError.data?.message || "保存失败，请检查输入后重试。";
  } finally {
    saving.value = false;
  }
}
function confirmDelete(kind: "category" | "group", id: string, name: string) {
  if (!isAdmin.value) return;
  deleteKind.value = kind;
  deleteId.value = id;
  deleteName.value = name;
  deleteError.value = "";
  deleteOpen.value = true;
}
async function remove() {
  if (!isAdmin.value || deleting.value) return;
  deleting.value = true;
  deleteError.value = "";
  try {
    await call(
      `/api/v1/admin/nav/${deleteKind.value === "category" ? "categories" : "groups"}/${deleteId.value}`,
      { method: "DELETE" },
    );
    deleteOpen.value = false;
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    deleteError.value =
      apiError.data?.message ||
      (deleteKind.value === "category"
        ? "请先清空分类下的主题。"
        : "请先移动或删除主题下的站点。");
  } finally {
    deleting.value = false;
  }
}
</script>

<template>
  <div class="space-y-6">
    <ManageHeader title="分类与主题">
      <template #subtitle
        >维护前台左侧导航结构；分类是一级入口，主题用于组织同类站点。</template
      >
      <template #actions
        ><UButton
          icon="i-tabler-folder-plus"
          label="新建分类"
          :disabled="!isAdmin"
          @click="openCategory()"
      /></template>
    </ManageHeader>
    <UAlert
      v-if="!isAdmin"
      color="warning"
      icon="i-tabler-lock"
      title="当前账号没有管理权限"
      description="分类与主题对当前账号只读。"
    />

    <ManageClientBoundary :rows="6">
    <ManageCollectionToolbar
      v-model:search="search"
      search-placeholder="搜索分类、主题或描述…"
    />
    <UAlert
      v-if="error"
      color="error"
      icon="i-tabler-alert-circle"
      title="分类结构加载失败"
      description="请检查 Nav API 与数据库状态。"
    />
    <SkeletonList v-else-if="pending" :rows="6" />
    <ManageEmpty
      v-else-if="!visibleCategories.length"
      icon="i-tabler-folders-off"
      :text="search ? '没有匹配的分类或主题' : '还没有分类'"
    />

    <div v-else class="space-y-4">
      <section
        v-for="category in visibleCategories"
        :key="category.id"
        class="overflow-hidden rounded-xl border border-default bg-default"
      >
        <header
          class="grid gap-3 border-b border-default bg-elevated/35 px-4 py-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center"
        >
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="grid size-10 shrink-0 place-items-center rounded-lg bg-primary/10 text-primary"
              ><UIcon :name="category.icon" class="size-5"
            /></span>
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <h2 class="font-medium text-highlighted">
                  {{ category.title }}
                </h2>
                <UBadge
                  :label="`${category.groups.length} 个主题`"
                  color="neutral"
                  variant="soft"
                  size="sm"
                />
              </div>
              <p class="mt-0.5 line-clamp-1 text-xs text-muted">
                {{ category.description || "未填写描述" }}
              </p>
            </div>
          </div>
          <div class="flex items-center justify-end gap-1">
            <UButton
              icon="i-tabler-plus"
              label="添加主题"
              color="neutral"
              variant="soft"
              size="sm"
              :disabled="!isAdmin"
              @click="openGroup(category)"
            />
            <UTooltip text="编辑分类"
              ><UButton
                icon="i-tabler-pencil"
                color="neutral"
                variant="ghost"
                size="sm"
                square
                :aria-label="`编辑分类：${category.title}`"
                :disabled="!isAdmin"
                @click="openCategory(category)"
            /></UTooltip>
            <UTooltip text="删除分类"
              ><UButton
                icon="i-tabler-trash"
                color="error"
                variant="ghost"
                size="sm"
                square
                :aria-label="`删除分类：${category.title}`"
                :disabled="!isAdmin"
                @click="confirmDelete('category', category.id, category.title)"
            /></UTooltip>
          </div>
        </header>

        <div v-if="category.groups.length" class="divide-y divide-default">
          <article
            v-for="group in category.groups"
            :key="group.id"
            class="grid gap-3 px-4 py-3 pl-8 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center"
          >
            <div class="flex min-w-0 items-start gap-3">
              <UIcon
                name="i-tabler-corner-down-right"
                class="mt-2 size-4 shrink-0 text-dimmed"
              />
              <div class="min-w-0">
                <p class="text-sm font-medium text-highlighted">
                  {{ group.title }}
                </p>
                <p class="mt-0.5 line-clamp-1 text-xs text-muted">
                  {{ group.description || "未填写描述" }}
                </p>
              </div>
            </div>
            <div class="flex items-center justify-end gap-1">
              <span class="mr-2 text-xs text-muted"
                >{{ group.linkCount }} 个站点 · 排序 {{ group.sortOrder }}</span
              >
              <UButton
                icon="i-tabler-pencil"
                color="neutral"
                variant="ghost"
                size="sm"
                square
                :aria-label="`编辑主题：${group.title}`"
                :disabled="!isAdmin"
                @click="openGroup(category, group)"
              />
              <UButton
                icon="i-tabler-trash"
                color="error"
                variant="ghost"
                size="sm"
                square
                :aria-label="`删除主题：${group.title}`"
                :disabled="!isAdmin"
                @click="confirmDelete('group', group.id, group.title)"
              />
            </div>
          </article>
        </div>
        <p v-else class="px-4 py-5 text-center text-sm text-muted">
          还没有主题，添加后才能创建站点。
        </p>
      </section>
    </div>
    </ManageClientBoundary>

    <USlideover
      v-model:open="panelOpen"
      :title="`${currentId ? '编辑' : '新建'}${entityKind === 'category' ? '分类' : '主题'}`"
    >
      <template #body>
        <UForm :schema="schema" :state="form" class="space-y-5" @submit="save">
          <UAlert
            v-if="saveError"
            color="error"
            variant="subtle"
            icon="i-tabler-alert-circle"
            title="保存失败"
            :description="saveError"
          />
          <UFormField
            v-if="entityKind === 'group'"
            name="categoryId"
            label="所属分类"
            required
            ><USelectMenu
              v-model="form.categoryId"
              :items="categoryItems"
              value-key="value"
              :search-input="{ placeholder: '搜索分类…' }"
              class="w-full"
          /></UFormField>
          <UFormField name="title" label="名称" required
            ><UInput v-model="form.title" class="w-full" autofocus
          /></UFormField>
          <UFormField
            v-if="entityKind === 'category'"
            name="icon"
            label="图标"
            help="选择一个易识别的分类图标"
          >
            <ManageIconPicker v-model="form.icon" compact />
          </UFormField>
          <UFormField name="description" label="描述"
            ><UTextarea v-model="form.description" :rows="4" class="w-full"
          /></UFormField>
          <UFormField name="sortOrder" label="排序"
            ><UInputNumber v-model="form.sortOrder" :min="0" class="w-full"
          /></UFormField>
          <UButton
            type="submit"
            block
            label="保存"
            icon="i-tabler-device-floppy"
            :loading="saving"
          />
        </UForm>
      </template>
    </USlideover>

    <UModal
      v-model:open="deleteOpen"
      :title="`删除${deleteKind === 'category' ? '分类' : '主题'}`"
      description="只有不再被下级内容使用时才能删除。"
    >
      <template #body
        ><p class="text-sm text-toned">确定删除「{{ deleteName }}」？</p>
        <UAlert
          v-if="deleteError"
          class="mt-4"
          color="error"
          icon="i-tabler-alert-circle"
          title="无法删除"
          :description="deleteError"
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
            :loading="deleting"
            @click="remove"
          /></div
      ></template>
    </UModal>
  </div>
</template>
