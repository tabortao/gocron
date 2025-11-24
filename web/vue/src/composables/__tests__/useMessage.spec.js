import { describe, it, expect, vi } from 'vitest'
import { useMessage } from '../useMessage'
import { ElMessage, ElMessageBox } from 'element-plus'

vi.mock('element-plus', () => ({
  ElMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  },
  ElMessageBox: {
    confirm: vi.fn()
  }
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key) => key
  })
}))

describe('useMessage', () => {
  it('should call ElMessage.success', () => {
    const { success } = useMessage()
    success('test message')
    expect(ElMessage.success).toHaveBeenCalledWith('test message')
  })

  it('should call ElMessage.error', () => {
    const { error } = useMessage()
    error('error message')
    expect(ElMessage.error).toHaveBeenCalledWith('error message')
  })

  it('should call ElMessageBox.confirm with default options', () => {
    const { confirm } = useMessage()
    confirm('Are you sure?')
    expect(ElMessageBox.confirm).toHaveBeenCalledWith(
      'Are you sure?',
      'common.tip',
      expect.objectContaining({
        confirmButtonText: 'common.confirm',
        cancelButtonText: 'common.cancel',
        type: 'warning',
        center: true
      })
    )
  })
})
