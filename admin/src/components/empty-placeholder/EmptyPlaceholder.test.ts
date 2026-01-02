import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'

import { EmptyPlaceholder } from '../index'

describe('EmptyPlaceholder Component', () => {
  it('render component', () => {
    const wrapper = mount(EmptyPlaceholder, {
      props: {
        show: true,
      },
    })

    expect(wrapper.exists()).toBe(true)
    expect(wrapper.find('.n-empty').exists()).toBe(true)
  })
})
