import { GetLocalIp } from '@wa/services/SystemService'

/* Copy text */
export const copyText = async (text: string): Promise<boolean> => {
  try {
    // Prefer Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      return true
    }
    // Fallback to document.execCommand
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    const result = document.execCommand('copy')
    textArea.remove()
    return result
  } catch (error) {
    console.error('Copy failed:', error)
    return false
  }
}

/* Recursively get n layers of IPs; stop if 127.0.0.1 */
export const getLocalIpsDepth = async (depth: number, excludeIps: string[] = []): Promise<string[]> => {
  // Return current results if depth is not -1 and <= 0
  if (depth !== -1 && depth <= 0) {
    return excludeIps
  }

  const ip = await GetLocalIp(excludeIps)
  if (ip === '127.0.0.1') {
    excludeIps.push(ip)
    return excludeIps
  }
  return getLocalIpsDepth(depth === -1 ? -1 : depth - 1, [...excludeIps, ip])
}
