import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js'
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js'
import http from 'http'
import { EventEmitter } from 'events'
import { parse } from 'url'
import { z } from 'zod'

const host = String(process.env.NOTE_HOST || 'http://localhost:1314').trim().replace(/^`+|`+$/g, '')
const token = process.env.NOTE_TOKEN || ''
let session = process.env.NOTE_SESSION || ''

const s = new McpServer({ name: 'ech0-noise-mcp', version: '0.1.0' })
const bus = new EventEmitter()

const authHeaders = () => {
  const h = { 'Content-Type': 'application/json' }
  if (token) h['Authorization'] = `Bearer ${token}`
  if (session) h['Cookie'] = session
  return h
}

async function searchTool(args) {
  try {
    const q = (args.keyword || '').trim()
    const page = Number(args.page || 1)
    const pageSize = Number(args.pageSize || 10)
    const fmt = String(args.format || '').toLowerCase()
    if (!q) {
      const url = `${host}/api/messages/page?page=${page}&pageSize=${pageSize}`
      const r = await fetch(url)
      if (!r.ok) throw new Error(`HTTP ${r.status}`)
      const j = await r.json()
      const arr = (j && j.data && Array.isArray(j.data.items)) ? j.data.items : (Array.isArray(j.items) ? j.items : [])
      const lines = Array.isArray(arr) ? arr.map((it) => {
        const content = String(it.content || '')
        const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
        const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(content)
        const imgMd = !hasMdImg && Array.isArray(imgs) && imgs.length ? imgs.map((u) => `![image](${u})`).join('\n') : ''
        const body = imgMd ? `${imgMd}\n${content}` : content
        const t = formatTimeShort(it.created_at || it.createdAt)
        return `[${it.id}] ${it.username || ''} ${t}\n\n${body}`
      }) : []
      const header = '以下是最新内容：'
      const imgCount = Array.isArray(arr) ? arr.filter((it) => {
        const c = String(it.content || '')
        const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(c)
        const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
        return hasMdImg || (Array.isArray(imgs) && imgs.length)
      }).length : 0
      const sample = Array.isArray(arr) ? arr.slice(0,3).map((it) => String(it.content || '').trim()).filter(Boolean) : []
      const summary = [`共${Array.isArray(arr) ? arr.length : 0}条`, imgCount ? `包含图片${imgCount}条` : ''].concat(sample.map((s) => `- ${s.slice(0,60)}`)).filter(Boolean).join('\n')
      const text = lines.length ? `${header}\n\n${lines.join('\n\n')}\n\n摘要：\n${summary}` : '无匹配结果'
      if (fmt === 'json') return { content: [{ type: 'text', text: JSON.stringify(j) }] }
      return { content: [{ type: 'text', text }] }
    }
    if (q.startsWith('#')) {
      const tag = encodeURIComponent(q.slice(1))
      const url = `${host}/api/messages/tags/${tag}?page=${page}&pageSize=${pageSize}`
      const r = await fetch(url)
      if (!r.ok) throw new Error(`HTTP ${r.status}`)
      const j = await r.json()
      const arr = (j && j.data && Array.isArray(j.data.items)) ? j.data.items : (Array.isArray(j.items) ? j.items : [])
      const lines = Array.isArray(arr) ? arr.map((it) => {
        const content = String(it.content || '')
        const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
        const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(content)
        const imgMd = !hasMdImg && Array.isArray(imgs) && imgs.length ? imgs.map((u) => `![image](${u})`).join('\n') : ''
        const body = imgMd ? `${imgMd}\n${content}` : content
        const t = formatTimeShort(it.created_at || it.createdAt)
        return `[${it.id}] ${it.username || ''} ${t}\n\n${body}`
      }) : []
      const header = `关于“${q}”的相关内容：`
      const imgCount = Array.isArray(arr) ? arr.filter((it) => {
        const c = String(it.content || '')
        const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(c)
        const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
        return hasMdImg || (Array.isArray(imgs) && imgs.length)
      }).length : 0
      const sample = Array.isArray(arr) ? arr.slice(0,3).map((it) => String(it.content || '').trim()).filter(Boolean) : []
      const summary = [`共${Array.isArray(arr) ? arr.length : 0}条`, imgCount ? `包含图片${imgCount}条` : ''].concat(sample.map((s) => `- ${s.slice(0,60)}`)).filter(Boolean).join('\n')
      const text = lines.length ? `${header}\n\n${lines.join('\n\n')}\n\n摘要：\n${summary}` : '无匹配结果'
      if (fmt === 'json') return { content: [{ type: 'text', text: JSON.stringify(j) }] }
      return { content: [{ type: 'text', text }] }
    }
    const url = `${host}/api/messages/search?keyword=${encodeURIComponent(q)}&page=${page}&pageSize=${pageSize}`
    const r = await fetch(url)
    if (!r.ok) throw new Error(`HTTP ${r.status}`)
    const j = await r.json()
    const arr = (j && j.data && Array.isArray(j.data.items)) ? j.data.items : (Array.isArray(j.items) ? j.items : [])
    const lines = Array.isArray(arr) ? arr.map((it) => {
      const content = String(it.content || '')
      const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
      const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(content)
      const imgMd = !hasMdImg && Array.isArray(imgs) && imgs.length ? imgs.map((u) => `![image](${u})`).join('\n') : ''
      const body = imgMd ? `${imgMd}\n${content}` : content
      const t = formatTimeShort(it.created_at || it.createdAt)
      return `[${it.id}] ${it.username || ''} ${t}\n\n${body}`
    }) : []
    const header = ((q && q.includes('欢迎')) || (Array.isArray(arr) && arr.some((it) => String(it.content || '').trim().startsWith('#欢迎')))) ? '欢迎内容如下：' : `关于“${q}”的相关内容：`
    const imgCount = Array.isArray(arr) ? arr.filter((it) => {
      const c = String(it.content || '')
      const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(c)
      const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
      return hasMdImg || (Array.isArray(imgs) && imgs.length)
    }).length : 0
    const sample = Array.isArray(arr) ? arr.slice(0,3).map((it) => String(it.content || '').trim()).filter(Boolean) : []
    const summary = [`共${Array.isArray(arr) ? arr.length : 0}条`, imgCount ? `包含图片${imgCount}条` : ''].concat(sample.map((s) => `- ${s.slice(0,60)}`)).filter(Boolean).join('\n')
    const text = lines.length ? `${header}\n\n${lines.join('\n\n')}\n\n摘要：\n${summary}` : '无匹配结果'
    if (fmt === 'json') return { content: [{ type: 'text', text: JSON.stringify(j) }] }
    return { content: [{ type: 'text', text }] }
  } catch (e) {
    const msg = String(e && e.message || e)
    return { content: [{ type: 'text', text: `error=${msg}` }] }
  }
}

