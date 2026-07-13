<template>
  <div class="space-y-8" aria-label="正在加载导航内容" aria-busy="true">
    <div class="grid gap-2 md:grid-cols-4 xl:grid-cols-7">
      <USkeleton v-for="index in 7" :key="index" class="h-12 rounded-xl" />
    </div>
    <div class="space-y-4">
      <USkeleton class="h-8 w-44" />
      <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
        <USkeleton v-for="index in 8" :key="index" class="h-48 rounded-xl" />
      </div>
    </div>
  </div>
</template>
