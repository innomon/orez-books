# Gemini: Frappe Books Re-architecture (Wails)

This document contains foundational mandates and engineering standards for this project. These instructions take absolute precedence over general defaults.

## 1. Core Mandates
- **Go Backend First:** All business logic, data processing, and system interactions must be moved to the Go backend.
- **Frontend Purity:** The frontend (Vue 3) should focus solely on UI and user interaction, calling the Go backend for all operations.
- **SQLite Integrity:** Ensure the SQLite database remains the source of truth, with robust transaction handling in the Go layer.
- **No Node.js in Production:** The final build must not depend on a Node.js runtime.

## 2. Engineering Standards

### 2.1 Backend (Go)
- **Typing:** Use strict, static typing for all Go methods and data structures.
- **Architecture:** Follow a clean architecture pattern (e.g., separating database models, services, and Wails handlers).
- **SQLite Driver:** Prefer `modernc.org/sqlite` for a Cgo-free build experience if possible, or `github.com/mattn/go-sqlite3` for performance where Cgo is acceptable.
- **Logging:** Use structured logging (e.g., `slog` or `zap`). Do not log sensitive financial or personal data (PII).
- **Concurrency:** Leverage Go routines for background tasks (e.g., data synchronization or heavy reports).

### 2.2 Frontend (Vue 3)
- **Composition API:** Use the Vue 3 Composition API with `<script setup>` for all new or ported components.
- **TypeScript:** Use TypeScript for all frontend logic to ensure type safety between Go and JS.
- **Wails Bindings:** Always use generated Wails bindings (found in `frontend/wailsjs/go`) to interact with the backend.

### 2.3 Migration Protocol
- **Surgical Porting:** When porting a feature from the original repository, first understand the Electron/Node.js implementation, then implement the equivalent Go version.
- **Verification:** Every ported feature must be verified against its original behavior in Frappe Books.
- **Data Safety:** Before performing any destructive operations on an SQLite file, ensure a backup mechanism is in place.

## 3. Tool Usage
- **Wails CLI:** Use `wails dev` for local development and `wails build` for production builds.
- **Go Mod:** Use Go modules for dependency management.
- **Vite:** Use the Vite development server provided by Wails for frontend development.

## 4. Security
- **Input Validation:** All data coming from the frontend (via bindings) must be validated in the Go backend before being used in SQLite queries or system commands.
- **Credential Protection:** Never hardcode secrets or API keys. Use environment variables or secure local storage (e.g., OS keychain).
