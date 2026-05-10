# Funkey Grab & Bite - **Unified Food Commerce Platform**

![Go](https://img.shields.io/badge/Backend-Go%201.24-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/API-Gin-008ECF)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-336791?logo=postgresql&logoColor=white)
![React](https://img.shields.io/badge/Admin-React%20%2B%20Vite-61DAFB?logo=react&logoColor=black)
![Nuxt](https://img.shields.io/badge/Web-Nuxt%204-00DC82?logo=nuxtdotjs&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/UI-Tailwind%20CSS-06B6D4?logo=tailwindcss&logoColor=white)
![TypeScript](https://img.shields.io/badge/Frontend-TypeScript-3178C6?logo=typescript&logoColor=white)
![AWS](https://img.shields.io/badge/Cloud-AWS-232F3E?logo=amazon-aws&logoColor=white)
![Terraform](https://img.shields.io/badge/IaC-Terraform-844FBA?logo=terraform&logoColor=white)

Funkey Grab & Bite is a full-stack digital operations platform for a modern fast-food and catering business. It unifies direct customer ordering, promotion-driven conversion, admin-level operations control, and operational data visibility into one system. The platform is designed to support both day-to-day service excellence and long-term scalability.

For contributors, this repository offers a clear modular architecture (API, public web, admin dashboard) with strongly typed frontend workflows and domain-oriented backend services. For investors and business stakeholders, it demonstrates a practical, revenue-oriented product surface: conversion funnel, retention mechanics (promotions), fulfillment visibility (orders/catering), and operational discipline (inventory, reporting, admin controls).

## 📚 Table Of Contents
- [Product Vision](#-product-vision)
- [Business Value](#-business-value)
- [System Architecture](#-system-architecture)
- [AWS Infrastructure Design](#-aws-infrastructure-design)
- [Terraform IaC Structure](#-terraform-iac-structure)
- [Repository Structure](#-repository-structure)
- [Core Platform Modules](#-core-platform-modules)
- [API Surface](#-api-surface)
- [Real-Time Event Contract](#-real-time-event-contract)
- [Data Model Snapshot](#-data-model-snapshot)
- [Security And Operational Controls](#-security-and-operational-controls)
- [Quick Start](#-quick-start)
- [Environment Variables](#-environment-variables)
- [Runbook Commands](#-runbook-commands)
- [Contribution Guide](#-contribution-guide)
- [Roadmap Signals](#-roadmap-signals)
- [Acknowledgments](#-acknowledgments)
- [Contact](#-contact)

## 🚀 Product Vision
Build a restaurant operations engine where:
- Customers can discover menu items, place orders quickly, and request catering confidently.
- Promotions become programmable growth levers, not ad-hoc discounts.
- Operations teams can manage demand, inventory, and customer service from one command center.
- Leadership has visibility into daily sales performance and operational health.

## 💼 Business Value
### Why this product is commercially compelling
- Multi-channel revenue capture:
Direct order flow + event catering + recurring promotions.
- Operational efficiency:
Centralized admin workflows reduce coordination friction between kitchen, fulfillment, and support.
- Customer retention:
Promotion validation and personalized order history encourage repeat transactions.
- Decision velocity:
Dashboard and reports support fast commercial and operational decisions.

### Investor-relevant product characteristics
- Category fit:
Digitized food commerce with integrated operations control.
- Expansion readiness:
Modular backend/frontend boundaries support feature extension and multi-location growth.
- Data leverage potential:
Order, inventory, user, and promotion events form the foundation for predictive analytics.

## 🏗️ System Architecture
```text
						+-----------------------------+
						|  Customers (Web / Mobile)  |
						+--------------+--------------+
									   |
									   | HTTPS
									   v
					  +----------------+----------------+
					  |    funkey-bite-web (Nuxt 4)     |
					  |  Landing, menu, order, catering |
					  +----------------+----------------+
									   |
									   | REST /api/v1
									   v
 +---------------------+   service/repo boundaries   +----------------------+
 | funkey-bite-admin   | <--------------------------> |  funkey-bite-api     |
 | React + Vite        |     auth, ops, reporting     |  Go + Gin + JWT      |
 | Ops dashboard       |                               |  Domain services     |
 +----------+----------+                               +----------+-----------+
			|                                                     |
			| optional realtime (Socket.IO contract)              | SQL
			v                                                     v
	  +-----+----------------------+                     +--------+-----------+
	  | Realtime event source      |                     | PostgreSQL         |
	  | (new_order, alerts, etc.)  |                     | users/orders/menu  |
	  +----------------------------+                     | promotions/inventory|
														 +--------------------+
```

## ☁️ AWS Infrastructure Design
Infrastructure is provisioned with Terraform and follows an edge-first AWS design with isolated service tiers.

- DNS and certificates: Route 53 + ACM.
- Edge protection: optional AWS WAF in front of CloudFront.
- Static delivery: dedicated CloudFront and private S3 origins for web and admin.
- API edge and compute: CloudFront API edge -> ALB -> ECS Fargate API service.
- Async workloads: ECS workers + Lambda + SQS + EventBridge Scheduler.
- Data layer: Aurora PostgreSQL Serverless v2 + ElastiCache Redis.
- Ops visibility: CloudWatch logs/alarms/dashboard + SNS notifications.

High-level infrastructure flow:

```text
Route53 -> WAF -> CloudFront(web/admin/api)
        -> S3 (web/admin) OR ALB -> ECS API
ECS/Lambda -> Aurora + Redis + SQS + Secrets Manager
EventBridge -> scheduled ECS tasks
CloudWatch/SNS -> monitoring and alerting
```

Implemented cloud architecture diagram:

```text
			┌──────────────────────────────────────────────────────┐
			│                   Route 53 (DNS)                    │
			│   funkeygrabandbite.com / api / admin / www         │
			└──────────────┬───────────────────────────────────────┘
				       │
			┌──────────────▼───────────────────────────────────────┐
			│             AWS WAF v2 + Shield Standard             │
			│      CRS · KnownBadInputs · SQLi · Rate-limiting     │
			└──────┬───────────────────────────┬────────────────────┘
			       │                           │
	      ┌────────────────▼──────┐       ┌───────────▼──────────────────┐
	      │  CloudFront (Web SPA) │       │  CloudFront (Admin SPA)      │
	      │  Nuxt storefront      │       │  React admin dashboard        │
	      └────────┬─────────────┘       └───────────┬──────────────────┘
		       │ OAC                              │ OAC
	      ┌────────▼─────────────┐       ┌───────────▼──────────────────┐
	      │  S3 Bucket (Web)     │       │  S3 Bucket (Admin)           │
	      │  Private + versioned │       │  Private + versioned         │
	      └──────────────────────┘       └──────────────────────────────┘

			┌──────────────────────────────────────┐
			│       CloudFront (API Edge)          │
			│ CachingDisabled · AllViewer policy   │
			│ Origin Shield · WebSocket support    │
			└────────────────┬─────────────────────┘
					 │ secret header
			┌────────────────▼─────────────────────┐
			│      Application Load Balancer       │
			│ HTTP->HTTPS redirect + header guard  │
			└────────────────┬─────────────────────┘
					 │
	  ┌──────────────────────────────┴────────────────────────────────┐
	  │                                                               │
┌─────────▼──────────────┐                               ┌───────────────▼───────────┐
│   ECS Fargate (API)    │                               │  ECS Fargate (Workers)    │
│   Go/Gin REST API      │                               │  Background task runners   │
│   Auto scaling         │                               │  FARGATE_SPOT preferred    │
└─────────┬──────────────┘                               └───────────────────────────┘
	  │
┌─────────▼────────────────────────────────────────────────────────────────────────┐
│                               Data and Messaging Layer                          │
│                                                                                 │
│  ┌──────────────────────────┐  ┌────────────────────┐  ┌──────────────────────┐ │
│  │  Aurora PostgreSQL Sv2   │  │  ElastiCache Redis │  │  SQS + DLQs          │ │
│  │  Multi-AZ + reader       │  │  7.1 encryption    │  │  Order + Catering    │ │
│  └──────────────────────────┘  └────────────────────┘  └──────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                            Async and Scheduled Work                              │
│                                                                                  │
│  ┌──────────────────────────────────┐        ┌──────────────────────────────────┐│
│  │  Lambda (Catering Notifier)      │        │ EventBridge Scheduler            ││
│  │  Container image + VPC           │        │ Promo expiry and daily reports   ││
│  │  SQS event source mapping        │        │ ECS task trigger                 ││
│  └──────────────────────────────────┘        └──────────────────────────────────┘│
└───────────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                         Security, Secrets and Monitoring                         │
│  Secrets Manager · IAM roles · CloudWatch logs/alarms/dashboard · SNS alerts    │
└───────────────────────────────────────────────────────────────────────────────────┘
```

Detailed infrastructure operations guide:
[terraform/README.md](terraform/README.md)

## 🧱 Terraform IaC Structure
Terraform implementation lives in [terraform](terraform) and is organized for reusable modules and environment compositions.

- Modules: [terraform/modules](terraform/modules)
Contains reusable units for networking, DNS, WAF, CloudFront, ALB, ECS, Lambda, Aurora, Redis, SQS, EventBridge, Secrets, and Monitoring.
- Environment roots: [terraform/environments/staging](terraform/environments/staging) and [terraform/environments/production](terraform/environments/production)
Each environment toggles components with feature flags, so you can deploy services individually or together.
- Scripts: [terraform/scripts](terraform/scripts)
Contains helper scripts for targeted deploy/plan and safe no-apply checks.

Common workflows:

```bash
cd terraform

# safe checks (no apply)
./scripts/tf-all.sh fmt
./scripts/tf-all.sh validate
./scripts/tf-all.sh plan

# regular workflows
make plan ENV=staging
make deploy-api ENV=production
```

## 🗂️ Repository Structure
```text
funkey-grab-and-bite/
├── README.md                         # Root product + engineering documentation
├── funkey-bite-api/                  # Go backend API (auth, orders, catering, admin)
│   ├── cmd/api/main.go               # Route wiring, middleware, service composition
│   ├── internal/
│   │   ├── handlers/                 # HTTP handlers + middleware
│   │   ├── services/                 # Business logic per domain
│   │   ├── repository/               # Postgres data access layer
│   │   ├── domain/models/            # Domain and transport models
│   │   ├── database/                 # DB initialization + table bootstrap
│   │   └── utils/                    # JWT, password, email, SMS, helpers
│   ├── migrations/                   # SQL migration scripts
│   └── go.mod                        # Backend dependency graph
├── funkey-bite-admin/                # Internal operations dashboard
│   ├── src/pages/                    # Orders, inventory, reports, settings, etc.
│   ├── src/api/                      # Axios config + typed API wrappers
│   ├── src/stores/                   # Zustand state domains
│   ├── src/contexts/                 # Auth/socket/theme providers
│   └── package.json                  # Admin build/runtime scripts
├── terraform/                        # AWS infrastructure as code (Terraform)
│   ├── modules/                      # Reusable infrastructure modules
│   ├── environments/                 # staging and production root compositions
│   ├── scripts/                      # fmt/validate/plan/deploy helper scripts
│   └── README.md                     # Infrastructure setup and operations guide
└── funkey-bite-web/                  # Public customer-facing web app (Nuxt)
	├── pages/                        # Menu, order, catering, promotions, legal, contact
	├── components/                   # Reusable domain UI blocks
	├── stores/                       # Pinia store modules (cart/menu)
	├── utils/                        # API helpers + shared utilities
	└── nuxt.config.ts                # Runtime config, SEO metadata, modules
```

## 🧩 Core Platform Modules
### 1) Public Commerce Experience (`funkey-bite-web`)
- Menu discovery, category browsing, and ordering workflows.
- Catering request intake for indoor/outdoor events.
- Promotion validation and customer conversion surfaces.
- SEO-first configuration and social metadata in Nuxt app head.

### 2) Admin Operations Console (`funkey-bite-admin`)
- Dashboard for daily metrics and operational overview.
- End-to-end order and catering management.
- Inventory alerts and stock control operations.
- Promotion lifecycle management.
- Admin account management, profile/security settings, and reports.

### 3) API And Business Logic (`funkey-bite-api`)
- REST API exposing customer, user, and admin workflows.
- Layered architecture: handlers -> services -> repository.
- JWT-based auth for users and admins.

## 📘 API Contract
The canonical route contract lives in [funkey-bite-api/openapi/openapi.json](funkey-bite-api/openapi/openapi.json).

- Backend maintenance command: `cd funkey-bite-api && make contract`
- Admin validation command: `cd funkey-bite-admin && npm run contract:check`
- Web validation command: `cd funkey-bite-web && npm run contract:check`

Frontend route helpers should only call endpoints declared in that contract.
- Centralized middleware for CORS, rate limiting, and auth checks.

## 🔌 API Surface
Base path: `/api/v1`

### Public / Customer endpoints
| Domain | Method | Endpoint | Purpose |
|---|---|---|---|
| Auth | POST | `/auth/register` | Register customer account |
| Auth | POST | `/auth/login` | Customer login |
| Auth | GET | `/auth/check` | Validate customer session |
| Menu | GET | `/menu/` | List available menu items |
| Menu | GET | `/menu/search` | Search menu |
| Menu | GET | `/menu/featured` | Get featured items |
| Menu | GET | `/menu/tags` | Filter menu by tags |
| Menu | GET | `/menu/:id` | Get menu item details |
| Orders | POST | `/orders/` | Place order |
| Orders | GET | `/orders/track/:orderNumber` | Track order |
| Orders | PATCH | `/orders/:id/cancel` | Cancel order |
| Catering | POST | `/catering/requests` | Submit catering request |
| Promotions | GET | `/promotions/active` | List active promotions |
| Promotions | GET | `/promotions/validate` | Validate promotion code |
| Settings | GET | `/settings` | Public business profile/settings |
| Settings | GET | `/settings/hours` | Opening hours |

### Admin endpoints
| Domain | Method | Endpoint | Purpose |
|---|---|---|---|
| Admin Auth | POST | `/admin/auth/login` | Admin login |
| Admin Auth | POST | `/admin/auth/logout` | Admin logout |
| Admin Auth | PATCH | `/admin/auth/password` | Change admin password |
| Dashboard | GET | `/admin/dashboard/stats` | Aggregate dashboard stats |
| Dashboard | GET | `/admin/dashboard/stats/today` | Today stats |
| Reports | GET | `/admin/reports/sales` | Sales reporting |
| Orders | GET | `/admin/orders` | List/manage orders |
| Orders | PATCH | `/admin/orders/:id/status` | Update order status |
| Users | GET | `/admin/users` | List customers |
| Users | PATCH | `/admin/users/:id/status` | Activate/deactivate customer |
| Menu | POST | `/admin/menu/items` | Create menu item |
| Menu | PUT | `/admin/menu/items/:id` | Update menu item |
| Menu | DELETE | `/admin/menu/items/:id` | Delete menu item |
| Catering | GET | `/admin/catering/requests` | List catering requests |
| Catering | PATCH | `/admin/catering/requests/:id/status` | Update catering status |
| Inventory | GET | `/admin/inventory` | Inventory list |
| Inventory | GET | `/admin/inventory/alerts` | Inventory alerts |
| Inventory | PATCH | `/admin/inventory/stock` | Update stock quantity |
| Promotions | GET | `/admin/promotions` | List promotions |
| Promotions | POST | `/admin/promotions` | Create promotion |
| Promotions | PUT | `/admin/promotions/:id` | Update promotion |
| Promotions | DELETE | `/admin/promotions/:id` | Delete promotion |
| Settings | GET | `/admin/settings` | Admin business settings |
| Settings | PUT | `/admin/settings` | Update business settings |
| Admin Users | GET | `/admin/users/admins` | List admin users |
| Admin Users | POST | `/admin/users/admins` | Create admin user |
| Admin Users | PUT | `/admin/users/admins/:id` | Update admin user |
| Admin Users | DELETE | `/admin/users/admins/:id` | Delete admin user |

## ⚡ Real-Time Event Contract
Admin app subscribes to Socket.IO events in `SocketContext` (URL from `VITE_WS_URL`).

Current consumed events:
- `new_order`
- `order_updated`
- `new_catering_request`
- `catering_request_updated`
- `inventory_alert`
- `inventory_updated`
- `new_customer`
- `system_notification`
- `menu_updated`
- `promotion_updated`

These events drive operational notifications and near-real-time dashboard/store updates.

## 🧱 Data Model Snapshot
Key backend tables include:
- `users`
- `admin_users`
- `orders`
- `order_items`
- `menu_categories`
- `menu_items`
- `catering_requests`
- `promotions`
- `promotion_usage`
- `inventory_items`
- `inventory_history`
- `inventory_alerts`
- `notifications`
- `business_settings`

The API also includes default admin bootstrapping logic to ensure first-run access exists.

## 🔐 Security And Operational Controls
Implemented controls:
- JWT-based authorization for customer and admin routes.
- Admin-only route isolation via admin auth middleware.
- Password hashing via bcrypt utilities.
- Rate limiting for public API group and tracking routes.
- CORS allowlist with localhost and production domain support.
- Request validation using struct validation in handlers.

Operational safeguards:
- Default admin account auto-creation can be configured via environment variables.
- Centralized API error handling in frontend clients.
- Session handling and auth token lifecycle controls in admin UI.

## ⚙️ Quick Start
### Prerequisites
- Go `1.24+`
- Node.js `18+` (recommended `20+`)
- npm `9+`
- PostgreSQL `14+`

### 1) Clone
```bash
git clone https://github.com/<your-org>/funkey-grab-and-bite.git
cd funkey-grab-and-bite
```

### 2) Backend setup (`funkey-bite-api`)
Linux/macOS:
```bash
cd funkey-bite-api
go mod download
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=funkey_grab_bite
export DB_SSLMODE=disable
go run cmd/api/main.go
```

Windows PowerShell:
```powershell
cd funkey-bite-api
go mod download
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_NAME="funkey_grab_bite"
$env:DB_SSLMODE="disable"
go run cmd/api/main.go
```

### 3) Admin setup (`funkey-bite-admin`)
```bash
cd funkey-bite-admin
npm install
npm run dev
```
Default local URL: `http://localhost:5173`

### 4) Public web setup (`funkey-bite-web`)
```bash
cd funkey-bite-web
npm install
npm run dev
```
Default local URL: `http://localhost:3000`

## 🔧 Environment Variables
### Backend (`funkey-bite-api`)
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`
- `DEFAULT_ADMIN_EMAIL` (optional)
- `DEFAULT_ADMIN_USERNAME` (optional)
- `DEFAULT_ADMIN_PASSWORD` (optional)
- `DEFAULT_ADMIN_ROLE` (optional)

### Admin frontend (`funkey-bite-admin`)
- `VITE_API_URL` (default `http://localhost:8080`)
- `VITE_WS_URL` (default `http://localhost:8080`)

### Public web (`funkey-bite-web`)
- `NUXT_PUBLIC_API_URL`
- `NUXT_PUBLIC_S3_URL`
- `NUXT_PUBLIC_SITE_URL`

## 🛠️ Runbook Commands
### Backend
```bash
cd funkey-bite-api
go test ./...
go run cmd/api/main.go
```

### Admin
```bash
cd funkey-bite-admin
npm run dev
npm run build
```

### Web
```bash
cd funkey-bite-web
npm run dev
npm run build
npm run generate
```

## 🤝 Contribution Guide
1. Create a feature branch from `main`.
2. Keep commits scoped and descriptive.
3. Validate impacted app(s) before opening a PR:
Go tests, frontend builds, and key auth/order flows.
4. In PR descriptions, include:
problem, approach, impact, test evidence, and rollback notes.
5. Avoid merging changes that break route contracts between frontend and API.

## 🧭 Roadmap Signals
- Realtime event producer alignment in backend for full socket parity.
- Advanced analytics: cohort retention, promotion ROI, fulfillment SLAs.
- Multi-branch support and per-location inventory partitioning.
- Payment provider integration hardening and finance reconciliation exports.

## 🙏 Acknowledgments
Built with Go, Gin, Nuxt, React, Tailwind, PostgreSQL, and the open-source ecosystem that enables fast product iteration.

## 📬 Contact
For partnership, product demos, or technical collaboration:
- Open an issue in this repository
- Or contact the core team via your project communication channel
