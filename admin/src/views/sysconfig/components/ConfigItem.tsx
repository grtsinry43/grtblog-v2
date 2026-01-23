import { NFormItem, NInput, NInputNumber, NSelect, NSwitch, NTag } from 'naive-ui'
import { defineComponent, type PropType } from 'vue'

import type { SysConfigItem } from '@/services/sysconfig'

export default defineComponent({
  name: 'ConfigItem',
  props: {
    item: { type: Object as PropType<SysConfigItem>, required: true },
    valueMap: { type: Object as PropType<Record<string, any>>, required: true },
    jsonBufferMap: { type: Object as PropType<Record<string, string>>, required: true },
    visible: { type: Function as PropType<(i: SysConfigItem) => boolean>, required: true },
  },
  setup(props) {
    return () => {
      const { item, valueMap, jsonBufferMap, visible } = props

      // 1. 可见性检查
      if (!visible(item)) return null

      // 2. 渲染 Label (Slot 内容)
      const renderLabel = () => (
        <div class='flex items-center gap-2'>
          <span>{item.label || item.key}</span>
          {item.isSensitive && (
            <NTag
              size='small'
              type='error'
              bordered={false}
              class='origin-left scale-90'
            >
              敏感
            </NTag>
          )}
           <NTag size="small" type="default" bordered={false} class="scale-90 opacity-50">{item.valueType}</NTag>
        </div>
      )

      // 3. 渲染控件
      let control = null
      switch (item.valueType) {
        case 'bool':
          control = <NSwitch v-model:value={valueMap[item.key]} />
          break
        case 'number':
          control = (
            <NInputNumber
              v-model:value={valueMap[item.key]}
              showButton={false}
              clearable
              placeholder='请输入数字'
            />
          )
          break
        case 'enum':
          const options =
            item.enumOptions?.map((opt) =>
              typeof opt === 'string' ? { label: opt, value: opt } : opt,
            ) || []
          control = (
            <NSelect
              v-model:value={valueMap[item.key]}
              options={options}
              clearable
            />
          )
          break
        case 'json':
          control = (
            <NInput
              v-model:value={jsonBufferMap[item.key]}
              type='textarea'
              autosize={{ minRows: 2, maxRows: 6 }}
              placeholder='请输入 JSON'
              style={{ fontFamily: 'monospace' }}
            />
          )
          break
        case 'string':
        default:
          const inputType = (item.meta as any)?.inputType === 'password' ? 'password' : 'text'
          control = (
            <NInput
              v-model:value={valueMap[item.key]}
              type={inputType}
              clearable
              showPasswordOn='click'
              placeholder={item.isSensitive ? '********** (留空不更新)' : ''}
            />
          )
          break
      }

      // 4. 组装 FormItem (使用 Slots)
      return (
        <NFormItem>
          {{
            label: renderLabel, // 这里的函数会被 Vue 渲染为 slot
            default: () => (
              <div class='w-full'>
                {control}
                {/* 5. 渲染描述文字 */}
                {item.description && (
                  <div class='mt-1.5 text-xs leading-relaxed text-gray-400'>{item.description}</div>
                )}
              </div>
            ),
          }}
        </NFormItem>
      )
    }
  },
})
