export default defineNuxtRouteMiddleware(async (to) => {
  const baseApi = useRuntimeConfig().public.baseApi || '/api'
  try {
    const res = await $fetch(`${baseApi}/user`, { credentials: 'include' })
    if ((res as any)?.code === 1) return
  } catch {}
  return navigateTo({ path: '/auth/login', query: { redirect: to.fullPath } })
})
