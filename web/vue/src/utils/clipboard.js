export async function copyText(text) {
  if (typeof text !== 'string') {
    text = String(text ?? '')
  }

  if (navigator?.clipboard?.writeText && window.isSecureContext) {
    await navigator.clipboard.writeText(text)
    return
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', '')
  textarea.style.position = 'fixed'
  textarea.style.top = '0'
  textarea.style.left = '0'
  textarea.style.opacity = '0'
  textarea.style.pointerEvents = 'none'

  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()

  const ok = document.execCommand('copy')
  document.body.removeChild(textarea)

  if (!ok) {
    throw new Error('copy failed')
  }
}
