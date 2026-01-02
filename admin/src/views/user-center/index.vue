<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'

import {
  NAlert,
  NButton,
  NCard,
  NDescriptions,
  NDescriptionsItem,
  NDivider,
  NForm,
  NFormItem,
  NGi,
  NGrid,
  NInput,
  NList,
  NListItem,
  NPopover,
  NSpace,
  NStatistic,
  NTabPane,
  NTabs,
  NTag,
  NThing,
  useMessage,
} from 'naive-ui'
import { ScrollContainer, UserAvatar } from '@/components'
import { toRefsUserStore, useUserStore } from '@/stores'
import {
  type OAuthBinding,
  changePassword,
  getAccessInfo,
  getOAuthBindings,
  updateProfile,
} from '@/services/auth'

import type { FormInst, FormItemRule } from 'naive-ui'

defineOptions({ name: 'UserCenter' })

const userStore = useUserStore()
const { user, token } = toRefsUserStore()
const message = useMessage()

const profileFormRef = ref<FormInst | null>(null)
const passwordFormRef = ref<FormInst | null>(null)
const oauthLoading = ref(false)
const oauthBindings = ref<OAuthBinding[]>([])

const profileForm = reactive({
  nickname: '',
  email: '',
  avatar: '',
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const profileRules: Record<string, FormItemRule[]> = {
  nickname: [{ required: true, message: '请输入昵称', trigger: ['blur', 'input'] }],
  email: [{ type: 'email', message: '请输入有效邮箱', trigger: ['blur', 'input'] }],
}

const passwordRules: Record<string, FormItemRule[]> = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: ['blur', 'input'] }],
  newPassword: [{ required: true, message: '请输入新密码', trigger: ['blur', 'input'] }],
  confirmPassword: [
    {
      required: true,
      trigger: ['blur', 'input'],
      validator: (_rule, value) => value === passwordForm.newPassword,
      message: '两次输入的密码不一致',
    },
  ],
}

const roles = computed(() => user.value.roles || [])
const permissions = computed(() => user.value.permissions || [])

async function loadAccessInfo() {
  const data = await getAccessInfo()
  userStore.setAuth({
    token: token.value || '',
    user: {
      id: data.user.id,
      username: data.user.username,
      nickname: data.user.nickname,
      email: data.user.email,
      avatar: data.user.avatar,
      roles: data.roles,
      permissions: data.permissions,
      createdAt: data.user.createdAt,
      updatedAt: data.user.updatedAt,
    },
  })
  profileForm.nickname = data.user.nickname
  profileForm.email = data.user.email
  profileForm.avatar = data.user.avatar
}

async function handleProfileSubmit() {
  profileFormRef.value?.validate(async (errors) => {
    if (errors) return
    const updated = await updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      avatar: profileForm.avatar,
    })
    userStore.setAuth({
      token: token.value || '',
      user: {
        id: updated.id,
        username: updated.username,
        nickname: updated.nickname,
        email: updated.email,
        avatar: updated.avatar,
        roles: user.value.roles,
        permissions: user.value.permissions,
        createdAt: updated.createdAt,
        updatedAt: updated.updatedAt,
      },
    })
    message.success('个人信息更新成功')
  })
}

async function handlePasswordSubmit() {
  passwordFormRef.value?.validate(async (errors) => {
    if (errors) return
    await changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    message.success('密码修改成功')
  })
}

async function loadOAuthBindings() {
  oauthLoading.value = true
  try {
    oauthBindings.value = await getOAuthBindings()
  } finally {
    oauthLoading.value = false
  }
}

function handleCopy(text: string) {
  navigator.clipboard.writeText(text)
  message.success('已复制到剪贴板')
}