async function pageTool(args) {
  try {
    const page = Number(args.page || args.page_number || 1)
    const pageSize = Number(args.pageSize || args.page_size || 10)
    const fmt = String(args.format || '').toLowerCase()
    const url = `${host}/api/messages/page?page=${page}&pageSize=${pageSize}`
    const r = await fetch(url)
    if (!r.ok) throw new Error(`HTTP ${r.status}`)
    const j = await r.json()
    const arr = (j && j.data && Array.isArray(j.data.items)) ? j.data.items : (Array.isArray(j.items) ? j.items : [])
    const lines = Array.isArray(arr) ? arr.map((it) => {
      const content = String(it.content || '')
      const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
      const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(content)
      const imgMd = !hasMdImg && Array.isArray(imgs) && imgs.length ? imgs.map((u) => `![image](${u})`).join('\n') : ''
      const body = imgMd ? `${imgMd}\n${content}` : content
      const t = formatTimeShort(it.created_at || it.createdAt)
      return `${body}\n\n[${it.id}] ${it.username || ''} ${t}`
    }) : []
    const header = (Array.isArray(arr) && arr.some((it) => String(it.content || '').trim().startsWith('#欢迎'))) ? '欢迎内容如下：' : '以下是分页内容：'
    const imgCount = Array.isArray(arr) ? arr.filter((it) => {
      const c = String(it.content || '')
      const hasMdImg = /!\[[^\]]*\]\([^\)]+\)/.test(c)
      const imgs = Array.isArray(it.images) ? it.images : (it.image ? [it.image] : (it.imageURL ? [it.imageURL] : []))
      return hasMdImg || (Array.isArray(imgs) && imgs.length)
    }).length : 0
    const sample = Array.isArray(arr) ? arr.slice(0,3).map((it) => String(it.content || '').trim()).filter(Boolean) : []
    const summary = [`共${Array.isArray(arr) ? arr.length : 0}条`, imgCount ? `包含图片${imgCount}条` : ''].concat(sample.map((s) => `- ${s.slice(0,60)}`)).filter(Boolean).join('\n')
    const text = lines.length ? `${header}\n\n${lines.join('\n\n')}\n\n摘要：\n${summary}` : '无匹配结果'
    if (fmt === 'json') return { content: [{ type: 'text', text: JSON.stringify(j) }] }
    return { content: [{ type: 'text', text }] }
  } catch (e) {
    const msg = String(e && e.message || e)
    return { content: [{ type: 'text', text: `error=${msg}` }] }
  }
}

