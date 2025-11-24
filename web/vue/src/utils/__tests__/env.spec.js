import { describe, it, expect, vi } from 'vitest'
import { env, devLog } from '../env'

describe('env utilities', () => {
  it('should have correct env properties', () => {
    expect(env).toHaveProperty('isDev')
    expect(env).toHaveProperty('isProd')
    expect(env).toHaveProperty('apiBaseUrl')
  })

  it('should have default apiBaseUrl', () => {
    expect(env.apiBaseUrl).toBe('/api')
  })

  it('devLog should only log in dev mode', () => {
    const consoleSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
    
    devLog('test message')
    
    if (env.isDev) {
      expect(consoleSpy).toHaveBeenCalledWith('[Dev]', 'test message')
    } else {
      expect(consoleSpy).not.toHaveBeenCalled()
    }
    
    consoleSpy.mockRestore()
  })
})
