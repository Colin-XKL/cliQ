---
title: Development Environment and Build Guide
description: Monorepo development environment setup and build guide
---

## Dependencies
- Node：`>=22`
- pnpm：`10`
- Go：`1.24`
- Nx：`19.8.3`（local）
- Wails CLI（for `cliq-app`）：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Workspace Structure
- Applications
  - `apps/cliq-app`：Wails hybrid application（Go + Vue）
  - `apps/cliq-hub-backend`：Pure Go backend service
  - `apps/cliq-frontend`：Pure Vue web frontend
- Packages
  - `packages/shared-go-lib`：Go shared library（template definition/parsing/validation, YAML tools）
  - `packages/shared-vue-ui`：Shared Vue components（currently includes `DynamicCommandForm`）
- Root files
  - `pnpm-workspace.yaml`、`package.json`、`nx.json`、`go.work`

## Initialization and Installation
- Install dependencies
  - `pnpm install`
- Verify Nx availability
  - `pnpm nx --version`
  - `pnpm nx graph --file=project-graph.html`
- Install Wails CLI（for `cliq-app` development/build）
  - `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Go workspace（`go.work` already configured with local module binding）
  - View bindings：`go work edit -print`

## Starting Development Servers
- Wails application（`cliq-app`）
  - `pnpm nx run cliq-app:serve`
  - Note：calls `wails dev`，frontend behavior driven by `apps/cliq-app/wails.json`
- Pure Go backend（`cliq-hub-backend`）
  - `pnpm nx run cliq-hub-backend:serve`
  - Equivalent：`go run ./apps/cliq-hub-backend/cmd/server`
- Pure web frontend（`cliq-frontend`）
  - `pnpm nx run cliq-frontend:serve`

## Project Building
- Wails application（frontend + desktop build）
  - `pnpm nx run cliq-app:build`
  - Note：first execute `frontend-build`（Vite build），then execute `wails build`
- Pure Go backend
  - `pnpm nx run cliq-hub-backend:build`
  - Equivalent：`go build ./apps/cliq-hub-backend/cmd/server`
- Pure web frontend
  - `pnpm nx run cliq-frontend:build`

## Testing and Validation
- Go shared library tests
  - `pnpm nx run shared-go-lib:test`
  - Equivalent：`go test ./packages/shared-go-lib/...`

## Common Issues
- `nx: command not found`
  - Run `pnpm install` to install workspace dependencies；verify with `pnpm nx --version`。
- `@nx-go/nx-go` version resolution failure
  - Use verified version：`3.3.1`（root `package.json` already configured）。
- Wails build failure
  - Confirm Wails CLI and system dependencies are installed locally；run `wails doctor` locally for troubleshooting。
- Shared UI import failure
  - Alias for `@repo/shared-vue-ui` already configured in `apps/cliq-app/frontend/vite.config.ts`；ensure `packages/shared-vue-ui/src` exists and exports components。