async function getTool(args) {
  const id = String(args.id)
  const r = await fetch(`${host}/api/messages/${id}`)
  const j = await r.json()
  return { content: [{ type: 'text', text: JSON.stringify(j) }] }
}

async function publishTool(args) {
  const priv = Boolean(args.private || false)
  const typeRaw = String(args.type || 'text').toLowerCase()
  const type = ['text','markdown','image','multipart','md'].includes(typeRaw) ? (typeRaw === 'md' ? 'markdown' : typeRaw) : 'text'
  const contentRaw = typeof args.content === 'string' ? args.content : (args.content ? JSON.stringify(args.content) : '')
  const content = String(contentRaw || '').trim()
  const images = Array.isArray(args.images) ? args.images.map(String) : []
  const image = args.image ? String(args.image) : (args.imageURL ? String(args.imageURL) : (images[0] || ''))
  const body = { type, private: priv }
  if (content) body.content = content
  if (type === 'image' || type === 'multipart') {
    if (images.length) body.images = images
    else if (image) body.image = image
  }
  const endpoint = token ? `${host}/api/token/messages` : `${host}/api/messages`
  const r = await fetch(endpoint, { method: 'POST', headers: authHeaders(), body: JSON.stringify(body) })
  const j = await r.json().catch(async () => ({ code: r.ok ? 1 : 0 }))
  return { content: [{ type: 'text', text: JSON.stringify(j) }] }
}

