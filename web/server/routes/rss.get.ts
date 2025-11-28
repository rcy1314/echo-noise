import { defineEventHandler, setHeader, getRequestURL } from 'h3'

export default defineEventHandler(async (event) => {
  try {
    const config = useRuntimeConfig()
    const baseApi = (config.public as any)?.baseApi || '/api'
    const url = getRequestURL(event)
    const origin = `${url.protocol}//${url.host}`

    const settings: any = await $fetch(`${baseApi}/frontend/config`, {
      credentials: 'include'
    }).catch(() => null)
    const frontend = settings?.data?.frontendSettings || {}
    const title = (frontend.rssTitle || frontend.siteTitle || '说说笔记').trim()
    const description = (frontend.rssDescription || frontend.description || '').trim()
    const icon = (frontend.rssFaviconURL || '/favicon.ico').trim()

    const msgRes: any = await $fetch(`${baseApi}/messages`, {
      query: { page: 1, pageSize: 20 },
      credentials: 'include'
    }).catch(() => null)
    const list = (msgRes?.data?.items || msgRes?.data?.messages || []).filter(Boolean)

    const imageUrl = icon.startsWith('http') ? icon : `${origin}${icon}`
    const now = new Date().toUTCString()
    const toItem = (m: any) => {
      const link = `${origin}/#/messages/${m.id}`
      const pubDate = m.created_at ? new Date(m.created_at).toUTCString() : now
      const title = (m.username || '用户') + ' 的消息'
      const content = (m.content || '').replace(/<!\-\-[\s\S]*?\-\->/g, '').trim()
      return `\n  <item>\n    <title><![CDATA[${title}]]></title>\n    <link>${link}</link>\n    <guid>${link}</guid>\n    <pubDate>${pubDate}</pubDate>\n    <description><![CDATA[${content}]]></description>\n  </item>`
    }

    const itemsXml = list.map(toItem).join('')
    const xml = `<?xml version="1.0" encoding="UTF-8"?>\n<rss version="2.0">\n<channel>\n<title><![CDATA[${title}]]></title>\n<link>${origin}/</link>\n<description><![CDATA[${description}]]></description>\n<image>\n  <url>${imageUrl}</url>\n  <title><![CDATA[${title}]]></title>\n  <link>${origin}/</link>\n</image>\n<lastBuildDate>${now}</lastBuildDate>${itemsXml}\n</channel>\n</rss>`

    setHeader(event, 'Content-Type', 'application/rss+xml; charset=utf-8')
    return xml
  } catch (e) {
    setHeader(event, 'Content-Type', 'text/plain; charset=utf-8')
    return 'RSS generate error'
  }
})

