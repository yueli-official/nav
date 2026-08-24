<script setup lang="ts">
import { PageHeader } from "@yueli/ui/admin";

defineOptions({ inheritAttrs: false });

const props = withDefaults(
  defineProps<{
    id: string;
    title: string;
    icon?: string;
    bodyClass?: string;
    mainId?: string;
  }>(),
  { bodyClass: "", mainId: "manage-main" },
);
const slots = useSlots();
const headingId = computed(() => `${props.id}-title`);
</script>

<template>
  <section
    v-bind="$attrs"
    :id="id"
    :aria-labelledby="headingId"
    class="min-w-0 space-y-5"
    data-manage-page
  >
    <PageHeader :title="title" :icon="icon" :heading-id="headingId">
      <template v-if="slots.actions" #actions><slot name="actions" /></template>
    </PageHeader>

    <div
      v-if="slots.toolbar || slots['toolbar-left'] || slots['toolbar-right']"
      class="yueli-card flex flex-wrap items-center justify-between gap-3 px-4 py-3 sm:px-5"
      data-manage-page-toolbar
    >
      <div class="flex min-w-0 flex-1 flex-wrap items-center gap-2">
        <slot name="toolbar-left" />
        <slot name="toolbar" />
      </div>
      <div v-if="slots['toolbar-right']" class="flex items-center gap-2">
        <slot name="toolbar-right" />
      </div>
    </div>

    <div class="min-w-0" :class="bodyClass"><slot /></div>
    <footer v-if="slots.footer"><slot name="footer" /></footer>
  </section>
</template>