async function deleteTool(args) {
  const idRaw = args && args.id
  const id = String(typeof idRaw === 'number' ? idRaw : (idRaw || '')).trim()
  if (!id || id === 'undefined' || id === 'null' || !/^[0-9]+$/.test(id)) {
    const kw = String(args && (args.keyword || args.content) || '').trim()
    if (!kw) {
      const text = '使用 MCP 删除失败：参数错误。请提供有效的 id（数字字符串），或传入 keyword/content'
      return { content: [{ type: 'text', text }] }
    }
    try {
      const url = `${host}/api/messages/search?keyword=${encodeURIComponent(kw)}&page=1&pageSize=1`
      const r = await fetch(url)
      const j = await r.json().catch(() => ({}))
      const arr = (j && j.data && Array.isArray(j.data.items)) ? j.data.items : (Array.isArray(j.items) ? j.items : [])
      const first = Array.isArray(arr) && arr[0]
      const fid = first && first.id
      if (!fid) {
        const text = '使用 MCP 删除失败：未找到匹配信息，无法删除'
        return { content: [{ type: 'text', text }] }
      }
      args.id = String(fid)
    } catch (e) {
      const text = `使用 MCP 删除失败：搜索错误 ${String(e && e.message || e)}`
      return { content: [{ type: 'text', text }] }
    }
  }
  try {
    const endpointToken = `${host}/api/token/messages/${id}`
    const endpointPublic = `${host}/api/messages/${id}`
    if (token) {
      const rtok = await fetch(endpointToken, { method: 'DELETE', headers: authHeaders() }).catch(() => null)
      if (rtok && rtok.ok) {
        const jt = await rtok.json().catch(() => ({}))
        const code = jt && typeof jt.code !== 'undefined' ? Number(jt.code) : undefined
        const msg = String(jt && jt.msg || '')
        if (code === 1 || msg.includes('删除成功')) {
          const text = `使用 MCP 删除成功：id=${id}`
          return { content: [{ type: 'text', text }] }
        }
      }
      if (rtok && !rtok.ok && (rtok.status === 401 || rtok.status === 403)) {
        const text = authMessage(rtok.status)
        return { content: [{ type: 'text', text }] }
      }
      if (rtok && !rtok.ok && rtok.status === 404) {
        const text = authMessage(404)
        return { content: [{ type: 'text', text }] }
      }
    }
    const checkAfterToken = await fetch(endpointPublic, { headers: authHeaders() }).catch(() => ({ ok: false, status: 500 }))
    if (checkAfterToken && (checkAfterToken.status === 404)) {
      const text = `使用 MCP 删除成功：id=${id}`
      return { content: [{ type: 'text', text }] }
    }
    if (checkAfterToken && checkAfterToken.ok) {
      const left = await checkAfterToken.json().catch(() => ({}))
      const code = left && typeof left.code !== 'undefined' ? Number(left.code) : undefined
      const msg = String(left && left.msg || '')
      if (code === 0 && (msg.includes('不存在') || msg.includes('获取信息失败'))) {
        const text = `使用 MCP 删除成功：id=${id}`
        return { content: [{ type: 'text', text }] }
      }
    }
    if (session) {
      const r = await fetch(endpointPublic, { method: 'DELETE', headers: authHeaders() })
      if (!r.ok && (r.status === 401 || r.status === 403)) {
        const text = authMessage(r.status)
        return { content: [{ type: 'text', text }] }
      }
      if (!r.ok && r.status === 404) {
        const text = `使用 MCP 删除成功：id=${id}`
        return { content: [{ type: 'text', text }] }
      }
      const checkPublic = await fetch(endpointPublic, { headers: authHeaders() }).catch(() => ({ ok: false, status: 500 }))
      if (checkPublic && (checkPublic.status === 404)) {
        const text = `使用 MCP 删除成功：id=${id}`
        return { content: [{ type: 'text', text }] }
      }
      if (checkPublic && checkPublic.ok) {
        const left = await checkPublic.json().catch(() => ({}))
        const code = left && typeof left.code !== 'undefined' ? Number(left.code) : undefined
        const msg = String(left && left.msg || '')
        if (code === 0 && (msg.includes('不存在') || msg.includes('获取信息失败'))) {
          const text = `使用 MCP 删除成功：id=${id}`
          return { content: [{ type: 'text', text }] }
        }
      }
    }
    try {
      const body = { content: '[已删除]' }
      const r3 = await fetch(endpointPublic, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(body) })
      const j3 = await r3.json().catch(() => ({}))
      const p3 = await fetch(endpointPublic, { headers: authHeaders() }).catch(() => ({ ok: false, status: 500 }))
      if (p3 && p3.ok) {
        const left3 = await p3.json().catch(() => ({}))
        const ok3 = left3 && left3.data && typeof left3.data.content === 'string' && left3.data.content.includes('已删除')
        if (ok3) {
          const text = `使用 MCP 删除失败，已改为逻辑删除：id=${id}`
          return { content: [{ type: 'text', text }] }
        }
      }
    } catch {}
    const text = `使用 MCP 删除失败：id=${id} 仍存在。`
    return { content: [{ type: 'text', text }] }
  } catch (e) {
    return { content: [{ type: 'text', text: JSON.stringify({ error: String(e && e.message || e) }) }] }
  }
}

async function updateTool(args) {
  const id = String(args.id)
  const content = String(args.content || '')
  const body = { content }
  if (!session) {
    const text = authMessage('need_session')
    return { content: [{ type: 'text', text }] }
  }
  try {
    const endpoint = token ? `${host}/api/token/messages/${id}` : `${host}/api/messages/${id}`
    const r = await fetch(endpoint, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(body) })
    if (!r.ok && (r.status === 401 || r.status === 403)) {
      const text = authMessage(r.status)
      return { content: [{ type: 'text', text }] }
    }
    if (!r.ok && r.status === 404 && token) {
      const text = authMessage(404)
      return { content: [{ type: 'text', text }] }
    }
    const j = await r.json()
    return { content: [{ type: 'text', text: JSON.stringify(j) }] }
  } catch (e) {
    return { content: [{ type: 'text', text: JSON.stringify({ error: String(e && e.message || e) }) }] }
  }
}

