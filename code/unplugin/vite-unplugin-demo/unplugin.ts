import { createUnplugin } from 'unplugin'

export const LogPrefixUnplugin = createUnplugin((options?: { prefix?: string }) => {
  const prefix = options?.prefix ?? '[DEBUG]:'
  return {
    name: 'log-prefix-unplugin',
    enforce: 'pre',
    transform(code, id) {
      if (!id.endsWith('.ts')) return

      return code.replace(/console\.log\(/g, `console.log('${prefix}', `)
    }
  }
})