onMounted(() => {
  profileForm.nickname = user.value.nickname
  profileForm.email = user.value.email
  profileForm.avatar = user.value.avatar
  loadAccessInfo()
  loadOAuthBindings()
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NGrid
      x-gap="24"
      y-gap="24"
      cols="1 800:3"
    >
      <!-- 左侧：个人概览卡片 -->
      <NGi span="1">
        <div class="space-y-4">
          <NCard
            class="h-full shadow-sm"
            content-style="padding: 24px;"
            :bordered="false"
          >
            <div class="flex flex-col items-center text-center">
              <UserAvatar
                :size="88"
                class="mb-4 shadow-md"
              />
              <div class="text-xl font-bold text-gray-800 dark:text-gray-100">
                {{ user.nickname || '未设置昵称' }}
              </div>
              <div class="text-sm text-gray-400">
                @{{ user.username }}
              </div>

              <div class="mt-4 flex gap-2">
                <NTag
                  :type="user.id ? 'success' : 'warning'"
                  size="small"
                  round
                  :bordered="false"
                >
                  {{ user.id ? '账号已激活' : '未激活' }}
                </NTag>
                <NTag
                  v-for="role in roles.slice(0, 2)"
                  :key="role"
                  type="primary"
                  size="small"
                  round
                  :bordered="false"
                >
                  {{ role }}
                </NTag>
              </div>

              <NDivider class="my-6" />

              <div class="w-full">
                <NGrid
                  cols="2"
                  x-gap="12"
                  class="text-left"
                >
                  <NGi>
                    <NStatistic
                      label="注册天数"
                      tabular-nums
                    >
                      <template #suffix>
                        天
                      </template>
                      {{ Math.floor((Date.now() - new Date(user.createdAt).getTime()) / (1000 * 60 * 60 * 24)) }}
                    </NStatistic>
                  </NGi>
                  <NGi>
                    <NPopover trigger="hover" scrollable style="max-height: 300px;">
                      <template #trigger>
                        <div class="cursor-pointer transition-opacity hover:opacity-80">
                          <NStatistic label="当前角色">
                            {{ roles[0] || '访客' }}
                            <template #suffix>
                              <span v-if="roles.length > 1" class="text-xs text-gray-400">
                                (+{{ roles.length - 1 }})
                              </span>
                            </template>
                          </NStatistic>
                        </div>
                      </template>
                      <div class="w-64">
                        <div class="mb-2 text-xs font-medium text-gray-500">拥有权限 ({{ permissions.length }})</div>
                        <div class="flex flex-wrap gap-1">
                          <NTag
                            v-for="perm in permissions"
                            :key="perm"
                            size="small"
                            :bordered="false"
                            type="info"
                          >
                            {{ perm }}
                          </NTag>
                          <span v-if="permissions.length === 0" class="text-xs text-gray-400">无特殊权限</span>
                        </div>
                      </div>
                    </NPopover>
                  </NGi>
                </NGrid>
              </div>
            </div>
          </NCard>

          <NCard
            title="详细信息"
            size="small"
            :bordered="false"
            class="shadow-sm"
          >
            <NDescriptions
              :column="1"
              label-placement="left"
              label-style="width: 80px; color: #888;"
            >
              <NDescriptionsItem label="ID">
                <span
                  class="cursor-pointer font-mono text-xs text-gray-500 hover:text-primary"
                  @click="handleCopy(String(user.id))"
                >
                  {{ user.id }}
                </span>
              </NDescriptionsItem>
              <NDescriptionsItem label="注册时间">
                {{ user.createdAt ? new Date(user.createdAt).toLocaleDateString() : '-' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="最后更新">
                {{ user.updatedAt ? new Date(user.updatedAt).toLocaleDateString() : '-' }}
              </NDescriptionsItem>
            </NDescriptions>
          </NCard>
        </div>
      </NGi>

      <!-- 右侧：设置选项卡 -->
      <NGi span="2">
        <NCard
          class="h-full shadow-sm"
          :bordered="false"
        >
          <NTabs
            type="line"
            size="medium"
            animated
            justify-content="start"
            pane-class="p-4 md:p-8"
          >
            <NTabPane
              name="profile"
              tab="个人资料"
            >
              <NGrid
                :x-gap="40"
                item-responsive
                responsive="screen"
              >
                <!-- 左侧表单 -->
                <NGi span="24 m:14 l:15">
                  <NForm
                    ref="profileFormRef"
                    :model="profileForm"
                    :rules="profileRules"
                    label-placement="top"
                    require-mark-placement="right-hanging"
                  >
                    <NGrid
                      :x-gap="24"
                      :cols="1"
                    >
                      <NGi>
                        <NFormItem
                          label="昵称"
                          path="nickname"
                        >
                          <NInput
                            v-model:value="profileForm.nickname"
                            placeholder="如何称呼您？"
                            size="large"
                          />
                        </NFormItem>
                      </NGi>
                      <NGi>
                        <NFormItem
                          label="邮箱"
                          path="email"
                        >
                          <NInput
                            v-model:value="profileForm.email"
                            placeholder="联系邮箱"
                            size="large"
                          />
                        </NFormItem>
                      </NGi>
                      <NGi>
                        <NFormItem label="头像链接">
                          <NInput
                            v-model:value="profileForm.avatar"
                            type="textarea"
                            :rows="2"
                            placeholder="请输入有效的图片 URL"
                          />
                        </NFormItem>
                      </NGi>
                    </NGrid>

                    <div class="mt-6">
                      <NButton
                        type="primary"
                        size="large"
                        strong
                        @click="handleProfileSubmit"
                      >
                        保存个人信息
                      </NButton>
                    </div>
                  </NForm>
                </NGi>

                <!-- 右侧头像预览 -->
                <NGi span="24 m:10 l:9">
                  <div class="flex h-full flex-col items-center justify-start rounded-2xl bg-gray-50 py-8 dark:bg-white/5">
                    <div class="mb-6 text-sm font-medium text-gray-500">
                      头像预览
                    </div>
                    <UserAvatar
                      :src="profileForm.avatar"
                      :size="160"
                      class="mb-6 shadow-xl ring-4 ring-white dark:ring-gray-700"
                    />
                    <div class="text-xs text-gray-400">
                      支持 JPG, PNG, GIF 格式
                    </div>
                  </div>
                </NGi>
              </NGrid>
            </NTabPane>

            <NTabPane
              name="security"
              tab="账号安全"
            >
              <div class="mx-auto max-w-lg py-4">
                <div class="mb-8 text-center">
                  <h3 class="text-lg font-medium text-gray-800 dark:text-gray-100">
                    修改登录密码
                  </h3>
                  <p class="mt-1 text-sm text-gray-400">
                    建议定期更换密码以保护您的账户安全
                  </p>
                </div>

                <NForm
                  ref="passwordFormRef"
                  :model="passwordForm"
                  :rules="passwordRules"
                  label-placement="left"
                  :label-width="100"
                  require-mark-placement="left"
                  size="large"
                >
                  <NFormItem
                    label="当前密码"
                    path="oldPassword"
                  >
                    <NInput
                      v-model:value="passwordForm.oldPassword"
                      type="password"
                      show-password-on="click"
                      placeholder="验证当前密码"
                    />
                  </NFormItem>
                  <NDivider />
                  <NFormItem
                    label="新密码"
                    path="newPassword"
                  >
                    <NInput
                      v-model:value="passwordForm.newPassword"
                      type="password"
                      show-password-on="click"
                      placeholder="设置新密码"
                    />
                  </NFormItem>
                  <NFormItem
                    label="确认密码"
                    path="confirmPassword"
                  >
                    <NInput
                      v-model:value="passwordForm.confirmPassword"
                      type="password"
                      show-password-on="click"
                      placeholder="再次输入以确认"
                    />
                  </NFormItem>

                  <div class="mt-8 flex justify-center">
                    <NButton
                      type="primary"
                      size="large"
                      class="w-full"
                      @click="handlePasswordSubmit"
                    >
                      确认修改
                    </NButton>
                  </div>
                </NForm>
              </div>
            </NTabPane>

            <NTabPane
              name="binding"
              tab="第三方绑定"
            >
              <div class="mx-auto max-w-4xl py-2">
                <NGrid
                  v-if="oauthBindings.length > 0"
                  :x-gap="16"
                  :y-gap="16"
                  cols="1 m:2"
                >
                  <NGi
                    v-for="item in oauthBindings"
                    :key="item.providerKey + item.oauthID"
                  >
                    <div class="group relative flex items-start gap-4 rounded-xl border border-gray-100 bg-white p-4 transition-all hover:border-primary-100 hover:shadow-lg dark:border-gray-700 dark:bg-gray-800">
                      <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg bg-primary-50 text-xl font-bold uppercase text-primary-600 dark:bg-primary-900/30">
                        {{ item.providerKey.slice(0, 1) }}
                      </div>
                      <div class="flex-1">
                        <div class="flex items-center justify-between">
                          <h4 class="font-bold text-gray-800 dark:text-gray-100">
                            {{ item.providerName || (item.providerKey.charAt(0).toUpperCase() + item.providerKey.slice(1)) }}
                          </h4>
                          <NTag
                            size="small"
                            type="success"
                            bordered
                            class="scale-90"
                          >
                            已绑定
                          </NTag>
                        </div>
                        <div class="mt-2 text-xs text-gray-400">
                          ID: {{ item.oauthID }}
                        </div>
                        <div class="mt-1 text-xs text-gray-400">
                          绑定于 {{ new Date(item.boundAt).toLocaleDateString() }}
                        </div>
                      </div>
                    </div>
                  </NGi>
                </NGrid>
                
                <!-- Empty State -->
                <div
                  v-else-if="!oauthLoading"
                  class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-200 py-16 text-center dark:border-gray-700"
                >
                  <div class="mb-4 text-4xl text-gray-300">🔗</div>
                  <p class="text-gray-500">暂无绑定的第三方账号</p>
                </div>
              </div>
            </NTabPane>
          </NTabs>
        </NCard>
      </NGi>
    </NGrid>
  </ScrollContainer>
</template>