async function pinTool(args) {
  const id = String(args.id)
  const pinned = Boolean(args.pinned)
  const body = { pinned }
  if (!session) {
    const text = authMessage('need_session')
    return { content: [{ type: 'text', text }] }
  }
  try {
    const endpoint = token ? `${host}/api/token/messages/${id}/pin` : `${host}/api/messages/${id}/pin`
    const r = await fetch(endpoint, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(body) })
    if (!r.ok && (r.status === 401 || r.status === 403)) {
      const text = authMessage(r.status)
      return { content: [{ type: 'text', text }] }
    }
    if (!r.ok && r.status === 404 && token) {
      const text = authMessage(404)
      return { content: [{ type: 'text', text }] }
    }
    const j = await r.json()
    return { content: [{ type: 'text', text: JSON.stringify(j) }] }
  } catch (e) {
    return { content: [{ type: 'text', text: JSON.stringify({ error: String(e && e.message || e) }) }] }
  }
}

async function settingsTool(args) {
  const body = {}
  if (typeof args.allowRegistration !== 'undefined') body.allowRegistration = Boolean(args.allowRegistration)
  if (typeof args.frontendSettings !== 'undefined') body.frontendSettings = args.frontendSettings
  if (!session) {
    const text = authMessage('need_session')
    return { content: [{ type: 'text', text }] }
  }
  try {
    const endpoint = token ? `${host}/api/token/settings` : `${host}/api/settings`
    const r = await fetch(endpoint, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(body) })
    if (!r.ok && (r.status === 401 || r.status === 403)) {
      const text = authMessage(r.status)
      return { content: [{ type: 'text', text }] }
    }
    if (!r.ok && r.status === 404 && token) {
      const text = authMessage(404)
      return { content: [{ type: 'text', text }] }
    }
    const j = await r.json()
    return { content: [{ type: 'text', text: JSON.stringify(j) }] }
  } catch (e) {
    return { content: [{ type: 'text', text: JSON.stringify({ error: String(e && e.message || e) }) }] }
  }
}

async function statusTool() {
  const r = await fetch(`${host}/api/status`)
  const j = await r.json()
  return { content: [{ type: 'text', text: JSON.stringify(j) }] }
}

async function calendarTool() {
  const r = await fetch(`${host}/api/messages/calendar`)
  const j = await r.json()
  return { content: [{ type: 'text', text: JSON.stringify(j) }] }
}

async function configTool() {
  const r = await fetch(`${host}/api/frontend/config`)
  const j = await r.json()
  return { content: [{ type: 'text', text: JSON.stringify(j) }] }
}

async function loginTool(args) {
  const username = String(args.username || '')
  const password = String(args.password || '')
  const r = await fetch(`${host}/api/login`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ username, password }) })
  const j = await r.json()
  const sc = r.headers.get('set-cookie') || ''
  let ck = ''
  if (sc) {
    ck = String(sc).split(',')[0].split(';')[0].trim()
    session = ck
  }
  return { content: [{ type: 'text', text: JSON.stringify({ cookie: ck || sc, response: j }) }] }
}

async function tokenTool() {
  const r = await fetch(`${host}/api/user/token/regenerate`, { method: 'POST', headers: authHeaders() })
  const j = await r.json()
  return { content: [{ type: 'text', text: JSON.stringify(j) }] }
}

async function rssTool() {
  const r = await fetch(`${host}/rss`)
  const t = await r.text()
  return { content: [{ type: 'text', text: t }] }
}

