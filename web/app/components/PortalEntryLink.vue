<script setup lang="ts">
import type { PortalEntry } from "~/data/portal";

defineProps<{ entry: PortalEntry }>();
</script>

<template>
  <a
    :href="entry.url"
    class="group flex min-w-0 items-center gap-3 rounded-lg px-3 py-2.5 transition-colors hover:bg-elevated focus-visible:outline-2 focus-visible:outline-offset-[-2px] focus-visible:outline-primary"
  >
    <span class="grid size-8 shrink-0 place-items-center rounded-lg bg-primary/10 text-primary">
      <UIcon :name="entry.icon" class="size-4" aria-hidden="true" />
    </span>
    <span class="min-w-0 flex-1">
      <span class="block break-words text-sm font-semibold text-highlighted group-hover:text-primary">{{ entry.title }}</span>
      <span class="mt-0.5 block text-xs text-muted">{{ entry.site }} · {{ entry.action }}</span>
    </span>
    <UIcon name="i-tabler-arrow-up-right" class="size-4 shrink-0 text-muted group-hover:text-primary" aria-hidden="true" />
  </a>
</template>
