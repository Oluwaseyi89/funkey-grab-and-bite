import { readFileSync } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const rootDir = path.resolve(__dirname, '..')
const specPath = path.join(rootDir, 'funkey-bite-api', 'openapi', 'openapi.json')

const spec = JSON.parse(readFileSync(specPath, 'utf8'))

const routeSets = {
  admin: [
    ['post', '/api/v1/admin/auth/login'],
    ['post', '/api/v1/admin/auth/logout'],
    ['patch', '/api/v1/admin/auth/password'],
    ['get', '/api/v1/admin/users/admins'],
    ['post', '/api/v1/admin/users/admins'],
    ['put', '/api/v1/admin/users/admins/{id}'],
    ['delete', '/api/v1/admin/users/admins/{id}'],
    ['get', '/api/v1/admin/dashboard/stats'],
    ['get', '/api/v1/admin/dashboard/stats/today'],
    ['get', '/api/v1/admin/reports/sales'],
    ['get', '/api/v1/admin/orders'],
    ['get', '/api/v1/admin/orders/{id}'],
    ['patch', '/api/v1/admin/orders/{id}/status'],
    ['get', '/api/v1/admin/users'],
    ['patch', '/api/v1/admin/users/{id}/status'],
    ['get', '/api/v1/admin/menu/items'],
    ['get', '/api/v1/admin/menu/items/{id}'],
    ['post', '/api/v1/admin/menu/items'],
    ['put', '/api/v1/admin/menu/items/{id}'],
    ['delete', '/api/v1/admin/menu/items/{id}'],
    ['get', '/api/v1/menu/categories'],
    ['post', '/api/v1/admin/menu/categories'],
    ['put', '/api/v1/admin/menu/categories/{id}'],
    ['get', '/api/v1/admin/inventory'],
    ['get', '/api/v1/admin/inventory/low-stock'],
    ['patch', '/api/v1/admin/inventory/stock'],
    ['post', '/api/v1/admin/inventory/restock'],
    ['get', '/api/v1/admin/inventory/alerts'],
    ['get', '/api/v1/admin/catering/requests'],
    ['patch', '/api/v1/admin/catering/requests/{id}/status'],
    ['get', '/api/v1/admin/promotions'],
    ['get', '/api/v1/admin/promotions/{id}'],
    ['post', '/api/v1/admin/promotions'],
    ['put', '/api/v1/admin/promotions/{id}'],
    ['delete', '/api/v1/admin/promotions/{id}'],
    ['get', '/api/v1/admin/settings'],
    ['put', '/api/v1/admin/settings'],
    ['get', '/api/v1/admin/realtime/ws']
  ],
  web: [
    ['get', '/health'],
    ['get', '/api/v1/menu/categories'],
    ['get', '/api/v1/menu'],
    ['get', '/api/v1/menu/{id}'],
    ['get', '/api/v1/menu/category'],
    ['post', '/api/v1/orders'],
    ['get', '/api/v1/orders/track/{orderNumber}'],
    ['post', '/api/v1/catering/requests'],
    ['get', '/api/v1/promotions/active']
  ]
}

const target = (process.argv[2] || 'all').toLowerCase()
const hasRoute = (method, routePath) => Boolean(spec.paths?.[routePath]?.[method])

const validateTarget = (name) => {
  const missing = routeSets[name].filter(([method, routePath]) => !hasRoute(method, routePath))
  if (missing.length > 0) {
    console.error(`OpenAPI contract validation failed for ${name}. Missing routes:`)
    for (const [method, routePath] of missing) {
      console.error(`- ${method.toUpperCase()} ${routePath}`)
    }
    return false
  }

  console.log(`OpenAPI contract validation passed for ${name}. ${routeSets[name].length} routes checked.`)
  return true
}

let ok = true
if (target === 'all') {
  ok = validateTarget('admin') && validateTarget('web')
} else if (target in routeSets) {
  ok = validateTarget(target)
} else {
  console.error(`Unknown validation target: ${target}`)
  process.exit(1)
}

if (!ok) {
  process.exit(1)
}