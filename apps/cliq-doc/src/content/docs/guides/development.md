---
title: Development Guide (Chinese)
description: Monorepo development guide (Chinese)
---

## 依赖环境
- Node：`>=22`
- pnpm：`10`
- Go：`1.24`
- Nx：`19.8.3`（本地）
- Wails CLI（用于 `cliq-app`）：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## 工作区结构
- 应用
  - `apps/cliq-app`：Wails 混合应用（Go + Vue）
  - `apps/cliq-hub-backend`：纯 Go 后端服务
  - `apps/cliq-frontend`：纯 Vue Web 前端
- 包
  - `packages/shared-go-lib`：Go 共享库（模板定义/解析/校验、YAML 工具）
  - `packages/shared-vue-ui`：共享 Vue 组件（目前包含 `DynamicCommandForm`）
- 根文件
  - `pnpm-workspace.yaml`、`package.json`、`nx.json`、`go.work`

## 初始化与安装
- 安装依赖
  - `pnpm install`
- 验证 Nx 可用
  - `pnpm nx --version`
  - `pnpm nx graph --file=project-graph.html`
- 安装 Wails CLI（用于 `cliq-app` 开发/构建）
  - `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Go 工作区（`go.work` 已配置本地模块绑定）
  - 查看绑定：`go work edit -print`

## 启动开发服务器
- Wails 应用（`cliq-app`）
  - `pnpm nx run cliq-app:serve`
  - 说明：调用 `wails dev`，前端行为由 `apps/cliq-app/wails.json` 驱动
- 纯 Go 后端（`cliq-hub-backend`）
  - `pnpm nx run cliq-hub-backend:serve`
  - 等价：`go run ./apps/cliq-hub-backend/cmd/server`
- 纯 Web 前端（`cliq-frontend`）
  - `pnpm nx run cliq-frontend:serve`

## 项目构建
- Wails 应用（前端 + 桌面构建）
  - `pnpm nx run cliq-app:build`
  - 说明：先执行 `frontend-build`（Vite 构建），再执行 `wails build`
- 纯 Go 后端
  - `pnpm nx run cliq-hub-backend:build`
  - 等价：`go build ./apps/cliq-hub-backend/cmd/server`
- 纯 Web 前端
  - `pnpm nx run cliq-frontend:build`

## 测试与校验
- Go 共享库测试
  - `pnpm nx run shared-go-lib:test`
  - 等价：`go test ./packages/shared-go-lib/...`

## 常见问题
- `nx: command not found`
  - 运行 `pnpm install` 安装工作区依赖；通过 `pnpm nx --version` 验证。
- `@nx-go/nx-go` 版本解析失败
  - 使用已验证版本：`3.3.1`（根 `package.json` 已配置）。
- Wails 构建失败
  - 确认本机已安装 Wails CLI 与系统依赖；可先在本地运行 `wails doctor` 排查。
- 共享 UI 导入失败
  - 已在 `apps/cliq-app/frontend/vite.config.ts` 配置 `@repo/shared-vue-ui` 的别名；确保 `packages/shared-vue-ui/src` 存在并导出组件。
