<template>
  <!-- 创建/编辑表单对话框 -->
  <YamiDialog
    v-model="show"
    :title="isEditMode ? '编辑分组' : '创建分组'"
    size="lg"
    :confirm-text="isEditMode ? '保存' : '提交'"
    :close-on-click-outside="false"
    @confirm="handleFormSubmit">
    <form class="p-2 space-y-4 text-muted-foreground">
      <!-- ID -->
      <div>
        <label class="label block mb-2">ID</label>
        <input
          v-model="formData.gid"
          type="text"
          class="input input-primary w-full"
          placeholder="请输入ID"
          :disabled="isEditMode"
          required />
      </div>

      <!-- Name -->
      <div>
        <label class="label block mb-2">名称</label>
        <input
          v-model="formData.title"
          type="text"
          class="input input-primary w-full"
          placeholder="请输入名称"
          required />
      </div>

      <!-- Order -->
      <div>
        <label class="label block mb-2">顺序</label>
        <input
          v-model="formData.order"
          type="number"
          class="input input-primary w-full"
          placeholder="请输入排序" />
      </div>
    </form>

    <template #footer>
      <div class="flex w-full">
        <button
          v-if="isEditMode && store.isLogin"
          @click="handleRemove"
          class="btn mr-auto btn-icon btn-destructive">
          <svg
            class="size-5"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24">
            <path
              fill="none"
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M10 11v6m4-6v6m5-11v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
          </svg>
        </button>
        <button type="submit" @click="handleFormSubmit" class="btn btn-primary ml-auto">
          确认
        </button>
        <button class="btn btn-accent ml-4" @click="show = false">取消</button>
      </div>
    </template>
  </YamiDialog>
</template>

<script setup lang="ts">
import type { Group } from '@/models'

const show = ref(false)
const isEditMode = ref(false)
const store = useAuthStore()

const formData = ref<Group>({
  gid: '',
  title: '',
  order: 0,
})

const contextData = ref({
  cid: '',
})

async function handleRemove() {
  const action = '删除分组'
  try {
    await navApi.deleteGroup(contextData.value.cid, formData.value.gid)
    toast.success(action + '成功')
    show.value = false
    location.reload()
  } catch (error: unknown) {
    handleApiError(error, action)
  }
}

async function handleFormSubmit() {
  show.value = false
  let action = '添加分组'
  if (isEditMode.value) {
    action = '修改分组'
  }
  try {
    if (!isEditMode.value) {
      await navApi.createGroup(contextData.value.cid, formData.value)
    } else {
      await navApi.updateGroup(contextData.value.cid, formData.value.gid, formData.value)
    }
    toast.success(action + '成功')
    show.value = false
    location.reload()
  } catch (error: unknown) {
    handleApiError(error, action)
  }
}

// 打开创建模式
function openCreate(cid: string) {
  contextData.value = { cid }
  isEditMode.value = false
  formData.value = {
    gid: '',
    title: '',
    order: 0,
  }
  show.value = true
}

// 打开编辑模式
function openEdit(cid: string, data: Group) {
  contextData.value = { cid }

  console.log(data)
  isEditMode.value = true
  formData.value = { ...data }
  show.value = true
}

defineExpose({
  openCreate,
  openEdit,
})
</script>
