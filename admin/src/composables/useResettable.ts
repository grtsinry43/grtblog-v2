import { reactive, ref } from 'vue'

import type { Reactive, Ref } from 'vue'

function clone<T>(value: T): T {
  if (value === null || typeof value !== 'object') return value

  try {
    return structuredClone(value)
  } catch (error) {
    console.warn('structuredClone failed, using JSON fallback:', error)
    return JSON.parse(JSON.stringify(value))
  }
}

export function useResettableRef<S>(
  initialState: S,
): [Ref<S>, () => void, (initialState: S) => void] {
  let initialValue = clone(initialState)
  const state = ref<S>(clone(initialState)) as Ref<S>

  const reset = () => {
    state.value = clone(initialValue)
  }

  const updateInitialState = (initialState: S) => {
    initialValue = clone(initialState)
  }

  return [state, reset, updateInitialState]
}

export function useResettableReactive<S extends object>(
  initialState: S,
): [Reactive<S>, (value: Partial<S>) => void, () => void, (initialState: S) => void] {
  let initialValue = clone(initialState)
  const state = reactive<S>(clone(initialState)) as Reactive<S>

  const setState = (value: Partial<S>) => {
    Object.assign(state, value)
  }

  const reset = () => {
    Object.keys(state).forEach((key) => delete state[key as keyof typeof state])
    Object.assign(state, clone(initialValue))
  }

  const updateInitialState = (initialState: S) => {
    initialValue = clone(initialState)
  }

  return [state, setState, reset, updateInitialState]
}
