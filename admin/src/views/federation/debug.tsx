import { NCard, NButton, NInput, NForm, NFormItem, NAlert, NCode } from 'naive-ui'
import { defineComponent, ref } from 'vue'

import { ScrollContainer } from '@/components'
import { useMessage } from 'naive-ui'

import { checkFederationRemote } from '@/services/federation-admin'

export default defineComponent({
  name: 'FederationDebug',
  setup() {
    const message = useMessage()
    const targetUrl = ref('')
    const loading = ref(false)
    const error = ref('')
    const result = ref<any>(null)

    const handleCheck = async () => {
      const url = targetUrl.value.trim()
      if (!url) {
        message.warning('请输入远端实例地址')
        return
      }
      loading.value = true
      error.value = ''
      result.value = null
      try {
        result.value = await checkFederationRemote(url)
        message.success('已获取远端信息')
      } catch (err: any) {
        error.value = err?.message || '请求失败'
      } finally {
        loading.value = false
      }
    }

    const renderBlock = (title: string, data: any) => (
      <div class='space-y-2'>
        <div class='text-sm font-medium text-neutral-600'>{title}</div>
        {data ? (
          <NCode
            code={JSON.stringify(data, null, 2)}
            language='json'
            wordWrap
          />
        ) : (
          <div class='text-xs text-neutral-400'>暂无数据</div>
        )}
      </div>
    )

    return () => (
      <ScrollContainer wrapperClass='p-4'>
        <NCard>
          <div class='space-y-6'>
            <div>
              <div class='text-base font-semibold'>联邦调试</div>
              <div class='text-xs text-neutral-500'>拉取远端 manifest / public-key / endpoints</div>
            </div>

            <NForm labelPlacement='left' labelWidth={120}>
              <NFormItem label='远端地址'>
                <div class='flex w-full gap-2'>
                  <NInput
                    v-model:value={targetUrl.value}
                    placeholder='https://example.com'
                  />
                  <NButton
                    type='primary'
                    loading={loading.value}
                    onClick={handleCheck}
                  >
                    检查
                  </NButton>
                </div>
              </NFormItem>
            </NForm>

            {error.value && <NAlert type='error' title='请求失败' content={error.value} />}

            {result.value && (
              <div class='space-y-6'>
                {renderBlock('Manifest', result.value.manifest)}
                {renderBlock('Public Key', result.value.public_key)}
                {renderBlock('Endpoints', result.value.endpoints)}
              </div>
            )}
          </div>
        </NCard>
      </ScrollContainer>
    )
  },
})
