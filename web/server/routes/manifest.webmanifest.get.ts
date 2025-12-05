import { send, setHeader } from 'h3'

export default defineEventHandler((event) => {
  setHeader(event, 'Content-Type', 'application/manifest+json; charset=utf-8')
  const manifest = {
    name: '说说笔记',
    short_name: '说说',
    start_url: '/',
    display: 'standalone',
    background_color: '#111111',
    theme_color: '#111111',
    icons: [
      { src: '/favicon.svg', sizes: 'any', type: 'image/svg+xml' },
      { src: '/favicon.svg', sizes: 'any', type: 'image/svg+xml', purpose: 'maskable any' }
    ]
  }
  return send(event, JSON.stringify(manifest))
})
