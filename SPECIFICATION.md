# Specification: Frappe Books Re-architecture (Electron to Wails)

## 1. Overview
The goal of this project is to re-architect **Frappe Books** by replacing the **Electron** desktop framework with **Wails**. This transition shifts the backend from **Node.js** to **Go**, while maintaining the existing **Vue.js 3** frontend.

## 2. Current Architecture (Electron)
- **Frontend:** Vue.js 3, Vite, Tailwind CSS.
- **Backend/Main Process:** Node.js (v16+).
- **Communication:** Electron IPC (`ipcMain`, `ipcRenderer`).
- **Database:** SQLite via `better-sqlite3` (Node native module) and `Knex.js`.
- **System Access:** Node.js standard libraries and Electron-specific APIs (e.g., `electron-store`, `dialog`, `menu`).
- **Packaging:** `electron-builder`.

## 3. Target Architecture (Wails)
- **Frontend:** Vue.js 3, Vite, Tailwind CSS (mostly unchanged).
- **Backend:** Go (v1.21+).
- **Communication:** Wails Bindings (Go methods exposed to JS) and Wails Events.
- **Database:** SQLite via Go drivers (e.g., `github.com/mattn/go-sqlite3` or `modernc.org/sqlite`) and a Go ORM/Query Builder (e.g., `GORM`).
- **System Access:** Go standard library and Wails Runtime APIs.
- **Packaging:** Wails CLI (produces native binaries).

## 4. Key Components for Migration

### 4.1 Backend Transition
- **Main Entry Point:** Replace `src/main/index.js` (or equivalent) with a Go `main.go` and `app.go`.
- **Logic Porting:** All business logic currently residing in the Electron main process must be ported to Go.
- **SQLite Integration:** Replace `better-sqlite3` with a Go-native SQLite implementation. This is critical as it eliminates the need for Node native module rebuilding.

### 4.2 Database Layer (Complex Porting)
- **Custom ORM:** Port the logic from `DatabaseCore.ts` (currently using Knex) to a Go-equivalent. This includes:
  - Dynamic table creation and schema migrations.
  - Handling of "SingleValue" tables for settings.
  - CRUD operations with nested "Table" (child) field support.
  - PRAGMA foreign_key management.
- **Bespoke Queries:** Port custom SQL logic from `bespoke.ts`.
- **Patches:** Implement the `runPatch.ts` logic in Go to handle data migrations between versions.

### 4.3 Native Features Migration
- **IPC Actions:** Map `ipcMain.handle` actions to Go `App` methods (Bindings).
- **IPC Messages:** Map `ipcMain.on` messages to Wails `Runtime` events or specific Go handlers.
- **Dialogs:** Port `dialog.showOpenDialog`, `dialog.showSaveDialog`, and `dialog.showMessageBox` to Wails Runtime Dialogs.
- **Window Management:** Port window lifecycle events (minimize, maximize, close) to Wails Runtime.
- **External Links:** Replace `shell.openExternal` with Go's `browser.OpenURL` or equivalent.

## 5. Performance and Resource Goals
- **Binary Size:** Reduce the distribution size significantly (Wails binaries are typically < 20MB vs. Electron's 100MB+).
- **Memory Usage:** Improve runtime memory efficiency by replacing the V8 backend with a Go binary.
- **Cold Start:** Faster application launch times.

## 6. Technical Constraints
- Must maintain offline-first capability.
- SQLite database schema must remain compatible or handled via robust migrations.
- Frontend UI should remain visually identical to the original app.
