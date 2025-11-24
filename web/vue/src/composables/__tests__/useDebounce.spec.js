import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useDebounceFn } from '../useDebounce'

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should debounce function calls', () => {
    const mockFn = vi.fn()
    const debouncedFn = useDebounceFn(mockFn, 300)

    // 快速调用多次
    debouncedFn('call 1')
    debouncedFn('call 2')
    debouncedFn('call 3')

    // 还没到延迟时间，不应该被调用
    expect(mockFn).not.toHaveBeenCalled()

    // 快进时间
    vi.advanceTimersByTime(300)

    // 应该只被调用一次，使用最后一次的参数
    expect(mockFn).toHaveBeenCalledTimes(1)
    expect(mockFn).toHaveBeenCalledWith('call 3')
  })

  it('should use custom delay', () => {
    const mockFn = vi.fn()
    const debouncedFn = useDebounceFn(mockFn, 500)

    debouncedFn('test')

    vi.advanceTimersByTime(300)
    expect(mockFn).not.toHaveBeenCalled()

    vi.advanceTimersByTime(200)
    expect(mockFn).toHaveBeenCalledWith('test')
  })
})
