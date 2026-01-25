import { NCard, NButton, NInput, NForm, NFormItem, NAlert, NCode } from 'naive-ui'
import { defineComponent, ref } from 'vue'

import { ScrollContainer } from '@/components'
import { useMessage } from 'naive-ui'

import { requestFederationFriendlink } from '@/services/federation-admin'

export default defineComponent({
  name: 'FederationInstances',
  setup() {
    const message = useMessage()
    const targetUrl = ref('')
    const rssUrl = ref('')
    const note = ref('')
    const loading = ref(false)
    const error = ref('')
    const result = ref<any>(null)

    const handleSubmit = async () => {
      const url = targetUrl.value.trim()
      if (!url) {
        message.warning('请输入远端实例地址')
        return
      }
      loading.value = true
      error.value = ''
      result.value = null
      try {
        result.value = await requestFederationFriendlink({
          target_url: url,
          message: note.value.trim() || undefined,
          rss_url: rssUrl.value.trim() || undefined,
        })
        message.success('已提交申请')
      } catch (err: any) {
        error.value = err?.message || '请求失败'
      } finally {
        loading.value = false
      }
    }

    return () => (
      <ScrollContainer wrapperClass='p-4'>
        <NCard>
          <div class='space-y-6'>
            <div>
              <div class='text-base font-semibold'>联邦实例</div>
              <div class='text-xs text-neutral-500'>对外发起友链申请</div>
            </div>

            <NForm labelPlacement='left' labelWidth={120}>
              <NFormItem label='远端地址'>
                <NInput
                  v-model:value={targetUrl.value}
                  placeholder='https://example.com'
                />
              </NFormItem>
              <NFormItem label='RSS 地址'>
                <NInput
                  v-model:value={rssUrl.value}
                  placeholder='https://example.com/rss'
                />
              </NFormItem>
              <NFormItem label='备注'>
                <NInput
                  v-model:value={note.value}
                  placeholder='可选'
                />
              </NFormItem>
            </NForm>

            <div class='flex justify-end'>
              <NButton
                type='primary'
                loading={loading.value}
                onClick={handleSubmit}
              >
                发起申请
              </NButton>
            </div>

            {error.value && <NAlert type='error' title='请求失败' content={error.value} />}

            {result.value && (
              <div class='space-y-2'>
                <div class='text-sm font-medium text-neutral-600'>响应</div>
                <NCode
                  code={JSON.stringify(result.value, null, 2)}
                  language='json'
                  wordWrap
                />
              </div>
            )}
          </div>
        </NCard>
      </ScrollContainer>
    )
  },
})
