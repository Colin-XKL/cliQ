---
title: Syntax Consistency
description: CLIQfile Syntax Consistency Guide
---

# CLIQfile Syntax Consistency Guide

This document outlines the key directories and files that must be maintained consistently to ensure proper CLIQfile syntax handling across the application.

## Core Model Definition

**File:** `/cliq/models/models.go`

**Purpose:** Defines the core data structures for CLIQfile parsing, including the `VariableDefinition` struct and type constants.

**Key Elements:**
- `VarTypeText`, `VarTypeFileInput`, `VarTypeFileOutput`, `VarTypeBoolean`, `VarTypeNumber`, `VarTypeSelect`
- All type constants must use consistent values across the system

## Backend Implementation

### Template Service
**File:** `/cliq/services/template_service.go`

**Purpose:** Handles parsing, validation, and generation of CLIQfile templates.

**Key Elements:**
- `determineVariableType()` function uses model constants
- `validateTemplate()` function validates against supported types
- Ensures type consistency between parsing and validation

### Validation Service
**File:** `/cliq-hub-backend/internal/validation/template.go`

**Purpose:** Validates CLIQfile structure and content on the backend.

**Key Elements:**
- `allowedTypes` map must include all supported types
- Type names must match model constants

### LLM Client
**File:** `/cliq-hub-backend/internal/llm/client.go`

**Purpose:** Generates CLIQfiles using LLM, includes syntax documentation in prompts.

**Key Elements:**
- Embedded syntax documentation (`cliqfile_syntax.md`) guides LLM generation
- Ensures generated templates follow consistent syntax

## Frontend Implementation

### Dynamic Form Component
**File:** `/cliq/frontend/src/components/DynamicCommandForm.vue`

**Purpose:** Renders UI components based on the variable type from CLIQfiles.

**Key Elements:**
- Type checking conditions (`variable.type === 'string'`, etc.)
- Must use same type names as backend models

### Template Editor
**File:** `/cliq/frontend/src/components/TemplateEditorModal.vue`

**Purpose:** Provides template editing and preview functionality.

## Documentation

### Syntax Reference
**File:** `/doc/cliqfile_syntax.md`

**Purpose:** User-facing documentation of CLIQfile syntax.

**Key Elements:**
- Must match actual supported types in implementation
- Used as reference for LLM generation

### Backend Syntax Documentation
**File:** `/cliq-hub-backend/asset/cliqfile_syntax.md`

**Purpose:** Embedded in LLM prompts to guide template generation.

**Note:** Same content as main documentation, embedded in backend binary.

## Key Consistency Rules

1. **Type Names:** All components must use the same type names (`string`, `file_input`, `file_output`, `number`, `boolean`, `select`)

2. **Backward Compatibility:** Changes to syntax must maintain compatibility with existing templates

3. **Validation First:** Add new types to validation before implementing in other components

4. **Frontend Update:** When adding new types, update the DynamicCommandForm to handle the new type

## Development Workflow

When adding or modifying CLIQfile syntax:

1. Update model constants in `models.go`
2. Add to validation logic in backend
3. Update frontend type handling
4. Update all documentation files
5. Test end-to-end
