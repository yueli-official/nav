<script setup lang="ts">
import { ManagePagination } from "@platform/manage/components";
import type {
  NavigationGroupResponse,
  NavigationResult,
} from "~/types/navigation";

definePageMeta({ width: "wide" });
const route = useRoute();
const router = useRouter();
const groupId = computed(() => String(route.params.groupId || ""));
const page = ref(Math.max(1, Number(route.query.page) || 1));
const size = ref(18);
const sort = ref(route.query.sort === "popular" ? "popular" : "default");
const sortItems = [
  { label: "默认排序", value: "default", icon: "i-tabler-sort-ascending" },
  { label: "最受欢迎", value: "popular", icon: "i-tabler-flame" },
];

const { data, status, error, refresh } =
  await useFetch<NavigationGroupResponse>(
    () => `/api/topics/${encodeURIComponent(groupId.value)}`,
    { query: { page, size, sort }, watch: [groupId, page, size, sort] },
  );
const totalPages = computed(() =>
  Math.max(1, Math.ceil((data.value?.total ?? 0) / size.value)),
);
const entries = computed<NavigationResult[]>(() => {
  const value = data.value;
  if (!value) return [];
  return value.items.map((item) => ({
    item,
    categoryId: value.category.id,
    categoryTitle: value.category.title,
    categoryIcon: value.category.icon,
    groupId: value.group.id,
    groupTitle: value.group.title,
    domain: domainFromUrl(item.url),
    searchText: "",
  }));
});

watch([page, sort], async () => {
  if (!import.meta.client) return;
  const query: Record<string, string> = {};
  if (page.value > 1) query.page = String(page.value);
  if (sort.value === "popular") query.sort = sort.value;
  await router.replace({ query });
});
watch(sort, () => {
  page.value = 1;
});

useSeoMeta({
  title: () =>
    data.value ? `${data.value.group.title} · ${data.value.site.name}` : "主题",
  description: () => data.value?.group.description,
});
</script>

<template>
  <div class="space-y-7">
    <nav class="flex items-center gap-2 text-sm text-muted" aria-label="面包屑">
      <NuxtLink to="/" class="hover:text-primary">首页</NuxtLink>
      <UIcon name="i-tabler-chevron-right" class="size-4" />
      <span>{{ data?.category.title }}</span>
    </nav>

    <header
      v-if="data"
      class="flex flex-col gap-5 border-b border-default pb-7 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="flex min-w-0 items-start gap-4">
        <span
          class="grid size-12 shrink-0 place-items-center rounded-xl bg-primary/10 text-primary"
        >
          <UIcon :name="data.category.icon" class="size-6" />
        </span>
        <div class="min-w-0">
          <p class="text-sm font-medium text-primary">
            {{ data.category.title }}
          </p>
          <h1
            class="font-display mt-1 text-3xl font-bold tracking-tight text-highlighted"
          >
            {{ data.group.title }}
          </h1>
          <p class="mt-2 max-w-[68ch] text-sm leading-6 text-muted">
            {{ data.group.description }}
          </p>
        </div>
      </div>
      <USelect
        v-model="sort"
        :items="sortItems"
        value-key="value"
        class="w-full sm:w-36"
        aria-label="站点排序"
      />
    </header>

    <UAlert
      v-if="error"
      color="error"
      icon="i-tabler-alert-circle"
      title="主题加载失败"
      description="请稍后重试。"
    >
      <template #actions
        ><UButton
          label="重试"
          color="neutral"
          variant="soft"
          @click="refresh()"
      /></template>
    </UAlert>
    <NavigationSkeleton v-else-if="status === 'pending'" />
    <template v-else-if="data">
      <div class="flex items-center justify-between text-sm text-muted">
        <p>共 {{ data.total }} 个站点</p>
        <p v-if="totalPages > 1">第 {{ page }} / {{ totalPages }} 页</p>
      </div>
      <div
        v-if="entries.length"
        class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3"
      >
        <NavigationLinkCard
          v-for="entry in entries"
          :key="entry.item.id"
          :entry="entry"
        />
      </div>
      <div
        v-else
        class="rounded-xl border border-dashed border-default py-14 text-center text-sm text-muted"
      >
        这个主题还没有公开站点。
      </div>
      <div
        v-if="totalPages > 1"
        class="flex justify-center border-t border-default pt-7"
      >
        <ManagePagination v-model="page" :total-pages="totalPages" />
      </div>
    </template>
  </div>
</template>
