# Agents: Frappe Books Re-architecture (Wails)

This document defines the specialized agent roles for the re-architecture of Frappe Books from Electron to Wails.

## 1. Go Backend Specialist
- **Expertise:** Go (Golang), Wails Framework, SQLite (Go-side), GORM, System Programming.
- **Responsibilities:** 
  - Designing the Go backend architecture.
  - Implementing SQLite data access layer in Go.
  - Porting Node.js business logic to Go.
  - Creating and maintaining Wails Bindings.
- **Goal:** Ensure a robust, performant, and type-safe Go backend.

## 2. Vue.js Frontend Specialist
- **Expertise:** Vue 3, Vite, Tailwind CSS, TypeScript, Wails Bindings Integration.
- **Responsibilities:**
  - Adapting the existing Frappe Books frontend to work within Wails.
  - Replacing `ipcRenderer` calls with Wails JS bindings.
  - Optimizing the frontend build for the Wails runtime.
  - Managing application state and UI logic.
- **Goal:** Maintain the visual and functional integrity of the UI while bridging it to the Go backend.

## 3. Database Migration Specialist
- **Expertise:** SQLite, SQL, Knex.js (reading), Go GORM (writing), Data Integrity.
- **Responsibilities:**
  - Analyzing existing SQLite schemas and Knex migrations.
  - Designing a compatible data layer in Go.
  - Ensuring seamless data migration for existing users.
  - Implementing robust data validation in the Go layer.
- **Goal:** Guarantee zero data loss and full compatibility with existing Frappe Books databases.

## 4. Native OS Integration Specialist
- **Expertise:** Wails Runtime APIs, OS-specific behavior (macOS, Windows, Linux), Menu/Tray APIs.
- **Responsibilities:**
  - Implementing native system menus and tray icons.
  - Handling file system dialogs and shell interactions.
  - Managing application window states.
  - Configuring cross-platform build artifacts.
- **Goal:** Provide a seamless, native-feeling experience across all desktop platforms.

## 5. Security & Compliance Agent
- **Expertise:** Secure coding practices in Go and JS, Privacy principles, PII handling.
- **Responsibilities:**
  - Auditing Go-JS communication for potential vulnerabilities.
  - Ensuring secure SQLite data handling.
  - Verifying that no sensitive data (PII) is exposed in logs or IPC.
  - Enforcing best practices for desktop application security.
- **Goal:** Maintain the security and privacy standards of Frappe Books throughout the re-architecture.
