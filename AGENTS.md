## Project Overview

**cliQ** is a cross-platform desktop application built with Wails that transforms complex CLI commands into intuitive graphical user interfaces. Users define command templates with variables in `.cliqfile.yaml` files, and cliQ automatically generates dynamic forms for parameter input.

## Architecture

This is an **Nx monorepo** with a hybrid architecture:
- **Wails v2**: Cross-platform desktop framework (Go backend + Vue frontend)
- **Go 1.24**: Backend for business logic, file operations, and command execution
- **Vue 3 + TypeScript**: Frontend providing the user interface
- **pnpm**: Package manager with workspace support
- **Task**: Task runner for build automation

### Monorepo Structure

```
├── apps/
│   ├── cliq-app/           # Main Wails desktop application
│   │   ├── frontend/       # Vue 3 frontend (PrimeVue UI)
│   │   ├── main.go         # Wails entry point
│   │   └── wails.json      # Wails configuration
│   ├── cliq-hub-backend/   # Go backend for template marketplace
│   ├── cliq-hub-frontend/  # Vue frontend for marketplace
│   └── cliq-doc/           # Documentation site (Astro)
├── packages/
│   ├── cliqfile/          # Go library for parsing .cliqfile.yaml
│   └── shared-vue-ui/     # Shared Vue components
└── doc/                   # Documentation and examples
```

## Development Commands

### Prerequisites
- Node.js >= 22
- pnpm 10
- Go 1.24
- Wails CLI

### Setup
```bash
pnpm install  # Install all dependencies
```

### Development
```bash
# Main desktop application (most common)
task dev:cliq
# or
pnpm nx run cliq-app:serve
cd apps/cliq-app && wails dev

# Backend service
task dev:cliq-hub-backend
# or
pnpm nx run cliq-hub-backend:serve

# Web frontend
pnpm nx run cliq-hub-frontend:serve
```

### Building
```bash
# Build desktop app
pnpm nx run cliq-app:build

# Build individual components
pnpm nx run cliq-hub-backend:build
pnpm nx run cliq-hub-frontend:build
```

### Testing
```bash
# Test Go cliqfile library
pnpm nx run cliqfile:test
# or
go test ./packages/cliqfile/...
```

## Key Technologies

### Frontend Stack
- **Vue 3 Composition API** with TypeScript
- **PrimeVue V4**: UI component library with auto-import
- **Vite**: Build tool with HMR
- **TailwindCSS v4**: Utility-first CSS
- **Monaco Editor**: Code editing capabilities

### Backend Stack
- **Wails IPC**: Go-Vue communication
- **Go Workspaces**: Multi-module Go development
- **YAML parsing**: gopkg.in/yaml.v3 for template parsing

## Template System

The core of cliQ is the `.cliqfile.yaml` template format:

```yaml
name: "FFmpeg Video Converter"
description: "Convert videos using FFmpeg"
version: "1.0"
author: "cliQ Team"
cliq_template_version: "1.0"

cmds:
  - name: "convert"
    description: "Convert video format"
    command: "ffmpeg -i {{input_file}} {{output_file}}"
    variables:
      - name: "input_file"
        type: "file_input"
        label: "Input Video"
        required: true
      - name: "output_file"
        type: "file_output"
        label: "Output Video"
        required: true
```

### Variable Types
- `string`: Text input
- `file_input`: File picker for input files
- `file_output`: File picker for output files
- `boolean`: Checkbox
- `number`: Number input
- `select`: Dropdown selection

## Architecture Patterns

### Wails Integration
- **Go Backend** (`apps/cliq-app/app.go`): Core business logic, file operations, command execution
- **Vue Frontend** (`apps/cliq-app/frontend/`): UI and form generation
- **IPC Communication**: Wails JS bindings for seamless Go-Vue communication

### Template Processing
1. **cliqfile Package**: Go library parses YAML templates into structured models
2. **Dynamic Form Generation**: Vue components auto-generate forms from template definitions
3. **Command Execution**: Go backend substitutes variables and executes commands

### Shared Components
- **DynamicCommandForm**: Auto-generates forms from template definitions
- **Shared UI**: Reusable Vue components in `packages/shared-vue-ui`

## Configuration Files

- **`wails.json`**: Wails app configuration, window settings, protocols
- **`go.work`**: Go workspace module definitions
- **`vite.config.ts`**: Frontend build with path aliases (`@/` for src)
- **`Taskfile.yml`**: Build automation and development tasks

## Development Workflow

1. **Feature Development**: Work primarily in `apps/cliq-app/` for desktop features
2. **Template Logic**: Update `packages/cliqfile/` for parsing/validation changes
3. **UI Components**: Add shared components to `packages/shared-vue-ui/`
4. **Testing**: Go unit tests in `packages/cliqfile/*_test.go`

