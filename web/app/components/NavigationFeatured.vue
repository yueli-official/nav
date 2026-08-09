<script setup lang="ts">
import type { NavigationResult } from "~/types/navigation";

const { entries } = defineProps<{ entries: NavigationResult[] }>();
const { openSearch } = useNavigationSearch();
const picks = computed(() => featuredNavigation(entries, 5));
const leadPick = computed(() => picks.value[0]);
const secondaryPicks = computed(() => picks.value.slice(1));

function rankLabel(index: number) {
  return String(index + 1).padStart(2, "0");
}

function recordClick(entry: NavigationResult) {
  recordNavigationClick(entry.item.id);
}
</script>

<template>
  <section
    v-if="leadPick"
    class="overflow-hidden rounded-2xl border border-default bg-default shadow-[var(--shadow-soft)] dark:border-transparent dark:bg-[var(--yueli-surface-card)]"
    aria-labelledby="featured-title"
  >
    <header
      class="flex flex-col gap-4 px-5 py-5 sm:flex-row sm:items-end sm:justify-between sm:px-6"
    >
      <div class="min-w-0">
        <p class="flex items-center gap-2 text-sm font-medium text-primary">
          <UIcon
            name="i-tabler-sparkles"
            class="size-4"
            aria-hidden="true"
          />编辑精选
        </p>
        <h1
          id="featured-title"
          class="font-display mt-1 text-2xl font-bold tracking-tight text-highlighted sm:text-3xl"
        >
          本周值得逛
        </h1>
        <p class="mt-1.5 max-w-2xl text-sm leading-6 text-muted">
          从 {{ entries.length }} 个互联网入口里，挑出适合立即开始的工具与资料。
        </p>
      </div>
      <UButton
        color="neutral"
        variant="soft"
        icon="i-tabler-search"
        label="搜索全部工具"
        class="self-start sm:self-auto"
        @click="openSearch"
      />
    </header>

    <div
      class="grid border-t border-default lg:grid-cols-[minmax(0,1.04fr)_minmax(0,1fr)]"
    >
      <NuxtLink
        :to="leadPick.item.url"
        external
        target="_blank"
        rel="noopener noreferrer"
        class="nav-hero-surface group relative flex min-h-56 min-w-0 flex-col justify-between overflow-hidden p-5 focus-visible:outline-2 focus-visible:outline-offset-[-3px] focus-visible:outline-primary sm:p-6 lg:min-h-64"
        @click="recordClick(leadPick)"
      >
        <span
          class="pointer-events-none absolute -right-2 -top-8 font-display text-[8rem] font-bold leading-none text-primary/8 sm:text-[10rem]"
          aria-hidden="true"
          >01</span
        >
        <div class="relative flex items-start justify-between gap-4">
          <span
            class="grid size-12 shrink-0 place-items-center overflow-hidden rounded-xl bg-default text-primary shadow-sm ring-1 ring-primary/15 dark:bg-[var(--yueli-surface-inset)]"
          >
            <NavigationFavicon
              :id="leadPick.item.id"
              :title="leadPick.item.title"
              :revision="leadPick.item.faviconRevision"
              eager
            />
          </span>
          <span
            class="rounded-full border border-primary/20 bg-default/80 px-3 py-1 text-xs font-semibold tracking-[0.12em] text-primary backdrop-blur dark:bg-[var(--yueli-surface-card)]"
            >TOP PICK</span
          >
        </div>
        <div class="relative mt-8 max-w-xl">
          <div class="flex items-center gap-2">
            <h2
              class="font-display text-2xl font-bold text-highlighted transition group-hover:text-primary sm:text-3xl"
            >
              {{ leadPick.item.title }}
            </h2>
            <UIcon
              name="i-tabler-arrow-up-right"
              class="size-5 shrink-0 text-primary transition group-hover:-translate-y-0.5 group-hover:translate-x-0.5"
              aria-hidden="true"
            />
          </div>
          <p
            class="mt-2 line-clamp-2 text-sm leading-6 text-toned sm:text-base"
          >
            {{ leadPick.item.description }}
          </p>
          <p class="mt-4 flex items-center gap-2 text-xs text-muted">
            <UIcon
              :name="leadPick.categoryIcon"
              class="size-4"
              aria-hidden="true"
            />
            {{ leadPick.categoryTitle }} / {{ leadPick.groupTitle }}
          </p>
        </div>
      </NuxtLink>

      <div
        class="grid gap-3 bg-elevated/45 p-3 dark:bg-[var(--yueli-surface-region)] sm:grid-cols-2 sm:p-4"
      >
        <NuxtLink
          v-for="(entry, index) in secondaryPicks"
          :key="entry.item.id"
          :to="entry.item.url"
          external
          target="_blank"
          rel="noopener noreferrer"
          class="group flex min-h-28 min-w-0 flex-col justify-between rounded-xl border border-default bg-default p-4 transition duration-200 hover:-translate-y-0.5 hover:border-primary/35 hover:shadow-[var(--shadow-soft)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:bg-[var(--yueli-surface-card)]"
          @click="recordClick(entry)"
        >
          <div class="flex min-w-0 items-start gap-3">
            <span
              class="font-display text-sm font-bold tabular-nums text-primary"
              >{{ rankLabel(index + 1) }}</span
            >
            <span
              class="grid size-9 shrink-0 place-items-center overflow-hidden rounded-lg bg-primary/10 text-primary ring-1 ring-primary/15 dark:bg-[var(--yueli-surface-inset)]"
            >
              <NavigationFavicon
                :id="entry.item.id"
                :title="entry.item.title"
                :revision="entry.item.faviconRevision"
                eager
              />
            </span>
            <span class="min-w-0 flex-1">
              <span
                class="line-clamp-1 font-display font-semibold text-highlighted transition group-hover:text-primary"
                >{{ entry.item.title }}</span
              >
              <span class="mt-0.5 line-clamp-1 text-xs text-muted">{{
                entry.domain
              }}</span>
            </span>
          </div>
          <div
            class="mt-4 flex items-center justify-between gap-3 text-xs text-muted"
          >
            <span class="line-clamp-1">{{ entry.categoryTitle }}</span>
            <UIcon
              name="i-tabler-arrow-up-right"
              class="size-4 shrink-0 transition group-hover:-translate-y-0.5 group-hover:translate-x-0.5 group-hover:text-primary"
              aria-hidden="true"
            />
          </div>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>
