import { useHead, useRuntimeConfig, useSeoMeta } from 'nuxt/app'

interface PageSeo {
  title: string
  description: string
  // The page's path, for its canonical link.
  path: string
  noindex?: boolean
}

// Title, description and canonical link for a page. The site name is added
// by the title template in nuxt.config.ts.
export function usePageSeo({ title, description, path, noindex = false }: PageSeo) {
  const siteUrl = String(useRuntimeConfig().public.siteUrl || 'https://funkeygrabandbite.com').replace(/\/$/, '')
  useSeoMeta({
    title,
    description,
    ogTitle: `${title} | Funkey Grab & Bite`,
    ogDescription: description,
    ogUrl: `${siteUrl}${path}`,
    robots: noindex ? 'noindex, nofollow' : undefined,
  })
  useHead({ link: [{ rel: 'canonical', href: `${siteUrl}${path}`, key: 'canonical' }] })
}