function wrap(name, fn) {
  return async (args) => {
    bus.emit('tool_start', { name, args })
    const res = await fn(args)
    bus.emit('tool_end', { name, res })
    return res
  }
}
function authMessage(status) {
  const hasToken = !!token
  const hasSession = !!session
  if (!hasSession) {
    if (hasToken) {
      if (status === 'need_session') return '需要登录：当前已设置令牌，但更新/置顶/设置仅支持会话；请先 登录。发布和删除可继续使用令牌。'
      return '需要登录：请先调用 登录 工具。更新/置顶/设置仅支持会话；发布和删除可使用令牌。'
    }
    return '需要登录或令牌：请先 登录 或设置 NOTE_TOKEN。更新/置顶/设置仅支持会话；发布和删除可使用令牌。'
  }
  if (status === 401 || status === 403) return '会话无效或已过期：请重新 登录。'
  if (hasToken && status === 404) return '后端未启用令牌路由：请更新后端或使用 登录 获取会话后操作。'
  if (hasToken && (status === 'token_invalid' || status === 'token_expired')) return '令牌无效或已过期：请在后台重新生成令牌，或使用会话操作。'
  return '权限不足或资源不存在：请确认权限后重试。'
}

function formatTime(v) {
  const d = new Date(String(v || ''))
  if (isNaN(d.getTime())) {
    const s = String(v || '').replace('T', ' ').replace('Z', '')
    return s.split('.')[0]
  }
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

function formatTimeShort(v) {
  const d = new Date(String(v || ''))
  if (!isNaN(d.getTime())) {
    const p = (n) => String(n).padStart(2, '0')
    return `${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
  }
  const s = String(v || '').replace('T', ' ').replace('Z', '')
  const m = s.match(/(\d{4})-(\d{2})-(\d{2})\s+(\d{2}):(\d{2})/)
  if (m) return `${m[2]}-${m[3]} ${m[4]}:${m[5]}`
  return s.slice(5, 16)
}

const searchSchema = z.object({
  keyword: z.string().optional(),
  query: z.string().optional(),
  page: z.number().optional(),
  pageSize: z.number().optional(),
  format: z.string().optional()
})

const pageSchema = z.object({ page: z.number().optional(), pageSize: z.number().optional(), format: z.string().optional() })
const deleteSchema = z.object({ id: z.string().optional(), keyword: z.string().optional(), content: z.string().optional() })
const idSchema = z.object({ id: z.string() })
const anyJson = z.any()
const publishSchema = z.object({
  content: anyJson.optional(),
  type: z.string().optional(),
  private: z.boolean().optional(),
  image: z.string().optional(),
  images: z.array(z.string()).optional(),
  imageURL: z.string().optional()
})
const updateSchema = z.object({ id: z.string(), content: z.string() })
const loginSchema = z.object({ username: z.string(), password: z.string() })
const pinSchema = z.object({ id: z.string(), pinned: z.boolean() })
const settingsSchema = z.object({ allowRegistration: z.boolean().optional(), frontendSettings: z.any().optional() })

s.registerTool('search', { description: '搜索消息：默认输出为人类可读文本；若需原始 JSON 传入 format=json。触发查询参数时必须调用本工具，只输出解析后的中文/多语言列表，不展示工具调用日志或未渲染标签。', inputSchema: searchSchema }, wrap('search', (args) => searchTool({ ...args, keyword: args.keyword || args.query || '' })))
s.registerTool('页面', { description: '分页列表：触发 page/pageSize 参数时必须调用本工具', inputSchema: pageSchema }, wrap('页面', pageTool))
s.registerTool('消息', { description: '获取消息详情：传入 id 时必须调用本工具', inputSchema: idSchema }, wrap('消息', getTool))
s.registerTool('publish', { description: '发布内容：触发发布参数时必须调用本工具；支持令牌或会话', inputSchema: publishSchema }, wrap('publish', publishTool))
s.registerTool('delete', { description: '删除消息：需要认证；支持 id 或 keyword/content', inputSchema: deleteSchema }, wrap('delete', deleteTool))
s.registerTool('搜索', { description: '搜索消息：默认输出为人类可读文本；若需原始 JSON 传入 format=json。触发查询参数时必须调用本工具，只输出解析后的中文/多语言列表，不展示工具调用日志或未渲染标签。', inputSchema: searchSchema }, wrap('搜索', (args) => searchTool({ ...args, keyword: args.keyword || args.query || '' })))
s.registerTool('发布', { description: '发布内容：触发发布参数时必须调用本工具；支持令牌或会话', inputSchema: publishSchema }, wrap('发布', publishTool))
s.registerTool('笔记', { description: '发布内容：触发发布参数时必须调用本工具；支持令牌或会话', inputSchema: publishSchema }, wrap('笔记', publishTool))
s.registerTool('说说', { description: '发布内容：触发发布参数时必须调用本工具；支持令牌或会话', inputSchema: publishSchema }, wrap('说说', publishTool))
s.registerTool('说说笔记', { description: '发布内容：触发发布参数时必须调用本工具；支持令牌或会话', inputSchema: publishSchema }, wrap('说说笔记', publishTool))
s.registerTool('删除', { description: '删除消息：需要认证；支持 id 或 keyword/content', inputSchema: deleteSchema }, wrap('删除', deleteTool))
s.registerTool('更新', { description: '更新消息：需要认证；传入 id/content 时必须调用本工具', inputSchema: updateSchema }, wrap('更新', updateTool))
s.registerTool('置顶消息', { description: '置顶/取消置顶：需要认证；传入 id/pinned 时必须调用本工具', inputSchema: pinSchema }, wrap('置顶消息', pinTool))
s.registerTool('设置', { description: '系统设置：需要认证；触发设置参数时必须调用本工具', inputSchema: settingsSchema }, wrap('设置', settingsTool))
s.registerTool('状态', { description: '系统状态：无入参；需要获取状态时调用本工具', inputSchema: z.object({}) }, statusTool)
s.registerTool('日历', { description: '发布日历：无入参；需要获取日历时调用本工具', inputSchema: z.object({}) }, calendarTool)
s.registerTool('配置', { description: '前端配置：无入参；需要获取配置时调用本工具', inputSchema: z.object({}) }, configTool)
s.registerTool('登录', { description: '会话登录：传入用户名与密码时必须调用本工具', inputSchema: loginSchema }, loginTool)
s.registerTool('令牌', { description: '令牌管理：需要认证；用于基于会话生成新令牌', inputSchema: z.object({}) }, tokenTool)
s.registerTool('token', { description: '令牌管理：需要认证；用于基于会话生成新令牌', inputSchema: z.object({}) }, tokenTool)
s.registerTool('RSS', { description: 'RSS：无入参；需要获取订阅时调用本工具', inputSchema: z.object({}) }, rssTool)
s.registerTool('rss', { description: 'RSS：无入参；需要获取订阅时调用本工具', inputSchema: z.object({}) }, rssTool)

const t = new StdioServerTransport()
await s.connect(t)

const httpPort = Number(process.env.NOTE_HTTP_PORT || 0)
if (httpPort) {
  const tools = new Map([
    ['search', searchTool], ['搜索', searchTool],
    ['publish', publishTool], ['发布', publishTool],
    ['delete', deleteTool], ['删除', deleteTool],
    ['update', updateTool], ['更新', updateTool],
    ['message', getTool], ['消息', getTool],
    ['page', pageTool], ['页面', pageTool],
    ['pin', pinTool], ['置顶消息', pinTool],
    ['settings', settingsTool], ['设置', settingsTool],
    ['status', statusTool], ['状态', statusTool],
    ['calendar', calendarTool], ['日历', calendarTool],
    ['config', configTool], ['配置', configTool],
    ['login', loginTool], ['令牌', tokenTool], ['token', tokenTool], ['rss', rssTool]
  ])

  const clients = new Set()

  const srv = http.createServer(async (req, res) => {
    const u = parse(req.url || '', true)
    const p = u.pathname || ''
    const cors = {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization'
    }
    if (req.method === 'OPTIONS') {
      res.writeHead(204, cors)
      res.end()
      return
    }
    if (req.method === 'GET' && p === '/mcp/tools') {
      res.writeHead(200, { 'Content-Type': 'application/json', ...cors })
      res.end(JSON.stringify(Array.from(tools.keys())))
      return
    }
    if (req.method === 'GET' && p === '/mcp/sse') {
      res.writeHead(200, {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        Connection: 'keep-alive',
        ...cors
      })
      res.write('retry: 5000\n\n')
      res.write(`event: mcp_hello\ndata: ${JSON.stringify({ name: 'ech0-noise-mcp', version: '0.1.0' })}\n\n`)
      res.write(`event: mcp_tools\ndata: ${JSON.stringify(Array.from(tools.keys()))}\n\n`)
      const onStart = (e) => res.write(`event: tool_start\ndata: ${JSON.stringify(e)}\n\n`)
      const onEnd = (e) => res.write(`event: tool_end\ndata: ${JSON.stringify(e)}\n\n`)
      const iv = setInterval(() => { try { res.write('event: keepalive\ndata: {}\n\n') } catch {} }, 30000)
      bus.on('tool_start', onStart)
      bus.on('tool_end', onEnd)
      clients.add(res)
      req.on('close', () => {
        bus.off('tool_start', onStart)
        bus.off('tool_end', onEnd)
        clients.delete(res)
        clearInterval(iv)
        res.end()
      })
      return
    }
    if (req.method === 'GET' && p.startsWith('/mcp/sse/tool/')) {
      const name = decodeURIComponent(p.replace('/mcp/sse/tool/', ''))
      const fn = tools.get(name)
      if (!fn) {
        res.writeHead(404, { 'Content-Type': 'text/event-stream', ...cors })
        res.write('retry: 5000\n\n')
        res.write(`event: error\ndata: ${JSON.stringify({ error: 'tool_not_found' })}\n\n`)
        res.end()
        return
      }
      res.writeHead(200, {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        Connection: 'keep-alive',
        ...cors
      })
      res.write('retry: 5000\n\n')
      let args = {}
      try {
        const q = u && u.query ? u.query : {}
        if (q && typeof q.input === 'string' && q.input.length) {
          try { args = JSON.parse(q.input) } catch {}
        } else if (q && typeof q.args === 'string' && q.args.length) {
          try { args = JSON.parse(q.args) } catch {}
        }
      } catch {}
      bus.emit('tool_start', { name, args })
      res.write(`event: tool_start\ndata: ${JSON.stringify({ name, args })}\n\n`)
      ;(async () => {
        try {
          const out = await fn(args)
          const txt = (() => {
            const c = out && out.content
            if (Array.isArray(c)) {
              const arr = c.filter((x) => x && x.type === 'text' && typeof x.text === 'string').map((x) => x.text)
              return arr.join('\n')
            }
            return typeof out === 'string' ? out : JSON.stringify(out)
          })()
          res.write(`event: tool_data\ndata: ${JSON.stringify({ text: txt })}\n\n`)
          bus.emit('tool_end', { name, res: out })
          res.write(`event: tool_end\ndata: ${JSON.stringify({ name, res: out })}\n\n`)
          res.end()
        } catch (e) {
          const msg = String(e && e.message || e)
          res.write(`event: error\ndata: ${JSON.stringify({ error: msg })}\n\n`)
          res.end()
        }
      })()
      return
    }
    if (req.method === 'POST' && p.startsWith('/mcp/tool/')) {
      const name = decodeURIComponent(p.replace('/mcp/tool/', ''))
      const fn = tools.get(name)
      if (!fn) {
        res.writeHead(404, { 'Content-Type': 'application/json', ...cors })
        res.end(JSON.stringify({ error: 'tool_not_found' }))
        return
      }
      let body = ''
      req.on('data', (chunk) => { body += chunk })
      req.on('end', async () => {
        try {
          const args = body ? JSON.parse(body) : {}
          bus.emit('tool_start', { name, args })
          const out = await fn(args)
          const txt = (() => {
            const c = out && out.content
            if (Array.isArray(c)) {
              const arr = c.filter((x) => x && x.type === 'text' && typeof x.text === 'string').map((x) => x.text)
              return arr.join('\n')
            }
            return typeof out === 'string' ? out : JSON.stringify(out)
          })()
          res.writeHead(200, { 'Content-Type': 'text/plain; charset=utf-8', ...cors })
          res.end(txt)
          bus.emit('tool_end', { name, res: out })
        } catch (e) {
          res.writeHead(200, { 'Content-Type': 'text/plain; charset=utf-8', ...cors })
          res.end(String(e && e.message || e))
        }
      })
      return
    }
    res.writeHead(404, cors)
    res.end()
  })
  srv.listen(httpPort)
}
