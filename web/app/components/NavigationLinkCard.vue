<script setup lang="ts">
import type { NavigationResult } from "~/types/navigation";

const { entry, showContext = false } = defineProps<{
  entry: NavigationResult;
  showContext?: boolean;
}>();
function recordClick() {
  recordNavigationClick(entry.item.id);
}
</script>

<template>
  <NuxtLink
    :to="entry.item.url"
    external
    target="_blank"
    rel="noopener noreferrer"
    class="group flex h-full min-w-0 flex-col rounded-xl border border-default bg-default p-4 transition duration-200 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-[var(--shadow-soft)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:bg-[var(--yueli-surface-card)]"
    @click="recordClick"
  >
    <div class="flex min-w-0 items-start gap-3">
      <span
        class="grid size-10 shrink-0 place-items-center overflow-hidden rounded-lg bg-primary/10 text-primary ring-1 ring-primary/15 dark:bg-[var(--yueli-surface-inset)]"
      >
        <NavigationFavicon
          :id="entry.item.id"
          :title="entry.item.title"
          :revision="entry.item.faviconRevision"
        />
      </span>

      <div class="min-w-0 flex-1">
        <div class="flex min-w-0 items-center gap-2">
          <h3
            class="line-clamp-1 font-display font-semibold text-highlighted transition group-hover:text-primary"
          >
            {{ entry.item.title }}
          </h3>
          <UIcon
            name="i-tabler-arrow-up-right"
            class="size-4 shrink-0 text-dimmed transition group-hover:translate-x-0.5 group-hover:-translate-y-0.5 group-hover:text-primary"
          />
        </div>
        <p class="mt-0.5 line-clamp-1 text-xs text-muted">{{ entry.domain }}</p>
      </div>
    </div>

    <p class="mt-3 line-clamp-2 min-h-[3rem] text-sm leading-6 text-toned">
      {{ entry.item.description }}
    </p>

    <div class="mt-4 flex min-w-0 items-center gap-2">
      <UBadge
        :label="kindLabel(entry.item.kind)"
        color="primary"
        variant="subtle"
        size="sm"
      />
      <UBadge
        v-if="entry.item.clickCount > 0"
        :label="`热门 ${entry.item.clickCount}`"
        icon="i-tabler-flame"
        color="warning"
        variant="subtle"
        size="sm"
      />
      <span
        v-for="tag in entry.item.tags.slice(0, 2)"
        :key="tag"
        class="truncate text-xs text-muted"
        >#{{ tag }}</span
      >
    </div>

    <p
      v-if="showContext"
      class="mt-3 flex items-center gap-1.5 border-t border-default pt-3 text-xs text-muted"
    >
      <UIcon :name="entry.categoryIcon" class="size-4" />
      {{ entry.categoryTitle }} / {{ entry.groupTitle }}
    </p>
  </NuxtLink>
</template>
