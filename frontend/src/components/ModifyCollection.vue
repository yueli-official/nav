<template>
  <!-- 创建/编辑表单对话框 -->
  <YamiDialog
    v-model="show"
    :title="isEditMode ? '编辑收藏' : '创建收藏'"
    size="lg"
    :confirm-text="isEditMode ? '保存' : '提交'"
    :close-on-click-outside="false"
    @confirm="handleFormSubmit">
    <form class="p-2 space-y-4 text-muted-foreground" @submit.prevent="handleFormSubmit">
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

      <!-- Description -->
      <div>
        <label class="label block mb-2">描述</label>
        <textarea
          v-model="formData.description"
          class="textarea textarea-primary w-full"
          rows="3"
          placeholder="请输入描述"
          required />
      </div>

      <!-- Link -->
      <div>
        <label class="label block mb-2">链接</label>
        <input
          v-model="formData.link"
          type="url"
          class="input input-primary w-full"
          placeholder="请输入链接" />
      </div>

      <!-- Tags -->
      <div>
        <label class="label block mb-2">标签 (逗号分隔)</label>
        <input
          v-model="tagsInput"
          type="text"
          class="input input-primary w-full"
          placeholder="请输入标签，用逗号分隔" />
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
const show = ref(false)
const isEditMode = ref(false)
const tagsInput = ref('')
const store = useAuthStore()

// 存储 cid 和 gid，用于提交时携带
const contextData = ref({
  cid: '',
  gid: '',
  oldLink: '',
})

const formData = ref<Collection>({
  title: '',
  description: '',
  link: '',
  tags: [],
})

async function handleRemove() {
  const action = '删除收藏'
  try {
    await navApi.deleteCollection(contextData.value.cid, contextData.value.gid, formData.value.link)
    toast.success(action + '成功')
    show.value = false
    location.reload()
  } catch (error: unknown) {
    handleApiError(error, action)
  }
}

async function handleFormSubmit() {
  show.value = false
  let action = '添加收藏'
  if (isEditMode.value) {
    action = '修改收藏'
  }
  try {
    if (!isEditMode.value) {
      await navApi.createCollection(contextData.value.cid, contextData.value.gid, formData.value)
    } else {
      await navApi.updateCollection(
        contextData.value.cid,
        contextData.value.gid,
        contextData.value.oldLink,
        formData.value,
      )
    }
    toast.success(action + '成功')
    show.value = false
    location.reload()
  } catch (error: unknown) {
    handleApiError(error, action)
  }
}

// 打开创建模式
function openCreate(cid: string, gid: string) {
  isEditMode.value = false
  contextData.value = { cid, gid, oldLink: '' }
  formData.value = {
    title: '',
    description: '',
    link: '',
    tags: [],
  }
  tagsInput.value = ''
  show.value = true
}

// 打开编辑模式
function openEdit(data: Collection, cid: string, gid: string, oldLink: string) {
  isEditMode.value = true
  contextData.value = { cid, gid, oldLink }
  formData.value = { ...data }
  tagsInput.value = (data.tags || []).join(', ')
  show.value = true
}

// 监听标签输入，自动转换为数组
watch(tagsInput, (val) => {
  formData.value.tags = val
    .split(',')
    .map((t) => t.trim())
    .filter(Boolean)
})

defineExpose({
  openCreate,
  openEdit,
})
</script>
