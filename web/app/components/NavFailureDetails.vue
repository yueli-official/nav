<script setup lang="ts">
import type { FailureFeedback } from "@yueli/http-runtime";
defineProps<{ feedback?: FailureFeedback }>();
</script>
<template>
  <details v-if="feedback?.failure" class="mt-2 text-xs text-muted">
    <summary class="cursor-pointer">技术详情</summary>
    <code class="block select-all break-all">{{ feedback.technical.code }}{{ feedback.technical.traceId ? ' · ' + feedback.technical.traceId : '' }}</code>
  </details>
</template>
