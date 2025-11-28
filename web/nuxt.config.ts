export default defineNuxtConfig({
  app: {
    head: {
      link: [
        { rel: 'stylesheet', href: 'https://www.noisework.cn/css/APlayer.min.css' },
        { rel: 'stylesheet', href: 'https://api.hypcvgm.top/NeteaseMiniPlayer/netease-mini-player-v2.css' },
        { rel: 'stylesheet', href: 'https://unpkg.com/@waline/client@v3/dist/waline.css' },
        { rel: 'stylesheet', href: 'https://cdn.jsdelivr.net/npm/@fancyapps/ui@5.0/dist/fancybox/fancybox.css' },
        { rel: 'icon', href: '/favicon.ico' },
        { rel: 'manifest', href: '/manifest.webmanifest' },
      ],
      script: [
        { src: 'https://cdn.jsdelivr.net/npm/jquery@3.2.1/dist/jquery.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/aplayer@1.10.1/dist/APlayer.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/meting@2.0.1/dist/Meting.min.js', body: true },
        { src: 'https://api.hypcvgm.top/NeteaseMiniPlayer/netease-mini-player-v2.js', body: true, defer: true },
        { 
          src: 'https://unpkg.com/@waline/client@v3/dist/waline.js',
          body: true,
          defer: true
        },
        { 
          src: 'https://cdn.jsdelivr.net/npm/@fancyapps/ui@5.0/dist/fancybox/fancybox.umd.js',
          body: true 
        },
        { src: 'https://unpkg.com/medium-zoom/dist/medium-zoom.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/html2canvas@1.4.1/dist/html2canvas.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/qrcode/build/qrcode.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/bcryptjs@2.4.3/dist/bcrypt.min.js', body: true },
      ],
      meta: [
        { name: "viewport", content: "width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" }
      ],
      title: '说说笔记'
    }
  },
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  plugins: [
    '~/plugins/fetch.ts'
  ],
  modules: [
    '@nuxt/ui',
    '@nuxt/fonts',
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
  ],
  css: [
    '@/assets/fonts/result.css',
  ],
  colorMode: {
    preference: 'light'
  },
  runtimeConfig: {
    public: {
      baseApi: process.env.BASE_API || '/api',
    }
  },
  // 添加以下配置
  nitro: {
    preset: 'node-server',
    devProxy: {
      '/api': {
        target: 'http://localhost:1315/api',
        changeOrigin: true
      },
      '/rss': {
        target: 'http://localhost:1315/rss',
        changeOrigin: true
      }
    },
    routeRules: {
      '/**': {
        headers: {
          'Content-Security-Policy': "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https:; style-src 'self' 'unsafe-inline' https:; img-src 'self' data: https: http:; font-src 'self' https:; connect-src 'self' http: https:; frame-src 'self' https:;"
        }
      }
    }
  },
  build: {
    transpile: ['@heroicons/vue'],
  },
  experimental: {
    payloadExtraction: false
  },
  devServer: {
    port: 3000,
    host: '0.0.0.0'
  },
  ssr: false,
})
