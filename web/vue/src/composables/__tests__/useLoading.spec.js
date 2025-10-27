import { describe, it, expect } from 'vitest'
import { useLoading } from '../useLoading'

describe('useLoading', () => {
  it('should initialize with false', () => {
    const { loading } = useLoading()
    expect(loading.value).toBe(false)
  })

  it('should set loading during async operation', async () => {
    const { loading, withLoading } = useLoading()
    
    const promise = withLoading(async () => {
      expect(loading.value).toBe(true)
      return 'result'
    })
    
    const result = await promise
    expect(loading.value).toBe(false)
    expect(result).toBe('result')
  })
})
