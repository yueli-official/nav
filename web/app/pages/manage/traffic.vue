<script setup lang="ts">
import { AdminMetricCard, AdminOverview, AdminTrendChart } from "@yueli/ui/admin";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "浏览统计" });

interface TrafficTotals {
  views: number;
  uniqueVisitorDays: number;
}

interface TrafficOverview {
  total: TrafficTotals;
  today: TrafficTotals;
  series: Array<TrafficTotals & { day: string }>;
}

const { call } = useApi();
const { data, error, refresh } = await useAsyncData(
  "nav-traffic-overview",
  () => call<TrafficOverview>("/api/v1/admin/nav/traffic"),
  { server: false },
);
const format = new Intl.NumberFormat("zh-CN");
const metrics = computed(() => data.value ? [
  { label: "今日浏览", value: data.value.today.views, icon: "i-tabler-eye" },
  { label: "今日访客", value: data.value.today.uniqueVisitorDays, icon: "i-tabler-users" },
  { label: "累计浏览", value: data.value.total.views, icon: "i-tabler-chart-line" },
] : []);
const trend = computed(() => (data.value?.series ?? []).map((day) => ({
  label: day.day,
  axisLabel: day.day.slice(5).replace("-", "/"),
  value: day.views,
})));
</script>

<template>
  <ManagePage id="traffic" title="浏览统计" icon="i-tabler-chart-bar" main-id="manage-main">
    <template #tools>
      <AdminOverview v-if="!error">
        <div data-admin-metrics>
          <template v-if="data">
            <AdminMetricCard
              v-for="metric in metrics"
              :key="metric.label"
              :label="metric.label"
              :value="format.format(metric.value)"
              :icon="metric.icon"
              tone="neutral"
            />
          </template>
          <template v-else>
            <USkeleton v-for="item in 3" :key="item" class="h-28 rounded-xl" />
          </template>
        </div>
      </AdminOverview>
    </template>
    <UAlert v-if="error" title="浏览统计暂时无法加载" color="error" variant="soft">
      <template #actions><UButton label="重试" color="neutral" variant="soft" @click="refresh()" /></template>
    </UAlert>
    <UCard v-else-if="data">
      <template #header><h2 class="font-semibold text-highlighted">近 14 天浏览趋势</h2></template>
      <AdminTrendChart :points="trend" metric-label="浏览" unit=" 次" empty-label="当前周期还没有浏览记录" />
    </UCard>
    <USkeleton v-else class="h-72 rounded-xl" />
  </ManagePage>
</template>
